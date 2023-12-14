<template>

  <div style="display: flex; align-items: center;" @click="ClickEve">
    <div style="cursor: pointer; display: flex; align-items: center;">
      <img :src="getImg" style="height: 20px;"><img>
    </div>
  </div>
</template>

<script>
import {GetImage} from "./image.js";
import {CallGoDo} from "./CallbackEventsOn.js";

export default {
  computed: {
    getImg() {
      if (this.theme) {
        this.applyImg()
      } else {
        this.applyImg()
      }
      return this.img
    }
  },
  data() {
    return {
      get theme() {
        return window.Theme.IsDark
      },
      set theme(newValue) {
        window.Theme.IsDark = newValue
      },
      state: 0,
      img: GetImage("空白")
    }
  }, methods: {
    ClickEve() {
      this.state++
      if (this.state > 2) {
        this.state = 0
      }
      this.applyImg()
    },
    applyImg() {
      CallGoDo("设置断点模式", {break: this.state})
      if (this.state === 0) {
        this.img = GetImage("空白")
      } else if (this.state === 1) {
        this.img = GetImage("拦截上行")
      } else {
        this.img = GetImage("拦截下行")
      }
    }
  }
}

</script>
<style scoped>

</style>