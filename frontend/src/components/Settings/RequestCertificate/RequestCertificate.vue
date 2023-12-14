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
      @cellClicked="onCellClicked"
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
import RequestCertificateType from "./RequestCertificateType.vue";
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
    'ag-grid-vue': AgGridVue, requestCertificateType: RequestCertificateType,
  },
  computed: {
    LineSelected() {
      if (this.ReadOnly) {
        this.MenuItems[0].visible = false
        this.MenuItems[1].visible = false
        this.MenuItems[2].visible = false
        this.MenuItems[3].visible = false
        this.columns[0].editable = false
      } else {
        this.MenuItems[1].visible = this.agSelectedLine !== null;
        this.MenuItems[2].visible = this.agSelectedLine !== null;
        this.MenuItems[3].visible = this.agSelectedLine !== null;
        this.MenuItems[0].visible = true
        this.columns[0].editable = true
      }
      if (this.MenuItems[3].visible) {
        this.MenuItems[3].name = this.agSelectedLine.data['载入状态'] === "已载入" ? "重新载入该证书" : "载入该证书"
      }
      if (this.MenuItems[2].visible) {
        this.MenuItems[2].visible = this.agSelectedLine.data['载入状态'] === "已载入"
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
            CallGoDo("创建证书管理器", null).then(res => {
              this.$nextTick(() => {
                this.AddLine("解析及发送", res)
              })
            })
            this.IsHasModify = true
          },
          disabled: false,
          visible: true
        },
        {
          name: '删除',
          action: () => {
            if (this.agSelectedLine !== null) {
              const obj = {}
              obj.context = []
              obj.context[0] = {id: this.agSelectedLine.data.context}
              CallGoDo("删除证书管理器", obj).then(res => {
                this.RowData.splice(this.agSelectedLine.rowIndex, 1)
                this.$nextTick(() => {
                  this.agGridApi.setRowData(this.RowData);
                  this.agSelectedLine = null
                  this.IsHasModify = true
                })
              })
            }
          },
          disabled: false,
          visible: false
        },
        {
          name: '获取证书上的主机名',
          action: () => {
            if (this.agSelectedLine !== null) {
              const obj = {
                context: this.agSelectedLine.data.context,
                rule: this.agSelectedLine.data['使用规则'],
              }
              CallGoDo("查询证书CommonName", obj).then(res => {
                if (res === '') {
                  ElMessage({
                    message: "证书上的主机名为空",
                    type: 'error',
                  })
                  return
                }
                this.agSelectedLine.data['主机名'] = res
                this.agGridApi.setRowData(this.RowData);
              })
            }
          },
          disabled: false,
          visible: false
        },
        {
          name: '载入该证书',
          action: this.LoadCertificate,
          disabled: false,
          visible: true
        },
      ],
      RowData: [],
      columns: [
        {
          field: "使用规则", tooltipField: '使用规则',
          minWidth: 110,
          maxWidth: 110,
          // 不显示过滤器
          menuTabs: [],
          //禁止列拖动
          suppressMovable: true,
          editable: true,
          sortable: true,
          cellRenderer: 'requestCertificateType',
          cellStyle: {'text-align': 'left'},
        },
        {
          field: "主机名", tooltipField: '主机名',
          // 不显示过滤器
          menuTabs: [],
          minWidth: 200,
          maxWidth: 200,
          //禁止列拖动
          suppressMovable: true,
          editable: true,
          sortable: true,
          cellStyle: {'text-align': 'left'},
        },
        {
          field: "证书文件", tooltipField: '证书文件',
          // 不显示过滤器
          menuTabs: [],
          //禁止列拖动
          suppressMovable: true,
          sortable: true,
          cellStyle: {'text-align': 'left'},
        },
        {
          field: "密码", tooltipField: '密码',
          minWidth: 130,
          maxWidth: 130,
          // 不显示过滤器
          menuTabs: [],
          //禁止列拖动
          suppressMovable: true,
          editable: true,
          sortable: true,
          cellStyle: {'text-align': 'left'},
        },
        {
          field: "载入状态", tooltipField: '载入状态',
          minWidth: 130,
          maxWidth: 130,
          // 不显示过滤器
          menuTabs: [],
          //禁止列拖动
          suppressMovable: true,
          sortable: true,
          cellStyle: {'text-align': 'left'},
        },
      ],
      RequestId: {MessageId: -1, Theology: -1},
      ReadOnly: false,
    }
  }
  ,
  methods: {
    AddLine(name, value) {
      let obj = {}
      obj[this.columns[0].field] = name
      obj[this.columns[4].field] = "未载入"
      obj.context = value
      this.RowData.push(obj)
      this.$nextTick(() => {
        this.agGridApi.setRowData(this.RowData);
      })
      this.IsHasModify = false
    },
    AddLines(objs) {
      const array = []
      for (let i = 0; i < objs.length; i++) {
        const rule = objs[i]['rule']
        const Host = objs[i]['Host']
        const FilePath = objs[i]['FilePath']
        const PassWord = objs[i]['PassWord']
        const context = objs[i]['context']
        let obj = {}
        obj[this.columns[0].field] = rule === 3 ? "仅解析" : rule === 2 ? "解析及发送" : "仅发送"
        obj[this.columns[1].field] = Host
        obj[this.columns[2].field] = FilePath
        obj[this.columns[3].field] = PassWord
        obj[this.columns[4].field] = "未载入"
        obj.context = context
        obj.Hash = (new Date().getTime()) + "";
        array.push(obj)
      }
      this.RowData = array
      this.$nextTick(() => {
        this.agGridApi.setRowData(this.RowData);
      })
      this.IsHasModify = false
    },
    Empty() {
      this.RowData = []
      this.agGridApi.setRowData(this.RowData);
      this.agSelectedLine = null
    },
    LoadCertificate() {
      if (this.agSelectedLine !== null) {
        CallGoDo("载入请求证书", {Data: this.agSelectedLine.data}).then(res => {
          this.agSelectedLine.data['载入状态'] = res ? "已载入" : "载入失败"
          this.agGridApi.setRowData(this.RowData);
        })
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
    onCellClicked(params) {
      const clickedColumn = params.column.colId;
      if (clickedColumn === "证书文件") {
        const obj = {
          Title: "请选择证书文件",
          Filters: [
            {Name: "P12证书文件", Pattern: "*.p12;*.pkcs12;*.pfx"},
            {Name: "PEM证书", Pattern: "*.pem"},
            {Name: "CER证书", Pattern: "*.cer"}]
        }
        CallGoDo("选择文件", obj).then(res => {
          if (res !== '') {
            params.data['证书文件'] = res
            this.agGridApi.setRowData(this.RowData);
          }
        })
      }
      // 在这里处理单元格点击事件
      // params 包含有关点击的信息，如 params.rowIndex、params.colDef 和 params.value 等
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