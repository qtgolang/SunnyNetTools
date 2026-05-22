<template>
  <div style="width: 100%;height: 100%;display: flex; align-items: center;">
    <McpStatus/>
    <PartitionOperator/>
    <SetIEProxy/>
    <PartitionOperator/>
    <HookState :isHideHook="isHideHook"/>
    <PartitionOperator/>
    <HookState2/>
    <PartitionOperator/>

    <div v-show="!IsShowExportProgress">
      <span ref="eLabel" style="font-size: 15px;">{{ Title }}</span>
      <div class="ag-header-cell-resize" style="position:absolute;height: 20px;width: 20px"></div>
    </div>
    <div v-show="IsShowExportProgress" style="width: 100%">
      <el-progress
          :text-inside="true"
          :stroke-width="24"
          :percentage="Progress"
          :duration="duration"
          status="warning"
          striped
          striped-flow
      >
        <template #default="{ percentage }">
          <span class="percentage-label">{{ Config_IsShowExportProgress.value }}：</span>
          <span class="percentage-value">{{ percentage }}%</span>
        </template>
      </el-progress>
    </div>

  </div>
</template>


<script>
import PartitionOperator from './PartitionOperator.vue';
import SetIEProxy from "./SetIEProxy.vue";
import McpStatus from "./McpStatus.vue";
import HookState from "./HomeState.vue";
import HookState2 from "./HomeState2.vue";
import {Config_IsDark, Config_IsShowExportProgress, Config_Status_Info} from "../../config/Config.js";
import {Events} from "@wailsio/runtime";

export default {
  props: ["isHideHook"],
  components: {
    McpStatus, SetIEProxy, PartitionOperator, HookState, HookState2
  },
  computed: {
    Config_IsShowExportProgress() {
      return Config_IsShowExportProgress
    },
    duration() {
      return Math.floor(this.Progress / 10)
    },
    backStyle() {
      let c = "height: 30px;background-color: #202020;position: relative;z-index: 999999;"
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
      Progress: 10,
      get IsShowExportProgress() {
        return Config_IsShowExportProgress.value !== ""
      },
      set IsShowExportProgress(a) {
        Config_IsShowExportProgress.value = a
      },
      get Title() {
        return Config_Status_Info.value;
      },
      set Title(value) {
        Config_Status_Info.value = value;
      },
      get theme() {
        return Config_IsDark.value
      },
      set theme(newValue) {
        if (Config_IsDark.value !== newValue) Config_IsDark.value = newValue
      }
    }
  }, methods: {},
  mounted() {
    Events.On("ExportProgress", (obj) => {
      this.Progress = parseInt(obj.data[0])
    })
    window.Footer = this;
  }
}
</script>
<style scoped>
</style>