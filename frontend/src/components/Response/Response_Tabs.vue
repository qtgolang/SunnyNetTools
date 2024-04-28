<template>
  <div class="ag-tabs ag-chart-tabbed-menu ag-focus-managed" style="position:relative">
    <div role="tablist" class="ag-tabs-header ag-chart-tabbed-menu-header" style="overflow: hidden;">
      <span v-if="DisplayTCPResponse===false" v-show="HTTPTabs[0].Show" @click="eHeaderClick(HTTPTabs[0])"
            :class="HTTPTabs[0].class"
            role="tab"> {{ HTTPTabs[0].name }} </span>
      <span v-if="DisplayTCPResponse===false" v-show="HTTPTabs[1].Show" @click="eHeaderClick(HTTPTabs[1])"
            :class="HTTPTabs[1].class"
            role="tab"> {{ HTTPTabs[1].name }} </span>
      <span v-if="DisplayTCPResponse===false" v-show="HTTPTabs[2].Show" @click="eHeaderClick(HTTPTabs[2])"
            :class="HTTPTabs[2].class"
            role="tab"> {{ HTTPTabs[2].name }} </span>
      <span v-if="DisplayTCPResponse===false" v-show="HTTPTabs[3].Show" @click="eHeaderClick(HTTPTabs[3])"
            :class="HTTPTabs[3].class"
            role="tab"> {{ HTTPTabs[3].name }} </span>
      <span v-if="DisplayTCPResponse===false" v-show="HTTPTabs[4].Show" @click="eHeaderClick(HTTPTabs[4])"
            :class="HTTPTabs[4].class"
            role="tab"> {{ HTTPTabs[4].name }} </span>
      <span v-if="DisplayTCPResponse===false" v-show="HTTPTabs[5].Show" @click="eHeaderClick(HTTPTabs[5])"
            :class="HTTPTabs[5].class"
            role="tab"> {{ HTTPTabs[5].name }} </span>
      <span v-if="DisplayTCPResponse===false" v-show="HTTPTabs[6].Show" @click="eHeaderClick(HTTPTabs[6])"
            :class="HTTPTabs[6].class"
            role="tab"> {{ HTTPTabs[6].name }} </span>
      <span v-if="DisplayTCPResponse===false" v-show="HTTPTabs[7].Show" @click="eHeaderClick(HTTPTabs[7])"
            :class="HTTPTabs[7].class"
            role="tab"> {{ HTTPTabs[7].name }} </span>
      <span v-if="DisplayTCPResponse" v-show="TCPTabs[0].Show" @click="eTcpClick(TCPTabs[0])" :class="TCPTabs[0].class"
            role="tab"> {{ TCPTabs[0].name }} </span>
      <span v-if="DisplayTCPResponse" v-show="TCPTabs[1].Show" @click="eTcpClick(TCPTabs[1])" :class="TCPTabs[1].class"
            role="tab"> {{ TCPTabs[1].name }} </span>
      <span v-if="DisplayTCPResponse" v-show="TCPTabs[2].Show" @click="eTcpClick(TCPTabs[2])" :class="TCPTabs[2].class"
            role="tab"> {{ TCPTabs[2].name }} </span>
      <span v-if="DisplayTCPResponse" v-show="TCPTabs[3].Show" @click="eTcpClick(TCPTabs[3])" :class="TCPTabs[3].class"
            role="tab"> {{ TCPTabs[3].name }} </span>
      <span v-if="DisplayTCPResponse" :class="UpdateIcon" unselectable="on" role="presentation"
            style="right: 3px;position:absolute;cursor:pointer;top:5px" @click="SetMaximize(2)"></span>
    </div>
    <div ref="BodyRect" role="presentation" class="ag-tabs-body ag-chart-tabbed-menu-body">
      <div :style="{width: '100%' ,height: BodyRectHeight,position: 'absolute'}">
        <div v-show="HTTPTabs[0].visible&&DisplayTCPResponse===false" style="width: 100%;height: 100%">
          <JavaScriptEdit ref="Raw" :height="BodyRectHeight" :glyphMargin="false" :readOnly="readOnly" Text=""
                          Name="Response"/>
        </div>
        <div v-show="HTTPTabs[1].visible&&DisplayTCPResponse===false" style="width: 100%;height: 100%">
          <List ref="Headers" :readOnly="readOnly"/>
        </div>
        <div v-show="HTTPTabs[2].visible||TCPTabs[0].visible" style="width: 100%;height: 100%">
          <VueText ref="RawText" :height="BodyRectHeight" :glyphMargin="false" :readOnly="readOnly" Language="'html'"
                   Name="Text"/>
        </div>
        <div v-show="HTTPTabs[3].visible&&DisplayTCPResponse===false" style="width: 100%;height: 100%">
          <IMGView ref="IMG"></IMGView>
        </div>
        <div v-show="HTTPTabs[4].visible&&DisplayTCPResponse===false" style="width: 100%;height: 100%">
          <div class="iframe-container">
            <iframe ref="iframe" sandbox="allow-same-origin allow-forms allow-top-navigation"
                    style="width: 100%; height: 100%" :srcdoc="GetHTML"></iframe>
          </div>
        </div>
        <div v-show="HTTPTabs[5].visible||TCPTabs[1].visible" style="width: 100%;height: 100%">
          <HexView ref="HexView" :Size="HexViewSize" :readOnly="readOnly" :raw="HexViewRaw"/>
        </div>
        <div v-show="HTTPTabs[6].visible&&DisplayTCPResponse===false" style="width: 100%;height: 100%">
          <Cookies ref="Cookies" :readOnly="readOnly"/>
        </div>
        <div v-show="HTTPTabs[7].visible||TCPTabs[2].visible" style="width: 100%;height: 100%">
          <JSon ref="Json" :height="BodyRectHeight" :width="BodyRectWidth" :readOnly="readOnly"/>
        </div>
        <div v-show="TCPTabs[3].visible&&DisplayTCPResponse" style="width: 100%;height: 100%">
          <Active ref="Active" :Height="BodyRectHeight"/>
        </div>
        <div v-if="updateSocketContent"/>
      </div>
    </div>
  </div>
</template>
<script>
import List from "../Request/List.vue";
import HexView from "../Request/HexView.vue";
import {
  Base64DecodeUint8,
  CallGoDo,
  SetUint8Array,
  StrBase64Encode,
  UInt8ToHex,
  UInt8ToStr
} from "../CallbackEventsOn.js";
import Cookies from "./Cookies.vue";
import VueText from "../Request/Text.vue";
import JavaScriptEdit from "../Request/Raw.vue";
import JSon from "../Request/JSon.vue";
import IMGView from "./IMGView.vue";
import Active from "./Active.vue";

const ClassMinName = "ag-icon ag-icon-minimize ag-panel-title-bar-button-icon"
const ClassMaxName = "ag-icon ag-icon-maximize ag-panel-title-bar-button-icon"
export default {
  props: ['Theology', 'RequestWay'],
  watch: {
    Theology(value) {
      this.RequestTheology = value

    },
    RequestWay(value) {
      this.DisplayTCPResponse = value.DisplayTCPResponse
    },
  },
  components: {
    Active,
    IMGView,
    JSon,
    JavaScriptEdit,
    Cookies,
    HexView,
    List,
    VueText
  },
  computed: {
    GetHTML() {
      return this.HTML
    },
    UpdateIcon() {
      return this.iconClass
    },
    updateSocketContent() {
      const SelectedLine = window.Socket.Line
      if (SelectedLine === null || SelectedLine === void 0) {
        if (this.$refs.Json !== null && this.$refs.Json !== void 0) {
          this.$refs.Json.SetCode(null)
        }
        if (this.$refs.RawText !== null && this.$refs.RawText !== void 0) {
          this.$refs.RawText.SetCode("")
        }
        this.HexViewRaw = SetUint8Array("", "")
        return false
      }
      const Index = SelectedLine.data["#"]
      window.Socket.Data = {
        Theology: window.Theology,
        Index: Index,
      }
      this.$nextTick(() => {
        CallGoDo("socket请求获取", window.Socket.Data).then(response => {
          if (response === null) {
            this.$nextTick(() => {
              this.$refs.Json.SetCode(null)
              this.$refs.RawText.SetCode("")
              this.$refs.Active.SetCode("")
              this.HexViewRaw = SetUint8Array("", "")
            })
            return
          }

          let Body = Base64DecodeUint8(response)
          let _Body = UInt8ToStr(Body, "utf-8")
          if (_Body.indexOf("�") !== -1) {
            _Body = UInt8ToStr(Body, "gbk")
          }
          let language = "plaintext"
          {
            try {
              const json = JSON.parse(_Body);
              if (typeof json === 'object' && json !== null) {
                language = "json"
              }
            } catch (error) {
            }
          }
          this.$nextTick(() => {
            this.$refs.RawText.SetLanguage(language)
            this.HexViewRaw = Body
            this.$refs.Active.SetCode(UInt8ToHex(Body, true))
            this.$refs.Json.SetCode(Body)
            this.$refs.RawText.SetCode(_Body)
          })

        })
      })
      return false
    },
  },
  data() {
    return {
      iconClass: ClassMaxName,
      HexViewSize: {w: 0, h: 0},
      BodyRectHeight: "0px",
      BodyRectWidth: "0px",
      readOnly: true,
      HexViewRaw: null,
      BodyUTF8: true,
      HTML: "",
      DisplayTCPResponse: true,
      RequestTheology: -1,
      HTTPTabs: [
        {id: 0, name: "原始响应", class: "ag-tab", visible: false, Show: true},
        {id: 1, name: "协议头", class: "ag-tab", visible: false, Show: true},
        {id: 2, name: "响应文本", class: "ag-tab", visible: false, Show: true},
        {id: 3, name: "图片视图", class: "ag-tab ag-tab-selected", visible: true, Show: true},
        {id: 4, name: "HTML视图", class: "ag-tab", visible: false, Show: true},
        {id: 5, name: "十六进制视图", class: "ag-tab", visible: false, Show: true},
        {id: 6, name: "Cookies", class: "ag-tab", visible: false, Show: true},
        {id: 7, name: "JSON视图", class: "ag-tab", visible: false, Show: true},
      ],
      TCPTabs: [
        {id: 0, name: "文本视图", class: "ag-tab", visible: false, Show: true},
        {id: 1, name: "十六进制视图", class: "ag-tab", visible: false, Show: true},
        {id: 2, name: "JSON视图", class: "ag-tab", visible: false, Show: true},
        {id: 3, name: "主动发送", class: "ag-tab ag-tab-selected", visible: true, Show: true},
      ],
      FinallyClick: {
        HTTPTabs: null,
        TCPTabs: null
      }
    };
  },
  methods: {
    SetMaximize(mode) {
      if (this.iconClass === ClassMaxName) {
        this.iconClass = ClassMinName
      } else {
        this.iconClass = ClassMaxName
      }
      window.vm.Tabs.ToolPanel.SetMaximize(mode)
    },
    async saveData() {
      if (this.readOnly) {
        return false
      }
      for (let i = 0; i < this.HTTPTabs.length; i++) {
        const tabs = this.HTTPTabs[i]
        if (tabs.visible) {
          switch (tabs.name) {
            case "原始响应":
              if (this.$refs.Raw.HasModify()) {
                let data = StrBase64Encode(this.$refs.Raw.GetCode())
                await CallGoDo("保存修改数据", {
                  Theology: this.RequestTheology,
                  Type: "Response",
                  Tabs: "Raw",
                  Data: data,
                  UTF8: this.BodyUTF8,
                })
                this.$refs.Raw.IsHasModify = false
                return true
              }
              break
            case "协议头":
              if (this.$refs.Headers.IsHasModify) {
                await CallGoDo("保存修改数据", {
                  Theology: this.RequestTheology,
                  Type: "Response",
                  Tabs: "Headers",
                  Data: this.$refs.Headers.RowData,
                  UTF8: this.BodyUTF8,
                })
                this.$refs.Headers.IsHasModify = false
                return true
              }
              break
            case "响应文本":
              if (this.$refs.RawText.HasModify()) {
                let mCode = StrBase64Encode(this.$refs.RawText.GetCode())
                await CallGoDo("保存修改数据", {
                  Theology: this.RequestTheology,
                  Type: "Response",
                  Tabs: "Body",
                  UTF8: this.BodyUTF8,
                  Data: mCode
                })
                this.$refs.RawText.IsHasModify = false
                return true
              }
              break
            case "Cookies":
              if (this.$refs.Cookies.IsHasModify) {
                await CallGoDo("保存修改数据", {
                  Theology: this.RequestTheology,
                  Type: "Response",
                  Tabs: "Cookies",
                  UTF8: this.BodyUTF8,
                  Data: this.$refs.Cookies.RowData
                })
                this.$refs.Cookies.IsHasModify = false
                return true
              }
              break
            case "十六进制视图":
              if (this.$refs.HexView.IsHasModify) {
                let data = this.$refs.HexView.CopyALL()
                await CallGoDo("保存修改数据", {
                  Theology: this.RequestTheology,
                  Type: "Response",
                  Tabs: "Hex",
                  Data: data
                })
                this.$refs.HexView.IsHasModify = false
                return true
              }
            case "JSON视图":
              let data = this.$refs.Json.GetCode()
              if (this.$refs.Json.IsHasModify()) {
                this.$refs.Json.jsonData = null
                await CallGoDo("保存修改数据", {
                  Theology: this.RequestTheology,
                  Type: "Response",
                  Tabs: "Json",
                  Data: StrBase64Encode(data),
                  UTF8: this.BodyUTF8,
                })
                return true
              }
              return true
          }
        }
      }
      return false
    },
    SetHTTPPagesShow(name, value) {
      for (let i = 0; i < this.HTTPTabs.length; i++) {
        if (this.HTTPTabs[i].name === name) {
          this.HTTPTabs[i].Show = value
        }
      }
    },
    UpdateData(response, Break) {
      if (window.vm.List.agSelectedLine === null) {
        return
      }
      this.readOnly = Break !== 2
      this.$refs.Headers.Empty()
      this.$refs.Cookies.Empty()
      this.SetHTTPPagesShow("Cookies", false)
      let language = "plaintext"
      let Body = Base64DecodeUint8(response.Body)
      let HexRaw = SetUint8Array("HTTP/1.1 " + response.StateCode + " " + response.StateText + "\r\n", "")
      let IsImage = false
      let IsHtml = false
      for (const key in response.Header) {
        const obj = response.Header[key]
        let value = ""
        if (obj.length > 0) {
          for (let i = 0; i < obj.length; i++) {
            HexRaw = SetUint8Array(HexRaw, key + ": " + obj[i] + "\r\n")
            this.$refs.Headers.AddLine(key, obj[i])
            value = obj[i]
            if (key.toUpperCase() === "SET-COOKIE") {
              this.SetHTTPPagesShow("Cookies", this.parsingCookie(value))
            } else if (key.toUpperCase() === 'CONTENT-TYPE') {
              IsImage = value.indexOf("image/") !== -1
              IsHtml = value.indexOf("html") !== -1
              const ar = (value + "/").replaceAll(";", "/").split("/")
              if (ar.length >= 2) {
                language = ar[1].toLowerCase()
              }
            }
          }
        } else {
          HexRaw = SetUint8Array(HexRaw, key + ": " + value + "\r\n")
          this.$refs.Headers.AddLine(key, value)
        }
      }
      let _hex = SetUint8Array(HexRaw, "\r\n")
      _hex = SetUint8Array(_hex, Body)
      this.HexViewRaw = _hex
      HexRaw = SetUint8Array(HexRaw, "\r\n")
      HexRaw = SetUint8Array(HexRaw, Body)
      let _HexRaw = UInt8ToStr(HexRaw, "utf-8")
      let _Body = UInt8ToStr(Body, "utf-8")
      this.BodyUTF8 = true
      if (_HexRaw.indexOf("�") !== -1) {
        _Body = UInt8ToStr(Body, "gbk")
        _HexRaw = UInt8ToStr(HexRaw, "gbk")
        this.BodyUTF8 = false
      }
      this.$refs.RawText.SetCode(_Body)
      this.$refs.RawText.SetLanguage(language)
      this.$refs.Raw.SetCode(_HexRaw)
      if (IsImage) {
        this.$refs.IMG.SetImg(response.Body, language)
        this.SetHTTPPagesShow("图片视图", true)
      } else {
        this.SetHTTPPagesShow("图片视图", false)
      }
      if (IsHtml) {
        this.HTML = _Body
        this.SetHTTPPagesShow("HTML视图", true)
      } else {
        this.SetHTTPPagesShow("HTML视图", false)
      }
      this.$refs.Json.SetReadOnly(this.readOnly)
      this.$refs.Json.SetCode(Body)
      this.SetHTTPPagesShow("JSON视图", true)

      this.$nextTick(() => {
        this.$refs.Headers.SelectedLine(0)
        this.$refs.Cookies.SelectedLine(0)
        this.SelectHTTPFolder()

        this.$refs.Raw.SetReadOnly(this.readOnly)
        this.$refs.RawText.SetReadOnly(this.readOnly)
      });
    },
    parsingCookie(Cookie) {
      let array1 = Cookie.split(";")
      if (array1.length > 1) {
        const array2 = array1[0].trim().split("=")
        array1.shift();
        if (array2.length === 1) {
          if (array2[0] !== "") {
            this.$refs.Cookies.AddLine(array2[0], "", "")
          }
        } else if (array2.length >= 2) {
          this.$refs.Cookies.AddLine(array2[0], array2[1], array1.join(";").trim())
        }
      }
      return this.$refs.Cookies.GetLine() !== 0
    },
    eHeaderClick(eve) {
      if (eve === null || eve === void 0) {
        eve = this.HTTPTabs[0]
      }
      this.FinallyClick.HTTPTabs = eve
      this.saveData()
      for (let i = 0; i < this.TCPTabs.length; i++) {
        this.TCPTabs[i].class = "ag-tab"
        this.TCPTabs[i].visible = false
      }
      for (let i = 0; i < this.HTTPTabs.length; i++) {
        if (this.HTTPTabs[i].id !== eve.id) {
          this.HTTPTabs[i].class = "ag-tab"
          this.HTTPTabs[i].visible = false
        } else {
          this.HTTPTabs[i].class = "ag-tab ag-tab-selected"
          this.HTTPTabs[i].visible = true
          if (this.HTTPTabs[i].name === "十六进制视图") {
            this.$nextTick(() => {
              this.$refs.HexView.Refresh()
            })
          }
        }
      }

    },
    eTcpClick(eve) {
      this.$refs.Active.IsWs = window.vm.Tabs.Request.DisplayHTTPHeader
      if (eve === null || eve === void 0) {
        eve = this.TCPTabs[0]
      }
      this.FinallyClick.TCPTabs = eve
      for (let i = 0; i < this.HTTPTabs.length; i++) {
        this.HTTPTabs[i].class = "ag-tab"
        this.HTTPTabs[i].visible = false
      }
      for (let i = 0; i < this.TCPTabs.length; i++) {
        if (this.TCPTabs[i].id !== eve.id) {
          this.TCPTabs[i].class = "ag-tab"
          this.TCPTabs[i].visible = false
        } else {
          this.TCPTabs[i].class = "ag-tab ag-tab-selected"
          this.TCPTabs[i].visible = true
          if (this.TCPTabs[i].name === "十六进制视图") {
            this.$nextTick(() => {
              this.$refs.HexView.Refresh()
            })
          }
        }
      }
    },
    SelectHTTPFolder() {
      this.$nextTick(() => {
        this.DisplayTCPResponse = window.vm.Tabs.Request.IsTCPorWebsocket() && window.vm.Tabs.Request.IsSelectedWebsocket()
        if (this.DisplayTCPResponse) {
          this.eTcpClick(this.FinallyClick.TCPTabs)
        }
      })
      let isTCP = window.vm.List.agSelectedLine.data["方式"].toUpperCase().indexOf("TCP") !== -1
      let isUDP = window.vm.List.agSelectedLine.data["方式"].toUpperCase().indexOf("UDP") !== -1
      if (isTCP || isUDP) {

        let nc = 0
        for (let i = 0; i < this.TCPTabs.length; i++) {
          if (!this.TCPTabs[i].visible) {
            nc++
          }
        }
        if (nc === this.TCPTabs.length) {
          this.eTcpClick(this.TCPTabs[0])
          return
        }
      } else {
        let nc = 0
        for (let i = 0; i < this.HTTPTabs.length; i++) {
          if (!this.HTTPTabs[i].visible) {
            nc++
          }
        }
        if (nc === this.HTTPTabs.length) {
          this.eHeaderClick(this.HTTPTabs[0])
          return
        }
      }
      for (let i = 0; i < this.HTTPTabs.length; i++) {
        if (this.HTTPTabs[i].visible && this.HTTPTabs[i].Show === false) {
          for (let n = 0; n < this.HTTPTabs.length; n++) {
            this.HTTPTabs[n].class = "ag-tab"
            this.HTTPTabs[n].visible = false
          }
          this.eHeaderClick(this.HTTPTabs[0])
          return
        }
      }
    }
  },
  mounted() {
    window.vm.Tabs.Response = this
    {
      const elementRef = this.$refs.BodyRect; // 获取元素的引用
      // 创建 ResizeObserver 实例并监听元素尺寸变化
      const resizeObserver = new ResizeObserver(entries => {
        for (const entry of entries) {
          const {width, height} = entry.contentRect;
          this.HexViewSize = {w: width, h: height}
          this.BodyRectHeight = height + "px"
          this.BodyRectWidth = width + "px"
        }
      });
      resizeObserver.observe(elementRef); // 开始监听元素尺寸变化
    }
    //禁止点击 iframe 内的所有元素
    {
      const iframe = this.$refs.iframe;
      // 禁止点击事件
      iframe.addEventListener('load', function () {
        const iframeDocument = iframe.contentWindow.document;
        iframeDocument.addEventListener('click', function (event) {
          event.stopPropagation();
          event.preventDefault();
        }, true);
      });
    }
  }
}
</script>
<style>
.iframe-container {
  position: relative;
  width: 100%;
  height: 100%;
}

.iframe-container iframe {
  pointer-events: auto;
  width: 100%;
  height: 100%;
}
</style>