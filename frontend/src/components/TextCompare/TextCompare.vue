<template>
  <div class="top-component" :style="WindowStyleZIndex" @mousedown="WindowClick">
    <div class="ag-panel ag-default-panel ag-dialog ag-ltr ag-popup-child" tabindex="-1" role="dialog"
         aria-label="Range Chart"
         :style="WindowStyle" ref="Window">
      <div ref="eTitleBar" class="ag-panel-title-bar ag-default-panel-title-bar ag-unselectable"
           @mousedown="HeaderClick">
        <span ref="eTitle" class="ag-panel-title-bar-title ag-default-panel-title-bar-title">文本对比</span>
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
            <div ref="eEmpty" class="ag-chart-empty-text ag-unselectable">
              <DiffEditor style="width: 100%;height: 100%"></DiffEditor>
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
import DiffEditor from "./diffeditor.vue";

function pxToInt(E) {
  const eve = E + ""
  return parseInt(eve.replaceAll("px", ""))
}

const WindowName = "TextCompare"
const ClassMinName = "ag-icon ag-icon-minimize ag-panel-title-bar-button-icon"
const ClassMaxName = "ag-icon ag-icon-maximize ag-panel-title-bar-button-icon"
export default {
  components: {DiffEditor},
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
        if (_top < 0) {
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
    }
  },
  computed: {
    WindowStyleZIndex() {
      return "z-index: " + window.UI.ZIndex[WindowName] + ";position: absolute;"
    },
    UpdateIcon() {
      return this.iconClass
    },
  },
  data() {
    return {
      thisWindowWidth: 950,
      thisWindowHeight: 470,
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
      Text1: "测试文本1 text\r\n测试文本",
      Text2: "测试文本2 text\r\n测试文本\r\n测试文本",
      OPTIONS: {
        automaticLayout: true,
        formatOnType: true,
        formatOnPaste: true,
        originalEditable: true,
        glyphMargin: false,
        get theme() {
          if (window.Theme.IsDark) {
            return 'vs-dark'
          }
          return 'vs'
        },
        set theme(newValue) {
          if (newValue === "vs") {
            window.Theme.IsDark = false
            return
          }
          window.Theme.IsDark = true
        }
      },
      diffEditorRef: null
    }
  },
  mounted() {
    window.addEventListener('mouseup', this.handleMouseUp);
    window.addEventListener('mousemove', this.handleMouseMove);
    window.addEventListener('resize', this.handleResize);
  },
  beforeUnmount() {
    window.removeEventListener('mouseup', this.handleMouseUp);
    window.removeEventListener('mousemove', this.handleMouseMove);
    window.removeEventListener('resize', this.handleResize);
  },
  methods: {
    handleResize(event) {
      if (this.MaxSize) {
        this.setWindowSize(30, 0, document.documentElement.clientWidth, document.documentElement.clientHeight - 60)
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
        this.setWindowSize(30, 0, document.documentElement.clientWidth, document.documentElement.clientHeight - 60)
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
          if (Y < 30) {
            Y = 30
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
          let clientY = event.clientY - 1
          let clientX = event.clientX - 1
          if (clientY < 0) {
            clientY = pxToInt(this.$refs.Window.style.top)
          }
          if (clientX < 0) {
            clientX = pxToInt(this.$refs.Window.style.left)
          }
          let height = this.HeaderClickPosition.top - clientY + this.HeaderClickPosition.height
          let width = this.HeaderClickPosition.left - clientX + this.HeaderClickPosition.width
          if (width < 220) {
            width = 220
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
          let g = event.clientY - 1
          if (h < 30) {
            h = 30
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
          if (width < 220) {
            width = 220
          }
          if (width + pxToInt(this.$refs.Window.style.left) > document.documentElement.clientWidth) {
            width = document.documentElement.clientWidth - pxToInt(this.$refs.Window.style.left)
          }
          let clientY = event.clientY - 1
          if (clientY < 0) {
            clientY = 0
          }
          let a1 = this.HeaderClickPosition.top - clientY
          if (width < 220) {
            width = 220
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
          if (width < 220) {
            width = 220
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
          if (width < 220) {
            width = 220
          }
          if (width + pxToInt(this.$refs.Window.style.left) > document.documentElement.clientWidth) {
            width = document.documentElement.clientWidth - pxToInt(this.$refs.Window.style.left)
          }
          if (width < 220) {
            width = 220
          }
          let a1 = event.clientY - pxToInt(this.$refs.Window.style.top)
          if (a1 + pxToInt(this.$refs.Window.style.top) > document.documentElement.clientHeight) {
            a1 = document.documentElement.clientHeight - pxToInt(this.$refs.Window.style.top)
          }
          if (a1 < 30) {
            a1 = 30
          }
          this.setWindowSize(pxToInt(this.$refs.Window.style.top), pxToInt(this.$refs.Window.style.left), width + 2, a1 + 2)

        } else if (this.WindowResizeMode === 6) {
          //调整底边
          let a1 = event.clientY - pxToInt(this.$refs.Window.style.top) + 2
          if (a1 + pxToInt(this.$refs.Window.style.top) > document.documentElement.clientHeight) {
            a1 = document.documentElement.clientHeight - pxToInt(this.$refs.Window.style.top)
          }
          if (a1 < 30) {
            a1 = 30
          }
          this.setWindowSize(pxToInt(this.$refs.Window.style.top), pxToInt(this.$refs.Window.style.left), pxToInt(this.$refs.Window.style.width), a1)
        } else if (this.WindowResizeMode === 7) {
          //调整左边和底边
          let a1 = event.clientY - pxToInt(this.$refs.Window.style.top) + 2
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
          if (width < 220) {
            width = 220
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
          if (width < 220) {
            width = 220
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
    handleMount(diffEditor) {
      this.diffEditorRef = diffEditor;
    }
  }
}
</script>
<style scoped>
.top-component {
  position: absolute;
  top: 0;
  left: 0;
}
</style>