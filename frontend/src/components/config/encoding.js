import {ref, watch} from "vue";

export function base64ToBytes(base64) {
    if (base64 === null) {
        return new Uint8Array(0);
    }
    if (base64 === '') {
        return new Uint8Array(0);
    }
    // 解码 Base64 得到二进制字符串
    const binaryString = atob(base64);
    const len = binaryString.length;
    const bytes = new Uint8Array(len);
    for (let i = 0; i < len; i++) {
        bytes[i] = binaryString.charCodeAt(i) & 0xFF;
    }
    return bytes;
}

export function StringToBytes(bs) {
    if (bs === null) {
        return "";
    }
    return StringOrBytesJoinToBytes(bs);
}

export function StringOrBytesJoinToBytes(...arrays) {
    const convertedArrays = arrays.map(array => {
        if (typeof array === 'string') {
            return new TextEncoder().encode(array); // 将字符串转换为 Uint8Array
        } else if (array instanceof Uint8Array) {
            return array; // 直接返回 Uint8Array
        } else {
            throw new TypeError('参数必须是 Uint8Array 或字符串');
        }
    });
    const totalLength = convertedArrays.reduce((sum, array) => sum + array.length, 0);
    const mergedArray = new Uint8Array(totalLength);
    let offset = 0;
    convertedArrays.forEach(array => {
        mergedArray.set(array, offset);
        offset += array.length;
    });
    return mergedArray;
}
window.bytesToString = bytesToString
window.bytesToBase64 = bytesToBase64 
export function bytesToBase64(uint8Array) {
    if (uint8Array === null || uint8Array.length === 0) {
        return ""
    }
    // 创建一个字符串
    let binaryString = '';
    const len = uint8Array.length;
    for (let i = 0; i < len; i++) {
        binaryString += String.fromCharCode(uint8Array[i]); // 将每个字节转换为字符
    }
    return btoa(binaryString); // 使用 btoa() 转换为 Base64
}

export function sanitizeHTML(htmlString) {
    // 1. 解析 HTML
    let parser = new DOMParser();
    let doc = parser.parseFromString(htmlString, "text/html");

    // 2. 删除 <script> 标签
    doc.querySelectorAll("script").forEach(script => script.remove());

    // 3. 移除 <a> 的 href，防止跳转
    doc.querySelectorAll("a").forEach(a => a.removeAttribute("href"));

    // 4. 移除 <img> 的 src，防止加载远程图片
    doc.querySelectorAll("img").forEach(img => img.removeAttribute("src"));

    // 5. 移除 <iframe>，防止嵌入远程页面
    doc.querySelectorAll("iframe").forEach(iframe => iframe.remove());

    // 6. 移除 <link> 的 href，防止加载外部 CSS
    doc.querySelectorAll("link[rel='stylesheet']").forEach(link => link.removeAttribute("href"));

    // 7. 移除 <video>、<audio> 的 src，防止播放外部资源
    doc.querySelectorAll("video, audio").forEach(media => media.removeAttribute("src"));

    // 8. 移除 <form> 的 action，防止表单提交
    doc.querySelectorAll("form").forEach(form => form.removeAttribute("action"));

    // 9. 移除 <meta http-equiv="refresh">，防止自动跳转
    doc.querySelectorAll("meta[http-equiv='refresh']").forEach(meta => meta.remove());

    // 10. 移除所有事件属性 (onload, onerror, onclick等)
    const eventAttributes = [
        "onload", "onerror", "onclick", "onmouseover", "onfocus",
        "onsubmit", "onblur", "onchange", "onkeyup", "onkeydown"
    ];
    doc.querySelectorAll("*").forEach(el => {
        eventAttributes.forEach(attr => el.removeAttribute(attr));
    });

    // 11. 返回处理后的 HTML
    return doc.documentElement.outerHTML;
}

export function bytesToString(uint8Array) {
    if (uint8Array === null || uint8Array.length === 0) {
        return ""
    }
    try {
        let _Body = toUTF8(uint8Array)
        if (_Body.indexOf("�") !== -1) {
            _Body = toGBK(uint8Array)
        }
        return _Body;
    } catch (e) {
        try {
            let _Body = toGBK(uint8Array)
            if (_Body.indexOf("�") !== -1) {
                _Body = toUTF8(uint8Array)
            }
            return _Body;
        } catch (e) {
        }
    }
    return toUTF8(uint8Array);
}

export function toGBK(uint8Array) {
    if (uint8Array === null || uint8Array.length === 0) {
        return ""
    }
    const decoder = new TextDecoder('gbk');
    if (typeof uint8Array === "string") {
        const encoder = new TextEncoder();
        const Array = encoder.encode(uint8Array);
        return decoder.decode(Array);
    }
    return decoder.decode(uint8Array);
}

export function toUTF8(uint8Array) {
    if (uint8Array === null || uint8Array.length === 0) {
        return ""
    }
    const decoder = new TextDecoder('utf-8');
    if (typeof uint8Array === "string") {
        const encoder = new TextEncoder();
        const Array = encoder.encode(uint8Array);
        return decoder.decode(Array);
    }
    return decoder.decode(uint8Array);
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

window.ToGBK = toGBK
window.ToUTF8 = toUTF8
window.Base64ToBytes = base64ToBytes
window.BytesJoin = StringOrBytesJoinToBytes
//当前编码 UTF8 / GBK
export const Config_Encoding_Current_ISUTF8 = ref(true);

const MenuContextKeys = [];
watch(Config_Encoding_Current_ISUTF8, (newValue, oldValue) => {
    MenuContextKeys.forEach((obj) => {
        obj.set(newValue)
    })
});


export function addMonacoEditorEncodingMenu(editor) {
    const Utf8_encoding = editor.createContextKey('encoding', Config_Encoding_Current_ISUTF8.value)
    editor.addAction({
        id: 'turnencoding',
        label: '显示GBK编码',
        contextMenuGroupId: 'my-commands1',
        contextMenuOrder: Number.MAX_SAFE_INTEGER,
        precondition: 'encoding',
        run: () => {
            Config_Encoding_Current_ISUTF8.value = false;
            Utf8_encoding.set(false);
        },
    })
    editor.addAction({
        id: 'turnencoding-utf',
        label: '显示UTF-8编码',
        contextMenuGroupId: 'my-commands1',
        contextMenuOrder: Number.MAX_SAFE_INTEGER,
        precondition: '!encoding',
        run: () => {
            Config_Encoding_Current_ISUTF8.value = true;
            Utf8_encoding.set(true);
        },
    })
    MenuContextKeys.push(Utf8_encoding)


    editor.updateOptions({
        renderValidationDecorations: "off", // 关闭错误和警告提示
        unicodeHighlight: {
            ambiguousCharacters: false, // 关闭不明确的 Unicode 字符警告
            invisibleCharacters: false, // 关闭不可见字符警告
            nonBasicASCII: false,       // 关闭非 ASCII 字符警告
        }
    });
}
