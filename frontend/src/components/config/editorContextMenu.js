import {nextTick} from 'vue';
import {ElMessageBox} from "element-plus";
import {AppGetEditorFontSize, AppSetEditorFontSize, ClipboardReadAll} from "../../../bindings/changeme/Service/appmain";
import {Events} from "@wailsio/runtime";
import * as monaco from 'monaco-editor'
import {Config_GOOS_IsWindows} from "./Config";

export function addMonacoEditorWordWrapMenu(editor, defaultOpts) {
    const key = editor.createContextKey('wordWrapOn', true)
    // 默认换行的动作
    editor.addAction({
        id: 'turnWordWrapOff',
        label: '关闭自动换行',
        contextMenuGroupId: 'my-commands',
        contextMenuOrder: Number.MAX_SAFE_INTEGER,
        precondition: 'wordWrapOn',
        run: () => {
            defaultOpts.wordWrap = 'off'
            editor.updateOptions({
                wordWrap: "off"
            });
            key.set(false)
        },
    })
    editor.addAction({
        id: 'turnWordWrapOn',
        label: '自动换行',
        contextMenuGroupId: 'my-commands',
        contextMenuOrder: Number.MAX_SAFE_INTEGER,
        precondition: '!wordWrapOn',
        run: () => {
            defaultOpts.wordWrap = 'on'
            editor.updateOptions({
                wordWrap: "on"
            });
            key.set(true)
        },
    })
    editor.addAction({
        id: 'presentData',
        label: '展示完整数据',
        contextMenuGroupId: 'my-presentData',
        contextMenuOrder: Number.MAX_SAFE_INTEGER,
        run: () => {
            if (defaultOpts.presentData) {
                defaultOpts.presentData()
            }
        },
    })
}

let FontSize = 18

export function addMonacoEditorFontSize(editor) {
    editor.addAction({
        id: 'editor.action.clipboardPasteAction', // 👈 用相同 ID 覆盖
        label: '粘贴　',
        contextMenuGroupId: '9_cutcopypaste',
        contextMenuOrder: 4,
        run: async () => {
            try {
                const text = await ClipboardReadAll()
                const pos = editor.getPosition()
                editor.executeEdits('', [{
                    range: new monaco.Range(pos.lineNumber, pos.column, pos.lineNumber, pos.column),
                    text,
                    forceMoveMarkers: true
                }])
            } catch (err) {
                console.error('发生错误：', err);
                console.error(err.stack);  // 打印详细的堆栈信息
            }
        },
    });
    editor.addAction({
        id: 'FontSize',
        label: '设置字体大小',
        contextMenuGroupId: 'my-FontSize',
        contextMenuOrder: Number.MAX_SAFE_INTEGER,
        run: () => {
            ElMessageBox.prompt('请输入字体大小（仅限数字）', '设置字体大小', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                inputPattern: /^\d+$/, // 只允许输入数字
                inputErrorMessage: '请输入有效的数字',
                inputValue: FontSize, // 默认值
                customClass: 'centered-message-box', // 自定义类名用于居中内容
            }).then(({value}) => {
                AppSetEditorFontSize(parseInt(value))
            })
        },
    })
    AppGetEditorFontSize().then(res => {
        FontSize = res
        editor.updateOptions({fontSize: FontSize});
    })
    Events.On("updateEditorFontSize", (value) => {
        let o=value.data[0];
        if (!o){
            o=value.data;
        }
        FontSize =o
        editor.updateOptions({fontSize: FontSize});
    })
}

export function RemoveMonacoEditorMenu(editor, isRemoveMenu, readOnly) {
    editor.onContextMenu(() => {
        nextTick(() => {
            const shadowHosts = document.querySelectorAll(".shadow-root-host");
            for (let i = 0; i < shadowHosts.length; i++) {
                const shadowHost = shadowHosts[i];
                if (shadowHost !== null && shadowHost !== undefined) {
                    const actionsContainer = shadowHost.shadowRoot.querySelectorAll(".actions-container")[0];
                    if (actionsContainer) {
                        if (Config_GOOS_IsWindows.value) {
                            const RequestDIV = document.getElementsByClassName("RequestDIV")[0]
                            const ResponseDIV = document.getElementsByClassName("ResponseDIV")[0]
                            if (RequestDIV && ResponseDIV) {
                                if (RequestDIV.contains(shadowHost)) {
                                    RequestDIV.style.zIndex = 2;
                                    ResponseDIV.style.zIndex = 1;
                                } else {
                                    RequestDIV.style.zIndex = 1;
                                    ResponseDIV.style.zIndex = 2;
                                }
                            }
                        }

                        const childNodes = actionsContainer.childNodes;
                        let hasRemoved = false;
                        let last = ""
                        while (true) {
                            for (let j = 0; j < childNodes.length; j++) {
                                const nodeValue = childNodes[j].innerHTML;
                                if (nodeValue.includes("粘贴") && !nodeValue.includes("粘贴　")) {
                                    childNodes[j].parentNode.removeChild(childNodes[j]);
                                    hasRemoved = true;
                                }
                                if (nodeValue.includes("转到") || nodeValue.includes("快速查看")|| nodeValue.includes("重命名符号")) {
                                    childNodes[j].parentNode.removeChild(childNodes[j]);
                                    hasRemoved = true;
                                }
                                if (readOnly()) {
                                    if (nodeValue.includes("粘贴　")) {
                                        childNodes[j].parentNode.removeChild(childNodes[j]);
                                        hasRemoved = true;
                                    }
                                }
                                if (j === 0 || j === childNodes.length - 1 || nodeValue === last) {
                                    if (nodeValue.indexOf("action-label codicon separator disabled") !== -1) {
                                        childNodes[j].parentNode.removeChild(childNodes[j]);
                                        hasRemoved = true;
                                        break;
                                    }
                                }

                                if (nodeValue.indexOf("action-label codicon separator disabled") !== -1) {
                                    last = nodeValue
                                } else {
                                    last = ""
                                }
                                if (isRemoveMenu(nodeValue, readOnly())) {
                                    childNodes[j].parentNode.removeChild(childNodes[j]);
                                    hasRemoved = true;
                                    break;
                                }
                            }
                            if (hasRemoved === true) {
                                hasRemoved = false;
                                continue;
                            }
                            break;
                        }
                        const menuContainer = shadowHost.shadowRoot.querySelectorAll(".monaco-menu-container")[0];
                        if (menuContainer) {
                            let mouseX = event.clientX + window.scrollX;
                            let mouseY = event.clientY + window.scrollY;
                            if (mouseY + menuContainer.offsetHeight > window.innerHeight) {
                                mouseY = window.innerHeight - menuContainer.offsetHeight;
                            }
                            if (mouseX + menuContainer.offsetWidth - 100 > window.innerWidth) {
                                mouseX = window.innerWidth - (menuContainer.offsetWidth - 100);
                            }
                            menuContainer.style.top = mouseY + "px";
                            menuContainer.style.left = mouseX + "px";
                        }
                    }
                }
            }
        }).then(r => {
            // 处理后续逻辑（如果需要的话）
        });
    });

}

function initMonacoEditorMenuEventListener() {
    document.addEventListener("mousemove", (event) => {
        const shadowHosts = document.querySelectorAll(".shadow-root-host");
        for (let i = 0; i < shadowHosts.length; i++) {
            const shadowHost = shadowHosts[i];
            if (shadowHost !== null && shadowHost !== undefined) {
                const actionContainers = shadowHost.shadowRoot.querySelectorAll(".actions-container");
                if (actionContainers.length < 1) {
                    break;
                }
                const shadowElementsUnderCursor = shadowHost.shadowRoot.elementsFromPoint(event.clientX, event.clientY);
                const targetActionItem = shadowElementsUnderCursor.find(el => el.classList.contains("action-item"));

                const childNodes = actionContainers[0].childNodes;
                childNodes.forEach((node) => {
                    if (node.nodeType === Node.ELEMENT_NODE) {
                        if (node.classList.contains('action-item') && node.classList.contains('focused')) {
                            node.className = "action-item";
                            const nodeChildren = node.childNodes;
                            if (nodeChildren.length > 0) {
                                nodeChildren[0].style = "color: var(--vscode-menu-foreground);";
                            }
                        }
                    }
                });

                if (targetActionItem) {
                    if (targetActionItem.classList.contains('action-item') &&
                        !targetActionItem.classList.contains('focused') &&
                        !targetActionItem.classList.contains('disabled')) {
                        targetActionItem.className = "action-item focused";
                        const targetChildren = targetActionItem.childNodes;
                        if (targetChildren.length > 0) {
                            targetChildren[0].style = "color: var(--vscode-menu-selectionForeground); background-color: var(--vscode-menu-selectionBackground); outline: 1px solid var(--vscode-menu-selectionBorder); outline-offset: -1px;";
                        }
                    }
                }
            }
        }

    });
}

initMonacoEditorMenuEventListener()
