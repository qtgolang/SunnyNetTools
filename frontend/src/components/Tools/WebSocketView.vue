<script>
import {AgGridVue} from "ag-grid-vue3";
import {
  Config_Find_Range_ALL,
  Config_Find_Window,
  Config_Find_Window_Hide,
  Config_Find_Window_Show, Config_HomeTextMark,
  Config_IsDark,
  Config_SelectedRow,
  Config_SocketSelectedRow,
  Config_Theme_agGrid
} from "../config/Config.js";
import {AG_GRID_LOCALE_CN} from "../config/AG_ZH_CN.js";
import {ElMessage} from "element-plus";
import {
  AppInsertDone,
  ClearAllSessionMessageIdArray,
  CopySessionMessageIdArray,
  DelSessionMessageIdArray,
  GOOS,
  StreamSearch
} from "../../../bindings/changeme/Service/appmain.js";
import {Events} from "@wailsio/runtime";
import ImageRenderer from './SocketImage.vue';
import Filter from "./Filter/filter.vue";

export default {
  props: ["Name"],
  components: {Filter, 'ag-grid-vue': AgGridVue, imageRenderer: ImageRenderer,},
  data() {
    return {
      get SelectedRow() {
        return Config_SelectedRow.value
      },
      set SelectedRow(value) {
        Config_SelectedRow.value = value
      },
      IsWindows: false,
      Filter: null,
      isWebsocket: false,
      MenuVisible: false,
      Theology: -1,
      agGridApi: null,
      rowData: [],
      copyTime: null,
      overlayNoRowsTemplate: `<span style="padding: 20px;" id="HookMessageText">没有可显示的数据</span>`,
      defaultColDef: {
        flex: 1,
        sortable: false,
        suppressHeaderMenuButton: true,
        suppressHeaderContextMenu: true
      },
      RowNodes: [],
      rowId: 0,
      gridOptions: {
        getRowStyle: this.onGetRowStyle,
        onContextMenuVisibleChanged: (event) => {
          this.MenuVisible = event.visible;
        },
        popupParent: document.getElementById("appMain"),
        enableAdvancedFilter: true,
        cellSelection: true,
        getRowId: (params) => params.data.MessageId,
        getContextMenuItems: this.MenuEvent,
        localeText: AG_GRID_LOCALE_CN,
        processCellForClipboard: (params) => {
          if (this.copyTime) {
            if (!this.isIntervalGreaterThan100ms(this.copyTime, new Date().getTime())) {
              return params.value;
            }
          }
          this.copyTime = new Date().getTime();
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
            field: "#", tooltipField: '#',
            minWidth: 50,
            width: 50,
            valueGetter: (params) => params.node.rowIndex + 1,
            maxWidth: 100,
            editable: false,
            cellRenderer: 'imageRenderer',
          },
          {
            field: "时间", tooltipField: '时间',
            minWidth: 120,
            width: 120,
            filter: true,
            editable: false,
          },
          {
            field: "类型", tooltipField: '类型',
            minWidth: 100,
            width: 100,
            filter: true,
            editable: false,
          },
          {
            field: "长度", tooltipField: '长度',
            minWidth: 60,
            maxWidth: 60,
            width: 60,
            filter: true,
            editable: false,
            filterParams: {filterOptions: ["equals", "notEqual", "greaterThan", "greaterThanOrEqual", "lessThan", "lessThanOrEqual"]},
          },
          {
            field: "数据", tooltipField: '数据',
            minWidth: 200,
            width: 200,
            filter: true,
            editable: false,
          },
          {
            field: "MessageId",
            minWidth: 0,
            width: 0,
            filter: false,
            hide: true,
            sortable: true
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
            if (selectedRowNodes.length > 0) {
              if (!___c) {
                Config_SocketSelectedRow.value = selectedRowNodes[0].data
              } else {
                Config_SocketSelectedRow.value = selectedRowNodes[selectedRowNodes.length - 1].data
              }
            } else {
              Config_SocketSelectedRow.value = [];
            }
            this.RowNodes = selectedRowNodes;
          }
        },
        isExternalFilterPresent: this.isExternalFilterPresent,
        doesExternalFilterPass: this.doesExternalFilterPass,
      },
      isGoFilter: false,//是否使用Go过滤
      GoFilterList: new Set([]),//只存放需要显示的 MessageId
      isSearch: false,  //是搜索结果
      GoSearchList: new Set([]),
      FollowDisplay: true,
      CopyMenuSubMenu: [
        {name: "复制（ 全部-仅可见数据 ）( HEX ) 到剪辑版", action: () => this.copyData("all1")},
        {name: "复制（ 全部-仅可见数据 ）(客户端发送) ( HEX ) 到剪辑版", action: () => this.copyData("send1")},
        {name: "复制（ 全部-仅可见数据 ）(客户端接收) ( HEX ) 到剪辑版", action: () => this.copyData("rec1")},
        "separator",
        {name: "复制（ 全部-含不可见数据 ）( HEX ) 到剪辑版", action: () => this.copyData("all")},
        {name: "复制（ 全部-含不可见数据 ）(客户端发送) ( HEX ) 到剪辑版", action: () => this.copyData("send")},
        {name: "复制（ 全部-含不可见数据 ）(客户端接收) ( HEX ) 到剪辑版", action: () => this.copyData("rec")}
      ],
      CopyMenuSubMenuSelected: {name: "复制（ 所选数据 ）( HEX ) 到剪辑版", action: () => this.copyData("selected")},
    }
  },
  computed: {
    agTheme() {
      return Config_Theme_agGrid.value
    },
  },
  watch: {
    "SelectedRow"(newValue) {
      try {
        this.Theology = parseInt(newValue["Theology"]);
      } catch (e) {
        this.Theology = -1;
      }
    }
  },
  mounted() {
    this.agGridApi = this.$refs.agGrid.api;
    this.isWebsocket = (this.Name === "Websocket")
    {
      Events.On("updateSocketStreamList", (obj) => {
        const array = obj?.data?.[0] ?? [];
        const isDone = obj?.data?.[1] ?? false;
        this.InsertSocketStream(array, isDone)
        //AppInsertDone()
      })
    }
    GOOS().then(isWindows => {
      this.IsWindows = isWindows;
    })
  },
  methods: {
    InsertSocketStream(array, isDone) {
      if (Array.isArray(array) && array.length > 0) {
        const api = this.agGridApi;
        const newArray = [];
        const _Theology_ = this.Theology;
        array.forEach((element) => {
          if (parseInt(element["Theology"]) !== _Theology_) {
            return
          }
          let _Type = '';
          let _Ico = '上行';
          if (element["IsClose"]) {
            _Type = '连接关闭'
            _Ico = 'stop';
          } else {
            const IsActiveSend = element["IsActiveSend"] ? "手动" : ""
            _Type = IsActiveSend + (element["IsSend"] ? "发送" : "接收");
            if (this.isWebsocket) {
              _Type += "(";
              switch (element["WebsocketType"]) {
                case 1:
                  _Type += "Text";
                  break
                case 2:
                  _Type += "Binary";
                  break
                case 8:
                  _Type += "Close";
                  break
                case 9:
                  _Type += "Ping";
                  break
                case 10:
                  _Type += "Pong";
                  break
                default:
                  _Type += "invalid";
                  break
              }
              _Type += ")";
            }
            _Ico = element["IsSend"] ? "上行" : "下行";
          }
          const obj = {
            "MessageId": element["MessageId"] + "",
            "时间": element["Time"],
            "数据": element["Body"],
            "IsClose": element["IsClose"],
            "类型": _Type,
            "长度": element['Length'],
            "ico": _Ico
          }
          if (this.isGoFilter) {
            if (element["Filter"]) {
              this.GoFilterList.add(element["MessageId"])
            }
          }
          newArray.push(obj)
        })
        if (newArray.length > 0) {
          const res = api.applyTransaction({add: newArray});
          //api.applyColumnState({state: [{colId: 'MessageId', sort: 'asc'}], defaultState: {sort: null},});
          //this.ensureNodeVisible(res.add[res.add.length - 1].data.MessageId);
          if (res.add.length > 0) {
            if (isDone && this.FollowDisplay) {
              this.ensureNodeVisible(res.add[res.add.length - 1].data.MessageId);
            }
          }
        }
      }
      AppInsertDone()
    },
    ensureNodeVisible(MessageIId) {
      requestAnimationFrame(() => {
        if (this.FollowDisplay && !this.MenuVisible) {
          const node = this.agGridApi.getRowNode(MessageIId + "");
          if (node) {
            setTimeout(() => {
              this.agGridApi.ensureNodeVisible(node);
              setTimeout(() => {
                this.agGridApi.ensureNodeVisible(node);
                setTimeout(() => {
                  this.agGridApi.ensureNodeVisible(node);
                }, 200)
              }, 200)
            }, 200)
            return
          }
          this.ensureNodeVisible(MessageIId)
        }
      })
    },
    onGetRowStyle(params) {
      let res = {
        fontFamily: "微软雅黑"
      }
      if (this.IsWindows) {
        res.fontFamily = `-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", "Helvetica Neue", Helvetica, Arial, sans-serif`;
      }
      if (params.data.backcolor) {
        const ColorID = params.data.backcolor;
        if (Config_HomeTextMark.has(ColorID)) {
          const value = Config_HomeTextMark.get(ColorID)
          res.color = Config_IsDark.value ? value["深色主题"] : value["浅色主题"]
          res.fontWeight = 'bold';
          return res
        }
      }
      return res
    },
    CancelSearch() {
      this.isSearch = false;
      this.agGridApi.onFilterChanged()
    },
    SearchDone(params, mode, Color, isCancelColor) {
      if (mode === "hide") {
        this.isSearch = true;
        for (let i = 0; i < params.length; i++) {
          this.GoSearchList.add(params[i])
        }
        this.agGridApi.onFilterChanged()
        Config_Find_Window_Hide.value()
        return
      }
      const array = []
      if (isCancelColor) {
        this.agGridApi.forEachNode((node) => {
          if (params.indexOf(parseInt(node.data["MessageId"])) !== -1) {
            node.data.backcolor = Color;
          } else {
            node.data.backcolor = "";
          }
          array.push(node)
        });
      } else {
        for (let i = 0; i < params.length; i++) {
          const node = this.agGridApi.getRowNode(params[i] + "");
          if (node) {
            node.data.backcolor = Color;
            array.push(node)
          }
        }
      }
      if (array.length > 0) {
        this.agGridApi.redrawRows({rowNodes: array});
      }
    },
    onFilterApply(rules) {
      if (!rules || Object.keys(rules).length === 0) {
        StreamSearch(this.Theology, "clear").then((array) => {
          this.isGoFilter = false;
          this.agGridApi.onFilterChanged()
        });
        return;
      }
      this.isGoFilter = true;
      this.GoFilterList.clear()
      StreamSearch(this.Theology, JSON.stringify(rules)).then((array) => {
        array.forEach((MessageId) => {
          this.GoFilterList.add(MessageId)
        })
        this.agGridApi.onFilterChanged()
      });
    },
    copyData(mode) {
      const array = [];
      let CopyType = "copy";
      const rowCount = this.agGridApi.getDisplayedRowCount(); // ✅ 获取可见行数
      switch (mode) {
        case "selected":
          for (let i = 0; i < this.RowNodes.length; i++) {
            array.push(parseInt(this.RowNodes[i].data.MessageId))
          }
          break
        case "all1":
          for (let i = 0; i < rowCount; i++) {
            const rowNode = this.agGridApi.getDisplayedRowAtIndex(i); // ✅ 获取每个可见行
            array.push(parseInt(rowNode.data.MessageId)); // ✅ 获取 id
          }
          break
        case "send1":
          for (let i = 0; i < rowCount; i++) {
            const rowNode = this.agGridApi.getDisplayedRowAtIndex(i); // ✅ 获取每个可见行
            if (rowNode.data["类型"].indexOf("发送") !== -1) {
              array.push(parseInt(rowNode.data.MessageId)); // ✅ 获取 id
            }
          }
          break
        case "rec1":
          for (let i = 0; i < rowCount; i++) {
            const rowNode = this.agGridApi.getDisplayedRowAtIndex(i); // ✅ 获取每个可见行
            if (rowNode.data["类型"].indexOf("发送") === -1) {
              array.push(parseInt(rowNode.data.MessageId)); // ✅ 获取 id
            }
          }
          break
        default:
          CopyType = mode;
          break
      }

      if (this.SelectedRow.Theology === undefined) {
        ElMessage({
          message: '未知错误',
          type: 'error',
        })
        return
      }
      CopySessionMessageIdArray(parseInt(this.SelectedRow.Theology), CopyType, array).then((res) => {
        if (res === "") {
          ElMessage({
            message: "复制数据成功",
            type: 'success',
          })
        } else {
          ElMessage({
            message: "复制出现错误:" + res,
            type: 'error',
          })
        }
      })
    },
    isExternalFilterPresent(params) {
      return this.isGoFilter || this.isSearch;
    },
    doesExternalFilterPass(node) {
      if (this.isSearch) {
        return this.GoSearchList.has(parseInt(node.data['MessageId']));
      }
      if (!this.isGoFilter) {
        return true;
      }
      return this.GoFilterList.has(parseInt(node.data['MessageId']));
    },
    //点击高级搜索
    onClickSearchButton() {
      if (Config_Find_Window.value.isSearchInProgress) {
        ElMessage({
          message: '正在查找其他值,请稍后..',
          type: 'error',
        })
        return
      }
      Config_Find_Range_ALL.value = false
      Config_Find_Window.value.Title = "🚀 [ " + this.Name + " ] 高级搜索";
      Config_Find_Window.value.CompleteCallback = this.agGridApi.ShowCancelSearch
      Config_Find_Window_Show.value()
    },
    onGridReady(params) {
      //向高级过滤器添加一个按钮
      {
        const customButtonContainer = this.$el.getElementsByClassName('ag-advanced-filter ag-advanced-filter-header-cell');
        if (customButtonContainer.length > 0) {
          if (customButtonContainer[0]) {
            const popup = customButtonContainer[0];
            popup.style.display = "none";
            this.Filter = popup.parentElement
          }
        }
      }
    },
    EventClick(event) {
      this.$refs.Filter.setFocusElement();
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
    Empty() {
      this.agGridApi.setGridOption("rowData", []);
    },
    MenuEvent(params) {
      const defaultMenuItems = params.defaultItems || [];
      const filteredMenuItems = defaultMenuItems.filter(item => {
        return item !== 'chartRange' && item !== "separator" && item !== 'export' && item !== 'paste' && item !== 'copyWithGroupHeaders'; // 禁用功能
      });

      //跟随显示
      {
        let Icon = '';
        if (this.FollowDisplay) {
          if (Config_IsDark.value) {
            Icon = '<div style="display: flex; align-items: center;width: 16px">' +
                `<svg xmlns="http://www.w3.org/2000/svg" style="top: 2px;position: relative;" width="16" height="14" viewBox="0 0 24 24" fill="none" stroke="#fff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">\\n' +
              '    <polyline points="20 6 9 17 4 12"/>' +
              '</svg>` +
                '</div>';
          } else {
            Icon = '<div style="display: flex; align-items: center;width: 16px">' +
                `<svg xmlns="http://www.w3.org/2000/svg" style="top: 2px;position: relative;" width="16" height="14" viewBox="0 0 24 24" fill="none" stroke="#000" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">\\n' +
              '    <polyline points="20 6 9 17 4 12"/>' +
              '</svg>` +
                '</div>';
          }
        }
        filteredMenuItems.push("separator");
        filteredMenuItems.push({
          name: "跟随显示",
          action: () => {
            this.FollowDisplay = !this.FollowDisplay
          },
          icon: Icon,
        });
      }

      const isEmpty = params.api.getDisplayedRowCount() === 0;
      if (!isEmpty) {
        filteredMenuItems.push("separator");
        const arrays = [];
        if (this.RowNodes.length > 0) {
          arrays.push(this.CopyMenuSubMenuSelected)
          arrays.push("separator");
        }
        this.CopyMenuSubMenu.forEach((o) => {
          arrays.push(o)
        })
        filteredMenuItems.push({
          name: "复制数据",
          subMenu: arrays,
          icon: `<span class="ag-icon ag-icon-copy"></span>`
        });
        filteredMenuItems.push("separator");
      }
      let delIcon = '';
      if (this.RowNodes.length > 0) {
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
            if (this.SelectedRow.Theology === undefined) {
              return
            }
            const array = []
            for (let i = 0; i < this.RowNodes.length; i++) {
              array.push(parseInt(this.RowNodes[i].data.MessageId))
            }
            this.RowNodes.forEach((node) => {
              node.setSelected(false);
            });
            DelSessionMessageIdArray(parseInt(this.SelectedRow.Theology), array).then((objs) => {
              const api = this.agGridApi;
              api.setGridOption("rowData", []);
              this.RowNodes = [];
              Config_SocketSelectedRow.value = []
              api.clearRangeSelection()
            })
          },
          icon: delIcon,
        });
      }
      if (!isEmpty) {
        filteredMenuItems.push({
          name: "删除全部",
          action: () => {
            if (this.SelectedRow.Theology === undefined) {
              return
            }
            ClearAllSessionMessageIdArray(parseInt(this.SelectedRow.Theology)).then((objs) => {
              const api = this.agGridApi;
              api.setGridOption("rowData", []);
              this.RowNodes = [];
              Config_SocketSelectedRow.value = []
              api.clearRangeSelection()
            })
          },
          icon: delIcon,
        });
      }
      return filteredMenuItems;
    },
    isIntervalGreaterThan100ms(time1, time2) {
      const difference = Math.abs(time1 - time2); // 计算时间差的绝对值
      return difference > 100; // 判断差值是否大于 100 毫秒
    },
  }
}
</script>

<template>
  <div style="width: 100%;height: 100%">
    <ag-grid-vue ref="agGrid"
                 @click.stop="EventClick"
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
    <Filter ref="Filter" :Name="Name" :Parent="Filter" :column="gridOptions.columnDefs" :apply="onFilterApply"
            :SearchDone="SearchDone" :CancelSearch="CancelSearch"/>
  </div>
</template>
<style>
.white-svg path {
  stroke: white;
}
</style>