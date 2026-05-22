import {createApp} from 'vue'
import ReplaceBody from './components/SideBar/Settings/ReplaceBody.vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'
import zhCn from 'element-plus/es/locale/lang/zh-cn' // 引入中文语言包
import {AllEnterpriseModule, IntegratedChartsModule, LicenseManager} from "ag-grid-enterprise";
import {AgChartsEnterpriseModule} from "ag-charts-enterprise";
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

window.themes = [
    {id: "themeQuartz", theme: themeQuartz},
    {id: "themeBalham", theme: themeBalham},
    {id: "themeMaterial", theme: themeMaterial},
    {id: "themeAlpine", theme: themeAlpine},
];

ModuleRegistry.registerModules([AllCommunityModule, IntegratedChartsModule.with(AgChartsEnterpriseModule), AllEnterpriseModule]);
LicenseManager.setLicenseKey(LicenseKey);

ModuleRegistry.registerModules([AllEnterpriseModule]);

// 配置 MonacoEnvironment
window.MonacoEnvironment = {
    getWorkerUrl: function (moduleId, label) {
        if (label === 'editorWorkerService') {
            return '/node_modules/monaco-editor/min/vs/base/worker/workerMain.js';
        }
        console.log("未配置 MonacoEnvironment getWorkerUrl ", moduleId, label)
        return '';
    }
};

document.addEventListener("mousemove", function (event) {
    window.mouseX = event.clientX;
    window.mouseY = event.clientY;
});

const app = createApp(ReplaceBody)
app.use(ElementPlus, {locale: zhCn})
app.mount('#app')
