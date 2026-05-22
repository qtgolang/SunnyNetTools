<script>
import {
  clearTextColorMap,
  Config_agGrid_API,
  Config_AutoRoll,
  Config_GOOS_IsWindows,
  Config_HTTP_Message_free,
  Config_IsDark,
  Config_Menu_isFileMenu,
  Config_SelectedRow,
} from "../../config/Config.js";
import {
  AppVersion,
  ClearAllSession,
  FreeAllRequest,
  McpFuncRes
} from "../../../../bindings/changeme/Service/appmain.js";
import {Keys_System_id_ALL_Release, Keys_System_id_Keys_Clear_Al, registerHotkeyFunction} from "../../config/Keys";
import WayContent from "./WayContent.vue";
import OpenSourceProtocol from "./OpenSourceProtocol.vue";
import OpenSource from "./OpenSource.vue";
import {Events, Window} from "@wailsio/runtime";

export default {
  components: {
    WayContent, OpenSource, OpenSourceProtocol
  },
  data() {
    return {
      get theme() {
        return Config_IsDark.value
      },
      set theme(newValue) {
        if (Config_IsDark.value !== newValue) Config_IsDark.value = newValue
      },
      get IsWindows() {
        return Config_GOOS_IsWindows.value
      },
      set IsWindows(value) {
        Config_GOOS_IsWindows.value = value
      },
      get Roll() {
        return Config_AutoRoll.value
      },
      set Roll(value) {
        Config_AutoRoll.value = value
      },
      version: "v1.0.0"
    }
  }
  ,
  mounted() {
    AppVersion().then(v => {
      this.version = v
    })
    registerHotkeyFunction(Keys_System_id_Keys_Clear_Al, this.clear)
    registerHotkeyFunction(Keys_System_id_ALL_Release, this.free)

    Events.On("mcp", async (evt) => {
      const mcp = evt?.data ?? {};
      const reply = (text) => {
        mcp.res = text;
        typeof McpFuncRes === "function" && McpFuncRes(mcp);
      };
      try {
        const page = String(mcp.page ?? "").toLowerCase();
        const tag = String(mcp.tag ?? "").toLowerCase();
        if (page !== "main" || tag !== "home") return;
        const msg = String(mcp.msg ?? "");
        switch (msg) {
          case "clear":
            this.clear()
            return reply("已全部删除");
          default:
            return;
        }
      } catch (e) {
        // 异常兜底：避免对端一直等
        try {
          reply("处理失败");
        } catch (_) {
        }
      }
    });
  }
  ,
  watch: {
    "Roll"(n, l) {
      this.$refs.AutoRoll.checked = n
    }
  }
  , computed: {
    getTheme() {
      return Config_IsDark.value ? "dark" : "light"
    },
    platform() {
      return this.IsWindows ? "windows" : "darwin"
    }
  },
  methods: {
    FileClick() {
      Config_Menu_isFileMenu.value = true;
      Config_agGrid_API.value.showContextMenu();
      setTimeout(() => {
        Config_Menu_isFileMenu.value = false;
      }, 200)
    }
    ,
    clear() {
      ClearAllSession().then((array) => {
        const api = Config_agGrid_API.value;
        api.setGridOption("rowData", []);
        Config_SelectedRow.value = [];
        api.clearCellSelection()
        clearTextColorMap()
      })
    }
    ,
    free() {
      Config_HTTP_Message_free.value(Config_SelectedRow.value)
      FreeAllRequest()
    }
    ,
    RollClick() {
      this.Roll = !this.Roll
    }
    ,
    handleClose() {
      Window.Close()
    }
    ,
    handleMinimize() {
      Window.Minimise()
    }
    ,
    handleMaximize() {
      Window.ToggleMaximise()
    }
    ,
    handleRestore() {

    }
    ,
  }
}
</script>

<template>
  <div style="width: 100%;height:30px;display: flex;position: relative;align-items: center">
    <v-titlebar
        :theme="getTheme"
        :show-icon="false"
        :platform="platform"
        ref="header"
        @close="handleClose"
        @minimize="handleMinimize"
        @maximize="handleMaximize"
        @unmaximize="handleRestore"
        style="margin-top: -2px"
    >
      <template v-slot:title>
        <div style="display: flex">
          <div style="left: 5px;position: relative; width: 80px">
            <button class="ag-button ag-advanced-filter-builder-button" data-ref="eBuilderFilterButton" tabindex="0"
                    @click="FileClick">
              <span data-ref="eBuilderFilterButtonIcon"><span class="ag-icon ag-icon-menu"></span></span>
              <span class="ag-advanced-filter-builder-button-label">文件</span>
            </button>
          </div>
          <div style="left: 5px;position: relative; width: 100px">
            <button class="ag-button ag-advanced-filter-builder-button" data-ref="eBuilderFilterButton" tabindex="0"
                    @click="clear">
              <span data-ref="eBuilderFilterButtonIcon"><span class="ag-icon ag-icon-cancel"></span></span>
              <span class="ag-advanced-filter-builder-button-label">全部删除</span>
            </button>
          </div>
          <div style="left: 5px;position: relative; width: 100px">
            <button class="ag-button ag-advanced-filter-builder-button" data-ref="eBuilderFilterButton" tabindex="0"
                    @click="free">
              <span data-ref="eBuilderFilterButtonIcon"><span class="ag-icon ag-icon-right"></span></span>
              <span class="ag-advanced-filter-builder-button-label">全部放行</span>
            </button>
          </div>
          <div style="left: 5px;position: relative; width: 100px">
            <div
                class="ag-labeled ag-label-align-right ag-checkbox ag-input-field ag-group-item ag-charts-advanced-settings-top-level-group-item"
                style="margin-top:0px;left: 5px;position: relative; width: 100px;font-size: 15px;display: flex;align-items: center;cursor: pointer;"
                @click="RollClick">
              <div class="ag-input-field-label ag-label ag-checkbox-label">
          <span
              class="notranslate immersive-translate-target-wrapper" style="font-size:14px">
            <span
                class="notranslate immersive-translate-target-translation-theme-none immersive-translate-target-translation-inline-wrapper-theme-none immersive-translate-target-translation-inline-wrapper">
              <span
                  class="notranslate immersive-translate-target-inner immersive-translate-target-translation-theme-none-inner">自动滚动</span>
            </span>
          </span>
              </div>
              <div class="ag-wrapper ag-input-wrapper ag-checkbox-input-wrapper">
                <input ref="AutoRoll" class="ag-input-field-input ag-checkbox-input" type="checkbox">
              </div>
            </div>
          </div>
        </div>
        <div style="display: flex; position: relative; width: 100%; --wails-draggable: drag">
          <!-- 居中的文字 -->
          <span style="width: 100%;justify-content: center;text-align: center;align-content: center;"
                @dblclick="handleMaximize">
            SunnyNet网络中间件 v{{ version }}
          </span>
          <div v-if="!this.IsWindows" style="width: 200px"></div>
        </div>
        <!-- 右侧按钮区域 -->
        <div style="margin-left: auto; display: flex; align-items: center; margin-right: 15px">
          <OpenSource/>
          <OpenSourceProtocol/>
          <WayContent style="width: 16px; height: 16px; margin-top: -5px; margin-left: 10px;"/>
        </div>

      </template>
    </v-titlebar>

  </div>
</template>
