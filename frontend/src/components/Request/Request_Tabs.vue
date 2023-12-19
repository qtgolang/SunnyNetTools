<template>
  <div class="ag-tabs ag-chart-tabbed-menu ag-focus-managed">
    <div role="tablist" class="ag-tabs-header ag-chart-tabbed-menu-header">
      <span v-if="DisplayHTTPHeader" v-show="HTTPTabs[0].Show" @click="eHeaderClick(HTTPTabs[0])"
            :class="HTTPTabs[0].class"
            role="tab"> {{ HTTPTabs[0].name }} </span>
      <span v-if="DisplayHTTPHeader" v-show="HTTPTabs[1].Show" @click="eHeaderClick(HTTPTabs[1])"
            :class="HTTPTabs[1].class"
            role="tab"> {{ HTTPTabs[1].name }} </span>
      <span v-if="DisplayHTTPHeader" v-show="HTTPTabs[2].Show" @click="eHeaderClick(HTTPTabs[2])"
            :class="HTTPTabs[2].class"
            role="tab"> {{ HTTPTabs[2].name }} </span>
      <span v-if="DisplayHTTPHeader" v-show="HTTPTabs[3].Show.Text" @click="eHeaderClick(HTTPTabs[3])"
            :class="HTTPTabs[3].class"
            role="tab"> {{ HTTPTabs[3].name }} </span>
      <span v-if="DisplayHTTPHeader" v-show="HTTPTabs[4].Show" @click="eHeaderClick(HTTPTabs[4])"
            :class="HTTPTabs[4].class"
            role="tab"> {{ HTTPTabs[4].name }} </span>
      <span v-if="DisplayHTTPHeader" v-show="HTTPTabs[5].Show" @click="eHeaderClick(HTTPTabs[5])"
            :class="HTTPTabs[5].class"
            role="tab"> {{ HTTPTabs[5].name }} </span>
      <span v-if="DisplayHTTPHeader" v-show="HTTPTabs[6].Show" @click="eHeaderClick(HTTPTabs[6])"
            :class="HTTPTabs[6].class"
            role="tab"> {{ HTTPTabs[6].name }} </span>
      <span v-if="DisplayHTTPHeader" v-show="HTTPTabs[7].Show" @click="eHeaderClick(HTTPTabs[7])"
            :class="HTTPTabs[7].class"
            role="tab"> {{ HTTPTabs[7].name }} </span>
      <div v-if="DisplayHTTPHeader===false" v-for="item in TCPTabs">
        <span role="tab" :class="item.class" aria-label="{{ item.name }}" @click="eHeaderClick(item)">
            <div>{{ item.name }}</div>
          </span>
      </div>
    </div>
    <div role="presentation" ref="BodyRect" class="ag-tabs-body ag-chart-tabbed-menu-body" style="">
      <div :style="{  width: '100%',  height: RawBodyRectHeight + 'px',  position: 'absolute'}">
        <div v-show="HTTPTabs[0].visible&&DisplayHTTPHeader" :style="GetBreakpointStateStyle">
          <JavaScriptEdit ref="Raw" :height="BodyRectHeight" :glyphMargin="false" :readOnly="readOnly"
                          :Text="RawStrings"
                          Name="Request"/>
        </div>
        <div v-show="HTTPTabs[1].visible&&DisplayHTTPHeader" :style="GetBreakpointStateStyle">
          <List ref="Headers" :readOnly="readOnly"/>
        </div>
        <div v-show="HTTPTabs[2].visible&&DisplayHTTPHeader" :style="GetBreakpointStateStyle">
          <List ref="Args" :readOnly="readOnly"/>
        </div>
        <div v-show="HTTPTabs[3].visible&&DisplayHTTPHeader" :style="GetBreakpointStateStyle">
          <div v-show="HTTPTabs[3].Show.BodyArgs"
               style="width: 100%; height: 30px; display: flex; justify-content: flex-start; align-items: center;">
            &nbsp;
            <div v-show="HTTPTabs[3].Show.Index===0">
              <el-button size="small" disabled>Raw</el-button>
              <el-button size="small" @click="TextTypeClick">UrlEncode</el-button>
            </div>
            <div v-show="HTTPTabs[3].Show.Index===1">
              <el-button size="small" @click="TextTypeClick">Raw</el-button>
              <el-button size="small" disabled>UrlEncode</el-button>
            </div>
          </div>
          <div v-show="HTTPTabs[3].Show.Index===0" style="width: 100%;height: 100%">
            <VueText ref="RawText" :height="GetTestRawHeight" :glyphMargin="false"
                     :readOnly="readOnly" Language="'html'"
                     Name="Text"/>
          </div>
          <div v-show="HTTPTabs[3].Show.Index===1" style="width: 100%;height: calc(100% - 30px)">
            <List ref="BodyArgs" :readOnly="readOnly"/>
          </div>
        </div>
        <div v-show="HTTPTabs[4].visible&&DisplayHTTPHeader" :style="GetBreakpointStateStyle">
          <HexView ref="HexView" :Size="HexViewSize" :readOnly="readOnly" :raw="HexViewRaw"/>
        </div>
        <div v-show="HTTPTabs[5].visible&&DisplayHTTPHeader" :style="GetBreakpointStateStyle">
          <List ref="Cookies" :readOnly="readOnly"/>
        </div>
        <div v-show="HTTPTabs[6].visible&&DisplayHTTPHeader" :style="GetBreakpointStateStyle">
          <JSon ref="Json" :height="BodyRectHeight" :width="BodyRectWidth" :readOnly="readOnly"/>
        </div>
        <div v-show="HTTPTabs[7].visible||DisplayHTTPHeader===false||DisplayHTTPHeaderWebSocket"
             :style="GetBreakpointStateStyle">
          <TcpView ref="WebSocket" :readOnly="readOnly" :width="BodyRectWidth"/>
        </div>
        <div style="width: 100%; height: 30px; display: flex; justify-content: center; align-items: center;">
          <div
              style="margin: 0 10px; background-color: #f64141; height: 28px; line-height: 28px; width: 110px; border-radius: 5px;">
            &nbsp;断点命中,然后->&nbsp;
          </div>
          <div v-show="BreakResponse===false">
            <el-button type="warning" style="height: 28px" @click="BreakClick(2)">中断响应</el-button>&nbsp;&nbsp;&nbsp;
          </div>
          <el-button type="success" style="height: 28px" @click="BreakClick(0)">运行到结束</el-button>
        </div>
      </div>

    </div>
  </div>
</template>
<script>
import JavaScriptEdit from "./Raw.vue";
import {
  Base64DecodeUint8,
  CallGoDo,
  HexToGbkStr,
  SetUint8Array,
  StrBase64Encode,
  UInt8ToStr
} from "../CallbackEventsOn.js";
import VueText from "./Text.vue";
import List from "./List.vue";
import JSon from "./JSon.vue";
import HexView from "./HexView.vue";
import TcpView from "./WebSocketView.vue";

export default {
  props: ['Theology', 'RequestWay'],
  watch: {
    Theology(value) {
      this.RequestTheology = value

    },
    RequestWay(value) {
      this.DisplayHTTPHeaderWebSocket = value.DisplayHTTPHeaderWebSocket
      this.DisplayHTTPHeader = value.DisplayHTTPHeader
    },
  },
  computed: {
    GetBreakpointStateStyle() {
      if (this.Breakpoint) {
        return "width: 100%;height: calc(100% - 30px)"
      }
      return "width: 100%;height: 100%"
    },
    GetTestRawHeight() {
      return this.HTTPTabs[3].Show.BodyArgs ? this.BodyRectHeight2 : this.BodyRectHeight
    }
  },
  components: {
    TcpView,
    HexView,
    JSon,
    List,
    VueText,
    JavaScriptEdit,
  },
  data() {
    return {
      HexViewSize: {w: 0, h: 0},
      //是否中断响应,如果是只能选择运行到结束
      BreakResponse: false,
      //是否显示拦截器,显示中断到响应,和运行到结束
      Breakpoint: false,
      //响应原始信息,
      RawInfo: null,
      BodyUTF8: true,
      BodyRectHeight: "0px",
      BodyRectWidth: "0px",
      BodyRectHeight2: "0px",
      RawBodyRectHeight: 0,
      RawBodyRectHeight2: 0,
      HTTPTabs: [
        {id: 0, name: "原始视图", class: "ag-tab ag-tab-selected", visible: true, Show: true},
        {id: 1, name: "协议头", class: "ag-tab", visible: false, Show: true},
        {id: 2, name: "参数视图", class: "ag-tab", visible: false, Show: true},
        {
          id: 3, name: "请求数据", class: "ag-tab", visible: false, Show: {
            BodyArgs: true,
            Text: true,
            Index: 0
          },
        },
        {id: 4, name: "十六进制视图", class: "ag-tab", visible: false, Show: true},
        {id: 5, name: "Cookies", class: "ag-tab", visible: false, Show: true},
        {id: 6, name: "JSON视图", class: "ag-tab", visible: false, Show: true},
        {id: 7, name: "WebSocket", class: "ag-tab", visible: false, Show: true}
      ],
      TCPTabs: [
        {id: 0, name: "TCP请求数据流", class: "ag-tab ag-tab-selected", visible: false, Show: true},
      ],
      RequestTheology: -1,
      readOnly: true,
      RawStrings: "",
      Body: {
        URL: ""
      },
      DisplayHTTPHeaderWebSocket: false,
      DisplayHTTPHeader: true,
      HexViewRaw: null,
    };
  },
  methods: {
    BreakClick(mode) {
      const cmm = window.vm.List.agSelectedLine.data["断点模式"]
      console.log(cmm, mode)
      if (cmm === 0) {
        window.vm.List.agSelectedLine.data["断点模式"] = 0
        return;
      }
      if (cmm === 1) {
        this.saveData().then(res => {
          CallGoDo("断点点击", {NextBreak: mode, Theology: this.RequestTheology})


          this.Breakpoint = false
          this.GetBodyRect(this.$refs.BodyRect.offsetHeight, this.$refs.BodyRect.offsetWidth)
          window.vm.List.RefreshRenderedNodes()
          //window.vm.List.agGridApi.setRowData(window.vm.List.RowData);


        })
        return
      }
      window.vm.Tabs.Response.saveData().then(res => {
        CallGoDo("断点点击", {NextBreak: mode, Theology: this.RequestTheology})

        this.Breakpoint = false
        this.GetBodyRect(this.$refs.BodyRect.offsetHeight, this.$refs.BodyRect.offsetWidth)
        window.vm.List.RefreshRenderedNodes()
        //window.vm.List.agGridApi.setRowData(window.vm.List.RowData);


      })
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
      this.SetHTTPPagesShow("WebSocket", window.vm.List.agSelectedLine.data["响应类型"].toUpperCase() === "WEBSOCKET")
      this.readOnly = Break !== 1
      {
        if (Break === 1 || Break === 2) {
          this.Breakpoint = true
          if (Break === 2) {
            this.BreakResponse = true
          } else {
            this.BreakResponse = false
          }
        } else {
          this.Breakpoint = false
        }
        this.GetBodyRect(this.$refs.BodyRect.offsetHeight, this.$refs.BodyRect.offsetWidth)
      }
      this.RawInfo = response
      this.$refs.Headers.Empty()
      this.$refs.Cookies.Empty()
      this.$refs.WebSocket.Empty()
      this.$refs.Args.Empty()
      this.$refs.WebSocket.MenuItems[2].selected = response.Options.StopSend
      this.$refs.WebSocket.MenuItems[3].selected = response.Options.StopRec
      this.$refs.WebSocket.MenuItems[4].selected = response.Options.StopALL
      this.SetHTTPPagesShow("Cookies", false)
      let language = "plaintext"
      let Body = Base64DecodeUint8(response.Body)
      let HexRaw = SetUint8Array(response.Method + " " + response.URL + " " + response.Proto + "\r\n", "")
      for (const key in response.Header) {
        HexRaw = SetUint8Array(HexRaw, key + ": ")
        const obj = response.Header[key]
        let value = ""
        if (obj.length > 0) {
          value = obj[0]
        }
        HexRaw = SetUint8Array(HexRaw, value + "\r\n")
        this.$refs.Headers.AddLine(key, value)
        if (key.toUpperCase() === "COOKIE") {
          this.SetHTTPPagesShow("Cookies", this.parsingCookie(value))
        } else if (key.toUpperCase() === 'CONTENT-TYPE') {
          const ar = (value + "/").replaceAll(";", "/").split("/")
          if (ar.length >= 2) {
            language = ar[1].toLowerCase()
          }
        }
      }
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
      //URL参数列表
      {
        let count = 0
        let array = response.URL.split("?")
        if (array.length === 2) {
          array = array[1].split("&")
          for (let i = 0; i < array.length; i++) {
            const arr = array[i].split("=")
            if (arr.length === 1) {
              if (arr[0] !== "") {
                this.$refs.Args.AddLine(arr[0], "", true)
                count++
              }
            } else if (arr.length === 2) {
              let _value = ""
              let _utf8 = true
              try {
                _value = decodeURIComponent(arr[1])
              } catch (e) {
                _value = HexToGbkStr(arr[1])
                _utf8 = false
              }
              this.$refs.Args.AddLine(arr[0], _value, _utf8)
              count++
            }
          }
        }
        this.SetHTTPPagesShow("参数视图", count !== 0)
      }
      //原始文本内容 ->JSON视图
      {
        this.HTTPTabs[3].Show.BodyArgs = false
        this.HTTPTabs[3].Show.Index = 0
        try {
          const json = JSON.parse(_Body);
          if (typeof json === 'object' && json !== null) {
            this.$refs.RawText.SetLanguage("json")
            language = "json"
            //JSON视图
            {
              this.$refs.Json.SetReadOnly(this.readOnly)
              this.$refs.Json.SetCode(Body)
            }
          } else {
            this.$refs.RawText.SetLanguage(language)
          }
        } catch (error) {
          this.$refs.RawText.SetLanguage(language)
        }
        this.$refs.RawText.SetCode(_Body)
        if (_Body === "") {
          this.HTTPTabs[3].Show.Text = false
          this.SetHTTPPagesShow("JSON视图", false)
        } else if (language !== "json") {
          this.HTTPTabs[3].Show.Text = true
          this.SetHTTPPagesShow("JSON视图", false)
          this.$refs.BodyArgs.Empty()
          const array = _Body.split("&")
          let ok = false
          for (let i = 0; i < array.length; i++) {
            const array2 = array[i].split("=")
            if (array2.length > 2) {
              ok = false
              break
            }
            if (array2.length === 1) {
              this.$refs.BodyArgs.AddLine(array2[0], "", true)
            } else {
              let _value = ""
              let _utf8 = true
              try {
                _value = decodeURIComponent(array2[1])
              } catch (e) {
                _value = HexToGbkStr(array2[1])
                _utf8 = false
              }
              this.$refs.BodyArgs.AddLine(array2[0], _value, _utf8)
            }
            ok = true
          }
          if (ok) {
            this.HTTPTabs[3].Show.BodyArgs = true
            this.HTTPTabs[3].Show.Index = 1
          } else {
            this.HTTPTabs[3].Show.BodyArgs = false
            this.HTTPTabs[3].Show.Index = 0
          }
        } else {
          this.HTTPTabs[3].Show.Text = true
          this.HTTPTabs[3].Show.BodyArgs = false
          this.HTTPTabs[3].Show.Index = 0
          this.SetHTTPPagesShow("JSON视图", true)
        }
      }
      //Hex视图
      {
        this.HexViewRaw = HexRaw
        this.$refs.HexView.IsHasModify = false
        this.$refs.Raw.SetCode(_HexRaw)
      }
      {
        //WebSocket数据
        this.$nextTick(() => {
          this.$refs.WebSocket.AddLines(response.SocketData)
          this.$refs.WebSocket.Refresh()
        });
      }
      this.$nextTick(() => {
        this.$refs.Headers.SelectedLine(0)
        this.$refs.Cookies.SelectedLine(0)
        this.$refs.Args.SelectedLine(0)
        this.$refs.BodyArgs.SelectedLine(0)
        this.SelectHTTPFolder()

        this.$refs.Raw.SetReadOnly(this.readOnly)
        this.$refs.RawText.SetReadOnly(this.readOnly)
      });

    },
    parsingCookie(Cookie) {
      let count = 0
      const array1 = Cookie.split(";")
      for (let i = 0; i < array1.length; i++) {
        const array2 = array1[i].trim().split("=")
        if (array2.length === 1) {
          if (array2[0] !== "") {
            this.$refs.Cookies.AddLine(array2[0], "")
            count++
          }
        } else if (array2.length === 2) {
          this.$refs.Cookies.AddLine(array2[0], array2[1])
          count++
        }
      }
      return count !== 0
    },
    async saveData() {
      if (this.readOnly) {
        return false
      }
      for (let i = 0; i < this.HTTPTabs.length; i++) {
        const tabs = this.HTTPTabs[i]
        if (tabs.visible) {
          switch (tabs.name) {
            case "原始视图":
              if (this.$refs.Raw.HasModify()) {
                let data = StrBase64Encode(this.$refs.Raw.GetCode())
                await CallGoDo("保存修改数据", {
                  Theology: this.RequestTheology,
                  Type: "Request",
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
                  Type: "Request",
                  Tabs: "Headers",
                  Data: this.$refs.Headers.RowData,
                  UTF8: this.BodyUTF8,
                })
                this.$refs.Headers.IsHasModify = false
                return true
              }
              break
            case "参数视图":
              if (this.$refs.Args.IsHasModify) {
                await CallGoDo("保存修改数据", {
                  Theology: this.RequestTheology,
                  Type: "Request",
                  Tabs: "Args",
                  Data: this.$refs.Args.RowData,
                  UTF8: this.BodyUTF8,
                })
                this.$refs.Args.IsHasModify = false
                return true
              }
              break
            case "请求数据":
              if (this.HTTPTabs[3].Show.BodyArgs && this.HTTPTabs[3].Show.Index === 1) {
                if (this.$refs.BodyArgs.IsHasModify) {
                  this.$refs.BodyArgs.IsHasModify = false
                  await CallGoDo("保存修改数据", {
                    Theology: this.RequestTheology,
                    Type: "Request",
                    Tabs: "BodyArgs",
                    UTF8: this.BodyUTF8,
                    Data: this.$refs.BodyArgs.RowData
                  }).then(res => {
                    this.$refs.RawText.SetCode(res)
                  })
                  return true
                }
              } else if (this.$refs.RawText.HasModify()) {
                let mCode = this.$refs.RawText.GetCode()
                const array = mCode.split("&")
                this.$refs.BodyArgs.Empty()
                for (let i = 0; i < array.length; i++) {
                  const array2 = array[i].split("=")
                  if (array2.length > 2) {
                    break
                  }
                  if (array2.length === 1) {
                    this.$refs.BodyArgs.AddLine(array2[0], "", true)
                  } else {
                    let _value = ""
                    let _utf8 = true
                    try {
                      _value = decodeURIComponent(array2[1])
                    } catch (e) {
                      _value = HexToGbkStr(array2[1])
                      _utf8 = false

                    }
                    this.$refs.BodyArgs.AddLine(array2[0], _value, _utf8)
                  }
                }
                await CallGoDo("保存修改数据", {
                  Theology: this.RequestTheology,
                  Type: "Request",
                  Tabs: "Body",
                  UTF8: this.BodyUTF8,
                  Data: StrBase64Encode(mCode)
                })
                this.$refs.RawText.IsHasModify = false
                return true
              }
              break
            case "Cookies":
              if (this.$refs.Cookies.IsHasModify) {
                await CallGoDo("保存修改数据", {
                  Theology: this.RequestTheology,
                  Type: "Request",
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
                  Type: "Request",
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
                  Type: "Request",
                  Tabs: "Json",
                  Data: StrBase64Encode(data),
                  UTF8: this.BodyUTF8,
                })
                return true
              }
              break
          }
        }
      }
      return false
    },
    eHeaderClick(eve) {
      this.saveData().then(res => {
        if (res) {
          CallGoDo("HTTP请求获取", {Theology: window.Theology}).then(response => {
            window.vm.Tabs.ToolPanel.Method = response.Method
            const Break = window.vm.List.agSelectedLine.data['断点模式']
            this.UpdateData(response, Break)
            for (let i = 0; i < this.HTTPTabs.length; i++) {
              if (this.HTTPTabs[i].id !== eve.id) {
                this.HTTPTabs[i].class = "ag-tab"
                this.HTTPTabs[i].visible = false
              } else {
                this.HTTPTabs[i].class = "ag-tab ag-tab-selected"
                this.HTTPTabs[i].visible = true
              }
            }
          })
          return
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
            window.vm.Tabs.ToolPanel.DisplayTCPResponse = this.HTTPTabs[i].name === "WebSocket";
            if (window.vm.Tabs.ToolPanel.DisplayTCPResponse) {
              window.vm.Tabs.Response.DisplayTCPResponse = this.IsTCPorWebsocket() && this.IsSelectedWebsocket()
              window.vm.Tabs.Response.eTcpClick(window.vm.Tabs.Response.FinallyClick.TCPTabs)
            } else {
              window.vm.Tabs.Response.eHeaderClick(window.vm.Tabs.Response.FinallyClick.HTTPTabs)
            }
          }
        }
      })
    },
    IsTCPorWebsocket() {
      if (this.DisplayHTTPHeaderWebSocket) {
        return true
      }
      return this.DisplayHTTPHeader === false
    },
    IsSelectedWebsocket() {
      return this.HTTPTabs[7].visible || this.DisplayHTTPHeader === false
    },
    TextTypeClick() {
      this.saveData().then(res => {
        if (this.HTTPTabs[3].Show.Index === 0) {
          this.HTTPTabs[3].Show.Index = 1
          return
        }
        this.HTTPTabs[3].Show.Index = 0
      })
    },
    GetBodyRect(height, width) {
      this.RawBodyRectHeight = height
      if (this.Breakpoint) {
        this.BodyRectHeight = (height - 30) + "px"
        this.BodyRectHeight2 = (height - 60) + "px"
        this.HexViewSize = {w: width, h: height - 30}
      } else {
        this.HexViewSize = {w: width, h: height}
        this.BodyRectHeight = height + "px"
        this.BodyRectHeight2 = (height - 30) + "px"
      }
      this.BodyRectWidth = width + "px"
      this.RawBodyRectHeight2 = this.BodyRectHeight2
    },
    SelectHTTPFolder() {
      for (let i = 0; i < this.HTTPTabs.length; i++) {
        if (this.HTTPTabs[i].visible && this.HTTPTabs[i].Show === false) {
          for (let n = 0; n < this.HTTPTabs.length; n++) {
            this.HTTPTabs[n].class = "ag-tab"
            this.HTTPTabs[n].visible = false
          }
          this.HTTPTabs[0].class = "ag-tab ag-tab-selected"
          this.HTTPTabs[0].visible = true
          return
        }
      }
    }
  },
  mounted() {
    window.vm.Tabs.Request = this
    // this.readOnly = false
    const elementRef = this.$refs.BodyRect; // 获取元素的引用
    // 创建 ResizeObserver 实例并监听元素尺寸变化
    const resizeObserver = new ResizeObserver(entries => {
      for (const entry of entries) {
        const {width, height} = entry.contentRect;
        this.GetBodyRect(height, width)
      }
    });
    resizeObserver.observe(elementRef); // 开始监听元素尺寸变化
  }
}
</script>