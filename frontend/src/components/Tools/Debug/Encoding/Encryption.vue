<script>
import VueText from "../../Theme/Text.vue";
import {Config_IsDark} from "../../../config/Config";
import {AppEncryptCall} from "../../../../../bindings/changeme/Service/appmain";
import {EncryptionInfo} from "../../../../../bindings/changeme/Service/Tools/DebugTools";
import {ElNotification} from "element-plus";
import {base64ToBytes, bytesToString} from "../../../config/encoding";

const children_hash = [
  {
    value: 'MD2',
    label: 'MD2',
  },
  {
    value: 'MD4',
    label: 'MD4',
  },
  {
    value: 'MD5',
    label: 'MD5',
  },
  {
    value: 'SHA-1',
    label: 'SHA-1',
  }
  ,
  {
    value: 'SHA-256',
    label: 'SHA-256',
  }
  ,
  {
    value: 'SHA-512',
    label: 'SHA-512',
  }
]
const children_des_3des = [
  {
    value: 'ECB',
    label: 'ECB',
  },
  {
    value: 'CBC',
    label: 'CBC',
  },
  {
    value: 'WithSalt',
    label: 'WithSalt',
  },
]
const children_aes = [
  {
    value: 'ECB',
    label: 'ECB',
  },
  {
    value: 'CBC',
    label: 'CBC',
  },
  {
    value: 'WithSalt-256',
    label: 'WithSalt-256',
  },
  {
    value: 'WithSalt-192',
    label: 'WithSalt-192',
  },
  {
    value: 'WithSalt-128',
    label: 'WithSalt-128',
  },
]
const children2 = [
  {
    value: '使用Hash',
    label: '使用Hash',
    children: children_hash
  },
  {
    value: '使用Hmac-Hash',
    label: '使用Hmac-Hash',
    children: children_hash
  },
  {
    value: '使用AES',
    label: '使用AES',
    children: children_aes
  },
  {
    value: '使用DES',
    label: '使用DES',
    children: children_des_3des
  },
  {
    value: '使用3DES',
    label: '使用3DES',
    children: children_des_3des
  },
  {
    value: '使用SM2',
    label: '使用SM2',
  },
  {
    value: '使用RC4',
    label: '使用RC4',
  },
  {
    value: '使用RC4-WithSalt',
    label: '使用RC4-WithSalt',
  }
]
const children = [
  {
    value: '到字符串',
    label: '到字符串',
    children: children2
  },
  {
    value: '到Hex',
    label: '到Hex',
    children: children2
  },
  {
    value: '到Base64',
    label: '到Base64',
    children: children2
  },
]

export default {
  components: {VueText},
  computed: {
    getKeyStyle() {
      if (this.isIv === false) {
        return "width:80%;"
      }
      return "width:40%;margin-right: 10px;"
    },
    getEncStyle() {
      if (this.isIv === false && this.isKey === false && this.isDecrypt === false) {
        return "width:100%;"
      }
      if (this.isDecrypt === false) {
        return "width:20%;margin-left: 10px"
      }
      return "width:10%;margin-left: 10px"
    },
    getBorderStyle() {
      if (this.theme) {
        this.borderStyle = '#c6bebe'
      } else {
        this.borderStyle = '#8a8484'
      }
      return 'width: calc(100%); height:calc(50% - 1px);overflow:hidden;border: 1px solid ' + this.borderStyle + ';margin-left: -1px'
    },
    BorderStyle() {
      return 'width: calc(100% -1px); height:calc(100% + 6px);overflow:hidden;border: 1px solid ' + this.borderStyle + ';margin-top: -6px;margin-left: 1px;'
    },
    isSm2() {
      for (let i = 0; i < this.LocalTmp.length; i++) {
        if (this.LocalTmp[i] === "使用SM2") {
          return true
        }
      }
      return false
    },
  },
  data() {
    return {
      get theme() {
        return Config_IsDark.value
      },
      set theme(newValue) {
        if (Config_IsDark.value !== newValue) Config_IsDark.value = newValue
      },
      borderStyle: "",
      DragFilePath: "",
      LocalTmp: [],
      sm2Mode: "C1C3C2",
      sm2Modes: [
        {
          value: "C1C2C3",
          label: "C1C2C3"
        },
        {
          value: "C1C3C2",
          label: "C1C3C2"
        }
      ],
      EncodingOptions: [
        {
          value: '输入为GBK字符串',
          label: '输入为GBK字符串',
          children: children
        },
        {
          value: '输入为UTF8字符串',
          label: '输入为UTF8字符串',
          children: children
        },
        {
          value: '从HEX',
          label: '从HEX',
          children: children
        },

        {
          value: '从Base64',
          label: '从Base64',
          children: children
        },
        {
          value: '从文件',
          label: '从文件',
          children: children
        },
      ],
      props: {expandTrigger: 'hover'},
      isKey: true,
      isIv: true,
      isDecrypt: true,
      keyType: "String",
      key: "",
      ivType: "String",
      iv: ""
    }
  },
  watch: {},
  methods: {
    DropFilesEvent(path) {
      this.$refs.Data.SetLanguage("text")
      this.$refs.Data.SetCode(path)
    },
    EncodingSelect(obj) {
      const a = obj[2] + ""
      if (a === "使用RC4-WithSalt" || a === "使用RC4" || a === "使用SM2") {
        this.isIv = false
        this.isKey = true
        this.isDecrypt = true
        return;
      }
      if (obj.length >= 4) {
        const b = obj[3] + ""
        if (b.indexOf("WithSalt") !== -1) {
          this.isIv = false
          this.isKey = true
          this.isDecrypt = true
          return
        }
        if (b.indexOf("CBC") !== -1) {
          this.isIv = true
          this.isKey = true
          this.isDecrypt = true
          return
        }
        if (b.indexOf("ECB") !== -1) {
          this.isIv = false
          this.isKey = true
          this.isDecrypt = true
          return
        }
      }
      this.isIv = false;
      this.isKey = a.indexOf("Hmac") !== -1;
      this.isDecrypt = this.isIv
    },
    Encrypt(mode) {
      if (this.LocalTmp === "") {
        ElNotification({
          position: 'bottom-right',
          message: '执行失败\n\n请选择要执行的 加密/解密 操作!!!',
          type: 'warning',
          customClass: 'multiline-message'
        })
        return
      }
      if (this.isKey && this.key === '') {
        ElNotification({
          position: 'bottom-right',
          message: '执行失败\n\n请输入要 加密/解密 KEY',
          type: 'warning',
          customClass: 'multiline-message'
        })
      }
      if (this.isIv && this.iv === '') {
        ElNotification({
          position: 'bottom-right',
          message: '执行失败\n\n请输入要 加密/解密 IV',
          type: 'warning',
          customClass: 'multiline-message'
        })
      }
      AppEncryptCall(EncryptionInfo.createFrom({
        children: this.LocalTmp,
        key: this.key,
        iv: this.iv,
        ivType: this.ivType,
        keyType: this.keyType,
        isEncrypt: mode === 1,
        SM2Mode: this.sm2Mode,
        Data: this.$refs.Data.GetCode()
      })).then(res => {
        this.$refs.Res.SetCode(bytesToString(base64ToBytes(res)))
      })
    }
  },
  mounted() {
    this.$refs.Data.SetLanguage("text")
    this.$refs.Data.SetReadOnly(false)
    this.$refs.Res.SetLanguage("text")
    this.$refs.Res.SetReadOnly(false)
  },
};
</script>

<template>
  <div :style="BorderStyle">
    <div
        style="height:40px;overflow:hidden; display: flex;justify-content: center;align-items: center;margin-left: 10px;margin-right: 10px;margin-top: 5px">
      <el-cascader
          v-model="LocalTmp"
          placeholder="加密配置选择"
          :options="EncodingOptions"
          :props="props"
          @change="EncodingSelect"
          style="width: 100%; display: flex;margin-top:-8px"
      />
    </div>

    <div
        style="height:50px;overflow:hidden; display: flex;margin-bottom: -12px;margin-left: 11px;margin-right: 10px;font-size: 15px;align-items: flex-start">
      <el-form-item label="请输入KEY" v-if="isKey" :style="getKeyStyle">
        <el-input
            style="display: flex;"
            class="input-with-select"
            v-model="key"
        >
          <template #prepend>
            <div v-show="!isSm2">
              <el-select v-model="keyType" placeholder="Select" style="width: 95px;margin-top: -2px">
                <el-option label="String" value="String"/>
                <el-option label="Base64" value="Base64"/>
                <el-option label="HEX" value="HEX"/>
              </el-select>
            </div>
            <div v-show="isSm2">
              Base64
            </div>
          </template>
        </el-input>
      </el-form-item>
      <el-form-item label="请输入IV" v-if="isIv" style="width: 40%;margin-left: 10px">
        <el-input
            style="display: flex;"
            size="small"
            class="input-with-select"
            v-model="iv"
        >
          <template #prepend>
            <el-select v-model="ivType" placeholder="Select" style="width: 95px">
              <el-option label="String" value="String"/>
              <el-option label="Base64" value="Base64"/>
              <el-option label="HEX" value="HEX"/>
            </el-select>
          </template>
        </el-input>
      </el-form-item>

      <div v-show="isSm2">
        <el-cascader
            v-model="sm2Mode"
            placeholder="SM2加密Mode选择"
            :options="sm2Modes"
            :props="props"
            style="width: 120px;margin-left: 10px"
        />
      </div>
      <el-button :style="getEncStyle" @click="Encrypt(1)">执行加密</el-button>
      <el-button v-if="isDecrypt" style="width: 10%" @click="Encrypt(2)">执行解密</el-button>
    </div>
    <div style=" width: 100%; height:calc(100% - 82px);overflow:hidden;">
      <div :style="getBorderStyle">
        <VueText ref="Data" :glyphMargin="false" Language="'text'" Name="data"/>
      </div>
      <div :style="getBorderStyle" style="margin-top: -1px">
        <VueText ref="Res" :glyphMargin="false" Language="'text'" Name="data"/>
      </div>

    </div>


  </div>
</template>
<style>
.el-form-item {
  margin-bottom: 4px
}

.el-cascader-menu__wrap.el-scrollbar__wrap {
  height: 220px !important;
}
</style>

<style>
.multiline-message {
  white-space: pre-line;
}
</style>