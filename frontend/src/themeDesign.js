import {createApp} from 'vue'
import App from './components/Tools/Theme/themeDesign.vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import {GOOS} from "../bindings/changeme/Service/appmain"; // 引入中文语言包
const app = createApp(App)

import { setLocaleData } from 'monaco-editor-nls';
import zh_CN from 'monaco-editor-nls/locale/zh-hans';

setLocaleData(zh_CN);

function SetMarginTop() {
    const aa = document.getElementsByClassName("css-1kyqri8")
    if (aa.length < 1) {
        requestAnimationFrame(() => {
            SetMarginTop()
        })
        return
    }
    aa[0].style = "margin-top:30px"
}

localStorage.setItem("theme-builder.atom.config.expandedEditors", ["General", "All Parameters", "所有参数"]);
GOOS().then(isWindows => {
    if (!isWindows) {
        SetMarginTop()
    }
})
app.use(ElementPlus, {locale: zhCn})
app.mount('#app')
