<template>
  <ag-grid-vue
      ref="agGrid"
      style="height: 100%;width: 100%"
      :defaultColDef="defaultColDef"
      :rowData="RowData"
      :columnDefs="columns"
      :enableRangeSelection="true"
      :enableCharts="true"
      :modules="leftModules"
      :grid-options="gridOptions"
      :overlayNoRowsTemplate="overlayNoRowsTemplate"
  >
  </ag-grid-vue>
</template>

<script>
import '@ag-grid-community/styles/ag-grid.css';
import '@ag-grid-community/styles/ag-theme-balham.css';
import {AgGridVue} from '@ag-grid-community/vue3';
import {ClipboardModule} from '@ag-grid-enterprise/clipboard';
import {SetFilterModule} from '@ag-grid-enterprise/set-filter';
import {ExcelExportModule} from '@ag-grid-enterprise/excel-export';
import ImageRenderer from './image.vue';

export default {
  props: ['readOnly'],
  watch: {
    readOnly(value) {
      this.ReadOnly = value
    },
  },
  components: {
    'ag-grid-vue': AgGridVue, imageRenderer: ImageRenderer,
  },
  computed: {},
  data() {
    return {
      //是否被修改
      IsHasModify: false,
      //当前选中行
      agSelectedLine: null,
      overlayNoRowsTemplate: `<span style="padding: 20px;" id="HookMessageText">无内容</span>`,
      leftModules: [SetFilterModule, ClipboardModule],
      rightModules: [ExcelExportModule],
      gridOptions: {
        //rowSelection: 'multiple',
        onRangeSelectionChanged: this.onRangeSelectionChanged,
        onRowClicked: this.onRowClicked,
        onCellFocused: this.onCellFocused,
        getContextMenuItems: this.onContextMenuItems,
        getRowStyle: this.onGetRowStyle,
        onRowDataUpdated: this.NewColumnsLoaded,
        onModelUpdated: this.NewColumnsLoaded,
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
        filter: true,
        floatingFilter: false,
        resizable: true,
        menuTabs: ['filterMenuTab'],
        suppressNavigable: false,
        cellClass: 'no-border'
      },
      MenuItems: [],
      RowData: [],
      columns: [
        {
          field: "名称", tooltipField: '名称',
          minWidth: 200,
          width: 200,
          maxWidth: 200,
          // 不显示过滤器
          menuTabs: [],
          //禁止列拖动
          suppressMovable: true,
          cellStyle: {'text-align': 'left'},
        },
        {
          field: "深色主题", tooltipField: '深色主题',
          // 不显示过滤器
          menuTabs: [],
          //禁止列拖动
          suppressMovable: true,
          cellRenderer: 'imageRenderer',
          cellStyle: {'text-align': 'left'},
          minWidth: 80,
          width: 80,
          maxWidth: 80,
        }, {
          field: " ", tooltipField: ' ',
          // 不显示过滤器
          menuTabs: [],
          //禁止列拖动
          suppressMovable: true,
          cellStyle: {'text-align': 'left'},
          minWidth: 2,
          width: 2,
          maxWidth: 2,
        }, {
          field: "浅色主题", tooltipField: '浅色主题',
          // 不显示过滤器
          menuTabs: [],
          //禁止列拖动
          suppressMovable: true,
          cellRenderer: 'imageRenderer',
          cellStyle: {'text-align': 'left'},
          minWidth: 80,
          width: 80,
          maxWidth: 80,
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
      RequestId: {MessageId: -1, Theology: -1},
      ReadOnly: true,
    }
  }
  ,
  methods: {
    Load() {
      console.log( window.interface.colorConfig)
      this.RowData = [
        {
          "名称": "TCP",
          "深色主题": window.interface.colorConfig.tcp.dark,
          "浅色主题": window.interface.colorConfig.tcp.right,
          "说明": "如果是TCP请求",
          "alias": "tcp"
        },
        {
          "名称": "UDP",
          "深色主题": window.interface.colorConfig.udp.dark,
          "浅色主题": window.interface.colorConfig.udp.right,
          "说明": "如果是UDP请求",
          "alias": "udp"
        },
        {
          "名称": "Websocket",
          "深色主题": window.interface.colorConfig.ws.dark,
          "浅色主题": window.interface.colorConfig.ws.right,
          "说明": "如果是WS/WSS请求", "alias": "ws"
        },
        {
          "名称": "CSS",
          "深色主题": window.interface.colorConfig.css.dark,
          "浅色主题": window.interface.colorConfig.css.right,
          "说明": "如果响应是 CSS 文件 ", "alias": "css"
        },
        {
          "名称": "JS",
          "深色主题": window.interface.colorConfig.js.dark,
          "浅色主题": window.interface.colorConfig.js.right,
          "说明": "如果响应是 js 文件 ",
          "alias": "js"
        },
        {
          "名称": "图片",
          "深色主题": window.interface.colorConfig.img.dark,
          "浅色主题": window.interface.colorConfig.img.right,
          "说明": "如果响应是 图片 文件 ", "alias": "img"
        },
        {
          "名称": "文档",
          "深色主题": window.interface.colorConfig.document.dark,
          "浅色主题": window.interface.colorConfig.document.right,
          "说明": "如果响应是 文档类型 ", "alias": "document"
        },
        {
          "名称": "-1",
          "深色主题": window.interface.colorConfig._1.dark,
          "浅色主题": window.interface.colorConfig._1.right,
          "说明": "如果响应状态码是 -1 的请求", "alias": "_1"
        },
        {
          "名称": "301",
          "深色主题": window.interface.colorConfig._301.dark,
          "浅色主题": window.interface.colorConfig._301.right,
          "说明": "如果响应状态码是 301 的请求", "alias": "_301"
        },
        {
          "名称": "302",
          "深色主题": window.interface.colorConfig._302.dark,
          "浅色主题": window.interface.colorConfig._302.right,
          "说明": "如果响应状态码是 302 的请求", "alias": "_302"
        },
        {
          "名称": "401",
          "深色主题": window.interface.colorConfig._401.dark,
          "浅色主题": window.interface.colorConfig._401.right,
          "说明": "如果响应状态码是 401 的请求", "alias": "_401"
        },
        {
          "名称": "403",
          "深色主题": window.interface.colorConfig._403.dark,
          "浅色主题": window.interface.colorConfig._403.right,
          "说明": "如果响应状态码是 403 的请求", "alias": "_403"
        },
        {
          "名称": "404",
          "深色主题": window.interface.colorConfig._404.dark,
          "浅色主题": window.interface.colorConfig._404.right,
          "说明": "如果响应状态码是 404 的请求", "alias": "_404"
        },
        {
          "名称": "500",
          "深色主题": window.interface.colorConfig._500.dark,
          "浅色主题": window.interface.colorConfig._500.right,
          "说明": "如果响应状态码是 500 的请求", "alias": "_500"
        },
        {
          "名称": "",
          "深色主题": "",
          "浅色主题": "",
          "说明": "恢复到默认", "alias": "reset"
        },
      ]
    },
    Empty() {
      this.RowData = []
      this.agGridApi.setRowData(this.RowData);
      this.agSelectedLine = null
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
      params.node.setSelected(true);
      this.agSelectedLine = params.node
    },
    SelectedLine(index) {
      const focusedRowNode = this.agGridApi.getRowNode(index);
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
      this.SelectedLine(event.rowIndex)
    },
  },
  mounted() {
    this.agGridApi = this.$refs.agGrid.gridOptions.api
    window.vm.ListColorManager = this
  }
}
</script>

<style>

</style>