package main

import (
	"bytes"
	"changeme/CommAnd"
	"changeme/MapHash"
	"compress/flate"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"github.com/andybalholm/brotli"
	"github.com/qtgolang/SunnyNet/SunnyNet"
	"github.com/qtgolang/SunnyNet/public"
	"github.com/qtgolang/SunnyNet/src/protobuf/JSON"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

func CallJs(command string, Args any) {
	app.CallDo(&Command{Command: command, Args: Args})
}

type AlertMsg struct {
	MessageText string `json:"msg"`
	Title       string `json:"title"`
}

func CallJsAlert(Title, MessageText string) {
	app.CallDo(&Command{Command: "弹出提示消息", Args: &AlertMsg{MessageText: MessageText, Title: Title}})
}

var HashMap = MapHash.NewHashMap()

func SaveData(args *JSON.SyJson) any {
	ConfigurationName := args.GetData("Type")
	Data := args.GetData("Data")
	switch ConfigurationName {
	case "列表颜色配置":
		e := GlobalConfig.SaveColorConfig(Data)
		if e != nil {
			CallJs("弹出错误提示", "保存配置失败:"+e.Error())
		}
		return true
	case "列数据":
		code, _ := base64.StdEncoding.DecodeString(strings.ReplaceAll(Data, "\\\\", "\\"))
		_TmpLock.Lock()
		GlobalConfig.Columns = string(code)
		_ = GlobalConfig.saveToFile()
		_TmpLock.Unlock()
		return true
	}
	return true
}

type ListInfo struct {
	MessageId int    `json:"MessageId"` //MessageId
	Theology  int    `json:"Theology"`  //唯一ID
	State     string `json:"状态"`        //状态文本
	URL       string `json:"请求地址"`      //请求地址
	HOST      string `json:"HOST"`      //请求地址
	ClientIP  string `json:"来源地址"`      //来源地址
	PID       string `json:"进程"`        //进程信息
	Method    string `json:"方式"`        //方式
	Ico       string `json:"ico"`       //显示图标
	Len       string `json:"响应长度"`      //显示图标
	Type      string `json:"响应类型"`      //显示图标
	SendTime  string `json:"请求时间"`      //显示图标
	RecTime   string `json:"响应时间"`      //显示图标
	Notes     string `json:"注释"`        //显示图标
	Color     struct {
		TagColor string `json:"TagColor"` //标记的文本颜色
		Search   string `json:"search"`   //搜索的背景颜色
	} `json:"color"` //显示图标
	Break uint8 `json:"断点模式"`
}
type UpdateCurrentResponse struct {
	Theology  int         `json:"Theology"` //唯一ID
	Header    http.Header `json:"Header"`
	Body      []byte      `json:"Body"`
	StateText string      `json:"StateText"`
	StateCode int         `json:"StateCode"`
	Break     bool        `json:"断点状态"`
	Error     bool        `json:"Error"`
	WebSocket bool        `json:"WebSocket"`
}
type UpdateICO struct {
	Theology int    `json:"Theology"` //唯一ID
	Ico      string `json:"ico"`      //显示图标
}

func init() {
	go InsertList()
	go UpdateResponseLength()
	GlobalConfig.LoadLocalFile()
}

var Insert sync.Mutex
var InsertData []any
var InsertDataMapTag = make(map[int]bool)
var UpdateData []any
var UpdateListICO []any
var SocketData []any
var currentlySelected = -1
var StatusText = "SunnyNetStatusText"
var breakpoint = 0
var SearchPercentage = -1

func AddInsertList(list *ListInfo) {
	Insert.Lock()
	if InsertDataMapTag[list.Theology] == false {
		InsertDataMapTag[list.Theology] = true
		InsertData = append(InsertData, list)
	}
	Insert.Unlock()
}
func InsertList() {
	for {
		time.Sleep(100 * time.Millisecond)
		Insert.Lock()
		if len(InsertData) > 0 {
			CallJs("插入列表", InsertData)
			InsertData = make([]any, 0)
			InsertDataMapTag = make(map[int]bool)
		}
		if len(UpdateData) > 0 {
			CallJs("更新列表", UpdateData)
			UpdateData = make([]any, 0)
		}
		if len(SocketData) > 0 {
			CallJs("更新Socket", SocketData)
			SocketData = make([]any, 0)
		}
		if len(UpdateListICO) > 0 {
			CallJs("更新ICO", UpdateListICO)
			UpdateListICO = make([]any, 0)
		}
		if SearchPercentage != -1 {
			CallJs("更新搜索进度", SearchPercentage)
		}
		if StatusText != "SunnyNetStatusText" {
			CallJs("更新状态文本", StatusText)
			StatusText = "SunnyNetStatusText"
		}
		Insert.Unlock()
	}
}
func SetStatusText(Text string) {
	Insert.Lock()
	StatusText = Text
	Insert.Unlock()
}
func UpdateResponseLength() {
	for {
		time.Sleep(1 * time.Second)
		am := HashMap.GetALLResponseLength()
		if len(am) < 1 {
			continue
		}
		Insert.Lock()
		CallJs("更新响应长度", am)
		Insert.Unlock()
	}
}
func UpdateIco(conn *MapHash.Request, _ContentType string) string {
	ContentType := strings.ToLower(_ContentType)
	Method := strings.ToUpper(conn.Method)
	if strings.Contains(ContentType, "image/") {
		return "img"
	}
	if strings.Contains(ContentType, "/javascript") {
		return "js"
	}
	if strings.Contains(ContentType, "/x-javascript") {
		return "js"
	}
	if strings.Contains(ContentType, "/css") {
		return "css"
	}
	if strings.Contains(ContentType, "/xml") {
		return "XML"
	}
	if strings.Contains(ContentType, "/json") {
		return "JSON"
	}
	if strings.Contains(ContentType, "/html") {
		return "HTML"
	}
	if strings.Contains(ContentType, "audio/") {
		return "audio"
	}
	if strings.Contains(ContentType, "video/") {
		return "video"
	}
	{
		URL := strings.ToLower(conn.URL)
		array := strings.Split(URL, "?")
		if len(array) > 1 {
			URL = array[0]
		}
		{
			if strings.HasSuffix(URL, ".woff2") ||
				strings.HasSuffix(URL, ".woff") ||
				strings.HasSuffix(URL, ".eot") ||
				strings.HasSuffix(URL, ".otf") ||
				strings.HasSuffix(URL, ".fon") ||
				strings.HasSuffix(URL, ".font") ||
				strings.HasSuffix(URL, ".ttc") ||
				strings.HasSuffix(URL, ".eotz") ||
				strings.HasSuffix(URL, ".dfont") ||
				strings.HasSuffix(URL, ".suit") ||
				strings.HasSuffix(URL, ".pfb") ||
				strings.HasSuffix(URL, ".ttf") {
				return "font"
			}
		}
		{
			if strings.HasSuffix(URL, ".swf") || strings.Contains(ContentType, "flash") {
				return "Flash"
			}
		}
	}
	if Method == "POST" || Method == "PUT" {
		return "POST"
	}
	if conn.Response.StateCode == 302 {
		return "302"
	}
	if conn.Response.StateCode == 401 {
		return "401"
	}
	if conn.Response.StateCode == 403 || conn.Response.StateCode == 404 || conn.Response.StateCode == 405 {
		return "stop"
	}
	return "generic"
}
func HttpCallback(Conn *SunnyNet.HttpConn) {
	SunnyNetMode := 0
	{
		if Conn.Type == public.HttpSendRequest {
			if Conn.Request != nil {
				if Conn.Request.Header != nil {
					SunnyNetMode, _ = strconv.Atoi(Conn.Request.Header.Get("SunnyNetMode"))
					Conn.Request.Header.Del("SunnyNetMode")
				}
			}
			HostsRulesUrl(Conn.Request.URL)
			u, b := ReplaceURL(Conn.Request.URL)
			if len(b) > 0 {
				Conn.Response = new(http.Response)
				Conn.Response.Body = io.NopCloser(bytes.NewBuffer(b))
				Conn.Response.Header = make(http.Header)
				Conn.Response.Header.Set("Server", "SunnyReplaceRules")
				Conn.Response.Header.Set("Accept-Ranges", "bytes")
				Conn.Response.Header.Set("Connection", "Close")
				Conn.Response.Header.Set("Content-Length", strconv.Itoa(len(b)))
				Conn.Response.ContentLength = int64(len(b))
				Conn.Response.StatusCode = 200
			} else {
				Conn.Request.URL = u
				ReplaceHeader(Conn.Request.Header)
				{
					_TmpLock.Lock()
					if DisableCache {
						delete(Conn.Request.Header, "If-None-Match")
						delete(Conn.Request.Header, "If-Modified-Since")
					}
					_TmpLock.Unlock()
				}
				Body := make([]byte, 0)
				if Conn.Request.Body != nil {
					Body, _ = io.ReadAll(Conn.Request.Body)
					_ = Conn.Request.Body.Close()
				}
				Body = ReplaceBody(Body)
				Conn.Request.Body = io.NopCloser(bytes.NewBuffer(Body))
			}
		} else if Conn.Type == public.HttpResponseOK {
			ReplaceHeader(Conn.Response.Header)
			{
				_TmpLock.Lock()
				if DisableCache {
					Conn.Response.Header["Cache-Control"] = []string{"no-cache, no-store, must-revalidate"}
					Conn.Response.Header["Pragma"] = []string{"no-cache"}
					Conn.Response.Header["Expires"] = []string{"0"}
				}
				_TmpLock.Unlock()
			}
			Body := make([]byte, 0)
			if Conn.Response.Body != nil {
				Body, _ = io.ReadAll(Conn.Response.Body)
				_ = Conn.Response.Body.Close()
			}
			Body = ReplaceBody(Body)
			Conn.Response.Body = io.NopCloser(bytes.NewBuffer(Body))
		}
		if !(GetWorkingState()) && Conn.Type == public.HttpSendRequest {
			if Conn.Type == public.HttpSendRequest {
				//执行发起请求脚本
				RunHTTPRequestScriptCode(Conn)
			} else if Conn.Type == public.HttpSendRequest {
				//执行请求响应脚本
				RunHTTPResponseScriptCode(Conn)
				h := HashMap.GetRequest(Conn.Theology)
				if h != nil {
					h.Conn = nil
				}
			} else {
				//执行请求错误脚本
				RunHTTPErrorScriptCode(Conn)
				h := HashMap.GetRequest(Conn.Theology)
				if h != nil {
					h.Conn = nil
				}
			}
			return
		}
	}
	if Conn.Type == public.HttpSendRequest {
		{
			if !RunHTTPRequestScriptCode(Conn) {
				return
			}
			h := HashMap.GetRequest(Conn.Theology)
			if h == nil {
				return
			}
			h.SendTime = time.Now().Format("15:04:05.000")
			if len(h.Response.Header) > 0 {
				h.RecTime = time.Now().Format("15:04:05.000")
				AddInsertList(&ListInfo{
					MessageId: -1,
					URL:       Conn.Request.URL.String(),
					HOST:      Conn.Request.URL.Host,
					ClientIP:  Conn.ClientIP,
					PID:       h.PID,
					Method:    Conn.Request.Method,
					Theology:  Conn.Theology,
					SendTime:  h.SendTime,
					RecTime:   h.RecTime,
					State:     "  -  ",
					Ico:       "上行",
					Break:     0,
				})
				time.Sleep(100 * time.Millisecond)
				Insert.Lock()
				isUpdateRequestInfo := currentlySelected == Conn.Theology
				Insert.Unlock()
				if isUpdateRequestInfo {
					CallJs("更新响应", &UpdateCurrentResponse{
						Theology:  Conn.Theology,
						Header:    h.Response.Header,
						Body:      h.Response.Body,
						StateText: http.StatusText(h.Response.StateCode),
						StateCode: h.Response.StateCode,
						Break:     false,
						Error:     false,
					})
				}
				Insert.Lock()
				ResponseType := ""
				if h.Response.Header != nil {
					_a := h.Response.Header["Content-Type"]
					if len(_a) > 0 {
						ResponseType = _a[0]
					} else {
						_a = h.Response.Header["content-type"]
						if len(_a) > 0 {
							ResponseType = _a[0]
						}
					}
					if ResponseType != "" {
						array := strings.Split(ResponseType+";", ";")
						if len(array) > 0 {
							ResponseType = array[0]
						}
					}
				}
				_tmp := &ListInfo{
					MessageId: -1,
					URL:       Conn.Request.URL.String(),
					HOST:      Conn.Request.URL.Host,
					ClientIP:  Conn.ClientIP,
					PID:       CommAnd.GetPidName(Conn.PID),
					Method:    Conn.Request.Method,
					Theology:  Conn.Theology,
					State:     strconv.Itoa(h.Response.StateCode),
					Len:       strconv.Itoa(len(h.Response.Body)),
					Type:      ResponseType,
					SendTime:  h.SendTime,
					RecTime:   h.RecTime,
					Ico:       UpdateIco(h, ResponseType),
					Break:     0,
				}
				UpdateData = append(UpdateData, _tmp)
				Insert.Unlock()
				return
			}
			Insert.Lock()
			IsBreak := uint8(0)
			if breakpoint == 1 || h.Break == 1 || SunnyNetMode == 1 {
				IsBreak = 1
			}
			Insert.Unlock()
			_tmp := &ListInfo{
				MessageId: -1,
				URL:       Conn.Request.URL.String(),
				HOST:      Conn.Request.URL.Host,
				ClientIP:  Conn.ClientIP,
				PID:       CommAnd.GetPidName(Conn.PID),
				Method:    Conn.Request.Method,
				Theology:  Conn.Theology,
				SendTime:  h.SendTime,
				State:     "  -  ",
				Ico:       "上行",
				Break:     IsBreak,
			}
			if len(h.Response.Body) > 0 || len(h.Response.Header) > 0 {
				ResponseType := "空白"
				if h.Response.Header != nil {
					_a := h.Response.Header["Content-Type"]
					if len(_a) > 0 {
						ResponseType = _a[0]
					} else {
						_a = h.Response.Header["content-type"]
						if len(_a) > 0 {
							ResponseType = _a[0]
						}
					}
					if ResponseType != "" {
						array := strings.Split(ResponseType+";", ";")
						if len(array) > 0 {
							ResponseType = array[0]
						}
					}
				}
				_tmp.Ico = ResponseType
				_tmp.Break = 0
			}
			AddInsertList(_tmp)
			if IsBreak == 1 {
				h.Wait.Add(1)
				h.Wait.Wait()
				_tmp = &ListInfo{
					MessageId: -1,
					URL:       Conn.Request.URL.String(),
					HOST:      Conn.Request.URL.Host,
					ClientIP:  Conn.ClientIP,
					PID:       CommAnd.GetPidName(Conn.PID),
					Method:    Conn.Request.Method,
					Theology:  Conn.Theology,
					SendTime:  h.SendTime,
					State:     "  -  ",
					Ico:       "上行",
					Break:     0,
				}
				AddInsertList(_tmp)
			}
			if SunnyNetMode == 2 {
				h.Break = 2
			}
		}
	} else if Conn.Type == public.HttpResponseOK {
		if Conn.Response.Header != nil {
			if Conn.Response.Body != nil {
				Body, _ := io.ReadAll(Conn.Response.Body)
				_ = Conn.Response.Body.Close()
				Encoding := ""
				ar := Conn.Response.Header["Content-Encoding"]
				if len(ar) > 0 {
					Encoding = ar[0]
				} else {
					ar = Conn.Response.Header["content-encoding"]
					if len(ar) > 0 {
						Encoding = ar[0]
					}
				}
				Encoding = strings.ToLower(Encoding)
				isOk := false
				if Encoding == "gzip" {
					gr, err := gzip.NewReader(io.NopCloser(bytes.NewBuffer(Body)))
					if err == nil {
						gr1, err1 := io.ReadAll(gr)
						if err1 == nil {
							Conn.Response.Body = io.NopCloser(bytes.NewBuffer(gr1))
							isOk = true
						}
					}
				} else if Encoding == "br" {
					br, err := io.ReadAll(brotli.NewReader(io.NopCloser(bytes.NewBuffer(Body))))
					if err == nil {
						Conn.Response.Body = io.NopCloser(bytes.NewBuffer(br))
						isOk = true
					}
				} else if Encoding == "deflate" {
					zr := flate.NewReader(io.NopCloser(bytes.NewBuffer(Body)))
					_ = zr.Close()
					bx, err := io.ReadAll(zr)
					if err == nil {
						Conn.Response.Body = io.NopCloser(bytes.NewBuffer(bx))
						isOk = true
					}
				}
				if isOk {
					delete(Conn.Response.Header, "content-encoding")
					delete(Conn.Response.Header, "Content-Encoding")
				} else {
					Conn.Response.Body = io.NopCloser(bytes.NewBuffer(Body))
				}
			}
			delete(Conn.Response.Header, "Transfer-Encoding")
		}
		Break := RunHTTPResponseScriptCode(Conn)
		Insert.Lock()
		isUpdateRequestInfo := currentlySelected == Conn.Theology
		Insert.Unlock()
		h := HashMap.GetRequest(Conn.Theology)
		if h == nil {
			return
		}
		IsBreak := uint8(0)
		if Break {
			h.Break = 2
		}
		h.RecTime = time.Now().Format("15:04:05.000")
		if breakpoint == 2 || h.Break == 2 {
			IsBreak = 2
		}
		if isUpdateRequestInfo {
			CallJs("更新响应", &UpdateCurrentResponse{
				Theology:  Conn.Theology,
				Header:    h.Response.Header,
				Body:      h.Response.Body,
				StateText: http.StatusText(h.Response.StateCode),
				StateCode: h.Response.StateCode,
				Break:     IsBreak == 2,
				Error:     false,
			})

		}
		Insert.Lock()
		ResponseType := ""
		if h.Response.Header != nil {
			_a := h.Response.Header["Content-Type"]
			if len(_a) > 0 {
				ResponseType = _a[0]
			} else {
				_a = h.Response.Header["content-type"]
				if len(_a) > 0 {
					ResponseType = _a[0]
				}
			}
			if ResponseType != "" {
				array := strings.Split(ResponseType+";", ";")
				if len(array) > 0 {
					ResponseType = array[0]
				}
			}
		}
		_tmp := &ListInfo{
			MessageId: -1,
			URL:       Conn.Request.URL.String(),
			HOST:      Conn.Request.URL.Host,
			ClientIP:  Conn.ClientIP,
			PID:       CommAnd.GetPidName(Conn.PID),
			Method:    Conn.Request.Method,
			Theology:  Conn.Theology,
			State:     strconv.Itoa(h.Response.StateCode),
			Len:       strconv.Itoa(len(h.Response.Body)),
			Type:      ResponseType,
			RecTime:   h.RecTime,
			Ico:       UpdateIco(h, ResponseType),
			Break:     IsBreak,
		}
		UpdateData = append(UpdateData, _tmp)
		Insert.Unlock()
		if IsBreak == 2 {
			h.Wait.Add(1)
			h.Wait.Wait()
			h.Break = 0
			_tmp = &ListInfo{
				MessageId: -1,
				URL:       Conn.Request.URL.String(),
				HOST:      Conn.Request.URL.Host,
				ClientIP:  Conn.ClientIP,
				PID:       CommAnd.GetPidName(Conn.PID),
				Method:    Conn.Request.Method,
				Theology:  Conn.Theology,
				State:     strconv.Itoa(h.Response.StateCode),
				Len:       strconv.Itoa(len(h.Response.Body)),
				Type:      ResponseType,
				RecTime:   h.RecTime,
				Ico:       UpdateIco(h, ResponseType),
				Break:     0,
			}
			UpdateData = append(UpdateData, _tmp)
		}
		h.Conn = nil
	} else if Conn.Type == public.HttpRequestFail {
		RunHTTPErrorScriptCode(Conn)
		h := HashMap.GetRequest(Conn.Theology)
		if h == nil {
			return
		}
		h.Break = 0
		h.RecTime = time.Now().Format("15:04:05.000")
		h.Response.Error = true
		h.Response.StateCode = -1
		h.Response.Body = []byte(Conn.GetError())
		Insert.Lock()
		isUpdateRequestInfo := currentlySelected == Conn.Theology
		Insert.Unlock()
		if isUpdateRequestInfo {
			CallJs("更新响应", &UpdateCurrentResponse{
				Theology:  Conn.Theology,
				Body:      h.Response.Body,
				StateText: "error",
				StateCode: h.Response.StateCode,
				Break:     false,
				Error:     true,
			})
		}
		Insert.Lock()
		_tmp := &ListInfo{
			MessageId: -1,
			URL:       Conn.Request.URL.String(),
			HOST:      Conn.Request.URL.Host,
			ClientIP:  Conn.ClientIP,
			PID:       CommAnd.GetPidName(Conn.PID),
			Method:    Conn.Request.Method,
			Theology:  Conn.Theology,
			State:     strconv.Itoa(h.Response.StateCode),
			Len:       strconv.Itoa(len(h.Response.Body)),
			Type:      "error",
			RecTime:   h.RecTime,
			Ico:       "error",
			Break:     0,
		}
		UpdateData = append(UpdateData, _tmp)
		Insert.Unlock()
		h.Conn = nil
	}
}
func WSCallback(Conn *SunnyNet.WsConn) {
	if Conn.Type == public.WebsocketDisconnect {
		time.Sleep(2 * time.Second)
	} else if Conn.Type == public.WebsocketUserSend || Conn.Type == public.WebsocketServerSend {
		Conn.SetMessageBody(ReplaceBody(Conn.GetMessageBody()))
	}
	Break := RunWebSocketScriptCode(Conn)
	if !(GetWorkingState()) {
		return
	}
	h := HashMap.GetRequest(Conn.Theology)
	if h == nil {
		return
	}
	if Conn.Type == public.WebsocketConnectionOK {
		h.WsConn = Conn
		h.Way = "Websocket"
		//Websocket连接成功
		Insert.Lock()
		_tmp := &ListInfo{
			MessageId: -1,
			URL:       Conn.Request.URL.String(),
			HOST:      Conn.Request.URL.Host,
			ClientIP:  Conn.ClientIP,
			PID:       CommAnd.GetPidName(Conn.Pid),
			Method:    "WebSocket",
			Theology:  Conn.Theology,
			Len:       "0/0",
			Type:      "WebSocket",
			RecTime:   h.RecTime,
			Ico:       "websocket_connect",
			Break:     0,
		}
		UpdateData = append(UpdateData, _tmp)
		Insert.Unlock()
		return
	}
	if Conn.Type == public.WebsocketUserSend || Conn.Type == public.WebsocketServerSend {
		h.RecTime = time.Now().Format("15:04:05.000")
		_tmp := &ListInfo{
			MessageId: -1,
			URL:       Conn.Request.URL.String(),
			HOST:      Conn.Request.URL.Host,
			ClientIP:  Conn.ClientIP,
			PID:       CommAnd.GetPidName(Conn.Pid),
			Method:    "WebSocket",
			Theology:  Conn.Theology,
			State:     "已连接",
			Len:       "0/0",
			Type:      "WebSocket",
			RecTime:   h.RecTime,
			Ico:       "websocket_connect",
			Break:     0,
		}
		AddInsertList(_tmp)
		if !Break {
			return
		}
		if Conn.Type == public.WebsocketUserSend {
			if h.Options.StopSend || h.Options.StopALL {
				return
			}
		} else {
			if h.Options.StopRec || h.Options.StopALL {
				return
			}
		}
		h.WsConn = Conn
		Body := Conn.GetMessageBody()
		if len(Body) < 1 {
			return
		}
		_WsType := Conn.GetMessageType()
		_Type := ""
		BodyHash := ""
		Ico := "上行"
		{
			if _WsType == 1 {
				_Type = "Text"
			} else if _WsType == 2 {
				_Type = "Binary"
			} else if _WsType == 8 {
				_Type = "Close"
			} else if _WsType == 9 {
				_Type = "Ping"
			} else if _WsType == 10 {
				_Type = "Pong"
			} else {
				_Type = "Invalid"
			}
			if len(Body) > 64 {
				BodyHash = fmt.Sprintf("% X", Body[:64]) + "..."
			} else {
				BodyHash = fmt.Sprintf("% X", Body)
			}
			if Conn.Type == public.WebsocketServerSend {
				Ico = "下行"
			}
		}
		_update := &MapHash.UpdateSocketData{
			Body: Body,
			Info: &MapHash.UpdateSocketList{
				Theology: Conn.Theology,
				Ico:      Ico,
				BodyHash: BodyHash,
				Length:   len(Body),
				Time:     time.Now().Format("15:04:05.000"),
				WsType:   _Type,
			},
		}
		HashMap.SetSocketData(Conn.Theology, _update, Conn.Type == public.WebsocketUserSend, len(Body))
		Insert.Lock()
		_update.Info.Index = len(h.SocketData)
		isUpdateRequestInfo := currentlySelected == Conn.Theology
		if isUpdateRequestInfo {
			SocketData = append(SocketData, _update.Info)
		}
		Insert.Unlock()
		//Websocket发送数据
		return
	}
	if Conn.Type == public.WebsocketDisconnect {
		h.WsConn = nil
		BodyHash := "已断开连接"
		_update := &MapHash.UpdateSocketData{
			Info: &MapHash.UpdateSocketList{
				Theology: Conn.Theology,
				Ico:      "websocket_close",
				BodyHash: BodyHash,
				Length:   0,
				Time:     "",
				Index:    -1,
			},
			Body: []byte("已断开连接"),
		}
		HashMap.SetSocketData(Conn.Theology, _update, false, 0)
		Insert.Lock()
		isUpdateRequestInfo := currentlySelected == Conn.Theology
		if isUpdateRequestInfo {
			SocketData = append(SocketData, _update.Info)
		}
		UpdateListICO = append(UpdateListICO, &UpdateICO{Theology: Conn.Theology, Ico: "websocket_close"})
		Insert.Unlock()
		//Websocket断开
		return
	}
}
func TcpCallback(Conn *SunnyNet.TcpConn) {
	h := HashMap.GetRequest(Conn.Theology)
	if Conn.Type == public.SunnyNetMsgTypeTCPClose {
		time.Sleep(2 * time.Second)
		if h != nil {
			h.TcpConn = nil
		}
	} else if Conn.Type == public.SunnyNetMsgTypeTCPClientSend || Conn.Type == public.SunnyNetMsgTypeTCPClientReceive {
		Conn.SetBody(ReplaceBody(Conn.GetBody()))
	} else if Conn.Type == public.SunnyNetMsgTypeTCPAboutToConnect {
		{
			h = HashMap.SetRequestTCP(Conn.Theology, Conn)
			h.URL = Conn.LocalAddr + "->" + Conn.RemoteAddr
			h.Method = string(Conn.GetBody())
			Conn.RemoteAddr = HostsRulesAddress(Conn.RemoteAddr)
			Conn.SetConnectionIP(Conn.RemoteAddr)
			return
		}
		//即将连接
	}

	if h == nil {
		return
	}
	Break := RunTcpScriptCode(Conn)
	if !(GetWorkingState()) {
		return
	}
	if Conn.Type == public.SunnyNetMsgTypeTCPConnectOK {
		//连接成功
		return
	}
	if Conn.Type == public.SunnyNetMsgTypeTCPClientSend || Conn.Type == public.SunnyNetMsgTypeTCPClientReceive {
		h.Way = h.Method
		if h.SendTime == "" {
			h.SendTime = time.Now().Format("15:04:05.000")
		}
		h.RecTime = time.Now().Format("15:04:05.000")
		_tmp := &ListInfo{
			MessageId: -1,
			URL:       h.URL,
			HOST:      Conn.RemoteAddr,
			ClientIP:  Conn.LocalAddr,
			PID:       CommAnd.GetPidName(Conn.Pid),
			Method:    h.Method,
			Theology:  Conn.Theology,
			State:     "已连接",
			Len:       "0/0",
			Type:      h.Method,
			SendTime:  h.SendTime,
			RecTime:   h.RecTime,
			Ico:       "websocket_connect",
			Break:     0,
		}
		AddInsertList(_tmp)
		h.PID = _tmp.PID
		h.ClientIP = Conn.LocalAddr
		//客户端发送\接收
		{
			if Conn.Type == public.SunnyNetMsgTypeTCPClientSend {
				if h.Options.StopSend || h.Options.StopALL {
					return
				}
			} else {
				if h.Options.StopRec || h.Options.StopALL {
					return
				}
			}
		}
		{
			if !Break {
				return
			}
			h.TcpConn = Conn
			Body := Conn.GetBody()
			if len(Body) < 1 {
				return
			}
			_Type := ""
			BodyHash := ""
			Ico := "上行"
			{
				if len(Body) > 64 {
					BodyHash = fmt.Sprintf("% X", Body[:64]) + "..."
				} else {
					BodyHash = fmt.Sprintf("% X", Body)
				}
				if Conn.Type == public.SunnyNetMsgTypeTCPClientReceive {
					Ico = "下行"
				}
			}
			_update := &MapHash.UpdateSocketData{
				Body: Body,
				Info: &MapHash.UpdateSocketList{
					Theology: Conn.Theology,
					Ico:      Ico,
					BodyHash: BodyHash,
					Length:   len(Body),
					Time:     time.Now().Format("15:04:05.000"),
					WsType:   _Type,
				},
			}
			HashMap.SetSocketData(Conn.Theology, _update, Conn.Type == public.SunnyNetMsgTypeTCPClientSend, len(Body))
			Insert.Lock()
			_update.Info.Index = len(h.SocketData)
			isUpdateRequestInfo := currentlySelected == Conn.Theology
			if isUpdateRequestInfo {
				SocketData = append(SocketData, _update.Info)
			}
			Insert.Unlock()
			//Websocket发送数据
			return
		}
		return
	}
	if Conn.Type == public.SunnyNetMsgTypeTCPClose {
		h.TcpConn = nil
		BodyHash := "已断开连接"
		_update := &MapHash.UpdateSocketData{
			Info: &MapHash.UpdateSocketList{
				Theology: Conn.Theology,
				Ico:      "websocket_close",
				BodyHash: BodyHash,
				Length:   0,
				Time:     "",
				Index:    -1,
			},
			Body: []byte("已断开连接"),
		}
		HashMap.SetSocketData(Conn.Theology, _update, false, 0)
		Insert.Lock()
		isUpdateRequestInfo := currentlySelected == Conn.Theology
		if isUpdateRequestInfo {
			SocketData = append(SocketData, _update.Info)
		}
		UpdateListICO = append(UpdateListICO, &UpdateICO{Theology: Conn.Theology, Ico: "websocket_close"})
		Insert.Unlock()
		//关闭
		return
	}

}
func UdpCallback(Conn *SunnyNet.UDPConn) {
	_TmpLock.Lock()
	if DisableUDP {
		_TmpLock.Unlock()
		Conn.Data = make([]byte, 0)
		return
	}
	_TmpLock.Unlock()
	if public.SunnyNetUDPTypeClosed == Conn.Type {
		time.Sleep(2 * time.Second)
	} else if public.SunnyNetUDPTypeSend == Conn.Type || public.SunnyNetUDPTypeReceive == Conn.Type {
		Conn.Data = ReplaceBody(Conn.Data)
	}
	Break := RunUdpScriptCode(Conn)
	if !(GetWorkingState()) {
		return
	}
	//在 Windows 捕获UDP需要加载驱动,并且设置进程名
	//其他情况需要设置Socket5代理,才能捕获到UDP
	//捕获到数据可以修改,修改空数据,取消发送/接收
	Theology := int(Conn.Theology)
	h := HashMap.GetRequest(Theology)
	if public.SunnyNetUDPTypeSend == Conn.Type || public.SunnyNetUDPTypeReceive == Conn.Type {
		{
			{
				if h == nil {
					h = HashMap.SetRequestUDP(Theology, Conn)
					h.URL = Conn.LocalAddress + "->" + Conn.RemoteAddress
					h.Method = "UDP"
				}
				if h.UdpConn == nil {
					h = HashMap.SetRequestUDP(Theology, Conn)
					h.URL = Conn.LocalAddress + "->" + Conn.RemoteAddress
					h.Method = "UDP"
				}
			}
			{
				if !Break {
					return
				}
				if Conn.Type == public.SunnyNetUDPTypeSend {
					if h.Options.StopSend || h.Options.StopALL {
						return
					}
				} else {
					if h.Options.StopRec || h.Options.StopALL {
						return
					}
				}
			}
			h.Way = h.Method
		}
		if h.SendTime == "" {
			h.SendTime = time.Now().Format("15:04:05.000")
		}
		h.RecTime = time.Now().Format("15:04:05.000")
		_tmp := &ListInfo{
			MessageId: -1,
			URL:       h.URL,
			HOST:      Conn.RemoteAddress,
			ClientIP:  Conn.LocalAddress,
			PID:       CommAnd.GetPidName(Conn.Pid),
			Method:    h.Method,
			Theology:  Theology,
			State:     "已连接",
			Len:       "0/0",
			Type:      h.Method,
			SendTime:  h.SendTime,
			RecTime:   h.RecTime,
			Ico:       "websocket_connect",
			Break:     0,
		}
		h.PID = _tmp.PID
		h.ClientIP = Conn.LocalAddress
		AddInsertList(_tmp)
		Body := Conn.Data
		if len(Body) < 1 {
			return
		}
		BodyHash := ""
		Ico := "上行"
		{
			if len(Body) > 64 {
				BodyHash = fmt.Sprintf("% X", Body[:64]) + "..."
			} else {
				BodyHash = fmt.Sprintf("% X", Body)
			}
			if Conn.Type == public.SunnyNetUDPTypeReceive {
				Ico = "下行"
			}
		}
		_update := &MapHash.UpdateSocketData{
			Body: Body,
			Info: &MapHash.UpdateSocketList{
				Theology: Theology,
				Ico:      Ico,
				BodyHash: BodyHash,
				Length:   len(Body),
				Time:     time.Now().Format("15:04:05.000"),
			},
		}
		HashMap.SetSocketData(Theology, _update, Conn.Type == public.SunnyNetUDPTypeSend, len(Body))
		Insert.Lock()
		_update.Info.Index = len(h.SocketData)
		isUpdateRequestInfo := currentlySelected == Theology
		if isUpdateRequestInfo {
			SocketData = append(SocketData, _update.Info)
		}
		Insert.Unlock()
		//Websocket发送数据
		return
	}
	if public.SunnyNetUDPTypeClosed == Conn.Type {
		if h == nil {
			return
		}
		h.UdpConn = nil
		BodyHash := "已断开连接"
		_update := &MapHash.UpdateSocketData{
			Info: &MapHash.UpdateSocketList{
				Theology: Theology,
				Ico:      "websocket_close",
				BodyHash: BodyHash,
				Length:   0,
				Time:     "",
				Index:    -1,
			},
			Body: []byte("已断开连接"),
		}
		HashMap.SetSocketData(Theology, _update, false, 0)
		Insert.Lock()
		isUpdateRequestInfo := currentlySelected == Theology
		if isUpdateRequestInfo {
			SocketData = append(SocketData, _update.Info)
		}
		UpdateListICO = append(UpdateListICO, &UpdateICO{Theology: Theology, Ico: "websocket_close"})
		Insert.Unlock()
		//关闭
		return
	}
}
