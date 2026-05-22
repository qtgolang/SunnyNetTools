<script>
import Table from "../../Tools/table.vue";
import {Config_IsRest, Config_MonacoEditorTheme, ObjString} from "../../config/Config";
import {ElMessage, ElNotification} from "element-plus";
import Dns from "./Dns.vue";
import * as monaco from "monaco-editor";
import {
  CreateProxyWay,
  GetProxyRoles,
  ProxyWayList,
  ProxyWayUpdate,
  SetProxyRoles
} from "../../../../bindings/changeme/Service/appmain";
import {addMonacoEditorFontSize, RemoveMonacoEditorMenu} from "../../config/editorContextMenu";
import {attachMcpSettingsReload} from "../../config/mcpSettingsSync.js";

export default {
  components: {Dns, Table},
  data() {
    return {
      set isRest(val) {
        Config_IsRest.value = val;
      },
      get isRest() {
        return Config_IsRest.value;
      },
      get theme() {
        return Config_MonacoEditorTheme.value
      },
      set theme(newValue) {
        Config_MonacoEditorTheme.value = newValue
      },
      // 主要配置
      defaultOpts: {
        value: '', // 编辑器的值
        theme: this.theme, // 编辑器主题：vs, hc-black, or vs-dark，更多选择详见官网
        roundedSelection: true, // 右侧不显示编辑器预览框
        scrollBeyondLastLine: false, // 禁止滚动超过最后一行
        autoIndent: true, // 自动缩进
        automaticLayout: true,
        formatOnType: true,
        formatOnPaste: true,
        originalEditable: true,
        glyphMargin: false,
        diffViewport: false,
        wordWrap: 'on', // 设置自动换行
        validationOptions: {
          validate: false, // 禁用语法错误提示
        },
        minimap: {enabled: false},
        fontSize: 13,
        readOnly: false,
        links: false,
        hover: {enabled: false},
        quickSuggestions: false,
        parameterHints: {enabled: false},
        suggestOnTriggerCharacters: false,
        wordBasedSuggestions: false,
        inlineSuggest: {enabled: false},
        colorDecorators: false,
        autoClosingBrackets: false,
        autoClosingQuotes: false,
        overviewRulerLanes: 0,
        lightbulb: {enabled: false},
      },
      Function: {
        GetCode: null,
        SetTheme: null,
        setValue: null,
        formatCode: null,
      },
      gid: 0,
      IsHasModify: false,
      isInit: false
    }
  },
  methods: {
    Empty() {
      this.agGridApi.setGridOption("rowData", []);
    },
    init() {
      //初始化表格信息
      {
        this.$refs.Way.agGridApi.hideOverlay();
        this.$refs.Way.agGridApi.setGridOption('overlayNoRowsTemplate', `<span style="padding: 20px;" id="HookMessageText">您还没有添加任何上游代理</span>`);
        this.$refs.Way.agGridApi.showNoRowsOverlay();
        this.$refs.Way.agGridApi.setGridOption('columnDefs', [
          {
            field: "状态", tooltipField: '状态',
            minWidth: 60,
            maxWidth: 60,
            width: 60,
          },
          {
            field: "上游代理", tooltipField: '上游代理',
            minWidth: 170,
            width: 170,
            maxWidth: 500,
            editable: true,
          },
          {
            field: "注释", tooltipField: '注释',
            minWidth: 150,
            maxWidth: 570,
            width: 500,
            editable: true,
          },
        ]);
        this.$refs.Way.Stopped = this.WayEdit
        this.$refs.Way.addValue = this.addWay
        this.$refs.Way.CellForClipboard = this.onCellForClipboard
      }
    },
    IsRemoveMenu(value) {
      let Menu = [
        "转到符号", "更改所有匹配项", "命令面板", "格式化"
      ]
      for (let i = 0; i < Menu.length; i++) {
        if (value.indexOf(Menu[i]) !== -1) {
          return true
        }
      }
      return false
    },
    reloadProxyWayFromBackend() {
      if (!this.$refs.Way?.agGridApi) {
        return Promise.resolve();
      }
      this.$refs.Way.Empty();
      return ProxyWayList().then((list) => {
        const array = [];
        for (let i = 0; i < list.length; i++) {
          const item = list[i];
          array.push({
            "上游代理": item["URL"],
            "注释": item["Note"],
            "状态": item["State"],
            "mid": item["ID"] + "",
            id: (this.gid++) + "",
          });
        }
        this.$refs.Way.agGridApi.applyTransaction({add: array});
      });
    },
    reloadProxyRolesFromBackend() {
      if (!this.Function?.setValue) {
        return Promise.resolve();
      }
      return GetProxyRoles().then((res) => this.Function.setValue(res));
    },
    SaveRoles() {
      this.IsHasModify = false;
      SetProxyRoles(this.Function.GetCode()).then(() => {
        ElNotification({
          showClose: true,
          message: '上游代理使用规则:已保存',
          type: 'success',
          position: 'bottom-right',
        })
      })
    },
    addWay() {
      CreateProxyWay().then(id => {
        this.$refs.Way.agGridApi.applyTransaction({
          add: [{
            "上游代理": "双击修改",
            "注释": '无',
            "状态": '禁用',
            "mid": id,
            id: (this.gid++) + "",
          }]
        });
        ElMessage.success('已添加,请双击修改,双击“状态”切换状态')
      })
    },
    onCellForClipboard(params) {
      const nodeData = params.node.data;
      if (nodeData["状态"] === "禁用") {
        nodeData["状态"] = "启用"
      } else {
        nodeData["状态"] = "禁用"
      }
      const array = [];
      array.push(nodeData)
      this.$refs.Way.agGridApi.forEachNode(node => {
        if (node.data.mid !== nodeData.mid) {
          node.data["状态"] = "禁用"
          array.push(node.data)
        }
      });
      this.$refs.Way.agGridApi.applyTransaction({update: array});
      this.WayEdit(params)
    },
    WayEdit(params) {
      const array = [];
      const isUserExist = (user) => {
        return array.some(item => item.Proxy === user);
      }
      let isExist = false;
      this.$refs.Way.agGridApi.forEachNode(node => {
        const Proxy = ObjString(node.data["上游代理"]).trim();
        const Status = ObjString(node.data["状态"]).trim();
        const Note = ObjString(node.data["状态"]).trim();
        const mid = node.data["mid"];
        if (!isUserExist(Proxy)) {
          array.push({Proxy: Proxy, Status: Status, Note: Note, mid: mid})
          return
        }
        isExist = true
      });
      if (isExist) {
        this.$refs.Way.Empty()
        const lodArray = [];
        array.forEach(item => {
          lodArray.push({
            "上游代理": item.Proxy,
            "状态": item.Status,
            "注释": item.Note,
            "mid": item.mid,
            id: (this.gid++) + "",
          });
        })
        this.$refs.Way.agGridApi.applyTransaction({add: lodArray});
        ElMessage.error('您添加的上游代理已存在,请重新添加！')
        return
      }
      const node = params.node.data;
      ProxyWayUpdate(parseInt(node.mid), node["上游代理"], node["状态"], node["注释"]).then((res) => {
        if (res) {
          ElNotification({
            showClose: true,
            message: '上游代理已更新',
            type: 'success',
            position: 'bottom-right',
          })
        } else {
          ElNotification({
            position: 'bottom-right',
            showClose: true,
            message: '上游代理更新失败',
            type: 'warning',
          })
        }
      })
    },
    initEdit() {
      // 初始化container的内容，销毁之前生成的编辑器
      this.$refs.container.innerHTML = ''
      // 生成 diff-editor 对象
      const editor = monaco.editor.create(this.$refs.container, this.defaultOpts)

      const model = editor.getModel()
      monaco.editor.setModelLanguage(model, "json")
      addMonacoEditorFontSize(editor)
      monaco.languages.typescript.javascriptDefaults.setDiagnosticsOptions({
        noSemanticValidation: true,  // 关闭 JS 语义检查
        noSyntaxValidation: true,    // 关闭 JS 语法检查
      });
      monaco.languages.typescript.typescriptDefaults.setDiagnosticsOptions({
        noSemanticValidation: true,
        noSyntaxValidation: true,
      });
      monaco.languages.json.jsonDefaults.setDiagnosticsOptions({
        validate: false, // 关闭 JSON 语法检查
      });

      editor.addAction({
        id: 'SaveConfiguration',
        label: '保存/应用修改',
        keybindings: [monaco.KeyMod.chord(monaco.KeyMod.CtrlCmd | monaco.KeyCode.KeyS)],
        contextMenuGroupId: 'my-commands2',
        contextMenuOrder: Number.MAX_SAFE_INTEGER,
        run: this.SaveRoles,
      })
      const key = editor.createContextKey('wordWrapOn', true)
      // 默认换行的动作
      editor.addAction({
        id: 'turnWordWrapOff',
        label: '关闭自动换行',
        contextMenuGroupId: 'my-commands',
        contextMenuOrder: Number.MAX_SAFE_INTEGER,
        precondition: 'wordWrapOn',
        run: () => {
          this.defaultOpts.wordWrap = 'off'
          editor.updateOptions({
            wordWrap: "off"
          });
          key.set(false)
        },
      })
      editor.addAction({
        id: 'turnWordWrapOn',
        label: '自动换行',
        contextMenuGroupId: 'my-commands',
        contextMenuOrder: Number.MAX_SAFE_INTEGER,
        precondition: '!wordWrapOn',
        run: () => {
          this.defaultOpts.wordWrap = 'on'
          editor.updateOptions({
            wordWrap: "on"
          });
          key.set(true)
        },
      })

      editor.addCommand(monaco.KeyCode.F1, function () {
        //注册一个F1快捷键，但是什么也不做，用来阻止弹出命令列表
      });
      RemoveMonacoEditorMenu(editor, this.IsRemoveMenu,()=>{return this.defaultOpts.readOnly;});
      this.Function.GetCode = () => {
        return editor.getValue()
      }

      this.Function.SetTheme = (newTheme) => {
        monaco.editor.setTheme(newTheme);
      }
      this.Function.setValue = (newContent) => {
        editor.setValue(newContent);
        editor.setScrollTop(0);
      }
      this.Function.formatCode = () => {
        editor?.getAction('editor.action.formatDocument').run()
      }
      // 监听内容变化事件
      editor.getModel().onDidChangeContent(() => {
        if (!this.IsHasModify) {
          if (this.isInit) {
            ElNotification({
              position: 'bottom-right',
              showClose: true,
              message: '您已修改上游代理使用规则,不要忘了右键选择保存！',
              type: 'warning',
            })
          }
        }
        this.IsHasModify = true
      });
    }
  },
  computed: {
    getTheme() {
      return this.theme
    }
  },
  watch: {
    isRest() {
      this.isInit = false
      this.$refs.Way.Empty();
      GetProxyRoles().then(res => {
        this.Function.setValue(res)
        ProxyWayList().then(list => {
          const array = []
          for (let i = 0; i < list.length; i++) {
            const item = list[i];
            array.push({
              "上游代理": item["URL"],
              "注释": item["Note"],
              "状态": item["State"],
              "mid": item["ID"] + "",
              id: (this.gid++) + "",
            })
          }
          this.$refs.Way.agGridApi.applyTransaction({add: array});
          this.isInit = true;
        })
      })
    },
    theme(n) {
      if (this.Function.SetTheme) {
        this.Function.SetTheme(n)
      }
    }
  },
  mounted() {
    this.init()
    this.initEdit()
    attachMcpSettingsReload("proxy_way", () => this.reloadProxyWayFromBackend());
    attachMcpSettingsReload("proxy_roles", () => this.reloadProxyRolesFromBackend());
    GetProxyRoles().then(res => {
      this.Function.setValue(res)
      ProxyWayList().then(list => {
        const array = []
        for (let i = 0; i < list.length; i++) {
          const item = list[i];
          array.push({
            "上游代理": item["URL"],
            "注释": item["Note"],
            "状态": item["State"],
            "mid": item["ID"] + "",
            id: (this.gid++) + "",
          })
        }
        this.$refs.Way.agGridApi.applyTransaction({add: array});
        this.isInit = true;
      })
    })
  }
}
</script>

<template>
  <div style="width: 100%">
    <Dns/>
    <Table ref="Way" style="height: 200px"/>
    <span style="justify-content: center;text-align: center;display: flex;margin-top: 10px;margin-bottom: 10px">以下地址不使用上游代理</span>
    <div
        ref="container"
        class="monaco-editor"
        style="text-align: left;height: 100px;width: 100%;"
        :drak="getTheme"
    ></div>
  </div>
</template>
