<template xmlns="http://www.w3.org/1999/html">
  <div class="top-component" :style="WindowStyleZIndex" @mousedown="WindowClick">
    <div class="ag-panel ag-default-panel ag-dialog ag-ltr ag-popup-child" tabindex="-1" role="dialog"
         aria-label="Range Chart"
         :style="WindowStyle" ref="Window">
      <div ref="eTitleBar" class="ag-panel-title-bar ag-default-panel-title-bar ag-unselectable"
           @mousedown="HeaderClick">
        <span ref="eTitle" class="ag-panel-title-bar-title ag-default-panel-title-bar-title"> {{ Title }}</span>
        <div ref="eTitleBarButtons" class="ag-panel-title-bar-buttons ag-default-panel-title-bar-buttons">
          <div class="ag-dialog-button ag-panel-title-bar-button" @click="SetSize">
            <span :class="UpdateIcon"></span>
          </div>
          <div class="ag-button ag-panel-title-bar-button" @click="Close">
            <span class="ag-icon ag-icon-cross ag-panel-title-bar-button-icon"></span>
          </div>
        </div>
      </div>
      <div ref="eContentWrapper" class="ag-panel-content-wrapper ag-default-panel-content-wrapper" style="height: 0px;">
        <div class="ag-chart ag-ltr" tabindex="-1">
          <div ref="eChartContainer" tabindex="-1" class="ag-chart-components-wrapper ag-chart-menu-visible">
            <div ref="eEmpty" class="ag-chart-empty-text ag-unselectable"
                 style="overflow-y: auto; overflow-x: hidden">
              <div :style="getSettingsRange">
                <el-container>
                  <el-container>
                    <el-main>
                      <el-collapse v-model="activeName" accordion>
                        <el-collapse-item name="常规设置">
                          <template #title>
                            <el-icon>
                              <setting/>
                            </el-icon>
                            &nbsp;
                            <span>常规设置</span>
                          </template>
                          <BasicSettings ref="Basic"/>
                        </el-collapse-item>
                        <el-collapse-item name="SSL证书">
                          <template #title>
                            <el-icon>
                              <Lock/>
                            </el-icon>
                            &nbsp;
                            <span>SSL证书</span>
                          </template>
                          <SSL ref="ssl"/>
                        </el-collapse-item>
                        <el-collapse-item name="强制走TCP">
                          <template #title>
                            <el-icon>
                              <Key/>
                            </el-icon>
                            &nbsp;
                            <span>强制走TCP</span>
                          </template>
                          <MustList ref="MustList"/>
                        </el-collapse-item>
                        <el-collapse-item name="上游网关">
                          <template #title>
                            <el-icon>
                              <Switch/>
                            </el-icon>
                            &nbsp;
                            <span>上游网关</span>
                          </template>
                          <agent ref="agent"/>
                        </el-collapse-item>
                        <el-collapse-item name="HOSTS设置">
                          <template #title>
                            <el-icon>
                              <Edit/>
                            </el-icon>
                            &nbsp;
                            <el-tooltip class="item" effect="dark" content="优先级高于替换规则,UDP请求无效"
                                        placement="top">
                              <span>HOSTS设置</span>
                            </el-tooltip>

                          </template>
                          <div style="width: 100%;height: 300px">
                            <List ref="hosts"/>
                          </div>
                        </el-collapse-item>
                        <el-collapse-item name="替换规则">
                          <template #title>
                            <el-icon>
                              <Grid/>
                            </el-icon>
                            &nbsp;
                            <el-tooltip class="item" effect="dark" content="优先级高于脚本编辑"
                                        placement="top">
                              <span> 替换规则</span>
                            </el-tooltip>
                          </template>
                          <div style="width: 100%;height: 300px">
                            <Replace ref="Replace"/>
                          </div>

                        </el-collapse-item>
                        <el-collapse-item name="脚本编辑">
                          <template #title>
                            <el-icon>
                              <EditPen/>
                            </el-icon>
                            &nbsp;
                            <span>脚本编辑</span>
                          </template>
                          <!--
                          <div style="width: 100%;height: 500px;position: relative">
                          </div>
                          -->
                          <SettingsJavaScriptEdit ref="GoJavaScriptEdit"/>
                        </el-collapse-item>
                        <el-collapse-item v-if="IsWindows" name="进程拦截">
                          <template #title>
                            <el-icon>
                              <Platform/>
                            </el-icon>
                            &nbsp;
                            <span>进程拦截</span>
                          </template>
                          <div v-show="LoadDrive" style="width: 100%;height: 300px;display: flex;">
                            <Proc/>
                          </div>
                          <div v-show="LoadDrive===false" style="width: 100%;height: 30px;display: flex;">
                            <div style="width: 100%;height: 100%">
                              <el-button v-show="DriveLoading===false" @click="loadDrive">加载驱动</el-button>
                              <el-button v-show="DriveLoading" loading>正在尝试加载驱动...</el-button>
                            </div>
                          </div>
                        </el-collapse-item>
                        <el-collapse-item name="请求证书">
                          <template #title>
                            <el-icon>
                              <Document/>
                            </el-icon>
                            &nbsp;
                            <span>请求证书</span>
                          </template>
                          <div style="width: 100%;height: 300px">
                            <RequestCertificate ref="RequestCertificate"/>
                          </div>

                        </el-collapse-item>
                        <el-collapse-item name="列表配色">
                          <template #title>
                            <el-icon>
                              <MagicStick/>
                            </el-icon>
                            &nbsp;
                            <span>列表配色</span>
                          </template>
                          <ListColorConfig ref="colorConfig" style="width: 100%;height: 456px"/>
                        </el-collapse-item>
                      </el-collapse>
                    </el-main>
                  </el-container>
                </el-container>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div class="ag-resizer-wrapper">
        <div ref="eTopLeftResizer" class="ag-resizer ag-resizer-topLeft" style="pointer-events: all;z-index: 1000"
             @mousedown="Resize($event,1)"></div>
        <div ref="eTopResizer" class="ag-resizer ag-resizer-top" style="pointer-events: all;z-index: 1000"
             @mousedown="Resize($event,2)"></div>
        <div ref="eTopRightResizer" class="ag-resizer ag-resizer-topRight" style="pointer-events: all;z-index: 1000"
             @mousedown="Resize($event,3)"></div>
        <div ref="eRightResizer" class="ag-resizer ag-resizer-right" style="pointer-events: all;z-index: 1000"
             @mousedown="Resize($event,4)"></div>
        <div ref="eBottomRightResizer" class="ag-resizer ag-resizer-bottomRight"
             style="pointer-events: all;z-index: 1000"
             @mousedown="Resize($event,5)"></div>
        <div ref="eBottomResizer" class="ag-resizer ag-resizer-bottom" style="pointer-events: all;z-index: 1000"
             @mousedown="Resize($event,6)"></div>
        <div ref="eBottomLeftResizer" class="ag-resizer ag-resizer-bottomLeft" style="pointer-events: all;z-index: 1000"
             @mousedown="Resize($event,7)"></div>
        <div ref="eLeftResizer" class="ag-resizer ag-resizer-left" style="pointer-events: all;z-index: 1000"
             @mousedown="Resize($event,8)"></div>
      </div>
    </div>
  </div>
</template>

<script>
import SettingsJavaScriptEdit from "./SettingsJavaScriptEdit.vue";
import {Base64DecodeStr, CallGoDo} from "./CallbackEventsOn.js";
import ListColorConfig from "./Settings/ListColorConfig/ListColorConfig.vue";
import BasicSettings from "./Settings/BasicSettings/BasicSettings.vue";
import SSL from "./Settings/SSL/SSL.vue";
import MustTCP from "./Settings/MustTCP/MustTCP.vue";
import Agent from "./Settings/Agent/agent.vue";
import List from "./Settings/HOSTS/List.vue";
import MustList from "./Settings/MustTCP/MustList.vue";
import Replace from "./Settings/ReplaceRules/Replace.vue";
import RequestCertificate from "./Settings/RequestCertificate/RequestCertificate.vue";
import Proc from "./Settings/process/proc.vue";
import {ElMessage} from "element-plus";

function pxToInt(E) {
  const eve = E + ""
  return parseInt(eve.replaceAll("px", ""))
}

const WindowName = "Settings"
const ClassMinName = "ag-icon ag-icon-minimize ag-panel-title-bar-button-icon"
const ClassMaxName = "ag-icon ag-icon-maximize ag-panel-title-bar-button-icon"
export default {
  props: ['show'],
  watch: {
    show(newValue) {
      if (newValue) {
        this.thisWindowWidth = 1000
        this.thisWindowHeight = 490
        this.MaxSize = false
        this.iconClass = ClassMaxName
        window.SetUILevel(WindowName)
        const h = document.documentElement.clientHeight
        const w = document.documentElement.clientWidth
        let _top = h - this.thisWindowHeight
        if (_top <= 60) {
          _top = 0
        } else {
          _top = _top / 2
        }
        let _left = w - this.thisWindowWidth
        if (_left < 0) {
          _left = 0
        } else {
          _left = _left / 2
        }
        this.setWindowSize(_top, _left, this.thisWindowWidth, this.thisWindowHeight)
      }
    },
    "activeName"(NewName, oldNmae) {
      const py = window.Theme.GOOS === "windows" ? 0 : -48
      let mTop = 0
      if (NewName !== "") {
        this.Title = "程序设置->" + NewName
      } else {
        this.Title = "程序设置"
      }
      if (NewName === "列表配色") {
        this.$refs.colorConfig.Load()
      } else if (NewName === "脚本编辑") {
        CallGoDo("获取脚本代码", null).then(res => {
          this.$refs.GoJavaScriptEdit.SetCode(Base64DecodeStr(res))
        })
        return
      }
      if (NewName === "常规设置") {
        mTop = 21
      } else if (NewName === "SSL证书") {
        mTop = 21 + 48
      } else if (NewName === "强制走TCP") {
        mTop = 21 + (2 * 48)
      } else if (NewName === "上游网关") {
        mTop = 21 + (3 * 48)
      } else if (NewName === "HOSTS设置") {
        mTop = 21 + (4 * 48)
      } else if (NewName === "替换规则") {
        mTop = 21 + (5 * 48)
      } else if (NewName === "脚本编辑") {
        mTop = 21 + (6 * 48)
      } else if (NewName === "进程拦截") {
        mTop = 21 + (7 * 48) + py
      } else if (NewName === "请求证书") {
        mTop = 21 + (8 * 48) + py
      } else if (NewName === "列表配色") {
        mTop = 21 + (9 * 48) + py
      }
      if (mTop === 0) {
        return;
      }
      let index = 0
      let lsat = -1;
      const SetScrollTop = () => {
        let MaxIndex = this.$refs.eEmpty.scrollHeight - this.$refs.eEmpty.clientHeight
        if (this.$refs.eEmpty.scrollTop !== mTop && MaxIndex !== lsat) {
          let mmTop = mTop
          if (mmTop > MaxIndex && mmTop !== 0) {
            mmTop = MaxIndex
          }
          //console.log("应该是", mTop, "之前是", this.$refs.eEmpty.scrollTop, "现在是", mmTop, this.$refs.eEmpty.clientHeight, this.$refs.eEmpty.scrollHeight, "Max", MaxIndex)
          this.$refs.eEmpty.scrollTop = mmTop
          lsat = MaxIndex
        }
        index++
        if (index > 50) {
          return
        }
        this.$nextTick(() => {
          setTimeout(() => {
            SetScrollTop()
          }, 20)
        })
      }
      this.$nextTick(() => {
        SetScrollTop()
      })
    }
  },
  components: {
    Proc,
    RequestCertificate,
    Replace,
    MustList,
    List, Agent, MustTCP, SSL, BasicSettings, ListColorConfig, SettingsJavaScriptEdit
  },
  computed: {
    IsWindows() {
      if (window.Theme) {
        return window.Theme.GOOS === "windows"
      }
      return false
    },
    WindowStyleZIndex() {
      return "z-index: " + window.UI.ZIndex[WindowName] + ";position: absolute; "
    },
    UpdateIcon() {
      return this.iconClass
    },
    svgStroke() {
      if (!this.theme) {
        return "#000"
      } else {
        return "#fff"
      }
    },
    getSettingsRange() {
      return "width: " + window.Size.Settings.Width + "px;height: " + window.Size.Settings.Height + "px;"
    }
  },
  data() {
    return {
      LoadDrive: false,
      DriveLoading: false,
      thisWindowWidth: 1000,
      thisWindowHeight: 490,
      get theme() {
        return window.Theme.IsDark
      },
      set theme(newValue) {
        window.Theme.IsDark = newValue
      },
      Title: "程序设置",
      activeName: "",
      iconClass: ClassMaxName,
      WindowStyle: "",
      MaxSize: false,
      RawSize: "",
      OptionWindowStyle: "",
      HeaderClickState: false,
      HeaderClickPosition: {
        left: 0,
        top: 0,
        width: 0,
        height: 0,
      },
      WindowResizeMode: 0,
      Show: {
        Settings: true,
        SSL: false,
        gateway: false,
        host: false,
        replace: false,
        script: false,
        process: false,
        certificate: false,
        color: false
      }
    }
  },
  mounted() {
    window.addEventListener('mouseup', this.handleMouseUp);
    window.addEventListener('mousemove', this.handleMouseMove);
    window.addEventListener('resize', this.handleResize);
    window.vm.Settings = this
    this.$nextTick(() => {
      // 获取元素的引用
      const elementRef = this.$refs.eEmpty;
      // 创建 ResizeObserver 实例并监听元素尺寸变化
      const resizeObserver = new ResizeObserver(entries => {
        for (const entry of entries) {
          const {width, height} = entry.contentRect;
          window.Size.Settings = {Width: width, Height: height}
        }
      });
      // 开始监听元素尺寸变化
      resizeObserver.observe(elementRef);
    })
  },
  beforeUnmount() {
    window.removeEventListener('mouseup', this.handleMouseUp);
    window.removeEventListener('mousemove', this.handleMouseMove);
    window.removeEventListener('resize', this.handleResize);
  },
  methods: {
    loadDrive() {
      this.DriveLoading = true
      CallGoDo("加载驱动", null).then(res => {
        if (res === true) {
          this.LoadDrive = true
        } else {
          ElMessage({
            message: "加载驱动失败:请检查是否有管理员权限？",
            type: 'error',
          })
        }
        this.DriveLoading = false
      })
    },
    deselect() {
      this.activeName = ""
    },
    handleResize(event) {
      if (this.MaxSize) {
        this.setWindowSize(0, 0, document.documentElement.clientWidth, document.documentElement.clientHeight - 60)
      }
    },
    Close() {
      window.UI[WindowName] = false
    },
    setWindowSize(top, left, width, height) {
      this.thisWindowWidth = width
      this.thisWindowHeight = height
      let t = "top: " + top + "px; left: " + left + "px; width: " + width + "px; max-width: " + width + "px; min-width: " + width + "px; height: " + height + "px; max-height: " + height + "px; min-height: " + height + "px;"
      t = t.replaceAll("pxpx;", "px;")
      this.WindowStyle = t
    },
    SetSize() {
      if (this.MaxSize === false) {
        this.iconClass = ClassMinName
        const top = pxToInt(this.$refs.Window.style.top)
        const left = pxToInt(this.$refs.Window.style.left)
        const width = pxToInt(this.$refs.Window.style.width)
        const height = pxToInt(this.$refs.Window.style.height)
        this.RawSize = "top: " + top + "px; left: " + left + "px; width: " + width + "px; max-width: " + width + "px; min-width: " + width + "px; height: " + height + "px; max-height: " + height + "px; min-height: " + height + "px;"
        this.setWindowSize(0, 0, document.documentElement.clientWidth, document.documentElement.clientHeight - 60)
        this.MaxSize = true
      } else {
        this.iconClass = ClassMaxName
        this.WindowStyle = this.RawSize
        this.MaxSize = false
      }
    },
    WindowClick(event) {
      window.SetUILevel(WindowName)
    },
    HeaderClick(event) {
      window.SetUILevel(WindowName)
      if (event.buttons !== 1) {
        return;
      }
      if (this.MaxSize) {
        this.HeaderClickState = false
        return
      }
      this.HeaderClickState = true
      this.HeaderClickPosition.left = pxToInt(event.clientX) - pxToInt(this.$refs.Window.style.left)
      this.HeaderClickPosition.top = pxToInt(event.clientY) - pxToInt(this.$refs.Window.style.top)
    },
    handleMouseUp(event) {
      this.HeaderClickState = false
      this.WindowResizeMode = 0
    },
    handleMouseMove(event) {
      {
        if (event.buttons !== 1) {
          return;
        }
        if (this.HeaderClickState) {
          let X = pxToInt(event.clientX) - this.HeaderClickPosition.left
          let Y = pxToInt(event.clientY) - this.HeaderClickPosition.top
          if (X < -(this.thisWindowWidth - 100)) {
            X = -(this.thisWindowWidth - 100)
          }
          if (X > document.documentElement.clientWidth - pxToInt(this.$refs.Window.style.width) + (this.thisWindowWidth - 100)) {
            X = document.documentElement.clientWidth - pxToInt(this.$refs.Window.style.width) + (this.thisWindowWidth - 100)
          }
          if (Y < 0) {
            Y = 0
          }
          if (Y > document.documentElement.clientHeight - pxToInt(this.$refs.Window.style.height) + (this.thisWindowHeight - 100)) {
            Y = document.documentElement.clientHeight - pxToInt(this.$refs.Window.style.height) + (this.thisWindowHeight - 100)
          }
          //this.$refs.Window.style.left = X + "px"
          //this.$refs.Window.style.top = Y + "px"
          const top = Y
          const left = X
          const width = pxToInt(this.$refs.Window.style.width)
          const height = pxToInt(this.$refs.Window.style.height)
          this.WindowStyle = "top: " + top + "px; left: " + left + "px; width: " + width + "px; max-width: " + width + "px; min-width: " + width + "px; height: " + height + "px; max-height: " + height + "px; min-height: " + height + "px;"

        } else if (this.WindowResizeMode === 1) {
          //调整左边和顶边
          let clientY = event.clientY - 1 - 30
          let clientX = event.clientX - 1
          if (clientY < 0) {
            clientY = pxToInt(this.$refs.Window.style.top)
          }
          if (clientX < 0) {
            clientX = pxToInt(this.$refs.Window.style.left)
          }
          let height = this.HeaderClickPosition.top - clientY + this.HeaderClickPosition.height
          let width = this.HeaderClickPosition.left - clientX + this.HeaderClickPosition.width
          if (width < 766) {
            width = 766
            clientX = pxToInt(this.$refs.Window.style.left)
          }
          if (height < 30) {
            height = 30
            clientY = pxToInt(this.$refs.Window.style.top)
          }
          this.setWindowSize(clientY, clientX, width, height)
        } else if (this.WindowResizeMode === 2) {
          //调整顶边
          if (event.clientY < 0) {
            return;
          }
          const a1 = this.HeaderClickPosition.top - event.clientY
          let h = this.HeaderClickPosition.height + a1
          let g = event.clientY - 1 - 30
          if (g < 0) {
            h = h + g
            g = 0
          }
          if (h < 0) {
            h = 0
            g = this.HeaderClickPosition.top2
          } else {
            this.HeaderClickPosition.top2 = g
          }

          this.setWindowSize(g, this.$refs.Window.style.left, this.$refs.Window.style.width, h)
        } else if (this.WindowResizeMode === 3) {
          //调整右边和顶边
          let clientX = event.clientX - 1
          if (clientX < 0) {
            clientX = 0
          }
          let width = clientX - this.HeaderClickPosition.left + this.HeaderClickPosition.width
          if (width < 766) {
            width = 766
          }
          if (width + pxToInt(this.$refs.Window.style.left) > document.documentElement.clientWidth) {
            width = document.documentElement.clientWidth - pxToInt(this.$refs.Window.style.left)
          }
          let clientY = event.clientY - 1 - 30
          if (clientY < 0) {
            clientY = 0
          }
          let a1 = this.HeaderClickPosition.top - clientY
          if (width < 766) {
            width = 766
            clientY = pxToInt(this.$refs.Window.style.left)
          }
          let height = this.HeaderClickPosition.height + a1
          if (height < 30) {
            height = 30
            clientY = this.HeaderClickPosition.top2
          } else {
            this.HeaderClickPosition.top2 = clientY
          }
          this.setWindowSize(clientY, pxToInt(this.$refs.Window.style.left), width + 2, height)

        } else if (this.WindowResizeMode === 4) {
          //调整右边
          let clientX = event.clientX - 1
          if (clientX < 0) {
            clientX = 0
          }
          let width = clientX - this.HeaderClickPosition.left + this.HeaderClickPosition.width
          if (width < 766) {
            width = 766
          }
          if (width + pxToInt(this.$refs.Window.style.left) > document.documentElement.clientWidth) {
            width = document.documentElement.clientWidth - pxToInt(this.$refs.Window.style.left)
          }
          this.setWindowSize(pxToInt(this.$refs.Window.style.top), pxToInt(this.$refs.Window.style.left), width + 2, this.HeaderClickPosition.height)
        } else if (this.WindowResizeMode === 5) {
          //调整右边和底边
          let clientX = event.clientX - 1
          if (clientX < 0) {
            clientX = 0
          }
          let width = clientX - this.HeaderClickPosition.left + this.HeaderClickPosition.width
          if (width < 766) {
            width = 766
          }
          if (width + pxToInt(this.$refs.Window.style.left) > document.documentElement.clientWidth) {
            width = document.documentElement.clientWidth - pxToInt(this.$refs.Window.style.left)
          }
          if (width < 766) {
            width = 766
          }
          let a1 = event.clientY - pxToInt(this.$refs.Window.style.top) - 30
          if (a1 + pxToInt(this.$refs.Window.style.top) > document.documentElement.clientHeight) {
            a1 = document.documentElement.clientHeight - pxToInt(this.$refs.Window.style.top)
          }
          if (a1 < 30) {
            a1 = 30
          }
          this.setWindowSize(pxToInt(this.$refs.Window.style.top), pxToInt(this.$refs.Window.style.left), width + 2, a1 + 2)

        } else if (this.WindowResizeMode === 6) {
          //调整底边
          let a1 = event.clientY - pxToInt(this.$refs.Window.style.top) + 2 - 30
          if (a1 + pxToInt(this.$refs.Window.style.top) > document.documentElement.clientHeight) {
            a1 = document.documentElement.clientHeight - pxToInt(this.$refs.Window.style.top)
          }
          if (a1 < 30) {
            a1 = 30
          }
          this.setWindowSize(pxToInt(this.$refs.Window.style.top), pxToInt(this.$refs.Window.style.left), pxToInt(this.$refs.Window.style.width), a1)
        } else if (this.WindowResizeMode === 7) {
          //调整左边和底边
          let a1 = event.clientY - pxToInt(this.$refs.Window.style.top) + 2 - 30
          if (a1 + pxToInt(this.$refs.Window.style.top) > document.documentElement.clientHeight) {
            a1 = document.documentElement.clientHeight - pxToInt(this.$refs.Window.style.top)
          }
          if (a1 < 30) {
            a1 = 30
          }
          let clientX = event.clientX - 1
          if (clientX < 0) {
            clientX = 0
          }
          let width = this.HeaderClickPosition.left - clientX + this.HeaderClickPosition.width
          if (width < 766) {
            width = 766
            clientX = pxToInt(this.$refs.Window.style.left)
          }
          if (a1 < 30) {
            a1 = 30
          }
          this.setWindowSize(pxToInt(this.$refs.Window.style.top), clientX, width, a1)
        } else if (this.WindowResizeMode === 8) {
          //调整左边
          let clientX = event.clientX - 1
          if (clientX < 0) {
            clientX = 0
          }
          let width = this.HeaderClickPosition.left - clientX + this.HeaderClickPosition.width
          if (width < 766) {
            width = 766
            clientX = pxToInt(this.$refs.Window.style.left)
          }
          this.setWindowSize(pxToInt(this.$refs.Window.style.top), clientX, width, this.HeaderClickPosition.height)
        }
      }
    },
    Resize(event, mode) {
      if (event.buttons !== 1) {
        return;
      }
      this.WindowResizeMode = mode
      this.HeaderClickPosition.left = pxToInt(event.clientX)
      this.HeaderClickPosition.top = pxToInt(event.clientY)
      this.HeaderClickPosition.width = pxToInt(this.$refs.Window.style.width)
      this.HeaderClickPosition.height = pxToInt(this.$refs.Window.style.height)
    },
  },

}
</script>
<style scoped>
.custom-menu-item {
  height: 50px; /* 调整组件之间的下边距 */
}
</style>
