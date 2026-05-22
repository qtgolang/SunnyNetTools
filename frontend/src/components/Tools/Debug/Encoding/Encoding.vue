<script>
import VueText from "../../Theme/Text.vue";
import {ElNotification} from "element-plus";
import {AppEncodingCall, ClipboardWriteAll} from "../../../../../bindings/changeme/Service/appmain";
import {EncodingInfo} from "../../../../../bindings/changeme/Service/Tools/DebugTools";
import {base64ToBytes, bytesToString} from "../../../config/encoding";
import {getImageType} from "../../../config/CheckImageType";

const compress = [
  {
    value: '压缩',
    label: '压缩',
    children: [
      {
        value: '输入为GBK',
        label: '输入为GBK',
        children: [
          {
            value: '到HEX',
            label: '到HEX',
          },
          {
            value: '到Base64',
            label: '到Base64',
          },
        ]
      },
      {
        value: '输入为UTF8',
        label: '输入为UTF8',
        children: [
          {
            value: '到HEX',
            label: '到HEX',
          },
          {
            value: '到Base64',
            label: '到Base64',
          },
        ]
      },
    ]
  },
  {
    value: '解压缩',
    label: '解压缩',
    children: [
      {
        value: '输入为HEX',
        label: '输入为HEX',
        children: [
          {
            value: '到HEX',
            label: '到HEX',
          },
          {
            value: '到Base64',
            label: '到Base64',
          },
          {
            value: '到文本',
            label: '到文本',
          },
        ]
      },
      {
        value: '输入为Base64',
        label: '输入为Base64',
        children: [
          {
            value: '到HEX',
            label: '到HEX',
          },
          {
            value: '到Base64',
            label: '到Base64',
          },
          {
            value: '到文本',
            label: '到文本',
          },
        ]
      },
    ]
  },
]

export default {
  components: {VueText},
  computed: {},
  data() {
    return {
      ImgBody: "",
      centerDialogVisible: false,
      LocalTmp1: [],
      LocalTmp2: [],
      EncodingOptions: [
        {
          value: '字符串',
          label: '字符串',
          children: [
            {
              value: 'URL编码',
              label: 'URL编码',
            },
            {
              value: 'URL解码',
              label: 'URL解码',
            },
            {
              value: '转到Base64',
              label: '转到Base64',
            },
            {
              value: '转到HEX',
              label: '转到HEX',
            },
          ]
        },
        {
          value: 'Base64',
          label: 'Base64',
          children: [
            {
              value: '解码到字符串',
              label: '解码到字符串',
            },
            {
              value: '解码到Hex',
              label: '解码到Hex',
            },
            {
              value: '解码到图片',
              label: '解码到图片',
            },
          ]
        },
        {
          value: 'HEX',
          label: 'HEX',
          children: [
            {
              value: '解码到字符串',
              label: '解码到字符串',
            },
            {
              value: '解码到Base64',
              label: '解码到Base64',
            },
          ]
        },
        {
          value: '文件',
          label: '文件',
          children: [
            {
              value: '到字符串',
              label: '到字符串',
            },
            {
              value: '到Base64',
              label: '到Base64',
            },
            {
              value: '到Hex',
              label: '到Hex',
            },
            {
              value: '效验',
              label: '效验',
              children: [
                {
                  value: 'CRC32',
                  label: 'CRC32',
                },
                {
                  value: 'MD5',
                  label: 'MD5',
                },
                {
                  value: 'SHA-1',
                  label: 'SHA-1',
                },
                {
                  value: 'SHA-256',
                  label: 'SHA-256',
                },
                {
                  value: 'SHA-512',
                  label: 'SHA-512',
                },
              ]
            },
            {
              value: '到图片预览',
              label: '到图片预览',
            },
          ]
        },
        {
          value: 'USC2转Ansi',
          label: 'USC2转Ansi',
        },
        {
          value: 'Ansi转USC2',
          label: 'Ansi转USC2',
        },
      ],
      otherOptions: [
        {
          value: '文本操作',
          label: '文本操作',
          children: [
            {
              value: '到大写',
              label: '到大写',
            },
            {
              value: '到小写',
              label: '到小写',
            },
            {
              value: '删除',
              label: '删除',
              children: [
                {
                  value: '全部空格',
                  label: '全部空格',
                },
                {
                  value: '全部换行',
                  label: '全部换行',
                },
                {
                  value: '全部换行及空格',
                  label: '全部换行及空格',
                },
              ]
            },
            {
              value: '取长度',
              label: '取长度',
            },
          ]
        },
        {
          value: '文本转换',
          label: '文本转换',
          children: [
            {
              value: '字节集',
              label: '字节集',
              children: [
                {
                  value: '文本到字节集',
                  label: '文本到字节集',
                },
                {
                  value: '字节集到文本',
                  label: '字节集到文本',
                },
                {
                  value: '字节集到HEX',
                  label: '字节集到HEX',
                },
                {
                  value: '字节集到Base64',
                  label: '字节集到Base64',
                },
              ]
            },
            {
              value: '字符',
              label: '字符',
              children: [
                {
                  value: '文本到字符',
                  label: '文本到字符',
                },
                {
                  value: '字符到文本',
                  label: '字符到文本',
                },
                {
                  value: '字符到HEX',
                  label: '字符到HEX',
                },
                {
                  value: '字符到Base64',
                  label: '字符到Base64',
                },
              ]
            },
          ]
        },
        {
          value: '压缩转换',
          label: '压缩转换',
          children: [
            {
              value: 'GZIP',
              label: 'GZIP',
              children: compress
            },
            {
              value: 'ZLIB',
              label: 'ZLIB',
              children: compress
            },
            {
              value: 'Brotli',
              label: 'Brotli',
              children: compress
            },
            {
              value: 'ZSTD',
              label: 'ZSTD',
              children: compress
            },
            {
              value: 'Deflate',
              label: 'Deflate',
              children: compress
            },
          ]
        },
        {
          value: '其他操作',
          label: '其他操作',
          children: [
            {
              value: '获取当前时间戳',
              label: '获取当前时间戳',
            },
            {
              value: '时间戳到时间',
              label: '时间戳到时间',
            },
            {
              value: '本地文件',
              label: '本地文件',
              children: [
                {
                  value: 'DEX重命名',
                  label: 'DEX重命名',
                },
                {
                  value: '取HEX',
                  label: '取HEX',
                },
                {
                  value: '取Base64',
                  label: '取Base64',
                },
              ]
            },
            {
              value: '参数排序',
              label: '参数排序',
              children: [
                {
                  value: 'a-zA-Z0-9',
                  label: 'a-zA-Z0-9',
                },
                {
                  value: 'z-aZ-A0-9',
                  label: 'z-aZ-A0-9',
                },
                {
                  value: 'Z-Az-a0-9',
                  label: 'Z-Az-a0-9',
                },
                {
                  value: 'A-Za-z0-9',
                  label: 'A-Za-z0-9',
                },
                {
                  value: '0-9z-aZ-A',
                  label: '0-9z-aZ-A',
                },
                {
                  value: '0-9Z-Az-a',
                  label: '0-9Z-Az-a',
                },
                {
                  value: '0-9a-zA-Z',
                  label: '0-9a-zA-Z',
                },
                {
                  value: '0-9A-Za-z',
                  label: '0-9A-Za-z',
                },
              ]
            },
          ]
        },
      ],
      props: {expandTrigger: 'hover'},
    }
  },
  watch: {
    LocalTmp1(n) {
      if (n === "" || n === [] || n.length < 1) {
        return
      }
      for (let i = 0; i < n.length; i++) {
        if (n[i] === "解码到图片") {
          this.LocalTmp1 = []
          this.showImg();
          this.centerDialogVisible = true
          return;
        }
      }
      AppEncodingCall(EncodingInfo.createFrom({
        children: n,
        Data: this.$refs.Data.GetCode()
      })).then(res => {
        for (let i = 0; i < n.length; i++) {
          if (n[i] === "到图片预览") {
            this.LocalTmp1 = []
            this.centerDialogVisible = true
            const ImageType = getImageType(base64ToBytes(res));
            this.ImgBody = "data:image/" + ImageType + ";base64," + res
            return;
          }
        }
        this.$refs.Res.SetCode(bytesToString(base64ToBytes(res)))
        this.LocalTmp1 = []
      })
      this.LocalTmp1 = []
    },
    LocalTmp2(n) {
      if (n === "" || n === [] || n.length < 1) {
        return
      }
      AppEncodingCall(EncodingInfo.createFrom({
        children: n,
        Data: this.$refs.Data.GetCode()
      })).then(res => {
        this.$refs.Res.SetCode(bytesToString(base64ToBytes(res)))
        this.LocalTmp2 = []
      })
      this.LocalTmp2 = []
    },
  },
  methods: {
    DropFilesEvent(path) {
      this.$refs.Data.SetLanguage("text")
      this.$refs.Data.SetCode(path)
    },
    exchange() {
      const a = this.$refs.Data.GetCode()
      const b = this.$refs.Res.GetCode()
      this.$refs.Data.SetCode(b)
      this.$refs.Res.SetCode(a)
    },
    empty() {
      this.$refs.Data.SetCode("")
      this.$refs.Res.SetCode("")
    },
    cover() {
      this.LocalTmp2 = ["引号转换"]
    },
    copyRes() {
      ClipboardWriteAll(this.$refs.Res.GetCode()).then(res => {
        if (res === "") {
          ElNotification({
            position: 'bottom-right',
            message: '复制成功\n\n已复制到剪辑版',
            type: 'success',
            customClass: 'multiline-message'
          })
        } else {
          ElNotification({
            position: 'bottom-right',
            showClose: true,
            message: '复制失败\n\n' + res,
            type: 'warning',
            customClass: 'multiline-message'
          })
        }
      })
    },
    showImg() {
      const ImageType = getImageType(base64ToBytes(this.$refs.Data.GetCode()));
      this.ImgBody = "data:image/" + ImageType + ";base64," + this.$refs.Data.GetCode()
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
  <div style=" width: 100%; height:100%;overflow:hidden">
    <div style=" width: 100%; height:calc(((100%)/2) - 30px );overflow:hidden">
      <VueText ref="Data" :glyphMargin="false" :readOnly="false" Language="'text'" Name="data"
               style="margin-top: -1px;margin-left: -1px;"/>
    </div>
    <div style=" width: 100%; height:60px;overflow:hidden; display: flex;justify-content: center;align-items: center;">
      <el-cascader
          v-model="LocalTmp1"
          placeholder="编码转换"
          :options="EncodingOptions"
          :props="props"
          style="width: 120px;margin-right: 10px"
      />
      <el-cascader
          v-model="LocalTmp2"
          placeholder="其他选项"
          :options="otherOptions"
          :props="props"
          style="width: 120px;margin-right: 10px"
      />
      <el-button style="width: calc(100% - 790px);" @click="exchange">交换</el-button>
      <el-button style="width: 150px;" @click="empty">清空</el-button>
      <el-button style="width: 150px;" @click="cover">引号转换</el-button>
      <el-button style="width: 150px;" @click="copyRes">复制结果</el-button>
    </div>
    <div style=" width: 100%; height:calc(((100%)/2) - 30px );overflow:hidden">
      <VueText ref="Res" :glyphMargin="false" :readOnly="false" Language="'text'" Name="data"
               style="margin-top: -1px;margin-left: -1px;"/>
    </div>
    <el-dialog
        v-model="centerDialogVisible"
        title="图片预览"
        width="800px"
        align-center
    >
      <div style="width: 100%; height: 600px; overflow: hidden;">
        <img
            ref="imgObj"
            :src="ImgBody"
            style="height: 100%; width: 100%; object-fit: contain; object-position: center; display: block;"
            alt=""/>
      </div>
    </el-dialog>

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