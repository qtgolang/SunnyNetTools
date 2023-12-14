<template>
  &nbsp;
  <div class="ag-picker-field ag-labeled ag-label-align-left ag-select ag-group-item "
       ref="familySelect" style="width: 132px;cursor: pointer; ">
    <div ref="eWrapper" class="ag-wrapper ag-picker-field-wrapper ag-picker-collapsed" aria-labelledby="ag-5650-label"
         tabindex="0"
         :style="getBackgroundColor">
      <div ref="eDisplayField" class="ag-picker-field-display" id="ag-5650-display" style="cursor: pointer; ">{{ title }}
      </div>
      <span class="ag-icon ag-icon-small-down"
            unselectable="on"
            role="presentation"></span>
    </div>
    <ComboBoxWindow ref="CB" v-show="Show" :size="getWindowSize" :func="getCallFunc"/>
  </div>


</template>
<script>
import ComboBoxWindow from './ComboBox/Window.vue';

export default {
  components: {
    ComboBoxWindow
  },
  computed: {
    getWindowSize() {
      return this.WindowSize
    },
    getCallFunc() {
      return this.CallFunction
    },
    getBackgroundColor() {
      const res = "flex: 1 1 auto;"
      return res + (this.theme ? "background-color:#202020" : "")
    }
  },
  data() {
    return {
      WindowSize: {
        top: 0,
        left: 0,
        width: 0,
        height: 0
      },
      Show: false,
      CallFunction: null,
      title: "未设置代理",
      get theme() {
        return window.Theme.IsDark
      },
      set theme(newValue) {
        window.Theme.IsDark = newValue
      },
    }
  },
  mounted() {
    const rect = this.$refs.eWrapper.getBoundingClientRect();
    this.WindowSize = {
      top: rect.top - 60,
      left: rect.left,
      width: rect.right - rect.left,
      height: 60
    }
    this.CallFunction = this.SetMode
    window.addEventListener('mousedown', this.handleMouseDown);
  },
  beforeUnmount() {
    // 组件卸载前，移除全局事件监听器
    window.removeEventListener('mousedown', this.handleMouseDown);
  },
  methods: {
    handleMouseDown(event) {
      let bool = false
      let eve = event.target
      while (true) {
        if (eve === this.$refs.familySelect) {
          bool = true
          break
        }
        if (eve.parentNode === null) {
          if (eve.parentElement === null) {
            break
          } else {
            eve = eve.parentElement
          }
        } else {
          eve = eve.parentNode
        }
      }
      this.Show = bool
    },
    SetMode(index) {
      this.Show = false
      if (index === 1) {
        this.title = "已设置系统代理"
      } else if (index === 2) {
        this.title = "未设置系统代理"
      }

    }
  }
}
</script>

<style scoped>
</style>