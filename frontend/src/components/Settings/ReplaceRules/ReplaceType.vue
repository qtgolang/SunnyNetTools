<script>
//

export default {
  props: ['params'],
  watch: {
    "value"(m, l) {
      this.Params.data['替换类型'] = m
      this.Params.api.applyTransaction({update: [this.Params.data]});
    },
  },
  data() {
    return {
      Params: null,
      value: "",
    }
  },
  created() {
    this.Params = this.params;
    this.value = this.Params.data['替换类型']
  },
  mounted() {

  }, methods: {}
};
</script>
<template>
  <el-select v-model="value" class="m-2" placeholder="Select" size="small"
             style="position: relative;top: -4px;left: -10px; width: 105px">
    <el-tooltip class="item" effect="dark" content="如果选择此项,那么 源内容及替换内容 都要求填写十六进制字符串"
                placement="top">
      <el-option key="HEX" label="HEX" value="HEX"/>
    </el-tooltip>
    <el-tooltip class="item" effect="dark" content="如果选择此项,那么 源内容及替换内容 都要求填写Base64字符串"
                placement="top">
      <el-option key="Base64" label="Base64" value="Base64"/>
    </el-tooltip>
    <el-tooltip class="item" effect="dark" content="如果选择此项,那么 源内容及替换内容 中的中文字符串都按UTF8编码处理"
                placement="top">
      <el-option key="String(UTF8)" label="String(UTF8)" value="String(UTF8)"/>
    </el-tooltip>
    <el-tooltip class="item" effect="dark" content="如果选择此项,那么 源内容及替换内容 中的中文字符串都按GBK编码处理"
                placement="top">
      <el-option key="String(GBK)" label="String(GBK)" value="String(GBK)"/>
    </el-tooltip>
    <el-tooltip class="item" effect="dark"
                content="如果选择此项,那么 源内容 应填写URL,替换内容应填写本地文件的绝对路径,如果本地文件不存在,则响应空"
                placement="top">
      <el-option key="响应文件" label="响应文件" value="响应文件"/>
    </el-tooltip>

  </el-select>
</template>
