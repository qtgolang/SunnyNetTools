<template>
  <div style="height: 18px; width: 18px; cursor:pointer;">
    <i v-if="IsDark" @click="SwitchTheme" style="font-size: 20px;">
      <el-icon><Sunrise /></el-icon>
    </i>
    <i v-if="IsDark===false" @click="SwitchTheme" style="font-size: 20px;">
      <el-icon><Sunny /></el-icon>
    </i>
  </div>
</template>
<script>

import {CallGoDo} from "./CallbackEventsOn.js";

export default {
  components: {},
  computed: {
    IsDark() {
      const m = this.Theme
      const htmlElement = document.querySelector('html');
      htmlElement.classList.remove('dark'); // 移除 'dark' 类
      if (m) {
        htmlElement.classList.add('dark');
      } else {
        htmlElement.classList.add('light');
      }
      this.$nextTick(() => {
        const ns = document.getElementsByClassName("el-divider__text")
        for (let i = 0; i < ns.length; i++) {
          if (m) {
            ns[i].style.backgroundColor = "rgb(45,52,54)"
            ns[i].style.Color = "#ffffff"
          } else {
            ns[i].style.backgroundColor = "rgb(255,255,255)"
            ns[i].style.Color = "#000000"
          }
        }
      })
      if (window.vm.List != null) {
        window.vm.List.RefreshRenderedNodes()
      }
      return !m
    }
  },
  data() {
    return {
      Theme: window.Theme.IsDark,
    }
  }, methods: {
    SwitchTheme() {
      if (window.Theme.IsDark) {
        this.Theme = false
        window.Theme.IsDark = false
      } else {
        this.Theme = true
        window.Theme.IsDark = true
      }
      CallGoDo("更新主题", {Dark: window.Theme.IsDark})
    }
  }, mounted() {
    window.vm.Theme = this
  }
}
</script>
<style scoped>
</style>