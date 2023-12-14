<template>
  <div class="ag-header-group-cell-label" :style="backStyle">
    <SetIEProxy/>
    <PartitionOperator/>
    <HookState/>
    <PartitionOperator/>
    <HookState2/>
    <PartitionOperator/>
    <span ref="eLabel">{{ Title }}</span>
    <div class="ag-header-cell-resize" style="position:absolute;height: 20px;width: 20px"></div>

    <SwitchTheme style="right: 0px;position:absolute;cursor:pointer"/>
  </div>
</template>


<script>
import agOption from './option.vue';
import SwitchTheme from './SwitchTheme.vue';
import PartitionOperator from './PartitionOperator.vue';
import SetIEProxy from "./SetIEProxy.vue";
import HookState from "./HookState.vue";
import HookState2 from "./HookState2.vue";

export default {
  components: {
    SetIEProxy,
    agOption, SwitchTheme, PartitionOperator, HookState, HookState2
  },
  computed: {
    backStyle() {
      let c = "height: 30px;background-color: #202020;"
      if (!this.theme) {
        c += "background-color: #f0f0f0;"
      } else {
        c += ""
      }
      return c
    }
  },
  data() {
    return {
      Title: "正在启动...",
      get theme() {
        return window.Theme.IsDark
      },
      set theme(newValue) {
        window.Theme.IsDark = newValue
      }
    }
  }, methods: {},
  mounted() {
    window.vm["Footer"] = this
  }
}
</script>
<style scoped>
</style>