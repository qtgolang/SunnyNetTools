<script>
//主题深色 ag-theme-balham-dark
//主题浅色 ag-theme-balham
import {reactive} from 'vue';
import '@ag-grid-community/styles/ag-grid.css';
import '@ag-grid-community/styles/ag-theme-balham.css';
import {AgGridVue} from '@ag-grid-community/vue3';
import {ClipboardModule} from '@ag-grid-enterprise/clipboard';
import {SetFilterModule} from '@ag-grid-enterprise/set-filter';
import {ExcelExportModule} from '@ag-grid-enterprise/excel-export';
import ImageRenderer from './image.vue';
import ToolPanel from './ToolPanel.vue';
import FindWindow from './Find/FindWindow.vue';
import AgHeaderGroup from './Header.vue';
import AgFooterGroup from './Footer.vue';
import {CallGoDo, deepCopy, SetTextColor, StrBase64Encode} from "./CallbackEventsOn.js";
import {ClipboardSetText} from "../../wailsjs/runtime/runtime.js";

window.Theme = reactive({
  IsDark: true,
  GOOS: Window.GOOS,
});
window.Socket = reactive({
  Line: null,
  Data: null
});
window.UI = reactive({
  ZIndex: {
    Settings: 0,
    TextCompare: 0,
    FindWindow: 0,
    DocCompare: 0,
    OpenSourceProtocol: 0,
  },
  Settings: false,
  TextCompare: false,
  FindWindow: false,
  DocCompare: false,
  OpenSourceProtocol: false
});
window.Size = reactive({
  Settings: {Width: 0, Height: 0},
  Doc: {Width: 0, Height: 0},
  OpenSourceProtocol: {Width: 0, Height: 0},
});
window.SetUILevel = function (name) {
  let max = 36
  for (let key in window.UI.ZIndex) {
    if (window.UI.ZIndex[key] > max) {
      max = window.UI.ZIndex[key]
    }
  }
  window.UI.ZIndex[name] = max + 1
}

export default {
  beforeMount() {

  }, components: {
    'ag-grid-vue': AgGridVue, imageRenderer: ImageRenderer, FindWindow,
    ToolPanel, AgHeaderGroup, AgFooterGroup
  },
  data: function () {
    return {
      index: 0,
      agGridApi: null,
      //选择的行数组
      agSelectedArray: [],
      //点击的列
      colId: "",
      //当前选中行
      agSelectedLine: null,
      RowData: [],
      RowDataHashMap: {},
      TagColorMap: [],
      ListDom: null,
      //鼠标是否在列表范围内
      IsListDomRange: false,
      leftModules: [SetFilterModule, ClipboardModule],
      rightModules: [ExcelExportModule],
      defaultColDef: {
        flex: 1,
        minWidth: 10,
        sortable: false,  // 禁用全部列的排序功能
        filter: "agTextColumnFilter",
        filterParams: {
          maxNumConditions: 10
        },
        resizable: true,
        menuTabs: ['filterMenuTab'],
        suppressNavigable: false,
        cellClass: 'no-border',
      },
      gridOptions: {
        getRowId: (params) => params.data.Theology,
        //rowSelection: 'multiple',
        onRangeSelectionChanged: this.onRangeSelectionChanged,
        onRowClicked: this.onRowClicked,
        onCellFocused: this.onCellFocused,
        getContextMenuItems: this.onContextMenuItems,
        getRowStyle: this.onGetRowStyle,
        onRowDataUpdated: this.NewColumnsLoaded,
        onModelUpdated: this.NewColumnsLoaded,
        suppressScrollOnNewData: true, // 禁用自动滚动到第一行
        onColumnVisible: this.onColumnChange,
        onColumnMoved: this.onColumnChange,
        onColumnResized: this.onColumnChange,
        onCellValueChanged: this.CellValueChanged,
      },
      columns: [
        {
          field: "序号", tooltipField: '序号',
          minWidth: 80,
          width: 80,
          maxWidth: 100,
          menuTabs: [], // 不显示过滤器
          cellRenderer: 'imageRenderer', cellStyle: {'text-align': 'left'},
        },
        {
          field: "方式", tooltipField: '方式',
          minWidth: 90,
          width: 90,
          maxWidth: 90, cellStyle: {'text-align': 'left'}
        },
        {
          field: "状态", tooltipField: '状态',
          minWidth: 80,
          width: 80,
          maxWidth: 80, cellStyle: {'text-align': 'left'}
        },
        {
          field: "请求地址", width: 200, minWidth: 200, tooltipField: '请求地址',
          maxWidth: 2000, cellStyle: {'text-align': 'left'}
        },
        {
          field: "响应长度", width: 120, minWidth: 120, tooltipField: '响应长度',
          maxWidth: 120, cellStyle: {'text-align': 'left'}
        },
        {
          field: "响应类型", width: 120, minWidth: 120, tooltipField: '响应类型',
          maxWidth: 120, cellStyle: {'text-align': 'left'}
        },
        {
          field: "进程", width: 100, minWidth: 100, tooltipField: '进程',
          maxWidth: 500, cellStyle: {'text-align': 'left'}
        },
        {
          field: "请求时间", width: 150, minWidth: 150, tooltipField: '请求时间', hide: true,
          maxWidth: 150, cellStyle: {'text-align': 'left'}
        },
        {
          field: "响应时间", width: 150, minWidth: 150, tooltipField: '响应时间', hide: true,
          maxWidth: 150, cellStyle: {'text-align': 'left'}
        },
        {
          field: "注释", width: 200, minWidth: 200, tooltipField: '注释', editable: true,
          maxWidth: 2000, cellStyle: {'text-align': 'left'}
        },
        {
          field: "来源地址", width: 150, minWidth: 150, tooltipField: '来源地址',
          maxWidth: 200,
        },
      ],
      sideBar: null,
      //是否跟随显示
      ListFollowShow: true,
      PageWidth: {
        Width: 0,
        min: 0,
        max: 0,
      },
      agTools: null,
      FindShow: false,
      MenuItems: [
        {
          name: '生成请求代码',
          subMenu: [
            {
              name: '易语言',
              disabled: false,
              visible: true,
              subMenu: [
                {
                  name: 'WinInet',
                  action: () => {
                    this.GenerateRequestCode("E", 'WinInet')
                  },
                  disabled: false,
                  visible: true
                },
                {
                  name: 'WinHttpW',
                  action: () => {
                    this.GenerateRequestCode("E", 'WinHttpW')
                  },
                  disabled: false,
                  visible: true
                },
                {
                  name: 'WinHttpR',
                  action: () => {
                    this.GenerateRequestCode("E", 'WinHttpR')
                  },
                  disabled: false,
                  visible: true
                },
                {
                  name: '网页_访问',
                  action: () => {
                    this.GenerateRequestCode("E", '网页_访问')
                  },
                  disabled: false,
                  visible: true
                },
                {
                  name: '网页_访问_对象',
                  action: () => {
                    this.GenerateRequestCode("E", '网页_访问_对象')
                  },
                  disabled: false,
                  visible: true
                },
                {
                  name: 'E2EE网站客户端',
                  action: () => {
                    this.GenerateRequestCode("E", 'e2ee')
                  },
                  disabled: false,
                  visible: true
                }
              ],
            },
            {
              name: 'GoLang',
              subMenu: [
                {
                  name: 'net/http',
                  action: () => {
                    this.GenerateRequestCode('Go', "net/http")
                  },
                  disabled: false,
                  visible: true
                },
              ],
              disabled: false,
              visible: true
            },
            {
              name: 'C#',
              subMenu: [
                {
                  name: 'HttpClient',
                  action: () => {
                    this.GenerateRequestCode('C#', "HttpClient")
                  },
                  disabled: false,
                  visible: true
                },
                {
                  name: 'RestSharp',
                  action: () => {
                    this.GenerateRequestCode('C#', "RestSharp")
                  },
                  disabled: false,
                  visible: true
                },
              ],
              disabled: false,
              visible: true
            },
            {
              name: 'Python',
              subMenu: [
                {
                  name: 'requests',
                  action: () => {
                    this.GenerateRequestCode('Python', "requests")
                  },
                  disabled: false,
                  visible: true
                },
              ],
              disabled: false,
              visible: true
            }
          ],
          disabled: false,
          visible: true
        },
        {
          name: '颜色标记',
          subMenu: [
            {
              name: '红色',
              action: () => {
                this.markerColor("红色")
              },
              disabled: false,
              visible: true
            },
            {
              name: '蓝色',
              action: () => {
                this.markerColor("蓝色")
              },
              disabled: false,
              visible: true
            },
            {
              name: '绿色',
              action: () => {
                this.markerColor("绿色")
              },
              disabled: false,
              visible: true
            },
            {
              name: '黄色',
              action: () => {
                this.markerColor("黄色")
              },
              disabled: false,
              visible: true
            },
            {
              name: '紫色',
              action: () => {
                this.markerColor("紫色")
              },
              disabled: false,
              visible: true
            },
            'separator',
            {
              name: '取消颜色标记',
              action: () => {
                this.markerColor("empty")
              },
              disabled: false,
              visible: true
            },
            'separator',
            {
              name: '取消搜索颜色标记',
              action: this.CancelSearch,
              disabled: false,
              visible: false
            },
          ],
          disabled: false,
          visible: true
        },
        'separator',
        {
          name: '删除选择',
          action: () => {
            this.delete()
          },
          disabled: false,
          visible: true
        },
        'separator',
        {
          name: '查找',
          action: () => {
            this.$nextTick(() => {
              window.UI.FindWindow = true
              window.SetUILevel("FindWindow")
            })
          },
          disabled: false,
          visible: true
        },
        {
          name: '重发',
          subMenu: [
            {
              name: '普通重发',
              action: () => {
                this.resend(3)
              },
              disabled: false,
              visible: true
            },
            {
              name: '重发 并 拦截上行',
              action: () => {
                this.resend(1)
              },
              disabled: false,
              visible: true
            },
            {
              name: '重发 并 拦截下行',
              action: () => {
                this.resend(2)
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
          name: '跟随显示',
          action: () => {
            this.ListFollowShow = !this.ListFollowShow
            window.vm.Header.SetAutoRollShow(this.ListFollowShow)
          },
          disabled: false,
          visible: true,
          icon: ''
        },
        {
          name: '断开选中的会话连接',
          action: this.CloseSession,
          disabled: false,
          visible: true,
          icon: ''
        },
      ], // 返回自定义的菜单项
      overlayNoRowsTemplate: `<span style="padding: 20px;" id="HookMessageText">还没有捕获到数据</span>`,
      get darkTheme() {
        return window.Theme.IsDark
      },
      set darkTheme(newValue) {
        window.Theme.IsDark = newValue
      },
    };
  },
  methods: {
    sendToClipboard(params) {
      ClipboardSetText(params.data)
    },
    SetEmptySearchMenuVisible(value) {
      this.MenuItems[1].subMenu[8].visible = value
    },
    CancelSearch() {
      this.SetEmptySearchMenuVisible(false)
      CallGoDo("取消搜索颜色标记", null).then((res) => {
        if (res) {
          for (let i = 0; i < res.length; i++) {
            const LastSearchResult = res[i]
            const obj = window.vm.List.RowDataHashMap[LastSearchResult]
            if (obj) {
              obj.data.color.search = null
              obj.setData(obj.data)
            }
          }
        }
        this.agSelectedLine = null
        this.RefreshRenderedNodes()
        this.$nextTick(() => {
          this.agSelectedLine = this.RowDataHashMap[window.Theology]
        })
      })
    },
    //隐藏工具面板
    HideToolsPanel() {
      // 创建一个 <style> 元素节点
      const styleElement = document.createElement('style');
      // 定义 CSS 规则
      const cssRules = `.ag-side-buttons {
                        padding-top: calc(var(--ag-grid-size) * 4);
                        width: 0;
                        position: relative;
                        overflow: hidden;
                        }`;
      // 将 CSS 规则写入 <style> 元素
      styleElement.innerHTML = cssRules;
      // 将 <style> 元素节点附加到文档的 <head> 元素中
      document.head.appendChild(styleElement);
    }
    ,
    updateWindowSize() {
      this.PageWidth.Width = window.innerWidth * 0.3;
      this.PageWidth.max = window.innerWidth - 50;
      this.PageWidth.min = 620//window.innerWidth * 0.3;
      if (this.PageWidth.Width < this.PageWidth.min) {
        this.PageWidth.Width = this.PageWidth.min
      }
      if (this.PageWidth.Width > this.PageWidth.max) {
        this.PageWidth.Width = this.PageWidth.max
      }
      for (let i = 0; i < this.sideBar.toolPanels.length; i++) {
        this.sideBar.toolPanels[i].width = this.PageWidth.Width
      }
      this.handleMouseMove()
    }
    ,
    getTools() {
      for (let i = 0; i < this.agTools.length; i++) {
        if (this.agTools[i].className.indexOf("ag-hidden") === -1) {
          return this.agTools[i]
        }
      }
      return null
    }
    ,
    handleMouseMove(event) {
      try {
        const tools = this.getTools()
        if (tools != null) {
          if (window.ToolsMaximize === true) {
            tools.style.width = window.innerWidth + "px"
            return
          }
          const w = parseInt(tools.style.width.replace('px', ''))
          if (w < this.PageWidth.min) {
            tools.style.width = this.PageWidth.min + "px"
          }
          if (w > this.PageWidth.max) {
            tools.style.width = this.PageWidth.max + "px"
          }
        }
      } catch (e) {
      }
    }
    ,
    handleKeyDown(event) {
      if (event.ctrlKey && (event.key === "f" || event.key === "F")) {
        window.UI.FindWindow = false
        this.$nextTick(() => {
          window.UI.FindWindow = true
          window.SetUILevel("FindWindow")
          this.$nextTick(() => {
            window.vm.Find.SetFocus()
          })
        })
      } else if (event.ctrlKey && (event.key === "c" || event.key === "C")) {
        if (this.IsListDomRange) {
          this.agGridApi.copySelectedRangeToClipboard();
        }
      }
    },
    onRangeSelectionChanged(event) {
      const rangeSelections = event.api.getCellRanges();
      const selectedRowNodes = [];
      rangeSelections.forEach(function (range) {
        const _a = range.startRow.rowIndex;
        const _b = range.endRow.rowIndex;
        let startRow = _a;
        let endRow = _b;
        if (_a > _b) {
          startRow = _b
          endRow = _a
        }
        event.api.forEachNodeAfterFilter(function (node) {
          if (node.rowIndex >= startRow && node.rowIndex <= endRow) {
            selectedRowNodes.push(node);
          }
        });
      });
      this.agSelectedArray = selectedRowNodes;
    },
    onRowClicked(params) {
      if (this.colId === "注释") {
        return
      }
      let array = []
      if (this.agSelectedLine) {
        const obj = this.RowDataHashMap[this.agSelectedLine.data.Theology]
        if (obj) {
          if (obj.data.color) {
            obj.data.color.selected = {}
          } else {
            obj.data.color = {selected: {}}
          }
          array.push(obj)
        }
      }
      this.agSelectedLine = params.node
      const obj = this.RowDataHashMap[this.agSelectedLine.data.Theology]
      if (obj) {
        const selectedColor = {dark: '#1b5168', right: "#b0c5cc"}
        if (obj.data.color) {
          obj.data.color.selected = selectedColor
        } else {
          obj.data.color = {selected: selectedColor}
        }
        array.push(obj)
      }

      this.$nextTick(() => {
        this.RefreshRenderedNodes()
        //this.agGridApi.refreshCells({rowNodes: array, force: true})
      })
    },
    CellValueChanged(event) {
      const newValue = event.newValue;
      const Theology = event.data['Theology'];
      CallGoDo("更新注释", {Theology: Theology, Data: newValue})
    },
    onCellFocused(event) {
      this.colId = event.column.colId
    },
    onContextMenuItems(params) {
      if (this.agSelectedArray.length < 1) {
        this.MenuItems[0].visible = false
        this.MenuItems[1].visible = false
        this.MenuItems[3].visible = false
        this.MenuItems[5].visible = false
        this.MenuItems[6].visible = false
        this.MenuItems[8].visible = true
        this.MenuItems[9].visible = false
        return getMenuItems(this.MenuItems)
      } else {
        this.MenuItems[0].visible = true
        this.MenuItems[1].visible = true
        this.MenuItems[3].visible = true
        this.MenuItems[5].visible = true
        this.MenuItems[6].visible = true
        this.MenuItems[8].visible = true
      }
      this.onRowClicked(params)
      this.onContextMenu();
      return getMenuItems(this.MenuItems)
    },
    onContextMenu() {
      let ok = false
      let showCF = true
      for (let i = 0; i < this.agSelectedArray.length; i++) {
        if (this.agSelectedArray[i].data) {
          const fs = ("" + this.agSelectedArray[i].data['方式']).toUpperCase()
          const zt = ("" + this.agSelectedArray[i].data['状态']).toUpperCase()
          if (fs.indexOf("TCP") !== -1 || fs.indexOf("UDP") !== -1 || fs.indexOf("WEBSOCKET") !== -1) {
            showCF = false
          }
          if (fs.indexOf("TCP") !== -1 && zt.indexOf("已连接") !== -1) {
            ok = true
            break
          }

          if (fs.indexOf("WEBSOCKET") !== -1 && zt.indexOf("已连接") !== -1) {
            ok = true
            break
          }
        }
      }
      if (this.MenuItems.length > 8) {
        this.MenuItems[0].visible = showCF
        this.MenuItems[6].visible = showCF
        this.MenuItems[9].visible = ok
      }
    },
    onColumnChange(event) {
      const GridColumns = this.$refs.agGrid.gridOptions.columnApi.getAllGridColumns()
      let Columns = []
      for (let i = 0; i < GridColumns.length; i++) {
        const m = GridColumns[i]
        const colDef = m.colDef
        let Column = {
          field: colDef.field, tooltipField: colDef.tooltipField,
          minWidth: colDef.minWidth,
          width: m.actualWidth,
          maxWidth: colDef.maxWidth,
          menuTabs: colDef.menuTabs,
          cellRenderer: colDef.cellRenderer, cellStyle: colDef.cellStyle,
          hide: !m.visible,
        }
        Columns[i] = Column
      }
      const ColumnsObjs = StrBase64Encode(JSON.stringify(Columns))
      CallGoDo("保存配置", {Type: "列数据", Data: ColumnsObjs})
    },
    GenerateRequestCode(Lang, module) {
      const array = []
      for (let i = 0; i < this.agSelectedArray.length; i++) {
        array.push(this.agSelectedArray[i].data['Theology'])
      }
      CallGoDo("创建请求代码", {Data: array, Lang: Lang, Module: module})
    },
    NewColumnsLoaded(params) {
      this.onRangeSelectionChanged(params)
      if (this.ListFollowShow) {
        const rowCount = this.agGridApi.getDisplayedRowCount() - 1
        if (rowCount > -1) {
          this.rowIndex = rowCount
          this.agGridApi.ensureIndexVisible(rowCount)
        }
      }
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
      if (params.data.color) {
        if (params.data.color.selected) {
          if (params.data.color.selected.dark) {
            //设置选中行的高亮背景色
            res.backgroundColor = this.darkTheme ? params.data.color.selected.dark : params.data.color.selected.right
          }
        }
        if (!res.backgroundColor) {
          if (params.data.color.search) {
            if (params.data.color.search !== '') {
              //设置搜索的背景颜色
              res.backgroundColor = params.data.color.search
            }
          }
        }
        if (!res.backgroundColor) {
          if (params.data.color.back !== null && params.data.color.back !== void 0) {
            const mm = this.darkTheme ? params.data.color.back.dark : params.data.color.back.right
            if (mm) {
              if (mm !== '') {
                //设置背景色
                res.backgroundColor = this.darkTheme ? params.data.color.back.dark : params.data.color.back.right
              }
            }
          }
        }
        if (!res.color) {
          if (params.data.color.TagColor) {
            //设置标记的文本颜色
            res.color = params.data.color.TagColor
            res.fontWeight = 'bold';
          }
        }
        if (!res.color) {
          SetTextColor(params.data)
          if (params.data.color.text !== null && params.data.color.text !== void 0) {
            //设置常规文本颜色
            res.color = this.darkTheme ? params.data.color.text.dark : params.data.color.text.right
          }
        }
      }
      return res
    }
    ,
    markerColor(Color) {
      /*
      *       agGridApi: null,
            agSelectedArray: [],//选择的行数组
            agSelectedLine: null,//选中行
      * */
      let TagColor = ""
      if (Color === "红色") {
        TagColor = "#ff0000"
      } else if (Color === "蓝色") {
        TagColor = "#006fff"
      } else if (Color === "绿色") {
        TagColor = "#02bd02"
      } else if (Color === "黄色") {
        TagColor = "#fab200"
      } else if (Color === "紫色") {
        TagColor = "#bc00f6"
      } else if (Color === "empty") {
        for (let i = 0; i < this.TagColorMap.length; i++) {
          const obj = this.RowDataHashMap[this.TagColorMap[i]]
          if (obj) {
            obj.data.color.TagColor = null
          }
        }
        CallGoDo("标记颜色", {Data: this.TagColorMap, empty: true})
        this.TagColorMap = []
        this.RefreshRenderedNodes()
        return
      }
      let array = []
      for (let i = 0; i < this.agSelectedArray.length; i++) {
        array.push({id: this.agSelectedArray[i].data['Theology'], color: TagColor})
        this.agSelectedArray[i].data.color.TagColor = TagColor
        this.TagColorMap.push(this.agSelectedArray[i].data.Theology)
      }
      CallGoDo("标记颜色", {Data: array, empty: false})
      this.RefreshRenderedNodes()
    },
    delete() {
      const array = []
      const array2 = []
      for (let i = 0; i < this.agSelectedArray.length; i++) {
        array.push(this.agSelectedArray[i].data['Theology'])
        array2.push(this.agSelectedArray[i].data)
      }
      CallGoDo("删除请求会话", {Data: array}).then(res => {
        window.vm.List.agGridApi.applyTransaction({remove: array2});
        this.agSelectedArray = []
        this.agSelectedLine = null
      })
    },
    CloseSession() {
      const array = []
      for (let i = 0; i < this.agSelectedArray.length; i++) {
        array.push(this.agSelectedArray[i].data['Theology'])
      }
      CallGoDo("关闭请求会话", {Data: array}).then(res => {
        let ay = []
        for (let i = 0; i < this.agSelectedArray.length; i++) {
          if (this.agSelectedArray[i].data) {
            const zt = ("" + this.agSelectedArray[i].data['状态']).toUpperCase()
            if (zt.indexOf("已连接") !== -1) {
              this.agSelectedArray[i].data['状态'] = "断开中"
              ay.push(this.agSelectedArray[i].data)
            }
          }
        }
        window.vm.List.agGridApi.applyTransaction({update: ay});
      })
    },
    resend(mode) {
      //mode=3 普通重新发送
      //mode=1 重新发送并且拦截上行
      //mode=2 重新发送并且拦截下行
      const array = []
      for (let i = 0; i < this.agSelectedArray.length; i++) {
        array.push(this.agSelectedArray[i].data['Theology'])
      }
      CallGoDo("重发请求", {Data: array, Mode: mode})
    },
  }
  ,
  computed: {
    IsDarkTheme() {
      const event = new Event('darkThemeChange');
      window.dispatchEvent(event);
      {
        {
          if (this.darkTheme) {
            if (this.ListFollowShow) {
              this.MenuItems[8].icon = `<svg xmlns="http://www.w3.org/2000/svg" style="top: 2px;position: relative;" width="16" height="14" viewBox="0 0 24 24" fill="none" stroke="#fff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">\\n' +
              '    <polyline points="20 6 9 17 4 12"/>' +
              '</svg>`
            } else {
              this.MenuItems[8].icon = `<svg xmlns="http://www.w3.org/2000/svg" style="top: 2px;position: relative;" width="16" height="14" viewBox="0 0 24 24" fill="none" stroke="#fff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
    <line x1="18" y1="6" x2="6" y2="18"/>
    <line x1="6" y1="6" x2="18" y2="18"/>
</svg>`
            }
          } else {
            if (this.ListFollowShow) {
              this.MenuItems[8].icon = `<svg xmlns="http://www.w3.org/2000/svg" style="top: 2px;position: relative;" width="16" height="14" viewBox="0 0 24 24" fill="none" stroke="#000" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">\\n' +
              '    <polyline points="20 6 9 17 4 12"/>' +
              '</svg>`
            } else {
              this.MenuItems[8].icon = `<svg xmlns="http://www.w3.org/2000/svg" style="top: 2px;position: relative;" width="16" height="14" viewBox="0 0 24 24" fill="none" stroke="#000" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
    <line x1="18" y1="6" x2="6" y2="18"/>
    <line x1="6" y1="6" x2="18" y2="18"/>
</svg>`
            }
          }
        }
      }
      return "example-wrapper " + (this.darkTheme ? "ag-theme-balham-dark" : "ag-theme-balham")
    }
    ,
    FindWindowShow() {
      return window.UI.FindWindow
    },
  }
  ,
  mounted() {
    this.agGridApi = this.$refs.agGrid.gridOptions.api
    //this.HideToolsPanel();
    window.ToolsMaximize = false
    this.sideBar = {
      toolPanels: [
        {
          id: 'reqs',
          labelDefault: '请求数据',
          labelKey: 'customStats',
          iconKey: 'columns',
          toolPanel: 'ToolPanel',
          width: this.PageWidth.Width,
        },
        {
          id: 'columns',
          labelDefault: '列选择',
          labelKey: 'columns',
          iconKey: 'columns',
          toolPanel: 'agColumnsToolPanel',
          width: this.PageWidth.Width,
          toolPanelParams: {
            suppressRowGroups: true,
            suppressValues: true,
            suppressPivots: true,
            suppressPivotMode: true,
            suppressColumnFilter: true,
            suppressColumnSelectAll: true,
            suppressColumnExpandAll: true,
          },
        },
      ],
      defaultToolPanel: 'reqs',
    }
    window.addEventListener('resize', this.updateWindowSize);
    this.updateWindowSize()
    this.agTools = document.getElementsByClassName("ag-tool-panel-wrapper")
    document.addEventListener('mousemove', this.handleMouseMove);
    window.addEventListener('keydown', this.handleKeyDown);
    //删除 ：root 下的 --el-menu-bg-color 样式
    document.documentElement.style.setProperty('--el-menu-bg-color', 'unset');
    window.vm.List = this
    this.$nextTick(() => {
      {
        const obj1 = this.$refs.agGrid.$el.childNodes
        if (obj1.length > 0) {
          obj1.forEach((element) => {
            if (element && element.classList) {
              const ClassName = Array.from(element.classList).join(" ")
              if (ClassName === "ag-root-wrapper ag-layout-normal ag-ltr") {
                const obj2 = element.childNodes
                if (obj2) {
                  obj2.forEach((element2) => {
                    if (element2 && element2.classList) {
                      const ClassName2 = Array.from(element2.classList).join(" ")
                      if (ClassName2 === "ag-root-wrapper-body ag-layout-normal ag-focus-managed") {
                        const obj3 = element2.childNodes
                        if (obj3) {
                          obj3.forEach((element3) => {
                            if (element3 && element3.classList) {
                              const ClassName3 = Array.from(element3.classList).join(" ")
                              if (ClassName3 === 'ag-root ag-unselectable ag-layout-normal') {
                                this.ListDom = element3
                              }
                            }
                          })
                        }
                      }
                    }
                  })
                }
              }
            }
          });
        }
      }
      if (this.ListDom) {
        // 添加鼠标移入事件监听器
        this.ListDom.addEventListener('mouseenter', (event) => {
          this.IsListDomRange = true
        });
        // 添加鼠标移出事件监听器
        this.ListDom.addEventListener('mouseleave', (event) => {
          this.IsListDomRange = false
        });
      }
    })
    const columnFilter = this.agGridApi.getFilterInstance('响应长度');
    columnFilter.setModel({
      type: 'notContains',
      filter: '0/0'
    });
    const responseTypeFilter = this.agGridApi.getFilterInstance('响应类型');
    responseTypeFilter.setModel({
      type: 'notEqual',
      filter: 'error'
    });

    window.vm.List.agGridApi.onFilterChanged();
  }
  ,
  beforeUnmount() {
    window.removeEventListener('resize', this.updateWindowSize); // 移除 resize 事件监听器
    document.removeEventListener('mousemove', this.handleMouseMove); // 移除 mousemove 事件监听器
    window.removeEventListener('keydown', this.handleKeyDown); // 移除 keydown 事件监听器
  }
  ,
}

function getMenuItems(Items) {
  let array = [];
  for (let i = 0; i < Items.length; i++) {
    if (Items[i].visible !== false) {
      let Item = deepCopy(Items[i]);
      if (Item.subMenu) {
        let array1 = [];
        for (let n = 0; n < Item.subMenu.length; n++) {
          if (Item.subMenu[n].visible !== false) {
            array1.push(Item.subMenu[n])
          }
        }
        if (array1.length > 0) {
          if (array1[array1.length - 1] === "separator") {
            array1.pop()
          }
        }
        Item.subMenu = array1
      }
      if (array.length < 1 && Item === "separator") {
        continue
      }
      if (array[array.length - 1] === "separator" && Item === "separator") {
        continue
      }
      array.push(Item)
    }

  }
  return array
}
</script>

<template>
  <div class="no-select" STYLE="width: 100%;height: 100%">
    <div :class="IsDarkTheme">
      <div class="inner-col">
        <AgHeaderGroup/>
        <div class="inner-col2">
          <ag-grid-vue
              ref="agGrid"
              style="height: 100%;"
              :defaultColDef="defaultColDef"
              :rowData="RowData"
              :columnDefs="columns"
              :enableRangeSelection="true"
              :enableCharts="true"
              :modules="leftModules"
              :sideBar="sideBar"
              :grid-options="gridOptions"
              :overlayNoRowsTemplate="overlayNoRowsTemplate"
              :allowContextMenuWithControlKey="true"
              :suppressCutToClipboard="true"
              :sendToClipboard="sendToClipboard"
          >
          </ag-grid-vue>
          <AgFooterGroup/>
          <FindWindow v-show="FindWindowShow" :show="FindWindowShow"/>
        </div>
      </div>
    </div>
    <div class="ag-resizer-wrapper">
      <div ref="eTopLeftResizer" class="ag-resizer ag-resizer-topLeft" style="pointer-events: all;z-index: 1000"></div>
      <div ref="eTopResizer" class="ag-resizer ag-resizer-top" style="pointer-events: all;z-index: 1000"></div>
      <div ref="eTopRightResizer" class="ag-resizer ag-resizer-topRight"
           style="pointer-events: all;z-index: 1000"></div>
      <div ref="eRightResizer" class="ag-resizer ag-resizer-right" style="pointer-events: all;z-index: 1000"></div>
      <div ref="eBottomRightResizer" class="ag-resizer ag-resizer-bottomRight"
           style="pointer-events: all;z-index: 1000"></div>
      <div ref="eBottomResizer" class="ag-resizer ag-resizer-bottom" style="pointer-events: all;z-index: 1000"></div>
      <div ref="eBottomLeftResizer" class="ag-resizer ag-resizer-bottomLeft"
           style="pointer-events: all;z-index: 1000"></div>
      <div ref="eLeftResizer" class="ag-resizer ag-resizer-left" style="pointer-events: all;z-index: 1000"></div>
    </div>
  </div>
</template>

<style scoped>

.ag-cell-focus, .ag-cell-no-focus {
  border: none !important;
}

/* This CSS is to not apply the border for the column having 'no-border' class */
.no-border.ag-cell:focus {
  border: none !important;
  outline: none;
}

.example-wrapper {
  display: flex;
  height: 100%;
  width: 100%;
  flex: 1 1;
  gap: 50px;
}

.inner-col {
  width: 100%;
  height: 100%;
}

.inner-col2 {
  width: 100%;
  height: calc(100% - 60px);
}


</style>

<style>
/* 设置悬停背景色为黄色，并提高优先级 */
:root {

}

.ag-theme-balham {
  --ag-row-hover-color: #024d6c;
}

.ag-theme-balham-dark {
  --ag-row-hover-color: rgb(139, 207, 238);
}
</style>