package Session

import "sync"

type SeMap struct {
	sync.Map
}

func (s *SeMap) Len() int {
	sx := 0
	s.Range(func(key, value interface{}) bool {
		sx++
		return true
	})
	return sx
}

var Session SeMap

var ConnectDisconnect = []byte(`连接已断开`)

type UpdateSocketStream struct {
	Theology      int
	MessageId     int
	WebsocketType int
	Length        int
	IsSend        bool
	IsActiveSend  bool
	Time          string
	Body          string
	Filter        bool
	IsClose       bool
}

type Stream struct {
	MessageId     int
	IsSend        bool
	IsActiveSend  bool
	Time          string
	Body          []byte
	IsClose       bool
	WebsocketType int
	Note          string
	ListFilter    *Filter `json:"-"`
}

func (s *Stream) GetBody() []byte {
	return s.Body
}
func (s *Stream) GetMessageTime() string {
	return s.Time
}
func (s *Stream) GetIsClose() bool {
	return s.IsClose
}
func (s *Stream) GetIsSend() bool {
	return s.IsSend
}
func (s *Stream) GetWebsocketType() int {
	if s.WebsocketType == 0 {
		return 255
	}
	return s.WebsocketType
}
func (s *Stream) GetIsActiveSend() bool {
	return s.IsActiveSend
}

func (s *Stream) GetMessageId() int {
	return s.MessageId
}

const hexChars = "0123456789ABCDEF"

func getHexWithSpaces(data []byte) string {
	n := len(data)
	if n == 0 {
		return ""
	}
	if n > 20 {
		n = 20
	}
	buf := make([]byte, n*3-1)
	j := 0
	for i := 0; i < n; i++ {
		b := data[i]
		buf[j], buf[j+1] = hexChars[b>>4], hexChars[b&0x0F]
		if i < n-1 {
			buf[j+2] = ' '
			j += 3
		} else {
			j += 2
		}
	}
	return string(buf)
}
func GetHexAllSpaces(data []byte) string {
	n := len(data)
	if n == 0 {
		return ""
	}
	buf := make([]byte, n*3-1)
	j := 0
	for i := 0; i < n; i++ {
		b := data[i]
		buf[j], buf[j+1] = hexChars[b>>4], hexChars[b&0x0F]
		if i < n-1 {
			buf[j+2] = ' '
			j += 3
		} else {
			j += 2
		}
	}
	return string(buf)
}

type AppSession interface {
	/*
		Filter 获取已设置的主列表过滤器
	*/
	GetListFilter() *Filter
	/*
		SetFilter 设置主列表过滤器
	*/
	SetListFilter(*Filter)
	/*
	  设置注释
	*/
	SetNote(note string)
	/*
	  设置注释
	*/
	GetNote() string
	/*
		IsWait 当前是否正在断点拦截,仅HTTP请求有效， tcp/dup 始终返回false
	*/
	IsWait() bool
	/*
		Filter 获取已设置的过滤器
	*/
	GetStreamFilter() *Filter
	/*
		SetFilter 设置过滤器
	*/
	SetStreamFilter(*Filter)
	/*
		RangeStream 遍历全部消息流
		返回true 表示继续,否则等到全部遍历完
	*/
	RangeStream(func(AppStream) bool)
	/*
		GetStream 获取指定消息流

		传入 MessageId
	*/
	GetStream(int) AppStream
	/*
		ClearAllSession 清空全部消息流

		发送/接收的数据长度也会重置为 0
	*/
	ClearAllSession() []int
	/*
		GetRecLength 当前消息流中全部已接收的数据长度 包括手动主动接收的数据

		使用 DeleteMessageIdArray 删除后 对应的长度会 同步 减少
	*/
	GetRecLength() int
	/*
		GetSenLength 当前消息流中全部已经发送的数据长度 包括手动主动发送的数据

		使用 DeleteMessageIdArray 删除后 对应的长度会 同步 减少
	*/
	GetSenLength() int
	/*
		DeleteMessageIdArray 删除消息流
		传入 MessageId 数组
		返回已删除的 MessageId 数组
	*/
	DeleteMessageIdArray([]int) []int
	/*
		GetMethod 获取请求方式
		HTTP请求	返回	[HTTP]
		Websocket	请求返回	[Websocket]
		UDP	请求返回	[UDP]
		TCP	请求返回	[TCP] 或 [TLS-TCP] 或 [TCP-MUST]
	*/
	GetMethod() string
	/*
		IsHTTP 请求是否为 http 请求

		IsHTTP 为[true]那么 IsWebsocket 不一定为 true

		但是 IsWebsocket 为 true , IsHTTP 一定为 true
	*/
	IsHTTP() bool
	/*
		IsWebsocket 请求是否为 ws 请求

		如果是 ws 请求那么 IsHTTP 也一定为 true

		但是 IsHTTP 为 true , IsWebsocket 不一定为 true
	*/
	IsWebsocket() bool
	/*
		IsTCP 请求是否为 TCP 请求
	*/
	IsTCP() bool
	/*
		IsUDP 请求是否为 UDP 请求
	*/
	IsUDP() bool
	/*
		IsDisconnect ws,tcp,udp 是否已经断开连接
	*/
	IsDisconnect() bool
	/*
	   ResendRequest 重新发送此请求N次
	*/
	ResendRequest(count, BreakMode, port int, OutRouterIP string)
	/*
		ActiveSend 主动发送数据
	*/
	ActiveSend(bool, int, []byte) string
	/*
		AddStream 添加消息流
	*/
	AddStream(*Stream, ...int)
	/*
		GetTheology 获取当前消息唯一ID
	*/
	GetTheology() int
	/*
		SetAnewInsert 设置标记为 是否需要重新插入到列表中
	*/
	SetAnewInsert(bool)
	/*
		ListMatch  传入指定的过滤器  此消息 是否需要在列表中显示
	*/
	ListMatch() bool
	Lock()
	Unlock()
}
type AppStream interface {
	/*
		MatchFilters 传入指定的过滤器 对此消息进行匹配

		如果满足过滤器添加 返回 true

		如果过滤器为 nil 也返回true
	*/
	MatchFilters(*Filter) bool
	/*
		GetMessageId 获取当前消息列的 MessageId
	*/
	GetMessageId() int
	/*
		ToUpdateStream 转为更新消息结构体
	*/
	ToUpdateStream(int, *Filter) UpdateSocketStream
	/*
		GetMessageTime 获取消息时间
	*/
	GetMessageTime() string
	/*
		获取消息流中的 字节数组
	*/
	GetBody() []byte

	/*
		GetIsClose 当前消息流是否标记为已断开
	*/
	GetIsClose() bool
	/*
		GetIsSend 当前消息流是否为发送到服务器的数据
	*/
	GetIsSend() bool
	/*
		GetIsActiveSend 当前消息流是否为 (主动/手动) [发送/接收] 的数据
	*/
	GetIsActiveSend() bool
	/*
		GetWebsocketType 当前消息流中的ws消息类型，如果当前消息不是ws消息返回 255
	*/
	GetWebsocketType() int
}
