<script>
import {AgGridVue} from "ag-grid-vue3";
import {
  Config_agGrid_API,
  Config_Find_Window,
  Config_Find_Window_Hide,
  Config_Find_Window_Show,
  Config_Theme_agGrid,
} from "../../config/Config.js";
import {AG_GRID_LOCALE_CN} from "../../config/AG_ZH_CN.js";
import Find from "./Find.vue";

export default {
  components: {Find, 'ag-grid-vue': AgGridVue},
  data() {
    return {
      agGridApi: null,
      gridOptions: {
        enableAdvancedFilter: true,
        suppressMovableColumns: true,
        localeText: AG_GRID_LOCALE_CN,
        popupParent: document.getElementById("appMain")
      },
      windowList: [],
      newPanel: null,
      topOffsetConfig: {
        original: 0,
        newPosition: 0,
        offsetDifference: 0,
      },
      Title: "🚀 高级搜索",
      setTitle: () => {
      },
      CloseButtonEvent: () => {
      },
      Height: 350,
      Width: 550,
      setNewHeight: () => {
      },
      CompleteCallback: () => {
      },
      isSearchInProgress: false
    }
  },
  mounted() {
    this.agGridApi = this.$refs.agGrid.api;
    Config_agGrid_API.value = this.$refs.agGrid.api;
    this.agGridApi.setSideBarPosition('right')
    Config_Find_Window_Show.value = () => {
      this.agGridApi.hideAdvancedFilterBuilder();
      this.windowList = Array.from(document.getElementsByClassName("ag-panel"));
      this.agGridApi.setGridOption('popupParent', document.getElementById("appMain"));
      this.agGridApi.showAdvancedFilterBuilder();
      this.hookFilterBuilder();
    }
    Config_Find_Window_Hide.value = this.agGridApi.hideAdvancedFilterBuilder
    Config_Find_Window.value = this
  },
  methods: {
    hookFilterBuilder(params) {
      // 获取当前所有 `ag-panel` 窗口
      const currentPanels = Array.from(document.getElementsByClassName("ag-panel"));

      // 过滤出新弹出的窗口
      const newPanels = currentPanels.filter(panel => !this.windowList.includes(panel));
      if (newPanels.length === 0) return;
      this.newPanel = newPanels[0];
      const closeButton = this.newPanel.querySelector(".ag-icon.ag-icon-cross.ag-panel-title-bar-button-icon");
      if (closeButton?.parentElement) {
        const parent = closeButton.parentElement; // ✅ 获取父元素
        parent.addEventListener("click", () => {
          delete this.newPanel.style.top; // 解除 `top` 劫持
          delete this.newPanel.style.width; // 解除 `width` 劫持
          delete this.newPanel.style.minWidth; // 解除 `minWidth` 劫持
          delete this.newPanel.style.minHeight; // 解除 `minHeight` 劫持
          delete this.newPanel.style.maxHeight; // 解除 `minHeight` 劫持
          this.CloseButtonEvent()
          this.isSearchInProgress = false
          this.newPanel = null;
        });
        parent.previousElementSibling.remove()
      }
      //禁止双击全屏
      const eTitleBarElement = this.newPanel.querySelector('[data-ref="eTitleBar"]');
      eTitleBarElement.addEventListener("dblclick", (event) => {
        event.stopImmediatePropagation();
      }, true);

      // ✅ 修改标题内容
      const titleElement = this.newPanel.querySelector('[data-ref="eTitle"]');
      if (titleElement) {
        titleElement.textContent = this.Title;
        this.setTitle = (eTitle) => {
          titleElement.textContent = eTitle;
        }
      }

      // ✅ 替换内容区域
      const contentElement = this.newPanel.querySelector('[data-ref="eList"]');
      if (contentElement) {
        contentElement.innerHTML = "";
        contentElement.appendChild(this.$refs.Find);
        const parent = contentElement.parentElement; // 获取父元素
        const lastChild = parent.lastElementChild; // 获取父元素的最后一个子元素
        if (lastChild) {
          parent.removeChild(lastChild); // 删除最后一个子元素
        }
      }
      // ✅ 计算并调整窗口高度
      const viewportHeight = window.innerHeight;
      const panelHeight = this.Height;
      const panelWidth = this.Width;
      const centeredTop = (viewportHeight - panelHeight) / 2;

      this.newPanel.style.setProperty("min-height", `${panelHeight}px`, "important");
      this.newPanel.style.setProperty("max-height", `${panelHeight}px`, "important");
      this.newPanel.style.setProperty("height", `${panelHeight}px`, "important");
      this.setNewHeight = (newHeight) => {
        this.Height = newHeight;
        this.newPanel.style.setProperty("min-height", `${newHeight}px`, "important");
        this.newPanel.style.setProperty("max-height", `${newHeight}px`, "important");
        this.newPanel.style.setProperty("height", `${newHeight}px`, "important");
      }
      // ✅ 计算 `top` 偏移量
      this.topOffsetConfig = {
        original: parseFloat(this.newPanel.style.top),  // 记录初始 `top`
        newPosition: centeredTop,                 // 计算出的居中 `top`
        offsetDifference: centeredTop - parseFloat(this.newPanel.style.top) // 差值
      };

      let adjustmentCount = 0;
      const topOffsetConfig = this.topOffsetConfig;
      // ✅ 劫持 `style.top`，确保 `top` 初始化时正确
      Object.defineProperty(this.newPanel.style, "top", {
        get() {
          return this.getPropertyValue("top");
        },
        set(value) {
          const numericValue = parseFloat(value);

          if (numericValue === topOffsetConfig.original) {
            if (adjustmentCount < 10) {
              this.setProperty("top", topOffsetConfig.newPosition + "px");
              return;
            }
            adjustmentCount = 10;
          }
          adjustmentCount++; //这里持续增加有没有什么问题，还有当 newPanel 元素删除时，会不会有什么问题
          this.setProperty("top", value);
        }
      });
      //锁定高度-宽度
      Object.defineProperty(this.newPanel.style, "minHeight", {
        get() {
          return this.getPropertyValue("minHeight");
        },
        set(value) {
          this.setProperty("minHeight", `${panelHeight}px`);
        }
      });
      Object.defineProperty(this.newPanel.style, "maxHeight", {
        get() {
          return this.getPropertyValue("maxHeight");
        },
        set(value) {
          this.setProperty("maxHeight", `${panelHeight}px`);
        }
      });
      Object.defineProperty(this.newPanel.style, "minWidth", {
        get() {
          return this.getPropertyValue("minWidth");
        },
        set(value) {
          this.setProperty("minWidth", `${panelWidth}px`);
        }
      });
      Object.defineProperty(this.newPanel.style, "width", {
        get() {
          return this.getPropertyValue("width");
        },
        set(value) {
          this.setProperty("width", `${panelWidth}px`);
        }
      });
      this.$nextTick(() => {
        this.$refs.FindMain.SetFocus()
      })
    }
  },
  watch: {},

  computed: {
    agTheme() {
      return Config_Theme_agGrid.value
    },
  },
}
</script>

<template>
  <div>
    <!-- 本文件重写ag-grid的弹出窗口，用来充当搜索窗口 -->
    <!-- 本文件重写ag-grid的弹出窗口，用来充当搜索窗口 -->
    <!-- 本文件重写ag-grid的弹出窗口，用来充当搜索窗口 -->
    <div style="width: 0;height: 0;position: absolute">
      <ag-grid-vue
          ref="agGrid"
          :enableCharts="true"
          :theme="agTheme"
          style="position: absolute"
          :grid-options="gridOptions"
          :loading="false"
          :allowContextMenuWithControlKey="true"
          :suppressCutToClipboard="true"
      />
    </div>
    <!-- 本文件重写ag-grid的弹出窗口，用来充当搜索窗口 -->
    <!-- 本文件重写ag-grid的弹出窗口，用来充当搜索窗口 -->
    <!-- 本文件重写ag-grid的弹出窗口，用来充当搜索窗口 -->
    <div ref="Find">
      <Find ref="FindMain" style="height: 100%;width: 100%;display: flex;justify-content: center"
            v-show="this.newPanel!==null"/>
    </div>
  </div>
</template>