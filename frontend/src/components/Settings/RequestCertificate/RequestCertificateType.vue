<script>
//

export default {
  props: ['params'],
  watch: {
    "value"(m, l) {
      this.Params.data['使用规则'] = m
    },
  },
  data() {
    return {
      Params: null,
      value: "解析及发送",
    }
  },
  created() {
    this.Params = this.params;
    this.value = this.Params.data['使用规则']
  },
  mounted() {

  }, methods: {}
};
</script>
<template>
  <el-select v-model="value" class="m-2" placeholder="Select" size="small"
             style="position: relative;top: -4px;left: -10px; width: 105px">

    <el-tooltip class="item" effect="dark" content="请求发送到服务器时使用该证书,并且与客户端建立的通信也使用该证书"
                placement="top">
      <el-option key="解析及发送" label="解析及发送" value="解析及发送"/>
    </el-tooltip>

    <el-tooltip class="item" effect="dark" content="仅与客户端建立的通信时使用该证书" placement="top">
      <el-option key="仅解析" label="仅解析" value="仅解析"/>
    </el-tooltip>

    <el-tooltip class="item" effect="dark" content="仅请求发送到服务器时使用该证书" placement="top">
      <el-option key="仅发送" label="仅发送" value="仅发送"/>
    </el-tooltip>

  </el-select>
</template>
