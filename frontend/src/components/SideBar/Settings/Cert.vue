<script>
import {AgGridVue} from "ag-grid-vue3";
import {Config_agGrid_API, Config_IsDark, Config_Theme_agGrid} from "../../config/Config.js";
import {AG_GRID_LOCALE_CN} from "../../config/AG_ZH_CN.js";
import {Dialogs, Events} from "@wailsio/runtime";
import RequestCertificateType from "./RequestCertificateType.vue";
import {
  CreateRequestCert,
  RequestCertGetCommonName,
  RequestCertRemove,
  RequestCertSetFile,
  RequestList
} from "../../../../bindings/changeme/Service/appmain";
import {ElMessage} from "element-plus";
import TitleBar from "../../TitleBar/TitleBar.vue";
import {attachMcpRequestCertReload} from "../../config/mcpCertSync.js";

export default {
  components: {TitleBar, 'ag-grid-vue': AgGridVue, "requestCertificateType": RequestCertificateType},
  data() {
    return {
      Stopped: null,
      agGridApi: null,
      rowData: [],
      addValue: null,
      overlayNoRowsTemplate: `<span style="padding: 20px;" id="HookMessageText">你还没有添加任何请求证书</span>`,
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
        onCellDoubleClicked: this.onCellDoubleClicked,
        columnDefs: [
          {
            field: "使用规则", tooltipField: '使用规则',
            minWidth: 110,
            maxWidth: 110,
            editable: true,
            cellRenderer: 'requestCertificateType',
            cellStyle: {'text-align': 'left'},
          },
          {
            field: "状态", tooltipField: '状态',
            minWidth: 100,
            width: 100,
            maxWidth: 100,
          },
          {
            field: "域名", tooltipField: '域名',
            minWidth: 170,
            width: 170,
            editable: true,
          },
          {
            field: "证书文件", tooltipField: '证书文件',
            minWidth: 200,
            width: 200,
          },
          {
            field: "证书密码", tooltipField: '证书密码',
            minWidth: 200,
            width: 200,
            editable: true,
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
          const columnKey = params.column.getColId();
          if (columnKey === "证书密码") {
            if (params.node.data["证书密码"] === "该文件类型无需密码") {
              this.agGridApi.stopEditing();
              return
            }
          }
          if (columnKey !== "证书文件" && columnKey !== "证书密码" && columnKey !== "域名" && columnKey !== "注释") {
            this.agGridApi.stopEditing();
          }
        },
        onCellEditingStopped: (params) => {
          const columnKey = params.column.getColId();
          if (columnKey === "注释") {
            this.loadCert(params, true)
            return
          }
          this.loadCert(params)
        },
      }
    }
  },
  watch: {
    isRest() {
      this.Empty()
      RequestList().then(async list => {
        const array = []
        for (let i = 0; i < list.length; i++) {
          const item = list[i];
          array.push({
            "使用规则": item["Role"],
            "域名": item["DoMain"],
            "证书文件": item["Path"],
            "证书密码": item["Pass"],
            "状态": item["LoadOk"] ? "已载入" : "未载入",
            "注释": item["Note"],
            id: item["ID"] + "",
          })
        }
        this.agGridApi.applyTransaction({add: array});
      })
    }
  },
  computed: {
    agTheme() {
      return Config_Theme_agGrid.value
    },
  },
  mounted() {
    this.agGridApi = this.$refs.agGrid.api;
    Config_agGrid_API.value = this.agGridApi;
    attachMcpRequestCertReload(() => this.reloadCertsFromBackend());
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
      this.Empty();
      RequestList().then(async list => {
        const array = []
        for (let i = 0; i < list.length; i++) {
          const item = list[i];
          array.push({
            "使用规则": item["Role"],
            "域名": item["DoMain"],
            "证书文件": item["Path"],
            "证书密码": item["Pass"],
            "状态": item["LoadOk"] ? "已载入" : "未载入",
            "注释": item["Note"],
            id: item["ID"] + "",
          })
        }
        this.agGridApi.applyTransaction({add: array});
      })
    })
  },
  methods: {
    async reloadCertsFromBackend() {
      if (!this.agGridApi) {
        return;
      }
      const list = await RequestList();
      const array = [];
      for (let i = 0; i < list.length; i++) {
        const item = list[i];
        let pass = item["Pass"];
        if (!pass && item["Path"] && (item["Path"].toLowerCase().endsWith(".pem") || item["Path"].toLowerCase().endsWith(".cer"))) {
          pass = "该文件类型无需密码";
        }
        array.push({
          "使用规则": item["Role"],
          "域名": item["DoMain"],
          "证书文件": item["Path"],
          "证书密码": pass || item["Pass"],
          "状态": item["LoadOk"] ? "已载入" : "未载入",
          "注释": item["Note"],
          id: item["ID"] + "",
        });
      }
      this.agGridApi.setGridOption("rowData", array);
      this.RowNodes = [];
      this.agGridApi.clearCellSelection();
    },
    addLine(rule, Domain, zt, filePath, password, node) {
      CreateRequestCert().then(id => {
        const data = [{
          "使用规则": rule,
          "域名": Domain,
          "证书文件": filePath,
          "证书密码": password,
          "状态": zt,
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
    onCellDoubleClicked(params) {
      const columnKey = params.column.getColId();
      if (columnKey === "证书文件") {
        const id = params.node.data.id
        const options = {
          CanChooseFiles: true,
          AllowsMultipleSelection: false,
          TreatsFilePackagesAsDirectories: true,
          AllowsOtherFiletypes: false,
          Filters: [
            {DisplayName: "P12文件", Pattern: "*.p12;*.pkcs12"},
            {DisplayName: "PEM文件", Pattern: "*.cer;*.pem"},
          ],
          Title: "选择证书文件",
        };
        Dialogs.OpenFile(options)
            .then((selectedFiles) => {
              let path = "";
              let p12 = false;
              if (selectedFiles.endsWith(".p12") || selectedFiles.endsWith(".pkcs12")) {
                path = selectedFiles;
                p12 = true;
              }
              if (selectedFiles.endsWith(".cer") || selectedFiles.endsWith(".pem")) {
                path = selectedFiles;
              }
              if (path === "") {
                ElMessage.error('请选择正确的证书文件类型');
                return
              }
              params.node.data["证书文件"] = path
              params.node.data["状态"] = "未载入"
              if (p12) {
                params.node.data["证书密码"] = "双击输入-证书密码"
              } else {
                params.node.data["证书密码"] = "该文件类型无需密码"
              }
              this.agGridApi.applyTransaction({update: [params.node.data]});
            })
            .catch((error) => {
              ElMessage.error('选择证书文件失败:' + error);
            });
      }
    },
    loadCert(params, noMessage) {
      const node = params.node.data;
      RequestCertSetFile(parseInt(node.id), node["使用规则"], node["域名"], node["证书文件"], node["证书密码"], node["注释"]).then(res => {
        if (res === "ok") {
          params.node.data["状态"] = "已载入"
          if (!noMessage) {
            ElMessage.success('证书载入成功')
          }
        } else {
          if (!noMessage) {
            ElMessage.error('证书载入失败');
          }
          params.node.data["状态"] = "载入失败"
        }
        this.agGridApi.applyTransaction({update: [params.node.data]});
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
                  await RequestCertRemove(parseInt(selectedRowNodes[i].id))
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
          }
        }
        if (this.RowNodes.length === 1) {
          //新增获取别名按钮
          {
            if (this.RowNodes[0].data["状态"] === "已载入") {
              filteredMenuItems.push({
                name: "获取证书上的别名",
                action: async () => {
                  RequestCertGetCommonName(parseInt(this.RowNodes[0].data.id)).then(name => {
                    this.RowNodes[0].data["域名"] = name
                    this.agGridApi.applyTransaction({update: [this.RowNodes[0].data]});
                    ElMessage.success('获取成功:' + name);
                    ElMessage.warning('如果不正确,请手动双击修改');
                  })
                },
              });
            }

            filteredMenuItems.push({
              name: "载入/重新载入该证书",
              action: async () => {
                this.loadCert({node: this.RowNodes[0]})
              },
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
              this.addLine("解析及发送", "www.test.com", "未载入", "双击选择证书文件", "", "")
            },
            icon: addIcon,
          });
        }
      }
      return filteredMenuItems;
    },
    Empty() {
      this.agGridApi.setGridOption("rowData", []);
    },
  }
}
</script>

<template>
  <div class="fullscreen-div" style="display: block;">
    <TitleBar Title="请求证书设置"></TitleBar>
    <ag-grid-vue ref="agGrid"
                 :theme="agTheme"
                 :rowData="rowData"
                 style="height: calc(100% - 29px);width:100%;margin-top: -1px;margin-left: -1px"
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