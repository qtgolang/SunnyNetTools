<script>
import {AgGridVue} from "ag-grid-vue3";
import {Config_IsDark, Config_Theme_agGrid, deleteThisObject} from "../../config/Config.js";
import {AG_GRID_LOCALE_CN} from "../../config/AG_ZH_CN.js";
import {Events} from "@wailsio/runtime";
import ReplaceType from "./ReplaceType.vue";
import ReplaceSourceType from "./ReplaceSourceType.vue";

import {ElNotification} from "element-plus";
import {
  CreateReplaceBody,
  ReplaceBodyList,
  ReplaceBodyRemove,
  ReplaceBodyUpdate
} from "../../../../bindings/changeme/Service/appmain";
import TitleBar from "../../TitleBar/TitleBar.vue";
import {attachMcpConfigReload} from "../../config/mcpRulesSync.js";

export default {
  components: {TitleBar, 'ag-grid-vue': AgGridVue, "replaceType": ReplaceType, "replaceSourceType": ReplaceSourceType},
  data() {
    return {
      Stopped: null,
      agGridApi: null,
      rowData: [],
      addValue: null,
      overlayNoRowsTemplate: `<span style="padding: 20px;" id="HookMessageText">你还没有添加任何数据替换/拦截</span>`,
      defaultColDef: {
        flex: 1,
        sortable: false,
        suppressHeaderMenuButton: true,
        suppressHeaderContextMenu: true
      },
      ReadOnly: false,
      _rulesReloading: false,
      _rulesReloadPromise: null,
      RowNodes: [],
      CellForClipboard: null,
      gridOptions: {
        cellSelection: true,
        suppressMovableColumns: true,
        stopEditingWhenCellsLoseFocus: true, // 失去焦点时自动结束编辑
        getRowId: (params) => params.data.id,
        getContextMenuItems: this.MenuEvent,
        localeText: AG_GRID_LOCALE_CN,
        columnDefs: [
          {
            field: "替换类型", tooltipField: '替换类型',
            minWidth: 110,
            maxWidth: 110,
            cellRenderer: 'replaceType',
            cellStyle: {'text-align': 'left'},
          },
          {
            field: "查找范围", tooltipField: '查找范围',
            minWidth: 145,
            maxWidth: 145,
            cellRenderer: 'replaceSourceType',
            cellStyle: {'text-align': 'left'},
          },
          {
            field: "旧数据", tooltipField: '旧数据',
            minWidth: 170,
            width: 170,
            editable: true,
          },
          {
            field: "新数据", tooltipField: '新数据',
            minWidth: 170,
            width: 170,
            editable: true,
          },
          {
            field: "状态", tooltipField: '状态',
            minWidth: 70,
            width: 70,
            maxWidth: 70,
          },
          {
            field: "注释", tooltipField: '注释',
            minWidth: 50,
            width: 200,
            editable: true,
          },
        ],
        rowSelection: {
          mode: 'multiRow',
          checkboxes: false,
          headerCheckbox: false,
          copySelectedRows: false,
          enableClickSelection: true,
        },
        onCellDoubleClicked: (event) => {
          const node = event.data;
          if (event.colDef.field !== "状态") {
            return;
          }
          const state = event.value === "已启用" ? "已禁用" : "已启用"
          ReplaceBodyUpdate(parseInt(node.id), node["替换类型"], node["查找范围"], node["旧数据"], node["新数据"], node["注释"], state).then(res => {
            if (res && (res.Ok ?? true)) {
              node["状态"] = res.State ?? "已启用";
              ElNotification({
                showClose: true,
                message: '请求拦截/替换规则已更新',
                type: 'success',
                position: 'bottom-right',
              })
            } else {
              node["状态"] = "已禁用"
              ElNotification({
                position: 'bottom-right',
                showClose: true,
                message: '请求拦截/替换规则更新失败,请检查数据格式是否正确?',
                type: 'warning',
              })
            }
            this.agGridApi.applyTransaction({update: [node]});
          })
        },
        onCellSelectionChanged: (params) => {
          if (params.started) {
            this.RowNodes.forEach((node) => {
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
            this.RowNodes = selectedRowNodes;
          }
        },
        onCellEditingStarted: (params) => {
        },
        onCellEditingStopped: (params) => {
          this.loadReplace(params)
        },
      }
    }
  },
  computed: {
    agTheme() {
      return Config_Theme_agGrid.value
    },
  },
  mounted() {
    this.agGridApi = this.$refs.agGrid.api;
    attachMcpConfigReload("all", () => this.reloadRulesFromBackend());
    Events.On("SetIsDark", (obj) => {
      try {
        const dark=obj.data[0]===true || obj.data===true;
        if (Config_IsDark.value === dark){
          return
        }
        const htmlElement = document.documentElement;
        if (dark) {
          htmlElement.setAttribute('data-dark-mode', 'true');
          htmlElement.setAttribute('data-ag-theme-mode', 'dark-blue');
        } else {
          htmlElement.setAttribute('data-dark-mode', '');
          htmlElement.setAttribute('data-ag-theme-mode', '');
        }
        Config_IsDark.value = dark
      } catch (e) {

      }
    })
    Events.On("onRest", (obj) => {
      this.reloadRulesFromBackend();
    })
    this.reloadRulesFromBackend();
  },
  methods: {
    reloadRulesFromBackend() {
      if (!this.agGridApi) {
        return Promise.resolve();
      }
      if (this._rulesReloading) {
        return this._rulesReloadPromise ?? Promise.resolve();
      }
      this._rulesReloading = true;
      this._rulesReloadPromise = ReplaceBodyList().then((list) => {
        const array = [];
        for (let i = 0; i < list.length; i++) {
          const item = list[i];
          array.push({
            "替换类型": item["Type"],
            "查找范围": item["Source"],
            "旧数据": item["Lod"],
            "新数据": item["New"],
            "注释": item["Note"],
            "状态": (item["State"] && String(item["State"]).trim() !== "")
                ? item["State"]
                : (item["Ok"] ? "已启用" : "已禁用"),
            id: item["ID"] + "",
          });
        }
        this.agGridApi.setGridOption("rowData", array);
      }).finally(() => {
        this._rulesReloading = false;
        this._rulesReloadPromise = null;
      });
      return this._rulesReloadPromise;
    },
    Empty() {
      this.agGridApi.setGridOption("rowData", []);
    },
    addLine(type, source, lod, New, node) {
      CreateReplaceBody().then(id => {
        const data = [{
          "替换类型": type,
          "查找范围": source,
          "旧数据": lod,
          "新数据": New,
          "状态": "已禁用",
          id: id + "",
          "注释": node
        }];
        this.agGridApi.applyTransaction({add: data});
      })
    },
    SetReadOnly(readOnly) {
      this.ReadOnly = readOnly;
      this.gridOptions.columnDefs.forEach((obj) => {
        obj.editable = !readOnly
      })
      this.agGridApi.setGridOption('columnDefs', this.gridOptions.columnDefs);
    },
    loadReplace(params, noMessage) {
      const node = params.node.data;
      ReplaceBodyUpdate(parseInt(node.id), node["替换类型"], node["查找范围"], node["旧数据"], node["新数据"], node["注释"], "已启用").then(res => {
        if (res) {
          node["状态"] = "已启用"
          ElNotification({
            showClose: true,
            message: '请求拦截/替换规则已更新',
            type: 'success',
            position: 'bottom-right',
          })
        } else {
          node["状态"] = "已禁用"
          ElNotification({
            position: 'bottom-right',
            showClose: true,
            message: '请求拦截/替换规则更新失败,请检查数据格式是否正确?',
            type: 'warning',
          })
        }
        this.agGridApi.applyTransaction({update: [node]});
      })
    },
    MenuEvent(params) {
      const filteredMenuItems = []
      if (!this.ReadOnly) {
        if (this.RowNodes.length > 0) {
          //新增删除菜单
          {
            let delIcon = '';
            if (!Config_IsDark.value) {
              delIcon = '<div style="display: flex; align-items: center;width: 16px">' +
                  '<svg clip-rule="evenodd" fill-rule="evenodd" stroke-linejoin="round" stroke-miterlimit="2" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path d="m12.002 2.005c5.518 0 9.998 4.48 9.998 9.997 0 5.518-4.48 9.998-9.998 9.998-5.517 0-9.997-4.48-9.997-9.998 0-5.517 4.48-9.997 9.997-9.997zm0 1.5c-4.69 0-8.497 3.807-8.497 8.497s3.807 8.498 8.497 8.498 8.498-3.808 8.498-8.498-3.808-8.497-8.498-8.497zm0 7.425 2.717-2.718c.146-.146.339-.219.531-.219.404 0 .75.325.75.75 0 .193-.073.384-.219.531l-2.717 2.717 2.727 2.728c.147.147.22.339.22.531 0 .427-.349.75-.75.75-.192 0-.384-.073-.53-.219l-2.729-2.728-2.728 2.728c-.146.146-.338.219-.53.219-.401 0-.751-.323-.751-.75 0-.192.073-.384.22-.531l2.728-2.728-2.722-2.722c-.146-.147-.219-.338-.219-.531 0-.425.346-.749.75-.749.192 0 .385.073.531.219z" fill-rule="nonzero"/></svg>' +
                  '</div>';
            } else {
              delIcon = '<div  class="white-svg" class="white-svg" style="display: flex; align-items: center;width: 16px">' +
                  '<svg clip-rule="evenodd" fill-rule="evenodd" stroke-linejoin="round" stroke-miterlimit="2" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path d="m12.002 2.005c5.518 0 9.998 4.48 9.998 9.997 0 5.518-4.48 9.998-9.998 9.998-5.517 0-9.997-4.48-9.997-9.998 0-5.517 4.48-9.997 9.997-9.997zm0 1.5c-4.69 0-8.497 3.807-8.497 8.497s3.807 8.498 8.497 8.498 8.498-3.808 8.498-8.498-3.808-8.497-8.498-8.497zm0 7.425 2.717-2.718c.146-.146.339-.219.531-.219.404 0 .75.325.75.75 0 .193-.073.384-.219.531l-2.717 2.717 2.727 2.728c.147.147.22.339.22.531 0 .427-.349.75-.75.75-.192 0-.384-.073-.53-.219l-2.729-2.728-2.728 2.728c-.146.146-.338.219-.53.219-.401 0-.751-.323-.751-.75 0-.192.073-.384.22-.531l2.728-2.728-2.722-2.722c-.146-.147-.219-.338-.219-.531 0-.425.346-.749.75-.749.192 0 .385.073.531.219z" fill-rule="nonzero"/></svg>' +
                  '</div>';
            }
            filteredMenuItems.push({
              name: "删除选中",
              action: async () => {
                const selectedRowNodes = [];
                this.RowNodes.forEach(function (range) {
                  selectedRowNodes.push(range.data)
                });
                for (let i = 0; i < selectedRowNodes.length; i++) {
                  await ReplaceBodyRemove(parseInt(selectedRowNodes[i].id))
                  deleteThisObject(selectedRowNodes[i].id + "|查找类型")
                }
                this.agGridApi.applyTransaction({remove: selectedRowNodes});
                this.RowNodes = [];
                this.agGridApi.clearCellSelection()
              },
              icon: delIcon,
            });
          }
        }
        //新增添加菜单
        {
          let addIcon = '';
          if (!Config_IsDark.value) {
            addIcon = '<div style="display: flex; align-items: center;width: 16px">' +
                '<svg clip-rule="evenodd" fill-rule="evenodd" stroke-linejoin="round" stroke-miterlimit="2" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path d="m12.002 2c5.518 0 9.998 4.48 9.998 9.998 0 5.517-4.48 9.997-9.998 9.997-5.517 0-9.997-4.48-9.997-9.997 0-5.518 4.48-9.998 9.997-9.998zm0 1.5c-4.69 0-8.497 3.808-8.497 8.498s3.807 8.497 8.497 8.497 8.498-3.807 8.498-8.497-3.808-8.498-8.498-8.498zm-.747 7.75h-3.5c-.414 0-.75.336-.75.75s.336.75.75.75h3.5v3.5c0 .414.336.75.75.75s.75-.336.75-.75v-3.5h3.5c.414 0 .75-.336.75-.75s-.336-.75-.75-.75h-3.5v-3.5c0-.414-.336-.75-.75-.75s-.75.336-.75.75z" fill-rule="nonzero"/></svg>' +
                '</div>';
          } else {
            addIcon = '<div class="white-svg" style="display: flex; align-items: center;width: 15px">' +
                '<svg clip-rule="evenodd" fill-rule="evenodd" stroke-linejoin="round" stroke-miterlimit="2" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path d="m12.002 2c5.518 0 9.998 4.48 9.998 9.998 0 5.517-4.48 9.997-9.998 9.997-5.517 0-9.997-4.48-9.997-9.997 0-5.518 4.48-9.998 9.997-9.998zm0 1.5c-4.69 0-8.497 3.808-8.497 8.498s3.807 8.497 8.497 8.497 8.498-3.807 8.498-8.497-3.808-8.498-8.498-8.498zm-.747 7.75h-3.5c-.414 0-.75.336-.75.75s.336.75.75.75h3.5v3.5c0 .414.336.75.75.75s.75-.336.75-.75v-3.5h3.5c.414 0 .75-.336.75-.75s-.336-.75-.75-.75h-3.5v-3.5c0-.414-.336-.75-.75-.75s-.75.336-.75.75z" fill-rule="nonzero"/></svg>' +
                '</div>';
          }
          filteredMenuItems.push({
            name: "新增一条",
            action: () => {
              this.addLine("字符串(UTF8)", "任意", "[双击修改-旧数据]", "[双击修改-新数据]", "[双击修改-注释]")
            },
            icon: addIcon,
          });
        }
      }
      return filteredMenuItems;
    },
  }
}
</script>

<template>
  <div class="fullscreen-div" style="display: block;">
    <TitleBar Title="请求拦截/数据替换设置"></TitleBar>
    <ag-grid-vue ref="agGrid"
                 :theme="agTheme"
                 :rowData="rowData"
                 style="height: calc(100% - 29px);width: 100%;margin-top: -1px"
                 :grid-options="gridOptions"
                 :loading="false"
                 :allowContextMenuWithControlKey="true"
                 :defaultColDef="defaultColDef"
                 :overlayNoRowsTemplate="overlayNoRowsTemplate"
                 :suppressCutToClipboard="true"
    />
  </div>
</template>
<style>
.white-svg path {
  stroke: white;
}
</style>

<style>
.fullscreen-div {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: stretch;
  overflow: hidden;
}
</style>