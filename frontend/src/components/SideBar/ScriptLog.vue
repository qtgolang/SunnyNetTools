<script>

import {Config_MonacoEditorTheme, Config_SunnyNetIsStart} from "../config/Config";
import * as monaco from "monaco-editor";
import {Events} from "@wailsio/runtime";
import {ElNotification} from "element-plus";
import {GetPort} from "../../../bindings/changeme/Service/appmain";
import {OpenTools} from "../CallbackEventsOn";
import {addMonacoEditorFontSize, RemoveMonacoEditorMenu} from "../config/editorContextMenu";

export default {
  data() {
    return {
      get theme() {
        return Config_MonacoEditorTheme.value
      },
      set theme(newValue) {
        Config_MonacoEditorTheme.value = newValue
      },
      maxSize: 1000,
      fifoArray: [],
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
        readOnly: true,
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
    }
  },
  methods: {
    IsRemoveMenu(value) {
      let Menu = [
        "转到符号", "更改所有匹配项", "命令面板", "保存修改", "格式化"
      ]
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
      editor.addAction({
        id: 'SaveConfiguration',
        label: '刷新日志',
        contextMenuGroupId: 'my-commands2',
        contextMenuOrder: Number.MAX_SAFE_INTEGER,
        run: this.Refresh,
      })
      editor.addAction({
        id: 'SaveConfigurationEXP',
        label: '清空日志',
        contextMenuGroupId: 'my-commands2',
        contextMenuOrder: Number.MAX_SAFE_INTEGER,
        run: () => {
          this.fifoArray = []
          this.Refresh()
        },
      })
      editor.addCommand(monaco.KeyCode.F1, function () {
        //注册一个F1快捷键，但是什么也不做，用来阻止弹出命令列表
      });
      RemoveMonacoEditorMenu(editor, this.IsRemoveMenu, () => {
        return this.defaultOpts.readOnly;
      });
      this.Function.GetCode = () => {
        return editor.getValue()
      }
      this.Function.setValue = (newContent) => {
        editor.setValue(newContent);
        // 获取 model
        const model = editor.getModel();

        // 获取最后一行和最后一列（光标放在末尾）
        const lineNumber = model.getLineCount();
        const column = model.getLineMaxColumn(lineNumber);

        // 设置光标位置
        editor.setPosition({lineNumber, column});

        // 可选：滚动视图跟随光标
        editor.revealPositionInCenter({lineNumber, column});
      }
      this.Function.formatCode = () => {
        editor?.getAction('editor.action.formatDocument').run()
      }
      this.Function.SetTheme = (newTheme) => {
        monaco.editor.setTheme(newTheme);
      }
    },
    addItem(item) {
      this.fifoArray.push(item);
      if (this.fifoArray.length > this.maxSize) {
        this.fifoArray.shift();
      }
    },
    Refresh() {
      let s = ''
      const max = this.fifoArray.length
      for (let i = 0; i < max; i++) {
        s += "日志时间:" + this.fifoArray[i].Time + "\n"
        if (this.fifoArray[i].Info.length < 1) {
          s += "日志内容:无" + "\n"
          continue
        }
        for (let n = 0; n < this.fifoArray[i].Info.length; n++) {
          s += "日志内容-参数(" + (n + 1) + "):" + this.fifoArray[i].Info[n] + "\n"
        }
        s += "\n"
      }
      if (s === "") {
        s = new Date().toLocaleString() + "\n\n暂无日志:你可以稍后右键刷新日志再试"
      } else {
        s += "--------------------------------------\n刷新时间:" + new Date().toLocaleString() + "\n你可以稍后继续右键刷新日志\n-> 只保留最近 " + this.maxSize + " 条日志\n-> 当前共: " + max + " 条日志\n--------------------------------------\n"
      }
      this.Function.setValue(s)
    },
    openScriptEdit() {
      if (Config_SunnyNetIsStart.value === false) {
        ElNotification({
          position: 'bottom-right',
          message: '打开脚本编辑失败\n您当前程序的工作端口未启动成功\n请修改端口后再试!!',
          type: 'warning',
          customClass: 'multiline-message'
        })
        return
      }
      GetPort().then(port => {
        const url = 'http://localhost:' + port + '/SunnyNetScriptEdit';
        OpenTools("脚本代码", true, url)
        return
        const windowName = 'my-ScriptEdit-popup';
        const width = 1000;
        const height = 600;
        const left = (window.screen.width - width) / 2;
        const top = (window.screen.height - height) / 2;
        const features = `width=${width},height=${height},left=${left},top=${top},resizable=yes,scrollbars=yes`;
        const win = window.open(url, windowName, features);
        if (win) {
          win.focus(); // 把窗口提到前面
        } else {
          ElNotification({
            position: 'bottom-right',
            message: '打开证书安装失败\n请允许弹出窗口或关闭浏览器拦截器!!',
            type: 'warning',
            customClass: 'multiline-message'
          })
        }
      })
    }
  },
  computed: {
    getTheme() {
      return this.theme
    }
  },
  watch: {
    theme(newValue) {
      if (this.Function.SetTheme) {
        this.Function.SetTheme(newValue)
      }
    }
  },
  mounted() {
    this.init()
    this.Function.setValue("你需要右键刷新日志")
    Events.On("addScriptLog", (obj) => {
      try {
        let o = obj.data;
        o.forEach((v) => {
          this.addItem(v)
        })
      } catch (e) {
      }
    })
  }
}
</script>

<template>
  <div style="height: 100%;width: 100%;">
    <div style="height: 30px;width: 100%;margin-top: 5px;margin-bottom: 5px">
      <el-button style="height: 100%;width: 100%;" size="small" @click="openScriptEdit">打开脚本编辑</el-button>
    </div>
    <div
        ref="container"
        class="monaco-editor"
        style="text-align: left;height: calc(100% - 40px);width: 100%;"
        :drak="getTheme"
    ></div>
  </div>
</template>
