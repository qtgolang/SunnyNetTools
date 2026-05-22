package Session

import (
	"bytes"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/qtgolang/SunnyNet/SunnyNet"
	"github.com/qtgolang/SunnyNet/src/SunnyProxy"
	"github.com/qtgolang/SunnyNet/src/public"
)

func GetTCPSession(Theology int) *TCPSession {
	O, OK := Session.Load(Theology)
	if OK == false {
		return nil
	}
	O1, OK1 := O.(*TCPSession)
	if OK1 == false {
		return nil
	}
	return O1
}
func (s *TCPSession) RangeStream(f func(AppStream) bool) {
	for _, key := range s.StreamKeys {
		if !f(s.Stream[key]) {
			return
		}
	}
}

func (s *TCPSession) AddStream(obj *Stream, _ ...int) {
	_, exists := s.Stream[obj.MessageId]
	if !exists {
		s.StreamKeys = append(s.StreamKeys, obj.MessageId)
	}
	s.Stream[obj.MessageId] = obj
}
func (s *Stream) ToSearchAny() map[string]interface{} {
	_t := "接收"
	if s.IsSend {
		_t = "发送"
	}
	if s.IsActiveSend {
		_t = "手动" + _t
	}
	if s.IsClose {
		_t = "连接关闭"
	}
	if bytes.Equal(s.Body, ConnectDisconnect) {
		_t = string(ConnectDisconnect)
	}
	testData := map[string]interface{}{
		"数据": s.Body,
		"类型": _t,
		"长度": float64(len(s.Body)),
		"时间": s.Time,
	}
	return testData
}

// MatchFilters 根据匹配规则显示是否匹配成功，如果无规则，则默认返回 true
func (s *Stream) MatchFilters(filter *Filter) bool {
	if filter == nil {
		return true
	}
	return matchFilters(s.ToSearchAny(), filter)
}

// ListIsMatch  此数据流
// 重新应用过滤器是否需要在主列表显示
// 继承插入时的过滤状态,因为在插入已经匹配其他项，所有现在只需要匹配该数据流就可以了
func (s *Stream) ListIsMatch() bool {
	if s.ListFilter == nil {
		return true
	}
	testData := map[string]interface{}{
		"全部数据": s.Body,
	}
	return matchFilters(testData, s.ListFilter)
}
func (s *Stream) ToUpdateStream(Theology int, filter *Filter) UpdateSocketStream {
	_stream := UpdateSocketStream{
		Theology:      Theology,
		MessageId:     s.MessageId,
		IsSend:        s.IsSend,
		Length:        len(s.Body),
		Time:          s.Time,
		Filter:        s.MatchFilters(filter),
		IsClose:       s.IsClose,
		IsActiveSend:  s.IsActiveSend,
		WebsocketType: s.GetWebsocketType(),
	}
	if s.IsClose {
		_stream.Body = string(ConnectDisconnect)
	} else {
		_stream.Body = getHexWithSpaces(s.Body)
	}
	return _stream
}

type TCPSession struct {
	Theology      int
	Method        string
	Ico           string
	Note          string //添加的注释
	ClientIP      string
	ProcessName   string
	RemoteAddress string //连接地址 IP文本
	Host          string //域名信息，没有域名的情况等于 RemoteAddress
	Time          string
	AnewInsert    bool             `json:"-"` //是否需要重新插入到列表
	RecLength     int              //接收数据长度
	SenLength     int              //发送数据长度
	Stream        map[int]*Stream  `json:"-"`
	StreamKeys    []int            `json:"-"`
	Conn          SunnyNet.ConnTCP `json:"-"`
	Disconnect    bool             `json:"-"`
	ListFilter    *Filter          `json:"-"` //主列表过滤器
	filter        *Filter
	UserName      string //身份验证账号
	sync.Mutex
}

func (s *TCPSession) GetNote() string {
	return s.Note
}

func (s *TCPSession) GetListFilter() *Filter {
	return s.ListFilter
}
func (s *TCPSession) SetListFilter(val *Filter) {
	s.ListFilter = val.Clone()
}
func (s *TCPSession) GetStreamFilter() *Filter {
	return s.filter.Clone()
}

func (s *TCPSession) SetStreamFilter(val *Filter) {
	s.filter = val.Clone()
}
func (s *TCPSession) GetStream(Theology int) AppStream {
	aa := s.Stream[Theology]
	return aa
}
func (s *TCPSession) ClearAllSession() []int {
	var array []int
	array = s.StreamKeys
	s.StreamKeys = make([]int, 0)
	s.Stream = make(map[int]*Stream, 0)
	s.RecLength = 0
	s.SenLength = 0
	return array
}
func (s *TCPSession) GetSenLength() int {
	return s.SenLength
}
func (s *TCPSession) GetRecLength() int {
	return s.RecLength
}
func (s *TCPSession) GetMethod() string {
	return s.Method
}
func (s *TCPSession) IsDisconnect() bool {
	return s.Disconnect
}

/*
ResendRequest 重新发送此请求N次
*/
func (s *TCPSession) ResendRequest(count, _, port int, OutRouterIP string) {
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
	proxyAddress, _ := SunnyProxy.ParseProxy(fmt.Sprintf("socks://%s:%d", "127.0.0.1", port), 30*1000)
	if count == 1 {
		go s.resendRequest(proxyAddress, outRouterIP)
		return
	}
	sem := make(chan struct{}, 10) //只允许10并发
	for i := 0; i < count; i++ {
		go func(id int) {
			sem <- struct{}{}
			defer func() {
				<-sem
			}()
			s.resendRequest(proxyAddress, outRouterIP)
		}(i)
	}
}
func (s *TCPSession) resendRequest(proxy *SunnyProxy.Proxy, OutRouterIP *net.TCPAddr) bool {
	wb, e := proxy.Dial("tcp", s.Host, OutRouterIP)
	if e != nil {
		return false
	}
	defer wb.Close()
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
			if _, er := wb.Write(stream.GetBody()); er != nil {
				return false
			}
		}
		return true
	})
	return true
}

// IsWait 当前是否正在断点拦截,仅HTTP请求有效， tcp/dup 始终返回false
func (s *TCPSession) IsWait() bool {
	return false
}
func (s *TCPSession) IsHTTP() bool      { return false }
func (s *TCPSession) IsWebsocket() bool { return false }
func (s *TCPSession) IsTCP() bool       { return true }
func (s *TCPSession) IsUDP() bool       { return false }
func (s *TCPSession) DeleteMessageIdArray(MessageIdArray []int) []int {
	var _tmpMap = make(map[int]bool, 0)
	for _, v := range MessageIdArray {
		_tmpMap[v] = true
		m := s.Stream[v]
		if m.IsSend {
			s.SenLength -= len(m.Body)
		} else {
			s.RecLength -= len(m.Body)
		}
		delete(s.Stream, v)
	}
	var newKeys []int
	var array []int
	for _, v := range s.StreamKeys {
		if !_tmpMap[v] {
			newKeys = append(newKeys, v)
		} else {
			array = append(array, v)
		}
	}
	_tmpMap = make(map[int]bool, 0)
	s.StreamKeys = newKeys
	return array
}
func (s *TCPSession) ActiveSend(isSendServer bool, _ int, data []byte) string {
	if isSendServer {
		if !s.Conn.SendToServer(data) {
			return "向服务器发送数据失败"
		}
		s.SenLength += len(data)
	} else {
		if !s.Conn.SendToClient(data) {
			return "向客户端发送数据失败"
		}
		s.RecLength += len(data)
	}
	return ""
}
func (s *TCPSession) GetTheology() int {
	return s.Theology
}
func (s *TCPSession) SetAnewInsert(anewInsert bool) {
	s.AnewInsert = anewInsert
}

func (s *TCPSession) ListMatch() bool {
	if s.ListFilter == nil {
		return true
	}
	return matchFilters(s.ToSearchAny(), s.ListFilter)
}
func (s *TCPSession) SetNote(note string) {
	s.Note = note
}
func (s *TCPSession) ToSearchAny() map[string]interface{} {
	testData := map[string]any{
		"方式":     s.Method,
		"响应类型":   s.Method,
		"进程":     s.ProcessName,
		"状态":     Or(s.Disconnect, "已连接", "已断开"),
		"HOST":   s.Host,
		"主机名":    s.Host,
		"请求时间":   s.Time,
		"响应时间":   "",
		"Path":   "",
		"参数":     "",
		"注释":     s.Note,
		"来源地址":   s.ClientIP,
		"响应IP":   s.RemoteAddress,
		"身份验证账号": s.UserName,
		"请求地址":   s.ClientIP + " -> " + s.Host,
	}
	testData["全部数据"] = func(f *Filter) bool {
		ok := false
		f.ColId = "数据"
		s.RangeStream(func(stream AppStream) bool {
			ok = stream.MatchFilters(f)
			return !ok
		})
		return ok
	}
	testData["响应长度"] = func(f *Filter) bool {
		ok := false
		s.RangeStream(func(stream AppStream) bool {
			ok = stream.MatchFilters(f)
			return !ok
		})
		return ok
	}
	return testData
}
