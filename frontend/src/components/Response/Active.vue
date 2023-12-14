<template>
  <VueText ref="RawText" :height="BodyRectHeight" :glyphMargin="false" :readOnly="false" Language="'html'"
           Name="ActiveView"/>
  <div :style="`position: absolute; top: ${RectHeight};display: flex;width: 100%;`">
    <el-select v-model="SendType" class="m-2" placeholder="Select" size="small" style="flex-grow: 0;width: 157px">
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

</template>

<script>
import VueText from "../Request/Text.vue";
import {CallGoDo, StrBase64Encode} from "../CallbackEventsOn.js";

export default {
  props: ['Height'],
  watch: {
    Height(value) {
      this._Height = value
    },
  },
  components: {VueText},
  computed: {
    BodyRectHeight() {
      return (parseInt(this._Height.replaceAll("px", "")) - 28) + "px"
    },
    RectHeight() {
      return (parseInt(this._Height.replaceAll("px", "")) + 5) + "px"
    },
  },
  data() {
    return {
      IsWs: false,
      IsSending: false,
      _Height: "0px",
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
    SetCode(value) {
      this.$refs.RawText.SetCode(value)
    }, Send() {
      this.IsSending = true
      let content = StrBase64Encode(this.$refs.RawText.GetCode())
      let way = window.vm.List.agSelectedLine.data["方式"].toUpperCase()
      const Data = {
        direction: this.direction,
        SendType: this.SendType,
        wsType: this.wsType,
        Data: content,
        Theology: window.Theology,
        IsWs: window.vm.Tabs.Request.DisplayHTTPHeader,
        IsTCP: way.indexOf("TCP") !== -1
      }
      CallGoDo("主动发送", Data).then(res => {
        this.IsSending = false
      })
    }
  }
}
</script>