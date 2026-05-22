<template>
  <div style="height: 100%;width: 100%;  overflow: hidden; position: relative" v-show="Show">
    <div style="height: 100%;width: 150px;display: grid;justify-content: center">
      <div v-html="GetInfo"></div>
      <br>
      <el-select v-model="imgStyle" placeholder="Select" size="small" style="width: 90px">
        <el-option
            v-for="item in options"
            :key="item.value"
            :label="item.label"
            :value="item.value"
        />
      </el-select>
      <br>
      <br>
      <el-button size="small" @click="Save" style="width: 90px;margin-left: 0px">保存到文件
      </el-button>
    </div>
    <div
        style="height: calc(100% - 10px);width: calc(100% - 157px);position: absolute;top:0px;left: 150px;    display: flex;    justify-content: center;  align-items: center;border:double">
      <img ref="imgObj" :src="Body" :style="imgStyle"/>
    </div>
  </div>
</template>
<script>
import {base64ToBytes} from "../config/encoding.js";
import {ElNotification} from "element-plus";
import {Config_SelectedRow} from "../config/Config";
import {AppSaveRequestImg} from "../../../bindings/changeme/Service/appmain";
import {Dialogs} from "@wailsio/runtime";

export default {
  props: ['isRequest'],
  computed: {
    GetInfo() {
      return this.Info
    }
  },
  data() {
    return {
      Type: "jpg",
      ReadOnly: false,
      WindowStyle: "",
      imgStyle: "height: 100%;width: 100%;object-fit: contain; object-position: center;",
      Body: "",
      Info: "",
      options: [
        {
          value: 'height: 100%;width: 100%;object-fit: contain; object-position: center;',
          label: '居中显示',
        },
        {
          value: "height: 100%;width: 100%;object-fit: fill; object-position: center;",
          label: '拉伸显示',
        },
      ],
      get SelectedRow() {
        return Config_SelectedRow.value
      },
      set SelectedRow(value) {
        Config_SelectedRow.value = value
      },
      Show: false,
    }
  },
  methods: {
    SetImg(B64IMG, geshi) {
      try {
        const size = base64ToBytes(B64IMG).length.toLocaleString();
        this.Body = "data:image/" + geshi + ";base64," + B64IMG
        this.Show = geshi !== "" && B64IMG !== "";
        let type = (geshi + "").toLowerCase()
        if (type.indexOf("svg") !== -1) {
          type = "svg"
        }
        this.Type = type
        this.$nextTick(() => {
          const width = this.$refs.imgObj.naturalWidth
          const height = this.$refs.imgObj.naturalHeight
          this.Info = `<br> <br> <span>类型: ${type} <span><br><br><span>尺寸: ${width} * ${height} <span><br><br><span>大小: ${size} Bytes<span><br><br>`
        })
      } catch (e) {
      }
    }
    ,
    Save() {
      const _SelectedTheology = parseInt(this.SelectedRow?.Theology ?? "0");

      const options = {
        TreatPackagesAsDirectories: true,
        CanCreateDirectories: true,
        Filters: [
          {DisplayName: "图片文件", Pattern: "*." + this.Type},
        ],
        Filename:"img.jpg",
        Title: "请选择要保存到的文件",
      };
      Dialogs.SaveFile(options).then((selectedFiles) => {
        const filePath = selectedFiles.substring(0, selectedFiles.lastIndexOf('/'));
        const fileName = selectedFiles.split('/').pop();
        const finalFileName = fileName.includes('.') ? fileName : `${fileName}.jpg`;
        const finalFilePath = `${filePath}/${finalFileName}`;
        AppSaveRequestImg(_SelectedTheology, this.Type, this.isRequest, finalFilePath).then(res => {
          if (res === "") {
            ElNotification({
              position: 'bottom-right',
              showClose: true,
              message: '保存成功\r\n已经储存至：' + selectedFiles,
              type: 'success',
              customClass: 'multiline-message'
            })
          } else {
            ElNotification({
              position: 'bottom-right',
              showClose: true,
              message: '保存失败\r\n' + res,
              type: 'error',
              customClass: 'multiline-message'
            })
          }
        })
      })
    }
  }
}
</script>
<style>
.multiline-message {
  white-space: pre-line;
}
</style>