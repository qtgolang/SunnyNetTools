<script>
//

import {CallGoDo, SetColorConfig} from "../../CallbackEventsOn.js";

export default {
  props: ['params'],
  data() {
    return {
      color1: "",
      Params: null,
      reset: false,
    }
  },
  created() {
    this.Params = this.params;
    if (this.params.data.alias === "reset") {
      this.reset = true
      return
    }
    this.color1 = this.params.value
    this.Params.eGridCell.addEventListener("click", () => {
      const inputElement = this.$refs.colorPicker;
      this.hide = inputElement.hide
      inputElement.show()
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
  }, methods: {
    onchange(eve) {
      this.Params.eGridCell.style.color = eve
      this.Params.eGridCell.style.backgroundColor = eve
      let res = {}
      if (this.Params.column.colId === "深色主题") {
        res.right = this.Params.data["浅色主题"]
        res.dark = eve
      } else {
        res.dark = this.Params.data["深色主题"]
        res.right = eve
      }
      SetColorConfig(this.Params.data.alias, res)
    },
    buttonChange() {
      CallGoDo("重置颜色列表", {dark: this.Params.column.colId === "深色主题"}).then(res => {
        window.interface.colorConfig = res
        window.vm.ListColorManager.Load()
        window.vm.List.RefreshRenderedNodes()
      })
    }
  }
};
</script>
<template>
  <div ref="abc" style="position: absolute;left: -20px;" v-if="reset===false">
    <el-color-picker v-model="color1" size="small" ref="colorPicker" @active-change="onchange"/>
  </div>
  <div ref="aaa" v-if="reset" style="position: absolute;top: -2px;">
    <el-button size="small" style="width: 100%" @click="buttonChange">默认</el-button>
  </div>
</template>