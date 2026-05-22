//内置常用工具
import {ref, watch} from "vue";
import {GetKeys, SetIsEditKeyDown, SetKeys} from "../../../bindings/changeme/Service/appmain";
import {Config_GOOS_IsWindows} from "./Config";

export const Keys_System_id_Boss = "Boss" //全部放行
export const Keys_System_id_ALL_Release = "system_0" //全部放行
export const Keys_System_id_Current_Release = "system_1"//当前请求放行
export const Keys_System_id_Cancel_IE_Agent = "system_2"//设置/取消IE代理
export const Keys_System_id_Keys_Clear_Al = "system_3"//清空全部记录
export const Keys_System_id_ResendRequest = "system_4"//重发请求
export const Keys_Tools_DebugTools = "tools_system_3"//调试工具
export const Keys_Tools_Theme_Color = "tools_system_4"//主题调色
export const Keys_Tools_Color_textMark = "tools_system_5"//主题调色
export const Keys_Tools_Cert = "tools_system_7"//证书安装
export const Keys_Tools_Export_Cert = "tools_system_10"//导出证书
export const Keys_Tools_Code_Create = "tools_system_11"//代码生成
export const Keys_Tools_Diff_Text = "tools_system_12"//代码生成

const systemKeys = [
    {
        Name: '老板键',
        ctrlKey: false,
        altKey: true,
        shiftKey: false,
        metaKey: false,
        keyCode: 0x51,
        key: "Q",
        value: "ALT + Q",
        ID: Keys_System_id_Boss
    },
    {
        Name: '全部放行',
        ctrlKey: false,
        altKey: false,
        shiftKey: true,
        metaKey: false,
        keyCode: 0x5A,
        key: "Z",
        value: "Shift + Z",
        ID: Keys_System_id_ALL_Release
    },
    {
        Name: '放行当前请求',
        ctrlKey: true,
        altKey: false,
        shiftKey: false,
        metaKey: false,
        key: "Z",
        keyCode: 0x5A,
        value: "CTRL + Z",
        ID: Keys_System_id_Current_Release
    },
    {
        Name: '设置/取消IE代理',
        ctrlKey: false,
        altKey: false,
        shiftKey: false,
        metaKey: false,
        key: "F12",
        keyCode: 0x7B,
        value: "F12",
        ID: Keys_System_id_Cancel_IE_Agent
    },
    {
        Name: '清空全部记录',
        ctrlKey: true,
        altKey: true,
        shiftKey: false,
        metaKey: false,
        key: "X",
        keyCode: 0x58,
        value: "CTRL + ALT + X",
        ID: Keys_System_id_Keys_Clear_Al
    },
    {
        Name: '批量重发',
        ctrlKey: true,
        altKey: false,
        shiftKey: false,
        metaKey: false,
        key: "R",
        keyCode: 0x52,
        value: "CTRL + R",
        ID: Keys_System_id_ResendRequest
    },
];
export const Config_Tools_SystemList = [
    {Name: "调试工具", Windows: false, thisOpen: false, register: false, ID: Keys_Tools_DebugTools},
    {Name: "主题调色", Windows: false, thisOpen: false, register: false, ID: Keys_Tools_Theme_Color},
    {Name: "颜色标记", Windows: false, thisOpen: true, register: false, ID: Keys_Tools_Color_textMark},
    {Name: "导出证书", Windows: false, thisOpen: false, register: false, ID: Keys_Tools_Export_Cert, noKeys: true},
    {Name: "代码生成", Windows: false, thisOpen: false, register: false, ID: Keys_Tools_Code_Create, noKeys: true},
    {Name: "证书安装", Windows: false, thisOpen: false, register: false, ID: Keys_Tools_Cert, noKeys: true},//noKeys 不设置快捷键
    {Name: "文本对比", Windows: false, thisOpen: false, register: false, ID: Keys_Tools_Diff_Text},
];

function newKeys(Name, ID) {
    return {
        Name: Name,
        ctrlKey: false,
        altKey: false,
        shiftKey: false,
        metaKey: false,
        key: "",
        value: "",
        ID: ID
    }
}

function checkKeys(obj) {
    if (obj.ctrlKey === undefined || obj.ctrlKey === null) {
        obj.ctrlKey = false
    }
    if (obj.altKey === undefined || obj.altKey === null) {
        obj.altKey = false
    }
    if (obj.shiftKey === undefined || obj.shiftKey === null) {
        obj.shiftKey = false
    }
    if (obj.metaKey === undefined || obj.metaKey === null) {
        obj.metaKey = false
    }
    if (obj.key === undefined || obj.key === null) {
        obj.key = ""
    }
    if (obj.keyCode === undefined || obj.keyCode === null) {
        obj.keyCode = 0
    }
    if (obj.value === undefined || obj.value === null) {
        obj.value = ""
    }
}

export const Config_Tools_CustomList = ref([]);
//快捷键信息
export const Config_Keys = new Map();
export const Config_Keys_ID = ref([]);

watch(Config_Tools_CustomList, (Array) => {

    for (let i = 0; i < Array.length; i++) {
        if (Config_Keys.has(Array[i].ID) || Array[i].noKeys === true || Array[i].thisOpen === true) {
            const row = Config_Keys.get(Array[i].ID)
            if (row) {
                if (row.Name !== Array[i].Name) {
                    row.Name = Array[i].Name
                }
            }
            continue
        }
        Config_Keys.set(Array[i].ID, newKeys(Array[i].Name, Array[i].ID))
    }
    const array = []
    Config_Keys.forEach((value, key) => {
        checkKeys(value)
        array.push(key)
    })
    Config_Keys_ID.value = array


}, {deep: true});


export function RestKeys() {
    console.log("RestKeys")
    Config_Keys.clear()
    Config_Keys_ID.value = []
    Config_Tools_CustomList.value = []
    for (let i = 0; i < systemKeys.length; i++) {
        Config_Keys.set(systemKeys[i].ID, systemKeys[i])
    }
    for (let i = 0; i < Config_Tools_SystemList.length; i++) {
        if (Config_Tools_SystemList[i].noKeys === true || Config_Tools_SystemList[i].thisOpen === true) {
            continue
        }
        Config_Keys.set(Config_Tools_SystemList[i].ID, newKeys(Config_Tools_SystemList[i].Name, Config_Tools_SystemList[i].ID))
    }
    for (let i = 0; i < Config_Tools_CustomList.value.length; i++) {
        if (Config_Tools_CustomList[i].noKeys === true || Config_Tools_CustomList[i].thisOpen === true) {
            continue
        }
        Config_Keys.set(Config_Tools_CustomList[i].ID, newKeys(Config_Tools_CustomList.value[i].Name, Config_Tools_CustomList.value[i].ID))
    }
    const array = []
    Config_Keys.forEach((value, key) => {
        checkKeys(value)
        array.push(key)
    })
    Config_Keys_ID.value = array
    SetKeys(JSON.stringify(GetKeysArray()))
}

export function GetKeysArray() {
    const array = []
    Config_Keys.forEach((value, key) => {
        checkKeys(value)
        array.push(value)
    })
    return array
}

function parseKeys(Array) {
    Config_Keys.clear()
    for (let i = 0; i < Array.length; i++) {
        if (Config_Keys.has(Array[i].ID)) {
            Config_Keys.delete(Array[i].ID)
        }
        Config_Keys.set(Array[i].ID, Array[i])
    }
    for (let i = 0; i < systemKeys.length; i++) {
        if (!Config_Keys.has(systemKeys[i].ID)) {
            Config_Keys.set(Array[i].ID, Array[i])
        }
    }
    for (let i = 0; i < Config_Tools_SystemList.length; i++) {
        if (!Config_Keys.has(Config_Tools_SystemList[i].ID) && Config_Tools_SystemList[i].noKeys !== true) {
            Config_Keys.set(Config_Tools_SystemList[i].ID, {
                Name: Config_Tools_SystemList[i].Name,
                ctrlKey: false,
                altKey: false,
                shiftKey: false,
                metaKey: false,
                key: "",
                value: "",
                ID: Config_Tools_SystemList[i].ID
            })
        }
        if (Config_Keys.has(Config_Tools_SystemList[i].ID) && Config_Tools_SystemList[i].thisOpen === true) {
            Config_Keys.delete(Config_Tools_SystemList[i].ID)
        }
    }
    for (let i = 0; i < Config_Tools_CustomList.value.length; i++) {
        if (!Config_Keys.has(Config_Tools_CustomList.value[i].ID) && Config_Tools_CustomList[i].noKeys !== true) {
            Config_Keys.set(Config_Tools_CustomList.value[i].ID, {
                Name: Config_Tools_CustomList.value[i].Name,
                ctrlKey: false,
                altKey: false,
                shiftKey: false,
                metaKey: false,
                key: "",
                value: "",
                ID: Config_Tools_CustomList.value[i].ID
            })
        }
        if (Config_Keys.has(Config_Tools_CustomList[i].ID) && Config_Tools_CustomList[i].thisOpen === true) {
            Config_Keys.delete(Config_Tools_CustomList[i].ID)
        }

    }
    const array = []
    Config_Keys.forEach((value, key) => {
        checkKeys(value)
        array.push(key)
    })
    Config_Keys_ID.value = array
}

export const Tools_CustomList_init_Func = ref(null);

function initKeys() {
    GetKeys().then(res => {
        if (res.length < 10) {
            RestKeys()
            return
        }
        try {
            parseKeys(JSON.parse(res))
        } catch (e) {
            RestKeys()
        }
    })
}

initKeys()

const Hotkey = new Map();

//注册快捷键回调函数
export function registerHotkeyFunction(id, func) {
    if (!Hotkey.has(id)) {
        Hotkey.set(id, []);
    }
    Hotkey.get(id).push(func);
}

export function keydownEventListener(event) {
    if (window.isEditKeyDown.value) {
        return;
    }
    if (event.key.toUpperCase() === "BACKSPACE") {
        return
    }
    const row = keydownEventToString(event)
    Config_Keys.forEach((item) => {
        if (item.key.toUpperCase() === row.key.toUpperCase() && item.altKey === row.altKey && item.ctrlKey === row.ctrlKey && item.shiftKey === row.shiftKey) {
            /*
            if (Config_GOOS_IsWindows.value) {
                if (!item.ID.startsWith("system_")) {
                    return
                }
            }
            */
            const array = Hotkey.get(item.ID)
            if (array) {
                array.forEach((func) => {
                    func(event);
                });
            }
        }
    })
}

export function ExternalKeydownEventListener(id) {
    const array = Hotkey.get(id)
    if (array) {
        array.forEach((func) => {
            func();
        });
    }
}

export function keydownEventToString(event) {
    let mKey = event.key
    if (mKey.startsWith("KEY")) {
        mKey = mKey.substring(3);
    }
    const row = {};
    row.altKey = event.altKey
    row.ctrlKey = event.ctrlKey
    row.shiftKey = event.shiftKey
    row.metaKey = event.metaKey
    row.key = mKey
    row.value = ""
    row.keyCode = event.keyCode
    let obj = ""
    let ControlCtrl = event.key === 'Control' || event.code === 'ControlLeft' || event.code === 'ControlRight'
    let ControlShift = event.key === 'Shift' || event.code === 'ShiftLeft' || event.code === 'ShiftRight'
    let ControlAlt = event.key === 'Alt' || event.code === 'AltLeft' || event.code === 'AltRight'
    if (Config_GOOS_IsWindows.value) {
        {
            if (row.ctrlKey) {
                obj = "CTRL"
            }
            if (row.altKey) {
                if (obj === "") {
                    obj = "ALT"
                } else {
                    obj += " + ALT"
                }
            }
            if (row.shiftKey) {
                if (obj === "") {
                    obj = "Shift"
                } else {
                    obj += " + Shift"
                }
            }
            if (!ControlCtrl && !ControlAlt && !ControlShift) {
                if (row.key) {
                    if (obj === "") {
                        obj = row.key
                    } else {
                        obj += " + " + row.key
                    }
                }
            }
            row.value = obj
        }
    } else {
        {
            if (row.ctrlKey) {
                obj = "Control"
            }
            if (row.altKey) {
                if (obj === "") {
                    obj = "Option"
                } else {
                    obj += " + Option"
                }
            }
            if (row.shiftKey) {
                if (obj === "") {
                    obj = "Shift"
                } else {
                    obj += " + Shift"
                }
            }
            if (!ControlCtrl && !ControlAlt && !ControlShift) {
                if (row.key) {
                    if (obj === "") {
                        obj = row.key
                    } else {
                        obj += " + " + row.key
                    }
                }
            }
        }
    }
    if (obj.toUpperCase() === "BACKSPACE") {
        row.value = ""
        return ""
    }
    row.value = obj
    return row
}

window.isEditKeyDown = ref(false);

watch(window.isEditKeyDown, (newValue, oldValue) => {
    if (newValue !== oldValue) {
        SetIsEditKeyDown(newValue).then();
    }
});