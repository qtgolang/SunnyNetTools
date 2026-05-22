let initialized = false;
let initPromise = null;

export function setupMonaco() {
    if (initialized) {
        return Promise.resolve();
    }
    if (initPromise) {
        return initPromise;
    }
    initPromise = (async () => {
        const {setLocaleData} = await import("monaco-editor-nls");
        const zhCN = await import("monaco-editor-nls/locale/zh-hans");
        setLocaleData(zhCN.default ?? zhCN);

        const workerUrl = (moduleId, label) => {
            if (label === "editorWorkerService" || moduleId === "workerMain.js") {
                return "./vs/base/worker/workerMain.js";
            }
            return "";
        };
        window.MonacoEnvironment = {getWorkerUrl: workerUrl};
        self.MonacoEnvironment = {getWorkerUrl: workerUrl};
        initialized = true;
    })();
    return initPromise;
}
