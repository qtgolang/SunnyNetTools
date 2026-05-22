<template>
  <div
      ref="container"
      class="monaco-editor"
      style="text-align: left;height: 100%;position: relative"
      v-show="getTheme"
  ></div>
</template>
<script>
import * as monaco from 'monaco-editor'
import {Config_IsDark, Config_MonacoEditorTheme} from "../../../config/Config";
import {addMonacoEditorFontSize, RemoveMonacoEditorMenu} from "../../../config/editorContextMenu";

let SaveTextUpdate = true;

function IsRemoveMenu(value, zd) {
  if (value.indexOf("保存修改") !== -1) {
    return !SaveTextUpdate
  }
  let Menu = [
    "转到符号", "更改所有匹配项", "命令面板"
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

export default {
  data() {
    return {
      get theme() {
        return window.Theme.IsDark
      },
      set theme(newValue) {
        window.Theme.IsDark = newValue
      },
      get MonacoEditorTheme() {
        return Config_MonacoEditorTheme.value
      },
      set MonacoEditorTheme(s) {
        Config_MonacoEditorTheme.value = s
      },
      // 主要配置
      defaultOpts: {
        value: '', // 编辑器的值
        theme: this.MonacoEditorTheme, // 编辑器主题：vs, hc-black, or vs-dark，更多选择详见官网
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
        readOnly: true
      },
      Function: {
        Save: null,
        GetCode: null,
        setValue: null,
        formatCode: null,
        SetReadOnly: null,
        setModelLanguage: null,
        SetTheme: null,
        setSaveTextUpdate: null
      },
      IsHasModify: false
    }
  },
  mounted() {
    this.init()
  },
  methods: {
    init() {
      // 初始化container的内容，销毁之前生成的编辑器
      this.$refs.container.innerHTML = ''
      // 生成 diff-editor 对象
      const editor = monaco.editor.create(this.$refs.container, this.defaultOpts)

      const model = editor.getModel()
      monaco.editor.setModelLanguage(model, "plaintext")

      addMonacoEditorFontSize(editor)
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
      RemoveMonacoEditorMenu(editor, IsRemoveMenu,()=>{return this.defaultOpts.readOnly;});
      this.Function.GetCode = () => {
        return editor.getValue()
      }
      this.Function.SetReadOnly = (readOnlyState) => {
        editor.updateOptions({
          readOnly: readOnlyState
        });
        this.defaultOpts.readOnly=readOnlyState;
      }
      this.Function.setModelLanguage = (lang) => {
        this.$nextTick(() => {
          monaco.editor.setModelLanguage(editor.getModel(), lang);
        })
      }
      this.Function.setSaveTextUpdate = (o) => {
        this.$nextTick(() => {
          SaveTextUpdate = o;
        })
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
        try {
          JSON.parse(editor.getValue());
          this.$nextTick(() => {
            monaco.editor.setModelLanguage(editor.getModel(), "json");
          })
        } catch (error) {

        }

        this.IsHasModify = true
      });
    },
    // 供父组件调用手动获取值
    GetCode() {
      if (this.Function.GetCode) {
        return this.Function.GetCode()
      }
      return null
    },
    SetReadOnly(v) {
      if (this.Function.SetReadOnly) {
        this.Function.SetReadOnly(v)
      }
    },
    SetLanguage(v) {
      if (this.Function.setModelLanguage) {
        this.Function.setModelLanguage(v)
      }
    },
    SetSaveTextUpdate(v) {
      if (this.Function.setSaveTextUpdate) {
        this.Function.setSaveTextUpdate(v)
      }
    },
    SetSaveFunc(v) {
      this.Function.Save = v
    },
    SetCode(newContent) {
      if (this.Function.setValue) {
        this.Function.setValue(newContent)
      }
      this.IsHasModify = false
    },
    formatCode() {
      if (this.Function.formatCode) {
        this.Function.formatCode()
      }
    },
    HasModify() {
      return this.IsHasModify
    }
  },
  watch: {
    MonacoEditorTheme(n) {
      if (this.Function.SetTheme) {
        this.Function.SetTheme(n)
      }
    }
  },
  computed: {
    getTheme() {
      if (this.Function.SetTheme) {
        this.Function.SetTheme(this.MonacoEditorTheme)
      }
      return true;
    }
  }
}
</script>

