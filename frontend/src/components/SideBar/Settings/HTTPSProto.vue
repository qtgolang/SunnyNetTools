<script>
import Table from "../../Tools/table.vue";
import {Config_IsDark, Config_IsRest, Config_MonacoEditorTheme} from "../../config/Config";
import * as monaco from "monaco-editor";
import {ElNotification} from "element-plus";
import {
  GetHTTPSProto,
  GetRandomJa3,
  GetSendIsHTTP1,
  SetHTTPSProto,
  SetRandomJa3,
  SetSendIsHTTP1
} from "../../../../bindings/changeme/Service/appmain";
import {addMonacoEditorFontSize, RemoveMonacoEditorMenu} from "../../config/editorContextMenu";
import {attachMcpSettingsReload} from "../../config/mcpSettingsSync.js";

const H1 = "http/1.1"
const H2 = 'h2'
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
      isInit: false,
      RandomJa3: false,
      get theme() {
        return Config_MonacoEditorTheme.value
      },
      set theme(newValue) {
        Config_MonacoEditorTheme.value = newValue
      },

      value: '',
      protoValue: H2,
      proto: [
        {
          value: H1,
          label: '仅使用 HTTP/1.1 发送',
        },
        {
          value: H2,
          label: 'HTTP/2.0 优先',
        },
      ],
      options: [
        {
          value: 'Firefox',
          label: 'Firefox',
          Config: "{\"ConnectionFlow\":12517377,\"HeaderPriority\":{\"StreamDep\":13,\"Exclusive\":false,\"Weight\":41},\"Priorities\":[{\"PriorityParam\":{\"StreamDep\":0,\"Exclusive\":false,\"Weight\":200},\"StreamID\":3},{\"PriorityParam\":{\"StreamDep\":0,\"Exclusive\":false,\"Weight\":100},\"StreamID\":5},{\"PriorityParam\":{\"StreamDep\":0,\"Exclusive\":false,\"Weight\":0},\"StreamID\":7},{\"PriorityParam\":{\"StreamDep\":7,\"Exclusive\":false,\"Weight\":0},\"StreamID\":9},{\"PriorityParam\":{\"StreamDep\":3,\"Exclusive\":false,\"Weight\":0},\"StreamID\":11},{\"PriorityParam\":{\"StreamDep\":0,\"Exclusive\":false,\"Weight\":240},\"StreamID\":13}],\"PseudoHeaderOrder\":[\":method\",\":path\",\":authority\",\":scheme\"],\"Settings\":{\"1\":65536,\"4\":131072,\"5\":16384},\"SettingsOrder\":[1,4,5]}",
        },
        {
          value: 'Opera',
          label: 'Opera',
          Config: "{\"ConnectionFlow\":15663105,\"HeaderPriority\":null,\"Priorities\":null,\"PseudoHeaderOrder\":[\":method\",\":authority\",\":scheme\",\":path\"],\"Settings\":{\"1\":65536,\"3\":1000,\"4\":6291456,\"6\":262144},\"SettingsOrder\":[1,3,4,6]}",
        },
        {
          value: 'Safari_IOS_17_0',
          label: 'Safari_IOS_17_0',
          Config: "{\"ConnectionFlow\":10485760,\"HeaderPriority\":null,\"Priorities\":null,\"PseudoHeaderOrder\":[\":method\",\":scheme\",\":path\",\":authority\"],\"Settings\":{\"2\":0,\"3\":100,\"4\":2097152},\"SettingsOrder\":[2,4,3]}",
        },
        {
          value: 'Safari_IOS_16_0',
          label: 'Safari_IOS_16_0',
          Config: "{\"ConnectionFlow\":10485760,\"HeaderPriority\":null,\"Priorities\":null,\"PseudoHeaderOrder\":[\":method\",\":scheme\",\":path\",\":authority\"],\"Settings\":{\"3\":100,\"4\":2097152},\"SettingsOrder\":[4,3]}",
        },
        {
          value: 'Chrome_117_120_124',
          label: 'Chrome_117_120_124',
          Config: "{\"ConnectionFlow\":10485760,\"HeaderPriority\":null,\"Priorities\":null,\"PseudoHeaderOrder\":[\":method\",\":scheme\",\":path\",\":authority\"],\"Settings\":{\"3\":100,\"4\":4194304},\"SettingsOrder\":[4,3]}",
        },
        {
          value: 'Chrome_106_116',
          label: 'Chrome_106_116',
          Config: "{\"ConnectionFlow\":15663105,\"HeaderPriority\":null,\"Priorities\":null,\"PseudoHeaderOrder\":[\":method\",\":authority\",\":scheme\",\":path\"],\"Settings\":{\"1\":65536,\"2\":0,\"4\":6291456,\"6\":262144},\"SettingsOrder\":[1,2,4,6]}",
        },
        {
          value: 'Chrome_103_105',
          label: 'Chrome_103_105',
          Config: "{\"ConnectionFlow\":15663105,\"HeaderPriority\":null,\"Priorities\":null,\"PseudoHeaderOrder\":[\":method\",\":authority\",\":scheme\",\":path\"],\"Settings\":{\"1\":65536,\"3\":1000,\"4\":6291456,\"6\":262144},\"SettingsOrder\":[1,3,4,6]}",
        },
      ],
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
    }
  },
  computed: {
    getTheme() {
      return this.theme
    },
    isHttp2() {
      return this.protoValue === H2
    },
  },
  mounted() {
    this.init()
    attachMcpSettingsReload("https", () => this.reloadHttpsFromBackend());
    GetRandomJa3().then(ok => {
      this.RandomJa3 = ok
      GetHTTPSProto().then(res => {
        this.Function.setValue(res)
        setTimeout(() => {
          this.Function.formatCode()
          setTimeout(() => {
            this.isInit = true;
            this.IsHasModify = false
          }, 1000)
        }, 1000)

      })
    })
    GetSendIsHTTP1().then(ok => {
      this.protoValue = ok ? H1 : H2
    })
  },
  watch: {
    protoValue(n) {
      SetSendIsHTTP1(n === H1)
    },
    theme(n) {
      if (this.Function.SetTheme) {
        this.Function.SetTheme(n)
      }
    },
    isRest() {
      this.isInit = false
      GetRandomJa3().then(ok => {
        this.RandomJa3 = ok
        GetHTTPSProto().then(res => {
          this.Function.setValue(res)
          setTimeout(() => {
            this.Function.formatCode()
            setTimeout(() => {
              this.isInit = true;
              this.IsHasModify = false
            }, 1000)
          }, 1000)

        })
      })
    },
    value(n) {
      if (n !== "") {
        setTimeout(() => {
          this.$nextTick(() => {
            this.value = "";
          })
        }, 100)
        for (let i = 0; i < this.options.length; i++) {
          if (this.options[i].value === n) {
            this.Function.setValue(this.options[i].Config)
            this.Function.formatCode()
            break
          }
        }
        ElNotification({
          position: 'bottom-right',
          message: '您可以修改其中的数值,实现修改HTTP2指纹的效果\n\n修改后,别忘了 `应用修改` 哦!!!',
          type: 'warning',
          customClass: 'multiline-message'
        })
        setTimeout(() => {
          ElNotification({
            position: 'bottom-right',
            message: '已经载入 HTTP2 指纹模板:\n[ ' + n + ' ]',
            type: 'success',
            customClass: 'multiline-message'
          })
        }, 100)
      }
    },
    RandomJa3(n) {
      if (this.isInit) {
        SetRandomJa3(n).then(() => {
          ElNotification({
            position: 'bottom-right',
            message: '更改:随机JA3指纹 配置已保存',
            type: 'success',
            customClass: 'multiline-message'
          })
        })
      }
    }
  },
  methods: {
    IsRemoveMenu(value, zd) {
      let Menu = [
        "转到符号", "更改所有匹配项", "命令面板", "保存修改"
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
      RemoveMonacoEditorMenu(editor,  this.IsRemoveMenu, ()=>{return this.defaultOpts.readOnly;});

      this.Function.GetCode = () => {
        return editor.getValue()
      }
      this.Function.setValue = (newContent) => {
        editor.setValue(newContent);
        editor.setScrollTop(0);
      }
      this.Function.formatCode = () => {
        editor?.getAction('editor.action.formatDocument').run()
      }
      this.Function.SetTheme = (newTheme) => {
        monaco.editor.setTheme(newTheme);
      }

      // 监听内容变化事件
      editor.getModel().onDidChangeContent(() => {
        if (!this.IsHasModify) {
          if (this.isInit) {
            ElNotification({
              position: 'bottom-right',
              message: '修改后,别忘了 `应用修改` 哦!!!',
              type: 'warning',
              customClass: 'multiline-message'
            })
          }
        }
        this.IsHasModify = true
      });
    }
    ,
    reloadHttpsFromBackend() {
      this.isInit = false;
      return GetRandomJa3().then((ok) => {
        this.RandomJa3 = ok;
        return GetSendIsHTTP1().then((h1) => {
          this.protoValue = h1 ? H1 : H2;
          return GetHTTPSProto().then((res) => {
            this.Function.setValue(res);
            setTimeout(() => {
              this.Function.formatCode();
              setTimeout(() => {
                this.isInit = true;
                this.IsHasModify = false;
              }, 500);
            }, 300);
          });
        });
      });
    },
    Save() {
      const data = this.Function.GetCode();
      try {
        JSON.parse(data)
      } catch (e) {
        ElNotification({
          position: 'bottom-right',
          message: '保存HTTP2指纹失败,请检查数据格式是否正确',
          type: 'error',
          customClass: 'multiline-message'
        })
        return
      }
      SetHTTPSProto(data).then(() => {
        ElNotification({
          position: 'bottom-right',
          message: 'HTTP2指纹已保存！',
          type: 'success',
          customClass: 'multiline-message'
        })
      })
    },
  }
}
</script>

<template>
  <div style="display: flex; flex-direction: column; gap: 5px; margin: 5px;">
    <div style="display: flex; position: relative; gap:30px; margin-left: 23px">
      <div
          class="ag-labeled ag-label-align-left ag-toggle-button ag-input-field ag-group-title-bar ag-charts-format-sub-level-group-title-bar ag-unselectable ag-selected"
          style="display: flex;width: 100%">
        <div class="ag-input-field-label ag-label ag-toggle-button-label">
          HTTPS请求发送时使用的协议：
        </div>

        <el-select v-model="protoValue" placeholder="HTTP/2.0 优先" style="" size="small">
          <el-option
              v-for="item in proto"
              :key="item.value"
              :label="item.label"
              :value="item.value"
          />
        </el-select>
      </div>
    </div>
    <div style="display: flex; position: relative; gap:30px; margin-left: 23px">
      <el-tooltip placement="bottom">
        <template #content>
          <div style="white-space: normal; line-height: 1.4;">
            开启后所有HTTPS的JA3指纹将随机变化
          </div>
        </template>
        <div
            class="ag-labeled ag-label-align-left ag-toggle-button ag-input-field ag-group-title-bar ag-charts-format-sub-level-group-title-bar ag-unselectable ag-selected">
          <div class="ag-input-field-label ag-label ag-toggle-button-label" style="margin-right: 10px">
            随机JA3
          </div>
          <el-switch v-model="RandomJa3" size="small" style="margin-right: 10px"/>
        </div>
      </el-tooltip>
      <div
          v-show="isHttp2"
          class="ag-labeled ag-label-align-left ag-toggle-button ag-input-field ag-group-title-bar ag-charts-format-sub-level-group-title-bar ag-unselectable ag-selected"
          style="display: flex;width: 100%">
        <div class="ag-input-field-label ag-label ag-toggle-button-label">
          HTTP2指纹设置
        </div>

        <el-select v-model="value" placeholder="选择模板" style="" size="small">
          <el-option
              v-for="item in options"
              :key="item.value"
              :label="item.label"
              :value="item.value"
          />
        </el-select>
      </div>
    </div>
    <div
        v-show="isHttp2"
        ref="container"
        class="monaco-editor"
        style="text-align: left;height: 400px;width: 100%;"
        :drak="getTheme"
    ></div>
  </div>
</template>
<style>
.multiline-message {
  white-space: pre-line;
}
</style>