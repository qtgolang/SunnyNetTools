package Session

import (
	"changeme/Service/clipboard"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

var splitBytes = []byte("*SunnyNetV4|.0.0.|SunnyNetV4*")

func ExportMessage(app AppSession) (string, error) {
	tempDir := os.TempDir()
	filePath := path.Join(tempDir, "SunnyNet_"+strconv.Itoa(app.GetTheology())+".bin")
	_ = os.Remove(filePath)
	f, e := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	if e != nil {
		return "", e
	}
	defer f.Close()
	isWebsocket := app.IsWebsocket()
	tmpBytes := make([]byte, 1)
	app.RangeStream(func(obj AppStream) bool {
		if obj.GetIsSend() {
			t := fmt.Sprintf("%d", parseFormattedTimeToTimestamp(obj.GetMessageTime()))
			if len(t) < 13 {
				_, _ = f.WriteString("0000000000000")
			} else {
				_, _ = f.WriteString(t[0:13])
			}
			if isWebsocket {
				tmpBytes[0] = byte(obj.GetWebsocketType())
				_, _ = f.Write(tmpBytes[:])
			}
			_, _ = f.Write(obj.GetBody())
			_, _ = f.Write(splitBytes)
		}
		return true
	})
	return filePath, nil
}
func parseFormattedTimeToTimestamp(formatted string) int64 {
	// 先按 ':' 拆出毫秒
	parts := strings.Split(formatted, ":")
	if len(parts) != 2 {
		return 0
	}

	hmsPart := parts[0] // "HH-MM-SS"
	msPart := parts[1]  // "mmm"

	// 替换 - 为 : 变成标准格式
	hmsFixed := strings.ReplaceAll(hmsPart, "-", ":")

	// 获取今天的年月日
	dateStr := time.Now().Format("2006-01-02")

	// 拼成完整时间字符串
	fullTimeStr := fmt.Sprintf("%s %s.%s", dateStr, hmsFixed, msPart)

	// 解析
	t, err := time.Parse("2006-01-02 15:04:05.000", fullTimeStr)
	if err != nil {
		return 0
	}

	return t.UnixNano() / int64(time.Millisecond)
}

func createECode(app AppSession, module string) error {
	/*
	       {Name: "网页_访问", Type: GenerateCodeType_HTTP},                1
	       {Name: "网页_访问_对象", Type: GenerateCodeType_HTTP},           1
	       {Name: "E2EE网站客户端", Type: GenerateCodeType_HTTP},           1
	       {Name: "WinHttpW", Type: GenerateCodeType_HTTP},           1
	       {Name: "WinInet", Type: GenerateCodeType_HTTP},           1
	       {Name: "WinHttpR", Type: GenerateCodeType_HTTP},           1
	       {Name: "SunnyHTTP", Type: GenerateCodeType_HTTP},
	       {Name: "SunnyWS客户端-异步", Type: GenerateCodeType_Websocket},
	       {Name: "SunnyWS客户端-同步", Type: GenerateCodeType_Websocket},
	       {Name: "SunnyTCP客户端-异步", Type: GenerateCodeType_TCP},
	       {Name: "SunnyTCP客户端-同步", Type: GenerateCodeType_TCP},
	       {Name: "Sunny-TLS-TCP客户端-异步", Type: GenerateCodeType_TCP},
	       {Name: "Sunny-TLS-客户端-同步", Type: GenerateCodeType_TCP},
	   ]
	*/
	code := ""
	switch module {
	case "SunnyWS客户端-同步", "SunnyWS客户端-异步":
		{
			if !app.IsWebsocket() {
				return errors.New("这个请求不支持生成 易语言 的 " + module + " 代码")
			}
			hp, ok := app.(*HttpSession)
			if !ok {
				return errors.New("这个请求不是HTTP请求")
			}
			URL, e := url.Parse(hp.Request.Url)
			if e != nil {
				return errors.New("解析URL失败")
			}
			tempFile, er := ExportMessage(hp)
			if er != nil {
				return errors.New("导出消息记录失败")
			}

			header := hp.Request.Header.Clone()
			Cookie := header.Get("Cookie")
			header.Del("Cookie")
			a1 := ""
			a2 := ".循环判断首 ()\n    发送数据 ＝ 客户端.接收数据 (2000, 返回数据类型)\n    调试输出 (发送数据, 返回数据类型)\n.循环判断尾 (发送数据 ≠ {  })"
			a3 := "真"
			a4 := "同步"
			if module == "SunnyWS客户端-异步" {
				a1 = "到整数 (&___回调函数参考)"
				a2 = `
延迟 (50000)  ' 确保数据已经全部被接收


.子程序 ___回调函数参考, , 公开
.参数 _Context, 整数型
.参数 消息类型, 整数型, , 1=接收消息 2=接收时连接被断开 3=发送时连接被断开
.参数 数据指针, 整数型, , 消息类型=2、3时 这里是错误信息
.参数 指针长度, 整数型
.参数 数据类型, 整数型, , #Sunny_WsMessage_ (当消息类型=1时有效)

.如果真 (消息类型 ≠ 1)
    调试输出 ("ws 连接被断开")
    返回 ()
.如果真结束
调试输出 (指针到字节集 (数据指针, 指针长度), "数据类型：", 数据类型)
`
				a3 = "假"
				a4 = "异步"
			}

			_header := ""
			for k, v := range header {
				if strings.ToUpper(k) == "CONTENT-LENGTH" {
					continue
				}
				if len(v) < 1 {
					_header += "局_协议头 ＝ 局_协议头 ＋ \"" + k + ": \"＋ #换行符\n"
				} else {
					_header += "局_协议头 ＝ 局_协议头 ＋ " + convertELangFormat(k+": "+v[0]) + "＋ #换行符\n"
				}
			}
			if Cookie != "" {
				_header += "局_协议头 ＝ 局_协议头 ＋ " + convertELangFormat("Cookie: "+Cookie) + "＋ #换行符\n"
			}
			code = `
.子程序 功能_WebSocket_` + a4 + `, , , 本子程序由Sunny中间件生成,请配合 [Sunny中间件模块] 使用。
.局部变量 客户端, SunnyWSS客户端
.局部变量 索引, 整数型
.局部变量 局_网址, 文本型
.局部变量 局_协议头, 文本型
.局部变量 发送数据数组, 字节集, , "0"
.局部变量 发送数据, 字节集
.局部变量 类型, 整数型
.局部变量 返回数据类型, 整数型
.局部变量 本次记录的时间, 长整数型
.局部变量 上次记录的时间, 长整数型

局_网址 ＝ "` + URL.String() + `"

` + _header + `


.如果真 (客户端.连接 (局_网址, 局_协议头, ` + a1 + `, ` + a3 + `) ＝ 假)
    调试输出 ("WS连接失败：", 局_网址, "错误信息", 客户端.取错误信息 ())
    返回 ()
.如果真结束

发送数据数组 ＝ 分割字节集 (读入文件 ("` + tempFile + `"), 到字节集 ("*SunnyNetV4|.0.0.|SunnyNetV4*"), )

.计次循环首 (取数组成员数 (发送数据数组), 索引)
    本次记录的时间 ＝ 到长整数 (到文本 (取字节集左边 (发送数据数组 [索引], 13)))
    .如果真 (上次记录的时间 ＞ 0)
        .如果 (本次记录的时间 － 上次记录的时间 ＞ 5000)
            Sleep (5000)' 防止太长时间的等待
        .否则
            Sleep (本次记录的时间 － 上次记录的时间)
        .如果结束

    .如果真结束
    上次记录的时间 ＝ 本次记录的时间
    发送数据 ＝ 取字节集右边 (发送数据数组 [索引], 取字节集长度 (发送数据数组 [索引]) － 13)
    .如果真 (取字节集长度 (发送数据) ≥ 2)
        类型 ＝ 发送数据 [1]
        发送数据 ＝ 取字节集右边 (发送数据, 取字节集长度 (发送数据) － 1)
        客户端.发送数据 (发送数据, 类型)
    .如果真结束

.计次循环尾 ()
` + a2 + ` 
`
		}
		break
	case "SunnyTCP客户端-异步", "SunnyTCP客户端-同步":
		{
			if !app.IsTCP() {
				return errors.New("这个请求不支持生成 易语言 的 " + module + " 代码")
			}
			hp, ok := app.(*TCPSession)
			if !ok {
				return errors.New("这个请求不是TCP请求")
			}
			tempFile, er := ExportMessage(hp)
			if er != nil {
				return errors.New("导出消息记录失败")
			}

			a1 := ""
			a2 := "\n.循环判断首 ()\n    接收数据 ＝ 客户端.取回数据 (2000, 错误信息)\n    调试输出 (接收数据)\n.循环判断尾 (接收数据 ≠ {  })"
			a3 := "真"
			a4 := "同步"
			a5 := ".局部变量 接收数据, 字节集\n"
			if module == "SunnyTCP客户端-异步" {
				a1 = "到整数 (&___回调函数参考)"
				a2 = `
延迟 (5000)  ' 确保数据被接收

.子程序 ___回调函数参考, , 公开
.参数 客户端标识, 整数型
.参数 消息类型, 整数型, , 1=接收消息 2=接收时连接被断开 3=发送时连接被断开
.参数 数据指针, 整数型, , 消息类型=2、3时 这里是错误信息
.参数 指针长度, 整数型

.如果真 (消息类型 ≠ 1)
    调试输出 ("TCP连接被断开", "错误信息", 到文本 (指针到字节集 (数据指针, 指针长度)))
    返回 ()
.如果真结束
调试输出 (指针到字节集 (数据指针, 指针长度))
`
				a3 = "假"
				a4 = "异步"
				a5 = ""
			}
			code = `
.子程序 功能_TCP_` + a4 + `, , , 本子程序由Sunny中间件生成,请配合 [Sunny中间件模块] 使用。
.局部变量 客户端, SunnyTCP客户端
.局部变量 索引, 整数型
.局部变量 局_地址, 文本型
.局部变量 错误信息, 文本型
.局部变量 发送数据数组, 字节集, , "0"
.局部变量 发送数据, 字节集
` + a5 + `.局部变量 本次记录的时间, 长整数型
.局部变量 上次记录的时间, 长整数型

局_地址 ＝ "` + hp.Host + `"
.如果真 (客户端.连接 (局_地址, ` + a1 + `, ` + a3 + `, 假) ＝ 假)
    调试输出 ("TCP连接失败：", 局_地址, "错误信息", 客户端.取错误信息 ())
    返回 ()
.如果真结束

发送数据数组 ＝ 分割字节集 (读入文件 ("` + tempFile + `"), 到字节集 ("*SunnyNetV4|.0.0.|SunnyNetV4*"), )

.计次循环首 (取数组成员数 (发送数据数组), 索引)
    本次记录的时间 ＝ 到长整数 (到文本 (取字节集左边 (发送数据数组 [索引], 13)))
    .如果真 (上次记录的时间 ＞ 0)
        .如果 (本次记录的时间 － 上次记录的时间 ＞ 5000)
            Sleep (5000)' 防止太长时间的等待
        .否则
            Sleep (本次记录的时间 － 上次记录的时间)
        .如果结束

    .如果真结束
    上次记录的时间 ＝ 本次记录的时间
    发送数据 ＝ 取字节集右边 (发送数据数组 [索引], 取字节集长度 (发送数据数组 [索引]) － 13)
    客户端.发送数据 (发送数据, 3000)

.计次循环尾 ()
` + a2 + ` 
`
		}
		break
	case "SunnyHTTP":
		{
			if !app.IsHTTP() {
				return errors.New("这个请求不支持生成 易语言 的 " + module + " 代码")
			}
			hp, ok := app.(*HttpSession)
			if !ok {
				return errors.New("这个请求不是HTTP请求")
			}
			URL, e := url.Parse(hp.Request.Url)
			if e != nil {
				return errors.New("解析URL失败")
			}
			header := hp.Request.Header.Clone()
			FuncName := "SunnyNetCreateRequest"
			Cookie := header.Get("Cookie")
			header.Del("Cookie")
			BytesType := IsBytesType(hp.Request.Body)
			code += ".子程序 " + FuncName + ", , 公开, 本子程序由Sunny中间件生成,请配合 [SunnyNet中间件模块] 使用。 " + URL.Path + "\n.局部变量 局_HTTP, " + module + "\n"
			if BytesType {
				code += ".局部变量 局_提交字节集, 字节集, , , \n"
			} else {
				code += ".局部变量 局_提交数据, 文本型, , , \n"
			}
			if Cookie != "" {
				code += ".局部变量 局_提交Cookie, 文本型, , , \n"
			}
			code += ".局部变量 局_响应字节集, 字节集, , , \n.局部变量 局_响应文本, 文本型, , , \n\n"
			if BytesType {
				code += "局_提交字节集 ＝ 编码_BASE64解码 (“" + base64.StdEncoding.EncodeToString(hp.Request.Body) + "”)\n\n"
			} else {
				ok1, d := isFormData(hp)
				if ok1 {
					code += d + "\n"
				} else {
					code += "局_提交数据 ＝ " + convertELangFormat(string(hp.Request.Body)) + "\n"
				}
			}
			if Cookie != "" {
				code += "局_提交Cookie ＝ 子文本替换 (" + convertELangFormat(Cookie) + ", “'”, #引号, , , 真)\n\n"
			}
			code += "局_HTTP.打开 (“" + hp.Request.Method + "”, “" + URL.String() + "”)\n"
			for k, v := range header {
				if strings.ToUpper(k) == "CONTENT-LENGTH" {
					continue
				}
				if k == "Accept-Encoding" {
					if len(v) < 1 {
						code += "' 局_HTTP.置协议头 (“" + k + "”, “”)\n"
					} else {
						code += "' 局_HTTP.置协议头 (“" + k + "”, " + convertELangFormat(v[0]) + ")\n"
					}
					continue
				}
				if len(v) < 1 {
					code += "局_HTTP.置协议头 (“" + k + "”, “”)\n"
				} else {
					code += "局_HTTP.置协议头 (“" + k + "”, " + convertELangFormat(v[0]) + ")\n"
				}
			}
			if Cookie != "" {
				code += "局_HTTP.置协议头 (“Cookie”, 局_提交Cookie)\n\n"
			}
			code += "局_HTTP.使用随机TLS指纹 (真)  ' 仅随机Ja3指纹\n"
			if hp.Response.ServerIP != "" && hp.Response.ServerIP != "本地响应" {
				code += "局_HTTP.设置实际连接地址 (“" + hp.Response.ServerIP + "”)  ' 直接使用这个地址,而不是使用DNS解析后的地址\n\n"
			}
			if BytesType {
				code += "局_HTTP.发送字节集 (局_提交字节集)\n"
			} else {
				code += "局_HTTP.发送 (局_提交数据)\n"
			}
			code += "局_响应字节集 ＝ 局_HTTP.取响应内容 ()\n"
			code += "局_响应文本 ＝ 编码_Utf8到Ansi (局_响应字节集)\n"
			code += "调试输出 (局_响应文本)\n\n"
		}
		break
	case "E2EE网站客户端":
		{
			if !app.IsHTTP() {
				return errors.New("这个请求不支持生成 易语言 的 " + module + " 代码")
			}
			hp, ok := app.(*HttpSession)
			if !ok {
				return errors.New("这个请求不是HTTP请求")
			}
			URL, e := url.Parse(hp.Request.Url)
			if e != nil {
				return errors.New("解析URL失败")
			}
			header := hp.Request.Header.Clone()
			FuncName := "SunnyNetCreateRequest"
			Cookie := header.Get("Cookie")
			header.Del("Cookie")
			BytesType := IsBytesType(hp.Request.Body)

			code += ".子程序 " + FuncName + ", , 公开, 本子程序由Sunny中间件生成,请配合 [E2EE] 使用。 " + URL.Path + "\n.局部变量 局_HTTP, 网站客户端\n"
			code += ".局部变量 局_网址, 文本型, , , \n"
			if BytesType {
				code += ".局部变量 局_提交数据, 字节集, , , \n"
			} else {
				code += ".局部变量 局_提交数据, 文本型, , , \n"
			}
			if Cookie != "" {
				code += ".局部变量 局_提交Cookie, 文本型, , , \n"
			}
			code += ".局部变量 局_响应字节集, 字节集, , , \n.局部变量 局_响应文本, 文本型, , , \n\n"

			code += "局_网址 ＝ " + convertELangFormat(URL.String()) + "\n\n"
			if BytesType {
				code += "局_提交数据 ＝ 编码_BASE64解码 (“" + base64.StdEncoding.EncodeToString(hp.Request.Body) + "”)\n\n"
			} else {
				ok1, d := isFormData(hp)
				if ok1 {
					code += d + "\n"
				} else {
					code += "局_提交数据 ＝ " + convertELangFormat(string(hp.Request.Body)) + "\n"
				}
			}
			if Cookie != "" {
				code += "局_提交Cookie ＝ 子文本替换 (" + convertELangFormat(Cookie) + ", “'”, #引号, , , 真)\n\n"
			}
			for k, v := range header {
				if strings.ToUpper(k) == "CONTENT-LENGTH" {
					continue
				}
				if k == "Accept-Encoding" {
					if len(v) < 1 {
						code += "' 局_HTTP.置请求头 (“" + k + "”, “”)\n"
					} else {
						code += "' 局_HTTP.置请求头 (“" + k + "”, " + convertELangFormat(v[0]) + ")\n"
					}
					continue
				}
				if len(v) < 1 {
					code += "局_HTTP.置请求头 (“" + k + "”, “”)\n"
				} else {
					code += "局_HTTP.置请求头 (“" + k + "”, " + convertELangFormat(v[0]) + ")\n"
				}
			}
			if Cookie != "" {
				code += "局_HTTP.置请求头 (“Cookie”, 局_提交Cookie)\n\n"
			}
			if hp.Request.Method == "GET" {
				code += "局_HTTP.执行GET (局_网址 , 局_响应字节集, 真, )\n"
			} else {
				code += "局_HTTP.执行POST (局_网址,局_提交数据, 局_响应字节集, 真, )\n"
			}
			code += "局_响应文本 ＝ 到文本 (局_响应字节集)\n"
			code += "调试输出 (局_响应文本)\n\n"
		}
		break
	case "网页_访问", "网页_访问_对象":
		{
			if !app.IsHTTP() {
				return errors.New("这个请求不支持生成 易语言 的 " + module + " 代码")
			}
			hp, ok := app.(*HttpSession)
			if !ok {
				return errors.New("这个请求不是HTTP请求")
			}
			URL, e := url.Parse(hp.Request.Url)
			if e != nil {
				return errors.New("解析URL失败")
			}
			header := hp.Request.Header.Clone()
			FuncName := "SunnyNetCreateRequest"
			Cookie := header.Get("Cookie")
			header.Del("Cookie")
			BytesType := IsBytesType(hp.Request.Body)

			code += ".子程序 " + FuncName + ", , 公开, 本子程序由Sunny中间件生成,请配合 [精易模块] 使用。 " + URL.Path + "\n"

			code += ".局部变量 局_网址, 文本型, , , \n"
			if BytesType {
				code += ".局部变量 局_提交数据, 字节集, , , \n"
			} else {
				code += ".局部变量 局_提交数据, 文本型, , , \n"
			}
			code += ".局部变量 局_协议头, 类_POST数据类, , , \n"

			if Cookie != "" {
				code += ".局部变量 局_提交Cookie, 文本型, , , \n"
			}
			code += ".局部变量 局_响应字节集, 字节集, , , \n.局部变量 局_响应文本, 文本型, , , \n\n"

			code += "局_网址 ＝ " + convertELangFormat(URL.String()) + "\n\n"

			if BytesType {
				code += "局_提交数据 ＝ 编码_BASE64解码 (“" + base64.StdEncoding.EncodeToString(hp.Request.Body) + "”)\n\n"
			} else {
				ok1, d := isFormData(hp)
				if ok1 {
					code += d + "\n"
				} else {
					code += "局_提交数据 ＝ " + convertELangFormat(string(hp.Request.Body)) + "\n"
				}
			}
			if Cookie != "" {
				code += "局_提交Cookie ＝ 子文本替换 (" + convertELangFormat(Cookie) + ", “'”, #引号, , , 真)\n\n"
			}

			for k, v := range header {
				if strings.ToUpper(k) == "CONTENT-LENGTH" {
					continue
				}
				if k == "Accept-Encoding" {
					if len(v) < 1 {
						code += "' 局_协议头.添加 (“" + k + "”, “”)\n"
					} else {
						code += "' 局_协议头.添加 (“" + k + "”, " + convertELangFormat(v[0]) + ")\n"
					}
					continue
				}
				if len(v) < 1 {
					code += "局_协议头.添加 (“" + k + "”, “”)\n"
				} else {
					code += "局_协议头.添加 (“" + k + "”, " + convertELangFormat(v[0]) + ")\n"
				}
			}
			if Cookie != "" {
				code += "局_协议头.添加 (“Cookie”, 局_提交Cookie)\n\n"
			}
			if hp.Request.Method == "GET" {
				code += "局_响应字节集 ＝ " + module + " (局_网址, 0, , , , 局_协议头.获取协议头数据 ())\n"
			} else {
				mod := "1"
				if hp.Request.Method == "POST" {
					mod = "1"
				} else if hp.Request.Method == "HEAD" {
					mod = "2"
				} else if hp.Request.Method == "PUT" {
					mod = "3"
				} else if hp.Request.Method == "OPTIONS" {
					mod = "4"
				} else if hp.Request.Method == "DELETE" {
					mod = "5"
				} else if hp.Request.Method == "TRACE" {
					mod = "6"
				} else if hp.Request.Method == "CONNECT" {
					mod = "7"
				}
				if BytesType {
					if module == "网页_访问_对象" {
						code += "局_响应字节集 ＝ " + module + " (局_网址, " + mod + ", , , , 局_协议头.获取协议头数据 (),,,,局_提交数据)\n"
					} else {
						code += "局_响应字节集 ＝ " + module + " (局_网址, " + mod + ", , , , 局_协议头.获取协议头数据 (),,,局_提交数据)\n"
					}

				} else {
					code += "局_响应字节集 ＝ " + module + " (局_网址, " + mod + ", 局_提交数据, , , 局_协议头.获取协议头数据 ())\n"
				}
			}
			code += "局_响应文本 ＝ 到文本 (局_响应字节集)\n"
			code += "调试输出 (局_响应文本)\n\n"
		}
		break
	case "WinInet", "WinHttpR", "WinHttpW":
		{
			if !app.IsHTTP() {
				return errors.New("这个请求不支持生成 易语言 的 " + module + " 代码")
			}
			hp, ok := app.(*HttpSession)
			if !ok {
				return errors.New("这个请求不是HTTP请求")
			}
			URL, e := url.Parse(hp.Request.Url)
			if e != nil {
				return errors.New("解析URL失败")
			}
			header := hp.Request.Header.Clone()
			FuncName := "SunnyNetCreateRequest"
			Cookie := header.Get("Cookie")
			header.Del("Cookie")
			BytesType := IsBytesType(hp.Request.Body)
			code += ".子程序 " + FuncName + ", , 公开, 本子程序由Sunny中间件生成,请配合 [WinHttp模块] 使用。 " + URL.Path + "\n.局部变量 局_HTTP, " + module + "\n"
			if BytesType {
				code += ".局部变量 局_提交字节集, 字节集, , , \n"
			} else {
				code += ".局部变量 局_提交数据, 文本型, , , \n"
			}
			if Cookie != "" {
				code += ".局部变量 局_提交Cookie, 文本型, , , \n"
			}
			code += ".局部变量 局_响应字节集, 字节集, , , \n.局部变量 局_响应文本, 文本型, , , \n\n"
			if BytesType {
				code += "局_提交字节集 ＝ 编码_BASE64解码 (“" + base64.StdEncoding.EncodeToString(hp.Request.Body) + "”)\n\n"
			} else {
				ok1, d := isFormData(hp)
				if ok1 {
					code += d + "\n"
				} else {
					code += "局_提交数据 ＝ " + convertELangFormat(string(hp.Request.Body)) + "\n"
				}
			}
			if Cookie != "" {
				code += "局_提交Cookie ＝ 子文本替换 (" + convertELangFormat(Cookie) + ", “'”, #引号, , , 真)\n\n"
			}
			code += "局_HTTP.Open (“" + hp.Request.Method + "”, “" + URL.String() + "”)\n"
			for k, v := range header {
				if strings.ToUpper(k) == "CONTENT-LENGTH" {
					continue
				}
				if k == "Accept-Encoding" {
					if len(v) < 1 {
						code += "' 局_HTTP.SetHeader (“" + k + "”, “”)\n"
					} else {
						code += "' 局_HTTP.SetHeader (“" + k + "”, " + convertELangFormat(v[0]) + ")\n"
					}
					continue
				}
				if len(v) < 1 {
					code += "局_HTTP.SetHeader (“" + k + "”, “”)\n"
				} else {
					code += "局_HTTP.SetHeader (“" + k + "”, " + convertELangFormat(v[0]) + ")\n"
				}
			}
			if Cookie != "" {
				code += "局_HTTP.SetHeader (“Cookie”, 局_提交Cookie)\n\n"
			}
			if BytesType {
				code += "局_HTTP.SendBin (局_提交字节集)\n"
			} else {
				code += "局_HTTP.Send (局_提交数据)\n"
			}
			code += "局_响应字节集 ＝ 局_HTTP.GetBody ()\n"
			code += "局_响应文本 ＝ 编码_Utf8到Ansi (局_响应字节集)\n"
			code += "调试输出 (局_响应文本)\n\n"
		}
		break
	default:
		return errors.New(fmt.Sprintf("无法生成[%s]的代码", module))
	}
	code = ".版本 2\n.支持库 spec\n\n" + code
	code = strings.ReplaceAll(code, "\r", "")
	code = strings.ReplaceAll(code, "\n", "\r\n")
	return clipboard.ClipboardWriteAll(code)
}

func createCShaPCode(app AppSession, module string) error {
	/*
		[
		        {Name: "HttpClient", Type: GenerateCodeType_HTTP},
		        {Name: "RestSharp", Type: GenerateCodeType_HTTP},
		]
	*/
	if !app.IsHTTP() {
		return errors.New("这个请求不支持生成 C# 的 " + module + " 代码")
	}
	hp, ok := app.(*HttpSession)
	if !ok {
		return errors.New("这个请求不是HTTP请求")
	}
	URL, e := url.Parse(hp.Request.Url)
	if e != nil {
		return errors.New("解析URL失败")
	}
	header := hp.Request.Header.Clone()
	FuncName := "SunnyNetCreateRequest"
	Cookie := header.Get("Cookie")
	header.Del("Cookie")
	BytesType := IsBytesType(hp.Request.Body)
	code := ""
	switch module {
	case "RestSharp":
		templateData := ""
		templateData1 := ""
		if len(hp.Request.Body) > 0 {
			if BytesType {
				templateData = `string base64String = "` + base64.StdEncoding.EncodeToString(hp.Request.Body) + `";  
            byte[] body = Convert.FromBase64String(base64String); `
			} else {
				templateData = `string String = "` + strReplaceAll(hp.Request.Body) + `";  
            byte[] body = Encoding.Default.GetBytes(String); `
			}
		}

		mod := "Method.Post"
		s := strings.ToUpper(hp.Request.Method)
		{
			if s == "POST" {
				mod = "Method.Post"
			}
			if s == "GET" {
				mod = "Method.Get"
			}
			if s == "PUT" {
				mod = "Method.Put"
			}
			if s == "DELETE" {
				mod = "Method.Delete"
			}
			if s == "HEAD" {
				mod = "Method.Head"
			}
			if s == "OPTIONS" {
				mod = "Method.Options"
			}
			if s == "PATCH" {
				mod = "Method.Patch"
			}
			if s == "MERGE" {
				mod = "Method.Merge"
			}
			if s == "COPY" {
				mod = "Method.Copy"
			}
			if s == "SEARCH" {
				mod = "Method.Search"
			}
		}
		_header := ""
		CONTENTMENT := "application/x-www-form-urlencoded"
		for k, v := range header {
			if strings.ToUpper(k) == "CONTENT-LENGTH" {
				continue
			}
			if strings.ToUpper(k) == "ACCEPT-ENCODING" {
				if len(v) > 0 {
					_header += `            //request.AddHeader("` + k + `","` + strReplaceAll([]byte(v[0])) + `")` + ";\n"
				}
				continue
			}
			if strings.ToUpper(k) == "CONTENT-TYPE" {
				if len(v) > 0 {
					ss := strings.Split(v[0]+";", ";")
					if len(ss) > 0 {
						CONTENTMENT = ss[0]
					}
				}

			}
			if len(v) < 1 {
				//request.AddHeader("cookie", "xxxxx");
				_header += `            request.AddHeader("` + k + `","")` + ";\n"
			} else {
				_header += `            request.AddHeader("` + k + `","` + strReplaceAll([]byte(v[0])) + `")` + ";\n"
			}
		}
		if Cookie != "" {
			_header += `            request.AddHeader("Cookie","` + strReplaceAll([]byte(Cookie)) + `")` + ";\n"
		}
		if len(hp.Request.Body) > 0 {
			templateData1 = `            request.AddParameter("` + CONTENTMENT + `", body, ParameterType.RequestBody); 
            `
		}
		_tmp := `/// <summary> 
        /// ` + FuncName + ` ` + URL.Path + `
        ///<para>本函数由SunnyNet网络中间件生成</para> 
        /// </summary> 
        public static void ` + FuncName + `()
        {
            string url = "` + URL.String() + `";
            ` + templateData + ` 
            var client = new RestClient(url);
            var request = new RestRequest("",` + mod + `);
` + _header + templateData1 + `var response = client.Execute(request); 
            Trace.WriteLine("Response StateCode:" + ((int)response.StatusCode)); 
            Trace.WriteLine("Response Text:\n" + response.Content); 
        }
`
		code = _tmp
		break
	case "HttpClient":

		templateData := `byte[] data = Encoding.Default.GetBytes(""); `
		if len(hp.Request.Body) > 0 {
			if BytesType {
				templateData = `string base64String = "` + base64.StdEncoding.EncodeToString(hp.Request.Body) + `";  
            byte[] data = Convert.FromBase64String(base64String); `
			} else {
				templateData = `string String = "` + strReplaceAll(hp.Request.Body) + `";  
            byte[] data = Encoding.Default.GetBytes(String); `
			}
		}
		mod := "Method.Post"
		CONTENTMENT := "application/x-www-form-urlencoded"
		s := strings.ToUpper(hp.Request.Method)
		{
			if s == "POST" {
				mod = "HttpMethod.Post"
			}
			if s == "GET" {
				mod = "HttpMethod.Get"
			}
			if s == "PUT" {
				mod = "HttpMethod.Put"
			}
			if s == "DELETE" {
				mod = "HttpMethod.Delete"
			}
			if s == "HEAD" {
				mod = "HttpMethod.Head"
			}
			if s == "OPTIONS" {
				mod = "HttpMethod.Options"
			}
			if s == "PATCH" {
				mod = "HttpMethod.Patch"
			}
			if s == "MERGE" {
				mod = "HttpMethod.Merge"
			}
			if s == "COPY" {
				mod = "HttpMethod.Copy"
			}
			if s == "SEARCH" {
				mod = "HttpMethod.Search"
			}
		}
		_header := ""
		for k, v := range header {
			if strings.ToUpper(k) == "CONTENT-LENGTH" {
				continue
			}
			if strings.ToUpper(k) == "CONTENT-TYPE" {
				if len(v) > 0 {
					ss := strings.Split(v[0]+";", ";")
					if len(ss) > 0 {
						CONTENTMENT = ss[0]
					}
				}
				continue
			}
			if strings.ToUpper(k) == "ACCEPT-ENCODING" {
				if len(v) > 0 {
					_header += `            //client.DefaultRequestHeaders.Add("` + k + `","` + strReplaceAll([]byte(v[0])) + `")` + ";\n"
				}
				continue
			}
			if len(v) < 1 {
				_header += `            client.DefaultRequestHeaders.Add("` + k + `","")` + ";\n"
			} else {
				_header += `            client.DefaultRequestHeaders.Add("` + k + `","` + strReplaceAll([]byte(v[0])) + `")` + ";\n"
			}
		}
		if Cookie != "" {
			_header += `            client.DefaultRequestHeaders.Add("Cookie","` + strReplaceAll([]byte(Cookie)) + `")` + ";\n"
		}
		_tmp := `/// <summary> 
        /// ` + FuncName + ` ` + URL.Path + `
        ///<para>本函数由SunnyNet网络中间件生成</para> 
        /// </summary> 
        public static async void ` + FuncName + `()
        {
            string url = "` + URL.String() + `";
            ` + templateData + `
            using (HttpClient client = new HttpClient())
            {
                ` + _header + `
                HttpContent content = new ByteArrayContent(data);
                content.Headers.ContentType = new System.Net.Http.Headers.MediaTypeHeaderValue("` + CONTENTMENT + `");
                HttpRequestMessage request = new HttpRequestMessage(` + mod + `, url);
                request.Content = content;
                HttpResponseMessage response = client.SendAsync(request).Result;
 
                Trace.WriteLine("Response Status code: " + response.StatusCode);
                byte[] responseBytes = response.Content.ReadAsByteArrayAsync().Result;
                Trace.WriteLine("Response bytes: " + BitConverter.ToString(responseBytes));
                Trace.WriteLine("Response Content: " + response.Content.ReadAsStringAsync().Result);
                 
            }
        }
`
		code = _tmp
		break
	default:
		return errors.New(fmt.Sprintf("无法生成[%s]的代码", module))
	}
	code = strings.ReplaceAll(code, "\r", "")
	code = strings.ReplaceAll(code, "\n", "\r\n")
	return clipboard.ClipboardWriteAll(code)
}

func createGoCode(app AppSession, module string) error {
	/*
			    {Name: "net/http", Type: GenerateCodeType_HTTP},
		        {Name: "SunnyNet/http", Type: GenerateCodeType_HTTP},
		        {Name: "Websocket-请求", Type: GenerateCodeType_Websocket},
		        {Name: "tcp-请求", Type: GenerateCodeType_TCP},
	*/
	FuncName := "SunnyNetCreateRequest"
	switch module {
	case "Websocket-请求":
		hp, isHttp := app.(*HttpSession)
		if !isHttp || !app.IsWebsocket() {
			return errors.New("这个请求不支持生成 GoLang 的 " + module + " 代码")
		}
		URL, e := url.Parse(hp.Request.Url)
		if e != nil {
			return errors.New("解析URL失败")
		}
		header := hp.Request.Header.Clone()
		Cookie := header.Get("Cookie")
		header.Del("Cookie")
		_header := "\t_headers := make(http.Header)\n"
		for k, v := range header {
			if strings.ToUpper(k) == "CONTENT-LENGTH" {
				continue
			}
			if strings.HasPrefix(strings.ToUpper(k), "SEC-WEBSOCKET") || strings.ToUpper(k) == "CONNECTION" || strings.ToUpper(k) == "UPGRADE" {
				if len(v) < 1 {
					_header += `	// _headers.Set("` + k + `","")` + "\n"
				} else {
					_header += `	// _headers.Set("` + k + `","` + strReplaceAll([]byte(v[0])) + `")` + "\n"
				}
				continue
			}
			if len(v) < 1 {
				_header += `	_headers.Set("` + k + `","")` + "\n"
			} else {
				_header += `	_headers.Set("` + k + `","` + strReplaceAll([]byte(v[0])) + `")` + "\n"
			}
		}
		if Cookie != "" {
			_header += `	_headers.Set("Cookie","` + strReplaceAll([]byte(Cookie)) + `")` + "\n"
		}
		if URL.Scheme == "http" {
			URL.Scheme = "ws"
		} else {
			URL.Scheme = "wss"
		}
		tempFile, er := ExportMessage(app)
		if er != nil {
			return errors.New("导出消息记录失败")
		}
		template := `
package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
) 

func main() {
	SunnyNetCreateRequest()
}

// ` + FuncName + ` 本函数由SunnyNet网络中间件生成  ` + URL.Path + `
func ` + FuncName + `() {
` + _header + `
	c, _, err := websocket.DefaultDialer.Dial("` + reText(URL.String()) + `", _headers) 
	if err != nil {
		log.Fatal("连接失败:", err)
	}
	defer c.Close()  
	bs, _ := os.ReadFile(` + "`" + tempFile + "`" + `)
	bytesArray := bytes.Split(bs, []byte("*SunnyNetV4|.0.0.|SunnyNetV4*"))
	Time := make([]byte, 13)
	StartTime := int64(0)
	for _, v := range bytesArray {
		if len(v) < 15 {
			continue
		}
		//计算与上一条数据间隔时间
		{
			copy(Time, v)
			ts, _ := strconv.ParseInt(string(Time), 10, 64)
			if StartTime != 0 {
				m := ts - StartTime
				if m > 5000 {
					//避免太长时间的模拟等待
					time.Sleep(5000 * time.Millisecond)
				} else {
					time.Sleep(time.Duration(m) * time.Millisecond)
				}
			}
			StartTime = ts
		}
		MessageType := v[13] //消息类型
		data := v[14:]       //真实发生数据
		_ = c.WriteMessage(int(MessageType), data)
	}
	for {
		//设置5秒读取超时
		_ = c.SetReadDeadline(time.Now().Add(time.Second * 5))
		// 接收消息
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Fatal("读取失败:", err)
		} 
		log.Printf("收到消息: %s", string(msg))
	}

`
		template += "\n}\n"
		code := strings.ReplaceAll(template, "\r", "")
		code = strings.ReplaceAll(code, "\n", "\r\n")
		return clipboard.ClipboardWriteAll(code)
	case "SunnyNet/http":
		hp, isHttp := app.(*HttpSession)
		if !isHttp {
			return errors.New("这个请求不支持生成 GoLang 的 " + module + " 代码")
		}
		URL, e := url.Parse(hp.Request.Url)
		if e != nil {
			return errors.New("解析URL失败")
		}
		header := hp.Request.Header.Clone()
		Cookie := header.Get("Cookie")
		header.Del("Cookie")
		BytesType := IsBytesType(hp.Request.Body)
		templateData := ""
		if BytesType {
			templateData = `	_Base64 := "` + base64.StdEncoding.EncodeToString(hp.Request.Body) + `"
	_Data, _ := base64.StdEncoding.DecodeString(_Base64)
	Body := io.NopCloser(bytes.NewBuffer(_Data))`
		} else {
			templateData = `	Body := io.NopCloser(bytes.NewBuffer([]byte("` + strReplaceAll(hp.Request.Body) + `")))`
		}
		_header := ""
		for k, v := range header {
			if strings.ToUpper(k) == "CONTENT-LENGTH" {
				continue
			}
			if _header == "" {
				if len(v) < 1 {
					_header += `req.Header.Set("` + k + `","")` + "\n"
				} else {
					_header += `req.Header.Set("` + k + `","` + strReplaceAll([]byte(v[0])) + `")` + "\n"
				}
				continue
			}
			if len(v) < 1 {
				_header += `	req.Header.Set("` + k + `","")` + "\n"
			} else {
				_header += `	req.Header.Set("` + k + `","` + strReplaceAll([]byte(v[0])) + `")` + "\n"
			}
		}
		if Cookie != "" {
			_header += `	req.Header.Set("Cookie","` + strReplaceAll([]byte(Cookie)) + `")` + "\n"
		}

		template := `
package main

import (
	"bytes"
	"fmt"
	"github.com/qtgolang/SunnyNet/src/http"
	"github.com/qtgolang/SunnyNet/src/httpClient"
	"io"
	"time"
)

func main() {
	SunnyNetCreateRequest()
}

// ` + FuncName + ` 本函数由SunnyNet网络中间件生成  ` + URL.Path + `
func ` + FuncName + `() {
` + templateData + `
	defer func() { _ = Body.Close() }()
	req, err := http.NewRequest("` + hp.Request.Method + `", "` + URL.String() + `", Body)
	if err != nil {
		panic(err)
	}
	` + _header + `
	res, _, err, Close := httpClient.Do(req, nil, false, nil, 30*time.Second, nil, nil)
	defer func() {
		if Close != nil {
			Close()
		}
	}()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = res.Body.Close()
	}()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("响应状态码", res.StatusCode)
	for k, v := range res.Header {
		fmt.Println(k, v)
	}
	fmt.Println(string(body))
}
`
		code := strings.ReplaceAll(template, "\r", "")
		code = strings.ReplaceAll(code, "\n", "\r\n")
		return clipboard.ClipboardWriteAll(code)
	case "net/http":
		hp, isHttp := app.(*HttpSession)
		if !isHttp {
			return errors.New("这个请求不支持生成 GoLang 的 " + module + " 代码")
		}
		URL, e := url.Parse(hp.Request.Url)
		if e != nil {
			return errors.New("解析URL失败")
		}
		header := hp.Request.Header.Clone()
		Cookie := header.Get("Cookie")
		header.Del("Cookie")
		BytesType := IsBytesType(hp.Request.Body)
		templateData := ""
		if BytesType {
			templateData = `	_Base64 := "` + base64.StdEncoding.EncodeToString(hp.Request.Body) + `"
	_Data, _ := base64.StdEncoding.DecodeString(_Base64)
	Body := io.NopCloser(bytes.NewBuffer(_Data))`
		} else {
			templateData = `	Body := io.NopCloser(bytes.NewBuffer([]byte("` + strReplaceAll(hp.Request.Body) + `")))`
		}
		_header := ""
		for k, v := range header {
			if strings.ToUpper(k) == "CONTENT-LENGTH" {
				continue
			}
			if _header == "" {
				if len(v) < 1 {
					_header += `req.Header.Set("` + k + `","")` + "\n"
				} else {
					_header += `req.Header.Set("` + k + `","` + strReplaceAll([]byte(v[0])) + `")` + "\n"
				}
				continue
			}
			if len(v) < 1 {
				_header += `	req.Header.Set("` + k + `","")` + "\n"
			} else {
				_header += `	req.Header.Set("` + k + `","` + strReplaceAll([]byte(v[0])) + `")` + "\n"
			}
		}
		if Cookie != "" {
			_header += `	req.Header.Set("Cookie","` + strReplaceAll([]byte(Cookie)) + `")` + "\n"
		}

		template := `
package main

func main() {
	SunnyNetCreateRequest()
}

// ` + FuncName + ` 本函数由SunnyNet网络中间件生成  //` + URL.Path + `
func ` + FuncName + `() {
` + templateData + `
	defer func() { _ = Body.Close() }()
	req, err := http.NewRequest("` + hp.Request.Method + `", "` + URL.String() + `", Body)
	if err != nil {
		panic(err)
	}
	` + _header + `
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer func() { _ = res.Body.Close() }()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("响应状态码", res.StatusCode)
	for k, v := range res.Header {
		fmt.Println(k, v)
	}
	fmt.Println(string(body))
}
`
		code := strings.ReplaceAll(template, "\r", "")
		code = strings.ReplaceAll(code, "\n", "\r\n")
		return clipboard.ClipboardWriteAll(code)
	case "tcp-请求":
		tp, isTcp := app.(*TCPSession)
		if !isTcp {
			return errors.New("这个请求不支持生成 GoLang 的 " + module + " 代码")
		}
		tempFile, er := ExportMessage(app)
		if er != nil {
			return errors.New("导出消息记录失败")
		}
		s := ""
		Name := "conn"
		if strings.Contains(strings.ToUpper(tp.Method), "TLS") {
			s = `
	config := &tls.Config{InsecureSkipVerify: true}
	tlsConn := tls.Client(conn, config)
	if er := tlsConn.Handshake(); er != nil {
		fmt.Println("握手失败:", er)
		return
	}
`
			Name = "tlsConn"
		}
		template := `
package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"time"
) 

func main() {
	SunnyNetCreateRequest()
}

// ` + FuncName + ` 本函数由SunnyNet网络中间件生成  -> ` + tp.RemoteAddress + `
func ` + FuncName + `() {
    conn, err := net.Dial("tcp", "` + tp.Host + `")
    if err != nil {
        fmt.Println("连接失败:", err)
        return
    }
    defer conn.Close()
` + s + `
	_ = ` + Name + `.SetDeadline(time.Time{})
	bs, _ := os.ReadFile(` + "`" + tempFile + "`" + `)
	bytesArray := bytes.Split(bs, []byte("*SunnyNetV4|.0.0.|SunnyNetV4*"))
	Time := make([]byte, 13)
	StartTime := int64(0)
	for _, v := range bytesArray {
		if len(v) < 14 {
			continue
		}
		//计算与上一条数据间隔时间
		{
			copy(Time, v)
			ts, _ := strconv.ParseInt(string(Time), 10, 64)
			if StartTime != 0 {
				m := ts - StartTime
				if m > 5000 {
					//避免太长时间的模拟等待
					time.Sleep(5000 * time.Millisecond)
				} else {
					time.Sleep(time.Duration(m) * time.Millisecond)
				}
			}
			StartTime = ts
		}
		data := v[13:] //真实发生数据
		_, _ = ` + Name + `.Write(data)
	}
	Reader := bufio.NewReader(` + Name + `)
	DataBuffer := make([]byte, 40960) //读取数据缓冲区
	for {
		//设置5秒读取超时
		_ = ` + Name + `.SetReadDeadline(time.Now().Add(time.Second * 5))
		// 接收消息
		Size, er := Reader.Read(DataBuffer)
		if er != nil {
			fmt.Println("读取失败:", er)
			return
		}
		message := DataBuffer[:Size]
		fmt.Println("收到消息:", string(message))
	}
`
		template += "\n}\n"
		code := strings.ReplaceAll(template, "\r", "")
		code = strings.ReplaceAll(code, "\n", "\r\n")
		return clipboard.ClipboardWriteAll(code)
	}
	return errors.New(fmt.Sprintf("无法生成[%s]的代码", module))
}

func createHSCode(app AppSession, module string) error {
	// {Name: "WinHttpW", Type: GenerateCodeType_HTTP},
	if !app.IsHTTP() {
		return errors.New("这个请求不支持生成 火山 的 " + module + " 代码")
	}
	hp, ok := app.(*HttpSession)
	if !ok {
		return errors.New("这个请求不是HTTP请求")
	}
	URL, e := url.Parse(hp.Request.Url)
	if e != nil {
		return errors.New("解析URL失败")
	}
	header := hp.Request.Header.Clone()
	FuncName := "SunnyNetCreateRequest"
	Cookie := header.Get("Cookie")
	header.Del("Cookie")
	BytesType := IsBytesType(hp.Request.Body)
	code := "<火山程序 类型 = \"通常\" 版本 = 1 />\r\n\r\n"
	if module == "WinHttpW" {
		code += "方法 " + FuncName + " <注释 = \"本函数由SunnyNet网络中间件生成,请搭配精易模块使用 " + URL.Path + "\">\n{\n"
		code += "    变量 局_HTTP <类型 = WinHttpW>\n"

		code += "    变量 局_请求地址 <类型 = 文本型>\n"
		if BytesType {
			code += "    变量 局_请求数据 <类型 = 字节集类>\n"
		} else {
			code += "    变量 局_请求数据 <类型 = 文本型>\n"
		}
		if Cookie != "" {
			code += "    变量 局_请求Cookie <类型 = 文本型>\n"
		}
		code += "    变量 局_响应字节集 <类型 = 字节集类>\n"
		code += "    变量 局_响应文本 <类型 = 文本型>\n"

		code += "    局_请求地址 ＝ \"" + reText(URL.String()) + "\"\n\n"

		if BytesType {
			code += "    局_请求数据 ＝ BASE64文本到字节集 (\"" + base64.StdEncoding.EncodeToString(hp.Request.Body) + "\")\n\n"
		} else {
			code += "    局_请求数据 ＝ \"" + reText(string(hp.Request.Body)) + "\"\n"
		}
		if Cookie != "" {
			code += "    局_请求Cookie ＝ \"" + reText(Cookie) + "\"\n"
		}

		code += "    局_HTTP.Open (\"" + hp.Request.Method + "\", 局_请求地址)\n"

		for k, v := range header {
			if strings.ToUpper(k) == "CONTENT-LENGTH" {
				continue
			}
			if k == "Accept-Encoding" {
				if len(v) < 1 {
					code += "    //局_HTTP.SetRequestHeader (\"" + k + "\",\"\")\n"
				} else {
					code += "    //局_HTTP.SetRequestHeader (\"" + k + "\",\"" + reText(v[0]) + "\")\n"
				}
				continue
			}
			if len(v) < 1 {
				code += "    局_HTTP.SetRequestHeader (\"" + k + "\",\"\")\n"
			} else {
				code += "    局_HTTP.SetRequestHeader (\"" + k + "\",\"" + reText(v[0]) + "\")\n"
			}
		}
		if Cookie != "" {
			code += "    局_HTTP.SetRequestHeader (\"Cookie\",局_请求Cookie)\n"
		}
		if hp.Request.Method == "GET" {
			code += "    局_HTTP.Send()\n"
		} else {
			if BytesType {
				code += "    局_HTTP.SendBin(局_请求数据)\n"
			} else {
				code += "    局_HTTP.Send(局_请求数据)\n"
			}
		}
		code += "    局_响应字节集 = 局_HTTP.GetResponseBody ()\n"
		code += "    局_响应文本 ＝ 多字节到文本 (局_响应字节集)\n"
		code += "    调试输出 (局_响应文本)\n\n"
	} else {
		return errors.New(fmt.Sprintf("无法生成[%s]的代码", module))
	}
	code += "}\n"
	code = strings.ReplaceAll(code, "\r", "")
	code = strings.ReplaceAll(code, "\n", "\r\n")
	return clipboard.ClipboardWriteAll(code)
}
func createPythonCode(app AppSession, module string) error {
	//  {Name: "requests", Type: GenerateCodeType_HTTP},
	if !app.IsHTTP() {
		return errors.New("这个请求不支持生成 Python 的 " + module + " 代码")
	}
	hp, ok := app.(*HttpSession)
	if !ok {
		return errors.New("这个请求不是HTTP请求")
	}
	URL, e := url.Parse(hp.Request.Url)
	if e != nil {
		return errors.New("解析URL失败")
	}
	header := hp.Request.Header.Clone()
	FuncName := "SunnyNetCreateRequest"
	Cookie := header.Get("Cookie")
	header.Del("Cookie")
	code := ""
	BytesType := IsBytesType(hp.Request.Body)
	sj := "    "
	if module == "requests" {
		_header := ""
		for k, v := range header {
			if strings.ToUpper(k) == "CONTENT-LENGTH" {
				continue
			}
			if len(v) < 1 {
				_header += fmt.Sprintf(`%s'%s': "",%s`, sj+sj, k, "\n")
			} else {
				_header += fmt.Sprintf(`%s'%s': "%s",%s`, sj+sj, k, strReplaceAll([]byte(v[0])), "\n")
			}
		}
		if Cookie != "" {
			_header += fmt.Sprintf(`%s'%s': "%s",%s`, sj, "Cookie", strReplaceAll([]byte(Cookie)), "\n")
		}
		payload := ""
		if BytesType {
			payload = fmt.Sprintf(`%sencoded_data = "%s"%spayload = base64.b64decode(encoded_data)%s`, sj, base64.StdEncoding.EncodeToString(hp.Request.Body), "\n"+sj, "\n")
		} else {
			payload = fmt.Sprintf(`%spayload = "%s"%s`, sj, strReplaceAll(hp.Request.Body), "\n")
		}
		code = fmt.Sprintf(`def %s():%s`, FuncName, "\n")
		code += fmt.Sprintf("%s\"\"\"\n", sj)
		code += fmt.Sprintf("%s[ %s ]\n", sj, URL.Path)
		code += fmt.Sprintf("%s本函数由SunnyNet网络中间件生成\n", sj)
		code += fmt.Sprintf("%s\"\"\"\n", sj)
		code += fmt.Sprintf("%surl = \"%s\"\n", sj, URL.String())
		code += payload
		code += fmt.Sprintf("%sheaders = {\n", sj)
		code += _header
		code += fmt.Sprintf("%s}\n", sj)
		code += fmt.Sprintf("%sresponse = requests.request(\"%s\", url, data = payload, headers = headers)\n", sj, hp.Request.Method)
		code += fmt.Sprintf("%sprint(response.text)\n", sj)

		code = strings.ReplaceAll(code, "\r", "")
		code = strings.ReplaceAll(code, "\n", "\r\n")
		return clipboard.ClipboardWriteAll(code)
	}
	return errors.New(fmt.Sprintf("无法生成[%s]的代码", module))
}
func CreateRequestCode(app AppSession, Lang string, module string) error {
	if app.IsUDP() {
		return errors.New("UDP不支持生成代码")
	}
	switch Lang {
	case "GoLang":
		return createGoCode(app, module)
	case "C#":
		return createCShaPCode(app, module)
	case "Python":
		return createPythonCode(app, module)
	case "火山":
		return createHSCode(app, module)
	case "易语言":
		return createECode(app, module)
	}
	return errors.New("找不到目标生成语言")
}
func IsBytesType(Body []byte) bool {
	for _, v := range Body {
		if v < 32 && v != 10 && v != 13 {
			return true
		}
	}
	return false
}
func strReplaceAll(body []byte) string {
	ss := strings.ReplaceAll(string(body), "\\", "\\\\")
	ss = strings.ReplaceAll(ss, "\"", "\\\"")
	ss = strings.ReplaceAll(ss, "\r", "\\r")
	ss = strings.ReplaceAll(ss, "\n", "\\n")
	return ss
}
func reText(string2 string) string {
	s := strings.ReplaceAll(string2, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\n", "\\n")
	return strings.ReplaceAll(s, "\"", "\\\"")
}
func isChinese(str string) bool {
	for _, v := range str {
		if v > 255 {
			return true
		}
	}
	return false
}

var ConvertELangFormat = convertELangFormat

func convertELangFormat(v string) string {
	str := v
	str = strings.ReplaceAll(str, "“", "\"")
	str = strings.ReplaceAll(str, "”", "\"")
	str = strings.ReplaceAll(str, "\r\n", "\n")
	str = strings.ReplaceAll(str, "\n", "\r\n")
	if v == "\r\n" || v == "\n" || v == "\r" {
		return "#换行符"
	}
	if !strings.Contains(str, "\"") && !strings.Contains(str, "\r\n") {
		return "“" + str + "”"
	}
	if strings.Contains(str, "\r\n") {
		arr := strings.Split(str, "\r\n")
		str = ""
		a := false
		for _, va := range arr {
			if !strings.Contains(va, "gzip, ") && !strings.Contains(va, "Content-Length:") && !strings.Contains(va, "Accept-Encoding: gzip") {
				str += va + "\r\n"
				a = true
			}
		}
		if a && strings.HasSuffix(str, "\r\n") {
			str = str[:len(str)-2]
		}
	}
	fh := ""
	issc := false
	if strings.Contains(str, "\"") {
		if !strings.Contains(str, "'") {
			fh = "'"
		} else if !strings.Contains(str, "#") {
			fh = "#"
		} else if !strings.Contains(str, "~") {
			fh = "~"
		} else if !strings.Contains(str, "!") {
			fh = "!"
		} else if !strings.Contains(str, "|") {
			fh = "|"
		} else if !strings.Contains(str, "/") {
			fh = "/"
		} else if !strings.Contains(str, "\\") {
			fh = "\\"
		} else if !strings.Contains(str, "&") {
			fh = "&"
		} else if !strings.Contains(str, "*") {
			fh = "*"
		} else {
			issc = true
			str = "\"" + strings.ReplaceAll(str, "\"", "\"+#引号+\"") + "\""
			fh = ""
		}
		if fh != "" {
			str = strings.ReplaceAll(str, "\"", fh)
		}
	}

	if strings.Contains(str, "\r\n") {
		if issc {
			str = strings.ReplaceAll(str, "\r\n", "\"+#换行符+\"")
		} else {
			str = "\"" + strings.ReplaceAll(str, "\r\n", "\"+#换行符+\"") + "\""
		}
	}
	str = strings.ReplaceAll(str, "+\"\"", "")
	str = strings.ReplaceAll(str, "\"\"", "")
	if strings.HasPrefix(str, "+") {
		str = str[1:]
	}
	if strings.HasPrefix(str, "+") {
		str = str[1:]
	}
	if strings.Contains(str, "++") {
		str = strings.ReplaceAll(str, "++", "+")
	}
	if fh != "" {
		if strings.HasPrefix(str, "\"") {
			str = "子文本替换 (" + str + ", \"" + fh + "\",#引号, , , 真)"
		} else {
			str = "子文本替换 (“" + str + "”, \"" + fh + "\",#引号, , , 真)"
		}
	}
	return str
}

func isFormData(e *HttpSession) (bool, string) {
	p := string(e.Request.Body)
	if !strings.Contains(p, "&") {
		return false, ""
	}
	if !strings.Contains(p, "=") {
		return false, ""
	}
	if strings.HasPrefix(p, "{") || strings.HasPrefix(p, "[") {
		return false, ""
	}
	Array := strings.Split(p, "&")
	Code := ""
	for index, v := range Array {
		Array1 := strings.Split(v, "=")
		if len(Array1) == 1 {
			Code += "局_提交数据 ＝ “" + Array1[0] + "=”\n"
		} else if len(Array1) == 2 {
			value, ex := url.QueryUnescape(Array1[1])
			if ex != nil {
				if index == 0 {
					Code += "局_提交数据 ＝ “" + Array1[0] + "=" + Array1[1] + "”\n"
				} else {
					Code += "局_提交数据 ＝ 局_提交数据 ＋ “&" + Array1[0] + "=" + Array1[1] + "”\n"
				}
				continue
			}
			if index == 0 {
				if isChinese(value) || value != Array1[1] {
					Code += "局_提交数据 ＝ “" + Array1[0] + "=” ＋ 编码_URL编码 (" + convertELangFormat(value) + ",真,真)\n"
				} else {
					Code += "局_提交数据 ＝ “" + Array1[0] + "=” ＋ " + convertELangFormat(value) + "\n"
				}
			} else {
				if isChinese(value) || value != Array1[1] {
					Code += "局_提交数据 ＝ 局_提交数据 ＋ “&" + Array1[0] + "=” ＋ 编码_URL编码 (" + convertELangFormat(value) + ",真,真)\n"
				} else {
					Code += "局_提交数据 ＝ 局_提交数据 ＋ “&" + Array1[0] + "=” ＋ " + convertELangFormat(value) + "\n"
				}
			}
		}
	}

	return true, Code
}
