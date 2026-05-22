<template>
  <div ref="splitA0" class="ag-column-panel" style="position: relative;width: 100%;height: 100%;display: block;">
    <div ref="splitA1" v-show="isShowRequest" class="ag-column-select ag-column-panel-column-select"
         :style="RequestStyle">
      <div aria-hidden="true"
           class="ag-column-drop-title-bar ag-column-drop-vertical-title-bar"
           style="    height: 30px;    top: -4px;    position: relative;zIndex: 0">
        <span class="ag-icon ag-icon-group ag-column-drop-icon ag-column-drop-vertical-icon"
        ></span><span
          class="ag-column-drop-title ag-column-drop-vertical-title">请求内容</span>
        <el-tag type="success" style="height: 20px;border-radius: 10px;margin-left: 10px">{{ Method }}</el-tag>
        <span :class="UpdateIcon" unselectable="on"
              style="right: 10px;position:absolute;cursor:pointer" @click="SetMaximize(1)"></span>
      </div>
      <div ref="RequestTabs" class="ag-tabs ag-chart-tabbed-menu ag-focus-managed"
           style="height: auto;zIndex: 0;text-align:center">
        <div ref="RequestTabsList" class="ag-tabs-header ag-chart-tabbed-menu-header"
             style="height: 30px;top:50px;display: ruby;">
           <span
               v-for="(tab, index) in Tabs"
               :key="index"
               @click="eHeaderClick(tab)"
               :class="tab.class"
               role="tab"
               v-show="tab.Show?.Text ?? tab.Show"
           >
    {{ tab.name }}
  </span>
        </div>
      </div>
      <div class="ag-tabs-body ag-chart-tabbed-menu-body" style="">
        <div :style="{  width: '100%',  height: '100%',  position: 'absolute'}" class="RequestDIV">
          <div v-show="Tabs[0].visible" :style="GetBreakpointStateStyle">
            <JavaScriptEdit ref="RequestRaw" :glyphMargin="false" Name="Request"/>
          </div>
          <div v-show="Tabs[1].visible" :style="GetBreakpointStateStyle">
            <Table ref="RequestHeader" :style="'height: 100%;width: 100%'"
                   :glyphMargin="false"
                   Name="Request"/>
          </div>
          <div v-show="Tabs[2].visible" :style="GetBreakpointStateStyle">
            <Table ref="RequestUrlArgs" style="height: 100%;width: 100%" :glyphMargin="false"
                   Name="Request"/>
          </div>
          <div v-show="Tabs[3].visible" :style="GetBreakpointStateStyle">
            <div style="width: 100%; height: 100%; justify-content: flex-start; align-items: center;">
              <div v-show="Tabs[3].Show.BodyArgs"
                   style="width: 100%; height: 30px; display: flex; justify-content: flex-start; align-items: center;">
                <div v-show="Tabs[3].Show.Index===0">
                  <el-button size="small" disabled>原始内容</el-button>
                  <el-button size="small" @click="RequestRawTextTypeClick">表单数据</el-button>
                </div>
                <div v-show="Tabs[3].Show.Index===1">
                  <el-button size="small" @click="RequestRawTextTypeClick">原始内容</el-button>
                  <el-button size="small" disabled>表单数据</el-button>
                </div>
              </div>
              <div v-show="Tabs[3].Show.Index===0"
                   :style="{width: '100%',height: (Tabs[3].Show.BodyArgs?'calc(100% - 30px)':'100%')}">
                <JavaScriptEdit ref="RequestRawText" :glyphMargin="false" Name="RequestRawText"/>
              </div>
              <div v-show="Tabs[3].Show.Index===1" style="width: 100%;height: calc(100% - 30px)">
                <Table ref="RequestBodyTable" style="height: 100%;width: 100%"
                       :glyphMargin="false"
                       Name="RequestBodyTable"/>
              </div>
            </div>
          </div>
          <div v-show="Tabs[4].visible" :style="GetBreakpointStateStyle">
            <HexView ref="RequestHex" style="height: 100%;width: 100%" Name="RequestHex"/>
          </div>

          <div v-show="Tabs[5].visible" :style="GetBreakpointStateStyle">
            <Table ref="RequestCookie" style="height: 100%;width: 100%" :glyphMargin="false" Name="RequestCookie"/>
          </div>
          <div v-show="Tabs[6].visible " :style="GetBreakpointStateStyle">
            <JSon ref="RequestJson" :height="'100%'" :width="'100%'"/>
          </div>

          <div v-show="Tabs[7].visible" :style="GetBreakpointStateStyle">
            <WebSocketView ref="RequestWebSocket" :height="'100%'" :width="'100%'" Name="Websocket"/>
          </div>

          <div v-show="Tabs[8].visible" :style="GetBreakpointStateStyle">
            <IMGView ref="RequestImg" style="height: 100%;width: 100%" :isRequest="true"/>
          </div>

          <div v-show="isBreak"
               style="width: 100%; height: 30px; display: flex; justify-content: center; align-items: center;">
            <div
                style="margin: 0 10px; background-color: #f64141; height: 28px; line-height: 28px; width: 110px; border-radius: 5px;">
              &nbsp;断点命中,然后->&nbsp;
            </div>
            <div v-show="BreakResponse===false">
              <el-button type="warning" style="height: 28px;width: 100px" @click="BreakClick(2)">中断响应</el-button>&nbsp;&nbsp;&nbsp;
            </div>
            <el-button type="success" style="height: 28px;width: 100px" @click="BreakClick(0)">运行到结束</el-button>
          </div>
        </div>

      </div>
      <div ref="split" @mousedown="handleMouseDown" class="ag-resizer ag-resizer-bottom"
           style="pointer-events: all;z-index: 4"></div>
    </div>
    <div ref="splitA3" v-show="isShowResponse&&!noResponse" :style="RangeStyle.splitA3">
      <div v-show="!isWsResponse" ref="ResponseTabs1"
           class="ag-column-drop-title-bar ag-column-drop-vertical-title-bar">
        <span class="ag-icon ag-icon-aggregation ag-column-drop-icon ag-column-drop-vertical-icon"></span>
        <span class="ag-column-drop-title ag-column-drop-vertical-title">响应内容</span>
        <el-tag class="ml-2" :type="ResponseCodeStateStyle" style="right: 30px;border-radius: 15px;margin-left: 10px">
          {{ ResponseCode }}
        </el-tag>
        <span :class="UpdateIcon" unselectable="on" style="right: 10px;position:absolute;cursor:pointer"
              @click="SetMaximize(2)"></span>
      </div>
      <div v-show="isWsResponse" ref="ResponseTabs2"
           class="ag-column-drop-title-bar ag-column-drop-vertical-title-bar">
        <div style="width: calc(100% - 30px);display: flex">
            <span
                v-for="(tab, index) in WebSocketTabs"
                :key="index"
                v-show="tab.Show"
                @click="eWebsocketClick(tab)"
                :class="tab.class"
                role="tab"
                style="width: auto"
            >
            {{ tab.name }}
            </span>
        </div>
        <span :class="UpdateIcon" unselectable="on" style="right: 10px;position:absolute;cursor:pointer"
              @click="SetMaximize(2)"></span>
      </div>
      <div v-show="!isWsResponse" ref="ResponseTabs3" class="ag-tabs ag-chart-tabbed-menu ag-focus-managed"
           style="height: auto;zIndex: 0;text-align:center">
        <div ref="ResponseTabsList" class="ag-tabs-header ag-chart-tabbed-menu-header"
             style="height: 30px;top:50px;display: ruby">
          <span
              v-for="(tab, index) in ResponseTabs"
              :key="index"
              v-show="tab.Show"
              @click="eResponseClick(tab)"
              :class="tab.class"
              role="tab"
              style="width: auto"
          >
    {{ tab.name }}
  </span>
        </div>
      </div>
      <div v-show="!isWsResponse" ref="ResponseTabs4" class="ag-tabs-body ag-chart-tabbed-menu-body"
           :style="RangeStyle.ResponseTabs4">
        <div :style="{  width: '100%',  height: '100%',  position: 'absolute'}" class="ResponseDIV">
          <div v-show="ResponseTabs[0].visible" style="width: 100%;height: 100%">
            <JavaScriptEdit ref="ResponseRaw" :glyphMargin="false" Name="Response"/>
          </div>
          <div v-show="ResponseTabs[1].visible" style="width: 100%;height: 100%">
            <Table ref="ResponseHeader" :style="'height: 100%;width: 100%'"
                   :glyphMargin="false"
                   Name="Response"/>
          </div>
          <div v-show="ResponseTabs[2].visible" style="width: 100%;height: 100%">
            <JavaScriptEdit ref="ResponseText" :glyphMargin="false" Name="ResponseText"/>
          </div>
          <div v-show="ResponseTabs[3].visible" style="width: 100%;height: 100%">
            <IMGView ref="ResponseImg" style="height: 100%;width: 100%"
                     Name="ResponseImg" :isRequest="false"/>
          </div>
          <div v-show="ResponseTabs[4].visible" style="width: 100%;height: 100%">
            <div class="iframe-container" style="width: 100%;height: 100%">
              <iframe ref="iframe" style="width: 100%; height: 100%"
                      :srcdoc="`${ResponseHTMLValue}`"></iframe>

            </div>
          </div>

          <div v-show="ResponseTabs[5].visible" style="width: 100%;height: 100%">
            <HexView ref="ResponseHEX" :glyphMargin="false" Name="ResponseHEX"/>
          </div>

          <div v-show="ResponseTabs[6].visible" style="width: 100%;height: 100%">
            <Table ref="ResponseCookie" style="height: 100%;width: 100%" :glyphMargin="false"
                   Name="ResponseCookie"/>
          </div>
          <div v-show="ResponseTabs[7].visible " style="width: 100%;height: 100%">
            <JSon ref="ResponseJson" :height="'100%'" :width="'100%'"/>
          </div>

        </div>
      </div>
      <div v-show="isWsResponse" class="ag-tabs-body ag-chart-tabbed-menu-body"
           style="width: 100%;height: 100%;">
        <div :style="{  width: '100%',  height: 'calc(100% - 33px)',  position: 'absolute'}">
          <div v-show="WebSocketTabs[0].visible" style="width: 100%;height: 100%">
            <JavaScriptEdit ref="WebSocketText" :glyphMargin="false" Name="WebSocketText"/>
          </div>
          <div v-show="WebSocketTabs[1].visible" style="width: 100%;height: 100%">
            <HexView ref="WebSocketHEX" :glyphMargin="false" Name="WebSocketHEX"/>
          </div>
          <div v-show="WebSocketTabs[2].visible" style="width: 100%;height: 100%">
            <JSon ref="WebSocketJson" :height="'100%'" :width="'100%'"/>
          </div>
          <div v-show="WebSocketTabs[3].visible" style="width: 100%;height: 100%">
            <WebsocketActive ref="WebSocketActive" style="width: 100%;height: 100%"/>
          </div>
        </div>
      </div>
    </div>
    <div ref="splitA2" v-show="noResponse" style="width: 100%;height: 100%;">
      <div v-show="Errors.isError===true" class="ag-column-drop-title-bar ag-column-drop-vertical-title-bar"
           style="width: 100%;height: 100%;display: flex;  justify-content: center; align-items: center; ">
        <el-result
            icon="error"
            title="错误"
            :sub-title="Errors.value"
            style="width: 100%;height: 100%;display: flex;  justify-content: center; align-items: center; "
        >
        </el-result>
      </div>
      <div v-show="!Errors.isError" class="ag-column-drop-title-bar ag-column-drop-vertical-title-bar"
           style="width: 100%;height: 100%;display: flex;  justify-content: center; align-items: center; ">
        {{ Errors.value }}
      </div>
    </div>
  </div>

</template>
<script>
import {
  Config_GOOS_IsWindows,
  Config_HTTP_Message_free,
  Config_SelectedRow,
  Config_SocketSelectedRow,
} from "../../../config/Config.js";
import {
  GetAllStream,
  GetHTTPRequestBody,
  GetHTTPResponseBody,
  GetHTTPSession,
  GetLocalServerPATH,
  GetSessionMessageBody,
  GetSocketFilter, McpFuncRes, ProtobufToJson,
  SetRequestNextBreakMode
} from "../../../../../bindings/changeme/Service/appmain.js";
import Table from "../../../Tools/table.vue"
import JavaScriptEdit from "../tool/Raw.vue"
import {
  base64ToBytes,
  bytesToString, PbJsonConvert,
  sanitizeHTML,
  StringOrBytesJoinToBytes,
  StringToBytes,
  toGBK,
  toUTF8
} from "../../../config/encoding.js";
import {
  bodyIsForm,
  ErrorReplace,
  headerArrayToLanguage,
  toHeader,
  toHeaderArray,
  toRequestCookiesHeader,
  toResponseCookiesHeader,
  urlToArgs
} from "../../../config/SunnyNetInfoApi.js";
import HexView from "../../../Tools/HexView.vue";
import JSon from "../../../Tools/JSon.vue";
import IMGView from "../../../Tools/IMGView.vue";
import {getImageType} from "../../../config/CheckImageType.js";
import WebSocketView from "../../../Tools/WebSocketView.vue";
import WebsocketActive from "../../../Tools/WebsocketActive.vue";
import {Keys_System_id_Current_Release, registerHotkeyFunction} from "../../../config/Keys";
const ClassMinName = "ag-icon ag-icon-minimize ag-panel-title-bar-button-icon"
const ClassMaxName = "ag-icon ag-icon-maximize ag-panel-title-bar-button-icon"
export default {
  components: {WebsocketActive, WebSocketView, JSon, IMGView, HexView, JavaScriptEdit, Table},
  computed: {
    UpdateIcon() {
      return this.iconClass
    },
    BreakResponse() {
      return this.SelectedRow['断点模式'] === 2;
    },
    isBreak() {
      return this.SelectedRow['断点模式'] !== 0;
    },
    GetBreakpointStateStyle() {
      if (this.SelectedRow['断点模式'] !== 0) {
        return "width: 100%;height: calc(100% - 30px);"
      }
      return "width: 100%;height: 100%"
    },
    isShowResponse() {
      return this.Show.showResponse
    },
    isWsResponse() {
      return this.Tabs[7].visible && this.Tabs[7].Show
    },
    isShowRequest() {
      return this.Show.showRequest;
    },
    Method() {
      return this.SelectedRow["方式"]
    },
    ResponseCode() {
      return this.SelectedRow["状态"]
    },
    RequestStyle() {
      return "height: 50%; flex: 0 0 auto;padding-top: calc(var(--ag-grid-size) * 2)"
    },
  },
  data() {
    return {
      RangeStyle: {
        splitA4: "width: 100%;height: 100%;",
        splitA3: "width: 100%;height: 100%;",
        ResponseTabs4: "width: 100%;height: 100%;",
      },
      Errors: {
        isError: false,
        value: "请求正在发送中...",
      },
      isChangeSize: false,
      iconClass: ClassMaxName,
      get SelectedRow() {
        return Config_SelectedRow.value
      },
      set SelectedRow(value) {
        Config_SelectedRow.value = value
      },
      get WebsocketSelectedRow() {
        return Config_SocketSelectedRow.value
      },
      set WebsocketSelectedRow(value) {
        Config_SocketSelectedRow.value = value
      },
      ResponseRangeStyle: "",
      Tabs: [
        {id: 0, name: "原始视图", class: "ag-tab  ag-tab-selected", visible: true, Show: true},
        {id: 1, name: "协议头", class: "ag-tab", visible: false, Show: true},
        {id: 2, name: "参数视图", class: "ag-tab", visible: false, Show: true},
        {
          id: 3, name: "请求数据", class: "ag-tab", visible: false, Show: {
            BodyArgs: false,
            Text: false,
            Index: 0
          },
        },
        {id: 4, name: "十六进制视图", class: "ag-tab", visible: false, Show: true},
        {id: 5, name: "Cookies", class: "ag-tab", visible: false, Show: true},
        {id: 6, name: "JSON视图", class: "ag-tab", visible: false, Show: true},
        {id: 7, name: "WebSocket", class: "ag-tab", visible: false, Show: true},
        {id: 8, name: "图片视图", class: "ag-tab", visible: false, Show: true}
      ],
      ResponseTabs: [
        {id: 0, name: "原始响应", class: "ag-tab ag-tab-selected", visible: true, Show: true},
        {id: 1, name: "协议头", class: "ag-tab", visible: false, Show: true},
        {id: 2, name: "响应文本", class: "ag-tab", visible: false, Show: true},
        {id: 3, name: "图片视图", class: "ag-tab ", visible: false, Show: true},
        {id: 4, name: "HTML视图", class: "ag-tab", visible: false, Show: true},
        {id: 5, name: "十六进制视图", class: "ag-tab", visible: false, Show: true},
        {id: 6, name: "Cookies", class: "ag-tab", visible: false, Show: true},
        {id: 7, name: "JSON视图", class: "ag-tab", visible: false, Show: true},
      ],
      WebSocketTabs: [
        {id: 0, name: "文本视图", class: "ag-tab", visible: false, Show: true},
        {id: 1, name: "十六进制视图", class: "ag-tab", visible: false, Show: true},
        {id: 2, name: "JSON视图", class: "ag-tab", visible: false, Show: true},
        {id: 3, name: "主动发送", class: "ag-tab ag-tab-selected", visible: true, Show: true},
      ],
      ResponseCodeStateStyle: "warning",
      ResponseHTMLValue: "",
      noResponse: true,
      Show: {
        bakSize: "",
        RequestMax: false,
        ResponseMax: false,
        showRequest: true,
        showResponse: true,
      },
    };
  },
  methods: {
    CalculationResponseRangeStyle() {
      this.$nextTick(() => {
        let too = 0;
        const el = this.$refs.splitA1;
        if (Config_GOOS_IsWindows.value) {
          too = el.clientHeight;
        } else {
          const footerEl = window.Footer?.$el;
          const a1Top = el.getBoundingClientRect().top + window.scrollY;
          const a1Height = el.clientHeight;
          const footerHeight = footerEl.clientHeight;
          too = a1Top + a1Height + footerHeight;
        }
        this.RangeStyle.splitA3 = "width: 100%;height: calc(100% - " + too + "px)";
        this.$refs.splitA2.style = this.RangeStyle.splitA3
        const ResponseTabs1 = parseInt(this.$refs.ResponseTabs1.clientHeight)
        const ResponseTabs2 = parseInt(this.$refs.ResponseTabs2.clientHeight)
        const ResponseTabs3 = parseInt(this.$refs.ResponseTabs3.clientHeight)
        this.RangeStyle.ResponseTabs4 = "width: 100%;height: calc(100% - " + (ResponseTabs1 + ResponseTabs2 + ResponseTabs3) + "px)";
      })
    },
    isWebsocketMessage() {
      const mod = this.SelectedRow["方式"] + "";
      return (mod.toLowerCase().indexOf("ws") !== -1) || (mod.toLowerCase().indexOf("websocket") !== -1);
    },
    handleMouseDown(event) {
      this.isChangeSize = true;
    },
    handleMouseMove(event) {
      if (!this.isChangeSize) return;

      const clientY = event.clientY;
      const docHeight = document.documentElement.clientHeight;
      const splitA0Height = parseInt(this.$refs.splitA0.clientHeight);

      // ✅ 处理顶部边界
      if (clientY + 5 < 60) {
        return;
      }

      // ✅ 处理 `noResponse` 的情况
      if (this.noResponse) {
        if (clientY > splitA0Height - 300) {
          return;
        }
      }

      // ✅ 处理底部边界
      if (clientY + 5 > docHeight - 80) {
        this.$refs.splitA1.style.height = (docHeight - 100) + 'px';
        return;
      }

      // ✅ 处理默认情况
      this.$refs.splitA1.style.height = (clientY + 5 - 30) + 'px';
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
      if (mode === 1) {
        if (this.iconClass !== ClassMaxName) {
          this.iconClass = ClassMaxName
          const Tools = this.getTools()
          if (Tools !== null) {
            Tools.style.width = this.RawSize
          }
          this.$refs.splitA1.style.height = this.Show.bakSize
          this.Show.showResponse = true
          this.Show.showRequest = true
          return;
        }
        this.iconClass = ClassMinName
        const Tools = this.getTools()
        if (Tools !== null) {
          this.RawSize = Tools.style.width
          Tools.style.width = window.innerWidth + "px"
        }
        this.Show.bakSize = this.$refs.splitA1.style.height
        this.$refs.splitA1.style.height = "100%"
        this.Show.showResponse = false
        this.Show.showRequest = true
        return
      }
      if (this.iconClass !== ClassMaxName) {
        this.iconClass = ClassMaxName
        const Tools = this.getTools()
        if (Tools !== null) {
          Tools.style.width = this.RawSize
        }
        this.$refs.splitA1.style.height = this.Show.bakSize
        this.Show.showResponse = true
        this.Show.showRequest = true
        return;
      }
      this.iconClass = ClassMinName
      const Tools = this.getTools()
      if (Tools !== null) {
        this.RawSize = Tools.style.width
        Tools.style.width = window.innerWidth + "px"
      }
      this.Show.bakSize = this.$refs.splitA1.style.height
      this.$refs.splitA1.style.height = "100%"
      this.Show.showResponse = true
      this.Show.showRequest = false
      return
    },
    eHeaderClick(eve) {
      this.Tabs.forEach((obj) => {
        if (obj.visible) {
          this.applyCheckRequestData(obj, this.SelectedRow)
        }
        obj.class = "ag-tab"
        obj.visible = false
      })
      eve.class = "ag-tab ag-tab-selected"
      eve.visible = true
    },
    eResponseClick(eve) {
      this.ResponseTabs.forEach((obj) => {
        if (obj.visible) {
          this.applyCheckResponse(obj, this.SelectedRow)
        }
        obj.class = "ag-tab"
        obj.visible = false
      })
      eve.class = "ag-tab ag-tab-selected"
      eve.visible = true
    },
    eWebsocketClick(eve) {
      this.WebSocketTabs.forEach((obj) => {
        obj.class = "ag-tab"
        obj.visible = false
      })
      eve.class = "ag-tab ag-tab-selected"
      eve.visible = true
    },
    RequestRawTextTypeClick() {
      this.applyCheckRequestData(this.Tabs[3], this.SelectedRow);
      if (this.Tabs[3].Show.Index === 0) {
        this.Tabs[3].Show.Index = 1;
      } else {
        this.Tabs[3].Show.Index = 0;
      }
    },
    initValue() {
      this.initValueFunc().then(() => {
        //Request tab
        {
          // 寻找当前选中的 tab
          const selectedTab = this.Tabs.find((tab) => tab.class === "ag-tab ag-tab-selected");
          // 如果当前选中的 tab 是"请求数据"且 Text 为空,或者 Show 为空,则切换到第一个 tab
          if (selectedTab && selectedTab.name === "请求数据" && (!selectedTab.Show || !selectedTab.Show.Text)) {
            this.eHeaderClick(this.Tabs[0]);
          } else if (selectedTab && !selectedTab.Show) {
            this.eHeaderClick(this.Tabs[0]);
          }
        }
        //Response tab
        {
          // 寻找当前选中的 tab
          const selectedTab = this.ResponseTabs.find((tab) => tab.class === "ag-tab ag-tab-selected");
          // 如果当前选中的 tab  Show 为空,则切换到第一个 tab
          if (selectedTab && !selectedTab.Show) {
            this.eResponseClick(this.ResponseTabs[0]);
          }
        }
        //响应区域高度
        this.ResponseRangeStyle = "width: 100%;height: calc(100% - " + (this.$refs.ResponseTabs1.clientHeight + this.$refs.ResponseTabs2.clientHeight + this.$refs.ResponseTabs3.clientHeight) + "px)";
        this.CalculationResponseRangeStyle()
      });
    },
    async initValueFull() {
      {
        if (this.SelectedRow.Theology === undefined) {
          return
        }
        if (this.Stream === null) {
          return
        }
        this.noResponse = true;
        this.Errors.isError = false
      }
      this.Errors.value = "完整数据正在加载中..."
      if (this.Stream.State === 3) {
        this.Errors.isError = true
        this.Errors.value = ErrorReplace(this.Stream.Error);
        this.$refs.splitA1.style.height = "50%";
      }
      this.Stream.Request.Body = await GetHTTPRequestBody(this.Stream.Theology, true)
      const RequestBody = base64ToBytes(this.Stream.Request.Body)
      const RequestHeader = toHeaderArray(this.Stream.Request.Header)
      let language = headerArrayToLanguage(RequestHeader)
      this.Stream.Request.IsMaxLength = false
      //请求
      {
        const ReadOnly = this.SelectedRow['断点模式'] !== 1
        {
          this.$refs.RequestJson.SetReadOnly(ReadOnly)
          this.$refs.RequestUrlArgs.SetReadOnly(ReadOnly)
          this.$refs.RequestCookie.SetReadOnly(ReadOnly)
          this.$refs.RequestHex.SetReadOnly(ReadOnly)
          this.$refs.RequestHeader.SetReadOnly(ReadOnly)
          this.$refs.RequestBodyTable.SetReadOnly(ReadOnly)
          this.$refs.RequestRaw.SetReadOnly(ReadOnly)
          this.$refs.RequestRawText.SetReadOnly(ReadOnly)
        }
        //原始视图 AND HEX
        {
          let message = `${this.Stream.Request.Method} ${this.Stream.Request.Url} ${this.Stream.Request.Proto}\r\n`;
          message += toHeader(this.Stream.Request.Header) + "\r\n"
          const bs = StringOrBytesJoinToBytes(message, RequestBody);
          this.$refs.RequestRaw.SetCode(bs)
          this.$refs.RequestHex.SetCode(bs)
          this.Tabs[0].Show = true;
          this.Tabs[4].Show = true;
        }
        //协议头视图
        {
          this.$refs.RequestHeader.Empty()
          await this.$refs.RequestHeader.AddLines(RequestHeader)
          this.Tabs[1].Show = true;
        }
        //参数视图
        {
          this.$refs.RequestUrlArgs.Empty()
          const array = urlToArgs(this.Stream.Request.Url);
          await this.$refs.RequestUrlArgs.AddLines(array)
          this.Tabs[2].Show = array.length > 0;
        }
        //POST请求参数视图 And JSON视图
        {
          this.Tabs[3].Show.Text = false
          this.$refs.RequestBodyTable.Empty()
          this.$refs.RequestRawText.SetCode(RequestBody)
          let json = null;
          try {
            json = JSON.parse(toGBK(RequestBody));
          } catch (e) {
            try {
              json = JSON.parse(toUTF8(RequestBody));
            } catch (e) {
            }
          }
          if (typeof json === 'object' && json !== null) {
            language = "json"
            this.Tabs[6].Show = true
          } else {
            this.Tabs[6].Show = false
            {
              //is pb 数据?
              for (let i = 0; i < Math.min(RequestBody.length, 20); i++) {
                if (RequestBody[i] < 20) {
                  this.Tabs[6].Show = true
                  break
                }
              }
            }
          }
          if (this.Tabs[6].Show) {
            this.$refs.RequestJson.SetCode(RequestBody)
          }
          this.$refs.RequestRawText.SetLanguage(language)
          this.$refs.RequestRaw.SetLanguage(language)
          if (RequestBody.length > 0) {
            const BodyTable = bodyIsForm(StringOrBytesJoinToBytes("data?", RequestBody))
            this.$refs.RequestBodyTable.AddLines(BodyTable)
            this.Tabs[3].Show.BodyArgs = BodyTable.length > 0
            this.Tabs[3].Show.Text = true
          }
        }
        // 请求Cookie
        {
          this.$refs.RequestCookie.Empty()
          const array = toRequestCookiesHeader(RequestHeader);
          this.Tabs[5].Show = array.length > 0;
          this.$refs.RequestCookie.AddLines(array)
        }
        // 图片视图
        {
          const ImageType = getImageType(RequestBody);
          this.$refs.RequestImg.SetImg(this.Stream.Request.Body, ImageType)
          this.Tabs[8].Show = ImageType !== null;
        }

      }
      //响应
      {

        if (this.Stream.Ico === "上行" || this.Stream.Ico === "拦截上行") {
          return
        }

        this.Stream.Response.Body = await GetHTTPResponseBody(this.Stream.Theology, true)
        const ResponseBody = base64ToBytes(this.Stream.Response.Body)
        const ResponseHeader = toHeaderArray(this.Stream.Response.Header)
        let ResponseLanguage = headerArrayToLanguage(ResponseHeader)
        this.Stream.Response.IsMaxLength = false

        this.noResponse = false;
        const ReadOnly = this.SelectedRow['断点模式'] !== 2
        {
          this.$refs.ResponseRaw.SetReadOnly(ReadOnly)
          this.$refs.ResponseHeader.SetReadOnly(ReadOnly)
          this.$refs.ResponseText.SetReadOnly(ReadOnly)
          this.$refs.ResponseHEX.SetReadOnly(ReadOnly)
          this.$refs.ResponseCookie.SetReadOnly(ReadOnly)
          this.$refs.ResponseJson.SetReadOnly(ReadOnly)
        }
        //响应 原始视图 AND 文本视图 AND HEX
        {
          let message = this.Stream.Response.Proto + " " + this.Stream.Response.Code + " " + this.Stream.Response.State + "\r\n" + toHeader(this.Stream.Response.Header) + "\r\n"
          const bs = StringOrBytesJoinToBytes(message, ResponseBody);
          this.$refs.ResponseRaw.SetCode(bs)
          this.$refs.ResponseText.SetCode(ResponseBody)
          this.$refs.ResponseText.SetLanguage(ResponseLanguage)
          this.$refs.ResponseRaw.SetLanguage(ResponseLanguage)
          this.$refs.ResponseHEX.SetCode(bs)
          this.ResponseTabs[0].Show = true;
          this.ResponseTabs[2].Show = true;
          this.ResponseTabs[5].Show = true;
        }
        //响应 协议头视图
        {
          this.$refs.ResponseHeader.Empty()
          this.$refs.ResponseHeader.AddLines(ResponseHeader)
          this.$refs.ResponseHeader.SetReadOnly(ReadOnly)
          this.ResponseTabs[1].Show = true;
        }
        //响应 图片视图
        {
          const ImageType = getImageType(ResponseBody);
          this.ResponseTabs[3].Show = ImageType !== null
          if (this.ResponseTabs[3].Show) {
            if (this.Stream.Response.IsMaxLength) {
              GetHTTPResponseBody(parseInt(this.SelectedRow.Theology), true).then(data => {
                this.$refs.ResponseImg.SetImg(data, ImageType)
              })
            } else {
              this.$refs.ResponseImg.SetImg(this.Stream.Response.Body, ImageType)
            }
          }
        }
        //响应 HTML
        {
          this.ResponseHTMLValue = sanitizeHTML(bytesToString(ResponseBody));
          this.ResponseTabs[4].Show = ResponseLanguage === "html";
        }

        //Cookie
        {
          this.$refs.ResponseCookie.Empty()
          const array = toResponseCookiesHeader(RequestHeader);
          this.ResponseTabs[6].Show = array.length > 0;
          this.$refs.ResponseCookie.AddLines(array)
        }

        //JSON
        {
          let json = null;
          try {
            json = JSON.parse(toGBK(ResponseBody));
          } catch (e) {
            try {
              json = JSON.parse(toUTF8(ResponseBody));
            } catch (e) {
            }
          }
          if (typeof json === 'object' && json !== null) {
            this.ResponseTabs[7].Show = true
          } else {
            this.ResponseTabs[7].Show = false
            {
              //is pb 数据?
              for (let i = 0; i < Math.min(ResponseBody.length, 20); i++) {
                if (ResponseBody[i] < 20) {
                  this.ResponseTabs[7].Show = true
                  break
                }
              }
            }
          }
          if (this.ResponseTabs[7].Show) {
            this.$refs.ResponseJson.SetCode(ResponseBody)
          }
        }

      }
    },
    async initValueFunc() {
      {
        if (this.SelectedRow.Theology === undefined) {
          return
        }
        if (this.Stream === null) {
          return
        }
        this.noResponse = true;
        this.Errors.isError = false
      }
      if (this.Stream.State === 1) {
        this.Errors.value = "请求正在发送中..."
      }
      if (this.Stream.State === 3) {
        this.Errors.isError = true
        this.Errors.value = ErrorReplace(this.Stream.Error);
        this.$refs.splitA1.style.height = "50%";
      }
      this.Stream.Request.Body = await GetHTTPRequestBody(this.Stream.Theology, false)
      const RequestBody = base64ToBytes(this.Stream.Request.Body)
      const RequestHeader = toHeaderArray(this.Stream.Request.Header)
      let language = headerArrayToLanguage(RequestHeader)
      //请求
      {
        const ReadOnly = this.SelectedRow['断点模式'] !== 1
        {
          this.$refs.RequestJson.SetReadOnly(ReadOnly)
          this.$refs.RequestUrlArgs.SetReadOnly(ReadOnly)
          this.$refs.RequestCookie.SetReadOnly(ReadOnly)
          this.$refs.RequestHex.SetReadOnly(ReadOnly)
          this.$refs.RequestHeader.SetReadOnly(ReadOnly)
          this.$refs.RequestBodyTable.SetReadOnly(ReadOnly)
          this.$refs.RequestRaw.SetReadOnly(ReadOnly)
          this.$refs.RequestRawText.SetReadOnly(ReadOnly)
        }
        //原始视图 AND HEX
        {
          let message = `${this.Stream.Request.Method} ${this.Stream.Request.Url} ${this.Stream.Request.Proto}\r\n`;
          message += toHeader(this.Stream.Request.Header) + "\r\n"
          const bs = StringOrBytesJoinToBytes(message, RequestBody);
          if (this.Stream.Request.IsMaxLength) {
            const bs1 = StringOrBytesJoinToBytes(RequestBody, "\n{.....数据过长未完全加载.....}\n{...您可以右键选择-展示完整数据...}\n{...展示完整数据,可能 UI 较卡顿...}");
            const bs2 = StringOrBytesJoinToBytes(message, bs1);
            this.$refs.RequestRaw.SetCode(bs2, () => {
              this.initValueFull()
            })
            this.$refs.RequestRaw.SetReadOnly(true)
          } else {
            this.$refs.RequestRaw.SetCode(bs)
          }
          this.$refs.RequestHex.SetCode(bs)
          this.Tabs[0].Show = true;
          this.Tabs[4].Show = true;
        }
        //协议头视图
        {
          this.$refs.RequestHeader.Empty()
          this.$refs.RequestHeader.AddLines(RequestHeader)
          this.Tabs[1].Show = true;
        }
        //参数视图
        {
          this.$refs.RequestUrlArgs.Empty()
          const array = urlToArgs(this.Stream.Request.Url);
          this.$refs.RequestUrlArgs.AddLines(array)
          this.Tabs[2].Show = array.length > 0;
        }
        //POST请求参数视图 And JSON视图
        {
          if (this.Stream.Request.IsMaxLength) {
            const bs1 = StringOrBytesJoinToBytes(RequestBody, "\n{.....数据过长未完全加载.....}\n{...您可以右键选择-展示完整数据...}\n{...展示完整数据,可能 UI 较卡顿...}");
            this.$refs.RequestRawText.SetCode(bs1, () => {
              this.initValueFull()
            })
            this.$refs.RequestRawText.SetReadOnly(true)
            this.Tabs[6].Show = false
            this.Tabs[3].Show.BodyArgs = false
            this.Tabs[3].Show.Text = true
            this.$refs.RequestJson.SetCode("{}")
          } else {
            {
              this.$refs.RequestRawText.SetCode(RequestBody)

              this.Tabs[3].Show.Text = false
              this.$refs.RequestBodyTable.Empty()
              let json = null;
              try {
                json = JSON.parse(toGBK(RequestBody));
              } catch (e) {
                try {
                  json = JSON.parse(toUTF8(RequestBody));
                } catch (e) {
                }
              }
              if (typeof json === 'object' && json !== null) {
                language = "json"
                this.Tabs[6].Show = true
              } else {
                this.Tabs[6].Show = false
                {
                  //is pb 数据?
                  for (let i = 0; i < Math.min(RequestBody.length, 20); i++) {
                    if (RequestBody[i] < 20) {
                      this.Tabs[6].Show = true
                      break
                    }
                  }
                }
              }
              if (this.Tabs[6].Show) {
                this.$refs.RequestJson.SetCode(RequestBody)
              }
              this.$refs.RequestRawText.SetLanguage(language)
              this.$refs.RequestRaw.SetLanguage(language)
              if (RequestBody.length > 0) {
                const BodyTable = bodyIsForm(StringOrBytesJoinToBytes("data?", RequestBody))
                this.$refs.RequestBodyTable.AddLines(BodyTable)
                this.Tabs[3].Show.BodyArgs = BodyTable.length > 0
                this.Tabs[3].Show.Text = true
              }
            }
          }
        }
        // 请求Cookie
        {
          this.$refs.RequestCookie.Empty()
          const array = toRequestCookiesHeader(RequestHeader);
          this.Tabs[5].Show = array.length > 0;
          this.$refs.RequestCookie.AddLines(array)
        }
        // 图片视图
        {
          if (this.Stream.Request.IsMaxLength) {
            this.Tabs[8].Show = false;
          } else {
            const ImageType = getImageType(RequestBody);
            this.$refs.RequestImg.SetImg(this.Stream.Request.Body, ImageType)
            this.Tabs[8].Show = ImageType !== null;
          }
        }
        //WebSocketView
        {
        }
      }
      //响应
      {
        if (this.Stream.Ico === "上行" || this.Stream.Ico === "拦截上行") {
          return
        }
        this.Stream.Response.Body = await GetHTTPResponseBody(this.Stream.Theology, false)
        const ResponseBody = base64ToBytes(this.Stream.Response.Body)
        const ResponseHeader = toHeaderArray(this.Stream.Response.Header)
        let ResponseLanguage = headerArrayToLanguage(ResponseHeader)

        this.noResponse = false;
        const ReadOnly = this.SelectedRow['断点模式'] !== 2
        {
          this.$refs.ResponseRaw.SetReadOnly(ReadOnly)
          this.$refs.ResponseHeader.SetReadOnly(ReadOnly)
          this.$refs.ResponseText.SetReadOnly(ReadOnly)
          this.$refs.ResponseHEX.SetReadOnly(ReadOnly)
          this.$refs.ResponseCookie.SetReadOnly(ReadOnly)
          this.$refs.ResponseJson.SetReadOnly(ReadOnly)
        }
        //响应 原始视图 AND 文本视图 AND HEX
        {
          let message = this.Stream.Response.Proto + " " + this.Stream.Response.Code + " " + this.Stream.Response.State + "\r\n" + toHeader(this.Stream.Response.Header) + "\r\n"
          const bs = StringOrBytesJoinToBytes(message, ResponseBody);
          if (this.Stream.Response.IsMaxLength) {
            if (ResponseLanguage.indexOf("video") !== -1 || ResponseLanguage.indexOf("octet-stream") !== -1 || ResponseLanguage.indexOf("audio") !== -1) {
              const bytes = "{.....较大的二进制数据.....}\n{...您可以右键选择-展示完整数据...}\n{...展示完整数据,可能 UI 较卡顿...}";
              const bs2 = StringOrBytesJoinToBytes(message, bytes);
              this.$refs.ResponseRaw.SetCode(bs2, () => {
                this.initValueFull()
              })
              this.$refs.ResponseText.SetCode(() => {
                this.initValueFull()
              })
            } else {
              const bs1 = StringOrBytesJoinToBytes(ResponseBody, "\n{.....数据过长未完全加载.....}\n{...您可以右键选择-展示完整数据...}\n{...展示完整数据,可能 UI 较卡顿...}");
              const bs2 = StringOrBytesJoinToBytes(message, bs1);
              this.$refs.ResponseText.SetCode(bs1, () => {
                this.initValueFull()
              })
              this.$refs.ResponseRaw.SetCode(bs2, () => {
                this.initValueFull()
              })
            }
            this.$refs.ResponseRaw.SetReadOnly(true)
            this.$refs.ResponseText.SetReadOnly(true)
          } else {
            this.$refs.ResponseRaw.SetCode(bs)
            this.$refs.ResponseText.SetCode(ResponseBody)
          }

          this.$refs.ResponseText.SetLanguage(ResponseLanguage)
          this.$refs.ResponseRaw.SetLanguage(ResponseLanguage)
          this.$refs.ResponseHEX.SetCode(bs)
          this.ResponseTabs[0].Show = true;
          this.ResponseTabs[2].Show = true;
          this.ResponseTabs[5].Show = true;
        }
        //响应 协议头视图
        {
          this.$refs.ResponseHeader.Empty()
          this.$refs.ResponseHeader.AddLines(ResponseHeader)
          this.$refs.ResponseHeader.SetReadOnly(ReadOnly)
          this.ResponseTabs[1].Show = true;
        }
        //响应 图片视图
        {
          const ImageType = getImageType(ResponseBody);
          this.ResponseTabs[3].Show = ImageType !== null
          if (this.ResponseTabs[3].Show) {
            if (this.Stream.Response.IsMaxLength) {
              GetHTTPResponseBody(parseInt(this.SelectedRow.Theology), true).then(data => {
                this.$refs.ResponseImg.SetImg(data, ImageType)
              })
            } else {
              this.$refs.ResponseImg.SetImg(this.Stream.Response.Body, ImageType)
            }
          }
        }
        //响应 HTML
        {
          if (this.Stream.Response.IsMaxLength) {
            this.ResponseTabs[4].Show = false
          } else {
            this.ResponseHTMLValue = sanitizeHTML(bytesToString(ResponseBody));
            this.ResponseTabs[4].Show = ResponseLanguage === "html";
          }
        }

        //Cookie
        {
          this.$refs.ResponseCookie.Empty()
          const array = toResponseCookiesHeader(RequestHeader);
          this.ResponseTabs[6].Show = array.length > 0;
          this.$refs.ResponseCookie.AddLines(array)
        }

        //JSON
        {
          if (this.Stream.Response.IsMaxLength) {
            this.$refs.ResponseJson.SetCode("{}")
            this.ResponseTabs[7].Show = false
          } else {
            let json = null;
            try {
              json = JSON.parse(toGBK(ResponseBody));
            } catch (e) {
              try {
                json = JSON.parse(toUTF8(ResponseBody));
              } catch (e) {
              }
            }
            if (typeof json === 'object' && json !== null) {
              this.ResponseTabs[7].Show = true
            } else {
              this.ResponseTabs[7].Show = false
              {
                //is pb 数据?
                for (let i = 0; i < Math.min(ResponseBody.length, 20); i++) {
                  if (ResponseBody[i] < 20) {
                    this.ResponseTabs[7].Show = true
                    break
                  }
                }
              }
            }
            if (this.ResponseTabs[7].Show) {
              this.$refs.ResponseJson.SetCode(ResponseBody)
            }
          }
        }

      }
    },
    async BreakClick(mode) {
      for (let i = 0; i < this.Tabs.length; i++) {
        if (this.Tabs[i].visible) {
          await this.applyCheckRequestData(this.Tabs[i], this.SelectedRow, true)
          break
        }
      }
      for (let i = 0; i < this.ResponseTabs.length; i++) {
        if (this.ResponseTabs[i].visible) {
          await this.applyCheckResponse(this.ResponseTabs[i], this.SelectedRow, true)
          break
        }
      }
      const _SelectedTheology = parseInt(this.SelectedRow?.Theology ?? "0");
      SetRequestNextBreakMode(_SelectedTheology, mode)
    },
    async applyCheckRequestData(obj, row, no) {
      const isBreakRequest = row['断点模式'] === 1;
      if (!isBreakRequest) {
        return
      }
      const _SelectedTheology = parseInt(row?.Theology ?? "0");
      let isEdit = false;
      switch (obj.name) {
        case "原始视图":
          const a1 = this.$refs.RequestRaw.GetRequest();
          if (a1) {
            this.Stream.Request.Body = a1.data;
            this.Stream.Request.Header = a1.Header;
            this.Stream.Request.Method = a1.Method;
            this.Stream.Request.Url = a1.URL;
            this.Stream.Request.Proto = a1.Proto;
            isEdit = true;
          }
          break
        case "Cookies":
          const a5 = this.$refs.RequestCookie.toRequestCookie();
          if (a5) {
            this.Stream.Request.Header = this.$refs.RequestHeader.updateRequestCookie(a5);
            isEdit = true;
          }
          break
        case "协议头":
          const a2 = this.$refs.RequestHeader.toRequestHeader();
          if (a2) {
            this.Stream.Request.Header = a2;
            isEdit = true;
          }
          break
        case "请求数据":
          if (obj.Show.Index === 0) {
            const a3 = this.$refs.RequestRawText.GetRawData();
            if (a3) {
              this.Stream.Request.Body = a3;
              isEdit = true;
            }
          }
          if (obj.Show.Index === 1) {
            const a3 = this.$refs.RequestBodyTable.toRawBodyTable();
            if (a3) {
              this.Stream.Request.Body = a3;
              isEdit = true;
            }
          }
          break
        case "JSON视图":
          const a3 = this.$refs.RequestJson.GetRawCode();
          if (a3) {
            this.Stream.Request.Body = a3;
            isEdit = true;
          }
          break
        case "十六进制视图":
          const a4 = this.$refs.RequestHex.GetALLBytes();
          if (a4) {
            this.Stream.Request.Body = a4.data;
            this.Stream.Request.Header = a4.Header;
            this.Stream.Request.Method = a4.Method;
            this.Stream.Request.Url = a4.URL;
            this.Stream.Request.Proto = a4.Proto;
            isEdit = true;
          }
          break
        case "WebSocket":
          //拦截状态不可能出现 Websocket 标签
          break
        case "图片视图":
          //图片视图 不支持修改
          break
      }
      if (isEdit) {
        await this.applyRequest(_SelectedTheology);
        if (no) {
          return
        }
        const _SelectedTheology2 = parseInt(this.SelectedRow?.Theology ?? "0");
        if (_SelectedTheology2 === _SelectedTheology) {
          if (!this.Stream.Request.IsMaxLength) {
            await this.initValueFull()
          } else {
            await this.initValueFunc()
          }
        }
      }
    },
    async applyCheckResponse(obj, row) {
      const isBreakResponse = row['断点模式'] === 2;
      if (!isBreakResponse) {
        return
      }
      const _SelectedTheology = parseInt(row?.Theology ?? "0");
      let isEdit = false;
      switch (obj.name) {
        case "原始响应":
          const a1 = this.$refs.ResponseRaw.GetResponse();
          if (a1) {
            this.Stream.Response.Body = a1.data;
            this.Stream.Response.Code = a1.Code;
            this.Stream.Response.State = a1.StateText;
            this.Stream.Response.Header = a1.Header;
            this.Stream.Response.Proto = a1.Proto;
            isEdit = true;
          }
          break
        case "协议头":
          const a2 = this.$refs.ResponseHeader.toRequestHeader();
          if (a2) {
            this.Stream.Response.Header = a2;
            isEdit = true;
          }
          break
        case "Cookies":
          const a5 = this.$refs.ResponseCookie.toResponseCookiesHeader();
          if (a5) {
            this.Stream.Response.Header = this.$refs.RequestHeader.updateResponseCookie(a5);
            isEdit = true;
          }
          break
        case "响应文本":
          const a6 = this.$refs.ResponseText.GetRawData();
          if (a6) {
            this.Stream.Response.Body = a6;
            isEdit = true;
          }
          break
        case "JSON视图":
          const a3 = this.$refs.ResponseJson.GetRawCode();
          if (a3) {
            this.Stream.Response.Body = a3;
            isEdit = true;
          }
          break
        case "十六进制视图":
          const a4 = this.$refs.ResponseHEX.GetResponseBytes();
          if (a4) {
            this.Stream.Response.Body = a4.data;
            this.Stream.Response.Code = a4.Code;
            this.Stream.Response.State = a4.StateText;
            this.Stream.Response.Header = a4.Header;
            this.Stream.Response.Proto = a4.Proto;
            isEdit = true;
          }
          break
        case "图片视图":
          //图片视图 不支持修改
          break
        case "HTML视图":
          //HTML视图 不支持修改
          break
      }
      if (isEdit) {
        await this.applyResponse(_SelectedTheology);
        const _SelectedTheology2 = parseInt(this.SelectedRow?.Theology ?? "0");
        if (_SelectedTheology2 === _SelectedTheology) {
          if (!this.Stream.Response.IsMaxLength) {
            await this.initValueFull()
          } else {
            await this.initValueFunc()
          }
        }
      }
    },
    async applyRequest(Theology) {
      const path = await GetLocalServerPATH();
      const formData = JSON.stringify(this.Stream.Request);
      return await fetch(path + '/UpdateHttpRequest?Theology=' + Theology, {
        method: 'POST',
        body: formData
      })
    },
    async applyResponse(Theology) {
      const path = await GetLocalServerPATH();
      const formData = JSON.stringify(this.Stream.Response);
      return await fetch(path + '/UpdateHttpResponse?Theology=' + Theology, {
        method: 'POST',
        body: formData
      })
    },
    free(lodRow) {
      this.Tabs.forEach((obj) => {
        if (obj.visible) {
          this.applyCheckRequestData(obj, lodRow)
        }
      })
      this.ResponseTabs.forEach((obj) => {
        if (obj.visible) {
          this.applyCheckResponse(obj, lodRow)
        }
      })
    }
  },
  watch: {
    "SelectedRow"(newValue, lodRow) {
      this.free(lodRow);
      const _SelectedTheology = parseInt(newValue?.Theology ?? "0");
      if (_SelectedTheology === 0) {
        this.Stream = null;
        return
      }
      if (!newValue["请求地址"].toLowerCase().startsWith("http")) {
        this.Stream = null;
        return
      }
      this.Errors.isError = false
      this.$refs.RequestWebSocket.Empty()
      GetSocketFilter(_SelectedTheology).then((Model) => {
        this.$refs.RequestWebSocket.$refs.Filter.setFilter(Model)
        GetHTTPSession(_SelectedTheology).then((res) => {
          this.Tabs.forEach((obj) => {
            switch (obj.name) {
              case "WebSocket":
                obj.Show = this.isWebsocketMessage()
                break
              case "协议头":
                obj.Show = true
                break
              case "请求数据":
                try {
                  const l = res.Request.Body.length;
                  if (l > 0) {
                    obj.Show.Text = true;
                    obj.Show.BodyArgs = false;
                    break
                  }
                } catch (e) {
                }
                obj.Show.Text = false;
                obj.Show.BodyArgs = false;
                break
            }
          })
          this.Stream = res;
          this.initValue()
          //WebSocketView
          {
            this.$refs.RequestWebSocket.Empty()
            this.Tabs[7].Show = this.isWebsocketMessage();
            if (this.Tabs[7].Show) {
              GetAllStream(parseInt(this.SelectedRow.Theology)).then((res) => {
                this.$refs.RequestWebSocket.InsertSocketStream(res, true)
              });
            }
          }
        })
      })
    },
    "WebsocketSelectedRow"(newValue) {
      const _SelectedTheology = parseInt(this.SelectedRow?.Theology ?? "0");
      if (_SelectedTheology === 0) {
        this.Stream = null;
        return
      }
      if (!this.SelectedRow["请求地址"].toLowerCase().startsWith("http")) {
        this.Stream = null;
        return
      }

      const MessageId = parseInt(newValue?.MessageId ?? "0");
      if (MessageId < 1) {
        const nil = StringToBytes("");
        this.$refs.WebSocketHEX.SetCode(nil)
        this.$refs.WebSocketJson.SetCode(nil)
        this.$refs.WebSocketText.SetCode(nil)
        this.$refs.WebSocketActive.SetCode(nil)
        return;
      }
      const IsClose = newValue?.IsClose ?? false;
      GetSessionMessageBody(parseInt(this.SelectedRow.Theology), MessageId).then((res) => {
        const bs = base64ToBytes(res);
        this.$refs.WebSocketHEX.SetCode(bs)
        this.$refs.WebSocketJson.SetCode(bs)
        this.$refs.WebSocketText.SetCode(bs)
        this.$refs.WebSocketActive.SetCode(bs, IsClose)
      })
    }
  },
  mounted() {
    {
      document.addEventListener('mouseup', (event) => {
        if (this.isChangeSize) {
          this.isChangeSize = false;
        }
      });
      document.addEventListener('mousemove', this.handleMouseMove);
    }
    {
      const elementRef = this.$refs.splitA1;
      const resizeObserver = new ResizeObserver(entries => {
        for (const entry of entries) {
          //const {width, height} = entry.contentRect;
          this.CalculationResponseRangeStyle()
        }
      });
      resizeObserver.observe(elementRef); // 开始监听元素尺寸变化
    }
    Config_HTTP_Message_free.value = this.free
    registerHotkeyFunction(Keys_System_id_Current_Release, () => {
      if (this.SelectedRow.Theology === undefined) {
        return
      }
      if (this.SelectedRow['断点模式'] !== 0) {
        this.BreakClick(0)
      }
    })

  }
}
</script>