<template>
  <div ref="mcpBar" class="mcp-bar" @click="openPanel" title="MCP 配置">
    <svg :class="['mcp-icon', enabled ? 'mcp-icon-on' : 'mcp-icon-off']" viewBox="0 0 24 24"
         xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
      <path
          d="M12 2a2 2 0 0 1 2 2v3h3a2 2 0 1 1 0 4h-3v3a2 2 0 1 1-4 0V11H7a2 2 0 1 1 0-4h3V4a2 2 0 0 1 2-2zm-6 9a6 6 0 0 0 6 6 6 6 0 0 0 6-6 6 6 0 0 0-6-6 6 6 0 0 0-6 6z"
          fill="currentColor"/>
    </svg>
    <span :class="['mcp-dot', enabled ? 'mcp-dot-on' : 'mcp-dot-off']"></span>
    <span :class="['mcp-label', enabled ? 'mcp-label-on' : 'mcp-label-off']">MCP</span>
  </div>

  <teleport to="body">
    <transition name="mcp-overlay-fade">
      <div v-if="panelOpen" class="mcp-overlay" @click.self="closePanel" @keydown.esc="closePanel">
        <transition name="mcp-panel-zoom" appear>
          <div v-if="panelOpen" :class="['mcp-panel', theme ? 'mcp-panel-dark' : 'mcp-panel-light']" @click.stop>
            <div class="mcp-panel-accent" :class="{ on: enabled }"></div>
            <div class="mcp-panel-head">
              <div class="mcp-panel-title-wrap">
                <span class="mcp-panel-title">MCP 配置</span>
                <span class="mcp-panel-sub">供 Cursor 等客户端连接本机抓包桥</span>
              </div>
              <button type="button" class="mcp-panel-close" aria-label="关闭" @click="closePanel">×</button>
            </div>
            <div class="mcp-panel-body">
              <div class="mcp-card mcp-card-switch">
                <div class="mcp-switch-row">
                  <div>
                    <div class="mcp-switch-label">MCP 服务</div>
                    <div class="mcp-switch-hint">{{ enabled ? "运行中，底部状态为绿色呼吸灯" : "已关闭，需手动启用" }}</div>
                  </div>
                  <el-switch
                      :model-value="enabled"
                      :loading="busy"
                      :disabled="busy"
                      inline-prompt
                      active-text="开"
                      inactive-text="关"
                      @change="onToggle"
                  />
                </div>
                <div class="mcp-status-row">
                  <div class="mcp-status-pill" :class="enabled ? 'on' : 'off'">
                    <span class="mcp-status-dot"></span>
                    {{ enabled ? "已启用" : "未启用" }}
                  </div>
                  <div v-if="!enabled" class="mcp-port-inline">
                    <label for="mcp-port" class="mcp-port-inline-label">端口</label>
                    <div class="mcp-port-field">
                      <input
                          id="mcp-port"
                          v-model.number="port"
                          type="number"
                          min="1"
                          max="65535"
                          class="mcp-port-input"
                          :disabled="busy"
                          inputmode="numeric"
                          aria-label="监听端口"
                      />
                    </div>
                  </div>
                </div>
              </div>

              <div v-if="!enabled" class="mcp-cap-bar">
                <button type="button" class="mcp-cap-btn" @click="openCapabilities">查看 MCP 能力</button>
              </div>

              <div v-if="enabled" class="mcp-code-wrap">
                <div class="mcp-code-head">
                  <span class="mcp-code-label">mcp.json 预览</span>
                  <div class="mcp-code-head-actions">
                    <button type="button" class="mcp-cap-btn mcp-cap-btn-inline" @click="openCapabilities">
                      查看 MCP 能力
                    </button>
                    <button
                        type="button"
                        class="mcp-copy-btn mcp-copy-btn-inline"
                        :disabled="busy"
                        @click="copyConfig"
                    >
                      {{ copied ? "已复制" : "复制 mcp.json 配置" }}
                    </button>
                  </div>
                </div>
                <pre class="mcp-code">{{ mcpJson }}</pre>
              </div>

              <p v-if="errMsg" class="mcp-err">{{ errMsg }}</p>
            </div>
          </div>
        </transition>
      </div>
    </transition>
  </teleport>

  <teleport to="body">
    <transition name="mcp-overlay-fade">
      <div v-if="capDialogOpen" class="mcp-cap-overlay" @click.self="closeCapabilities">
        <div :class="['mcp-cap-dialog', theme ? 'mcp-cap-dialog-dark' : 'mcp-cap-dialog-light']" @click.stop>
          <div class="mcp-cap-dialog-head">
            <span class="mcp-cap-dialog-title">MCP 能力说明</span>
            <div class="mcp-cap-dialog-actions">
              <button
                  v-if="capDocUrl"
                  type="button"
                  class="mcp-cap-link-btn"
                  @click="openDocExternal"
              >
                在浏览器打开
              </button>
              <button type="button" class="mcp-cap-dialog-close" aria-label="关闭" @click="closeCapabilities">×</button>
            </div>
          </div>
          <div v-if="capLoading" class="mcp-cap-loading">加载中…</div>
          <iframe v-else-if="capDocUrl" class="mcp-cap-iframe" :src="capDocUrl" title="MCP 文档"/>
          <div v-else class="mcp-cap-list">
            <div v-for="item in capOps" :key="item.op" class="mcp-cap-item">
              <div class="mcp-cap-op">{{ item.op }}</div>
              <p v-if="item.args" class="mcp-cap-meta"><span>参数</span>{{ item.args }}</p>
              <p class="mcp-cap-desc">{{ item.description }}</p>
              <p v-if="item.returns" class="mcp-cap-meta"><span>返回</span>{{ item.returns }}</p>
            </div>
          </div>
        </div>
      </div>
    </transition>
  </teleport>
</template>

<script>
import {Config_IsDark} from "../../config/Config.js";
import {Events} from "@wailsio/runtime";
import {Tour_Add} from "../Tour";
import {
  MCPEnable,
  MCPDisable,
  MCPStatusJSON,
  MCPListOpsJSON,
  MCPDocURL,
} from "../../../../bindings/changeme/Service/appmain.js";

const COPY_FEEDBACK_MS = 1200;

export default {
  data() {
    return {
      enabled: false,
      port: 6987,
      streamURL: "",
      panelOpen: false,
      busy: false,
      errMsg: "",
      copied: false,
      copyTimer: null,
      capDialogOpen: false,
      capLoading: false,
      capDocUrl: "",
      capOps: [],
    };
  },
  computed: {
    theme() {
      return Config_IsDark.value;
    },
    mcpJson() {
      const url = this.streamURL || `http://127.0.0.1:${this.port}/sunnynet/tools/mcp`;
      return JSON.stringify({mcpServers: {SunnyNetTools: {url}}}, null, 2);
    },
  },
  mounted() {
    Tour_Add(this.$refs.mcpBar, 4, "MCP", "查看 MCP 状态；启用后可复制 mcp.json、查看 MCP 能力说明");
    this.refreshStatus();
    Events.On("mcpBridgeChanged", (obj) => {
      const raw = obj?.data?.[0] ?? obj?.data ?? "";
      this.applyStatusPayload(raw);
    });
  },
  beforeUnmount() {
    if (this.copyTimer) {
      clearTimeout(this.copyTimer);
    }
  },
  methods: {
    applyStatusPayload(raw) {
      let st = null;
      try {
        st = typeof raw === "string" ? JSON.parse(raw) : raw;
      } catch (_) {
        return;
      }
      if (!st || typeof st !== "object") {
        return;
      }
      this.enabled = !!(st.enabled || st.httpEnabled);
      if (st.lastPort > 0) {
        this.port = st.lastPort;
      } else if (st.defaultPort > 0) {
        this.port = st.defaultPort;
      }
      if (st.mcpStreamableURL) {
        this.streamURL = st.mcpStreamableURL;
      } else if (st.httpListenAddr) {
        this.streamURL = `http://${st.httpListenAddr}/sunnynet/tools/mcp`;
      } else if (!this.enabled) {
        this.streamURL = "";
      }
    },
    async refreshStatus() {
      try {
        const json = await MCPStatusJSON();
        this.applyStatusPayload(json);
      } catch (_) {
        this.enabled = false;
        this.streamURL = "";
      }
    },
    openPanel() {
      this.panelOpen = true;
      this.errMsg = "";
      this.refreshStatus();
    },
    closePanel() {
      this.panelOpen = false;
      this.errMsg = "";
      this.closeCapabilities();
    },
    closeCapabilities() {
      this.capDialogOpen = false;
      this.capLoading = false;
    },
    async openCapabilities() {
      this.capDialogOpen = true;
      this.capLoading = true;
      this.capDocUrl = "";
      this.capOps = [];
      try {
        const [opsRaw, docUrl] = await Promise.all([MCPListOpsJSON(), MCPDocURL()]);
        this.capDocUrl = (docUrl || "").trim();
        let env = null;
        try {
          env = typeof opsRaw === "string" ? JSON.parse(opsRaw) : opsRaw;
        } catch (_) {
          env = null;
        }
        this.capOps = Array.isArray(env?.capabilities) ? env.capabilities : [];
      } catch (e) {
        this.errMsg = String(e?.message || e);
        this.closeCapabilities();
      } finally {
        this.capLoading = false;
      }
    },
    openDocExternal() {
      if (!this.capDocUrl) {
        return;
      }
      try {
        const win = window.open(this.capDocUrl, "_blank");
        if (win) {
          win.focus();
        }
      } catch (_) {
        /* ignore */
      }
    },
    async onToggle(wantOn) {
      if (wantOn === this.enabled) {
        return;
      }
      if (wantOn) {
        await this.enableMcp();
      } else {
        await this.disableMcp();
      }
    },
    async enableMcp() {
      this.busy = true;
      this.errMsg = "";
      try {
        const msg = await MCPEnable(this.port || 6987);
        if (msg) {
          this.errMsg = msg;
          return;
        }
        await this.refreshStatus();
      } catch (e) {
        this.errMsg = String(e?.message || e);
      } finally {
        this.busy = false;
      }
    },
    async disableMcp() {
      this.busy = true;
      this.errMsg = "";
      try {
        const msg = await MCPDisable();
        if (msg) {
          this.errMsg = msg;
        }
        await this.refreshStatus();
      } catch (e) {
        this.errMsg = String(e?.message || e);
      } finally {
        this.busy = false;
      }
    },
    async copyConfig() {
      const text = this.mcpJson;
      try {
        await navigator.clipboard.writeText(text);
      } catch (_) {
        const ta = document.createElement("textarea");
        ta.value = text;
        ta.style.position = "fixed";
        ta.style.left = "-9999px";
        document.body.appendChild(ta);
        ta.select();
        document.execCommand("copy");
        document.body.removeChild(ta);
      }
      this.copied = true;
      if (this.copyTimer) {
        clearTimeout(this.copyTimer);
      }
      this.copyTimer = setTimeout(() => {
        this.copied = false;
        this.copyTimer = null;
      }, COPY_FEEDBACK_MS);
    },
  },
};
</script>

<style scoped>
.mcp-bar {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 0 8px;
  min-width: 72px;
  cursor: pointer;
  font-size: 15px;
  user-select: none;
  border-radius: 4px;
  transition: background 0.15s;
}

.mcp-bar:hover {
  background: rgba(148, 163, 184, 0.12);
}

.mcp-icon {
  width: 16px;
  height: 16px;
  flex-shrink: 0;
}

.mcp-icon-on {
  color: #22c55e;
}

.mcp-icon-off {
  color: #9ca3af;
}

.mcp-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.mcp-dot-on {
  background: #22c55e;
  animation: mcp-breathe 2s ease-in-out infinite;
}

.mcp-dot-off {
  background: #9ca3af;
}

.mcp-label {
  font-weight: 600;
  letter-spacing: 0.02em;
}

.mcp-label-on {
  color: #22c55e;
}

.mcp-label-off {
  color: #9ca3af;
}

@keyframes mcp-breathe {
  0%, 100% {
    opacity: 1;
    box-shadow: 0 0 0 0 rgba(34, 197, 94, 0.55);
  }
  50% {
    opacity: 0.55;
    box-shadow: 0 0 8px 2px rgba(34, 197, 94, 0.45);
  }
}

/* 遮罩：不覆盖顶栏 30px、底栏 30px（与 Home Header / Footer 一致） */
.mcp-overlay {
  position: fixed;
  top: 30px;
  left: 0;
  right: 0;
  bottom: 30px;
  z-index: 99999;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
  background: transparent;
}

.mcp-overlay-fade-enter-active,
.mcp-overlay-fade-leave-active {
  transition: opacity 0.22s ease;
}

.mcp-overlay-fade-enter-from,
.mcp-overlay-fade-leave-to {
  opacity: 0;
}

.mcp-panel-zoom-enter-active {
  transition: transform 0.28s cubic-bezier(0.22, 1, 0.36, 1), opacity 0.22s ease;
}

.mcp-panel-zoom-leave-active {
  transition: transform 0.18s ease, opacity 0.18s ease;
}

.mcp-panel-zoom-enter-from,
.mcp-panel-zoom-leave-to {
  transform: scale(0.94) translateY(8px);
  opacity: 0;
}

.mcp-panel {
  position: relative;
  width: min(480px, 100%);
  max-height: min(82vh, 620px);
  border-radius: 14px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  box-shadow:
      0 0 0 1px rgba(255, 255, 255, 0.06),
      0 24px 64px rgba(0, 0, 0, 0.45),
      0 8px 24px rgba(0, 0, 0, 0.25);
}

.mcp-panel-accent {
  height: 3px;
  background: linear-gradient(90deg, #64748b, #475569);
  transition: background 0.3s;
}

.mcp-panel-accent.on {
  background: linear-gradient(90deg, #22c55e, #4ade80, #22c55e);
  background-size: 200% 100%;
  animation: mcp-accent-flow 3s linear infinite;
}

@keyframes mcp-accent-flow {
  0% {
    background-position: 0% 50%;
  }
  100% {
    background-position: 200% 50%;
  }
}

.mcp-panel-dark {
  background: linear-gradient(165deg, #1e293b 0%, #0f172a 100%);
  color: #e2e8f0;
}

.mcp-panel-light {
  background: linear-gradient(165deg, #ffffff 0%, #f1f5f9 100%);
  color: #0f172a;
  box-shadow:
      0 0 0 1px rgba(15, 23, 42, 0.08),
      0 24px 48px rgba(15, 23, 42, 0.12);
}

.mcp-panel-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  padding: 16px 18px 12px;
  gap: 12px;
}

.mcp-panel-title-wrap {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.mcp-panel-title {
  font-size: 17px;
  font-weight: 600;
}

.mcp-panel-sub {
  font-size: 12px;
  opacity: 0.65;
}

.mcp-panel-close {
  border: none;
  background: rgba(148, 163, 184, 0.15);
  width: 32px;
  height: 32px;
  border-radius: 8px;
  font-size: 20px;
  line-height: 1;
  cursor: pointer;
  color: inherit;
  flex-shrink: 0;
  transition: background 0.15s;
}

.mcp-panel-close:hover {
  background: rgba(148, 163, 184, 0.28);
}

.mcp-panel-body {
  padding: 0 18px 18px;
  overflow: auto;
}

.mcp-card {
  border-radius: 10px;
  padding: 14px;
  margin-bottom: 12px;
  background: rgba(15, 23, 42, 0.35);
  border: 1px solid rgba(148, 163, 184, 0.12);
}

.mcp-panel-light .mcp-card {
  background: rgba(255, 255, 255, 0.7);
  border-color: rgba(15, 23, 42, 0.08);
}

.mcp-card-switch {
  margin-bottom: 12px;
}

.mcp-switch-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.mcp-switch-label {
  font-size: 15px;
  font-weight: 600;
}

.mcp-switch-hint {
  font-size: 12px;
  opacity: 0.65;
  margin-top: 4px;
}

.mcp-status-row {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 10px 12px;
  margin-top: 12px;
}

.mcp-status-pill {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 500;
  flex-shrink: 0;
}

.mcp-status-pill.on {
  background: rgba(34, 197, 94, 0.15);
  color: #4ade80;
}

.mcp-status-pill.off {
  background: rgba(148, 163, 184, 0.12);
  color: #94a3b8;
}

.mcp-status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: currentColor;
}

.mcp-status-pill.on .mcp-status-dot {
  animation: mcp-breathe 2s ease-in-out infinite;
}

.mcp-port-inline {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-left: auto;
  flex-shrink: 0;
}

.mcp-port-inline-label {
  font-size: 12px;
  font-weight: 500;
  opacity: 0.7;
  white-space: nowrap;
  user-select: none;
}

.mcp-port-field {
  display: flex;
  align-items: center;
  min-width: 92px;
  padding: 0 10px;
  border-radius: 10px;
  background: rgba(15, 23, 42, 0.45);
  border: 1px solid rgba(148, 163, 184, 0.22);
  box-shadow: inset 0 1px 2px rgba(0, 0, 0, 0.18);
  transition: border-color 0.2s ease, box-shadow 0.2s ease, background 0.2s ease;
}

.mcp-port-field:focus-within {
  border-color: rgba(34, 197, 94, 0.55);
  background: rgba(15, 23, 42, 0.6);
  box-shadow: inset 0 1px 2px rgba(0, 0, 0, 0.12), 0 0 0 3px rgba(34, 197, 94, 0.14);
}

.mcp-panel-light .mcp-port-field {
  background: rgba(255, 255, 255, 0.95);
  border-color: rgba(15, 23, 42, 0.12);
  box-shadow: inset 0 1px 2px rgba(15, 23, 42, 0.04);
}

.mcp-panel-light .mcp-port-field:focus-within {
  background: #fff;
  border-color: rgba(34, 197, 94, 0.45);
  box-shadow: inset 0 1px 2px rgba(15, 23, 42, 0.04), 0 0 0 3px rgba(34, 197, 94, 0.12);
}

.mcp-port-input {
  width: 72px;
  padding: 8px 0;
  border: none;
  background: transparent;
  font-size: 14px;
  font-weight: 600;
  font-variant-numeric: tabular-nums;
  text-align: center;
  outline: none;
  color: inherit;
}

.mcp-port-input:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.mcp-port-input::-webkit-outer-spin-button,
.mcp-port-input::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}

.mcp-port-input[type="number"] {
  -moz-appearance: textfield;
  appearance: textfield;
}

.mcp-url-hint {
  margin: 0 0 12px;
  padding: 8px 10px;
  font-size: 11px;
  word-break: break-all;
  opacity: 0.55;
  border-radius: 8px;
  background: rgba(15, 23, 42, 0.25);
}

.mcp-panel-light .mcp-url-hint {
  background: rgba(15, 23, 42, 0.04);
}

.mcp-copy-btn {
  border: none;
  border-radius: 8px;
  padding: 8px 18px;
  font-size: 13px;
  cursor: pointer;
  transition: opacity 0.15s, transform 0.1s;
  background: rgba(34, 197, 94, 0.15);
  color: #22c55e;
  border: 1px solid rgba(34, 197, 94, 0.4);
}

.mcp-panel-light .mcp-copy-btn {
  background: #ecfdf5;
  color: #15803d;
  border-color: rgba(21, 128, 61, 0.35);
}

.mcp-copy-btn-inline {
  flex-shrink: 0;
  padding: 5px 12px;
  font-size: 12px;
  border-radius: 7px;
  white-space: nowrap;
}

.mcp-copy-btn:active:not(:disabled) {
  transform: scale(0.98);
}

.mcp-copy-btn:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}

.mcp-err {
  color: #f87171;
  font-size: 13px;
  margin: 0 0 10px;
  padding: 8px 10px;
  border-radius: 8px;
  background: rgba(248, 113, 113, 0.1);
}

.mcp-code-wrap {
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid rgba(148, 163, 184, 0.16);
  background: rgba(15, 23, 42, 0.35);
  box-shadow: inset 0 1px 2px rgba(0, 0, 0, 0.12);
}

.mcp-panel-light .mcp-code-wrap {
  background: rgba(248, 250, 252, 0.9);
  border-color: rgba(15, 23, 42, 0.08);
  box-shadow: inset 0 1px 2px rgba(15, 23, 42, 0.04);
}

.mcp-code-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  padding: 10px 12px;
  border-bottom: 1px solid rgba(148, 163, 184, 0.12);
  background: rgba(15, 23, 42, 0.25);
}

.mcp-code-head-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.mcp-cap-bar {
  display: flex;
  justify-content: center;
  margin-bottom: 12px;
}

.mcp-cap-btn {
  border: 1px solid rgba(96, 165, 250, 0.45);
  border-radius: 8px;
  padding: 8px 16px;
  font-size: 13px;
  cursor: pointer;
  background: rgba(59, 130, 246, 0.14);
  color: #60a5fa;
  transition: opacity 0.15s, transform 0.1s;
}

.mcp-panel-light .mcp-cap-btn {
  background: #eff6ff;
  color: #2563eb;
  border-color: rgba(37, 99, 235, 0.35);
}

.mcp-cap-btn-inline {
  padding: 5px 12px;
  font-size: 12px;
  border-radius: 7px;
  white-space: nowrap;
}

.mcp-cap-btn:active {
  transform: scale(0.98);
}

.mcp-panel-light .mcp-code-head {
  background: rgba(15, 23, 42, 0.04);
}

.mcp-code-label {
  font-size: 12px;
  font-weight: 500;
  opacity: 0.72;
}

.mcp-code {
  margin: 0;
  padding: 14px 16px;
  font-family: Consolas, Monaco, "Cascadia Code", "SF Mono", monospace;
  font-size: 12px;
  line-height: 1.6;
  letter-spacing: 0.02em;
  overflow: auto;
  max-height: 168px;
  white-space: pre;
  color: #cbd5e1;
  background: rgba(2, 6, 23, 0.45);
  scrollbar-width: thin;
  scrollbar-color: rgba(148, 163, 184, 0.35) transparent;
}

.mcp-panel-light .mcp-code {
  color: #334155;
  background: #f1f5f9;
}

/* 能力层：顶 28px、底栏 30px 不遮挡，宽 100% */
.mcp-cap-overlay {
  position: fixed;
  top: 28px;
  left: 0;
  right: 0;
  bottom: 30px;
  width: 100%;
  height: calc(100vh - 58px);
  z-index: 100000;
  display: flex;
  flex-direction: column;
  align-items: stretch;
  justify-content: stretch;
  padding: 0;
  background: transparent;
  box-sizing: border-box;
}

.mcp-cap-dialog {
  width: 100%;
  height: 100%;
  min-height: 0;
  max-height: none;
  display: flex;
  flex-direction: column;
  border-radius: 0;
  overflow: hidden;
  border: none;
  box-shadow: none;
}

.mcp-cap-dialog-dark {
  background: #0f172a;
  color: #e2e8f0;
}

.mcp-cap-dialog-light {
  background: #f8fafc;
  color: #0f172a;
}

.mcp-cap-dialog-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 12px 16px;
  border-bottom: 1px solid rgba(148, 163, 184, 0.15);
  flex-shrink: 0;
}

.mcp-cap-dialog-title {
  font-size: 15px;
  font-weight: 600;
}

.mcp-cap-dialog-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.mcp-cap-link-btn {
  border: none;
  border-radius: 7px;
  padding: 5px 10px;
  font-size: 12px;
  cursor: pointer;
  background: rgba(59, 130, 246, 0.15);
  color: #60a5fa;
}

.mcp-panel-light .mcp-cap-link-btn {
  color: #2563eb;
  background: #dbeafe;
}

.mcp-cap-dialog-close {
  border: none;
  background: rgba(148, 163, 184, 0.15);
  width: 30px;
  height: 30px;
  border-radius: 8px;
  font-size: 18px;
  line-height: 1;
  cursor: pointer;
  color: inherit;
}

.mcp-cap-loading {
  padding: 48px 16px;
  text-align: center;
  opacity: 0.7;
}

.mcp-cap-iframe {
  flex: 1;
  min-height: 0;
  width: 100%;
  border: none;
  background: #020617;
}

.mcp-cap-list {
  flex: 1;
  min-height: 0;
  overflow: auto;
  padding: 12px 14px 16px;
  scrollbar-gutter: stable;
  scrollbar-width: thin;
  scrollbar-color: rgba(96, 165, 250, 0.55) rgba(15, 23, 42, 0.35);
}

.mcp-cap-list::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

.mcp-cap-list::-webkit-scrollbar-track {
  background: rgba(15, 23, 42, 0.35);
  border-radius: 8px;
  margin: 4px 0;
}

.mcp-cap-list::-webkit-scrollbar-thumb {
  background: linear-gradient(180deg, rgba(96, 165, 250, 0.65), rgba(34, 197, 94, 0.45));
  border-radius: 8px;
  border: 2px solid transparent;
  background-clip: padding-box;
}

.mcp-cap-list::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(180deg, rgba(125, 211, 252, 0.85), rgba(74, 222, 128, 0.65));
  background-clip: padding-box;
}

.mcp-panel-light .mcp-cap-list {
  scrollbar-color: rgba(37, 99, 235, 0.45) rgba(226, 232, 240, 0.9);
}

.mcp-panel-light .mcp-cap-list::-webkit-scrollbar-track {
  background: rgba(226, 232, 240, 0.95);
}

.mcp-panel-light .mcp-cap-list::-webkit-scrollbar-thumb {
  background: linear-gradient(180deg, rgba(59, 130, 246, 0.55), rgba(34, 197, 94, 0.4));
  background-clip: padding-box;
}

.mcp-cap-item {
  padding: 12px 0;
  border-bottom: 1px solid rgba(148, 163, 184, 0.12);
}

.mcp-cap-item:last-child {
  border-bottom: none;
}

.mcp-cap-op {
  font-family: Consolas, Monaco, "Cascadia Code", monospace;
  font-size: 13px;
  font-weight: 600;
  color: #4ade80;
  margin-bottom: 6px;
}

.mcp-panel-light .mcp-cap-op {
  color: #15803d;
}

.mcp-cap-desc {
  margin: 0;
  font-size: 12px;
  line-height: 1.55;
  opacity: 0.88;
  white-space: pre-wrap;
  word-break: break-word;
}

.mcp-cap-meta {
  margin: 4px 0 0;
  font-size: 11px;
  line-height: 1.45;
  opacity: 0.72;
}

.mcp-cap-meta span {
  display: inline-block;
  min-width: 2.5em;
  margin-right: 6px;
  font-weight: 600;
  opacity: 0.9;
}
</style>
