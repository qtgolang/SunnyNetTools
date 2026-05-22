import {Events} from "@wailsio/runtime";
import {
    AppGenerateCodeInterface,
    AppGetGenerateCodeInterface,
    ClipboardWriteAll
} from "../../../bindings/changeme/Service/appmain";
import {ElNotification} from "element-plus";
import {base64ToBytes} from "./encoding";

export const GenerateCodeType_HTTP = "HTTP"
export const GenerateCodeType_Websocket = "Websocket"
export const GenerateCodeType_TCP = "TCP"

const GenerateCodeListMenu = {
    "GoLang": [
        {Name: "net/http", Type: GenerateCodeType_HTTP},
        {Name: "SunnyNet/http", Type: GenerateCodeType_HTTP},
        {Name: "Websocket-请求", Type: GenerateCodeType_Websocket},
        {Name: "tcp-请求", Type: GenerateCodeType_TCP},
    ],
    "C#": [
        {Name: "HttpClient", Type: GenerateCodeType_HTTP},
        {Name: "RestSharp", Type: GenerateCodeType_HTTP},
    ],
    "Python": [
        {Name: "requests", Type: GenerateCodeType_HTTP},
    ],
    "火山": [
        {Name: "WinHttpW", Type: GenerateCodeType_HTTP},
    ],
    "易语言": [
        {Name: "网页_访问", Type: GenerateCodeType_HTTP},
        {Name: "网页_访问_对象", Type: GenerateCodeType_HTTP},
        {Name: "E2EE网站客户端", Type: GenerateCodeType_HTTP},
        {Name: "WinHttpW", Type: GenerateCodeType_HTTP},
        {Name: "WinInet", Type: GenerateCodeType_HTTP},
        {Name: "WinHttpR", Type: GenerateCodeType_HTTP},
        {Name: "SunnyHTTP", Type: GenerateCodeType_HTTP},
        {Name: "SunnyWS客户端-异步", Type: GenerateCodeType_Websocket},
        {Name: "SunnyWS客户端-同步", Type: GenerateCodeType_Websocket},
        {Name: "SunnyTCP客户端-异步", Type: GenerateCodeType_TCP},
        {Name: "SunnyTCP客户端-同步", Type: GenerateCodeType_TCP},
    ],
}
window.CodeInterface = []


export function GetGenerateCodeListDebugToolsMenu(func, codeDebugToolsFunc) {
    const array = []
    const isHttp = true

    Object.keys(GenerateCodeListMenu).forEach(key => {
        const obj = GenerateCodeListMenu[key];
        const arr = [];
        obj.forEach(value => {
            if ((isHttp && value.Type === GenerateCodeType_HTTP)) {
                arr.push({
                    value: value.Name,
                    label: value.Name,
                    action: () => {
                        func(key, value.Name)
                    },
                })
            }
        })
        if (arr.length > 0) {
            array.push({
                value: key,
                label: key,
                children: arr,
            })
        }
    });
    const array2 = []
    const array3 = []
    Object.keys(window.CodeInterface).forEach(Language => {
        const _Module = window.CodeInterface[Language]
        try {
            if (_Module.func) {
                for (let i = 0; i < array2.length; i++) {
                    const element = array2[i];
                    if (element) {
                        if (element.value === Language) {
                            return
                        }
                    }
                }
                if (_Module.isHTTP === isHttp) {
                    array3.push({
                        value: Language,
                        label: Language,
                        action: () => {
                            codeDebugToolsFunc(_Module.func)
                        },
                    })
                }
                return
            }
        } catch (e) {
        }
        Object.keys(_Module).forEach(Module => {
            try {
                const __Module = window.CodeInterface[Language][Module]
                const fc = __Module.func
                if (fc) {
                    if (__Module.isHTTP === isHttp) {
                        const obj = {
                            value: Module,
                            label: Module,
                            action: () => {
                                codeDebugToolsFunc(fc)
                            },
                        }
                        for (let i = 0; i < array2.length; i++) {
                            const element = array2[i];
                            if (element.value === Module) {
                                array2[i].children.push(obj)
                                return
                            }
                        }
                        array2.push({
                            value: Language,
                            label: Language,
                            children: [obj],
                        })
                    }
                }
            } catch (e) {
            }
        })
    });
    if (array2.length > 0) {
        for (let i = 0; i < array2.length; i++) {
            const element = array2[i];
            array.push(element)
        }
    }
    if (array3.length > 0) {
        for (let i = 0; i < array3.length; i++) {
            const element = array3[i];
            array.push(element)
        }
    }
    return array
}


export function GetGenerateCodeListMenu(Type, func, Theology) {
    const array = []
    const isUDP = Type.indexOf("udp") !== -1
    if (isUDP) {
        return array;
    }
    const isWebsocket = Type.indexOf("websocket") !== -1
    const isTCP = Type.indexOf("tcp") !== -1
    const isHttp = !isWebsocket && !isTCP

    /*
        isHTTP: true,//只有HTTP请求显示这个选项
        isWebsocket: false,
        isTCP: false
    * */
    Object.keys(GenerateCodeListMenu).forEach(key => {
        const obj = GenerateCodeListMenu[key];
        const arr = [];
        obj.forEach(value => {
            if ((isHttp && value.Type === GenerateCodeType_HTTP) || (isWebsocket && value.Type === GenerateCodeType_Websocket) || (isTCP && value.Type === GenerateCodeType_TCP)) {
                arr.push({
                    name: value.Name,
                    action: () => {
                        func(key, value.Name)
                    },
                })
            }
        })
        if (arr.length > 0) {
            array.push({
                name: key,
                subMenu: arr,
            })
        }
    });
    const array2 = []
    const array3 = []
    Object.keys(window.CodeInterface).forEach(Language => {
        const _Module = window.CodeInterface[Language]
        try {
            if (_Module.func) {
                for (let i = 0; i < array2.length; i++) {
                    const element = array2[i];
                    if (element) {
                        if (element.name === Language) {
                            return
                        }
                    }
                }
                if (_Module.isHTTP === isHttp && _Module.isWebsocket === isWebsocket && _Module.isTCP === isTCP) {
                    array3.push({
                        name: Language,
                        action: () => {
                            codeFunc({
                                isHTTP: isHttp,
                                isWebsocket: isWebsocket,
                                isTCP: isTCP,
                                Theology: Theology
                            }, _Module.func)
                        },
                    })
                }
                return
            }
        } catch (e) {
        }
        Object.keys(_Module).forEach(Module => {
            try {
                const __Module = window.CodeInterface[Language][Module]
                const fc = __Module.func
                if (fc) {
                    if (__Module.isHTTP === isHttp && __Module.isWebsocket === isWebsocket && __Module.isTCP === isTCP) {
                        const obj = {
                            name: Module,
                            action: () => {
                                codeFunc({
                                    isHTTP: isHttp,
                                    isWebsocket: isWebsocket,
                                    isTCP: isTCP,
                                    Theology: Theology
                                }, fc)
                            },
                        }
                        for (let i = 0; i < array2.length; i++) {
                            const element = array2[i];
                            if (element.name === Module) {
                                array2[i].subMenu.push(obj)
                                return
                            }
                        }
                        array2.push({
                            name: Language,
                            subMenu: [obj],
                        })
                    }
                }
            } catch (e) {
            }
        })
    });
    if (array2.length > 0 || array3.length > 0) {
        array.push("separator")
    }
    if (array2.length > 0) {
        for (let i = 0; i < array2.length; i++) {
            const element = array2[i];
            array.push(element)
        }
    }
    if (array3.length > 0) {
        for (let i = 0; i < array3.length; i++) {
            const element = array3[i];
            array.push(element)
        }
    }
    return array
}

function codeFunc(args, func) {
    AppGenerateCodeInterface(args.Theology).then((r) => {
        let res = ""
        if (args.isTCP) {
            res = func({
                isTCP: args.isTCP,
                isWebsocket: args.isWebsocket,
                isHTTP: args.isHTTP,
                Request: {
                    Host: r.URL,
                    IP: r.IP,
                    TmpFile: r.Path
                }
            })
        } else {
            res = func({
                isTCP: args.isTCP,
                isWebsocket: args.isWebsocket,
                isHTTP: args.isHTTP,
                Request: {
                    Method: r.Method,
                    URL: r.URL,
                    Header: r.Header,
                    Body: base64ToBytes(r.Body),
                    ServerIP: r.IP,
                    TmpFile: r.Path
                }
            })
        }
        if (res === null || res === undefined || res === "") {
            ElNotification({
                position: 'bottom-right',
                message: '生成代码失败\n请检查接口函数是否编写正确？',
                type: 'warning',
                customClass: 'multiline-message'
            })
            return
        }
        ClipboardWriteAll(res)
        ElNotification({
            position: 'bottom-right',
            showClose: true,
            message: '代码生成成功\n\n已复制到剪辑版',
            type: 'success',
            customClass: 'multiline-message'
        })
    })
}

function initGenerateCodeInterface(code) {
    try {
        window.CodeInterface = []
        const tmp = `function run() {
    let CodeInterface = {};
` + code + `
    init()
    return CodeInterface;
}
return run();
`
        const f = new Function(tmp)
        window.CodeInterface = f()
    } catch (e) {
        window.CodeInterface = []
    }

    if (window.GenerateCodeList) {
        window.GenerateCodeList()
    }
}

Events.On("GenerateCodeInterface", (obj) => {
    let o = obj.data[0];
    if (!o) {
        o = obj.data;
    }
    initGenerateCodeInterface(o)
})

AppGetGenerateCodeInterface().then(res => {
    initGenerateCodeInterface(res)
})