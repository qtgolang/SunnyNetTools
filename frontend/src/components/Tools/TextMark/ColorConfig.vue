<script>

export default {
  props: ['params'],
  data() {
    return {
      color1: "",
      color2: "",
      Params: null,
      init: false,
    }
  },
  created() {
    this.Params = this.params;
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
    this.$refs.abc.getElementsByClassName("el-color-picker__trigger")[0].className = ""
    const mm = document.getElementsByClassName("el-button el-button--small is-text el-color-dropdown__link-btn")
    for (let i = 0; i < mm.length; i++) {
      mm[i].remove()
    }
    const mm1 = document.getElementsByClassName("el-button el-button--small is-plain el-color-dropdown__btn")
    for (let i = 0; i < mm1.length; i++) {
      mm1[i].innerHTML = `<span class="">确定</span>`
    }
    this.init = true
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
    onchange(eve) {
      this.Params.eGridCell.style.color = eve
      this.Params.eGridCell.style.backgroundColor = eve
    },
  },
  watch: {
    color1(color) {
      if (this.init) {
        this.Params.data["浅色主题"] = color
        try {
          this.Params.data.save()
        } catch (e) {
        }
      }
    },
  }
};
</script>
<template>
  <div ref="abc" style="position: absolute;left: -20px;">
    <el-color-picker v-model="color1" size="small" ref="colorPicker" @active-change="onchange"/>
  </div>
</template>