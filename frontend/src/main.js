import {createApp} from 'vue'
import App from './App.vue'
import './style.css';
import {ClientSideRowModelModule} from '@ag-grid-community/client-side-row-model';
import {ColumnsToolPanelModule} from '@ag-grid-enterprise/column-tool-panel';
import {MenuModule} from '@ag-grid-enterprise/menu';
import {GridChartsModule} from '@ag-grid-enterprise/charts';
import {ModuleRegistry} from '@ag-grid-community/core';
import {SideBarModule} from '@ag-grid-enterprise/side-bar';
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'

import * as ElementPlusIconsVue from '@element-plus/icons-vue'

// Register shared Modules globally
ModuleRegistry.registerModules([
    ClientSideRowModelModule,
    MenuModule,
    GridChartsModule,
    SideBarModule,
    ColumnsToolPanelModule,
]);
let app = createApp(App)
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
}

document.addEventListener("mousemove", function (event) {
    window.mouseX = event.clientX;
    window.mouseY = event.clientY;
});
//禁止出现滚动条
app.directive('no-scroll', {
    mounted(el) {
        const overflow = document.body.style.overflow;
        el.cacheStyle = overflow;
        document.body.style.overflow = 'hidden';
    },
    unmounted(el) {
        document.body.style.overflow = el.cacheStyle || '';
    }
})
app.use(ElementPlus)

app.mount('#app')
