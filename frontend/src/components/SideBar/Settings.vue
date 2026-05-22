<script>
import {AgGridVue} from "ag-grid-vue3";
import Tools from "../Tools/Tools.vue";
import BaseSettings from "./Settings/BaseSettings.vue";
import Keys from "./Settings/Keys.vue";
import HTTPSProto from "./Settings/HTTPSProto.vue";
import MustTcp from "./Settings/MustTcp.vue";
import Way from "./Settings/Way.vue";
import Host from "./Settings/Host.vue";
import ReplaceBody from "./Settings/ReplaceBody.vue";
import Dns from "./Settings/Dns.vue";
import FormatSection from "./Settings/section.vue";
import {OpenTools} from "../CallbackEventsOn";
import {Config_IsDark} from "../config/Config.js";

export default {
  components: {
    FormatSection,
    Dns,
    Replace: ReplaceBody,
    Host,
    Way,
    MustTcp,
    HTTPSProto,
    Keys,
    BaseSettings,
    "ag-grid-vue": AgGridVue,
    Tools,
  },

  data() {
    return {
      Show: {
        Base: false,
        Keys: false,
        HTTP: false,
        MustTCP: false,
        Way: false,
        Host: false,
        Cert: false,
        ReplaceBody: false,
      },
    };
  },

  computed: {
    // 关键：theme 必须放在 computed，才能跟随 Config_IsDark 自动更新
    theme: {
      get() {
        return Config_IsDark.value;
      },
      set(v) {
        if (Config_IsDark.value !== v) Config_IsDark.value = v;
      },
    },

  },

  watch: {
    // 保持原逻辑：任意一个打开就把其他关掉
    Show: {
      deep: true,
      handler(n, o) {
        const oldVal = o || {};
        const openedKey = Object.keys(n).find((k) => n[k] && !oldVal[k]);
        if (openedKey) this.setShow(openedKey);
      },
    },
  },

  methods: {
    OpenTools,
    setShow(key) {
      Object.keys(this.Show).forEach((k) => {
        if (k === key) return;
        if (k === "ReplaceBody" || k === "Cert") OpenTools(k, false, "");
        this.Show[k] = false;
      });
    },

    ShowCert() {
      this.setShow("Cert");
      this.Show.Cert = true;
      OpenTools("Cert", true, "");
    },

    ShowRep() {
      this.setShow("ReplaceBody");
      this.Show.ReplaceBody = true;
      OpenTools("ReplaceBody", true, "");
    },

  },
};
</script>

<template>
  <div style="width: 100%; overflow: auto; scrollbar-width: none">
    <div class="ag-chart-tab ag-chart-format" style="overflow: auto; scrollbar-width: none">
      <div class="ag-chart-format-wrapper">
        <!-- 常规设置 -->
        <FormatSection
            v-model="Show.Base"
            title="常规设置"
            :bodyStyle="{ height: 'auto', width: '100%', overflow: 'hidden' }"
        >
          <BaseSettings/>
        </FormatSection>

        <!-- 快捷键设置 -->
        <FormatSection v-model="Show.Keys" title="快捷键设置">
          <Keys/>
        </FormatSection>

        <!-- HTTPS协议设置 -->
        <FormatSection v-model="Show.HTTP" title="HTTPS协议设置">
          <HTTPSProto/>
        </FormatSection>

        <!-- Host设置 -->
        <FormatSection v-model="Show.Host" title="Host设置">
          <Host/>
        </FormatSection>

        <!-- 强制TCP设置 -->
        <FormatSection v-model="Show.MustTCP" title="强制TCP设置">
          <MustTcp/>
        </FormatSection>

        <!-- 上游代理/二级代理 -->
        <FormatSection v-model="Show.Way" title="上游代理/二级代理">
          <Way/>
        </FormatSection>

        <!-- 数据替换设置：点击打开工具，不折叠 -->
        <FormatSection :modelValue="false" :collapsible="false" :useButtonTitle="true" @toggle="ShowRep">
          <template #title>
            <el-button size="small">拦截/数据替换设置</el-button>
          </template>
        </FormatSection>

        <!-- 请求证书设置：点击打开工具，不折叠 -->
        <FormatSection :modelValue="false" :collapsible="false" :useButtonTitle="true" @toggle="ShowCert">
          <template #title>
            <el-button size="small">请求证书设置</el-button>
          </template>
        </FormatSection>
      </div>
    </div>
  </div>
</template>

<style scoped>
.ag-chart-tabbed-menu-body::after {
  content: "";
  display: block;
  height: 0;
}
</style>
