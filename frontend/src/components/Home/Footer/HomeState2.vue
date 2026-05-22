<template>

  <div ref="Bark" style="display: flex; align-items: center;" @click="ClickEve">
    <div style="cursor: pointer; display: flex; align-items: center;">
      <img :src="getImg" style="height: 20px;"><img>
    </div>
  </div>
</template>

<script>
import {GetImage} from "../../Tools/image.js";
import {Config_IsDark} from "../../config/Config.js";
import {SetBreakMode} from "../../../../bindings/changeme/Service/appmain.js";
import {Tour_Add, Tour_Start} from "../Tour";

export default {
  computed: {
    getImg() {
      return this.img
    }
  },
  mounted() {
    Tour_Add(this.$refs.Bark, 3, "拦截模式", "点击这里切换拦截状态,仅支持HTTP请求的拦截")
  },
  data() {
    return {
      get theme() {
        return Config_IsDark.value
      },
      set theme(newValue) {
        if (Config_IsDark.value !== newValue) Config_IsDark.value = newValue
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
      SetBreakMode(this.state)
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