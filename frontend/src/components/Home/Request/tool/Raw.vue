<template>
  <div
      ref="container"
      class="monaco-editor"
      style="text-align: left;height: 100%;"
      v-show="getTheme"
  ></div>
</template>
<script>
import * as monaco from 'monaco-editor'
import {Config_IsDark, Config_MonacoEditorTheme} from "../../../config/Config.js";
import {
  addMonacoEditorEncodingMenu,
  bytesToBase64,
  Config_Encoding_Current_ISUTF8,
  StringToBytes,
  toGBK,
  toUTF8
} from "../../../config/encoding.js";
import {
  addMonacoEditorFontSize,
  addMonacoEditorWordWrapMenu,
  RemoveMonacoEditorMenu
} from "../../../config/editorContextMenu.js";
import {bytesToRequest, bytesToResponse} from "../../../config/SunnyNetInfoApi.js";

export default {
  props: ["Overflow", "NoShowEncoding"],
  data() {
    return {
      get theme() {
        return Config_IsDark.value
      },
      set theme(newValue) {
        if (Config_IsDark.value !== newValue) Config_IsDark.value = newValue
      },
      get ShowUTF8Encoding() {
        return Config_Encoding_Current_ISUTF8.value
      },
      set ShowUTF8Encoding(newValue) {
        Config_Encoding_Current_ISUTF8.value = newValue
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
        fixedOverflowWidgets: this.getOverflow(),
        glyphMargin: false,
        diffViewport: false,
        wordWrap: 'on', // 设置自动换行
        validationOptions: {
          validate: false, // 禁用语法错误提示
        },
        minimap: {enabled: false},
        fontSize: 13,
        readOnly: true,
        presentData: null,
        links: false, // 禁用 URL 链接提示
        hover: {enabled: false}, // 禁用悬停提示
        quickSuggestions: false,  // 禁用输入时的代码建议
        parameterHints: {enabled: false}, // 禁用参数提示
        suggestOnTriggerCharacters: false, // 禁用触发字符（如 `.` 触发补全）
        wordBasedSuggestions: false, // 禁用基于单词的建议
        inlineSuggest: {enabled: false}, // 禁用内联建议
        colorDecorators: false, //关闭颜色选择器
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
        SetTheme: null,
        SetNewValue: null,
        __switchEncoding: null,
      },
      IsHasModify: false,
      IsSetCode: false,
      value: null,
      inn: 0
    }
  },
  mounted() {
    this.init()
  },
  watch: {
    ShowUTF8Encoding(n, l) {
      if (this.Function.__switchEncoding) {
        this.Function.__switchEncoding(n);
      }
    },
    MonacoEditorTheme(n) {
      if (this.Function.SetTheme) {
        this.Function.SetTheme(n)
      }
    }
  },
  methods: {
    getOverflow() {
      return this.Overflow === "true" || this.Overflow === true
    },
    IsRemoveMenu(value, zd) {
      let Menu = [
        "转到符号", "更改所有匹配项", "命令面板"
      ]
      if (this.NoShowEncoding === true) {
        Menu.push("编码")
      }

      if (zd) {
        Menu.push("格式化")
      }
      if (this.defaultOpts.presentData === null) {
        Menu.push("展示完整数据")
      }
      for (let i = 0; i < Menu.length; i++) {
        if (value.indexOf(Menu[i]) !== -1) {
          return true
        }
      }
      return false
    },
    init() {
      // 初始化container的内容，销毁之前生成的编辑器
      this.$refs.container.innerHTML = ''
      // 生成 diff-editor 对象
      const editor = monaco.editor.create(this.$refs.container, this.defaultOpts)

      const model = editor.getModel()
      monaco.editor.setModelLanguage(model, "plaintext")
      addMonacoEditorFontSize(editor)
      addMonacoEditorEncodingMenu(editor);
      addMonacoEditorWordWrapMenu(editor, this.defaultOpts);
      editor.addCommand(monaco.KeyCode.F1, function () {
        //注册一个F1快捷键，但是什么也不做，用来阻止弹出命令列表
      });

      RemoveMonacoEditorMenu(editor, this.IsRemoveMenu,()=>{return this.defaultOpts.readOnly;});
      this.Function.GetCode = () => {
        return editor.getValue()
      }
      this.Function.SetReadOnly = (readOnlyState) => {
        editor.updateOptions({
          readOnly: readOnlyState,
        });
        this.defaultOpts.readOnly=readOnlyState;
      }
      this.Function.__switchEncoding = (newEncoding) => {
        this.$nextTick(() => {
          if (this.Function.setValue) {
            const position = editor.getPosition(); // 获取光标位置
            const scrollPosition = editor.getScrollTop(); // 获取滚动条位置
            if (newEncoding) {
              this.Function.setValue(toUTF8(this.value))
            } else {
              this.Function.setValue(toGBK(this.value))
            }
            this.$nextTick(() => {
              editor.setPosition(position); // 设置光标位置
              editor.setScrollTop(scrollPosition); // 设置滚动条位置
            })
          }
        })
      }
      this.Function.setModelLanguage = (lang) => {
        this.$nextTick(() => {
          monaco.editor.setModelLanguage(editor.getModel(), lang);
        })
      }
      this.Function.setValue = (newContent) => {
        this.$nextTick(() => {
          if (newContent.length > 1024000) {
            editor.trigger('keyboard', 'turnWordWrapOff', {});
          } else {
            editor.trigger('keyboard', 'turnWordWrapOn', {});
          }
          editor.setValue(newContent);
          editor.setScrollTop(0);
        })
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
          if (i === 0) {
            let ar1 = array[i].split(" ")
            if (ar1.length === 3) {
              if (ar1[1].toUpperCase().startsWith('HTTP')) {
                this.Function.deltaDecorations({lineNumber: 1, column: 1}, "Method", ar1[0].length)
                this.Function.deltaDecorations({lineNumber: 1, column: ar1[0].length + 2}, "URL", ar1[1].length)
                this.Function.deltaDecorations({
                  lineNumber: 1,
                  column: ar1[0].length + ar1[1].length + 3
                }, "SS", ar1[2].length)
                continue
              }
            }
            if (ar1.length >= 2) {
              const SS = ar1[0]
              this.Function.deltaDecorations({lineNumber: 1, column: 1}, "SS", SS.length)
              this.Function.deltaDecorations({
                lineNumber: 1,
                column: ar1[0].length + 1
              }, "SS2", array[i].length - SS.length)
            }
            continue
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
          monaco.editor.setModelLanguage(editor.getModel(), language);
        }
      }
      // 监听内容变化事件
      editor.getModel().onDidChangeContent(() => {
        this.Function.SetNewValue(editor.getValue())
        if (!this.IsSetCode) {
          this.IsHasModify = true
          return
        }
        this.IsHasModify = false
        this.IsSetCode = false
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
    SetCode(newContent, presentData) {
      this.IsSetCode = true;
      this.value = newContent
      if (this.Function.setValue) {
        if (Config_Encoding_Current_ISUTF8.value) {
          this.Function.setValue(toUTF8(newContent));
        } else {
          this.Function.setValue(toGBK(newContent));
        }
      }
      if (presentData) {
        this.defaultOpts.presentData = () => {
          presentData()
          this.defaultOpts.presentData = null;
        }
      } else {
        this.defaultOpts.presentData = null
      }
    },
    formatCode() {
      if (this.Function.formatCode) {
        this.Function.formatCode()
      }
    },
    GetRequest() {
      if (!this.IsHasModify) {
        return null;
      }
      const value = this.GetCode();
      return bytesToRequest(value);
    },
    GetResponse() {
      if (!this.IsHasModify) {
        return null;
      }
      const value = this.GetCode();
      return bytesToResponse(value);
    },
    GetRawData() {
      if (!this.IsHasModify) {
        return null;
      }
      const value = this.GetCode()
      return bytesToBase64(StringToBytes(value));
    },
  },
  computed: {
    getTheme() {
      if (this.Function.SetTheme) {
        this.Function.SetTheme(this.MonacoEditorTheme)
      }
      this.$nextTick(() => {
        this.$nextTick(() => {
          if (this.Function.GetCode) {
            if (this.Function.SetNewValue) {
              const code = this.Function.GetCode()
              this.IsSetCode = true
              const IsHasModify = this.IsHasModify
              this.Function.SetNewValue(code)
              this.IsHasModify = IsHasModify
            }
          }
        })
      })
      return true
    }
  }
}
</script>


<style>
.monaco-editor {
  --vscode-editor-background: var(--ag-background-color) !important;
  --vscode-editorGutter-background: var(--ag-background-color) !important; /* 行号的背景色 */
}

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