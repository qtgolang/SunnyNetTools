<template>
  <div ref="splitA0" class="ag-column-panel" style="position: relative;width: 100%;height: 100%;display: block;">
    <div ref="splitA1" v-show="isShowRequest" class="ag-column-select ag-column-panel-column-select"
         style="height: 50%; flex: 0 0 auto;padding-top: calc(var(--ag-grid-size) * 2)">
      <div aria-hidden="true"
           class="ag-column-drop-title-bar ag-column-drop-vertical-title-bar"
           style="height: 30px;    top: -4px;    position: relative;z-index: 0">
        <el-tag type="success" style="height: 20px;border-radius: 10px;margin-right: 10px">{{ Method }}</el-tag>
        <span
            class="ag-column-drop-title ag-column-drop-vertical-title">消息流内容</span>
        <span :class="UpdateIcon"
              style="right: 10px;position:absolute;cursor:pointer" @click="SetMaximize(1)"></span>
      </div>

      <div class="ag-tabs-body ag-chart-tabbed-menu-body" style="">
        <div style="width: 100%;height: 100%">
          <WebSocketView ref="Stream" :height="'100%'" :width="'100%'" :Name="Method" Name="TCP_UDP"/>
        </div>

      </div>
      <div ref="split" @mousedown="handleMouseDown" class="ag-resizer ag-resizer-bottom"
           style="pointer-events: all;z-index: 4"></div>
    </div>
    <div ref="splitA3" v-show="isShowResponse" :style="RangeStyle">
      <div class="ag-column-drop-title-bar ag-column-drop-vertical-title-bar">
        <div style="width: calc(100% - 30px);display: flex">
            <span
                v-for="(tab, index) in Tabs"
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
        <span :class="UpdateIcon" style="right: 10px;position:absolute;cursor:pointer"
              @click="SetMaximize(2)"></span>
      </div>
      <div class="ag-tabs-body ag-chart-tabbed-menu-body" style="width: 100%;height: 100%;">
        <div :style="{  width: '100%',  height: 'calc(100% - 33px)',  position: 'absolute',zIndex:'1',}">
          <div v-show="Tabs[0].visible" style="width: 100%;height: 100%">
            <JavaScriptEdit ref="WebSocketText" :glyphMargin="false" Name="WebSocketText"/>
          </div>
          <div v-show="Tabs[1].visible" style="width: 100%;height: 100%">
            <HexView ref="WebSocketHEX" :glyphMargin="false" Name="WebSocketHEX"/>
          </div>
          <div v-show="Tabs[2].visible" style="width: 100%;height: 100%">
            <JSon ref="WebSocketJson" :height="'100%'" :width="'100%'"/>
          </div>
          <div v-show="Tabs[3].visible" style="width: 100%;height: 100%">
            <WebsocketActive ref="WebSocketActive" style="width: 100%;height: 100%"/>
          </div>
        </div>
      </div>
    </div>
  </div>

</template>
<script>
import {Config_SelectedRow, Config_SocketSelectedRow} from "../../../config/Config.js";
import {
  GetAllStream,
  GetSessionMessageBody,
  GetSocketFilter,
} from "../../../../../bindings/changeme/Service/appmain.js";
import Table from "../../../Tools/table.vue"
import JavaScriptEdit from "../tool/Raw.vue"
import {base64ToBytes, StringToBytes,} from "../../../config/encoding.js";
import HexView from "../../../Tools/HexView.vue";
import JSon from "../../../Tools/JSon.vue";
import WebSocketView from "../../../Tools/WebSocketView.vue";
import WebsocketActive from "../../../Tools/WebsocketActive.vue";

const ClassMinName = "ag-icon ag-icon-minimize ag-panel-title-bar-button-icon"
const ClassMaxName = "ag-icon ag-icon-maximize ag-panel-title-bar-button-icon"
export default {
  components: {WebsocketActive, WebSocketView, JSon, HexView, JavaScriptEdit, Table},
  computed: {
    isShowResponse() {
      return this.Show.showResponse
    },
    isShowRequest() {
      return this.Show.showRequest;
    },
    Method() {
      return this.SelectedRow["方式"]
    },
    UpdateIcon() {
      return this.iconClass
    },
  },
  data() {
    return {
      RangeStyle: "",
      isChangeSize: false,
      iconClass: ClassMaxName,
      get SelectedRow() {
        return Config_SelectedRow.value
      },
      set SelectedRow(value) {
        Config_SelectedRow.value = value
      },
      get SocketSelectedRow() {
        return Config_SocketSelectedRow.value
      },
      set SocketSelectedRow(value) {
        Config_SocketSelectedRow.value = value
      },
      Tabs: [
        {id: 0, name: "文本视图", class: "ag-tab", visible: false, Show: true},
        {id: 1, name: "十六进制视图", class: "ag-tab", visible: false, Show: true},
        {id: 2, name: "JSON视图", class: "ag-tab", visible: false, Show: true},
        {id: 3, name: "主动发送", class: "ag-tab ag-tab-selected", visible: true, Show: true},
      ],
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
        const r1 = parseInt(this.$refs.splitA1.clientHeight)
        this.RangeStyle = "width: 100%;height: calc(100% - " + r1 + "px)";
      })
    },
    handleMouseDown() {
      this.isChangeSize = true;
    },
    handleMouseMove(event) {
      if (!this.isChangeSize) return;

      const clientY = event.clientY;
      const docHeight = document.documentElement.clientHeight;

      // ✅ 处理顶部边界
      if (clientY + 5 < 60) {
        return;
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
    },
    eWebsocketClick(eve) {
      this.Tabs.forEach((obj) => {
        obj.class = "ag-tab"
        obj.visible = false
      })
      eve.class = "ag-tab ag-tab-selected"
      eve.visible = true
    },
  },
  watch: {
    "SelectedRow"(newValue) {
      const _SelectedTheology = parseInt(newValue?.Theology ?? "0");
      if (_SelectedTheology === 0) {
        return
      }
      if (newValue["方式"].toLowerCase().indexOf("tcp") === -1 && newValue["方式"].toLowerCase().indexOf("udp") === -1) {
        return
      }
      this.$refs.Stream.Empty()
      GetSocketFilter(_SelectedTheology).then((Model) => {
        this.$refs.Stream.$refs.Filter.setFilter(Model)
        GetAllStream(_SelectedTheology).then((res) => {
          this.$refs.Stream.InsertSocketStream(res, true)
        });
      })
    },
    "SocketSelectedRow"(newValue) {
      const _SelectedTheology = parseInt(this.SelectedRow?.Theology ?? "0");
      if (_SelectedTheology === 0) {
        return
      }
      if (this.SelectedRow["方式"].toLowerCase().indexOf("tcp") === -1 && this.SelectedRow["方式"].toLowerCase().indexOf("udp") === -1) {
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
        if (res) {
          const bs = base64ToBytes(res);
          this.$refs.WebSocketHEX.SetCode(bs)
          this.$refs.WebSocketJson.SetCode(bs)
          this.$refs.WebSocketText.SetCode(bs)
          this.$refs.WebSocketActive.SetCode(bs, IsClose)
        }
      })
    }
  },
  mounted() {
    {
      document.addEventListener('mouseup', () => {
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

  }
}
</script>