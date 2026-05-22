<script>
import {AgGridVue} from "ag-grid-vue3";
import {Config_IsDark, Config_Theme_agGrid, ObjString} from "../../config/Config.js";
import {AG_GRID_LOCALE_CN} from "../../config/AG_ZH_CN.js";
import {Dialogs} from "@wailsio/runtime";
import {ProcessAddName, ProcessDelName} from "../../../../bindings/changeme/Service/appmain";
import {attachMcpDeviceName} from "../../config/mcpDeviceSync.js";
import {ElMessage} from "element-plus";

export default {
  components: {'ag-grid-vue': AgGridVue},
  data() {
    return {
      Stopped: null,
      agGridApi: null,
      rowData: [],
      addValue: null,
      RowId: 0,
      overlayNoRowsTemplate: `<span style="padding: 20px;" >你还没有添加任何进程名,你可以右键进行编辑</span>`,
      defaultColDef: {
        flex: 1,
        sortable: false,
        suppressHeaderMenuButton: true,
        suppressHeaderContextMenu: true
      },
      ReadOnly: false,
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
            field: "进程名称", tooltipField: '进程名称',
            minWidth: 310,
            width: 310,
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
          this.Save()
        },
      }
    }
  },
  watch: {},
  computed: {
    agTheme() {
      return Config_Theme_agGrid.value
    },
  },
  mounted() {
    this.agGridApi = this.$refs.agGrid.api;
    attachMcpDeviceName(this);
  },
  methods: {
    addLine(Name) {
      const data = [{
        "进程名称": Name,
        id: (this.RowId++) + "",
      }];
      this.agGridApi.applyTransaction({add: data});
    },
    SetReadOnly(readOnly) {
      this.ReadOnly = readOnly;
      this.gridOptions.columnDefs.forEach((obj) => {
        obj.editable = !readOnly
      })
      this.agGridApi.setGridOption('columnDefs', this.gridOptions.columnDefs);
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
                  await ProcessDelName(selectedRowNodes[i]["进程名称"])
                }
                this.agGridApi.applyTransaction({remove: selectedRowNodes});
                this.RowNodes = [];
                this.agGridApi.clearCellSelection()
                if (this.Stopped) {
                  this.Stopped()
                }
              },
              icon: delIcon,
            });
            filteredMenuItems.push({
              name: "清空全部",
              action: async () => {
                this.Empty()
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
              if (this.addValue) {
                this.addValue()
                return
              }
              this.addLine("双击编辑：例如:abc.exe")
            },
            icon: addIcon,
          });
          filteredMenuItems.push({
            name: "一键添加雷电模拟器相关",
            action: () => {
              this.addLine("dnplayer.exe")
              this.addLine("dnplayer.exe")
              this.addLine("LdVBoxHeadless.exe")
              this.addLine("LdBoxHeadless.exe")
              this.addLine("LdVBoxHeadless.exe")
              this.addLine("Ld9BoxHeadless.exe")
              this.addLine("VBoxNetNat.exe")
              this.Save()
            },
            icon: addIcon,
          });
          filteredMenuItems.push({
            name: "一键添加微信相关",
            action: () => {
              this.addLine("WeChat.exe")
              this.addLine("wechatweb.exe")
              this.addLine("WechatAppLauncher.exe")
              this.addLine("WeChatAppEx.exe")
              this.addLine("WeChatPlayer.exe")
              this.addLine("WechatBrowser.exe")
              this.addLine("WeChatXFile.exe")
              this.Save()
            },
            icon: addIcon,
          });
        }
      }
      return filteredMenuItems;
    },
    Empty() {
      this.agGridApi.forEachNode((node) => {
        const name = ObjString(node.data["进程名称"]).trim().toLowerCase();
        if (name !== '') {
          ProcessDelName(name)
        }
      });
      this.RowNodes = [];
      this.agGridApi.setGridOption("rowData", []);
      this.Save()
    },
    Save() {
      let array = new Map();
      let rmArray = [];
      this.agGridApi.forEachNode((node) => {
        const name = ObjString(node.data["进程名称"]).trim().toLowerCase();
        if (name !== '') {
          if (array.has(name, node.data)) {
            rmArray.push(node.data)
          } else {
            array.set(name, node.data)
          }
        } else {
          rmArray.push(node.data)
        }
      });
      this.agGridApi.applyTransaction({remove: rmArray});
      this.agGridApi.forEachNode((node) => {
        const name = ObjString(node.data["进程名称"]).trim().toLowerCase();
        ProcessAddName(name)
      })
    }
  }
}
</script>

<template>
  <ag-grid-vue ref="agGrid"
               :theme="agTheme"
               :rowData="rowData"
               style="height: 400px;width: 100%"
               :grid-options="gridOptions"
               :loading="false"
               :allowContextMenuWithControlKey="true"
               :defaultColDef="defaultColDef"
               :overlayNoRowsTemplate="overlayNoRowsTemplate"
               :suppressCutToClipboard="true"
  />
</template>
<style>
.white-svg path {
  stroke: white;
}
</style>

<style>
.fullscreen-div {
  position: fixed; /* 让 div 相对于视口固定 */
  top: 0;
  left: 0;
  width: 100vw; /* 100% 视口宽度 */
  height: 100vh; /* 100% 视口高度 */
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  color: white;
}
</style>