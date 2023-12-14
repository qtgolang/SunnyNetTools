<template>
  <div>
    <el-radio-group v-model="MustTcpMode" class="ml-4">
      <el-radio label="MustTcp" size="large">全部强制走TCP</el-radio>
      <el-radio label="CancelMustTcp" size="large">自定义配置</el-radio>
    </el-radio-group>
  </div>
  <div v-show="MustTcpMode==='CancelMustTcp'" style="width: 100%;height: 150px;position: relative;top: 3px">
    <TextAge ref="Rules" Name="agent" :Text="DefaultRules" :readOnly="false"/>
  </div>
</template>

<script>

import {CallGoDo} from "../../CallbackEventsOn.js";
import TextAge from "../Agent/TextAge.vue";
import {ElMessage} from "element-plus";

export default {
  components: {TextAge},
  watch: {
    "MustTcpMode"(n, e) {
      CallGoDo("保存强制TCP使用规则", {Data: n})
    }
  },
  data() {
    return {
      MustTcpMode: "CancelMustTcp",
      DefaultRules: ""
    }
  }
  ,
  methods: {
    Save() {
      CallGoDo("保存强制TCP使用规则", {Data: this.$refs.Rules.GetCode()}).then(res => {
        if (res) {
          ElMessage({
            message: "保存成功",
            type: 'success',
          })
        } else {
          ElMessage({
            message: "保存失败",
            type: 'error',
          })
        }
      })
    }
  },
  mounted() {
    this.$refs.Rules.SetSaveFunc(this.Save)
    CallGoDo("CancelMustTcp", null)
  }
}
</script>

<style>

</style>