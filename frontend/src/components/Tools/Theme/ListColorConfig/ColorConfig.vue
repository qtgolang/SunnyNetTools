<script>

import {DefaultColor, SetListColor} from "../../../../../bindings/changeme/Service/appmain";
import {Config_IsDark} from "../../../config/Config";

export default {
  props: ['params'],
  data() {
    return {
      color1: "",
      color2: "",
      Params: null,
      reset: false,
    }
  },
  created() {
    this.Params = this.params;
    if (this.params.data["说明"].indexOf("恢复到") !== -1) {
      this.reset = true
      return
    }
    this.color1 = this.hexToRgb(this.params.value)
    this.Params.eGridCell.addEventListener("click", () => {
      const inputElement = this.$refs.colorPicker;
      if (inputElement) {
        this.hide = inputElement.hide
        inputElement.show()
      }
    });

    this.Params.eGridCell.style.color = this.color1
    this.Params.eGridCell.style.backgroundColor = this.color1
  },
  mounted() {
    if (!this.reset) {
      this.$refs.abc.getElementsByClassName("el-color-picker__trigger")[0].className = ""
      const mm = document.getElementsByClassName("el-button el-button--small is-text el-color-dropdown__link-btn")
      for (let i = 0; i < mm.length; i++) {
        mm[i].remove()
      }
      const mm1 = document.getElementsByClassName("el-button el-button--small is-plain el-color-dropdown__btn")
      for (let i = 0; i < mm1.length; i++) {
        mm1[i].innerHTML = `<span class="">确定</span>`
      }
    }

  },
  methods: {
    hexToRgb(hex) {
      if (hex.startsWith("rgb")) {
        return hex;
      }
      hex = hex.replace(/^#/, "");
      let r = parseInt(hex.substring(0, 2), 16);
      let g = parseInt(hex.substring(2, 4), 16);
      let b = parseInt(hex.substring(4, 6), 16);
      return `rgb(${r}, ${g}, ${b})`;
    },
    colorToHex(color) {
      // 如果以 # 开头，直接返回（大写处理）
      if (color.startsWith("#")) {
        return color.toUpperCase();
      }

      // 处理 "rgb(r, g, b)" 格式
      let result = color.match(/\d+/g);
      if (!result || result.length < 3) return null;

      let r = parseInt(result[0]).toString(16).padStart(2, "0");
      let g = parseInt(result[1]).toString(16).padStart(2, "0");
      let b = parseInt(result[2]).toString(16).padStart(2, "0");

      return `#${r}${g}${b}`.toUpperCase();
    },
    onchange(eve) {
      this.Params.eGridCell.style.color = eve
      this.Params.eGridCell.style.backgroundColor = eve
      this.watchColor(eve)
    },
    buttonChange() {
      DefaultColor(Config_IsDark.value)
    },
    watchColor(color) {
      if (color === "#111") {
        return
      }
      SetListColor(Config_IsDark.value, (this.Params.data.id), this.colorToHex(color))
    }
  },
  watch: {
    color1(color) {
      //this.watchColor(color)
    },
  }
};
</script>
<template>
  <div ref="abc" style="position: absolute;left: -20px;" v-if="reset===false">
    <el-color-picker v-model="color1" size="small" ref="colorPicker" @active-change="onchange"/>
  </div>
  <div ref="aaa" v-if="reset"
       style="position: relative;top: 3px;display: flex;justify-content: center; /* 水平居中 */align-items: center;">
    <el-button size="small" style="width: 100%;margin-left: 0" @click="buttonChange">默认</el-button>
  </div>
</template>