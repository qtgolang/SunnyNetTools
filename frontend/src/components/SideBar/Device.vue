<script>

import {
  GOOS,
  IsLoadDevice,
  LoadDevice,
  ProcessAny,
  SetDeviceStopUpdate
} from "../../../bindings/changeme/Service/appmain";
import {ElNotification} from "element-plus";
import MustTcp from "./Settings/MustTcp.vue";
import Host from "./Settings/Host.vue";
import BaseSettings from "./Settings/BaseSettings.vue";
import Keys from "./Settings/Keys.vue";
import HTTPSProto from "./Settings/HTTPSProto.vue";
import Way from "./Settings/Way.vue";
import Name from "./Device/Name.vue";
import Pid from "./Device/Pid.vue";
import {attachMcpDeviceRoot} from "../config/mcpDeviceSync.js";

export default {
  components: {Pid, Name, Way, HTTPSProto, Keys, BaseSettings, Host, MustTcp},
  data() {
    return {
      LoadDrive: false,//是否已经加载
      DriveLoading: false,//加载中
      LoadMode: "Tun",
      Show: {
        Name: true,
        Pid: false,
      },
      ProcessesAny: false,
      _root_width: 0,
      _root_height: 0,
      IsWindows: false,
    }
  },
  watch: {
    "Show.Name"(n) {
      if (n) {
        this.Show.Pid = false
      }
    },
    "Show.Pid"(n) {
      if (n) {
        this.Show.Name = false
      }
    },
    ProcessesAny(n) {
      SetDeviceStopUpdate(n)
      ProcessAny(n, false)
      this.$refs.devName.Empty()
      this.$refs.devPid.deselectAll()
    }
  },
  methods: {
    AnyStopNetwork() {
      ProcessAny(this.ProcessesAny, true)
    },
    loadDrive() {
      this.DriveLoading = true
      IsLoadDevice().then(ok => {
        if (!ok) {
          const mode = this.LoadMode === "NFAPI" ? 1 : (this.LoadMode === "Proxifier" ? 0 : 2)
          LoadDevice(mode).then(ok2 => {
            if (!ok2) {
              this.LoadDrive = false
              this.DriveLoading = false
              ElNotification({
                position: 'bottom-right',
                message: '载入驱动失败\n\n请检查是否已管理员权限运行?',
                type: 'error',
                customClass: 'multiline-message'
              })
              return
            }
            this.LoadDrive = true
            this.DriveLoading = false
          })
          return
        }
        this.LoadDrive = true
        this.DriveLoading = false
      })
    },
    updateBoxHeight() {
      if (this.LoadDrive) {
        this.LoadDrive = false
        requestAnimationFrame(() => {
          this.LoadDrive = true
        })
      }
    }
  },
  computed: {
    getMainStyle() {
      if (this.LoadDrive === false || this.ProcessesAny) {
        return "display: flex; justify-content: center; align-items: center;width: 100%;"
      }
      return "width: 100%; height:" + this.$refs.device.parentElement.clientHeight + "px"
    },
    PidStyle() {
      if (this.Show) {
        return this.Show.Pid ? "height: calc(100% - 110px)" : "height: 40px"
      }

      return "height: 40px"
    },
    NameStyle() {
      if (this.Show) {
        return this.Show.Name ? "height: calc(100% - 80px)" : "height: 40px"
      }
      return "height: 40px"
    }
  },
  mounted() {
    attachMcpDeviceRoot(this);
    this.observer = new IntersectionObserver(entries => {
      if (this.ProcessesAny) {
        return
      }
      SetDeviceStopUpdate(!entries[0].isIntersecting)
    })
    GOOS().then((IsWindows) => {
      this.IsWindows = IsWindows
      window.addEventListener('resize', this.updateBoxHeight)
    })
  }
}
</script>

<template>
  <div :style="getMainStyle" ref="device">
    <div v-show="LoadDrive === false" style="width: 100%; height: 60px; ">
      <div class="mb-2 ml-4" style="display: flex; justify-content: center; align-items: center;">
        <el-radio-group v-model="LoadMode" :disabled="DriveLoading">
          <el-tooltip placement="top" v-if="IsWindows">
            <template #content>
              <div style="white-space: normal; line-height: 1.4;">
                使用 NetFilter 驱动
                <br><br>
                兼容性不是很好
                <br>
                部分系统可能无法加载
              </div>
            </template>
            <el-radio value="NFAPI" size="large">NFAPI</el-radio>
          </el-tooltip>

          <el-tooltip placement="top" v-if="IsWindows">
            <template #content>
              <div style="white-space: normal; line-height: 1.4;">
                使用 Proxifier 驱动
                <br><br>
                兼容性好
                <br>
                但是无法捕获UDP数据
              </div>
            </template>
            <el-radio value="Proxifier" size="large">Proxifier</el-radio>
          </el-tooltip>

          <el-tooltip placement="top">
            <template #content>
              <div style="white-space: normal; line-height: 1.4;">
                使用 Tun 驱动
                <br><br>
                兼容性好
                <br>
                但是无法捕获127.0.0.1数据
              </div>
            </template>
            <el-radio value="Tun" size="large">Tun</el-radio>
          </el-tooltip>

        </el-radio-group>
      </div>
      <div style="width: 100%; height: 100%; display: flex; justify-content: center; align-items: center;">
        <el-button v-show="DriveLoading === false" @click="loadDrive" style="width: 200px;">
          加载驱动
        </el-button>
        <el-button v-show="DriveLoading" loading style="width: 200px;">
          正在尝试加载驱动...
        </el-button>
      </div>
    </div>
    <div v-show="LoadDrive " class="ag-chart-tab ag-chart-format"
         style="overflow: hidden;scrollbar-width: none;height: 100%;">
      <div style="display: flex; justify-content: center; align-items: center;">
        <el-checkbox v-model="ProcessesAny" label="捕获任意进程" size="large"/>
        <el-button label="" size="small" style="margin-left: 20px" @click="AnyStopNetwork" v-if="IsWindows">对所有进程执行一次断网
        </el-button>
      </div>
      <div class="ag-chart-format-wrapper" v-show="ProcessesAny===false" style="height: 100%;">
        <div class="ag-chart-format-section" ref="_name" :style="NameStyle">
          <div class="ag-group ag-charts-format-top-level-group ag-group-item-alignment-center"
               data-ref="axisGroup">
            <div class="ag-group-title-bar ag-charts-format-top-level-group-title-bar ag-unselectable"
                 style="cursor: pointer" @click="Show.Name=!Show.Name">
              <span class="ag-group-title-bar-icon ag-charts-format-top-level-group-title-bar-icon">
                <span
                    :class="Show.Name?'ag-icon ag-icon ag-icon-tree-open':'ag-icon ag-icon ag-icon-tree-closed'"></span>
              </span>
              <span class="ag-group-title ag-charts-format-top-level-group-title" style="cursor: pointer">
              进程名称设置
              </span>
            </div>
            <div style="height: calc(100% - 30px);width: 100%" v-show="Show.Name">
              <Name ref="devName" style="height: 100%"></Name>
            </div>
          </div>
        </div>
        <div class="ag-chart-format-section" :style="PidStyle">
          <div class="ag-group ag-charts-format-top-level-group ag-group-item-alignment-center" style="height: 100%;"
               data-ref="axisGroup">
            <div class="ag-group-title-bar ag-charts-format-top-level-group-title-bar ag-unselectable"
                 style="cursor: pointer" @click="Show.Pid=!Show.Pid" ref="_pid">
              <span class="ag-group-title-bar-icon ag-charts-format-top-level-group-title-bar-icon">
                <span
                    :class="Show.Pid?'ag-icon ag-icon ag-icon-tree-open':'ag-icon ag-icon ag-icon-tree-closed'"></span>
              </span>
              <span class="ag-group-title ag-charts-format-top-level-group-title" style="cursor: pointer">
              进程PID设置
              </span>
            </div>
            <div style="height: 100%;width: 100%" v-show="Show.Pid">
              <Pid ref="devPid" style="height: 100%"></Pid>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<style>
.ag-chart-tabbed-menu-body::after {
  content: "";
  display: block;
  height: 0;
}
</style>