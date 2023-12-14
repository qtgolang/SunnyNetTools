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
        this.columns[2].editable = false
      } else {
        this.MenuItems[1].visible = this.agSelectedLine !== null;
        this.MenuItems[0].visible = true
        this.columns[0].editable = true
        this.columns[1].editable = true
        this.columns[2].editable = true
      }
      //务必返回 false
      return false
    }
  },
  data() {
    return {
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
      MenuItems: [
        {
          name: '增加',
          action: () => {
            this.AddLine("Header" + (this.RowData.length + 1), "Header Value")
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
            }
          },
          disabled: false,
          visible: false
        },
      ],
      RowData: [],
      columns: [
        {
          field: "名称", tooltipField: '名称',
          minWidth: 80,
          // 不显示过滤器
          menuTabs: [],
          //禁止列拖动
          suppressMovable: true,
          editable: true,
          sortable: true,
          cellStyle: {'text-align': 'left'},
        },
        {
          field: "值", tooltipField: '值',
          minWidth: 80,
          // 不显示过滤器
          menuTabs: [],
          //禁止列拖动
          suppressMovable: true,
          editable: true,
          cellStyle: {'text-align': 'left'}
        },
        {
          field: "其他值", tooltipField: '其他值',
          minWidth: 80,
          // 不显示过滤器
          menuTabs: [],
          //禁止列拖动
          suppressMovable: true,
          editable: true,
          cellStyle: {'text-align': 'left'}
        },
      ],
      RequestId: {MessageId: -1, Theology: -1},
      ReadOnly: true,
    }
  }
  ,
  methods: {
    AddLine(name, value, v2) {
      this.RowData.push({名称: name, 值: value, 其他值: v2})
      this.agGridApi.setRowData(this.RowData);
      this.IsHasModify = false
    },
    GetLine() {
      return this.RowData.length
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
  }
}
</script>

<style>

</style>