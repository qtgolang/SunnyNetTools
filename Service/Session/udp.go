package Session

import (
	"sync"

	"github.com/qtgolang/SunnyNet/SunnyNet"
)

type UDPSession struct {
	Theology      int
	Method        string
	Ico           string
	Note          string //添加的注释
	ClientIP      string
	ProcessName   string
	RemoteAddress string
	Time          string
	AnewInsert    bool             `json:"-"` //是否需要重新插入到列表
	RecLength     int              //接收数据长度
	SenLength     int              //发送数据长度
	Stream        map[int]*Stream  `json:"-"`
	StreamKeys    []int            `json:"-"`
	Conn          SunnyNet.ConnUDP `json:"-"`
	Disconnect    bool             `json:"-"`
	ListFilter    *Filter          `json:"-"`
	filter        *Filter
	UserName      string //身份验证账号 UDP无法获取始终为空字符串
	sync.Mutex
}

func (s *UDPSession) GetNote() string {
	return s.Note
}

func GetUDPSession(Theology int) *UDPSession {
	O, OK := Session.Load(Theology)
	if OK == false {
		return nil
	}
	O1, OK1 := O.(*UDPSession)
	if OK1 == false {
		return nil
	}
	return O1
}
func (s *UDPSession) IsDisconnect() bool {
	return s.Disconnect
}

/*
ResendRequest 重新发送此请求N次
*/
func (s *UDPSession) ResendRequest(_, _, _ int, _ string) {
	//UDP不支持重新发送
	panic("implement me")
}
func (s *UDPSession) AddStream(obj *Stream, _ ...int) {
	_, exists := s.Stream[obj.MessageId]
	if !exists {
		s.StreamKeys = append(s.StreamKeys, obj.MessageId)
	}
	s.Stream[obj.MessageId] = obj
}
func (s *UDPSession) RangeStream(f func(AppStream) bool) {
	for _, key := range s.StreamKeys {
		if !f(s.Stream[key]) {
			return
		}
	}
}

// IsWait 当前是否正在断点拦截,仅HTTP请求有效， tcp/dup 始终返回false
func (s *UDPSession) IsWait() bool {
	return false
}
func (s *UDPSession) GetStreamFilter() *Filter {
	return s.filter.Clone()
}
func (s *UDPSession) SetStreamFilter(val *Filter) {
	s.filter = val.Clone()
}
func (s *UDPSession) GetListFilter() *Filter {
	return s.ListFilter
}
func (s *UDPSession) SetListFilter(val *Filter) {
	s.ListFilter = val.Clone()
}
func (s *UDPSession) GetStream(Theology int) AppStream {
	return s.Stream[Theology]
}
func (s *UDPSession) ClearAllSession() []int {
	var array []int
	array = s.StreamKeys
	s.StreamKeys = make([]int, 0)
	s.Stream = make(map[int]*Stream, 0)
	s.RecLength = 0
	s.SenLength = 0
	return array
}
func (s *UDPSession) GetSenLength() int {
	return s.SenLength
}
func (s *UDPSession) GetRecLength() int {
	return s.RecLength
}
func (s *UDPSession) GetMethod() string {
	return "UDP"
}
func (s *UDPSession) IsHTTP() bool      { return false }
func (s *UDPSession) IsWebsocket() bool { return false }
func (s *UDPSession) IsTCP() bool       { return false }
func (s *UDPSession) IsUDP() bool       { return true }
func (s *UDPSession) DeleteMessageIdArray(MessageIdArray []int) []int {
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
func (s *UDPSession) SetNote(note string) {
	s.Note = note
}
func (s *UDPSession) ActiveSend(isSendServer bool, _ int, data []byte) string {
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
func (s *UDPSession) GetTheology() int {
	return s.Theology
}

func (s *UDPSession) SetAnewInsert(anewInsert bool) {
	s.AnewInsert = anewInsert
}

func (s *UDPSession) ListMatch() bool {
	if s.ListFilter == nil {
		return true
	}
	return matchFilters(s.ToSearchAny(), s.ListFilter)
}
func (s *UDPSession) ToSearchAny() map[string]interface{} {

	testData := map[string]interface{}{
		"方式":     s.Method,
		"响应类型":   s.Method,
		"进程":     s.ProcessName,
		"状态":     Or(s.Disconnect, "已连接", "已断开"),
		"HOST":   s.RemoteAddress,
		"主机名":    s.RemoteAddress,
		"请求时间":   s.Time,
		"响应时间":   "",
		"Path":   "",
		"参数":     "",
		"注释":     s.Note,
		"来源地址":   s.ClientIP,
		"响应IP":   s.RemoteAddress,
		"身份验证账号": s.UserName,
		"请求地址":   s.ClientIP + " -> " + s.RemoteAddress,
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
func Or(b bool, s1, s2 string) string {
	if b {
		return s1
	}
	return s2
}
