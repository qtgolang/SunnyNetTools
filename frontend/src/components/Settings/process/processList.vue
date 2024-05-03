<template>
  <ag-grid-vue
      ref="agGrid"
      style="height: 100%;"
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
  <div v-show="LineSelected"/>
</template>

<script>
import '@ag-grid-community/styles/ag-grid.css';
import '@ag-grid-community/styles/ag-theme-balham.css';
import {AgGridVue} from '@ag-grid-community/vue3';
import {ClipboardModule} from '@ag-grid-enterprise/clipboard';
import {SetFilterModule} from '@ag-grid-enterprise/set-filter';
import {ExcelExportModule} from '@ag-grid-enterprise/excel-export';
import {CallGoDo} from "../../CallbackEventsOn.js";

export default {
  props: ['readOnly'],
  watch: {
    readOnly(value) {
      this.ReadOnly = value
    },
  },
  components: {
    'ag-grid-vue': AgGridVue,
  },
  computed: {
    LineSelected() {
      if (this.ReadOnly) {
        this.MenuItems[0].visible = false
        //this.columns[0].editable = false
      } else {
        this.MenuItems[0].visible = this.agSelectedLine !== null;
        //this.columns[0].editable = true
      }
      //务必返回 false
      return false
    }
  },
  data() {
    return {
      //是否被修改
      IsHasModify: false,
      //当前选中行
      agSelectedLine: null,
      overlayNoRowsTemplate: `<span style="padding: 20px;" id="HookMessageText">未添加任何进程名</span>`,
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
      MenuItems: [
        {
          name: '删除',
          action: () => {
            if (this.agSelectedLine !== null) {
              this.RowData.splice(this.agSelectedLine.rowIndex, 1)
              this.agGridApi.setRowData(this.RowData);
              this.agSelectedLine = null
              this.IsHasModify = true
            }
          },
          disabled: false,
          visible: false
        },
        "separator",
        {
          name: '一键添加微信相关',
          action: () => {
            this.AddProcessName("WeChat.exe")
            this.AddProcessName("wechatweb.exe")
            this.AddProcessName("WechatAppLauncher.exe")
            this.AddProcessName("WeChatAppEx.exe")
            this.AddProcessName("WechatBrowser.exe")
            this.AddProcessName("WeChatPlayer.exe")
            this.AddProcessName("WeChatXFile.exe")
          },
          disabled: false,
          visible: true
        },

        {
          name: '一键添加雷电模拟器相关',
          action: () => {
            this.AddProcessName("dnplayer.exe")
            this.AddProcessName("LdBoxHeadless.exe")
            this.AddProcessName("LdVBoxHeadless.exe")
            this.AddProcessName("Ld9BoxHeadless.exe")
          },
          disabled: false,
          visible: true
        },
      ],
      RowData: [],
      columns: [
        {
          field: "进程名", tooltipField: '进程名',
          minWidth: 80,
          // 不显示过滤器
          menuTabs: [],
          //禁止列拖动
          suppressMovable: true,
          //editable: true,
          cellStyle: {'text-align': 'left'},
        },
      ],
      RequestId: {MessageId: -1, Theology: -1},
      ReadOnly: false,
    }
  }
  ,
  methods: {
    AddProcessName(name) {
      for (let i = 0; i < this.RowData.length; i++) {
        const n = this.RowData[i]['进程名'].toLowerCase()
        const ne = name.toLowerCase()
        if (n === ne) {
          return
        }
      }
      CallGoDo("进程驱动添加进程名", {Name: name, isSet: true}).then(res => {
        this.RowData.push({'进程名': name})
        this.agGridApi.setRowData(this.RowData);
      })
    },
    AddLine(name, value) {
      let obj = {}
      obj[this.columns[0].field] = name
      this.RowData.push(obj)
      this.agGridApi.setRowData(this.RowData);
      this.IsHasModify = false
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
    handleDocumentClick(event) {
      try {
        const editingCells = this.agGridApi.getEditingCells();
        if (editingCells.length > 0) {
          if (event.target) {
            const mm = event.target
            const ma = this.$refs.agGrid.$el.getElementsByClassName("ag-input-field-input ag-text-field-input")
            for (let i = 0; i < ma.length; i++) {
              if (ma[i] === mm) {
                return
              }
            }
          }
          this.agGridApi.stopEditing(); // 退出单元格编辑
        }
      } catch (e) {
        console.log(e)
      }
    },
  },
  mounted() {
    this.agGridApi = this.$refs.agGrid.gridOptions.api
    document.addEventListener('click', this.handleDocumentClick);
  }
}
</script>

<style>

</style>