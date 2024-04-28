<template>

  <div style="display: flex; align-items: center;" @click="ClickEve">
    <div style="cursor: pointer; display: flex; align-items: center;">
      <img :src="img" style="height: 20px;"><img>
      &nbsp;
      {{ Title }}
    </div>
  </div>
</template>

<script>
import {GetImage} from "./image.js";
import {CallGoDo} from "./CallbackEventsOn.js";
window.isHideHook = false;
export default {
  data() {
    return {
      state: true,
      Title: "正在捕获",
      img: GetImage("开始捕获")
    }
  }, methods: {
    ClickEve() {
      this.state = !this.state
      this.img = GetImage(this.state ? "开始捕获" : "停止捕获")
      this.Title = this.state ? "正在捕获" : "隐藏捕获"
      if (this.state) {
        document.getElementById("HookMessageText").innerText = "还没有捕获到数据"
        window.isHideHook = false;
      } else {
        document.getElementById("HookMessageText").innerText = "您隐藏了捕获数据"
        window.isHideHook = true;
      }
      CallGoDo("工作状态", {State: this.state})
    }
  }
}
</script>
<style scoped>

</style>