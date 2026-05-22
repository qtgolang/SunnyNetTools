import {bytesToBase64, bytesToString, StringToBytes, toGBK, toUTF8} from "./encoding.js";
import {NewBytesIO} from "./BytesIO.js";


export function toHeader(Header) {
    let message = '';
    Object.entries(Header).forEach(([key, value]) => {
        try {
            value.forEach((val) => {
                message += key + ": " + val + "\r\n"
            })
        } catch (e) {
            message += key + ": " + value.join("; ") + "\r\n"
        }
    });
    return message
}

export function toHeaderArray(Header) {
    const message = [];
    try {
        Object.entries(Header).forEach(([key, value]) => {
            try {
                value.forEach((val) => {
                    message.push({name: key, value: val});
                })
            } catch (e) {
                message.push({name: key, value: value.join("; ")});
            }
        });
    } catch (e) {
    }
    return message
}

export function headerArrayToLanguage(Array) {
    for (let i = 0; i < Array.length; i++) {
        if (Array[i].name.toUpperCase() === "CONTENT-TYPE") {
            const ar = (Array[i].value + "/").replaceAll(";", "/").split("/")
            if (ar.length >= 2) {
                if (ar[0].toLowerCase() === "image") {
                    return (ar[0] + "/" + ar[1]).toLowerCase()
                }
                return ar[1].toLowerCase()
            }
        }
    }
    return "plaintext"
}

export function toResponseCookiesHeader(Header) {
    const _Header = toHeaderArray(Header)
    const message = [];
    _Header.forEach((obj) => {
        if ((obj.name + "").toLowerCase() === "set-cookie") {
            const array = (obj.value + "").replaceAll(";", "=").split("=")
            const v = {name: "", value: "", value2: ""};
            if (array.length > 0) {
                v.name = array[0]
            }
            if (array.length > 1) {
                v.value = array[1]
            }
            if (v.name !== '' && v.value !== '') {
                v.value2 = (obj.value + "").replaceAll(v.name + "=" + v.value + ";", "").trim()
                message.push(v)
            }
        }
    })
    return message;
}


export function toRequestCookiesHeader(Headers) {
    const cookieValues = [];
    for (let i = 0; i < Headers.length; i++) {
        const header = Headers[i];
        if (header.name.toLowerCase().startsWith('cookie')) {
            const cookies = header.value.split(';');
            for (let j = 0; j < cookies.length; j++) {
                const cookie = cookies[j].trim();
                const nameIndex = cookie.indexOf('=');
                if (nameIndex !== -1) {
                    const name = cookie.substring(0, nameIndex);
                    const val = cookie.substring(nameIndex + 1);
                    cookieValues.push({name: name, value: val});
                }
            }
            return cookieValues;
        }
    }
    return cookieValues;
}

export function urlToArgs(u) {
    let message = [];
    let array = (u + "").split("?");
    if (array.length === 2) {
        array = array[1].split("&");
        for (let i = 0; i < array.length; i++) {
            const array1 = array[i].split("=");
            const val = {name: "", value: ""}
            if (array1.length > 1) {
                val.name = array1[0]
                val.value = array1.slice(1).map(str => str.trim()).join("=");
                message.push(val)
                continue
            }
            if (array1.length > 0) {
                val.name = array1[0]
                message.push(val)
            }
        }
    }
    return message
}

//数据是否表单格式
export function bodyIsForm(_Body) {
    let formParams = '';
    try {
        formParams = toGBK(_Body);
    } catch (e) {
        try {
            formParams = toUTF8(_Body);
        } catch (e) {
            return [];
        }
    }

    // 过滤掉不符合表单参数格式的情况
    if (formParams.indexOf('=') === -1 || formParams.match(/["{}[\]':]/) !== null) {
        return [];
    }

    return urlToArgs(formParams);
}

//替换错误信息
export function ErrorReplace(Error) {
    let err = Error + ""
    err = err.replaceAll("[SunnyNet]", "")
    if (err.indexOf("use of closed network connection") !== -1) {
        return "接收数据失败";
    }
    return err;
}

const spaceCode = ' '.charCodeAt(0);
const COLON_CHAR_CODE = ':'.charCodeAt(0);

export function parseHTTPHeaders(obj) {
    const headers = {};
    while (obj.Buffered() > 0) {
        const LineBytes = obj.ReadLine(true)
        if (LineBytes.length === 0) {
            break
        }
        const Line = NewBytesIO(LineBytes);
        const key = bytesToString(Line.ReadBytes(COLON_CHAR_CODE, true));
        const value = bytesToString(Line.ReadAll()).trim();
        if (!headers[key]) {
            headers[key] = [];
        }
        headers[key].push(value)
    }
    return headers;
}

export function bytesToRequest(data) {
    const obj = NewBytesIO(StringToBytes(data));
    const res = {Method: "", URL: "", Proto: "", Header: {}, data: ""};
    const FirstLineBytes = obj.ReadLine(true)
    if (FirstLineBytes) {
        const FirstLine = NewBytesIO(FirstLineBytes);
        res.Method = bytesToString(FirstLine.ReadBytes(spaceCode, true));
        res.URL = bytesToString(FirstLine.ReadBytes(spaceCode, true));
        res.Proto = bytesToString(FirstLine.ReadAll());
    }
    res.Header = parseHTTPHeaders(obj)
    res.data = bytesToBase64(obj.ReadAll())
    return res;
}

export function bytesToResponse(data) {
    const obj = NewBytesIO(StringToBytes(data));
    const res = {Proto: "", Code: "", StateText: "", Header: {}, data: ""};
    const FirstLineBytes = obj.ReadLine(true)
    if (FirstLineBytes) {
        const FirstLine = NewBytesIO(FirstLineBytes);
        res.Proto = bytesToString(FirstLine.ReadBytes(spaceCode, true));
        res.Code = bytesToString(FirstLine.ReadBytes(spaceCode, true));
        res.StateText = bytesToString(FirstLine.ReadAll());
    }
    res.Header = parseHTTPHeaders(obj)
    res.data = bytesToBase64(obj.ReadAll())
    return res;
}