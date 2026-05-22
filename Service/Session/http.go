package Session

import (
	"bytes"
	"fmt"
	"net"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/qtgolang/SunnyNet/SunnyNet"
	"github.com/qtgolang/SunnyNet/src/SunnyProxy"
	"github.com/qtgolang/SunnyNet/src/crypto/tls"
	"github.com/qtgolang/SunnyNet/src/http"
	"github.com/qtgolang/SunnyNet/src/httpClient"
	"github.com/qtgolang/SunnyNet/src/loop"
	"github.com/qtgolang/SunnyNet/src/public"
	"github.com/qtgolang/SunnyNet/src/websocket"
)

func (h HttpSessionRequest) FirstLine() string {
	return h.Method + " " + h.Url + " " + h.Proto
}

func (h HttpSessionRequest) GetHeader() map[string][]string {
	return h.Header
}

func (h HttpSessionRequest) GetBody() []byte {
	return h.Body
}
func (h HttpSessionResponse) FirstLine() string {
	return h.Proto + " " + h.Code + " " + h.State
}

func (h HttpSessionResponse) GetHeader() map[string][]string {
	return h.Header
}

func (h HttpSessionResponse) GetBody() []byte {
	return h.Body
}

type HttpMessage interface {
	FirstLine() string
	GetHeader() map[string][]string
	GetBody() []byte
}
type HttpSessionRequest struct {
	Method       string
	Url          string
	ClientIP     string
	Time         string
	ProcessName  string
	Proto        string
	Header       http.Header
	Body         []byte `json:"-"`
	IsUpdateBody bool   `json:"-"`
	BodyLength   int
	IsMaxLength  bool
}

type HttpSessionResponse struct {
	Time         string
	Code         string //响应状态码
	State        string //响应状态文本
	Type         string
	Length       int
	Proto        string
	Header       http.Header
	Body         []byte
	ResponseType string
	ServerIP     string `json:"-"` //是否本地响应
	IsUpdateBody bool   `json:"-"`
	BodyLength   int
	IsMaxLength  bool
}
type HttpSession struct {
	Theology            int
	Ico                 string
	Request             HttpSessionRequest
	Response            HttpSessionResponse
	Error               string
	State               int                      //请求状态
	RecLength           int                      //接收数据长度
	SenLength           int                      //发送数据长度
	Note                string                   //添加的注释
	WebsocketStream     map[int]*WebsocketStream `json:"-"`
	WebsocketStreamKeys []int                    `json:"-"`
	IsWebsocketRequest  bool                     `json:"-"` //是否未ws请求
	WebsocketDisconnect bool                     `json:"-"` //ws请求是否已断开
	AnewInsert          bool                     `json:"-"` //是否需要重新插入到列表
	ListFilter          *Filter                  `json:"-"`
	NextBreakMode       uint32                   `json:"-"` //下一个拦截状态
	Wg                  WaitGroup                `json:"-"`
	WebsocketConn       SunnyNet.ConnWebSocket   `json:"-"`
	filter              *Filter
	UserName            string //身份验证账号
	sync.Mutex
}

func (s *HttpSession) GetNote() string {
	return s.Note
}

// IsWait 当前是否正在断点拦截,仅HTTP请求有效， tcp/dup 始终返回false
func (s *HttpSession) IsWait() bool {
	return s.Wg.IsWait()
}
func (s *HttpSession) GetListFilter() *Filter {
	return s.ListFilter
}
func (s *HttpSession) SetListFilter(val *Filter) {
	s.ListFilter = val.Clone()
}
func (s *HttpSession) RangeStream(f func(AppStream) bool) {
	for _, key := range s.WebsocketStreamKeys {
		if !f(s.WebsocketStream[key]) {
			return
		}
	}
}
func (s *HttpSession) IsDisconnect() bool {
	return s.WebsocketDisconnect
}

/*
ResendRequest 重新发送此请求N次
*/
func (s *HttpSession) ResendRequest(count, BreakMode, port int, OutRouterIP string) {
	proxyAddress, _ := SunnyProxy.ParseProxy(fmt.Sprintf("socks://%s:%d", "127.0.0.1", port), 30*1000)
	if count == 1 {
		go s.resendRequest(BreakMode, proxyAddress, OutRouterIP)
		return
	}
	sem := make(chan struct{}, 10) //只允许10并发
	for i := 0; i < count; i++ {
		go func(id int) {
			sem <- struct{}{}
			defer func() {
				<-sem
			}()
			s.resendRequest(0, proxyAddress, OutRouterIP) //批量重新发送，不支持设置断点拦截
		}(i)
	}
}

func (s *HttpSession) resendRequest(BreakMode int, Proxy *SunnyProxy.Proxy, OutRouterIP string) {
	if s.IsWebsocket() {
		var dialer websocket.Dialer
		mUrl := strings.ToLower(s.Request.Url)
		if strings.HasPrefix(mUrl, "https") || strings.HasPrefix(mUrl, "wss") {
			t := &tls.Config{InsecureSkipVerify: true}
			dialer = websocket.Dialer{TLSClientConfig: t}
		} else {
			dialer = websocket.Dialer{}
		}
		var outRouterIP *net.TCPAddr
		_, ip := public.IsLocalIP(OutRouterIP)
		if ip != nil {
			if ip.To4() != nil {
				localAddr, err := net.ResolveTCPAddr("tcp", OutRouterIP+":0")
				if err == nil {
					outRouterIP = localAddr
				}
			} else {
				localAddr, err := net.ResolveTCPAddr("tcp", "["+OutRouterIP+"]:0")
				if err == nil {
					outRouterIP = localAddr
				}
			}
		}
		Header := s.Request.Header.Clone()
		Header.Set(public.HTTPClientTags, "true")
		if BreakMode == 1 {
			// 发送断点拦截
			Header.Set(public.HTTPClientTags+"_BreakMode", "Send")
		} else if BreakMode == 2 {
			// 接收断点拦截
			Header.Set(public.HTTPClientTags+"_BreakMode", "Rec")
		}
		Header.Set(public.HTTPClientTags+"_GuaranteeDisplay", "true")
		wb, resq, _, err := dialer.Dial(s.Request.Url, Header, Proxy, outRouterIP)
		if err != nil || resq == nil {
			return
		}
		_sleep := int64(0)
		s.RangeStream(func(stream AppStream) bool {
			if stream.GetIsSend() {
				_time := parseFormattedTimeToTimestamp(stream.GetMessageTime())
				if _sleep > 0 {
					p := _time - _sleep
					if p > 5000 {
						p = 5000
					}
					time.Sleep(time.Duration(p) * time.Millisecond)
				}
				_sleep = _time
				if er := wb.WriteMessage(stream.GetWebsocketType(), stream.GetBody()); er != nil {
					return false
				}
			}
			return true
		})
		_ = wb.Close()
		return
	}
	req, err := http.NewRequest(s.Request.Method, s.Request.Url, bytes.NewBuffer(s.Request.Body))
	if err != nil {
		return
	}
	req.Header = s.Request.Header.Clone()
	req.Header.Set(public.HTTPClientTags, "true")
	if BreakMode == 1 {
		// 发送断点拦截
		req.Header.Set(public.HTTPClientTags+"_BreakMode", "Send")
	} else if BreakMode == 2 {
		// 接收断点拦截
		req.Header.Set(public.HTTPClientTags+"_BreakMode", "Rec")
	}
	req.Header.Set(public.HTTPClientTags+"_GuaranteeDisplay", "true")
	if OutRouterIP != "" {
		req.SetContext(public.OutRouterIPKey, OutRouterIP)
	}
	_connection := func(conn net.Conn) {
		l, _ := extractPorts(conn)
		loop.AddLoopFilter(l)
	}
	_close := func(conn net.Conn) {
		l, _ := extractPorts(conn)
		loop.UnLoopFilter(l)
	}
	op := httpClient.Options{
		RequestProxy:  Proxy,
		CheckRedirect: false,
		TLSConfig:     nil,
		OutTime:       (30 * 1000) * time.Millisecond,
		GetTLSValues:  public.GetTLSValues,
		MConn:         nil,
		Event:         httpClient.Event{Connection: _connection, Close: _close},
	}
	r := httpClient.DoOptions(req, op)
	if r.Close != nil {
		r.Close()
	}
}
func extractPorts(conn net.Conn) (local uint16, remote uint16) {
	if conn == nil { //连接为空直接报错
		return 0, 0
	}
	tcpLocal, ok1 := conn.LocalAddr().(*net.TCPAddr)   //本地TCP地址
	tcpRemote, ok2 := conn.RemoteAddr().(*net.TCPAddr) //对端TCP地址
	if !ok1 || !ok2 {                                  //非TCP连接
		return 0, 0
	}
	return uint16(tcpLocal.Port), uint16(tcpRemote.Port)
}
func (s *HttpSession) GetStreamFilter() *Filter {
	return s.filter.Clone()
}
func (s *HttpSession) SetStreamFilter(val *Filter) {
	s.filter = val.Clone()
}
func NewHttpSession() *HttpSession {
	return &HttpSession{}
}
func GetHttpSession(Theology int) *HttpSession {
	O, OK := Session.Load(Theology)
	if OK == false {
		return nil
	}
	O1, OK1 := O.(*HttpSession)
	if OK1 == false {
		return nil
	}
	return O1
}

func GetAppSession(theology int) AppSession {
	if value, ok := Session.Load(theology); ok {
		if obj, ok1 := value.(AppSession); ok1 {
			return obj
		}
	}
	return nil
}
func (s *HttpSession) GetStream(Theology int) AppStream {
	return s.WebsocketStream[Theology]
}
func (s *HttpSession) ClearAllSession() []int {
	var array []int
	array = s.WebsocketStreamKeys
	s.WebsocketStreamKeys = make([]int, 0)
	s.WebsocketStream = make(map[int]*WebsocketStream, 0)
	s.RecLength = 0
	s.SenLength = 0
	return array
}

func (s *HttpSession) GetSenLength() int {
	return s.SenLength
}
func (s *HttpSession) GetRecLength() int {
	return s.RecLength
}
func (s *HttpSession) GetMethod() string {
	if s.IsWebsocketRequest {
		return "Websocket"
	}
	return "HTTP"
}
func (s *HttpSession) ActiveSend(isSendServer bool, wsType int, data []byte) string {
	if isSendServer {
		if !s.WebsocketConn.SendToServer(wsType, data) {
			return "向服务器发送数据失败"
		}
		s.SenLength += len(data)
	} else {
		if !s.WebsocketConn.SendToClient(wsType, data) {
			return "向客户端发送数据失败"
		}
		s.RecLength += len(data)
	}
	return ""
}
func (s *HttpSession) IsHTTP() bool      { return true }
func (s *HttpSession) IsWebsocket() bool { return s.IsWebsocketRequest }
func (s *HttpSession) IsTCP() bool       { return false }
func (s *HttpSession) IsUDP() bool       { return false }
func (s *HttpSession) DeleteMessageIdArray(MessageIdArray []int) []int {
	var _tmpMap = make(map[int]bool)
	for _, v := range MessageIdArray {
		_tmpMap[v] = true
		m := s.WebsocketStream[v]
		if m.IsSend {
			s.SenLength -= len(m.Body)
		} else {
			s.RecLength -= len(m.Body)
		}
		delete(s.WebsocketStream, v)
	}
	var newKeys []int
	var array []int
	for _, v := range s.WebsocketStreamKeys {
		if !_tmpMap[v] {
			newKeys = append(newKeys, v)
		} else {
			array = append(array, v)
		}
	}
	_tmpMap = make(map[int]bool)
	s.WebsocketStreamKeys = newKeys
	return array
}
func (s *HttpSession) GetTheology() int {
	return s.Theology
}
func (s *HttpSession) SetAnewInsert(anewInsert bool) {
	s.AnewInsert = anewInsert
}

func (s *HttpSession) ListMatch() bool {
	if s.ListFilter == nil {
		return true
	}
	return matchFilters(s.ToSearchAny(), s.ListFilter)
}
func (s *HttpSession) SetNote(note string) {
	s.Note = note
}
func (s *HttpSession) ToSearchAny() map[string]interface{} {
	_u, _ := url.Parse(s.Request.Url)
	_Host := ""
	_Path := ""
	_RawQuery := ""
	_ClientIP := ""
	if _u != nil {
		_Host = _u.Host
		_Path = _u.Path
		_RawQuery = _u.RawQuery
		_ClientIP = s.Request.ClientIP + " -> " + _u.Host
	}
	testData := map[string]interface{}{
		"方式":     Or(s.IsWebsocketRequest, "Websocket", Or(s.State != public.HttpRequestFail, s.Request.Method, "错误")),
		"响应类型":   Or(s.IsWebsocketRequest, "Websocket", s.Response.ResponseType),
		"进程":     s.Request.ProcessName,
		"HOST":   _Host,
		"主机名":    _Host,
		"Path":   _Path,
		"参数":     _RawQuery,
		"请求地址":   _ClientIP,
		"请求时间":   s.Request.Time,
		"响应时间":   s.Response.Time,
		"注释":     s.Note,
		"状态码":    s.Response.Code,
		"来源地址":   s.Request.ClientIP,
		"响应IP":   s.Response.ServerIP,
		"身份验证账号": s.UserName,
	}
	if s.State == public.HttpRequestFail {
		testData["状态"] = "错误"
	} else {
		if s.IsWebsocketRequest {
			if s.WebsocketDisconnect {
				testData["状态"] = "已断开"
			} else {
				testData["状态"] = "已连接"
			}
		} else if s.State == public.HttpResponseOK {
			testData["状态"] = s.Response.Code
		} else {
			testData["状态"] = "  -  "
		}
	}
	if !s.IsWebsocketRequest {
		testData["响应长度"] = float64(s.Response.Length)
	} else {
		testData["响应长度"] = func(f *Filter) bool {
			var ok bool
			fs := f.Clone()
			s.RangeStream(func(stream AppStream) bool {
				ok = stream.MatchFilters(fs)
				return !ok
			})
			return ok
		}
	}
	testData["全部数据"] = func(f *Filter) bool {
		_data := fmt.Sprintf("%s %s %s\r\n", s.Request.Method, s.Request.Url, s.Request.Proto)
		for key, v := range s.Request.Header {
			for _, val := range v {
				_data += key + ": " + val + "\r\n"
			}
		}
		_data += string(s.Request.Body)
		Data := map[string]interface{}{"全部数据": _data}
		if f.Type == "notContains" {
			if !MatchFilters(Data, f.Clone()) {
				return false
			}
		} else {
			if MatchFilters(Data, f.Clone()) {
				return true
			}
		}

		_data = fmt.Sprintf("%s %s %s\r\n", s.Response.Proto, s.Response.Code, s.Response.State)
		for key, v := range s.Response.Header {
			for _, val := range v {
				_data += key + ": " + val + "\r\n"
			}
		}

		_data += string(s.Response.Body)
		Data = map[string]interface{}{"全部数据": _data}
		if f.Type == "notContains" {
			//返回true 表示不包含
			if !MatchFilters(Data, f.Clone()) {
				return false
			}
		} else {
			//返回 true 表示 包含
			if MatchFilters(Data, f.Clone()) {
				return true
			}
		}
		fs := f.Clone()
		fs.ColId = "数据"
		var ok bool
		if fs.Type == "notContains" {
			ok = true
			s.RangeStream(func(stream AppStream) bool {
				ok = stream.MatchFilters(fs)
				return ok
			})
		} else {
			s.RangeStream(func(stream AppStream) bool {
				ok = stream.MatchFilters(fs)
				return !ok
			})
		}
		return ok
	}
	return testData
}
