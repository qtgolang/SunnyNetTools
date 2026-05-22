import {createApp} from 'vue'
import App from './App.vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import ContextMenu from "@imengyu/vue3-context-menu"
import * as ElementPlusIconsVue from "@element-plus/icons-vue";
import '@imengyu/vue3-context-menu/lib/vue3-context-menu.css'
import {AllEnterpriseModule, IntegratedChartsModule, LicenseManager} from "ag-grid-enterprise";
import {AgChartsEnterpriseModule} from "ag-charts-enterprise";

import TitlebarComponent from "./components/TitleBar/VUETitlebar/vueTitlebar.vue"
import {
    AllCommunityModule,
    ModuleRegistry,
    themeAlpine,
    themeBalham,
    themeMaterial,
    themeQuartz,
} from "ag-grid-community";

import { setLocaleData } from 'monaco-editor-nls';
import zh_CN from 'monaco-editor-nls/locale/zh-hans';

setLocaleData(zh_CN);
import {LicenseKey} from "./AGLicenseKey";
/*
{
    console.log = (...args) => {
    };
    console.error = (...args) => {
    };
    console.warn = (...args) => {
    };
    window.onerror = (message, source, lineno, colno, error) => {
        return true; // 阻止错误向上冒泡
    };

}
*/
window.themes = [
    {id: "themeQuartz", theme: themeQuartz},
    {id: "themeBalham", theme: themeBalham},
    {id: "themeMaterial", theme: themeMaterial},
    {id: "themeAlpine", theme: themeAlpine},
];

ModuleRegistry.registerModules([AllCommunityModule, IntegratedChartsModule.with(AgChartsEnterpriseModule), AllEnterpriseModule]);

LicenseManager.setLicenseKey(LicenseKey);
 

{
// 配置 MonacoEnvironment
    const Worker = (moduleId, label) => {
        if (label === 'editorWorkerService') {
            return './vs/base/worker/workerMain.js';
        }
        if (moduleId === 'workerMain.js') {
            return './vs/base/worker/workerMain.js';
        }
        console.log("未配置 MonacoEnvironment getWorkerUrl ", moduleId, label)
        return '';
    }
    window.MonacoEnvironment = {
        getWorkerUrl: Worker
    };
    self.MonacoEnvironment = {
        getWorkerUrl: Worker
    };
}
document.addEventListener("mousemove", function (event) {
    window.mouseX = event.clientX;
    window.mouseY = event.clientY;
});

const app = createApp(App)
app.component(TitlebarComponent.name, TitlebarComponent);
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
}
app.use(ContextMenu)
app.use(ElementPlus, {locale: zhCn})
app.mount('#app')
