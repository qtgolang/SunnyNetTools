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
import ImageRenderer from './SocketImage.vue';
import {CallGoDo} from "../CallbackEventsOn.js";

export default {
  props: ['readOnly', 'width'],
  watch: {
    readOnly(value) {
      this.ReadOnly = value
    },
    width(value) {
      this._Width = parseInt(value.replaceAll("px", ""))
    },
  },
  components: {
    'ag-grid-vue': AgGridVue, imageRenderer: ImageRenderer,
  },
  computed: {
    LineSelected() {
      this.MenuItems[6].subMenu[0].disabled = this.agSelectedLine === null
      if (this.darkTheme) {
        for (let i = 0; i < this.MenuItems.length; i++) {
          if (typeof this.MenuItems[i] === 'string') {
            continue
          }
          if (this.MenuItems[i].selected) {
            this.MenuItems[i].icon = `<svg xmlns="http://www.w3.org/2000/svg" style="top: 2px;position: relative;" width="16" height="14" viewBox="0 0 24 24" fill="none" stroke="#fff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">\\n' +
              '    <polyline points="20 6 9 17 4 12"/>' +
              '</svg>`
          } else {
            this.MenuItems[i].icon = ""
          }
        }
      } else {
        for (let i = 0; i < this.MenuItems.length; i++) {
          if (typeof this.MenuItems[i] === 'string') {
            continue
          }
          if (this.MenuItems[i].selected) {
            this.MenuItems[i].icon = `<svg xmlns="http://www.w3.org/2000/svg" style="top: 2px;position: relative;" width="16" height="14" viewBox="0 0 24 24" fill="none" stroke="#000" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">\\n' +
              '    <polyline points="20 6 9 17 4 12"/>' +
              '</svg>`
          } else {
            this.MenuItems[i].icon = ""
          }
        }
      }
      window.Socket.Line = this.agSelectedLine
      //务必返回 false
      return false
    },
  },
  data() {
    return {
      _Width: 0,
      //是否被修改
      IsHasModify: false,
      //当前选中行
      agSelectedLine: null,
      overlayNoRowsTemplate: `<span style="padding: 20px;" id="HookMessageText">暂无发送、接收的数据</span>`,
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
          name: '跟随显示',
          action: () => {
            this.MenuItems[0].selected = !this.MenuItems[0].selected
          },
          //禁止点击
          disabled: false,
          selected: true
        },
        'separator',
        {
          name: '停止插入发送',
          action: () => {
            this.MenuItems[2].selected = !this.MenuItems[2].selected
            if (this.MenuItems[2].selected) {
              this.MenuItems[4].selected = false
            }
            let way = window.vm.List.agSelectedLine.data["方式"].toUpperCase()
            CallGoDo("设置右键菜单配置", {
              StopSend: this.MenuItems[2].selected,
              StopRec: this.MenuItems[3].selected,
              StopALL: this.MenuItems[4].selected,
              Theology: window.Theology,
              IsWs: window.vm.Tabs.Request.DisplayHTTPHeader,
              IsTCP: way.indexOf("TCP") !== -1
            })
          },
          //禁止点击
          disabled: false,
          selected: false
        },
        {
          name: '停止插入接收',
          action: () => {
            this.MenuItems[3].selected = !this.MenuItems[3].selected
            if (this.MenuItems[3].selected) {
              this.MenuItems[4].selected = false
            }
            let way = window.vm.List.agSelectedLine.data["方式"].toUpperCase()
            CallGoDo("设置右键菜单配置", {
              StopSend: this.MenuItems[2].selected,
              StopRec: this.MenuItems[3].selected,
              StopALL: this.MenuItems[4].selected,
              Theology: window.Theology,
              IsWs: window.vm.Tabs.Request.DisplayHTTPHeader,
              IsTCP: way.indexOf("TCP") !== -1
            })
          },
          //禁止点击
          disabled: false,
          selected: false
        },
        {
          name: '全部停止插入',
          action: () => {
            this.MenuItems[4].selected = !this.MenuItems[4].selected
            if (this.MenuItems[4].selected) {
              this.MenuItems[2].selected = false
              this.MenuItems[3].selected = false
            }
            let way = window.vm.List.agSelectedLine.data["方式"].toUpperCase()
            CallGoDo("设置右键菜单配置", {
              StopSend: this.MenuItems[2].selected,
              StopRec: this.MenuItems[3].selected,
              StopALL: this.MenuItems[4].selected,
              Theology: window.Theology,
              IsWs: window.vm.Tabs.Request.DisplayHTTPHeader,
              IsTCP: way.indexOf("TCP") !== -1
            })
          },
          //禁止点击
          disabled: false,
          selected: false
        },
        'separator',
        {
          name: '复制',
          subMenu: [
            {
              name: '复制选中HEX到剪辑版',
              action: () => {
                let way = window.vm.List.agSelectedLine.data["方式"].toUpperCase()
                CallGoDo("socket点击右键菜单", {
                  Type: "Selected",
                  Theology: window.Theology,
                  IsWs: window.vm.Tabs.Request.DisplayHTTPHeader,
                  IsTCP: way.indexOf("TCP") !== -1,
                  SelectedID: this.agSelectedLine.data["#"]
                })
              },
              disabled: false,
              visible: true
            },
            {
              name: '复制所有HEX到剪辑版',
              action: () => {
                let way = window.vm.List.agSelectedLine.data["方式"].toUpperCase()
                CallGoDo("socket点击右键菜单", {
                  Type: "AllHEX",
                  Theology: window.Theology,
                  IsWs: window.vm.Tabs.Request.DisplayHTTPHeader,
                  IsTCP: way.indexOf("TCP") !== -1,
                })
              },
              disabled: false,
              visible: true
            },
            {
              name: '复制所有发送数据HEX到剪辑版',
              action: () => {
                let way = window.vm.List.agSelectedLine.data["方式"].toUpperCase()
                CallGoDo("socket点击右键菜单", {
                  Type: "sendHEX",
                  Theology: window.Theology,
                  IsWs: window.vm.Tabs.Request.DisplayHTTPHeader,
                  IsTCP: way.indexOf("TCP") !== -1,
                })
              },
              disabled: false,
              visible: true
            },
            {
              name: '复制所有接收数据HEX到剪辑版',
              action: () => {
                let way = window.vm.List.agSelectedLine.data["方式"].toUpperCase()
                CallGoDo("socket点击右键菜单", {
                  Type: "recHEX",
                  Theology: window.Theology,
                  IsWs: window.vm.Tabs.Request.DisplayHTTPHeader,
                  IsTCP: way.indexOf("TCP") !== -1,
                })
              },
              disabled: false,
              visible: true
            },
          ],
          disabled: false,
          visible: true
        },
        'separator',
        {
          name: '清空',
          action: () => {
            let way = window.vm.List.agSelectedLine.data["方式"].toUpperCase()
            CallGoDo("socket点击右键菜单", {
              Type: "empty",
              Theology: window.Theology,
              IsWs: window.vm.Tabs.Request.DisplayHTTPHeader,
              IsTCP: way.indexOf("TCP") !== -1
            }).then(res => {
              if (res) {
                this.Empty()
              }
            })
          },
          //禁止点击
          disabled: false,
          selected: false
        },
      ],
      RowData: [],
      columns: [
        {
          field: "#", tooltipField: '#',
          minWidth: 80,
          maxWidth: 80,
          // 不显示过滤器
          menuTabs: [],
          //禁止列拖动
          suppressMovable: true, cellRenderer: 'imageRenderer',
          cellStyle: {'text-align': 'left'},
        },
        {
          field: "时间", tooltipField: '时间',
          minWidth: 97,
          maxWidth: 97,
          // 不显示过滤器
          menuTabs: [],
          //禁止列拖动
          suppressMovable: true,
        },
        {
          field: "类型", tooltipField: '类型',
          minWidth: 60,
          maxWidth: 60,
          // 不显示过滤器
          menuTabs: [],
          //禁止列拖动
          suppressMovable: true,
          cellStyle: {'text-align': 'left'}
        },
        {
          field: "长度", tooltipField: '长度',
          minWidth: 60,
          maxWidth: 60,
          // 不显示过滤器
          menuTabs: [],
          //禁止列拖动
          suppressMovable: true,
          cellStyle: {'text-align': 'left'}
        },
        {
          field: "数据", tooltipField: '数据',
          // 不显示过滤器
          menuTabs: [],
          //禁止列拖动
          suppressMovable: true,
          cellStyle: {'text-align': 'left'}
        },
      ],
      RequestId: {MessageId: -1, Theology: -1},
      ReadOnly: true,
      rowIndex: 0,
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
      this.agGridApi.applyTransaction({add: objs});
    },
    RefreshRenderedNodes() {
      const visibleRowNodes = this.agGridApi.getRenderedNodes();
      let array = []
      visibleRowNodes.forEach(node => {
        array.push(node)
      });
      this.$nextTick(() => {
        //this.agGridApi.refreshCells({rowNodes: array, force: true})
        this.agGridApi.redrawRows({rowNodes: array});
      })
    },
    onGetRowStyle(params) {
      let res = {
        fontFamily: "微软雅黑"
      }
      if (params.data.background) {
        res.backgroundColor = params.data.background
      }
      return res
    },
    SetColumnsMode(ws) {
      if (ws) {
        this.columns = [
          {
            field: "#", tooltipField: '#',
            minWidth: 80,
            maxWidth: 80,
            // 不显示过滤器
            menuTabs: [],
            //禁止列拖动
            suppressMovable: true, cellRenderer: 'imageRenderer',
            cellStyle: {'text-align': 'left'},
          },
          {
            field: "时间", tooltipField: '时间',
            minWidth: 97,
            maxWidth: 97,
            // 不显示过滤器
            menuTabs: [],
            //禁止列拖动
            suppressMovable: true,
          },
          {
            field: "类型", tooltipField: '类型',
            minWidth: 60,
            maxWidth: 60,
            // 不显示过滤器
            menuTabs: [],
            //禁止列拖动
            suppressMovable: true,
            cellStyle: {'text-align': 'left'}
          },
          {
            field: "长度", tooltipField: '长度',
            minWidth: 60,
            maxWidth: 60,
            // 不显示过滤器
            menuTabs: [],
            //禁止列拖动
            suppressMovable: true,
            cellStyle: {'text-align': 'left'}
          },
          {
            field: "数据", tooltipField: '数据',
            // 不显示过滤器
            menuTabs: [],
            //禁止列拖动
            suppressMovable: true,
            cellStyle: {'text-align': 'left'}
          },
        ]
      } else {
        this.columns = [
          {
            field: "#", tooltipField: '#',
            minWidth: 80,
            maxWidth: 80,
            // 不显示过滤器
            menuTabs: [],
            //禁止列拖动
            suppressMovable: true, cellRenderer: 'imageRenderer',
            cellStyle: {'text-align': 'left'},
          },
          {
            field: "时间", tooltipField: '时间',
            minWidth: 97,
            maxWidth: 97,
            // 不显示过滤器
            menuTabs: [],
            //禁止列拖动
            suppressMovable: true,
          },
          {
            field: "长度", tooltipField: '长度',
            minWidth: 60,
            maxWidth: 60,
            // 不显示过滤器
            menuTabs: [],
            //禁止列拖动
            suppressMovable: true,
            cellStyle: {'text-align': 'left'}
          },
          {
            field: "数据", tooltipField: '数据',
            // 不显示过滤器
            menuTabs: [],
            //禁止列拖动
            suppressMovable: true,
            cellStyle: {'text-align': 'left'}
          },
        ]
      }
    },
    Refresh() {
      setTimeout(() => {
        this.$nextTick(() => {
          this.agGridApi.applyTransaction({add: []});
        })
      }, 500)
    },
    Empty() {
      this.RowData = []
      this.agGridApi.setRowData(this.RowData);
      this.agSelectedLine = null
    },
    onContextMenuItems() {
      let array = [];
      for (let i = 0; i < this.MenuItems.length; i++) {
        array.push(this.MenuItems[i])
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
    NewColumnsLoaded(params) {
      if (this.MenuItems[0].selected) {
        const rowCount = this.agGridApi.getDisplayedRowCount() - 1
        if (rowCount > -1) {
          this.rowIndex = rowCount
          this.agGridApi.ensureIndexVisible(rowCount)
        }
      }
    }
    ,
  },
  mounted() {
    this.agGridApi = this.$refs.agGrid.gridOptions.api

  }
}
</script>

<style>

</style>