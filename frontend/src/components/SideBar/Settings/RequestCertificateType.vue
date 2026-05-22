<script>
import {ElMessage} from "element-plus";

export default {
  props: ['params'],
  watch: {
    "value"(m, l) {
      if (m === this.value2) {
        return
      }
      this.Params.data['使用规则'] = m
      ElMessage.warning("修改使用规则:请右键重新载入证书")
    },
  },
  data() {
    return {
      Params: null,
      value: "解析及发送",
      value2: "解析及发送",
    }
  },
  created() {
    this.Params = this.params;
    this.Params = this.params;
    this.value2 = this.Params.data['使用规则']
    this.value = this.Params.data['使用规则']
  },
  mounted() {

  }, methods: {}
};
</script>
<template>
  <el-select v-model="value" class="m-2" placeholder="Select" size="small"
             style="position: relative;top: -2px;left: -8px; width: 105px">

    <el-tooltip class="item" effect="dark" content="请求发送到服务器时使用该证书,并且与客户端建立的通信也使用该证书"
                placement="right">
      <el-option key="解析及发送" label="解析及发送" value="解析及发送"/>
    </el-tooltip>

    <el-tooltip class="item" effect="dark" content="仅与客户端建立的通信时使用该证书" placement="right">
      <el-option key="仅解析" label="仅解析" value="仅解析"/>
    </el-tooltip>

    <el-tooltip class="item" effect="dark" content="仅请求发送到服务器时使用该证书" placement="right">
      <el-option key="仅发送" label="仅发送" value="仅发送"/>
    </el-tooltip>

  </el-select>
</template>
