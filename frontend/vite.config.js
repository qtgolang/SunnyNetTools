import {defineConfig} from "vite";
import vue from "@vitejs/plugin-vue";
import wails from "@wailsio/runtime/plugins/vite";

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [vue(), wails("./bindings")],
    build: {
        rollupOptions: {
            input: {
                main: 'index.html',
                replace: 'ReplaceBody.html',
                other: 'Other.html',
                debugTools: 'debugTools.html',
                cert: 'Cert.html',
                theme: 'Theme.html',
                themeDesign: 'ThemeDesign.html'
            }
        }
    }
});
