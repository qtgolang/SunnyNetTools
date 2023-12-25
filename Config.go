package main

import (
	"bytes"
	"changeme/CommAnd"
	"changeme/MapHash"
	"changeme/Resource"
	"encoding/json"
	"fmt"
	"github.com/qtgolang/SunnyNet/SunnyNet"
	"github.com/qtgolang/SunnyNet/public"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"go/format"
	"io"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

func FormatCode(code string) string {
	formatted, err := format.Source([]byte(code))
	if err != nil {
		return code
	}
	return string(formatted)
}

const ScriptCode = `
func NewHttpSunny(uniqueId,Type,PID int, Url, Method string, Header http.Header, Body []byte, SetAgent func(ProxyUrl string) bool,Header2 http.Header, Body2 []byte, StateCode int, Display bool) (string, string, http.Header, []byte, http.Header, []byte, int, bool,bool) {
	inter := &HTTP{uniqueId:uniqueId,Type: Type,PID:PID, URL: Url, Method: Method, Header: Header, Body: Body, Display: Display,Break:false,SetAgent:SetAgent}
	inter.Response.Header = Header2
	inter.Response.Body = Body2
	inter.Response.StateCode = StateCode
	HttpRequestCallbackFunction(inter)
	return inter.URL, inter.Method, inter.Header, inter.Body, inter.Response.Header, inter.Response.Body, inter.Response.StateCode, inter.Display,inter.Break
}

func NewWebsocketSunny(uniqueId, Type,PID int, Url, Method string, Header http.Header, MessageType int, Body []byte, SendDataToServer func(MessageType int, data []byte) bool, SendDataToClient func(MessageType int, data []byte) bool, Close func() bool) ([]byte, bool) {
	inter := &WebSocket{uniqueId: uniqueId, Type: Type,PID:PID, URL: Url, Method: Method, Header: Header, MessageType: MessageType, Body: Body, Display: true}
	inter.SendDataToServer = SendDataToServer
	inter.SendDataToClient = SendDataToClient
	inter.Close = Close
	WebSocketCallbackFunction(inter)
	return inter.Body, inter.Display
}

func NewTCPSunny(uniqueId, Type,PID int, Body []byte, LocalAddress, RemoteAddress string, SetConnectionIP  func(NewAddress string) bool  ,	SetAgent func(ProxyUrl string) bool  ,SendDataToServer func(data []byte) bool, SendDataToClient func(data []byte) bool, Close func() bool) ([]byte, bool) {
	inter := &TcpSocket{uniqueId: uniqueId, Type: Type,PID:PID, Body: Body, Display: true, LocalAddress: LocalAddress, RemoteAddress: RemoteAddress,SetConnectionIP:SetConnectionIP,SetAgent:SetAgent}
	inter.SendDataToServer = SendDataToServer
	inter.SendDataToClient = SendDataToClient
	inter.Close = Close
	tcpCallbackFunction(inter)
	return inter.Body, inter.Display
}
func NewUDPSunny(uniqueId, Type,PID int, Body []byte, LocalAddress, RemoteAddress string) ([]byte, bool) {
	inter := &UdpSocket{uniqueId: uniqueId, Type: Type,PID:PID, Body: Body, Display: true, LocalAddress: LocalAddress, RemoteAddress: RemoteAddress} 
	udpCallbackFunction(inter)
	return inter.Body, inter.Display
}
var __lock sync.Mutex
var scriptLog = make([]any, 0, 1000)
func Log(msg ...interface{}) {
	__lock.Lock()
	if len(msg) > 1 {
		scriptLog = append(scriptLog, msg...)
	} else if len(msg) < 1 {
		__lock.Unlock()
		return
	} else {
		scriptLog = append(scriptLog, msg[0])
	} 
	if len(scriptLog) > 1000 {
		scriptLog = scriptLog[1:len(scriptLog)]
	}
	__lock.Unlock()
}
 `

var NewEval *interp.Interpreter
var lock sync.Mutex
var httpFunc func(uniqueId, Type, PID int, Url, Method string, Header http.Header, Body []byte, SetAgent func(ProxyUrl string) bool, Header2 http.Header, Body2 []byte, StateCode int, Display bool) (string, string, http.Header, []byte, http.Header, []byte, int, bool, bool)
var wsFunc func(uniqueId, Type, PID int, Url, Method string, Header http.Header, MessageType int, Body []byte, SendDataToServer func(MessageType int, data []byte) bool, SendDataToClient func(MessageType int, data []byte) bool, Close func() bool) ([]byte, bool)
var tcpFunc func(uniqueId, Type, PID int, Body []byte, LocalAddress, RemoteAddress string, SetConnectionIP func(NewAddress string) bool, SetAgent func(ProxyUrl string) bool, SendDataToServer func(data []byte) bool, SendDataToClient func(data []byte) bool, Close func() bool) ([]byte, bool)
var udpFunc func(uniqueId, Type, PID int, Body []byte, LocalAddress, RemoteAddress string) ([]byte, bool)

func extractImport(s string) map[string]bool {
	arrayMap := make(map[string]bool)
	str := ""
	start := false
	start2 := false
	for _, v := range s {
		if v == '\n' {
			if strings.HasPrefix(str, "import ") {
				str = strings.ReplaceAll(str, "import ", "")
				arrayMap[strings.TrimSpace(str)] = true
				str = ""
				continue
			}
			if str != "" && start {
				arrayMap[strings.TrimSpace(str)] = true
			}
			str = ""
		} else {
			if start2 && string(v) == "\"" {
				arrayMap[strings.TrimSpace(" \""+str+"\"")] = true
				str = ""
				start2 = false
				continue
			}
			str += string(v)
			if start && str == ")" {
				start = false
			}
			if str == "import (" {
				start = true
				str = ""
				continue
			}
			if str == "import \"" {
				start2 = true
				str = ""
				continue
			}

		}
	}
	return arrayMap
}
func extractCodeBody(s string) string {
	str := ""
	res := ""
	statr1 := false
	statr2 := false
	statr3 := false
	for _, v := range s {
		if v == '\n' {
			if statr3 {
				str = ""
				statr3 = false
				continue
			}
			if strings.HasPrefix(str, "package ") {
				str = ""
				continue
			}
			if strings.HasPrefix(str, "import (") {
				str = ""
				statr1 = true
				continue
			}
			if strings.HasPrefix(str, "import \"") {
				str = ""
				statr2 = true
				continue
			}
			if strings.HasPrefix(str, "import ") {
				str = ""
				statr3 = true
				continue
			}
			if !statr1 && !statr2 && !statr3 {
				res += str + "\n"
				str = ""
			}
		} else {
			if statr1 {
				if string(v) == ")" {
					statr1 = false
					str = ""
					continue
				}
			}
			if statr2 {
				if string(v) == "\"" {
					statr2 = false
					str = ""
					continue
				}
			}
			str += string(v)
		}
	}
	return res
}
func RunCode() (SErr string) {
	var iEval = interp.New(interp.Options{})
	iEval.Use(stdlib.Symbols)
	ca := string(GlobalConfig.GoScriptCode) + ScriptCode
	//分析出用户编写的脚本中引用的包
	UserImport := extractImport(ca)
	src := string(Resource.GoBuiltFuncCode)
	//分析内置函数引用的包
	SystemPort := extractImport(src)
	for k, _ := range SystemPort {
		if UserImport[k] == false {
			_, _ = iEval.Eval("import " + k)
		}
	}
	CodeBody := extractCodeBody(src)
	S := ca + CodeBody
	_, err := iEval.Eval(S)
	if err != nil {
		errorSrc := strings.ReplaceAll(err.Error(), "_.go:", "")
		ar := strings.Split(errorSrc, "error: unable to find source related to:")
		if len(ar) >= 2 {
			ar1 := strings.Split(ar[0], ": import")
			if len(ar1) >= 2 {
				return "错误位置:" + ar1[0] + " 找不到引入包 [ " + ar1[1] + " ]"
			}
		}
		ar = strings.Split(errorSrc, ": expected declaration, found")
		if len(ar) >= 2 {
			return "错误位置:" + ar[0] + " 无效的字符 [ " + ar[1] + " ]"
		}
		ar = strings.Split(errorSrc, ": expected ';', found")
		if len(ar) >= 2 {
			ar1 := strings.Split(ar[1], " (and")
			if len(ar1) > 1 {
				return "错误位置:" + ar[0] + " 无效的字符 [ " + ar1[0] + " ]"
			}
			return "错误位置:" + ar[0] + " 无效的字符 [ " + ar[1] + " ]"
		}
		ar = strings.Split(errorSrc, ": undefined: ")
		if len(ar) >= 2 {
			return "错误位置:" + ar[0] + " 未定义的 [ " + ar[1] + " ]"
		}
		ar = strings.Split(errorSrc, ": expected operand, found")
		if len(ar) >= 2 {
			return "错误位置:" + ar[0] + " 参数不正确 请检查传递的参数"
		}
		ar = strings.Split(errorSrc, ": undefined selector: ")
		if len(ar) >= 2 {
			return "错误位置:" + ar[0] + " 未定义的属性 [ " + ar[1] + " ]"
		}
		ar = strings.Split(errorSrc, ":")
		if len(ar) >= 2 {
			like, _ := strconv.Atoi(ar[0])
			like2 := len(strings.Split(string(GlobalConfig.GoScriptCode), "\n"))
			if like > like2 {
				return "错误: 默认结构体已被更改,请检查代码"
			}
		}
		ar = strings.Split(errorSrc, ": illegal character ")
		if len(ar) >= 2 {
			ar1 := strings.Split(errorSrc, " ")
			if len(ar1) >= 1 {
				return "错误位置:" + ar[0] + " 非法字符 " + ar1[len(ar1)-1] + " "
			}
			return "错误位置:" + ar[0] + " 非法字符 " + ar[1]
		}
		if strings.Index(errorSrc, ": package ") != -1 && strings.Index(errorSrc, "has no symbol ") != -1 {
			ar = strings.Split(errorSrc, ": package ")
			if len(ar) >= 2 {
				pos := ar[0]
				ar = strings.Split(ar[1], " ")
				if len(ar) >= 2 {
					ar = strings.Split(ar[1], " ")
					pack := ar[0]
					ar = strings.Split(errorSrc, "has no symbol ")
					if len(ar) >= 2 {
						funcName := ar[1]
						return "错误位置:" + pos + " 在包 " + pack + " 中 找不到函数 -> \"" + funcName + "\""
					}
				}
			}
		}
		return "错误位置:" + errorSrc
	}
	v, err := iEval.Eval("main.NewHttpSunny")
	if err != nil {
		return err.Error()
	}

	_httpFunc := v.Interface().(func(uniqueId, Type, PID int, Url, Method string, Header http.Header, Body []byte, SetAgent func(ProxyUrl string) bool, Header2 http.Header, Body2 []byte, StateCode int, Display bool) (string, string, http.Header, []byte, http.Header, []byte, int, bool, bool))
	if _httpFunc == nil {
		return "找不到NewHttpSunny"
	}
	defer func() {
		if p := recover(); p != nil {
			SErr = fmt.Sprintf("%v", p)
		}
	}()
	_httpFunc(0, 1, 0, "https://test.com", "GET", make(http.Header), []byte("点击测试脚本代码:HTTP回调函数"), func(ProxyUrl string) bool {
		return false
	}, make(http.Header), make([]byte, 0), 0, false)

	v, err = iEval.Eval("main.NewWebsocketSunny")
	if err != nil {
		return err.Error()
	}
	_wsFunc := v.Interface().(func(uniqueId, Type, PID int, Url, Method string, Header http.Header, MessageType int, Body []byte, SendDataToServer func(MessageType int, data []byte) bool, SendDataToClient func(MessageType int, data []byte) bool, Close func() bool) ([]byte, bool))
	if _wsFunc == nil {
		return "找不到NewHttpSunny"
	}
	defer func() {
		if p := recover(); p != nil {
			SErr = fmt.Sprintf("%v", p)
		}
	}()
	_wsFunc(0, 1, 0, "https://test.com", "GET", make(http.Header), 1, []byte("点击测试脚本代码:WS回调函数"), func(MessageType int, data []byte) bool {
		return false
	}, func(MessageType int, data []byte) bool {
		return false
	}, func() bool {
		return false
	})

	v, err = iEval.Eval("main.NewTCPSunny")
	if err != nil {
		return err.Error()
	}
	_tcpFunc := v.Interface().(func(uniqueId, Type, PID int, Body []byte, LocalAddress, RemoteAddress string, SetConnectionIP func(NewAddress string) bool, SetAgent func(ProxyUrl string) bool, SendDataToServer func(data []byte) bool, SendDataToClient func(data []byte) bool, Close func() bool) ([]byte, bool))
	if _tcpFunc == nil {
		return "找不到NewHttpSunny"
	}
	defer func() {
		if p := recover(); p != nil {
			SErr = fmt.Sprintf("%v", p)
		}
	}()
	_tcpFunc(0, 1, 0, []byte("点击测试脚本代码:tcp回调函数"), "TCP_Test_LocalAddress", "TCP_Test_RemoteAddress", func(NewAddress string) bool {
		return false
	}, func(ProxyUrl string) bool {
		return false
	}, func(data []byte) bool {
		return false
	}, func(data []byte) bool {
		return false
	}, func() bool {
		return false
	})

	v, err = iEval.Eval("main.NewUDPSunny")
	if err != nil {
		return err.Error()
	}
	_udpFunc := v.Interface().(func(uniqueId, Type, PID int, Body []byte, LocalAddress, RemoteAddress string) ([]byte, bool))
	if _tcpFunc == nil {
		return "找不到NewHttpSunny"
	}
	defer func() {
		if p := recover(); p != nil {
			SErr = fmt.Sprintf("%v", p)
		}
	}()
	_udpFunc(0, 1, 0, []byte("点击测试脚本代码:udp回调函数"), "UDP_Test_LocalAddress", "UDP_Test_RemoteAddress")

	lock.Lock()
	NewEval = iEval
	httpFunc = _httpFunc
	wsFunc = _wsFunc
	tcpFunc = _tcpFunc
	udpFunc = _udpFunc
	lock.Unlock()
	return ""
}
func RunCodeLog() (Str string) {
	defer func() {
		if p := recover(); p != nil {
			Str = fmt.Sprintf("error: %v", p)
		}
	}()
	lock.Lock()
	defer lock.Unlock()
	__log, err := NewEval.Eval("scriptLog")
	if err != nil {
		return Str
	}
	scriptLog := __log.Interface().([]any)
	for _, v := range scriptLog {
		Str += fmt.Sprintf("%v", v) + "\r\n"
	}
	return Str
}

func RunHTTPRequestScriptCode(Conn *SunnyNet.HttpConn) bool {
	h := &MapHash.Request{Display: true}
	h.URL = Conn.Request.URL.String()
	h.Header = Conn.Request.Header.Clone()
	h.Method = Conn.Request.Method
	h.Proto = "HTTP/1.1"
	if Conn.Request.Body != nil {
		Body, _ := io.ReadAll(Conn.Request.Body)
		_ = Conn.Request.Body.Close()
		h.Body = Body
		Conn.Request.Body = io.NopCloser(bytes.NewBuffer(Body))
	}
	h.Conn = Conn

	lock.Lock()
	_Call := httpFunc
	lock.Unlock()
	if _Call == nil {
		h.Break = 0
		h.Display = true
		h.Way = "HTTP"
		h.PID = CommAnd.GetPidName(Conn.PID)
		h.ClientIP = Conn.ClientIP
		HashMap.SetRequest(Conn.Theology, h)
		return true
	}
	if Conn.Response != nil {
		if Conn.Response.Body != nil {
			Body, _ := io.ReadAll(Conn.Response.Body)
			_ = Conn.Response.Body.Close()
			Conn.Response.Body = io.NopCloser(bytes.NewBuffer(Body))
			h.Response.Body = Body
		} else {
			h.Response.Body = []byte{}
		}
		h.Response.StateCode = Conn.Response.StatusCode
		h.Response.Header = Conn.Response.Header
	}
	_URL, _Method, _Header, _Body, _Header2, _Body2, _StateCode, Display, _Break := _Call(Conn.Theology, Conn.Type, Conn.PID, h.URL, h.Method, h.Header, h.Body, Conn.SetAgent, h.Response.Header, h.Response.Body, h.Response.StateCode, h.Display)
	h.Response.StateCode = _StateCode
	if len(_Header2) > 0 {
		Conn.Response = new(http.Response)
		h.Response.Header = _Header2
		Conn.Response.Header = _Header2
		if h.Response.StateCode < 100 {
			h.Response.StateCode = 200
		}
		Conn.Response.StatusCode = h.Response.StateCode
	}
	if len(_Body2) > 0 {
		if Conn.Response == nil {
			Conn.Response = new(http.Response)
		}
		h.Response.Body = _Body2
		if Conn.Response.Body != nil {
			_ = Conn.Response.Body.Close()
		}
		Conn.Response.Body = io.NopCloser(bytes.NewBuffer(_Body2))
		if h.Response.StateCode < 100 {
			h.Response.StateCode = 200
		}
		Conn.Response.StatusCode = h.Response.StateCode
		if Conn.Response.Header == nil {
			Conn.Response.Header = make(http.Header)
			Conn.Response.Header.Set("Server", "Sunny")
			Conn.Response.Header.Set("Accept-Ranges", "bytes")
			Conn.Response.Header.Set("Connection", "Close")
		}
		Conn.Response.Header.Set("Content-Length", strconv.Itoa(len(_Body2)))
		h.Response.Header = Conn.Response.Header
		Conn.Response.ContentLength = int64(len(_Body2))
	}
	h.Display = Display
	if _URL != h.URL {
		h.URL = _URL
		Conn.Request.URL, _ = url.Parse(_URL)
	}
	if _Break {
		h.Break = 1
	} else {
		h.Break = 0
	}
	if _Method != h.Method {
		h.Method = _Method
		Conn.Request.Method = _Method
	}
	if !mapsEqual(_Header, h.Header) {
		h.Header = _Header
		Conn.Request.Header = _Header
	}
	if !bytes.Equal(_Body, h.Body) {
		h.Body = _Body
		if Conn.Request.Body != nil {
			_ = Conn.Request.Body.Close()
		}
		Conn.Request.Body = io.NopCloser(bytes.NewBuffer(_Body))
	}
	if Display {
		h.Way = "HTTP"
		h.PID = CommAnd.GetPidName(Conn.PID)
		h.ClientIP = Conn.ClientIP
		HashMap.SetRequest(Conn.Theology, h)
	}
	return Display
}
func RunHTTPResponseScriptCode(Conn *SunnyNet.HttpConn) bool {
	URL := Conn.Request.URL.String()
	Header := Conn.Request.Header.Clone()
	Method := Conn.Request.Method
	Body := make([]byte, 0)
	if Conn.Request.Body != nil {
		Body, _ = io.ReadAll(Conn.Request.Body)
		_ = Conn.Request.Body.Close()
		Conn.Request.Body = io.NopCloser(bytes.NewBuffer(Body))
	}
	StateCode := Conn.Response.StatusCode
	Body2 := make([]byte, 0)
	Header2 := Conn.Response.Header.Clone()
	if Conn.Response.Body != nil {
		Body2, _ = io.ReadAll(Conn.Response.Body)
		_ = Conn.Response.Body.Close()
		Conn.Response.Body = io.NopCloser(bytes.NewBuffer(Body2))
	}
	lock.Lock()
	_Call := httpFunc
	lock.Unlock()
	if _Call == nil {
		return false
	}
	_, _, _, _, _Header, _Body, _StateCode, _, _Break := _Call(Conn.Theology, Conn.Type, Conn.PID, URL, Method, Header, Body, Conn.SetAgent, Header2, Body2, StateCode, true)

	if Conn.Response.StatusCode != _StateCode {
		Conn.Response.StatusCode = _StateCode
	}
	if !mapsEqual(_Header, Header2) {
		Header2 = _Header
		Conn.Response.Header = _Header
	}
	if !bytes.Equal(_Body, Body2) {
		Body2 = _Body
		if Conn.Response.Body != nil {
			_ = Conn.Response.Body.Close()
		}
		Conn.Response.Body = io.NopCloser(bytes.NewBuffer(_Body))
	}
	h1 := HashMap.GetRequest(Conn.Theology)
	if h1 == nil {
		return false
	}
	h1.Response.Header = _Header
	h1.Response.Body = _Body
	h1.Response.StateCode = _StateCode
	h1.Response.Conn = Conn
	return _Break
}
func RunHTTPErrorScriptCode(Conn *SunnyNet.HttpConn) {
	URL := Conn.Request.URL.String()
	Header := Conn.Request.Header.Clone()
	Method := Conn.Request.Method
	Body := make([]byte, 0)
	if Conn.Request.Body != nil {
		Body, _ = io.ReadAll(Conn.Request.Body)
		_ = Conn.Request.Body.Close()
		Conn.Request.Body = io.NopCloser(bytes.NewBuffer(Body))
	}
	Body2 := []byte(Conn.GetError())
	lock.Lock()
	_Call := httpFunc
	lock.Unlock()
	if _Call == nil {
		return
	}
	_Call(Conn.Theology, Conn.Type, Conn.PID, URL, Method, Header, Body, Conn.SetAgent, make(http.Header), Body2, -1, true)
}
func RunWebSocketScriptCode(Conn *SunnyNet.WsConn) bool {
	URL := Conn.Request.URL.String()
	Header := Conn.Request.Header.Clone()
	Method := Conn.Request.Method
	Body := Conn.GetMessageBody()
	lock.Lock()
	_Call := wsFunc
	lock.Unlock()
	if _Call == nil {
		return true
	}
	Body2, Break := _Call(Conn.Theology, Conn.Type, Conn.Pid, URL, Method, Header, Conn.GetMessageType(), Body, Conn.SendToServer, Conn.SendToClient, Conn.Close)
	Conn.SetMessageBody(Body2)
	return Break
}
func RunTcpScriptCode(Conn *SunnyNet.TcpConn) bool {
	lock.Lock()
	_Call := tcpFunc
	lock.Unlock()
	if _Call == nil {
		return true
	}
	_Type := 0
	if Conn.Type == public.SunnyNetMsgTypeTCPAboutToConnect {
		_Type = 5
	} else if Conn.Type == public.SunnyNetMsgTypeTCPConnectOK {
		_Type = 1
	} else if Conn.Type == public.SunnyNetMsgTypeTCPClientSend {
		_Type = 2
	} else if Conn.Type == public.SunnyNetMsgTypeTCPClientReceive {
		_Type = 3
	} else if Conn.Type == public.SunnyNetMsgTypeTCPClose {
		_Type = 4
	}
	Body := Conn.GetBody()
	Body2, Break := _Call(Conn.Theology, _Type, Conn.Pid, Body, Conn.LocalAddr, Conn.RemoteAddr, Conn.SetConnectionIP, Conn.SetAgent, func(data []byte) bool { return Conn.SendToServer(data) != 0 }, func(data []byte) bool { return Conn.SendToClient(data) != 0 }, Conn.Close)
	if _Type == 2 || _Type == 3 || _Type == 5 {
		Conn.SetBody(Body2)
	}
	return Break
}
func RunUdpScriptCode(Conn *SunnyNet.UDPConn) bool {
	lock.Lock()
	_Call := udpFunc
	lock.Unlock()
	if _Call == nil {
		return true
	}
	_Type := 0
	if Conn.Type == public.SunnyNetUDPTypeClosed {
		_Type = 3
	} else if Conn.Type == public.SunnyNetUDPTypeSend {
		_Type = 1
	} else {
		_Type = 2
	}
	Body2, Break := _Call(int(Conn.Theology), _Type, Conn.Pid, Conn.Data, Conn.LocalAddress, Conn.RemoteAddress)
	Conn.Data = Body2
	return Break
}
func mapsEqual(a, b http.Header) bool {
	if len(a) != len(b) {
		return false
	}
	for key, valueA := range a {
		valueB, ok := b[key]
		if !ok || len(valueA) != len(valueB) {
			return false
		}
		for i := range valueA {
			if valueA[i] != valueB[i] {
				return false
			}
		}
	}
	return true
}

type Color struct {
	Dark  string `json:"dark"`
	Right string `json:"right"`
}
type ConfigReplaceRules struct {
	Type string `json:"Type"`
	Src  string `json:"Src"`
	Dest string `json:"Dest"`
	Hash string `json:"Hash"`
}
type ConfigRequestCertManager struct {
	Rule     uint8  `json:"rule"`
	FilePath string `json:"FilePath"`
	PassWord string `json:"PassWord"`
	Host     string `json:"Host"`
}
type UserConfig struct {
	ColorConfig struct {
		TCP      Color `json:"tcp"`
		UDP      Color `json:"udp"`
		WS       Color `json:"ws"`
		CSS      Color `json:"css"`
		JS       Color `json:"js"`
		IMG      Color `json:"img"`
		Text     Color `json:"document"`
		State1   Color `json:"_1"`
		State301 Color `json:"_301"`
		State302 Color `json:"_302"`
		State401 Color `json:"_401"`
		State403 Color `json:"_403"`
		State404 Color `json:"_404"`
		State500 Color `json:"_500"`
	} `json:"ColorConfig"`
	GoScriptCode           []byte               `json:"ScriptCode"`
	Port                   int                  `json:"Port"`
	DisableUDP             bool                 `json:"DisableUDP"`
	DisableCache           bool                 `json:"DisableCache"`
	Authentication         bool                 `json:"OpenAuthentication"`
	AuthenticationUserInfo map[string]string    `json:"AuthenticationUserInfo"`
	GlobalProxy            string               `json:"GlobalProxy"`
	GlobalProxyRules       string               `json:"GlobalProxyRules"`
	ReplaceRules           []ConfigReplaceRules `json:"ReplaceRules"`
	HostsRules             []ConfigReplaceRules `json:"HostsRules"`
	DarkTheme              uint8                `json:"DarkTheme"`
	Filter                 string               `json:"Filter"`
	KeysStrings            string               `json:"KeysStrings"`
	Columns                string               `json:"Columns"`
	MustTcp                struct {
		Open  bool   `json:"open"`
		Rules string `json:"Rules"`
	} `json:"MustTcp"`
	Size struct {
		Width  int `json:"Width"`
		Height int `json:"Height"`
	} `json:"Size"`
	Cert struct {
		Default bool   `json:"Default"`
		CaPath  string `json:"CaPath"`
		KeyPath string `json:"KeyPath"`
	} `json:"Cert"`
	RequestCertManager map[int]ConfigRequestCertManager `json:"RequestCertManager"`
	GOOS               string                           `json:"GOOS"`
}

func (c *UserConfig) loadDefaultValue() {
	c.GOOS = runtime.GOOS
	if c.RequestCertManager == nil {
		c.RequestCertManager = make(map[int]ConfigRequestCertManager)
	}
	if c.Size.Width < 822 {
		c.Size.Width = 1540
	}
	if c.Size.Height < 387 {
		c.Size.Height = 700
	}
	if c.Filter == "" {
		c.Filter = "{\"响应长度\":{\"filterType\":\"text\",\"type\":\"notContains\",\"filter\":\"0/0\"},\"响应类型\":{\"filterType\":\"text\",\"type\":\"notEqual\",\"filter\":\"error\"}}"
	}
	if c.DarkTheme == 0 {
		c.DarkTheme = 1
	}
	if c.ReplaceRules == nil {
		c.ReplaceRules = make([]ConfigReplaceRules, 0)
	}
	if c.HostsRules == nil {
		c.HostsRules = make([]ConfigReplaceRules, 0)
	}
	if c.AuthenticationUserInfo == nil {
		c.AuthenticationUserInfo = make(map[string]string)
	}
	//证书选择使用
	{
		if c.Cert.Default == false {
			if c.Cert.CaPath == "" {
				c.Cert.Default = true
			}
			if c.Cert.KeyPath == "" {
				c.Cert.Default = true
			}
		}
		if c.Cert.Default == false {
			bs1, _ := os.ReadFile(c.Cert.CaPath)
			bs2, _ := os.ReadFile(c.Cert.KeyPath)
			if len(bs1) < 1 && len(bs2) < 1 {
				c.Cert.Default = true
			}
		}
	}
	//判断加载默认的列表颜色配置
	{
		if c.ColorConfig.TCP.Dark == "" {
			c.ColorConfig.TCP.Dark = "#B570FF"
		}
		if c.ColorConfig.TCP.Right == "" {
			c.ColorConfig.TCP.Right = "#7B00FF"
		}

		if c.ColorConfig.UDP.Dark == "" {
			c.ColorConfig.UDP.Dark = "#B570FF"
		}
		if c.ColorConfig.UDP.Right == "" {
			c.ColorConfig.UDP.Right = "#7B00FF"
		}

		if c.ColorConfig.WS.Dark == "" {
			c.ColorConfig.WS.Dark = "#B570FF"
		}
		if c.ColorConfig.WS.Right == "" {
			c.ColorConfig.WS.Right = "#7B00FF"
		}

		if c.ColorConfig.CSS.Dark == "" {
			c.ColorConfig.CSS.Dark = "#E676EA"
		}
		if c.ColorConfig.CSS.Right == "" {
			c.ColorConfig.CSS.Right = "#F700FF"
		}

		if c.ColorConfig.JS.Dark == "" {
			c.ColorConfig.JS.Dark = "#4AC432"
		}
		if c.ColorConfig.JS.Right == "" {
			c.ColorConfig.JS.Right = "#26A34E"
		}

		if c.ColorConfig.IMG.Dark == "" {
			c.ColorConfig.IMG.Dark = "#617D73"
		}
		if c.ColorConfig.IMG.Right == "" {
			c.ColorConfig.IMG.Right = "#9A9A9A"
		}

		if c.ColorConfig.Text.Dark == "" {
			c.ColorConfig.Text.Dark = "#5387F8"
		}
		if c.ColorConfig.Text.Right == "" {
			c.ColorConfig.Text.Right = "#024AFF"
		}

		if c.ColorConfig.State1.Dark == "" {
			c.ColorConfig.State1.Dark = "#FE5151"
		}
		if c.ColorConfig.State1.Right == "" {
			c.ColorConfig.State1.Right = "#FF0000"
		}

		if c.ColorConfig.State401.Dark == "" {
			c.ColorConfig.State401.Dark = "#FE5151"
		}
		if c.ColorConfig.State401.Right == "" {
			c.ColorConfig.State401.Right = "#FF0000"
		}

		if c.ColorConfig.State403.Dark == "" {
			c.ColorConfig.State403.Dark = "#FE5151"
		}
		if c.ColorConfig.State403.Right == "" {
			c.ColorConfig.State403.Right = "#FF0000"
		}

		if c.ColorConfig.State404.Dark == "" {
			c.ColorConfig.State404.Dark = "#FE5151"
		}
		if c.ColorConfig.State404.Right == "" {
			c.ColorConfig.State404.Right = "#FF0000"
		}

		if c.ColorConfig.State500.Dark == "" {
			c.ColorConfig.State500.Dark = "#FE5151"
		}
		if c.ColorConfig.State500.Right == "" {
			c.ColorConfig.State500.Right = "#FF0000"
		}
		if c.ColorConfig.State302.Dark == "" {
			c.ColorConfig.State302.Dark = "#DBE617"
		}
		if c.ColorConfig.State302.Right == "" {
			c.ColorConfig.State302.Right = "#B7A807"
		}
		if c.ColorConfig.State301.Dark == "" {
			c.ColorConfig.State301.Dark = "#DBE617"
		}
		if c.ColorConfig.State301.Right == "" {
			c.ColorConfig.State301.Right = "#B7A807"
		}

	}
	//判断加载默认的脚本代码
	{
		if len(c.GoScriptCode) < 1 {
			c.GoScriptCode = []byte(string(Resource.GoCode) + Resource.ScriptAnnotation)
		}
		if RunCode() != "" {
			c.GoScriptCode = []byte(string(Resource.GoCode) + Resource.ScriptAnnotation)
		}
	}
	//基本设置
	{
		if c.Port < 1 || c.Port > 65535 {
			c.Port = 2024
		}
	}
	//上游代理
	{
		if c.GlobalProxyRules == "" {
			c.GlobalProxyRules = "www.baidu.com;*.ip138.com"
		}
	}
	//上游代理
	{
		if c.MustTcp.Rules == "" {
			c.MustTcp.Rules = "*.qqww.com;ww.baidu.com;google.com"
		}
	}
}
func (c *UserConfig) ResetAll() string {
	configLock.Lock()
	homeDir, err1 := os.UserHomeDir()
	if err1 != nil {
		configLock.Unlock()
		return "获取用户目录失败:" + err1.Error()
	}
	os.Remove(homeDir + "/Sunny/Config.json")
	configLock.Unlock()
	return ""
}
func (c *UserConfig) saveToFile() error {
	homeDir, err1 := os.UserHomeDir()
	if err1 != nil {
		panic("获取用户目录失败:" + err1.Error())
	}
	_ = os.Mkdir(homeDir+"/Sunny", 0777)
	if fileExists(homeDir + "/Sunny/Config.json") {
		err := os.Remove(homeDir + "/Sunny/Config.json")
		if err != nil {
			return err
		}
	}
	bs, e := json.MarshalIndent(c, "", "\t")
	if e != nil {
		return e
	}
	e = os.WriteFile(homeDir+"/Sunny/Config.json", bs, 777)
	return e
}
func (c *UserConfig) LoadLocalFile() {
	configLock.Lock()
	defer configLock.Unlock()
	homeDir, err1 := os.UserHomeDir()
	if err1 != nil {
		panic("获取用户目录失败:" + err1.Error())
	}
	_ = os.Mkdir(homeDir+"/Sunny", 0777)
	bs, _ := os.ReadFile(homeDir + "/Sunny/Config.json")
	json.Unmarshal(bs, &c)
	c.loadDefaultValue()
}
func (c *UserConfig) SaveColorConfig(Data string) error {
	configLock.Lock()
	defer configLock.Unlock()
	_ = json.Unmarshal([]byte(Data), &c.ColorConfig)
	c.loadDefaultValue()
	return c.saveToFile()
}
func (c *UserConfig) ResetColorConfig(dark bool) error {
	configLock.Lock()
	defer configLock.Unlock()
	if dark {
		c.ColorConfig.TCP.Dark = ""
		c.ColorConfig.UDP.Dark = ""
		c.ColorConfig.WS.Dark = ""
		c.ColorConfig.CSS.Dark = ""
		c.ColorConfig.JS.Dark = ""
		c.ColorConfig.IMG.Dark = ""
		c.ColorConfig.Text.Dark = ""
		c.ColorConfig.State1.Dark = ""
		c.ColorConfig.State301.Dark = ""
		c.ColorConfig.State302.Dark = ""
		c.ColorConfig.State401.Dark = ""
		c.ColorConfig.State403.Dark = ""
		c.ColorConfig.State404.Dark = ""
		c.ColorConfig.State500.Dark = ""
	} else {
		c.ColorConfig.TCP.Right = ""
		c.ColorConfig.UDP.Right = ""
		c.ColorConfig.WS.Right = ""
		c.ColorConfig.CSS.Right = ""
		c.ColorConfig.JS.Right = ""
		c.ColorConfig.IMG.Right = ""
		c.ColorConfig.Text.Right = ""
		c.ColorConfig.State1.Right = ""
		c.ColorConfig.State301.Right = ""
		c.ColorConfig.State302.Right = ""
		c.ColorConfig.State401.Right = ""
		c.ColorConfig.State403.Right = ""
		c.ColorConfig.State404.Right = ""
		c.ColorConfig.State500.Right = ""
	}
	c.loadDefaultValue()
	return c.saveToFile()
}

var configLock sync.Mutex
var GlobalConfig *UserConfig

func init() {
	GlobalConfig = &UserConfig{}
	GlobalConfig.LoadLocalFile()
}
func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
