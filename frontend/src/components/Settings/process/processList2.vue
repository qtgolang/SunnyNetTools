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
  <div v-show="LineSelected" :class="IsDarkTheme"/>
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
        this.columns[0].editable = false
      } else {
        this.columns[0].editable = true
      }
      //务必返回 false
      return false
    },
    IsDarkTheme() {
      const event = new Event('darkThemeChange');
      window.dispatchEvent(event);
      {
        if (this.darkTheme) {
          if (this.ListFollowShow) {
            this.MenuItems[0].icon = `<svg xmlns="http://www.w3.org/2000/svg" style="top: 2px;position: relative;" width="16" height="14" viewBox="0 0 24 24" fill="none" stroke="#fff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">\\n' +
              '    <polyline points="20 6 9 17 4 12"/>' +
              '</svg>`
          } else {
            this.MenuItems[0].icon = `<svg xmlns="http://www.w3.org/2000/svg" style="top: 2px;position: relative;" width="16" height="14" viewBox="0 0 24 24" fill="none" stroke="#fff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
    <line x1="18" y1="6" x2="6" y2="18"/>
    <line x1="6" y1="6" x2="18" y2="18"/>
</svg>
`
          }

        } else {
          if (this.ListFollowShow) {
            this.MenuItems[0].icon = `<svg xmlns="http://www.w3.org/2000/svg" style="top: 2px;position: relative;" width="16" height="14" viewBox="0 0 24 24" fill="none" stroke="#000" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">\\n' +
              '    <polyline points="20 6 9 17 4 12"/>' +
              '</svg>`
          } else {
            this.MenuItems[0].icon = `<svg xmlns="http://www.w3.org/2000/svg" style="top: 2px;position: relative;" width="16" height="14" viewBox="0 0 24 24" fill="none" stroke="#000" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
    <line x1="18" y1="6" x2="6" y2="18"/>
    <line x1="6" y1="6" x2="18" y2="18"/>
</svg>
`
          }
        }
      }
      return "example-wrapper " + (this.darkTheme ? "ag-theme-balham-dark" : "ag-theme-balham")
    }
  },
  data() {
    return {
      //是否被修改
      IsHasModify: false,
      //当前选中行
      agSelectedLine: null,
      ListFollowShow: true,
      RowDataHashMap: {},
      overlayNoRowsTemplate: `<span style="padding: 20px;" id="HookMessageText">无内容</span>`,
      leftModules: [SetFilterModule, ClipboardModule],
      rightModules: [ExcelExportModule],
      gridOptions: {
        getRowId: (params) => params.data.PID,
        rowSelection: 'multiple',
        onRangeSelectionChanged: this.onRangeSelectionChanged,
        onRowClicked: this.onRowClicked,
        onCellFocused: this.onCellFocused,
        getContextMenuItems: this.onContextMenuItems,
        getRowStyle: this.onGetRowStyle,
        onRowDataUpdated: this.NewColumnsLoaded,
        onModelUpdated: this.NewColumnsLoaded,
        onRowSelected: this.onRowSelected,
        suppressScrollOnNewData: true,
      },
      defaultColDef: {
        flex: 1,
        minWidth: 10,
        sortable: false,
        floatingFilter: false,
        resizable: true,
        filter:  true,
        suppressNavigable: false,
        cellClass: 'no-border'
      },
      MenuItems: [
        {
          name: '保证显示最后一行',
          action: () => {
            this.ListFollowShow = !this.ListFollowShow
          },
          disabled: false,
          visible: true,
          icon: ''
        },
      ],
      RowData: [],
      columns: [
        {
          checkboxSelection: true,
          headerCheckboxSelection: true,
          field: "PID", tooltipField: 'PID',
          minWidth: 110,
          maxWidth: 110,
          // 不显示过滤器
          menuTabs: ['filterMenuTab'],
          //禁止列拖动
          suppressMovable: true,
          editable: true,
          sortable: true,
          cellStyle: {'text-align': 'left'},
        },
        {
          field: "进程名", tooltipField: '进程名',
          minWidth: 80,
          // 不显示过滤器
          menuTabs: ['filterMenuTab'],
          //禁止列拖动
          suppressMovable: true,
          editable: true,
          sortable: true,
          cellStyle: {'text-align': 'left'},
        },
      ],
      RequestId: {MessageId: -1, Theology: -1},
      ReadOnly: false,
      get darkTheme() {
        return window.Theme.IsDark
      },
      set darkTheme(newValue) {
        window.Theme.IsDark = newValue
      },
    }
  }
  ,
  methods: {
    AddLines(objs) {
      const res = this.agGridApi.applyTransaction({add: objs});
      if (res.add) {
        res.add.forEach((rowNode) => {
          this.RowDataHashMap[rowNode.data.PID] = rowNode
        });
      }
      this.IsHasModify = false
    },
    Update(objs) {
      this.agGridApi.applyTransaction({update: objs});
      this.IsHasModify = false
    },
    Delete(objs) {
      this.agGridApi.applyTransaction({remove: objs});
      this.IsHasModify = false
    },
    Empty() {
      this.RowData = []
      this.agGridApi.setRowData(this.RowData);
      this.agSelectedLine = null
    },
    NewColumnsLoaded(params) {
      if (this.ListFollowShow) {
        const rowCount = this.agGridApi.getDisplayedRowCount() - 1
        if (rowCount > -1) {
          this.rowIndex = rowCount
          this.agGridApi.ensureIndexVisible(rowCount)
        }
      }
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
    onRowSelected(event) {
      const node = event.node;
      const gx = node.isSelected()
      CallGoDo("进程驱动添加PID", {PID: node.data['PID'], isSelected: gx})
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