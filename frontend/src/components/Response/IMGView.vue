<template>
  <div style="height: 100%;width: 100%;  overflow: hidden; position: relative">
    <div style="height: 100%;width: 150px">
      <div v-html="GetInfo"></div>
      <br>
      <el-select v-model="imgStyle" placeholder="Select" size="small" style="width: 90px" v-show="this.Body!==''">
        <el-option
            v-for="item in options"
            :key="item.value"
            :label="item.label"
            :value="item.value"
        />
      </el-select>
      <br>
      <br>
      <el-button size="small" @click="Save" v-show="this.Body!==''">保存到桌面</el-button>
    </div>
    <div v-show="this.Body!==''"
         style="height: 100%;width: calc(100% - 150px);position: absolute;top:0px;left: 150px;    display: flex;    justify-content: center;  align-items: center;">
      <img ref="imgObj" :src="Body" :style="imgStyle"/>
    </div>
  </div>
</template>
<script>
import {Base64Decode, CallGoDo} from "../CallbackEventsOn.js";

export default {
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
      ]
    }
  },
  methods: {
    SetImg(B64IMG, geshi) {
      this.Body = "data:image/" + geshi + ";base64," + B64IMG
      this.$nextTick(() => {
        let type = geshi.toLowerCase()
        if (type.indexOf("svg") !== -1) {
          type = "svg"
        }
        this.Type = type
        const size = Base64Decode(B64IMG).length.toLocaleString();
        const width = this.$refs.imgObj.naturalWidth
        const height = this.$refs.imgObj.naturalHeight
        this.Info = `<br> <br> <span>类型: ${type} <span><br><br><span>尺寸: ${width} * ${height} <span><br><br><span>大小: ${size} Bytes<span><br><br>`

      })
    }
    ,
    Save() {
      CallGoDo("保存响应图片", {Theology: window.Theology, type: this.Type})
    }
  }
}
</script>
<style>


</style>