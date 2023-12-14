<template>
  <div style="width: 100%;height: 100%">
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
  </div>
</template>

<script>
import '@ag-grid-community/styles/ag-grid.css';
import '@ag-grid-community/styles/ag-theme-balham.css';
import {AgGridVue} from '@ag-grid-community/vue3';
import {ClipboardModule} from '@ag-grid-enterprise/clipboard';
import {SetFilterModule} from '@ag-grid-enterprise/set-filter';
import {ExcelExportModule} from '@ag-grid-enterprise/excel-export';
import {CallGoDo} from "../../CallbackEventsOn.js";
import {ElMessage} from "element-plus";

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
        this.MenuItems[1].visible = false
        this.columns[0].editable = false
        this.columns[1].editable = false
      } else {
        this.MenuItems[1].visible = this.agSelectedLine !== null;
        this.MenuItems[2].visible = this.agSelectedLine !== null;
        this.MenuItems[0].visible = true
        this.columns[0].editable = true
        this.columns[1].editable = true
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
      overlayNoRowsTemplate: `<span style="padding: 20px;" id="HookMessageText">未设置任何账号密码</span>`,
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
          name: '增加',
          action: () => {
            this.AddLine("User", "Pass")
            this.IsHasModify = true
          },
          disabled: false,
          visible: true
        },
        {
          name: '删除',
          action: () => {
            if (this.agSelectedLine !== null) {
              this.RowData.splice(this.agSelectedLine.rowIndex, 1)
              this.agGridApi.setRowData(this.RowData);
              this.agSelectedLine = null
              this.IsHasModify = true
              this.Apply()
            }
          },
          disabled: false,
          visible: false
        },
        {
          name: '清空',
          action: this.Empty,
          disabled: false,
          visible: false
        },
        {
          name: '应用/使其生效',
          action: this.Apply,
          disabled: false,
          visible: true
        },
      ],
      RowData: [],
      columns: [
        {
          field: "账号", tooltipField: '账号',
          // 不显示过滤器
          menuTabs: [],
          //禁止列拖动
          suppressMovable: true,
          editable: true,
          sortable: true,
          cellStyle: {'text-align': 'left'},
        },
        {
          field: "密码", tooltipField: '密码',
          // 不显示过滤器
          menuTabs: [],
          //禁止列拖动
          suppressMovable: true,
          editable: true,
          cellStyle: {'text-align': 'left'}
        },
      ],
      ReadOnly: false,
    }
  }
  ,
  methods: {
    AddLine(name, value) {
      this.RowData.push({账号: name, 密码: value})
      this.agGridApi.setRowData(this.RowData);
      this.IsHasModify = false
    },
    Empty() {
      this.RowData = []
      this.agGridApi.setRowData(this.RowData);
      this.agSelectedLine = null
      this.Apply()
    },
    Apply() {
      let max = this.RowData.length
      for (let n = 0; n < max; n++) {
        max = this.RowData.length
        for (let i = 0; i < max; i++) {
          let user = ''
          let pass = ''
          if (this.RowData[i].账号) {
            user = this.RowData[i]['账号'].trim()
          }
          if (this.RowData[i].密码) {
            pass = this.RowData[i]['密码'].trim()
          }
          if (user === '' || pass === '') {
            this.RowData.splice(i, 1)
            break
          }
        }
      }
      CallGoDo("更新身份验证账号信息", {Data: this.RowData}).then(res => {
        if (res) {
          ElMessage({
            message: "身份验证账号信息,已更新",
            type: 'success',
          })
        }
        this.agGridApi.setRowData(this.RowData);
        this.agSelectedLine = null
      })
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