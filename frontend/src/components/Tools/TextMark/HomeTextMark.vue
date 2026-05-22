<template>
  <div style="height: 100%;width: 100%">
    <ag-grid-vue
        ref="agGrid"
        style="height: 90%;width: 100%;margin-bottom: 20px"
        :defaultColDef="defaultColDef"
        :rowData="RowData"
        :columnDefs="columns"
        :cellSelection="true"
        :enableCharts="true"
        :grid-options="gridOptions"
    >
    </ag-grid-vue>
    <el-text type="warning" v-show="isIsDark">你可以在主列表-选中请求后</el-text>
    <el-text type="danger" v-show="!isIsDark">你可以在主列表-选中请求后</el-text>
    <br>
    <br>
    <el-text type="warning" v-show="isIsDark">右键点击选中颜色标记，点击右键菜单中的 “对应的颜色即可”</el-text>
    <el-text type="danger" v-show="!isIsDark">右键点击选中颜色标记，点击右键菜单中的 “对应的颜色即可”</el-text>
  </div>
</template>

<script>
import {AgGridVue} from "ag-grid-vue3";
import {ExcelExportModule} from '@ag-grid-enterprise/excel-export';
import ColorConfigRenderer from './ColorConfig.vue';
import ColorConfigDarkRenderer from './ColorConfigDark.vue';
import {Config_HomeTextMark, Config_IsDark, DefaultRowData, getThisObject} from "../../config/Config";
import {SetHomeTextMark} from "../../../../bindings/changeme/Service/appmain";

export default {
  components: {
    'ag-grid-vue': AgGridVue, imageRenderer: ColorConfigRenderer, imageDarkRenderer: ColorConfigDarkRenderer
  },
  computed: {
    isIsDark() {
      return this.IsDark
    }
  },
  data() {
    return {
      get IsDark() {
        return Config_IsDark.value
      },
      set IsDark(value) {
        Config_IsDark.value = value
      },
      //当前选中行
      agSelectedRowNodes: [],
      overlayNoRowsTemplate: `<span style="padding: 20px;" id="HookMessageText">无内容</span>`,
      rightModules: [ExcelExportModule],
      gridOptions: {
        getRowId: (params) => params.data.id,
        rowSelection: {
          mode: 'singleRow',
          checkboxes: false,  // 启用行选择的复选框
        },
        getContextMenuItems: this.onContextMenuItems,
        onCellSelectionChanged: (params) => {
          if (params.started) {
            this.agSelectedRowNodes.forEach((node) => {
              node.setSelected(false);
            });
          }
          if (params.finished) {
            const rangeSelections = params.api.getCellRanges();
            const selectedRowNodes = [];
            let ___c = true;
            rangeSelections.forEach(function (range) {
              const _a = range.startRow.rowIndex;
              const _b = range.endRow.rowIndex;

              let startRow = _a;
              let endRow = _b;
              if (_a > _b) {
                startRow = _b
                endRow = _a
                ___c = false;
              } else {
                ___c = true;
              }
              params.api.forEachNodeAfterFilter(function (node) {
                if (node.rowIndex >= startRow && node.rowIndex <= endRow) {
                  selectedRowNodes.push(node);
                  node.setSelected(true);
                }
              });
            });
            this.agSelectedRowNodes = selectedRowNodes;
          }
        },
        // 禁用自动滚动到第一行
        suppressScrollOnNewData: true,
        onCellEditingStopped: (params) => {
          this.Save()
        },
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
      MenuItems: [
        {
          name: "新增一条",
          action: () => {
            this.addLine()
          },
        },
      ],
      RowData: [],
      columns: [
        {
          field: "名称", tooltipField: '名称',
          minWidth: 100,
          width: 100,
          maxWidth: 100,
          // 不显示过滤器
          menuTabs: [],
          //禁止列拖动
          editable: true,
          suppressMovable: true,
          cellStyle: {'text-align': 'left'},
        },
        {
          field: "深色主题", tooltipField: '深色主题',
          menuTabs: [],
          suppressMovable: true,
          cellRenderer: 'imageDarkRenderer',
          cellStyle: {'text-align': 'left'},
          minWidth: 120,
          width: 120,
          maxWidth: 120,
        },
        {
          field: "浅色主题", tooltipField: '浅色主题',
          menuTabs: [],
          suppressMovable: true,
          cellRenderer: 'imageRenderer',
          cellStyle: {'text-align': 'left'},
          minWidth: 120,
          width: 120,
          maxWidth: 120,
        },

      ],
    }
  },
  methods: {
    Save() {
      let array = []
      this.agGridApi.forEachNode((node) => {
        const obj = {
          id: node.data.id,
          "名称": node.data["名称"],
          "深色主题": node.data["深色主题"],
          "浅色主题": node.data["浅色主题"],
        }
        Config_HomeTextMark.set(node.data.id, obj)
        array.push(obj)
      })
      if (array.length < 1) {
        const array = [];
        for (let i = 0; i < DefaultRowData.length; i++) {
          const mm = DefaultRowData[i]
          mm.save = this.Save
          array.push(mm)
        }
        this.agGridApi.applyTransaction({add: array});
      }
      SetHomeTextMark(JSON.stringify(array))
      const func = getThisObject("HomeListRefresh")
      if (func) {
        func()
      }
    },
    addLine() {
      const data = [{
        "名称": "自定义",
        "深色主题": "#ff0707",
        "浅色主题": "#cc7c7c",
        id: (new Date().getTime() + ""),
        save: this.Save,
      }];
      this.agGridApi.applyTransaction({add: data});
      requestAnimationFrame(() => {
        this.Save()
      })
    },
    onContextMenuItems() {
      let array = [];
      for (let i = 0; i < this.MenuItems.length; i++) {
        if (this.MenuItems[i].visible !== false) {
          array.push(this.MenuItems[i])
        }
      }
      if (this.agSelectedRowNodes.length > 0) {
        array.push({
          name: "删除",
          action: () => {
            const arr = [];
            this.agSelectedRowNodes.forEach((node) => {
              arr.push(node.data)
            });
            this.agGridApi.applyTransaction({remove: arr});
            requestAnimationFrame(() => {
              this.Save()
            })
          },
        })
      }
      return array
    },
  },
  mounted() {
    this.agGridApi = this.$refs.agGrid.api;

    if (Config_HomeTextMark.size < 1) {

      const array = [];
      for (let i = 0; i < DefaultRowData.length; i++) {
        const mm = DefaultRowData[i]
        mm.save = this.Save
        array.push(mm)
      }
      for (let i = 0; i < array.length; i++) {
        Config_HomeTextMark.set(array[i].id, array[i])
      }
      this.agGridApi.applyTransaction({add: array});
    } else {
      const array = []
      Config_HomeTextMark.forEach((item) => {
        item.save = this.Save
        array.push(item)
      })
      this.agGridApi.applyTransaction({add: array});
    }
  }
}
</script>

<style>

</style>