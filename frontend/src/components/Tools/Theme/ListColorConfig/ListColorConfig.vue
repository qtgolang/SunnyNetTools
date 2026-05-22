<template>
  <ag-grid-vue
      ref="agGrid"
      style="height: 100%;width: 100%"
      :defaultColDef="defaultColDef"
      :rowData="RowData"
      :columnDefs="columns"
      :cellSelection="true"
      :enableCharts="true"
      :grid-options="gridOptions"
      :theme="theme"
  >
  </ag-grid-vue>
</template>

<script>
import {AgGridVue} from "ag-grid-vue3";
import {ExcelExportModule} from '@ag-grid-enterprise/excel-export';
import ColorConfigRenderer from './ColorConfig.vue';
import {Config_IsDark} from "../../../config/Config";
import {GetListColor} from "../../../../../bindings/changeme/Service/appmain";
import {Events} from "@wailsio/runtime";

export default {
  props: ['readOnly', "theme"],
  watch: {
    readOnly(value) {
      this.ReadOnly = value
    },
    IsDark(n) {
      this.applyInit(n)
    },
  },
  components: {
    'ag-grid-vue': AgGridVue, imageRenderer: ColorConfigRenderer,
  },
  computed: {},
  data() {
    return {
      get IsDark() {
        return Config_IsDark.value
      },
      set IsDark(value) {
        Config_IsDark.value = value
      },
      //是否被修改
      IsHasModify: false,
      //当前选中行
      agSelectedLine: null,
      agSelectedXY: {
        top: 0,
        right: 0
      },
      overlayNoRowsTemplate: `<span style="padding: 20px;" id="HookMessageText">无内容</span>`,
      rightModules: [ExcelExportModule],
      gridOptions: {
        getRowId: (params) => params.data.id,
        rowSelection: {
          mode: 'singleRow',
          checkboxes: false,  // 启用行选择的复选框
        },
        onRowClicked: this.onRowClicked,
        onCellFocused: this.onCellFocused,
        getContextMenuItems: this.onContextMenuItems,
        onCellValueChanged: (event) => {
          this.IsHasModify = true
        },
        // 禁用自动滚动到第一行
        suppressScrollOnNewData: true,
      },
      defaultColDef: {
        flex: 1,
        minWidth: 10,
        // 禁用全部列的排序功能
        sortable: false,
        filter: false,
        suppressHeaderMenuButton: true,
        floatingFilter: false,
        resizable: true,
        cellClass: 'no-border'
      },
      MenuItems: [],
      RowData: [
        {"名称": "TCP请求", "字体颜色": "#111", "说明": "如果是 TCP/TLS-TCP/TCP-Must/ 请求", id: "1"},
        {"名称": "UDP请求", "字体颜色": "#111", "说明": "如果是UDP请求", id: "2"},
        {"名称": "Websocket", "字体颜色": "#111", "说明": "如果是WS/WSS请求", id: "3"},
        {"名称": "CSS", "字体颜色": "#111", "说明": "如果响应是 CSS 文件 ", id: "4"},
        {"名称": "javaScript", "字体颜色": "#111", "说明": "如果响应是 js 文件", id: "5"},
        {"名称": "图片", "字体颜色": "#111", "说明": "如果响应是 图片 文件 ", id: "6"},
        {"名称": "文档", "字体颜色": "#111", "说明": "如果响应是 TXT/HTML/...文档类型 ", id: "7"},
        {"名称": "错误请求", "字体颜色": "#111", "说明": "如果响应状态码是 -1 的请求", id: "8"},
        {"名称": "重定向请求", "字体颜色": "#111", "说明": "如果响应状态码是 301/302 的重定向请求", id: "9"},
        {"名称": "40x", "字体颜色": "#111", "说明": "如果响应状态码是 401/403/404/500 的请求", id: "10"},
        {"名称": "常规项", "字体颜色": "#111", "说明": "如果不是上列的任何一项", id: "99"},
        {"名称": "", "字体颜色": "", "说明": "恢复到默认", id: "gg"}
      ],
      columns: [
        {
          field: "名称", tooltipField: '名称',
          minWidth: 100,
          width: 100,
          maxWidth: 100,
          // 不显示过滤器
          menuTabs: [],
          //禁止列拖动
          suppressMovable: true,
          cellStyle: {'text-align': 'left'},
        },
        {
          field: "字体颜色", tooltipField: '字体颜色',
          // 不显示过滤器
          menuTabs: [],
          //禁止列拖动
          suppressMovable: true,
          cellRenderer: 'imageRenderer',
          cellStyle: {'text-align': 'left'},
          minWidth: 120,
          width: 120,
          maxWidth: 120,
        },
        {
          field: "说明", tooltipField: '说明',
          // 不显示过滤器
          menuTabs: [],
          //禁止列拖动
          suppressMovable: true,
          cellStyle: {'text-align': 'left'},
        },
      ],
      ReadOnly: true,
    }
  }
  ,
  methods: {
    applyInit(IsDark) {
      GetListColor(IsDark).then((res) => {
        const array = []
        Object.entries(res).forEach(([key, value]) => {
          const focusedRowNode = this.agGridApi.getRowNode(key);
          if (focusedRowNode) {
            focusedRowNode.data["字体颜色"] = value
            array.push(focusedRowNode)
          }
        });
        const focusedRowNode = this.agGridApi.getRowNode("gg");
        if (focusedRowNode) {
          focusedRowNode.data["说明"] = "恢复到默认" + (IsDark ? "暗色配置" : "亮色配置")
          array.push(focusedRowNode)
        }
        this.agGridApi.redrawRows({rowNodes: array});
      })
    },
    runTask() {
      requestAnimationFrame(() => {
        const elements = document.querySelectorAll('.el-popper.is-light.el-color-picker__panel.el-color-dropdown');
        for (let i = 0; i < elements.length; i++) {
          const left = parseInt(elements[i].style.left)
          if (left !== this.agSelectedXY.right) {
            elements[i].style.left = this.agSelectedXY.right + "px";
          }
        }
        this.runTask();
      });
    },
    onContextMenuItems() {
      let array = [];
      for (let i = 0; i < this.MenuItems.length; i++) {
        if (this.MenuItems[i].visible !== false) {
          array.push(this.MenuItems[i])
        }
      }
      return array
    },
    onRowClicked(params) {
      const rect = params.event.target.getBoundingClientRect();
      this.agSelectedXY.top = rect.top;
      this.agSelectedXY.right = rect.right;
      params.node.setSelected(true);
      this.agSelectedLine = params.node
    },
    SelectedLine(index) {
      const focusedRowNode = this.agGridApi.getRowNode(index + "");
      if (focusedRowNode) {
        if (this.agSelectedLine === null) {
          focusedRowNode.setSelected(true);
          this.agSelectedLine = focusedRowNode
          return
        }
        if (focusedRowNode.rowIndex !== this.agSelectedLine.rowIndex && focusedRowNode.id !== this.agSelectedLine.id) {
          focusedRowNode.setSelected(true);
          this.agSelectedLine = focusedRowNode
        }
      }
    },
    onCellFocused(event) {
      const rowNode = event.api.getDisplayedRowAtIndex(event.rowIndex);
      this.SelectedLine(rowNode.data.id)

    },
  },
  mounted() {
    this.agGridApi = this.$refs.agGrid.api;
    const observer = new MutationObserver(() => {
      requestAnimationFrame(() => {
        const elements = document.querySelectorAll('.el-popper.is-light.el-color-picker__panel.el-color-dropdown');
        elements.forEach(el => {
          if (parseInt(el.style.left) !== this.agSelectedXY.right) {
            el.style.left = this.agSelectedXY.right + "px";
          }
        });
      });
    });
    observer.observe(document.body, {attributes: true, subtree: true, attributeFilter: ['style']});
    this.applyInit(this.IsDark)
    Events.On("RestColor", (obj) => {
      let o = obj.data[0];
      if (!o) {
        o = obj.data;
      }
      this.applyInit(o)
    })
  }
}
</script>

<style>

</style>