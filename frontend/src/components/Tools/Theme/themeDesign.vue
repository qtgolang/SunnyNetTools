<script>
import {Events} from "@wailsio/runtime";
import {IsDark} from "../../../../bindings/changeme/Service/appmain.js";
import TitleBar from "../../TitleBar/TitleBar.vue";
import {Config_IsDark} from "../../config/Config.js";

export default {
  components: {TitleBar},
  data() {
    return {}
  },
  mounted() {
    Events.On("SetIsDark", (obj) => {
      try {
        const dark=obj.data[0]===true || obj.data===true;
        if (Config_IsDark.value === dark){
          return
        }
        const htmlElement = document.documentElement;
        if (dark) {
          htmlElement.setAttribute('data-dark-mode', 'true');
          htmlElement.setAttribute('data-ag-theme-mode', 'dark-blue');
        } else {
          htmlElement.setAttribute('data-dark-mode', '');
          htmlElement.setAttribute('data-ag-theme-mode', '');
        }
      } catch (e) {

      }
    })
    IsDark().then((res) => {
      try {
        const htmlElement = document.documentElement;
        if (res) {
          htmlElement.setAttribute('data-dark-mode', 'true');
          htmlElement.setAttribute('data-ag-theme-mode', 'dark-blue');
        } else {
          htmlElement.setAttribute('data-dark-mode', '');
          htmlElement.setAttribute('data-ag-theme-mode', '');
        }
      } catch (e) {

      }
    })

  }
}
</script>

<template>
  <TitleBar Title="主题调色" style="height: 30px;z-index: 9999"></TitleBar>
</template>
