<script>
import Table from "../../Tools/table.vue";
import {Config_IsRest, Config_MonacoEditorTheme} from "../../config/Config";
import * as monaco from "monaco-editor";
import {ElNotification} from "element-plus";
import {GetMustTcpRoles, GetMustTcpType, SetMustTcpRoles} from "../../../../bindings/changeme/Service/appmain";
import {addMonacoEditorFontSize, RemoveMonacoEditorMenu} from "../../config/editorContextMenu";
import {attachMcpSettingsReload} from "../../config/mcpSettingsSync.js";

export default {
  components: {Table},
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
        setValue: null,
        SetTheme: null,
        formatCode: null,
      },
      model: "MustTcpLei",
      isInit: false,
    }
  },
  computed: {
    getTheme() {
      return this.theme
    }
  },
  methods: {
    IsRemoveMenu(value, zd) {
      let Menu = [
        "转到符号", "更改所有匹配项", "命令面板",
      ]
      if (zd) {
        Menu.push("格式化")
      }
      for (let i = 0; i < Menu.length; i++) {
        if (value.indexOf(Menu[i]) !== -1) {
          return true
        }
      }
      return false
    }
    ,
    init() {
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
        run: () => {
          if (this.Save != null) {
            this.Save()
          }
        },
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
      RemoveMonacoEditorMenu(editor, this.IsRemoveMenu, ()=>{return this.defaultOpts.readOnly;});

      this.Function.GetCode = () => {
        return editor.getValue()
      }
      this.Function.setValue = (newContent) => {
        editor.setValue(newContent);
        editor.setScrollTop(0);
      }
      this.Function.SetTheme = (newTheme) => {
        monaco.editor.setTheme(newTheme);
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
              message: '您已修改强制TCP规则,不要忘了右键选择保存！',
              type: 'warning',
            })
          }
        }
        this.IsHasModify = true
      });
    }
    ,
    Save() {
      this.IsHasModify = false
      let __type = 1;
      if (this.model === "ALLMustTcp") {
        __type = 0
      } else if (this.model === "MustTcpLei") {
        __type = 1
      } else if (this.model === "MustTcpWai") {
        __type = 2
      }
      SetMustTcpRoles(__type, this.Function.GetCode()).then(() => {
        ElNotification({
          position: 'bottom-right',
          showClose: true,
          message: '强制TCP规则已更新',
          type: 'success',
        })
      })
    },
    reloadMustTcpFromBackend() {
      this.isInit = false;
      return GetMustTcpRoles().then((res) => {
        this.Function.setValue(res);
        return GetMustTcpType().then((ty) => {
          if (ty === 0) {
            this.model = "ALLMustTcp";
          } else if (ty === 1) {
            this.model = "MustTcpLei";
          } else if (ty === 2) {
            this.model = "MustTcpWai";
          }
          this.isInit = true;
          this.IsHasModify = false;
        });
      });
    },
  },
  watch: {
    model() {
      this.Save()
    },
    theme(n) {
      if (this.Function.SetTheme) {
        this.Function.SetTheme(n)
      }
    },
    isRest() {
      this.isInit = false
      GetMustTcpRoles().then(res => {
        this.Function.setValue(res)
        GetMustTcpType().then(ty => {
          if (ty === 0) {
            this.model = "ALLMustTcp"
          } else if (ty === 1) {
            this.model = "MustTcpLei"
          } else if (ty === 2) {
            this.model = "MustTcpWai"
          }
          this.isInit = true
          this.IsHasModify = false
        })
      })
    }
  },
  mounted() {
    this.init()
    attachMcpSettingsReload("musttcp", () => this.reloadMustTcpFromBackend());
    GetMustTcpRoles().then(res => {
      this.Function.setValue(res)
      GetMustTcpType().then(ty => {
        if (ty === 0) {
          this.model = "ALLMustTcp"
        } else if (ty === 1) {
          this.model = "MustTcpLei"
        } else if (ty === 2) {
          this.model = "MustTcpWai"
        }
        this.isInit = true
        this.IsHasModify = false
      })
    })
  }
}
</script>


<template>
  <div style="display: flex; flex-direction: column; gap: 5px; margin: 5px;">
    <div style="display: flex; position: relative; gap:30px; justify-content: center;">
      <el-radio-group v-model="model" style="display: flex; position: relative; gap:30px; justify-content: center;">
        <el-tooltip placement="bottom">
          <template #content>
            <div style="white-space: normal; line-height: 1.4;">
              开启后所有请求都转为TCP请求,HTTPS将不会解密数据
            </div>
          </template>
          <div
              class="ag-labeled ag-label-align-left ag-toggle-button ag-input-field ag-group-title-bar ag-charts-format-sub-level-group-title-bar ag-unselectable ag-selected">
            <el-radio value="ALLMustTcp">全部走TCP</el-radio>
          </div>
        </el-tooltip>
        <el-tooltip placement="bottom">
          <template #content>
            <div style="white-space: normal; line-height: 1.4;">
              不在规则中的地址,将强制转为TCP请求
            </div>
          </template>
          <div
              class="ag-labeled ag-label-align-left ag-toggle-button ag-input-field ag-group-title-bar ag-charts-format-sub-level-group-title-bar ag-unselectable ag-selected">
            <el-radio value="MustTcpWai">规则外走TCP</el-radio>
          </div>
        </el-tooltip>
        <el-tooltip placement="bottom">
          <template #content>
            <div style="white-space: normal; line-height: 1.4;">
              规则内的地址,将强制转为TCP请求
            </div>
          </template>
          <div
              class="ag-labeled ag-label-align-left ag-toggle-button ag-input-field ag-group-title-bar ag-charts-format-sub-level-group-title-bar ag-unselectable ag-selected">
            <el-radio value="MustTcpLei">规则内走TCP</el-radio>
          </div>
        </el-tooltip>
      </el-radio-group>
    </div>
    <div style="display: flex; position: relative; gap:30px; justify-content: center;" v-show="model!=='ALLMustTcp'">
      <div
          ref="container"
          class="monaco-editor"
          style="text-align: left;height: 400px;width: 100%;"
          :drak="getTheme"
      ></div>
    </div>
  </div>
</template>
