<script>

import {Events} from "@wailsio/runtime";
import TitleBar from "../../TitleBar/TitleBar.vue";
import {Config_IsDark} from "../../config/Config";
import CertInstall from "./CertInstall/CertInstall.vue";
import {IsDark} from "../../../../bindings/changeme/Service/appmain";
import CodeCreate from "./CodeCreate/CodeCreate.vue";
import DiffText from "./DiffText/diffeditor.vue";

export default {
  components: {CodeCreate, CertInstall, TitleBar, DiffText},
  data() {
    return {
      Title: "",
      Url: "",
    }
  },
  watch: {},
  computed: {
    isCertInstall() {
      return this.Title === '证书安装'
    },
    isCodeCreate() {
      return this.Title === '代码生成'
    },
    isDiffText() {
      return this.Title === '文本对比'
    },
    isMCP() {
      return this.Title === 'MCO能力描述'
    },
  },
  mounted() {
    Events.On("SetIsDark", (obj) => {
      const dark = obj.data[0] === true || obj.data === true;
      if (Config_IsDark.value === dark) {
        return
      }
      this.setIsDark(dark)
    })
    Events.On("LoadUrl", (obj) => {
      try {
        console.log(obj.data)
        this.Title = obj.data[0]
        if ((obj.data[1] + "").startsWith("http")) {
          this.Url = obj.data[1]
        } else {
          this.Url = "about:blank"
        }
        setTimeout(() => {
          try {
            const iframe = document.getElementById('myIframe');
            iframe.contentWindow.document.addEventListener('contextmenu', function (e) {
              e.preventDefault();
            });
          } catch (err) {
            // 如果 iframe 跨域，就无法访问其 contentWindow.document
            console.warn('无法访问 iframe 内容，可能是跨域限制');
          }
        }, 1000)
      } catch (e) {
      }
    })
    IsDark().then(isDark => {
      this.setIsDark(isDark)
    })
  },
  methods: {
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
          htmlElement.style.backgroundColor = 'rgb(255, 255, 255)';
        }
        Config_IsDark.value = IsDark
      } catch (e) {

      }
    },
  }
}
</script>

<template>
  <div class="fullscreen-div"
       style="display: block;width: 100%;height: 100%;position: absolute;user-select: none;overflow:hidden">
    <div style="width: 100%;height: 100%;display: block">
      <TitleBar :Title="Title"></TitleBar>
      <div style="width: 100% ;height:100% ;position: relative;overflow:hidden">
        <iframe id="myIframe" v-show="(!isCertInstall&&!isCodeCreate&&!isDiffText)||isMCP" :src="Url"
                style="width:  calc(100% - 5px);height:  calc(100% - 25px);">
        </iframe>
        <div v-show="isCertInstall"
             style="width:  100%;height:  calc(100% - 45px);position: relative;justify-content: center;align-content: center;text-align: center">
          <CertInstall>
          </CertInstall>
        </div>
        <div v-show="isDiffText"
             style="width:  100%;height:  calc(100% - 30px);position: relative;justify-content: center;align-content: center;text-align: center">
          <DiffText>
          </DiffText>
        </div>
        <div v-show="isCodeCreate"
             style="width:  100%;height:  calc(100% - 30px);position: relative;justify-content: center;align-content: center;text-align: center">
          <CodeCreate>
          </CodeCreate>
        </div>
      </div>
    </div>
  </div>
</template>
<style>
.white-svg path {
  stroke: white;
}
</style>
