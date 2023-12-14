<template>
  <div
      ref="container"
      class="monaco-editor"
      style="text-align: left;height: 100%;position: relative"
  ></div>
</template>
<script>
import * as monaco from 'monaco-editor'

function IsRemoveMenu(value, zd) {

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
      // 主要配置
      defaultOpts: {
        value: '', // 编辑器的值
        Language: "go",//默认语言
        theme: 'vs-dark', // 编辑器主题：vs, hc-black, or vs-dark，更多选择详见官网
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
        readOnly: false
      },
      Function: {
        Save: null,
        GetCode: null,
        setValue: null,
        SetReadOnly: null,
        setModelLanguage: null
      }
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

      editor.onContextMenu(() => {
        this.$nextTick(() => {
          const asbb = document.querySelectorAll(".shadow-root-host")
          for (let i = 0; i < asbb.length; i++) {
            const as = asbb[i]
            if (as !== null && as !== undefined) {
              const aacs = as.shadowRoot.querySelectorAll(".actions-container")[0]
              if (aacs) {
                const aas = aacs.childNodes
                let bl = false
                while (true) {
                  for (let i = 0; i < aas.length; i++) {
                    for (let i = 0; i < aas.length; i++) {
                      const v = aas[i].innerHTML
                      if (i === 0 || i === aas.length - 1) {
                        if (v.indexOf("action-label codicon separator disabled") !== -1) {
                          aas[i].parentNode.removeChild(aas[i]);
                          bl = true
                          break
                        }
                      }
                      if (IsRemoveMenu(v, this.defaultOpts.readOnly)) {
                        aas[i].parentNode.removeChild(aas[i]);
                        bl = true
                        break
                      }
                    }
                  }
                  if (bl === true) {
                    bl = false
                    continue
                  }
                  break
                }
                const menu = as.shadowRoot.querySelectorAll(".monaco-menu-container")[0]
                if (menu) {
                  let mouseX = event.clientX + window.scrollX;
                  let mouseY = event.clientY + window.scrollY;
                  if (mouseY + menu.offsetHeight > window.innerHeight) {
                    mouseY = window.innerHeight - menu.offsetHeight
                  }
                  if (mouseX + menu.offsetWidth - 100 > window.innerWidth) {
                    mouseX = window.innerWidth - (menu.offsetWidth - 100)
                  }
                  menu.style.top = mouseY + "px"
                  menu.style.left = mouseX + "px"
                }
              }
            }
          }
        });
      });
      this.Function.GetCode = () => {
        return editor.getValue()
      }
      this.Function.SetReadOnly = (readOnlyState) => {
        editor.updateOptions({
          readOnly: readOnlyState
        });
      }
      this.Function.setModelLanguage = (lang) => {
        monaco.editor.setModelLanguage(editor.getModel(), lang);
      }
      this.Function.setValue = (newContent) => {
        editor.setValue(newContent);
      }
      // 监听内容变化事件
      editor.getModel().onDidChangeContent(() => {
      });

      monaco.editor.setModelLanguage(editor.getModel(), this.defaultOpts.Language);
      window.VsCodeEdit.JavaScriptEdit = true
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
    }
  }
}
</script>

