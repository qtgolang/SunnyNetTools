<template>
  <div
      ref="container"
      class="monaco-editor"
      v-show="getTheme"
      :style="`text-align: left;`"

  ></div>
</template>
<script>
import * as monaco from 'monaco-editor'

export default {
  data() {
    return {
      // 主要配置
      defaultOpts: {
        value: '', // 编辑器的值
        theme: 'vs-dark', // 编辑器主题：vs, hc-black, or vs-dark，更多选择详见官网
        roundedSelection: true, // 右侧不显示编辑器预览框
        scrollBeyondLastLine: false, // 禁止滚动超过最后一行
        autoIndent: true, // 自动缩进
        automaticLayout: true,
        formatOnType: true,
        formatOnPaste: true,
        originalEditable: true,
        glyphMargin: false,
        diffViewport: false
      },
      Function: {
        Save: null,
        GetCode: null,
        SetTheme: null
      }
    }
  },
  watch: {},
  mounted() {
    this.init()
  },
  computed: {
    getTheme() {
      const th = window.Theme.IsDark ? "vs-dark" : "vs"
      if (this.Function.SetTheme) {
        this.Function.SetTheme(th)
      }
      return true;
    }
  },
  methods: {
    init() {
      // 初始化container的内容，销毁之前生成的编辑器
      this.$refs.container.innerHTML = ''
      // 生成 diff-editor 对象
      const diffMediator = monaco.editor.createDiffEditor(this.$refs.container, this.defaultOpts)
      const str1 = "Hello Word\r\n" + "123" + "\r\n456"
      const str2 = "Hello Word\r\n" + "" + "\r\n456"
      const lhsModel = monaco.editor.createModel(str1, 'text/plain');
      const rhsModel = monaco.editor.createModel(str2, 'text/plain');
      diffMediator.setModel({
        original: lhsModel,
        modified: rhsModel
      });
      diffMediator.updateOptions({
        diffOverviewRulerOptions: false
      });
      this.monacoEditor = diffMediator
      this.Function.SetTheme = (newTheme) => {
        monaco.editor.setTheme(newTheme);
      }
    },

  }
}
</script>

