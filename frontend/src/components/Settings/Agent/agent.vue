<template>
  <div class="mt-4">
    <el-input v-model="input1" placeholder="请输入上游代理地址 (账号:密码@IP:端口) 例如(admin:pass@127.0.0.1:8888)"
              class="input-with-select" :disabled="inputDisabled">
      <template #prepend>
        <el-select v-model="select" placeholder="Select" style="width: 115px">
          <el-option label="Socket5://" value="1"/>
          <el-option label="HTTP://" value="2"/>
        </el-select>
      </template>
      <template #append>
        <el-switch @click="switchClick" v-model="switchValue" active-text="开启" inactive-text="关闭"/>
      </template>
    </el-input>
    <el-tag class="ml-2" type="info" style="width: 100%;position: relative;top: 3px;text-align: left">
      以下地址,不使用上游代理
    </el-tag>
    <div style="width: 100%;height: 150px;position: relative;top: 3px">
      <TextAge ref="agentInput" Name="agent" :Text="DefaultRules" :readOnly="false"/>
    </div>
  </div>
</template>

<script>

import TextAge from "./TextAge.vue";
import {CallGoDo} from "../../CallbackEventsOn.js";
import {ElMessage} from "element-plus";

export default {
  components: {TextAge},
  watch: {},
  data() {
    return {
      inputDisabled: false,
      switchValue: false,
      input1: "",
      select: "Socket5://",
      DefaultRules: ""
    }
  }, mounted() {
    this.$refs.agentInput.SetSaveFunc(this.Save)
  }, methods: {
    switchClick() {
      const n = this.switchValue
      let uri = ''
      if (n) {
        if (this.input1.trim() === '') {
          ElMessage({
            message: "请先输入上游代理地址",
            type: 'error',
          })
          this.switchValue = false
          this.inputDisabled = false
          return
        }
        uri = this.select + this.input1
      }
      CallGoDo("设置上游代理", {Data: uri, Set: n}).then(res => {
        if (res) {
          if (!n) {
            ElMessage({
              message: "已取消",
              type: 'success',
            })
            this.inputDisabled = false
            return
          }
          ElMessage({
            message: "设置上游代理成功",
            type: 'success',
          })
          this.inputDisabled = true
          return
        }
        ElMessage({
          message: "请检查你输入的上游代理格式是否正确",
          type: 'error',
        })
        this.switchValue = false
        this.inputDisabled = false
      })
    },
    Save() {
      CallGoDo("保存上游代理使用规则", {Data: this.$refs.agentInput.GetCode()}).then(res => {
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
  }
}
</script>

<style>
.input-with-select .el-input-group__prepend {
  background-color: var(--el-fill-color-blank);
}
</style>
