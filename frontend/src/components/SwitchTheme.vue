<template>
  <div style="height: 20px; width: 20px; cursor:pointer">
    <i v-if="IsDark" @click="SwitchTheme">
      <svg class="v-icon__svg" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" role="img" aria-hidden="true">
        <path
            d="M17.75,4.09L15.22,6.03L16.13,9.09L13.5,7.28L10.87,9.09L11.78,6.03L9.25,4.09L12.44,4L13.5,1L14.56,4L17.75,4.09M21.25,11L19.61,12.25L20.2,14.23L18.5,13.06L16.8,14.23L17.39,12.25L15.75,11L17.81,10.95L18.5,9L19.19,10.95L21.25,11M18.97,15.95C19.8,15.87 20.69,17.05 20.16,17.8C19.84,18.25 19.5,18.67 19.08,19.07C15.17,23 8.84,23 4.94,19.07C1.03,15.17 1.03,8.83 4.94,4.93C5.34,4.53 5.76,4.17 6.21,3.85C6.96,3.32 8.14,4.21 8.06,5.04C7.79,7.9 8.75,10.87 10.95,13.06C13.14,15.26 16.1,16.22 18.97,15.95M17.33,17.97C14.5,17.81 11.7,16.64 9.53,14.5C7.36,12.31 6.2,9.5 6.04,6.68C3.23,9.82 3.34,14.64 6.35,17.66C9.37,20.67 14.19,20.78 17.33,17.97Z"></path>
      </svg>
    </i>
    <i v-if="IsDark===false" @click="SwitchTheme">
      <svg style="fill: white;" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" role="img" aria-hidden="true">
        <path
            d="M12,7A5,5 0 0,1 17,12A5,5 0 0,1 12,17A5,5 0 0,1 7,12A5,5 0 0,1 12,7M12,9A3,3 0 0,0 9,12A3,3 0 0,0 12,15A3,3 0 0,0 15,12A3,3 0 0,0 12,9M12,2L14.39,5.42C13.65,5.15 12.84,5 12,5C11.16,5 10.35,5.15 9.61,5.42L12,2M3.34,7L7.5,6.65C6.9,7.16 6.36,7.78 5.94,8.5C5.5,9.24 5.25,10 5.11,10.79L3.34,7M3.36,17L5.12,13.23C5.26,14 5.53,14.78 5.95,15.5C6.37,16.24 6.91,16.86 7.5,17.37L3.36,17M20.65,7L18.88,10.79C18.74,10 18.47,9.23 18.05,8.5C17.63,7.78 17.1,7.15 16.5,6.64L20.65,7M20.64,17L16.5,17.36C17.09,16.85 17.62,16.22 18.04,15.5C18.46,14.77 18.73,14 18.87,13.21L20.64,17M12,22L9.59,18.56C10.33,18.83 11.14,19 12,19C12.82,19 13.63,18.83 14.37,18.56L12,22Z"></path>
      </svg>
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