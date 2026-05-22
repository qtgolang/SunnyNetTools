<script>

import {Window} from "@wailsio/runtime";
import {Config_GOOS_IsWindows, Config_IsDark} from "../config/Config";
import VTitlebar from "./VUETitlebar/vueTitlebar.vue";

export default {
  components: {VTitlebar},
  props: ["Title"],
  data() {
    return {
      get IsWindows() {
        return Config_GOOS_IsWindows.value
      },
      set IsWindows(value) {
        Config_GOOS_IsWindows.value = value
      },
      get theme() {
        return Config_IsDark.value
      },
      set theme(newValue) {
        if (Config_IsDark.value !== newValue) Config_IsDark.value = newValue
      },
    }
  }
  , computed: {
    getTheme() {
      return Config_IsDark.value ? "dark" : "light"
    },
    platform() {
      return this.IsWindows ? "windows" : "darwin"
    }
  },
  methods: {
    handleClose() {
      Window.Close()
    }
    ,
    handleMinimize() {
      Window.Minimise()
    }
    ,
    handleMaximize() {
      Window.ToggleMaximise()
    }
    ,
    handleRestore() {

    }
    ,
  }
}
</script>

<template>
  <div style="width: 100%;height:30px;display: flex;position: relative;align-items: center;cursor: default;user-select: none;margin-top: -1px">
    <v-titlebar
        :theme="getTheme"
        :show-icon="false"
        :platform="platform"
        ref="header"
        @close="handleClose"
        @minimize="handleMinimize"
        @maximize="handleMaximize"
        @unmaximize="handleRestore"
    >
      <template v-slot:title>
        <div style="display: flex; position: relative; width: 100%;">
          <!-- 居中的文字 -->
          <span style="width: 100%;justify-content: center;text-align: center;align-content: center;margin-left: 90px;--wails-draggable: drag"  @dblclick="handleMaximize">
            {{ Title }}
          </span>
          <div v-if="!this.IsWindows" style="width: 125px"></div>
        </div>
      </template>
    </v-titlebar>

  </div>
</template>
