<template>
  <div
      ref="container"
      class="monaco-editor"
      style="text-align: left;height: 100%;"
      :drak="getTheme"
  ></div>
</template>
<script>
import * as monaco from 'monaco-editor'
import {Config_IsDark, Config_MonacoEditorTheme} from "../../../config/Config";
import {addMonacoEditorFontSize, RemoveMonacoEditorMenu} from "../../../config/editorContextMenu";


function IsRemoveMenu(value, zd) {

  let Menu = [
    "转到符号", "更改所有匹配项", "命令面板", "格式化"
  ]
  for (let i = 0; i < Menu.length; i++) {
    if (value.indexOf(Menu[i]) !== -1) {
      return true
    }
  }
  return false
}

export default {
  props: ["Overflow"],
  data() {
    return {
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
        fixedOverflowWidgets: this.getOverflow(),
        diffViewport: false,
        wordWrap: 'on', // 设置自动换行
        validationOptions: {
          validate: false, // 禁用语法错误提示
        },
        minimap: {enabled: false},
        fontSize: 13,
        readOnly: false
      },
      Function: {
        Save: null,
        GetCode: null,
        setValue: null,
        formatCode: null,
        SetReadOnly: null,
        setModelLanguage: null,
        DelAllDecorations: null,
        deltaDecorations: null,
        codiconClose: null,
        SetTheme: null,
        SetNewValue: null
      },
      IsHasModify: false,
      inn: 0
    }
  },
  mounted() {
    this.init()
  },
  methods: {
    getOverflow() {
      return this.Overflow === "true" || this.Overflow === true
    },
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
      this.Function.setValue = (newContent) => {
        editor.setValue(newContent);
        editor.setScrollTop(0);
      }
      this.Function.formatCode = () => {
        editor?.getAction('editor.action.formatDocument').run()
      }


      this.Function.DelAllDecorations = () => {
        if (editor) {
          let ids = [];
          let decorations = editor.getModel().getAllDecorations();
          for (let decoration of decorations) {
            if (decoration !== undefined && decoration !== null) {
              if (decoration.options !== undefined && decoration.options !== null) {
                if (decoration.options.inlineClassName !== undefined && decoration.options.inlineClassName !== null) {
                  if (decoration.options.inlineClassName.indexOf("vscode-editor") !== -1) {
                    ids.push(decoration.id);
                  }
                }
              }
            }
          }
          if (ids && ids.length) {
            editor.deltaDecorations(ids, []);
          }
        }
      }
      this.Function.SetTheme = (newTheme) => {
        monaco.editor.setTheme(newTheme);
      }
      this.Function.codiconClose = () => {
        setTimeout(() => {
          this.inn++
          try {
            document.getElementsByClassName("codicon-widget-close")[0].click()
          } catch (e) {
          }
          if (this.inn < 30) {
            this.Function.codiconClose()
          }
        }, 100)
      }
      this.Function.deltaDecorations = (position, color, length) => {
        const dark = Config_IsDark.value ? "-dark" : ""
        const decorations = editor.deltaDecorations([], [
          {
            range: new monaco.Range(position.lineNumber, position.column, position.lineNumber, position.column + length),
            options: {inlineClassName: `vscode-editor-${color}` + dark, hoverMessage: 'Custom colored text'},
          },
        ]);

        // 更新装饰器
        editor.deltaDecorations(decorations, [
          {
            range: new monaco.Range(position.lineNumber, position.column, position.lineNumber, position.column + length),
            options: {inlineClassName: `vscode-editor-${color}` + dark, hoverMessage: 'Custom colored text'},
          },
        ]);
      }
      this.Function.SetNewValue = (newValue) => {
        this.IsHasModify = true
        this.inn = 0
        this.Function.DelAllDecorations()
        let array = newValue.replaceAll("\r", "").replaceAll(" ", " ").split("\n\n")
        array = array[0].split("\n")
        let language = "plaintext"
        for (let i = 0; i < array.length; i++) {
          const languageHeader = array[i].toUpperCase().replaceAll(";", "/") + "/"
          if (languageHeader.indexOf('CONTENT-TYPE') !== -1) {
            const ar = languageHeader.split("/")
            if (ar.length >= 2) {
              language = ar[1].toLowerCase()
            }
          }
          let ar1 = array[i].split(":")
          if (ar1.length >= 1) {
            this.Function.deltaDecorations({lineNumber: i + 1, column: 1}, "HeaderName", ar1[0].length)
            this.Function.deltaDecorations({
              lineNumber: i + 1,
              column: ar1[0].length + 1
            }, "HeaderValue", array[i].length - ar1[0].length)
          }
        }
        if (editor) {
          this.$nextTick(() => {
            monaco.editor.setModelLanguage(editor.getModel(), language);
          })
        }
        this.$nextTick(() => {
          this.Function.codiconClose()
        })
      }

      // 监听内容变化事件
      editor.getModel().onDidChangeContent(() => {
        this.Function.SetNewValue(editor.getValue())
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
    theme(n) {
      if (this.Function.SetTheme) {
        this.Function.SetTheme(n)
      }
    }
  },
  computed: {
    getTheme() {
      this.$nextTick(() => {
        this.$nextTick(() => {
          if (this.Function.GetCode) {
            if (this.Function.SetNewValue) {
              const code = this.Function.GetCode()
              const IsHasModify = this.IsHasModify
              this.Function.SetNewValue(code)
              this.IsHasModify = IsHasModify
            }
          }
        })
      })
      return this.theme
    }
  }
}
</script>


<style>
.vscode-editor-Method {
  color: #A626A4 !important;
}

.vscode-editor-Method-dark {
  color: #D52665 !important;
}

.vscode-editor-URL {
  color: #50A14F !important;
}

.vscode-editor-URL-dark {
  color: #88AD6E !important;
}

.vscode-editor-SS {
  color: #4078F2 !important;
}

.vscode-editor-SS-dark {
  color: #61AEEE !important;
}

.vscode-editor-SS2 {
  color: #A57C22 !important;
}

.vscode-editor-SS2-dark {
  color: #AF835A !important;
}

.vscode-editor-HeaderName {
  color: #279425 !important;
}

.vscode-editor-HeaderName-dark {
  color: #88AD6E !important;
}

.vscode-editor-HeaderValue {
  color: #000000 !important;
}

.vscode-editor-HeaderValue-dark {
  color: #dcdcdc !important;
}
</style>