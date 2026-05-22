import {createApp} from 'vue'
import App from './components/Tools/Debug/DebugTools.vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'

import zhCn from 'element-plus/es/locale/lang/zh-cn' // 引入中文语言包

import { setLocaleData } from 'monaco-editor-nls';
import zh_CN from 'monaco-editor-nls/locale/zh-hans';

setLocaleData(zh_CN);

const app = createApp(App)
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
app.use(ElementPlus, {locale: zhCn})
app.mount('#app')
