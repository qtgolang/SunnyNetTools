<script>
import {AgGridVue} from "ag-grid-vue3";
import {
  Config_agGrid_API,
  Config_AutoRoll,
  Config_Find_Window_Hide,
  Config_Focus_Element,
  Config_HomeTextMark,
  Config_IsDark,
  Config_IsShowExportProgress,
  Config_Menu_isFileMenu,
  Config_MenuVisible,
  Config_SelectedRow,
  Config_SunnyNetIsStart,
  Config_Status_Info,
  Config_HTTP_Message_free,
  Config_Theme_agGrid,
  DisableClick,
  getThisObject,
  registerThisObject,
} from "./config/Config.js";
import {
  AppCheckSunnyNet,
  AppDeleteSession,
  AppDisconnectTCPRequest,
  AppExport,
  AppGenerateCode,
  AppImport,
  AppResendRequest,
  FreeAllRequest,
  GetColumnState,
  GetPort,
  GOOS,
  ListSearch,
  McpFuncRes,
  OpenSunnyFile,
  SaveSunnyFile,
  SetColumnState,
  SetRequestNextBreakMode,
  UpdateNote
} from "../../bindings/changeme/Service/appmain.js";
import {AG_GRID_LOCALE_CN} from "./config/AG_ZH_CN.js";
import ListenOn from "./config/ListenOn.vue";
import reqsPanel from './Home/reqsPanel.vue';
import HomeFooter from './Home/Footer/Footer.vue';
import Tools from './Tools/Tools.vue';
import Settings from './SideBar/Settings.vue';
import Device from './SideBar/Device.vue';
import ScriptLog from './SideBar/ScriptLog.vue';
import Header from "./Home/Header/Header.vue";
import ImageRenderer from "./Home/imageRenderer.vue";
import FindWindow from "./Tools/Find/FindWindow.vue";
import Filter from "./Tools/Filter/filter.vue";
import {keydownEventListener, Keys_System_id_ResendRequest, registerHotkeyFunction} from "./config/Keys";
import {GetGenerateCodeListMenu} from "./config/GenerateCode";
import {ElLoading, ElMessageBox, ElNotification} from "element-plus";
import {nextTick} from "vue";
import TitleBar from "./TitleBar/TitleBar.vue";
import VTitlebar from "./TitleBar/VUETitlebar/vueTitlebar.vue";
import {OpenTools} from "./CallbackEventsOn";
import {Events} from "@wailsio/runtime";

export default {
  components: {
    VTitlebar,
    TitleBar,
    Filter,
    Settings,
    FindWindow,
    Header,
    'ag-grid-vue': AgGridVue,
    'toolsSideBar': Tools,
    'ScriptLogSideBar': ScriptLog,
    reqsPanel,
    ImageRenderer,
    HomeFooter,
    ListenOn,
    Device
  },
  data() {
    return {
      agGridApi: null,
      agGridApiMain: null,
      agSelectedRowNodes: [],
      rowData: [],
      colDefs: [],
      isHideHook: false,
      AdvancedFilterModel: null,
      mainFilter: null,
      isGoFilter: false,//是否使用Go过滤
      GoFilterList: new Set([]),
      isSearch: false,  //是搜索结果
      GoSearchList: new Set([]),
      HideHookList: new Set([]),
      defaultColDef: {
        flex: 1,
        sortable: false,
        resizable: true,
        suppressHeaderMenuButton: true,
      },
      MoveColumn: false,
      overlayNoRowsTemplate: `<span style="padding: 20px;" id="HookMessageText">还没有捕获到数据</span>`,
      sideBar: {
        toolPanels: [
          {
            id: 'reqs',
            labelDefault: '请求数据',
            labelKey: 'customStats',
            iconKey: 'linked',
            toolPanel: 'reqsPanel',
            minWidth: 400,
            width: 600,
          },
          {
            id: 'columns',//这个列看不见,因为我隐藏了面板，所以无所谓了,但是不能删，删了就不对，如果删除了，请求数据面板就，找不到一些样式
            labelDefault: '列调整',
            labelKey: 'columns',
            iconKey: 'columns',
            toolPanel: 'agColumnsToolPanel',
            width: 200,
            minWidth: 200,
            maxWidth: 200,
            hide: true,
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
        defaultToolPanel: '',//reqs
      },
      FilterColumn: [],
      ResendRequestCountShow: false,
      ResendRequestCount: 1,
      IsWindows: false,
      gridOptions: {
        enableAdvancedFilter: true,
        suppressMovableColumns: true,
        onContextMenuVisibleChanged: (event) => {
          Config_MenuVisible.value = event.visible;
        },
        getRowId: (params) => {
          return String(params.data.Theology)
        },
        stopEditingWhenCellsLoseFocus: true, // 失去焦点时自动结束编辑
        columnDefs: [
          {
            field: "序号", tooltipField: '序号',
            minWidth: 85,
            width: 85,
            maxWidth: 85,
            menuTabs: [], // 不显示过滤器
            pinned: "left",
            valueGetter: (params) => {
              if (!this.isGoFilter) {
                return params.node.rowIndex + 1;
              }
              let pso = params.node.data.rowIndex;
              if (pso) {
                return pso + 1;
              }
              this.agGridApi.forEachNode((node, index) => {
                if (node === params.node) {
                  params.node.data.rowIndex = index;
                }
              });
              return params.node.data.rowIndex + 1;
            },
            cellRenderer: 'ImageRenderer',
            cellStyle: {'text-align': 'left'},
            suppressFiltersToolPanel: true,
            suppressColumnsToolPanel: true,
          },
          {
            field: 'Theology',
            hide: true, // 不显示
            suppressColumnsToolPanel: true, // 不在“列面板”中显示
            suppressFiltersToolPanel: true, // 不在“过滤面板”中显示
            suppressMovable: true,          // 不可被拖拽
            suppressSizeToFit: true,        // 不随容器缩放
            suppressHeaderKeyboardEvent: true, // 禁用 header 的交互
            suppressHeaderContextMenu: true,   // 禁止右键菜单
            sortable: true, //（可选）仍然可以参与排序
            filter: false   //（可选）不允许筛选
          },
          {
            field: "方式", tooltipField: '方式',
            minWidth: 80,
            width: 100,
            filter: true,
            maxWidth: 120, cellStyle: {'text-align': 'left'}
          },
          {
            field: "状态", tooltipField: '状态',
            minWidth: 90,
            width: 90,
            filter: true,
            maxWidth: 90, cellStyle: {'text-align': 'center'}
          },
          {
            field: "主机名", width: 120, minWidth: 120, tooltipField: '主机名', hide: true,
            filter: true,
            maxWidth: 500, cellStyle: {'text-align': 'left'}
          },
          {
            field: "路径", width: 200, minWidth: 200, tooltipField: '路径', hide: true,
            filter: true,
            maxWidth: 2000, cellStyle: {'text-align': 'left'}
          },
          {
            field: "请求地址", width: 400, minWidth: 200, tooltipField: '请求地址',
            filter: true,
            maxWidth: 2000, cellStyle: {'text-align': 'left'},
          },
          {
            field: "响应长度", width: 120, minWidth: 120, tooltipField: '响应长度',
            filter: true,
            maxWidth: 120, cellStyle: {'text-align': 'left'},
          },
          {
            field: "响应类型", width: 120, minWidth: 120, tooltipField: '响应类型',
            filter: true, hide: true,
            maxWidth: 120, cellStyle: {'text-align': 'left'}
          },
          {
            field: "进程", width: 100, minWidth: 100, tooltipField: '进程',
            filter: true,
            maxWidth: 500, cellStyle: {'text-align': 'left'}
          },
          {
            field: "注释", width: 200, minWidth: 200, tooltipField: '注释', editable: true,
            filter: true,
            maxWidth: 2000, cellStyle: {'text-align': 'left'}
          },
          {
            field: "请求时间", width: 150, minWidth: 150, tooltipField: '请求时间', hide: true,
            filter: true,
            maxWidth: 150, cellStyle: {'text-align': 'left'}
          },
          {
            field: "响应时间", width: 150, minWidth: 150, tooltipField: '响应时间', hide: true,
            filter: true,
            maxWidth: 150, cellStyle: {'text-align': 'left'}
          },
          {
            field: "来源地址", width: 150, minWidth: 150, tooltipField: '来源地址',
            filter: true,
            maxWidth: 200,
          },
          {
            field: "响应IP", width: 150, minWidth: 150, tooltipField: '响应IP',
            filter: true,
            maxWidth: 150, cellStyle: {'text-align': 'left'}
          },
          {
            field: "身份验证账号", width: 150, minWidth: 150, tooltipField: '身份验证账号', hide: true,
            filter: true,
            maxWidth: 150, cellStyle: {'text-align': 'left'}
          },
          {
            field: "参数", width: 150, minWidth: 150, tooltipField: '参数', hide: true,
            filter: true,
            maxWidth: 2000, cellStyle: {'text-align': 'left'}
          },
        ],
        localeText: AG_GRID_LOCALE_CN,
        getContextMenuItems: this.MenuEvent,
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
            if (selectedRowNodes.length > 0) {
              if (!___c) {
                Config_SelectedRow.value = selectedRowNodes[0].data
              } else {
                Config_SelectedRow.value = selectedRowNodes[selectedRowNodes.length - 1].data
              }
              this.agGridApi.openToolPanel("reqs");
            } else {
              Config_SelectedRow.value = [];
            }
            this.agSelectedRowNodes = selectedRowNodes;
          }
        },
        onModelUpdated: this.onModelUpdated,
        getRowStyle: this.onGetRowStyle,
        isExternalFilterPresent: this.isExternalFilterPresent,
        doesExternalFilterPass: this.doesExternalFilterPass,
        onColumnVisible: this.onColumnVisible,
        onColumnMoved: this.saveColumn,
        onColumnResized: this.saveColumn,
        onCellEditingStopped: (params) => {
          const _Theology = parseInt(params.data["Theology"])
          UpdateNote(_Theology, params.data["注释"]).then((isShow) => {
            this.addFilter(_Theology, isShow)
            this.agGridApi.applyTransaction({update: [params.data]});
          })
        },
        onFilterChanged: () => {
          this.agGridApi.refreshCells({
            columns: ['序号'],
            force: true
          });
        }
      },
      gridSideBar: {},
      cellSelection: {
        suppressMultiRanges: true,
      },
      rowSelection: {
        mode: 'multiRow',
        checkboxes: false,
        headerCheckbox: false,
        copySelectedRows: false,
        enableClickSelection: true,
      },
      IsShowNoRowsOverlay: true,//是否显示没有行的提示信息
      debounceTimer: null,
      textMark: new Map(),
      get isDisableClick() {
        return DisableClick.value
      },
      set isDisableClick(v) {
        DisableClick.value = v
      },
      get IsDark() {
        return Config_IsDark.value
      },
      set IsDark(v) {
        Config_IsDark.value = v
      }
    }
  },
  computed: {
    agTheme() {
      return Config_Theme_agGrid.value
    },
    isNoDisableClick() {
      return (this.isDisableClick ? "pointer-events: none;" : "") + "width: 100%;height: 100%"
    },
  },
  mounted() {
    this.agGridApi = this.$refs.agGrid.api;
    this.agGridApiMain = this.$refs.agGridSide.api
    {

      GOOS().then((IsWindows) => {
        this.IsWindows = IsWindows;
        let ToolsWidth = 200;
        let ToolsMaxWidth = 400;
        if (!IsWindows) {
          ToolsWidth = 100;
          ToolsMaxWidth = 100;
        }
        const sideBar = {
          toolPanels: [
            {
              id: 'Settings',
              labelDefault: '设置',
              labelKey: 'Settings',
              iconKey: 'settings',
              toolPanel: 'Settings',
              minWidth: 400,
              width: 400,
              maxWidth: 400,
            },
            {
              id: 'Tools',
              labelDefault: '常用工具',
              labelKey: 'customStats',
              iconKey: 'aggregation',
              toolPanel: 'toolsSideBar',
              minWidth: ToolsWidth,
              width: ToolsWidth,
              maxWidth: ToolsMaxWidth,
            },
            {
              id: 'GoCodeLog',
              labelDefault: '脚本日志',
              labelKey: 'customStats',
              iconKey: 'eye',
              toolPanel: 'ScriptLogSideBar',
              minWidth: 40,
              width: 400,
              maxWidth: 400,
            },
          ],
          position: 'left',
          defaultToolPanel: '',//默认不打开，如果需要打开 例如添加 id [Tools]
        }
        sideBar.toolPanels.push(
            {
              id: 'Device',
              labelDefault: '驱动加载',
              labelKey: 'customStats',
              iconKey: 'linked',
              toolPanel: 'Device',
              minWidth: 400,
              width: 400,
              maxWidth: 400,
            })
        this.agGridApiMain.setGridOption('sideBar', sideBar);
      })
    }
    Config_agGrid_API.value = this.$refs.agGrid.api;
    this.agGridApi.setSideBarPosition('right')
    this.onColumnVisible()
    registerThisObject("registerDropFiles", (el, path) => {
      if (!path.toLowerCase().endsWith(".sy4")) {
        ElNotification({
          position: 'bottom-right',
          message: '打开记录文件失败\n\n文件格式不正确！',
          type: 'warning',
          customClass: 'multiline-message'
        })
        return;
      }
      this.openSy4File(path)
    })
    registerThisObject("HomeListRefresh", this.RefreshVisibleNodes)
    registerThisObject("ColorMarking", () => this.ColorMarking('Clear'))
    registerThisObject("MCPRowRefresh", (theology) => {
      const th = parseInt(theology, 10);
      const cur = parseInt(Config_SelectedRow.value?.Theology ?? "0", 10);
      if (!Number.isFinite(th) || th !== cur) {
        return;
      }
      const api = Config_agGrid_API.value;
      const node = api?.getRowNode(String(th));
      if (!api || !node) {
        return;
      }
      api.deselectAll();
      api.clearCellSelection();
      this.agSelectedRowNodes = [];
      Config_SelectedRow.value = [];
      this.$nextTick(() => {
        node.setSelected(true);
        this.agSelectedRowNodes = [node];
        Config_SelectedRow.value = node.data;
        api.ensureNodeVisible(node);
      });
    })
    registerHotkeyFunction(Keys_System_id_ResendRequest, this.OpenResendRequest)
    {
      // 监听键盘按键事件
      document.addEventListener("keydown", (event) => {
        const placeholder = event.target.getAttribute('placeholder') + ""
        if (placeholder === "请按下快捷键") {
          return;
        }
        // 检测是否按下 Ctrl + F 组合键
        if (event.ctrlKey && event.key === "f") {
          // 如果焦点在 Monaco Editor 中，则不执行后续逻辑
          if (event.target.closest(".monaco-editor")) {
            return;
          }
          // 获取需要聚焦的输入元素（如果不存在，则调用 setFocusElement 获取）
          const focusElement = Config_Focus_Element.value ?? this.$refs.Filter.setFocusElement();
          // 如果找到了可用的输入框，则进行聚焦并触发点击事件
          if (focusElement) {
            focusElement.focus();
            focusElement.click();
          }
          return;
        }
        keydownEventListener(event)
      });
    }
    {
      // 获取 AG Grid 的根容器
      const gridRoot = document.querySelector("[data-ref='eGridRoot']");
      if (gridRoot) {
        // 清空右侧内容，仅保留左侧面板
        gridRoot.innerHTML = "";
        gridRoot.appendChild(this.$refs.agGrid.$el);
        gridRoot.id = "appMain"; // 重新赋 ID
        this.agGridApi.setGridOption("popupParent", gridRoot.parentElement)
      }
      // 隐藏侧边栏按钮（ag-side-buttons）
      const sideButtons = this.$refs.agGrid.$el.getElementsByClassName("ag-side-buttons");
      for (const buttonContainer of sideButtons) {
        if (buttonContainer.childNodes.length > 0) {
          buttonContainer.style.display = "none";
        }
      }
    }
    setTimeout(this.isInstallCert, 10)
    const style = document.createElement('style');
    style.textContent = `
  .ag-chart-tabbed-menu-body::after {
    display: none !important;
  }
`;
    document.head.appendChild(style);
    Events.On("mcp", async (evt) => {
      const mcp = evt?.data ?? {};
      // 统一回包：确保只要进到 handler，就能在需要时返回
      const reply = (text) => {
        mcp.res = text;
        typeof McpFuncRes === "function" && McpFuncRes(mcp);
      };
      try {
        const page = String(mcp.page ?? "").toLowerCase();
        const tag = String(mcp.tag ?? "").toLowerCase();
        if (page !== "main") return;
        switch (tag) {
          case "rowrefresh": {
            const body = JSON.parse(String(mcp.msg ?? "{}"));
            const th = parseInt(body.theology, 10);
            const fn = getThisObject("MCPRowRefresh");
            if (typeof fn === "function" && Number.isFinite(th)) {
              fn(th);
            }
            return reply("ok");
          }
          case "rowupdatesend": {
            const body = JSON.parse(String(mcp.msg ?? "{}"));
            const fn = getThisObject("MCPApplyHttpSendRowUpdate");
            if (typeof fn === "function") {
              fn(body);
            }
            return reply("ok");
          }
          case "rowmark": {
            const body = JSON.parse(String(mcp.msg ?? "{}"));
            const mark = body.mark ?? "";
            const ids = Array.isArray(body.ids) ? body.ids : [];
            for (const id of ids) {
              const th = parseInt(id, 10);
              if (Number.isFinite(th) && th !== 0) {
                this.textMark.set(th, mark);
              }
            }
            this.RefreshVisibleNodes();
            return reply("已标记 " + ids.length + " 行");
          }
          case "rownote": {
            const body = JSON.parse(String(mcp.msg ?? "{}"));
            const note = body.note ?? "";
            const ids = Array.isArray(body.ids) ? body.ids : [];
            const api = Config_agGrid_API.value;
            if (!api) {
              return reply("已更新 0 行注释");
            }
            const updates = [];
            for (const id of ids) {
              const th = parseInt(id, 10);
              if (!Number.isFinite(th) || th === 0) {
                continue;
              }
              const rowNode = api.getRowNode(String(th));
              if (!rowNode) {
                continue;
              }
              rowNode.data["注释"] = note;
              updates.push(rowNode.data);
              UpdateNote(th, note).then((isShow) => {
                this.addFilter(th, isShow);
              });
            }
            if (updates.length > 0) {
              api.applyTransaction({update: updates});
            }
            return reply("已更新 " + updates.length + " 行注释");
          }
          case "listindextotheology": {
            let want = null;
            try {
              const parsed = JSON.parse(String(mcp.msg ?? "[]"));
              if (Array.isArray(parsed)) {
                want = new Set(
                    parsed.map((x) => parseInt(x, 10)).filter((n) => n > 0),
                );
              }
            } catch (_) {
              return reply("{}");
            }
            const out = {};
            const api = Config_agGrid_API.value;
            if (!api) {
              return reply("{}");
            }
            // 遍历全部行（含被过滤器隐藏的行），序号与 rowIndex+1 一致
            api.forEachNode((node) => {
              const listIndex = node.rowIndex + 1;
              if (want === null || want.size === 0 || want.has(listIndex)) {
                const th = parseInt(node.data["Theology"], 10);
                out[listIndex] = Number.isFinite(th) ? th : 0;
              }
            });
            if (want !== null && want.size > 0) {
              for (const k of want) {
                if (!(k in out)) {
                  out[k] = 0;
                }
              }
            }
            return reply(JSON.stringify(out));
          }
          case "breakrelease": {
            const body = JSON.parse(String(mcp.msg ?? "{}"));
            const ids = Array.isArray(body.ids) ? body.ids : [];
            const mode = parseInt(body.mode ?? 0, 10);
            const api = Config_agGrid_API.value;
            for (const id of ids) {
              const th = parseInt(id, 10);
              if (!Number.isFinite(th)) {
                continue;
              }
              const rowNode = api?.getRowNode(String(th));
              if (rowNode?.data) {
                rowNode.data["断点模式"] = mode;
                SetRequestNextBreakMode(th, mode);
              }
            }
            if (ids.length > 0) {
              Config_HTTP_Message_free(
                  ids.map((x) => parseInt(x, 10)).filter((n) => Number.isFinite(n)),
              );
            }
            api?.refreshCells({columns: ["断点模式"], force: true});
            return reply("ok");
          }
          case "breakreleaseall":
            FreeAllRequest();
            Config_HTTP_Message_free(Config_SelectedRow.value ?? []);
            Config_agGrid_API.value?.refreshCells({columns: ["断点模式"], force: true});
            return reply("ok");
          case "enginestatus": {
            const msg = JSON.parse(String(mcp.msg ?? "{}"));
            if (msg.running != null) {
              Config_SunnyNetIsStart.value = !!msg.running;
            }
            if (msg.statusText != null) {
              Config_Status_Info.value = String(msg.statusText);
            }
            return reply("ok");
          }
          case "recordsimport": {
            const body = JSON.parse(String(mcp.msg ?? "{}"));
            const rows = Array.isArray(body.rows) ? body.rows : [];
            if (rows.length > 0 && this.$refs.listen) {
              this.$refs.listen.insertArray(rows, true, () => {});
            }
            return reply("ok");
          }
          case "delreq":
            const msg = JSON.parse(String(mcp.msg ?? ""));
            const arrayID = []
            const array = []
            for (const item of msg) {
              arrayID.push(parseInt(item))
              const rowNode = Config_agGrid_API.value.getRowNode(item + "");
              if (rowNode) {
                array.push(rowNode.data)
                rowNode.setSelected(false);
              }
            }
            AppDeleteSession(arrayID).then(() => {
              this.agGridApi.applyTransaction({remove: array});
              this.agSelectedRowNodes = []
              Config_SelectedRow.value = [];
              this.agGridApi.clearCellSelection();
              this.agGridApi.refreshCells({
                columns: ['序号'],
                force: true,
                suppressFlash: true, // 不需要闪烁动画，可以不加
              });
            })
            return reply("成功删除 " + arrayID.length + "条 记录");
          default:
            return;
        }
      } catch (e) {
        try {
          reply("处理失败");
        } catch (_) {
        }
      }
    });
  },
  methods: {
    setIsInstallCertDrag() {
      requestAnimationFrame(() => {
        const messageBox = document.getElementsByClassName("el-overlay-message-box");
        if (messageBox.length > 0) {
          if (messageBox[0]) {
            messageBox[0].style = "--wails-draggable: drag"
            return
          }
        }
        this.setIsInstallCertDrag()
      })
    },
    isInstallCert() {
      requestAnimationFrame(() => {
        const obj = document.getElementById('rootList')
        const loading = ElLoading.service({
          lock: true,
          text: '正在检测是否已经安装 SunnyNet 证书',
          background: 'rgba(0, 0, 0, 0.7)',
          target: obj,
        })
        AppCheckSunnyNet().then((isInstall) => {
          if (isInstall) {
            loading.close()
          } else {
            this.setIsInstallCertDrag()
            ElMessageBox.confirm('当前系统中未发现 SunnyNet 证书<br>部分功能可能异常！！<br>是否立即查看证书安装教程？', '证书未安装', {
              // if you want to disable its autofocus
              // autofocus: false,
              confirmButtonText: '立即查看',
              cancelButtonText: '关闭这个提醒',
              dangerouslyUseHTMLString: true,
              closeOnClickModal: false,     // 禁止点击遮罩关闭
              closeOnPressEscape: false,    // 可选：禁止按 ESC 关闭
              showClose: false              // 可选：不显示右上角关闭按钮
            }).then(() => {
              GetPort().then(port => {
                const url = 'http://localhost:' + port + '/install.html'
                OpenTools("证书安装", true, url)
                ElNotification({
                  position: 'bottom-right',
                  message: 'SunnyNet 证书未安装！\n\n部分功能可能异常！！',
                  type: 'warning',
                  customClass: 'multiline-message'
                })
                loading.close()
              })
            }).catch(() => {
              ElNotification({
                position: 'bottom-right',
                message: 'SunnyNet 证书未安装！\n\n部分功能可能异常！！',
                type: 'warning',
                customClass: 'multiline-message'
              })
              loading.close()
            })


          }
        })
      })
    },
    ExportMessage(ids) {
      SaveSunnyFile("").then((selectedFiles, err) => {
        if (selectedFiles === "") {
          return
        }
        ElNotification({
          position: 'bottom-right',
          showClose: true,
          message: '正在保存记录到文件...\n这可能需要一定时间...',
          type: 'success',
          customClass: 'multiline-message'
        })
        Config_IsShowExportProgress.value = "请稍后,正在保存记录"
        AppExport(ids, selectedFiles).then((err) => {
          Config_IsShowExportProgress.value = ""
          if (err === "") {
            ElNotification({
              position: 'bottom-right',
              message: '导出成功！',
              type: 'success',
              customClass: 'multiline-message'
            })
            return
          }
          ElNotification({
            position: 'bottom-right',
            message: '导出失败！\n\n' + err,
            type: 'warning',
            customClass: 'multiline-message'
          })
        }).catch((error) => {
          ElNotification({
            position: 'bottom-right',
            message: '导出失败！',
            type: 'warning',
            customClass: 'multiline-message'
          })
        });
      }).catch((error) => {
        ElNotification({
          position: 'bottom-right',
          message: '保存记录失败\n\n未选择文件！',
          type: 'warning',
          customClass: 'multiline-message'
        })
      });
    },
    setSelectedRow(addRows, GuaranteeDisplay) {
      // 清除之前选中的行
      this.agSelectedRowNodes.forEach((node) => {
        node.setSelected(false);
      });
      const rangeSelections = this.agGridApi.getCellRanges();
      const colIds = [];
      rangeSelections.forEach(function (range) {
        range.columns.forEach(function (column) {
          colIds.push(column.colId);
        });
      });

      let targetNode = null;
      for (let i = 0; i < addRows.length; i++) {
        const node = addRows[i];
        if (node.data.Theology === GuaranteeDisplay) {
          targetNode = node;
          this.agSelectedRowNodes = [node];
          Config_SelectedRow.value = node.data;
          this.agGridApi.ensureNodeVisible(node);
          node.setSelected(true);
          break;
        }
      }
      if (!targetNode) return;
      this.agGridApi.clearCellSelection();
      this.agGridApi.addCellRange({
        rowStartIndex: targetNode.rowIndex,
        rowEndIndex: targetNode.rowIndex,
        columns: colIds,
      });
      this.agGridApi.setFocusedCell(targetNode.rowIndex, colIds);
    },
    saveColumn() {
      if (Config_SunnyNetIsStart.value) {
        clearTimeout(this.debounceTimer);
        this.debounceTimer = setTimeout(() => {
          const colState = this.agGridApi.getColumnState();
          SetColumnState(JSON.stringify(colState))
        }, 1000);
      }
    },
    CancelSearch() {
      this.isSearch = false;
      this.GoSearchList.clear();
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
        const mm = new Map();
        for (let i = 0; i < params.length; i++) {
          const node = this.agGridApi.getRowNode(params[i] + "");
          if (node) {
            node.data.backcolor = Color;
            array.push(node)
            mm.set(parseInt(node.data["Theology"]), node)
          }
        }
        this.agGridApi.forEachNode((node) => {
          if (!mm.has(parseInt(node.data["Theology"]))) {
            node.data.backcolor = "";
            array.push(node)
          }
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
        ListSearch("clear").then((array) => {
          this.isGoFilter = false;
          this.agGridApi.onFilterChanged()
          this.agGridApi.refreshCells({
            force: true,               // 强制刷新所有字段
            suppressFlash: true        // 不闪烁高亮（可选）
          });
        });
        return;
      }
      this.agGridApi.forEachNode((node, index) => {
        node.data.rowIndex = index;
      });
      this.isGoFilter = true;
      this.GoFilterList.clear()
      ListSearch(JSON.stringify(rules)).then((array) => {
        array.forEach((Theology) => {
          this.GoFilterList.add(Theology)
        })
        this.agGridApi.onFilterChanged()
      });
    },
    HideHook(params) {
      this.isHideHook = params
      if (params) {
        this.AdvancedFilterModel = this.$refs.Filter.agGridApi.getAdvancedFilterModel();
        this.HideHookList.clear()
        this.agGridApi.forEachNode((node) => {
          const theology = parseInt(node.data["Theology"]);
          this.HideHookList.add(theology)
        });
        this.agGridApi.onFilterChanged();
      } else {
        this.HideHookList.clear()
        this.onFilterApply(this.AdvancedFilterModel)
      }
    },
    addFilter(__Theology__, Filter) {
      if (!this.isGoFilter || Filter == null) return;
      const theology = parseInt(__Theology__);
      const exists = this.GoFilterList.has(theology);
      if (Filter) {
        if (!exists) {
          this.GoFilterList.add(theology);
        }
      } else {
        if (exists) {
          this.GoFilterList.delete(theology);
        }
      }
    },
    isExternalFilterPresent(params) {
      return this.isGoFilter || this.isHideHook || this.isSearch;
    },
    doesExternalFilterPass(node) {
      if (this.isHideHook) {
        if (!this.HideHookList.has(parseInt(node.data['Theology']))) {
          return false;
        }
      }
      if (this.isSearch) {
        return this.GoSearchList.has(parseInt(node.data['Theology']));
      }
      if (!this.isGoFilter) {
        return true;
      }
      return this.GoFilterList.has(parseInt(node.data['Theology']));
    },
    onModelUpdated() {
      const visibleRowCount = this.agGridApi.getDisplayedRowCount();//可见行
      if (visibleRowCount === 0) {
        let totalRowCount = 0;//总行数
        this.agGridApi.forEachNode(() => {
          totalRowCount++
        });
        let Template = "";
        if (totalRowCount !== 0) {
          if (this.isHideHook) {
            Template = `<span style="padding: 20px;" id="HookMessageText">有数据,但您隐藏了捕获数据</span>`;
          } else {
            Template = `<span style="padding: 20px;" id="HookMessageText">根据当前过滤器中条件：没有符合条件的数据</span>`;
          }
        } else {
          this.agGridApi.hideOverlay();
          Template = `<span style="padding: 20px;" id="HookMessageText">还没有捕获到数据</span>`;
        }
        this.overlayNoRowsTemplate = Template;
        this.agGridApi.showNoRowsOverlay();
        this.IsShowNoRowsOverlay = true
      } else {
        //有行
        if (this.IsShowNoRowsOverlay) {
          //当前正在显示没有行的提示信息

          //不显示没有行的提示信息，因为有行
          this.IsShowNoRowsOverlay = false;
          //执行API隐藏提示
          this.agGridApi.hideOverlay();
        }
      }
    },
    onGetRowStyle(params) {
      let res = {
        fontFamily: "微软雅黑"
      }
      if (this.IsWindows) {
        res.fontFamily = `-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", "Helvetica Neue", Helvetica, Arial, sans-serif`;
      }
      const id = parseInt(params.data['Theology'])
      if (this.textMark.has(id)) {
        const ColorID = this.textMark.get(id)
        if (Config_HomeTextMark.has(ColorID)) {
          const value = Config_HomeTextMark.get(ColorID)
          res.color = Config_IsDark.value ? value["深色主题"] : value["浅色主题"]
          res.fontWeight = 'bold';
          return res
        }
      }
      if (params.data.color) {
        res.color = params.data.color
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
    onGridReady(params) {
      this.initColumnWidth()
      //向高级过滤器添加一个按钮
      {
        const customButtonContainer = this.$el.getElementsByClassName('ag-advanced-filter ag-advanced-filter-header-cell');
        let mm = null;
        // 由于这里是 home 也就是主页面,也会有其他的高级过滤器,但是刚加载的时候肯定只有当前主列表会显示,其他的都是隐藏的,所以获取不到宽度
        // 并且主窗口的高级过滤器宽度肯定是最长的,所以这里获取宽度最长的一个
        for (let i = 0; i < customButtonContainer.length; i++) {
          if (mm == null) {
            mm = customButtonContainer[i]
            continue
          }
          const width1 = mm.getBoundingClientRect().width;
          const width2 = customButtonContainer[i].getBoundingClientRect().width;
          if (width1 < width2) {
            mm = customButtonContainer[i]
          }
        }
        if (mm !== null) {
          mm.style.display = "none";
          this.mainFilter = mm.parentElement
        }
      }
    },
    initColumnWidth() {
      const ColumnInfo = JSON.parse(JSON.stringify(this.gridOptions.columnDefs))
      this.initColumnWidthRunTask(ColumnInfo)
    },
    initColumnWidthRunTask(ColumnInfo) {
      requestAnimationFrame(() => {
        if (!this.agGridApi) return;
        GetColumnState().then(res => {
          let colState = [];
          try {
            if (res === "") {
              colState = this.agGridApi.getColumnState();
            } else {
              colState = JSON.parse(res)
              this.agGridApi.applyColumnState({
                state: colState,
                applyOrder: true,
              });
              return
            }
          } catch (e) {
            colState = this.agGridApi.getColumnState();
          }
          const colDefMap = new Map(ColumnInfo.map(col => [col.field, col.width, col.hide]));
          let hasChanges = false;
          for (let col of colState) {
            const newWidth = colDefMap.get(col.colId);
            const newHide = colDefMap.get(col.hide);
            if (newWidth !== undefined && col.width !== newWidth) {
              col.width = newWidth;
              col.hide = newHide;
              hasChanges = true;
            }
          }
          if (hasChanges) {

          } else {
            setTimeout(() => this.initColumnWidthRunTask(ColumnInfo), 100);
          }
        })

      });
    },
    handleClick(event) {
      this.agGridApiMain.closeToolPanel()
      this.$refs.Filter.setFocusElement();
      if (event.target) {
        if (event.target.className === "ag-menu-option-part ag-menu-option-text" && event.target.innerText === "选择列") {
          this.MoveColumnEvent()
          return;
        }
        if ((event.target.className + "").indexOf("ag-viewport") !== -1) {
          //如果点击没有的行，隐藏面板
          try {
            const parentElement = event.target.parentNode.parentNode.parentNode.parentNode.childNodes;
            const ariaHiddenValue = parentElement[6].getAttribute('aria-hidden');
            if (ariaHiddenValue === "false") {
              this.agGridApi.closeToolPanel()
              this.agSelectedRowNodes.forEach((node) => {
                node.setSelected(false);
              });
              Config_SelectedRow.value = [];
              this.agSelectedRowNodes = [];
            }
          } catch (e) {
          }
          return;
        }
      }
    },
    RefreshVisibleNodes() {
      const visibleNodes = this.agGridApi.getRenderedNodes();
      this.agGridApi.redrawRows({rowNodes: visibleNodes});
    },
    ColorMarking(ColorID) {
      if (ColorID === "Clear") {
        this.textMark.clear()
        this.RefreshVisibleNodes()
        return
      }
      for (let i = 0; i < this.agSelectedRowNodes.length; i++) {
        const id = parseInt(this.agSelectedRowNodes[i].data['Theology'])
        this.textMark.set(id, ColorID)
      }
      this.RefreshVisibleNodes()
    },
    SubResendRequest() {
      this.ResendRequestCountShow = false
      this.ResendRequest(parseInt(this.ResendRequestCount), 0)
    },
    OpenResendRequest() {
      this.ResendRequestCountShow = true
      requestAnimationFrame(() => {
        this.$refs.ResendRequest.focus()
        nextTick(() => {
          const inputEl = this.$refs.ResendRequest.input
          if (inputEl) {
            inputEl.select()
          }
        })
      })
    },
    ResendRequest(n, BreakMode) {
      if (this.agSelectedRowNodes.length !== 1) {
        ElNotification({
          position: 'bottom-right',
          showClose: true,
          message: '重发请求失败\n\n没有选择请求',
          type: 'warning',
          customClass: 'multiline-message'
        })
        return
      }
      const Method = (this.agSelectedRowNodes[0].data["方式"]).toLowerCase()
      const isUDP = Method.indexOf("udp") !== -1
      if (isUDP) {
        ElNotification({
          position: 'bottom-right',
          showClose: true,
          message: '重发请求失败\n\nUDP请求不支持重发',
          type: 'warning',
          customClass: 'multiline-message'
        })
        return
      }
      if (Config_SunnyNetIsStart.value === false) {
        ElNotification({
          position: 'bottom-right',
          showClose: true,
          message: '重发请求失败\n\n您当前程序的工作端口未启动成功\n请修改端口后再试!!',
          type: 'warning',
          customClass: 'multiline-message'
        })
        return
      }
      const id = parseInt(this.agSelectedRowNodes[0].data['Theology'])
      let res = null;
      if (n === 1) {
        res = AppResendRequest(id, 1, BreakMode)
      } else {
        res = AppResendRequest(id, n, BreakMode)
      }
      res.then((ok) => {
        if (ok === false) {
          ElNotification({
            position: 'bottom-right',
            showClose: true,
            message: '重发请求失败\n\n未找到此请求',
            type: 'error',
            customClass: 'multiline-message'
          })
          return
        }
        ElNotification({
          position: 'bottom-right',
          showClose: true,
          message: '重发请求\n\n重新发送选中请求已提交',
          type: 'success',
          customClass: 'multiline-message'
        })
      })
    },
    GenerateCode(Language, Type) {
      if (this.agSelectedRowNodes.length < 1) {
        ElNotification({
          position: 'bottom-right',
          showClose: true,
          message: '代码生成失败\n\n没有选择请求',
          type: 'warning',
          customClass: 'multiline-message'
        })
        return
      }
      const id = parseInt(this.agSelectedRowNodes[0].data['Theology'])
      AppGenerateCode(id, Language, Type).then(res => {
        if (res === "") {
          ElNotification({
            position: 'bottom-right',
            showClose: true,
            message: '代码生成成功\n\n已复制到剪辑版',
            type: 'success',
            customClass: 'multiline-message'
          })
        } else {
          ElNotification({
            position: 'bottom-right',
            showClose: true,
            message: '代码生成失败\n\n' + res,
            type: 'warning',
            customClass: 'multiline-message'
          })
        }
      })
    },
    HomeListMenu(params) {
      const defaultMenuItems = params.defaultItems || [];
      const filteredMenuItems = defaultMenuItems.filter(item => {
        return item !== 'chartRange' && item !== 'export' && item !== 'paste' && item !== 'copyWithGroupHeaders'; // 禁用功能
      });
      if (this.agSelectedRowNodes.length > 0) {
        filteredMenuItems.push("separator");
        filteredMenuItems.push({
          name: "删除选中的请求",
          action: () => {
            const array = []
            const arrayID = []
            this.agSelectedRowNodes.forEach((node) => {
              arrayID.push(parseInt(node.data['Theology']))
              array.push(node.data)
              node.setSelected(false);
            });
            AppDeleteSession(arrayID).then(() => {
              this.agGridApi.applyTransaction({remove: array});
              this.agSelectedRowNodes = []
              Config_SelectedRow.value = [];
              this.agGridApi.clearCellSelection();
              this.agGridApi.refreshCells({
                columns: ['序号'],
                force: true,
                suppressFlash: true, // 不需要闪烁动画，可以不加
              });
            })
          },
        });
        filteredMenuItems.push({
          name: "删除非选中的请求",
          action: () => {
            const array = []
            const arrayID = []
            const arrMap = new Map();
            this.agSelectedRowNodes.forEach((node) => {
              arrMap.set(parseInt(node.data['Theology']), true);
            });
            this.agGridApi.forEachNode((node) => {
              const Theology = parseInt(node.data["Theology"]);
              if (!arrMap.has(Theology)) {
                array.push(node.data)
                arrayID.push(Theology)
              }
            });
            arrMap.clear();
            AppDeleteSession(arrayID).then(() => {
              this.agGridApi.applyTransaction({remove: array});
              this.agSelectedRowNodes = []
              Config_SelectedRow.value = [];
              this.agGridApi.clearCellSelection();
              this.agGridApi.refreshCells({
                columns: ['序号'],
                force: true,
                suppressFlash: true, // 不需要闪烁动画，可以不加
              });
            })
          },
        });
        filteredMenuItems.push("separator");
        const arrays = []
        Config_HomeTextMark.forEach((value) => {
          arrays.push({
            name: value["名称"],
            action: () => this.ColorMarking(value["id"])
          })
        })
        if (this.textMark.size > 0) {
          arrays.push("separator");
          arrays.push({name: "清除之前颜色标记", action: () => this.ColorMarking("Clear")})
        }
        filteredMenuItems.push({
          name: "颜色标记",
          subMenu: arrays,
        });
        if (this.agSelectedRowNodes.length === 1) {
          const Method = (this.agSelectedRowNodes[0].data["方式"]).toLowerCase()
          const isUDP = Method.indexOf("udp") !== -1
          if (!isUDP) {
            const isTCP = Method.indexOf("tcp") !== -1
            if (isTCP) {
              filteredMenuItems.push({
                name: "重发请求",
                subMenu: [
                  {
                    name: "普通重发",
                    action: () => {
                      this.ResendRequest(1, 0)
                    },
                  },
                  "separator"
                  ,
                  {
                    name: "批量重发",
                    action: this.OpenResendRequest,
                  },
                ],
              });
            } else {
              filteredMenuItems.push({
                name: "重发请求",
                subMenu: [
                  {
                    name: "普通重发",
                    action: () => {
                      this.ResendRequest(1, 0)
                    },
                  },
                  {
                    name: "重发 并 拦截上行",
                    action: () => {
                      this.ResendRequest(1, 1)
                    },
                  },
                  {
                    name: "重发 并 拦截下行",
                    action: () => {
                      this.ResendRequest(1, 2)
                    },
                  },
                  "separator"
                  ,
                  {
                    name: "批量重发",
                    action: this.OpenResendRequest,
                  },
                ],
              });
            }
            const ares = GetGenerateCodeListMenu(Method, this.GenerateCode, parseInt(this.agSelectedRowNodes[0].data['Theology']))
            if (ares.length > 0) {
              filteredMenuItems.push({
                name: "代码生成",
                subMenu: ares,
              });
            }
          }
        }
        filteredMenuItems.push("separator");
      }
      let i, t = '';
      if (Config_AutoRoll.value) {
        if (Config_IsDark.value) {
          i = '<div style="display: flex; align-items: center;width: 16px">' +
              `<svg xmlns="http://www.w3.org/2000/svg" style="top: 2px;position: relative;" width="16" height="14" viewBox="0 0 24 24" fill="none" stroke="#fff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">\\n' +
              '    <polyline points="20 6 9 17 4 12"/>' +
              '</svg>` +
              '</div>';
        } else {
          i = '<div style="display: flex; align-items: center;width: 16px">' +
              `<svg xmlns="http://www.w3.org/2000/svg" style="top: 2px;position: relative;" width="16" height="14" viewBox="0 0 24 24" fill="none" stroke="#000" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">\\n' +
              '    <polyline points="20 6 9 17 4 12"/>' +
              '</svg>` +
              '</div>';
        }

        t = "关闭自动滚动"
      } else {
        t = "开启自动滚动"
      }
      if (this.agSelectedRowNodes.length > 0) {
        let Methods = []
        for (let i = 0; i < this.agSelectedRowNodes.length; i++) {
          const _Method = (this.agSelectedRowNodes[i].data["方式"]).toLowerCase()
          const _State = (this.agSelectedRowNodes[i].data["状态"]).toLowerCase()
          if (_Method.indexOf("tcp") !== -1 && _State.indexOf("断开") === -1) {
            const id = parseInt(this.agSelectedRowNodes[i].data['Theology'])
            Methods.push(id)
          }
        }
        if (Methods.length > 0) {
          filteredMenuItems.push({name: "断开所有选中的TCP请求", action: () => AppDisconnectTCPRequest(Methods)})
          filteredMenuItems.push("separator");
        }
      }
      filteredMenuItems.push({
        name: t,
        action: () => {
          Config_AutoRoll.value = !Config_AutoRoll.value;
        },
        icon: i,
      });
      return filteredMenuItems;
    },
    openSy4File(selectedFiles) {
      ElNotification({
        position: 'bottom-right',
        showClose: true,
        message: '正在还原文件...\n这可能需要一定时间...',
        type: 'success',
        customClass: 'multiline-message'
      })
      Config_IsShowExportProgress.value = "请稍后,正在还原文件"
      AppImport(selectedFiles).then((obj) => {
        const err = obj[0]
        const res = obj[1]
        this.$refs.listen.insertArray(res, true, (arr) => {
          Config_IsShowExportProgress.value = ""
          if (err === "") {
            ElNotification({
              position: 'bottom-right',
              message: '记录文件还原成功！\n\n还原记录：' + arr.length + " 条",
              type: 'success',
              customClass: 'multiline-message'
            })
            return
          }
          ElNotification({
            position: 'bottom-right',
            message: '记录文件还原失败！\n\n' + err,
            type: 'warning',
            customClass: 'multiline-message'
          })
        })

      })
    },
    HeaderFileMenu() {
      const filteredMenuItems = [];
      filteredMenuItems.push({
        name: '打开记录文件',
        action: () => {
          OpenSunnyFile("").then((selectedFiles, err) => {
            this.openSy4File(selectedFiles)
          })
              .catch((error) => {
              });
        },
        icon: '<div style="display: flex; align-items: center;">' +
            '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-file-plus"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path><polyline points="14 2 14 8 20 8"></polyline><line x1="12" y1="18" x2="12" y2="12"></line><line x1="9" y1="15" x2="15" y2="15"></line></svg>' +
            '</div>',
      });
      const saveMenuItems = [{
        name: '保存全部记录',
        action: () => {
          const array = []
          this.agGridApi.forEachNodeAfterFilterAndSort((node) => {
            const theology = parseInt(node.data["Theology"]);
            array.push(theology)
          });
          if (array.length < 1) {
            ElNotification({
              position: 'bottom-right',
              showClose: true,
              message: '保存记录失败\n\n当前没有任何请求',
              type: 'error',
              customClass: 'multiline-message'
            })
            return
          }
          this.ExportMessage(array)
        }
      }];
      if (this.agSelectedRowNodes.length > 0) {
        saveMenuItems.push({
          name: '当前选择的所有记录',
          action: () => {
            const array = []
            for (let i = 0; i < this.agSelectedRowNodes.length; i++) {
              array.push(parseInt(this.agSelectedRowNodes[i].data["Theology"]))
            }
            if (array.length < 1) {
              ElNotification({
                position: 'bottom-right',
                showClose: true,
                message: '保存记录失败\n\n你没有选中任何请求',
                type: 'error',
                customClass: 'multiline-message'
              })
              return
            }
            this.ExportMessage(array)
          }
        })
      }
      filteredMenuItems.push({
        name: '保存',
        subMenu: saveMenuItems,
        icon: '<div style="display: flex; align-items: center;">' +
            '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-save"><path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"></path><polyline points="17 21 17 13 7 13 7 21"></polyline><polyline points="7 3 7 8 15 8"></polyline></svg>' +
            '</div>',
      });
      return filteredMenuItems;
    },
    MenuEvent(params) {
      if (Config_Menu_isFileMenu.value) {
        return this.HeaderFileMenu();
      }
      return this.HomeListMenu(params);
    },
    onColumnVisible(params) {
      const array = [];
      array.push({
        field: "全部数据",
        tooltipField: "全部数据",
        filter: true,
        headerName: "全部数据",
        filterParams: {
          filterOptions: ["contains", "notContains"]
        },
      })
      const allColumns = this.agGridApi.getColumns();
      allColumns.forEach((col) => {
        if (col.isVisible() && col.isFilterAllowed()) {
          array.push({
            field: col.getId(),
            tooltipField: col.getId(),
            headerName: col.getId(),
            filter: true,
          })
        }
      });
      this.FilterColumn = array
      this.saveColumn();
    },
    MoveColumnEvent() {
      if (this.MoveColumn) {
        return
      }
      this.MoveColumn = true
      this.agGridApi.setGridOption('suppressMovableColumns', false);
      this.ColumnMoveEventFunc()
    },
    ColumnMoveEventFunc() {
      requestAnimationFrame(() => {
        if (document.getElementsByClassName("ag-panel ag-default-panel ag-dialog ag-ltr ag-popup-child ag-focus-managed").length < 1) {
          this.agGridApi.setGridOption('suppressMovableColumns', true);
          this.MoveColumn = false
          return
        }
        this.ColumnMoveEventFunc()
      })
    }
  },
  watch: {
    ResendRequestCount(n) {
      if (parseInt(n) < 1 || isNaN(parseInt(n))) {
        this.ResendRequestCount = 1
        requestAnimationFrame(() => {
          this.$refs.ResendRequest.focus()
          nextTick(() => {
            const inputEl = this.$refs.ResendRequest.input
            if (inputEl) {
              inputEl.select()
            }
          })
        })
      }
    },
    IsDark() {
      requestAnimationFrame(() => {
        this.RefreshVisibleNodes()
      })
    },
  }
}
</script>

<template>
  <div :style="isNoDisableClick" id="ROOT">
    <div style="width: 100%;height:calc(100%);position: relative">
      <div style="width: 100%;height:30px;position: relative"
           class="home-root ag-theme-params-1 ag-theme-buttonStyle-1 ag-theme-iconSet-4">
        <Header></Header>
      </div>

      <div style="width: 100%;height: calc(100% - 60px)" id="rootList">
        <ag-grid-vue ref="agGridSide"
                     :theme="agTheme"
                     :sideBar="gridSideBar"
                     style="margin-left: -1px;height: 100% ;width: calc(100% + 2px);margin-top: -2px"
        />
        <ag-grid-vue @click="handleClick"
                     ref="agGrid"
                     :theme="agTheme"
                     :rowData="rowData"
                     :sideBar="sideBar"
                     :onGridReady="onGridReady"
                     style="margin-left: -1px;margin-top: -1px;margin-bottom: -1px;height: calc(100% + 3px) !important;width: calc(100% + 3px);"
                     :defaultColDef="defaultColDef"
                     :cellSelection="cellSelection"
                     :rowSelection="rowSelection"
                     :enableCharts="true"
                     :grid-options="gridOptions"
                     :loading="false"
                     :allowContextMenuWithControlKey="true"
                     :overlayNoRowsTemplate="overlayNoRowsTemplate"
                     :suppressCutToClipboard="true"
        />
      </div>
      <div style="width: 100%;height:32px;font-size: 18px" class="home-root ag-theme-params-1">
        <HomeFooter style="width: 100%;top: 0" :isHideHook="HideHook"></HomeFooter>
        <el-dialog
            v-model="ResendRequestCountShow"
            title="批量重放请求"
            width="360"
            align-center
        >
          <el-form label-position="top">
            <el-form-item label="您想要重发多少次？">
              <el-input-number
                  ref="ResendRequest"
                  v-model="ResendRequestCount"
                  :min="1"
                  :max="100000"
                  controls-position="right"
                  style="width: 100%;"
                  @keydown.enter="SubResendRequest"
              />
            </el-form-item>
          </el-form>

          <template #footer>
            <div class="dialog-footer">
              <el-button @click="SubResendRequest">确定重发</el-button>
            </div>
          </template>
        </el-dialog>

      </div>
      <ListenOn ref="listen" style="display: none" :addFilter="addFilter" :setSelectedRow="setSelectedRow"/>
      <FindWindow style="display: none"/>
      <Filter ref="Filter" Name="主列表" :Parent="mainFilter" :column="FilterColumn" :apply="onFilterApply"
              :SearchDone="SearchDone" :CancelSearch="CancelSearch"/>


    </div>
  </div>
</template>