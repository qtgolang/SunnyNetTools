<template>
  &nbsp;&nbsp;&nbsp;
  <div v-show="isSettings===false" style="width: 100px;cursor: pointer; " @click="Click">{{ title }}</div>
  <div v-show="isSettings" style="width: 100px;cursor: pointer; ">正在设置中...</div>
</template>

<script>

import {CallGoDo} from "./CallbackEventsOn.js";

export default {
  data() {
    return {
      state: false,
      title: "未设置系统IE代理",
      isSettings: false
    }
  }, methods: {
    Click() {
      this.state = !this.state
      this.isSettings = true
      CallGoDo("设置IE代理", {Set: this.state}).then(res => {
        this.isSettings = false
        if (!res) {
          return
        }
        if (this.state) {
          this.title = "已设置系统IE代理"
        } else {
          this.title = "未设置系统IE代理"
        }
      })
    }
  }
}
</script>
<style scoped>

</style>