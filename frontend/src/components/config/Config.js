import {ref, watch} from "vue";
import {
    GetAgGridDarkTheme,
    GetAgGridLightTheme,
    GetHomeTextMark,
    GetListColor,
    GOOS,
    IsDark,
    SetIsDark
} from "../../../bindings/changeme/Service/appmain.js";
import {themeQuartz} from "ag-grid-community";

export let mousemoveMouseX = 0;
export let mousemoveMouseY = 0;
window.addEventListener('mousemove', (e) => {
    mousemoveMouseX = e.clientX;
    mousemoveMouseY = e.clientY;
});

export function getCurrentElement(x, y, func) {
    if (x !== mousemoveMouseX && y !== mousemoveMouseY) {
        func(document.elementFromPoint(mousemoveMouseX, mousemoveMouseY))
        return
    }
    requestAnimationFrame(() => {
        getCurrentElement(x, y, func)
    })
}

export const DropFilesEvent = (func) => {
    getCurrentElement(mousemoveMouseX, mousemoveMouseY, func)
}
//是否显示导出进度
export const Config_IsShowExportProgress = ref("");
//是否点击菜单的文件
export const Config_Menu_isFileMenu = ref(false);
//是否自动滚动
export const Config_AutoRoll = ref(false);
//当前选择的行
export const Config_SelectedRow = ref({});
//Websocket or TCP or UDP 当前选择的行
export const Config_SocketSelectedRow = ref({});
//是否为Windows
export const Config_GOOS_IsWindows = ref(false);
//是否为暗色模式
export const Config_IsDark = ref(true);
//状态信息
export const Config_Status_Info = ref("正在启动...");
//agGrid 的API对象
export const Config_agGrid_API = ref(null);

//当前列表区域右键菜单是否弹出/是否可见
export const Config_MenuVisible = ref(false);
//ag-Grid 主题
export const Config_Theme_agGrid = ref(null);

//查询窗口显示
export const Config_Find_Window_Show = ref(() => {
});
//查询窗口隐藏
export const Config_Find_Window_Hide = ref(() => {
});
//SunnyNet 是否启动
export const Config_SunnyNetIsStart = ref(false);
export const Config_Find_Window = ref(null);
export const DisableClick = ref(false);
//是否在所有区域查找 包括HTTP？
export const Config_Find_Range_ALL = ref(true);
//查询窗口配置信息
export const Config_Find = ref({
    value: "",
    action: "hide",
    Type: "UTF8",
    check: ["不区分大小写", "取消之前的颜色标记"],
    Color: "红色"
});
export const Config_HTTP_Message_free = ref(() => {
});

export const Config_Find_Type_Options = [
    {
        value: 'auto',
        label: '自动查找 [ GBK/UTF8/HEX/Base64/Number ]',
    },
    {
        value: 'UTF8',
        label: '将查找内容转为(UTF8)字符串,进行查找',
    },
    {
        value: 'GBK',
        label: '将查找内容转为( GBK )字符串,进行查找',
    },
    {
        value: 'Hex',
        label: '将查找内容进行十六进制解码后,进行查找',
    },
    {
        value: 'Base64',
        label: '将查找内容进行 Base64 解码后,进行查找',
    },
    {
        value: '整数4',
        label: '将查找内容输入的数值转为( 4字节整数 ),进行查找',
    },
    {
        value: '整数8',
        label: '将查找内容输入的数值转为( 8字节整数 ),进行查找',
    },
    {
        value: '浮点数4',
        label: '将查找内容输入的数值转为( 4字节浮点数 ),进行查找',
    },
    {
        value: '浮点数8',
        label: '将查找内容输入的数值转为( 8字节浮点数 ),进行查找',
    },
    {
        value: 'pb',
        label: '将数据进行尝试 ProtoBuf 解码,若解码成功,将查找输入的{查找内容}',
    },
]
export const Config_Find_Range_Options_All = [
    {
        value: '在全部范围寻找',
        label: '在全部范围寻找',
    },
    {
        value: 'HTTP请求',
        label: '在 HTTP/HTTPS 请求数据中寻找',
    },
    {
        value: 'HTTP响应',
        label: '在 HTTP/HTTPS 响应数据中寻找',
    },
    {
        value: 'socketSend',
        label: '在 TCP/UDP/WebSocket 发送数据中寻找',
    },
    {
        value: 'socketRec',
        label: '在 TCP/UDP/WebSocket 接收数据中寻找',
    },

    {
        value: 'socketAll',
        label: '在 TCP/UDP/WebSocket 发送/接收 数据中寻找',
    },
]
export const Config_Find_Range_Options_Websocket = [
    {
        value: '在全部范围寻找',
        label: '在全部范围寻找',
    },
    {
        value: 'socketSend',
        label: '在发送数据中寻找',
    },
    {
        value: 'socketRec',
        label: '在接收数据中寻找',
    },
]
export const Config_Focus_Element = ref(null);
export const Config_MonacoEditorTheme = ref("");
//状态信息
export const Config_IsRest = ref(0);

watch(Config_IsDark, (newValue, oldValue) => {
    if (newValue !== oldValue) {
        SetIsDark(newValue).then(() => {
            initTheme()
        })
        GetListColor(newValue).then((res) => {
            Object.entries(res).forEach(([ColorID, Color]) => {
                SetTextColor(ColorID, Color)
            });
        })
    }
});
//状态信息
export const Config_HomeTextMark = new Map();

export function Find_Type_Del_Option(obj, name) {
    const array = [];
    for (let i = 0; i < obj.length; i++) {
        if (obj[i].value !== name) {
            array.push(obj[i])
        }
    }
    return array
}


function initTheme() {
    IsDark().then((_isDark) => {
        Config_MonacoEditorTheme.value = _isDark ? 'vs-dark' : 'vs';
        document.documentElement.className = _isDark ? "dark" : "light";
        if (_isDark) {
            GetAgGridDarkTheme().then((res) => {
                setTheme(res, _isDark)
            })
        } else {
            GetAgGridLightTheme().then((res) => {
                setTheme(res, _isDark)
            })
        }
        GetListColor(_isDark).then((res) => {
            Object.entries(res).forEach(([ColorID, Color]) => {
                SetTextColor(ColorID, Color)
            });
        })
        if (Config_IsDark.value !== _isDark) {
            Config_IsDark.value = _isDark
        }
    })
    GetHomeTextMark().then((res) => {
        Config_HomeTextMark.clear()
        try {
            JSON.parse(res).forEach((item) => {
                Config_HomeTextMark.set(item.id, item)
            })
        } catch (e) {
        }
        if (Config_HomeTextMark.size < 1) {
            for (let i = 0; i < DefaultRowData.length; i++) {
                Config_HomeTextMark.set(DefaultRowData[i].id, DefaultRowData[i])
            }
        }
    })
}

export const DefaultRowData = [
    {"名称": "红色", "深色主题": "#DC2651", "浅色主题": "#F20808", id: "1"},
    {"名称": "蓝色", "深色主题": "#1578F1", "浅色主题": "#0871F2", id: "2"},
    {"名称": "绿色", "深色主题": "#0F8405", "浅色主题": "#12740A", id: "3"},
    {"名称": "黄色", "深色主题": "#C69005", "浅色主题": "#efbe47", id: "4"},
    {"名称": "紫色", "深色主题": "#47056D", "浅色主题": "#5F02A6", id: "5"},
]

export function setTheme(obj, isDark) {
    window.temp = (new Function("return " + obj))();
    Config_Theme_agGrid.value = themeQuartz.withParams(window['temp']);
}

async function initGOOS() {
    Config_GOOS_IsWindows.value = await GOOS()
}


initGOOS().then(() => {
    initTheme();
})

window.Clipboard_SetText = function (text, func) {
    if (window.__Clipboard_SetText) {
        window.__Clipboard_SetText(text, func)
    }
}

export function getSecureRandomString(length) {
    const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
    const array = new Uint8Array(length);
    crypto.getRandomValues(array);
    return Array.from(array, (byte) => chars[byte % chars.length]).join("");
}

export const Config_TextColorID1 = ref(""); // TCP请求
export const Config_TextColorID2 = ref(""); // UDP请求
export const Config_TextColorID3 = ref(""); // Websocket
export const Config_TextColorID4 = ref(""); // CSS
export const Config_TextColorID5 = ref(""); // JavaScript
export const Config_TextColorID6 = ref(""); // 图片
export const Config_TextColorID7 = ref(""); // 文档
export const Config_TextColorID8 = ref(""); // 错误请求
export const Config_TextColorID9 = ref(""); // 重定向请求
export const Config_TextColorID10 = ref(""); // 40x 错误
export const Config_TextColorID99 = ref(""); // 常规项

const TextColorMap = new Map();
const configMap = {
    "1": Config_TextColorID1,
    "2": Config_TextColorID2,
    "3": Config_TextColorID3,
    "4": Config_TextColorID4,
    "5": Config_TextColorID5,
    "6": Config_TextColorID6,
    "7": Config_TextColorID7,
    "8": Config_TextColorID8,
    "9": Config_TextColorID9,
    "10": Config_TextColorID10,
    "99": Config_TextColorID99
};
Object.entries(configMap).forEach(([colorID, refValue]) => {
    watch(refValue, (newValue) => {
        try {
            const api = Config_agGrid_API.value;
            const array = [];
            for (const [key, value] of TextColorMap) {
                if (value.toString() === colorID) {
                    const node = api.getRowNode(key);
                    if (node) {
                        node.data.color = newValue;
                        array.push(node.data);
                    }
                }
            }
            api.applyTransaction({update: array});
        } catch (e) {
        }

    });
});
// 统一监听所有颜色 ID 变化


// **状态码映射表**
const statusToColorID = {
    "301": 9, "302": 9, // 重定向
    "401": 10, "403": 10, "404": 10, "500": 10, // 40x 错误
    "-1": 8, "错误": 8 // 错误请求
};

// **响应类型映射表**
const typeToColorID = [
    {key: "/css", id: 4},
    {key: "text/", id: 7},
    {key: "image/", id: 6},
    {key: "javascript", id: 5}
];

// **请求方式映射表**
const modeToColorID = [
    {key: "websocket", id: 3},
    {key: "udp", id: 2},
    {key: "tcp", id: 1}
];

// **获取文本颜色**
export function GetTextColor(node) {
    const status = String(node["状态"]);
    if (statusToColorID[status]) {
        node.color = configMap[statusToColorID[status]].value;
        setTextColorMap(node, statusToColorID[status]);
        return;
    }

    const responseType = String(node["响应类型"]).toLowerCase();
    for (const {key, id} of typeToColorID) {
        if (responseType.includes(key)) {
            node.color = configMap[id].value;
            setTextColorMap(node, id);
            return;
        }
    }

    const requestMode = String(node["方式"]).toLowerCase();
    for (const {key, id} of modeToColorID) {
        if (requestMode.includes(key)) {
            node.color = configMap[id].value;
            setTextColorMap(node, id);
            return;
        }
    }

    //常规项,不满足任何一项
    node.color = configMap["99"].value;
    setTextColorMap(node, 99);
}

// **设置颜色映射**
function setTextColorMap(node, id) {
    TextColorMap.set(String(node["Theology"]), id);
}

// **清空颜色映射**
export function clearTextColorMap() {
    TextColorMap.clear();
}

// **删除某个颜色映射**
export function deleteTextColorMap(Theology) {
    TextColorMap.delete(String(Theology));
}

// **设置文本颜色**
export function SetTextColor(ColorID, Color) {
    if (configMap[ColorID]) {
        configMap[ColorID].value = Color;
    }
}

export function ObjString(o) {
    if (o === undefined) {
        return ""
    }
    if (o === null) {
        return ""
    }
    return String(o)
}

const thisObjectMap = new Map();

export function registerThisObject(id, func) {
    deleteThisObject(id);
    thisObjectMap.set(id, func);
}

export function deleteThisObject(id) {
    if (!thisObjectMap.has(id)) {
        return
    }
    thisObjectMap.delete(id)
}

export function getThisObject(id) {
    if (!thisObjectMap.has(id)) {
        return () => {
        }
    }
    return thisObjectMap.get(id)
}