import {Do} from "../../wailsjs/go/main/App.js";

import {Base64} from 'js-base64';
import {ElMessage, ElMessageBox} from "element-plus";
import {BrowserOpenURL} from "../../wailsjs/runtime/runtime.js";

export function Base64Encode(str) {
    return Base64.fromUint8Array(Uint8Array.from(StringToBytes(str)));
}

export function StrBase64Encode(str) {
    const Encoder = new TextEncoder();
    return Base64.fromUint8Array(Uint8Array.from(Encoder.encode(str)));
}

export function Base64BytesEncode(bytes) {
    return Base64.fromUint8Array(Uint8Array.from(bytes));
}

export function PbJsonConvert(pbJson) {
    let Json = {}
    for (let index = 0; index < pbJson.length; index++) {
        if (pbJson[index].Type === "Object") {
            Json[pbJson[index].tag] = PbJsonConvert(pbJson[index].value)
        } else if (pbJson[index].Type === "Varint") {
            if (parseInt(pbJson[index].value).toString() !== pbJson[index].value) {
                Json[pbJson[index].tag] = "超出精度的大整数:" + pbJson[index].value
            } else {
                Json[pbJson[index].tag] = parseInt(pbJson[index].value)
            }
        } else if (pbJson[index].Type === "String") {
            Json[pbJson[index].tag] = pbJson[index].value
        } else if (pbJson[index].Type === "StringRaw") {
            const v = pbJson[index].note + ""
            if (v.indexOf("\\ufffd\\ufffd") !== -1) {
                Json[pbJson[index].tag] = "字节集数据 ->Base64:" + pbJson[index].value
            } else {
                Json[pbJson[index].tag] = pbJson[index].value
            }
        } else {
            Json[pbJson[index].tag] = pbJson[index].value
        }
    }
    return Json;
}

export function HexToGbkStr(str) {
    const encoder = new TextEncoder();
    const decoder = new TextDecoder("gbk");
    let a = [];
    for (let i = 0; i < str.length; i++) {
        const char = str[i];
        if (char === '%') {
            // 获取下2个字符
            const hexCode = str.substring(i + 1, i + 3);
            const decimalCode = parseInt(hexCode, 16);
            a.push(decimalCode);
            i += 2; // 跳过已读取的字符
        } else {
            const im = char.charCodeAt()
            if (im < 256) {
                a.push(im);
            } else {

            }
        }
    }
    return UInt8ToStr(new Uint8Array(a), "gbk");
}

//byte转字符串
export function bytesToString(bytes) {
    if (bytes === null) {
        return ""
    }
    let str = [];
    for (let i = 0; i < bytes.length; i++) {
        const b = bytes[i]
        const c = String.fromCharCode(b)
        str.push(c);
    }
    return str.join("");
}

export function Base64Decode(str) {
    return bytesToString(Base64.toUint8Array(str));
}

export function Base64DecodeStr(str) {
    const decoder = new TextDecoder('utf-8');
    return decoder.decode(Base64.toUint8Array(str));
}

export function Base64DecodeUint8(str) {
    if (str === null) {
        return new Uint8Array([])
    }
    return Base64.toUint8Array(str);
}

export function StringToBytes(str) {
    if (str === null) {
        return []
    }
    const bytes = [];
    for (let i = 0; i < str.length; i++)
        bytes.push(str.charCodeAt(i) & 0xFF);
    return bytes;
}

export function StringToUInt8Array(str) {
    if (str === null) {
        return new Uint8Array([])
    }
    for (let i = 0; i < str.length; i++) {
        let c = str.charCodeAt(i)
        if (c > 127) {
            console.log(c)
        }
    }
    const Encoder = new TextEncoder();
    return Encoder.encode(str);
}

export function stringToHex(str) {
    let hex = '';
    for (let i = 0; i < str.length; i++) {
        const charCode = str.charCodeAt(i);
        const hexCode = charCode.toString(16);
        hex += hexCode.padStart(2, '0'); // 补齐为两位十六进制数
    }
    return hex;
}

export function UInt8ToHex(str, space) {
    if (space) {
        let hex = '';
        for (let i = 0; i < str.length; i++) {
            const charCode = str[i];
            const hexCode = charCode.toString(16);
            hex += hexCode.padStart(2, '0') + " "; // 补齐为两位十六进制数
        }
        return hex.toUpperCase().trim();
    }
    let hex = '';
    for (let i = 0; i < str.length; i++) {
        const charCode = str[i];
        const hexCode = charCode.toString(16);
        hex += hexCode.padStart(2, '0'); // 补齐为两位十六进制数
    }
    return hex.toUpperCase();
}

export function UInt8ToStr(bs, LABEL) {
    const decoder = new TextDecoder(LABEL);
    return decoder.decode(bs);
}

export async function protobufToJson(pb, skip) {
    return CallGoDo("protobufToJson", {Data: Base64.fromUint8Array(pb), skip: skip});
}

export function SetUint8Array(arg1, arg2) {
    let a1 = null;
    if (arg1 instanceof Uint8Array) {
        a1 = arg1
    } else if (Array.isArray(arg1)) {
        a1 = new Uint8Array(arg1)
    } else if (typeof arg1 === "string") {
        const encoder = new TextEncoder();
        a1 = encoder.encode(arg1);
    } else {
        debugger
        return arg1
    }
    let a2 = null;
    if (arg2 instanceof Uint8Array) {
        a2 = arg2
    } else if (Array.isArray(arg2)) {
        a2 = new Uint8Array(arg2)
    } else if (typeof arg2 === "string") {
        const encoder = new TextEncoder();
        a2 = encoder.encode(arg2);
    } else {
        debugger
        return arg2
    }
    const newArray = new Uint8Array(a1.length + a2.length);
    newArray.set(a1);
    newArray.set(a2, a1.length);
    return newArray
}

export async function CallGoDo(arg1, arg2) {
    return await Do({Command: arg1, Args: arg2});
}

export function SunnyErrorReplaceAll(res) {
    let Body = Base64DecodeUint8(res)
    let _Body = UInt8ToStr(Body, "utf-8")
    if (_Body.indexOf("�") !== -1) {
        _Body = UInt8ToStr(Body, "gbk")
    }
    _Body = _Body.replaceAll("dial tcp", "连接到")
    _Body = _Body.replaceAll("connectex: No connection could be made because the target machine actively refused it.", "由于目标计算机主动拒绝，因此无法建立连接。")
    if (_Body.indexOf("read") !== -1 && _Body.indexOf("i/o timeout") !== -1) {
        _Body = "读取数据超时"
    }
    _Body = _Body.replaceAll("websocket: bad handshake", " Websocket 握手失败")
    _Body = _Body.replaceAll(":443: i/o timeout", " 超时")
    _Body = _Body.replaceAll(":443:i/o timeout", " 超时")
    _Body = _Body.replaceAll("Handshake failed. ", " 握手失败")
    _Body = _Body.replaceAll("A connection attempt failed because the connected party did not properly respond after a period of time, or established connection failed because connected host has failed to respond.", " 连接尝试失败，因为连接方在一段时间后没有正确响应，或者由于连接的主机没有响应而建立的连接失败。")
    _Body = _Body.replaceAll(": bind: Only one usage of each socket address (protocol/network address/port) is normally permitted", " 该端口已经被使用,请更换其他端口")
    _Body = _Body.replaceAll("listen tcp 0.0.0.0", "监听端口")
    //listen tcp 0.0.0.0:2022: bind: Only one usage of each socket address (protocol/network address/port) is normally permitted
    return _Body
}

export function deepCopy(obj) {
    if (obj === null || typeof obj !== "object") {
        return obj;
    }

    let copy;

    if (obj instanceof Array) {
        copy = [];
        for (let i = 0; i < obj.length; i++) {
            copy[i] = deepCopy(obj[i]);
        }
    } else if (obj instanceof Function) {
        copy = obj; // 直接复制函数引用
    } else {
        copy = {};
        for (let key in obj) {
            if (obj.hasOwnProperty(key)) {
                copy[key] = deepCopy(obj[key]);
            }
        }
    }

    return copy;
}

window.WindowOpenURL = function (a) {
    BrowserOpenURL(a)
}

export async function EventsDo(arg) {
    let Args = arg["Args"]
    switch (arg.Command) {
        case "弹出成功提示":
        case "弹出成功信息":
        case "弹出成功消息":
            ElMessage({
                message: Args,
                type: 'success',
            })
            return
        case "弹出提示消息":
        case "弹出提示信息":
            const a = Args.msg.replaceAll("\r", "").replaceAll("\n", "<br>")
            await ElMessageBox.alert(
                a,
                Args.title,
                {
                    dangerouslyUseHTMLString: true,
                    confirmButtonText: '好的',
                    closeOnClickModal: true, // 设置点击遮罩层关闭消息框
                    closeOnPressEscape: true, // 设置按下 ESC 键关闭消息框
                }
            )
            return
        case "有新版本":
            await ElMessageBox.alert(
                Args,
                "有可用更新",
                {
                    dangerouslyUseHTMLString: true,
                    confirmButtonText: '朕知道了',
                }
            )
            return
        case "弹出错误提示":
        case "弹出错误信息":
        case "弹出错误消息":
            ElMessage({
                message: Args,
                dangerouslyUseHTMLString: true,
                type: 'error',
            })
            return
        case "更新状态文本":
            window.vm.Footer.Title = (Args + "")
            return
        case "启动状态":
            console.log("启动状态", Args)
            if (Args === null || Args === "") {
                CallGoDo("获取运行端口", null).then(port => {
                    window.vm.Footer.Title = "启动成功:" + port
                })
            } else {
                window.vm.Footer.Title = "启动失败:  " + SunnyErrorReplaceAll(Args)
            }
            return
        case "插入列表":
            let array = []
            for (let i = 0; i < Args.length; i++) {
                let obj = window.vm.List.RowDataHashMap[Args[i].Theology]
                if (!obj) {
                    Args[i]["序号"] = (window.vm.List.index++) + 1
                    SetTextColor(Args[i])
                    array.push(Args[i])
                }
                //window.vm.List.RowData.push(Args[i])
            }
            if (array.length > 0) {
                const res = window.vm.List.agGridApi.applyTransaction({add: array});
                if (res.add) {
                    //let array = []
                    res.add.forEach(function (rowNode) {
                        window.vm.List.RowDataHashMap[rowNode.data.Theology] = rowNode
                    });
                    //console.log(array)
                    //refreshCells(array)
                    setTimeout(() => {
                        window.vm.List.agGridApi.applyTransaction({add: []});
                    }, 500)
                }
                IsRefreshList = true
                //window.vm.List.agGridApi.setRowData(window.vm.List.RowData);
            }
            return
        case "更新响应长度": {
            //let array = []
            for (let i = 0; i < Args.length; i++) {
                const Theology = Args[i]["Theology"]
                let obj = window.vm.List.RowDataHashMap[Theology]
                if (obj) {
                    obj.data["响应长度"] = "" + Args[i]["Send"] + "/" + Args[i]["Rec"]
                    //array.push(obj)
                    obj.setData(obj.data)
                }
            }
            IsRefreshList = true
            //refreshCells(array)
            //window.vm.List.agGridApi.applyTransaction({update: array});
            //window.vm.List.agGridApi.setRowData(window.vm.List.RowData);
        }
            return
        case "更新响应":
            if (window.Theology === Args["Theology"]) {
                if (Args["Error"]) {
                    window.vm.List.agSelectedLine.data["断点模式"] = 0
                    window.vm.Tabs.Request.BreakResponse = false
                    window.vm.Tabs.Request.Breakpoint = false
                    let Body = SunnyErrorReplaceAll(Args["Body"])
                    window.vm.Tabs.ToolPanel.NoData = true
                    window.vm.Tabs.ToolPanel.TitleMsg = `<span style="color: red">${Body}</span>`
                    return
                }
                window.vm.Tabs.ToolPanel.NoData = false
                window.vm.Tabs.ToolPanel.StateCode = Args.StateCode
                window.vm.Tabs.ToolPanel.CodeState = "success"
                if (Args.StateCode === 403 || Args.StateCode === 404) {
                    window.vm.Tabs.ToolPanel.CodeState = "danger"
                } else if (Args.StateCode !== 200) {
                    window.vm.Tabs.ToolPanel.CodeState = "warning"
                }

                if (Args["断点状态"]) {
                    window.vm.List.agSelectedLine.data["断点模式"] = 2
                    window.vm.Tabs.Request.BreakResponse = false
                    window.vm.Tabs.Request.Breakpoint = true
                    window.vm.Tabs.Request.GetBodyRect(window.vm.Tabs.Request.$refs.BodyRect.offsetHeight)
                    window.vm.Tabs.Response.UpdateData(Args, 2)
                } else {
                    window.vm.Tabs.Response.UpdateData(Args, 0)
                }

            }
            return
        case "更新Socket":
            let objs = []
            for (let i = 0; i < Args.length; i++) {
                if (window.Theology === Args[i]["Theology"]) {
                    objs.push(Args[i])
                }
            }
            window.vm.Tabs.Request.$refs.WebSocket.AddLines(objs)
            window.vm.Tabs.Request.$refs.WebSocket.Refresh()
            return
        case "更新列表": {
            //let array = []
            for (let i = 0; i < Args.length; i++) {
                const Theology = Args[i]["Theology"]
                let obj = window.vm.List.RowDataHashMap[Theology]
                if (obj) {
                    obj.data["响应类型"] = Args[i]["响应类型"]
                    obj.data["响应时间"] = Args[i]["响应时间"]
                    obj.data["响应长度"] = Args[i]["响应长度"]
                    obj.data["断点模式"] = Args[i]["断点模式"]
                    obj.data["请求地址"] = Args[i]["请求地址"]
                    obj.data["方式"] = Args[i]["方式"]
                    obj.data["状态"] = Args[i]["状态"]
                    obj.data["ico"] = Args[i]["ico"]
                    obj.data["WebSocket"] = Args[i]["WebSocket"]
                    if (obj.data["ico"] === "websocket_connect") {
                        obj.data["状态"] = "已连接"
                    } else if (obj.data["ico"] === "websocket_close") {
                        obj.data["状态"] = "已断开"
                    }
                    SetTextColor(obj.data)
                    //array.push(obj)
                    obj.setData(obj.data)
                }
            }
            IsRefreshList = true
            // refreshCells(array)
            //window.vm.List.agGridApi.applyTransaction({update: array});
            //window.vm.List.agGridApi.setRowData(window.vm.List.RowData);
        }
            return
        case "更新ICO":
            //let array = []
            for (let i = 0; i < Args.length; i++) {
                const Theology = Args[i]["Theology"]
                let obj = window.vm.List.RowDataHashMap[Theology]
                if (obj) {
                    obj.data["ico"] = Args[i]["ico"]
                    if (obj.data["ico"] === "websocket_connect") {
                        obj.data["状态"] = "已连接"
                    } else if (obj.data["ico"] === "websocket_close") {
                        obj.data["状态"] = "已断开"
                    }

                    obj.setData(obj.data)
                    //array.push(obj)
                }
            }
            IsRefreshList = true
            //refreshCells(array)
            //window.vm.List.agGridApi.applyTransaction({update: array});
            //window.vm.List.agGridApi.setRowData(window.vm.List.RowData);
            return
        case "加载配置":
            if (Args['ColorConfig']) {
                window.interface.colorConfig = Args['ColorConfig']
            }
            //常规设置
        {
            if (Args['Port']) {
                window.vm.Settings.$refs.Basic.Port = Args['Port']
            }
            if (Args['DisableUDP']) {
                window.vm.Settings.$refs.Basic.Option.DisableUDP = Args['DisableUDP']
            }
            if (Args['DisableCache']) {
                window.vm.Settings.$refs.Basic.Option.DisableBrowserCache = Args['DisableCache']
            }
            if (Args['OpenAuthentication']) {
                window.vm.Settings.$refs.Basic.Option.authentication = Args['OpenAuthentication']
            }
            if (Args['AuthenticationUserInfo']) {
                const obj = Args['AuthenticationUserInfo']
                window.vm.Settings.$refs.Basic.$refs.UserInfo.RowData = []
                Object.keys(obj).forEach(function (prop) {
                    if (obj && obj[prop]) {
                        const user = prop.trim()
                        const pass = obj[prop].trim()
                        if (user !== '' && pass !== '') {
                            window.vm.Settings.$refs.Basic.$refs.UserInfo.RowData.push({'账号': user, '密码': pass})
                        }
                    }
                });
                window.vm.Settings.$refs.Basic.$refs.UserInfo.agGridApi.setRowData(window.vm.Settings.$refs.Basic.$refs.UserInfo.RowData);
            }
        }
            //上游代理设置
        {
            if (Args['GlobalProxy']) {
                let agent = Args['GlobalProxy'] + ""
                let nAgent = agent.replace("Socket5://", "");
                nAgent = nAgent.replace("HTTP://", "");
                if (agent.startsWith("Socket5://")) {
                    window.vm.Settings.$refs.agent.select = "Socket5://"
                } else {
                    window.vm.Settings.$refs.agent.select = "HTTP://"
                }
                window.vm.Settings.$refs.agent.input1 = nAgent
            }
            if (Args['GlobalProxyRules']) {
                window.vm.Settings.$refs.agent.$refs.agentInput.SetCode(Args['GlobalProxyRules'])
            }
        }
            //强制TCP
        {
            if (Args['MustTcp']) {
                const C = Args['MustTcp']['open']
                const value = Args['MustTcp']['Rules']
                window.vm.Settings.$refs.MustList.MustTcpMode = C ? "MustTcp" : "CancelMustTcp"
                window.vm.Settings.$refs.MustList.$refs.Rules.SetCode(value)
            }

        }
            //证书选择
        {
            if (Args['Cert']) {
                const Default = Args['Cert']['Default']
                const CaPath = Args['Cert']['CaPath']
                const KeyPath = Args['Cert']['KeyPath']
                window.vm.Settings.$refs.ssl.sslMode = Default ? "默认证书" : "自定义证书"
                window.vm.Settings.$refs.ssl.CaFilePath = CaPath
                window.vm.Settings.$refs.ssl.KeyFilePath = KeyPath
            }
        }
            //替换规则
        {
            const objs = Args['ReplaceRules']
            if (objs) {
                window.vm.Settings.$refs.Replace.AddLines(objs)
            }
        }
            //Hosts规则
        {
            const objs = Args['HostsRules']
            if (objs) {
                window.vm.Settings.$refs.hosts.AddLines(objs)
            }
        }
            //请求证书管理
        {

            if (Args['RequestCertManager']) {
                const objs = Args['RequestCertManager']
                const array = []
                Object.keys(objs).forEach(function (prop) {
                    if (objs && objs[prop]) {
                        let obj = {}
                        obj.context = parseInt(prop)
                        obj.rule = objs[prop].rule
                        obj.Host = objs[prop].Host
                        obj.FilePath = objs[prop].FilePath
                        obj.PassWord = objs[prop].PassWord
                        array.push(obj)
                    }
                });
                window.vm.Settings.$refs.RequestCertificate.AddLines(array);
            }
        }
            //主题
        {
            window.Theme.IsDark = (Args['DarkTheme'] + "") === "1"
            window.vm.Theme.Theme = (Args['DarkTheme'] + "") === "1"
        }
            //其他
        {
            window.Theme.GOOS = Args.GOOS
            window.vm.Doc.activeName = Args.GOOS === "windows" ? "Windows" : "MacOs"
            if (Args.GOOS !== "windows") {
                window.vm.List.columns[6].hide = true
            }
        }
            //恢复过滤器
        {
            try {
                const Filter = JSON.parse(Args.Filter)
                Object.keys(Filter).forEach((key) => {
                    const responseTypeFilter = window.vm.List.agGridApi.getFilterInstance(key);
                    responseTypeFilter.setModel(Filter[key]);
                });
                window.vm.List.agGridApi.onFilterChanged();
            } catch (e) {
            }
        }
            //恢复快捷键
        {
            try {
                if (Args.GOOS !== "windows") {
                    Object.keys(window.KeysStrings).forEach((key) => {
                        if (window.KeysStrings[key]) {
                            if (window.KeysStrings[key].value) {
                                let obj = window.KeysStrings[key].value
                                obj = obj.replaceAll("CTRL", "Control")
                                window.KeysStrings[key].value = obj.replaceAll("ALT", "Option")
                            }
                        }
                    });
                }
                const KeysStrings = JSON.parse(Args.KeysStrings)
                Object.keys(KeysStrings).forEach((key) => {
                    window.KeysStrings[key].ctrlKey = KeysStrings[key].ctrlKey
                    window.KeysStrings[key].altKey = KeysStrings[key].altKey
                    window.KeysStrings[key].shiftKey = KeysStrings[key].shiftKey
                    window.KeysStrings[key].key = KeysStrings[key].key
                    if (Args.GOOS !== "windows") {
                        let obj = KeysStrings[key].value
                        obj = obj.replaceAll("CTRL", "Control")
                        window.KeysStrings[key].value = obj.replaceAll("ALT", "Option")
                    } else {
                        let obj = KeysStrings[key].value
                        obj = obj.replaceAll("Control", "CTRL")
                        window.KeysStrings[key].value = obj.replaceAll("Option", "ALT")
                    }
                });
                window.vm.List.agGridApi.onFilterChanged();
            } catch (e) {
            }
        }
            //恢复列数据
        {
            try {
                const Columns = JSON.parse(Args.Columns)
                Object.keys(Columns).forEach((key) => {
                    Columns[key].editable = false
                    if (Columns[key]) {
                        if (Columns[key].field === '注释') {
                            Columns[key].editable = true
                        }
                    }
                    window.vm.List.columns[key] = Columns[key]
                });

            } catch (e) {
            }

        }
            //111111111111111111111
            return
        case "更新搜索进度":
            window.vm.Find.per = Args
            return
        default:
            console.log(arg.Command, arg)
            return
    }
}

export function SetColorConfig(Name, Color) {
    const obj = window.interface.colorConfig[Name]
    if (obj) {
        obj.dark = Color.dark
        obj.right = Color.right
        IsRefreshRenderedNodes = true
        CallGoDo("保存配置", {Type: "列表颜色配置", Data: window.interface.colorConfig}).then(r => {
        })
    }
}

let IsRefreshList = false
let IsRefreshRenderedNodes = false
window.interface = {
    colorConfig: {
        tcp: {dark: '', right: ""},
        udp: {dark: '', right: ""},
        ws: {dark: '', right: ""},
        css: {dark: '', right: ""},
        js: {dark: '', right: ""},
        img: {dark: '', right: ""},
        document: {dark: '', right: ""},
        _1: {dark: '', right: ""},
        _301: {dark: '', right: ""},
        _302: {dark: '', right: ""},
        _401: {dark: '', right: ""},
        _403: {dark: '', right: ""},
        _404: {dark: '', right: ""},
        _500: {dark: '', right: ""},
    }
}
let isExecuting = false;
setInterval(() => {
    if (!isExecuting) {
        {
            if (window.vm.List && IsRefreshList) {
                window.vm.List.agGridApi.applyTransaction({add: []});
                IsRefreshList = false
            }
        }
        {
            if (IsRefreshRenderedNodes) {
                window.vm.List.RefreshRenderedNodes()
                IsRefreshRenderedNodes = false
            }
        }
        isExecuting = false;
    }
}, 200)

export function SetTextColor(obj) {
    if (obj) {
        if (obj.color === null || obj.color === void 0) {
            obj.color = {}
        }
        const Method = obj.方式.toUpperCase()
        const Type = obj.响应类型.toUpperCase()
        const State = obj.状态.toUpperCase()
        if (Method.indexOf("TCP") !== -1) {
            obj.color.text = window.interface.colorConfig.tcp
        } else if (Method.indexOf("UDP") !== -1) {
            obj.color.text = window.interface.colorConfig.udp
        } else if (Method.indexOf("WEBSOCKET") !== -1) {
            obj.color.text = window.interface.colorConfig.ws
        } else if (Method.indexOf("POST") !== -1) {

        }
        if (Type.indexOf("CSS") !== -1) {
            obj.color.text = window.interface.colorConfig.css
        } else if (Type.indexOf("JAVASCRIPT") !== -1) {
            obj.color.text = window.interface.colorConfig.js
        } else if (Type.indexOf("IMAGE") !== -1) {
            obj.color.text = window.interface.colorConfig.img
        } else if (Type.indexOf("TEXT/") !== -1) {
            obj.color.text = window.interface.colorConfig.document
        }


        if (State === "-1") {
            obj.color.text = window.interface.colorConfig._1
        } else if (State === "301") {
            obj.color.text = window.interface.colorConfig._301
        } else if (State === "302") {
            obj.color.text = window.interface.colorConfig._302
        } else if (State === "401") {
            obj.color.text = window.interface.colorConfig._401
        } else if (State === "403") {
            obj.color.text = window.interface.colorConfig._403
        } else if (State === "404") {
            obj.color.text = window.interface.colorConfig._404
        } else if (State === "500") {
            obj.color.text = window.interface.colorConfig._500
        }


    }
}

function refreshCells(array) {
    window.vm.List.agGridApi.refreshCells({rowNodes: array, force: true})
}

function formatTime(timestamp) {
    const date = new Date(timestamp * 1000); // 将时间戳转换为毫秒级别

    const hours = date.getHours();
    const minutes = "0" + date.getMinutes();
    const seconds = "0" + date.getSeconds();
    return hours + ':' + minutes.substr(-2) + ':' + seconds.substr(-2);
}