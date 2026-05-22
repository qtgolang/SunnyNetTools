<script>
import {AgGridVue} from "ag-grid-vue3";
import {Config_IsDark, Config_Theme_agGrid, ObjString} from "../config/Config.js";
import {AG_GRID_LOCALE_CN} from "../config/AG_ZH_CN.js";
import {ElMessage} from "element-plus";
import {URLQueryUnescape} from "../../../bindings/changeme/Service/appmain.js";
import {base64ToBytes, bytesToBase64, StringToBytes, toGBK, toUTF8} from "../config/encoding.js";

export default {
  components: {'ag-grid-vue': AgGridVue},
  data() {
    return {
      Stopped: null,
      DeleteID: null,
      agGridApi: null,
      rowData: [],
      copyTime: null,
      addValue: null,
      rowId: 0,
      overlayNoRowsTemplate: `<span style="padding: 20px;" id="HookMessageText">没有可显示的数据</span>`,
      defaultColDef: {
        flex: 1,
        sortable: false,
        suppressHeaderMenuButton: true,
        suppressHeaderContextMenu: true
      },
      ReadOnly: false,
      RowNodes: [],
      IsHasModify: false,
      CellForClipboard: null,
      gridOptions: {
        cellSelection: true,
        suppressMovableColumns: true,
        stopEditingWhenCellsLoseFocus: true, // 失去焦点时自动结束编辑
        getRowId: (params) => params.data.id,
        getContextMenuItems: this.MenuEvent,
        localeText: AG_GRID_LOCALE_CN,
        processCellForClipboard: (params) => {
          if (this.copyTime) {
            if (!this.isIntervalGreaterThan100ms(this.copyTime, new Date().getTime())) {
              return params.value;
            }
          }
          this.copyTime = new Date().getTime();
          if (this.CellForClipboard) {
            this.CellForClipboard(params)
            return params.value;
          }
          ElMessage({
            message: "复制成功,已将所选内容复制到剪辑版",
            type: 'success',
          })
          return params.value;
        },
        // sendToClipboard: this.sendToClipboard,
        onCellDoubleClicked: this.onCellDoubleClicked,
        columnDefs: [
          {
            field: "名称", tooltipField: '名称',
            minWidth: 200,
            width: 200,
            editable: true,
          },
          {
            field: "值", tooltipField: '值',
            minWidth: 200,
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
        onCellEditingStopped: (params) => {
          this.IsHasModify = true;
          if (this.Stopped) {
            this.Stopped(params)
          }
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
  },
  methods: {
    onGridReady(params) {
    },
    SetResponseCookiesType() {
      this.Empty()
      this.gridOptions.columnDefs = [
        {
          field: "名称", tooltipField: '名称',
          minWidth: 150,
          width: 150,
          editable: true,
        },
        {
          field: "值", tooltipField: '值',
          minWidth: 100,
          width: 100,
          editable: true,
        },
        {
          field: "其他值", tooltipField: '其他值',
          minWidth: 100,
          width: 100,
          editable: true,
        },
      ];
      this.agGridApi.setGridOption('columnDefs', this.gridOptions.columnDefs);
    },
    SetHeaderType() {
      this.Empty()
      this.gridOptions.columnDefs = [
        {
          field: "名称", tooltipField: '名称',
          minWidth: 150,
          width: 150,
          editable: true,
        },
        {
          field: "值", tooltipField: '值',
          minWidth: 150,
          width: 150,
          editable: true,
        },
      ];
      this.agGridApi.setGridOption('columnDefs', this.gridOptions.columnDefs);
    },
    async AddLines(array) {
      const data = [];
      for (const obj of array) {
        const m = await this.addLine(obj.name, obj.value, "UTF-8", true);
        data.push(m);
      }
      this.agGridApi.applyTransaction({add: data});
      this.agGridApi.redrawRows();       // 强制重新绘制行
      this.agGridApi.refreshCells();     // 或更新所有 cell
      this.agGridApi.ensureIndexVisible(0); // 滚动到第一行确保渲染
      this.IsHasModify = false;
    },
    async addLine(name, value, value2, i) {
      let val = "";
      let ending = "UTF8";
      let u2 = [];

      try {
        val = decodeURIComponent(value);
      } catch (e) {
        try {
          const res = await URLQueryUnescape(value);
          u2 = base64ToBytes(res);
          if (u2.length > 0) {
            try {
              val = decodeURIComponent(toGBK(u2));
              ending = "GBK";
            } catch (e) {
              try {
                val = decodeURIComponent(toUTF8(u2));
              } catch (e) {
                val = value;
              }
            }
          } else {
            val = value;
          }
        } catch (e) {
          val = value;
        }
      }

      const data = {"名称": name, "值": val, "其他值": value2, id: (this.rowId++) + "", "编码": ending};
      this.IsHasModify = true;
      if (i === true) {
        return data;
      }
      this.agGridApi.applyTransaction({add: [data]});
    },
    Empty() {
      const allData = [];
      this.agGridApi.forEachNode((node) => {
        allData.push(node.data);
      });
      this.agGridApi.applyTransaction({remove: allData});
      this.rowId = 0;
    },
    SetReadOnly(readOnly) {
      this.ReadOnly = readOnly;
      this.gridOptions.columnDefs.forEach((obj) => {
        obj.editable = !readOnly
      })
      this.agGridApi.setGridOption('columnDefs', this.gridOptions.columnDefs);
    },
    onCellDoubleClicked(params) {
      const isEditable = typeof params.colDef.editable === "function"
          ? params.colDef.editable(params)
          : params.colDef.editable;
      if (!isEditable) {
        params.api.copySelectedRangeToClipboard()
      }
    },
    MenuEvent(params) {
      const defaultMenuItems = params.defaultItems || [];
      const filteredMenuItems = defaultMenuItems.filter(item => {
        return item !== 'chartRange' && item !== "separator" && item !== 'export' && item !== 'paste' && item !== 'copyWithGroupHeaders'; // 禁用功能
      });
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
              action: () => {
                const selectedRowNodes = [];
                this.RowNodes.forEach((range) => {
                  selectedRowNodes.push(range.data)
                  if (this.DeleteID) {
                    this.DeleteID(range.data.id)
                  }
                });
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
              if (this.addValue) {
                this.addValue()
                return
              }
              this.addLine("新内容", "新内容", "")
            },
            icon: addIcon,
          });
        }
      }
      return filteredMenuItems;
    },
    toRequestHeader() {
      if (!this.IsHasModify) {
        return null;
      }
      let Header = {};
      this.agGridApi.forEachNode((node) => {
        const name = ObjString(node.data["名称"]).trim();
        const value = ObjString(node.data["值"]).trim();
        if (name !== '') {
          // 检查并更新 Header
          const key = this.hasCaseInsensitiveKey(Header, name);
          if (key !== null) {
            Header[key].push(value);
          } else {
            Header[name] = [value];
          }
        }
      });
      return Header;
    },
    updateResponseCookie(cookie) {
      let Header = {};
      this.agGridApi.forEachNode(node => {
        const name = ObjString(node.data["名称"]).trim();
        if (name.toLowerCase() === "set-cookie") {
          return
        }
        const val = ObjString(node.data["值"]).trim();
        const key = this.hasCaseInsensitiveKey(Header, name);
        if (key !== null) {
          Header[key].push(val);
        } else {
          Header[name] = [val];
        }
      });
      cookie.forEach((val) => {
        const key = this.hasCaseInsensitiveKey(Header, "Set-Cookie");
        if (key !== null) {
          Header[key].push(val);
        } else {
          Header["Set-Cookie"] = [val];
        }
      })
      return Header;
    },
    updateRequestCookie(cookie) {
      let Header = {};
      this.agGridApi.forEachNode(node => {
        const name = ObjString(node.data["名称"]).trim();
        let val = ObjString(node.data["值"]).trim();
        if (name.toLowerCase() === "cookie") {
          val = cookie;
        }
        // 检查并更新 Header
        const key = this.hasCaseInsensitiveKey(Header, name);
        if (key !== null) {
          Header[key].push(val);
        } else {
          Header[name] = [val];
        }
      });
      return Header;
    },
    toRequestCookie() {
      if (!this.IsHasModify) {
        return null;
      }
      let Cookies = "";
      this.agGridApi.forEachNode((node) => {
        const name = ObjString(node.data["名称"]).trim();
        if (name !== '') {
          const value = encodeURIComponent(node.data["值"].trim());
          if (Cookies === "") {
            Cookies = name + "=" + value;
          } else {
            Cookies += "; " + name + "=" + value;
          }
        }
      });
      return Cookies;
    },
    toRawBodyTable() {
      if (!this.IsHasModify) {
        return null;
      }
      let Table = "";
      this.agGridApi.forEachNode((node) => {
        const name = ObjString(node.data["名称"]).trim();
        if (name !== '') {
          const value = encodeURIComponent(node.data["值"].trim());
          if (Table === "") {
            Table += name + "=" + value;
          } else {
            Table += "&" + name + "=" + value
          }
        }
      });
      return bytesToBase64(StringToBytes(Table));
    },
    toResponseCookiesHeader() {
      if (!this.IsHasModify) {
        return null;
      }
      const arr = [];
      this.agGridApi.forEachNode((node) => {
        const name = ObjString(node.data["名称"]).trim();
        if (name !== '') {
          arr.push(name + "=" + node.data["值"].trim() + "; " + node.data["其他值"].trim())
        }
      });
      return arr;
    },
    hasCaseInsensitiveKey(object, key) {
      const lowerCaseKey = key.toLowerCase();
      for (const k of Object.keys(object)) {
        if (k.toLowerCase() === lowerCaseKey) {
          return k; // 返回原始键名
        }
      }
      return null; // 如果没有找到，返回 null
    },
    isIntervalGreaterThan100ms(time1, time2) {
      const difference = Math.abs(time1 - time2); // 计算时间差的绝对值
      return difference > 100; // 判断差值是否大于 100 毫秒
    }
  }
}
</script>

<template>
  <ag-grid-vue ref="agGrid"
               :theme="agTheme"
               :rowData="rowData"
               style="height: 100%;"
               :onGridReady="onGridReady"
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