<script>
import RequestTabs from './Request/Request_Tabs.vue';
import ResponseTabs from './Response/Response_Tabs.vue';
import {CallGoDo, SunnyErrorReplaceAll} from "./CallbackEventsOn.js";

const ClassMinName = "ag-icon ag-icon-minimize ag-panel-title-bar-button-icon"
const ClassMaxName = "ag-icon ag-icon-maximize ag-panel-title-bar-button-icon"

export default {
  components: {
    RequestTabs,
    ResponseTabs,
  },

  data() {
    return {
      TitleMsg: "<span>此请求正在发送中</span>",
      NoRequest: true,
      StateCode: "",
      CodeState: "",
      NoData: true,
      isChangeSize: false,
      clientY: 0,
      RawSize: "",
      iconClass: ClassMaxName,
      Method: "",
      Show: {
        Request: true,
        Response: true,
        Height: {
          Request: 0,
          Response: 0,
        }
      }
      ,
      DisplayHTTPHeaderWebSocket: false,
      DisplayHTTPHeader: true,
      DisplayTCPResponse: false,
    };
  },
  computed: {
    GetTitleMsg() {
      return this.TitleMsg
    }
    ,
    GetTCPResponseHeader() {
      if (this.DisplayTCPResponse) {
        return "flex: 1 1 auto; padding-top: 0px;"
      }
      return "flex: 1 1 auto;"
    },
    UpdateIcon() {
      return this.iconClass
    },
    GetRequestWay() {
      window.mv = window.vm.List.agSelectedLine
      if (window.vm.List.agSelectedLine != null) {
        let way = window.vm.List.agSelectedLine.data["方式"].toUpperCase()
        let wayType = window.vm.List.agSelectedLine.data["响应类型"].toUpperCase()
        if (way.indexOf("TCP") !== -1 || way.indexOf("UDP") !== -1) {
          this.DisplayHTTPHeader = false
          this.DisplayHTTPHeaderWebSocket = false
          this.DisplayTCPResponse = true
          this.$refs.Request.$refs.WebSocket.SetColumnsMode(false)
        } else if (wayType.indexOf("WEBSOCKET") !== -1) {
          this.DisplayHTTPHeader = true
          this.DisplayHTTPHeaderWebSocket = true
          this.$refs.Request.$refs.WebSocket.SetColumnsMode(true)
        } else {
          this.DisplayHTTPHeader = true
          this.DisplayHTTPHeaderWebSocket = false
          this.DisplayTCPResponse = false
        }
        return {
          DisplayHTTPHeader: this.DisplayHTTPHeader,
          DisplayHTTPHeaderWebSocket: this.DisplayHTTPHeaderWebSocket,
          DisplayTCPResponse: this.DisplayTCPResponse
        }
      } else {
        this.NoRequest = true
        return {
          DisplayHTTPHeader: false,
          DisplayHTTPHeaderWebSocket: false,
          DisplayTCPResponse: false
        }
      }
    },
    GetTheology() {
      this.NoRequest = true
      if (window.vm.List.agSelectedLine != null) {
        let Theology = parseInt(window.vm.List.agSelectedLine.data.Theology)
        const Break = window.vm.List.agSelectedLine.data['断点模式']
        if (Theology === -1) {
          return Theology
        }
        this.NoRequest = false
        CallGoDo("HTTP请求获取", {Theology: Theology}).then(response => {
          this.Method = response.Method
          this.StateCode = response.Response.StateCode
          this.CodeState = "success"
          if (this.StateCode === 403 || this.StateCode === 404) {
            this.CodeState = "danger"
          } else if (this.StateCode !== 200 && this.StateCode !== 101) {
            this.CodeState = "warning"
          }
          window.Theology = Theology
          this.$refs.Request.UpdateData(response, Break)
          if (this.StateCode === 0 || response.Response.Header === null) {
            this.NoData = true
            if (Break === 1) {
              this.TitleMsg = "<span>您正在编辑此请求</span>"
            } else {
              this.TitleMsg = "<span>此请求正在发送中</span>"
            }
          } else {
            this.NoData = false
          }
          if (response.Response.Error) {
            window.vm.List.agSelectedLine.data["断点模式"] = 0
            window.vm.Tabs.Request.BreakResponse = false
            window.vm.Tabs.Request.Breakpoint = false
            let Body = SunnyErrorReplaceAll(response.Response.Body)
            window.vm.Tabs.ToolPanel.NoData = true
            window.vm.Tabs.ToolPanel.TitleMsg = `<span style="color: red">${Body}</span>`
            return
          }
          this.$refs.Response.UpdateData(response.Response, Break)
        })
        return Theology
      } else {
        this.NoRequest = true
        return -1
      }
    }
    ,
    TopStyle() {
      if (this.DisplayHTTPHeader) {
        return "height: 50%; flex: 0 0 auto;padding-top: calc(var(--ag-grid-size) * 2)"
      }
      return "height: 50%; flex: 0 0 auto;"
    },
    TcpClass() {
      if (this.DisplayHTTPHeader) {
        return "ag-unselectable ag-column-drop-vertical"
      }
      return "ag-unselectable"
    }
  },
  methods: {
    getElementSize(Element) {
      const rect = Element.getBoundingClientRect();
      const left = rect.left;
      const top = rect.top;
      const width = rect.width;
      const height = rect.height;
      return {left: left, top: top, width: width, height: height}
    },
    handleMouseDown(event) {
      this.isChangeSize = true;
    },
    getTools() {
      const agTools = document.getElementsByClassName("ag-tool-panel-wrapper")
      for (let i = 0; i < agTools.length; i++) {
        if (agTools[i].className.indexOf("ag-hidden") === -1) {
          return agTools[i]
        }
      }
      return null
    },
    SetMaximize(mode) {
      if (window.ToolsMaximize) {
        this.iconClass = ClassMaxName
        const Tools = this.getTools()
        if (Tools !== null) {
          Tools.style.width = this.RawSize
        }
        window.ToolsMaximize = false
        this.$refs.splitA1.style.height = this.Show.Height.Request
      } else {
        this.iconClass = ClassMinName
        window.ToolsMaximize = true
        const Tools = this.getTools()
        if (Tools !== null) {
          this.RawSize = Tools.style.width
          Tools.style.width = window.innerWidth + "px"
        }
        this.Show.Height.Request = this.$refs.splitA1.style.height
        this.$refs.splitA1.style.height = "100%"
      }
      if (mode === 1) {
        this.Show.Request = true
        this.Show.Response = !window.ToolsMaximize
      } else {
        this.Show.Request = !window.ToolsMaximize
        this.Show.Response = true
      }
    },
    handleMouseMove(event) {
      if (this.isChangeSize) {
        this.clientY = event.clientY;
        if (this.clientY + 5 < 60) {
          return
        }
        if (this.clientY + 5 > document.documentElement.clientHeight - 60) {
          this.$refs.splitA1.style.height = (document.documentElement.clientHeight - 90) + 'px';
          return
        }
        this.$refs.splitA1.style.height = (this.clientY + 5 - 30) + 'px';
      }
    }
  },
  mounted() {
    document.addEventListener('mouseup', (event) => {
      if (this.isChangeSize) {
        this.isChangeSize = false;
      }
    });
    document.addEventListener('mousemove', this.handleMouseMove);
    window.vm.Tabs.ToolPanel = this
  },
};
</script>

<template>
  <div class="ag-column-panel" style="position: relative">
    <div ref="splitA1" v-show="Show.Request&&NoRequest===false" class="ag-column-select ag-column-panel-column-select"
         :style="TopStyle">
      <div v-if="DisplayHTTPHeader" aria-hidden="true"
           class="ag-column-drop-title-bar ag-column-drop-vertical-title-bar">
        <span class="ag-icon ag-icon-group ag-column-drop-icon ag-column-drop-vertical-icon" unselectable="on"
              role="presentation"></span><span
          class="ag-column-drop-title ag-column-drop-vertical-title">请求内容</span>
        &nbsp;&nbsp;
        <el-tag class="ml-2" type="success" style="right: 30px;border-radius: 15px;">{{ Method }}</el-tag>
        <span :class="UpdateIcon" unselectable="on" role="presentation"
              style="right: 0;position:absolute;cursor:pointer" @click="SetMaximize(1)"></span>
      </div>
      <RequestTabs ref="Request" :Theology="GetTheology" :RequestWay="GetRequestWay"/>
      <div ref="splitA2" v-show="Show.Response" @mousedown="handleMouseDown" class="ag-resizer ag-resizer-bottom"
           style="pointer-events: all;z-index: 1000"></div>

    </div>
    <div ref="splitA3" v-show="Show.Response&& (NoData===false || DisplayTCPResponse)&&NoRequest===false"
         :class="TcpClass" role="presentation"
         :style="GetTCPResponseHeader">
      <div v-show="DisplayHTTPHeader && DisplayTCPResponse===false" aria-hidden="true"
           class="ag-column-drop-title-bar ag-column-drop-vertical-title-bar">
        <span class="ag-icon ag-icon-aggregation ag-column-drop-icon ag-column-drop-vertical-icon" unselectable="on"
              role="presentation"></span>
        <span class="ag-column-drop-title ag-column-drop-vertical-title">响应内容</span>
        &nbsp;&nbsp;
        <el-tag class="ml-2" :type="CodeState" style="right: 30px;border-radius: 15px;">{{ StateCode }}</el-tag>
        <span :class="UpdateIcon" unselectable="on" role="presentation"
              style="right: 0px;position:absolute;cursor:pointer" @click="SetMaximize(2)"></span>
      </div>
      <ResponseTabs ref="Response" :Theology="GetTheology" :RequestWay="GetRequestWay"/>
    </div>
    <div v-show="NoRequest===false&& NoData&&DisplayTCPResponse===false"
         style="width: 100%;height: 100%;display: flex;align-items: center;justify-content: center;">
      <div v-html="GetTitleMsg"></div>
    </div>
    <div v-show="NoRequest"
         style="width: 100%;height: 100%;display: flex;align-items: center;justify-content: center;position: absolute">
      <div>
        <span>请选择一个会话</span>
      </div>
    </div>
  </div>

</template>