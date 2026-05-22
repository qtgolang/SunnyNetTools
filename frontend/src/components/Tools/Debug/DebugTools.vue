<template>
  <div class="fullscreen-div"
       style="display: block;width: 100%;height: 100%;position: absolute;user-select: none;overflow:hidden">
    <div style="width: 100%;height: 100%;display: block">
      <TitleBar Title="调试工具"></TitleBar>
      <div style="width: 100% ;height:100% ;position: relative;overflow:hidden">
        <div
            style="width:  100%;height:  calc(100% - 45px);position: relative;justify-content: center;text-align: center">
          <div ref="eEmpty" class="ag-chart-empty-text ag-unselectable" style="display: block;width: 100%;height: 100%">
            <div style="width: 100%; height: 99px; display: flex; justify-content: center; align-items: center;">
              <div style="display: flex;">
                <Icon name="HTTP调试" @click="SetSelect('HTTP调试')"/>
                <Icon name="编码转换" @click="SetSelect('编码转换')"/>
                <Icon name="加密解密" @click="SetSelect('加密解密')"/>
                <Icon name="中英翻译" @click="SetSelect('中英翻译')"/>
              </div>
            </div>
            <div
                style="width: 100%;height: 1px;display: block;margin-bottom: 5px;margin-top: -5px;background-color: #817e7e;">
            </div>
            <div style="width: calc(100% - 2px);height: calc(100% - 88px)">
              <PostMan ref="httpDebug" v-show="Module==='HTTP调试'"/>
              <Encoding ref="encoding" v-show="Module==='编码转换'"/>
              <Encryption ref="Encryption" v-show="Module==='加密解密'"/>
              <Translation ref="Translation" v-if="Module==='中英翻译'" style="width: 100%;height: 100%"/>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import {Config_IsDark} from "../../config/Config";
import {Events} from "@wailsio/runtime";
import {IsDark} from "../../../../bindings/changeme/Service/appmain";
import CertInstall from "../Other/CertInstall/CertInstall.vue";
import TitleBar from "../../TitleBar/TitleBar.vue";
import Icon from "../icon.vue";
import Translation from "./Encoding/Translation.vue";
import Encryption from "./Encoding/Encryption.vue";
import Encoding from "./Encoding/Encoding.vue";
import PostMan from "./HTTP/PostMan.vue";

export default {
  components: {
    Icon, Translation, Encryption, Encoding,PostMan,
    TitleBar, CertInstall
  },
  computed: {},
  watch: {
    DropFiles(n) {
      if (n !== "") {
        if (this.Module === "加密解密") {
          this.$refs.Encryption.DropFilesEvent(n)
        } else if (this.Module === "编码转换") {
          this.$refs.encoding.DropFilesEvent(this.DropFiles)
        }
      }
    },
  },
  data() {
    return {
      DropFiles: "",
      Module: "HTTP调试",
    }
  },
  mounted() {
    Events.On("DropFilesEvent", (e) => {
      if (e.data[0] !== "调试工具") {
        //其他窗口的拖放事件
        return;
      }
      this.DropFiles = e.data[1][0]
    })
    Events.On("SetIsDark", (obj) => {
      const dark=obj.data[0]===true || obj.data===true;
      if (Config_IsDark.value === dark){
        return
      }
      this.setIsDark(dark)
    })
    IsDark().then(isDark => {
      this.setIsDark(isDark)
    })
  },
  beforeUnmount() {
  },
  methods: {
    SetSelect(m) {
      this.Module = m;
    },
    setIsDark(IsDark) {
      try {
        const htmlElement = document.documentElement;
        if (IsDark) {
          htmlElement.setAttribute('data-dark-mode', 'true');
          htmlElement.setAttribute('data-ag-theme-mode', 'dark-blue');
          htmlElement.style.backgroundColor = 'rgb(23,21,21)';
        } else {
          htmlElement.setAttribute('data-dark-mode', '');
          htmlElement.setAttribute('data-ag-theme-mode', '');
          htmlElement.style.backgroundColor = 'rgb(241,241,241)';
        }
        Config_IsDark.value = IsDark
      } catch (e) {

      }
    },
  },

}
</script>
<style>
.fullscreen-div {
  position: fixed; /* 让 div 相对于视口固定 */
  top: 0;
  left: 0;
  width: 100vw; /* 100% 视口宽度 */
  height: 100vh; /* 100% 视口高度 */
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  color: white;
}
</style>
<script setup lang="ts">
</script>
