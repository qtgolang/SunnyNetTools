<template>
  <div style="width: 100%;height:calc(100% - 40px);">
    <div :style="`top: 0;left: 0;width: 100%;height: calc(100% - 40px);`">
      <VueText ref="RawText" Language="'html'" Name="ActiveView" :HideEncodingConvert="true"
               :NoRedrawingAllowed="true"/>
    </div>
    <div style="height: 40px;display: flex;width: 100%;justify-content: center; align-items: center; flex-grow: 1;">
      <el-select v-model="SendType" class="m-2" placeholder="Select" size="small" style="flex-grow: 0;width: 157px;">
        <el-option
            v-for="item in SendTypeOptions"
            :key="item.value"
            :label="item.label"
            :value="item.value"
        />
      </el-select>
      <el-select v-model="wsType" class="m-2" placeholder="Select" size="small" style="width: 122px;" v-show="IsWs">
        <el-option
            v-for="item in wsTypeOptions"
            :key="item.value"
            :label="item.label"
            :value="item.value"
        />
      </el-select>
      <el-select v-model="direction" class="m-2" placeholder="Select" size="small" style="width: 115px;">
        <el-option
            v-for="item in directionOptions"
            :key="item.value"
            :label="item.label"
            :value="item.value"
        />
      </el-select>
      <el-button size="small" style="flex-grow: 1;" @click="Send" v-show="IsSending===false">立即发送</el-button>
      <el-button size="small" style="flex-grow: 1;" disabled v-show="IsSending">正在发送</el-button>
    </div>
  </div>
</template>

<script>
import VueText from "../Home/Request/tool/Text.vue";
import {Config_SelectedRow} from "../config/Config.js";
import {GoGetHex, SessionActiveSend} from "../../../bindings/changeme/Service/appmain.js";
import {ElMessage} from "element-plus";
import {bytesToBase64, bytesToString, StringToBytes} from "../config/encoding.js";

export default {
  components: {VueText},
  computed: {
    IsWs() {
      const m = Config_SelectedRow?.value?.["方式"] ?? "";
      this.IsSending = false
      return m.toLowerCase().indexOf("web") !== -1
    }
  },
  data() {
    return {
      IsSending: false,
      SendType: "HEX",
      direction: "Server",
      wsType: "Text",
      SendTypeOptions: [
        {
          value: 'HEX',
          label: '发送类型:HEX',
        },
        {
          value: 'Base64',
          label: '发送类型:Base64',
        },
        {
          value: 'GBK',
          label: '发送类型:String(GBK)',
        },
        {
          value: 'UTF8',
          label: '发送类型:String(UTF8)',
        },
      ],
      directionOptions: [
        {
          value: 'Client',
          label: '向 客户端 发送',
        },
        {
          value: 'Server',
          label: '向 服务器 发送',
        },
      ],
      wsTypeOptions: [
        {
          value: 'Text',
          label: '数据类型:Text',
        },
        {
          value: 'Binary',
          label: '数据类型:Binary',
        },
      ],
    }
  },
  methods: {
    SetCode(value, isClose) {
      if (this.SendType === "HEX" && isClose !== true) {
        GoGetHex(bytesToBase64(StringToBytes(value))).then((res) => {
          this.$refs.RawText.Function.setValue(res)
          this.$refs.RawText.Function.DelAllDecorations()
        })
        return
      }
      this.$refs.RawText.Function.setValue(bytesToString(value))
      this.$refs.RawText.Function.DelAllDecorations()
    }, Send() {
      this.IsSending = true
      const content = bytesToBase64(StringToBytes(this.$refs.RawText.GetCode()))
      SessionActiveSend(parseInt(Config_SelectedRow.value["Theology"]), this.direction === "Server", this.SendType, this.wsType === "Text" ? 1 : 2, content).then((res) => {
        this.IsSending = false
        if (res === "") {
          ElMessage({
            message: "发送数据成功",
            type: 'success',
          })
        } else {
          ElMessage({
            message: "主动发送,出现错误:" + res,
            type: 'error',
          })
        }
      })
    }
  },
  mounted() {
    this.$refs.RawText.SetLanguage("plaintext")
    this.$refs.RawText.SetReadOnly(false)
  },
}
</script>