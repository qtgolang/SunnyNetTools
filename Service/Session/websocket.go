package Session

func (v *WebsocketStream) ToSearchAny() map[string]interface{} {
	_t := "接收"
	if v.IsSend {
		_t = "发送"
	}
	if v.IsActiveSend {
		_t = "手动" + _t
	}
	if v.IsClose {
		_t = "连接关闭"
	}
	{
		switch v.WebsocketType {
		case 1:
			_t += "(Text)"
			break
		case 2:
			_t += "(Binary)"
			break
		case 8:
			_t += "(Close)"
			break
		case 9:
			_t += "(Ping)"
			break
		case 10:
			_t += "(Pong)"
			break
		default:
			_t += "(invalid)"
			break
		}
	}
	testData := map[string]interface{}{
		"数据":   v.Body,
		"类型":   _t,
		"长度":   float64(len(v.Body)),
		"响应长度": float64(len(v.Body)),
		"时间":   v.Time,
	}
	return testData
}

// MatchFilters 根据匹配规则显示是否匹配成功，如果无规则，则默认返回 true
func (v *WebsocketStream) MatchFilters(filter *Filter) bool {
	if filter == nil {
		return true
	}
	return matchFilters(v.ToSearchAny(), filter)
}

func (s *HttpSession) AddStream(obj *Stream, wsType ...int) {
	if len(wsType) != 1 {
		return
	}
	stream := &WebsocketStream{MessageId: obj.MessageId, Time: obj.Time, WebsocketType: wsType[0], Body: obj.Body, IsSend: obj.IsSend, IsActiveSend: obj.IsActiveSend, IsClose: obj.IsClose}
	_, exists := s.WebsocketStream[obj.MessageId]
	if !exists {
		s.WebsocketStreamKeys = append(s.WebsocketStreamKeys, obj.MessageId)
	}
	s.WebsocketStream[obj.MessageId] = stream
}

type WebsocketStream struct {
	MessageId     int
	IsSend        bool
	IsActiveSend  bool
	Time          string
	Body          []byte
	WebsocketType int
	IsClose       bool
	Note          string
	ListFilter    *Filter `json:"-"`
}

func (v *WebsocketStream) GetNote() string {
	return v.Note
}

func (v *WebsocketStream) SetNote(s string) {
	v.Note = s
}

func (v *WebsocketStream) GetMessageTime() string {
	return v.Time
}
func (v *WebsocketStream) GetBody() []byte {
	return v.Body
}
func (v *WebsocketStream) GetIsClose() bool {
	return v.IsClose
}
func (v *WebsocketStream) GetIsSend() bool {
	return v.IsSend
}
func (v *WebsocketStream) GetWebsocketType() int {
	return v.WebsocketType
}

func (v *WebsocketStream) GetIsActiveSend() bool {
	return v.IsActiveSend
}

func (v *WebsocketStream) GetMessageId() int {
	return v.MessageId
}

func (v *WebsocketStream) ToUpdateStream(Theology int, filter *Filter) UpdateSocketStream {
	_stream := UpdateSocketStream{
		Theology:      Theology,
		MessageId:     v.MessageId,
		WebsocketType: v.WebsocketType,
		IsSend:        v.IsSend,
		Length:        len(v.Body),
		Time:          v.Time,
		Filter:        v.MatchFilters(filter),
		IsClose:       v.IsClose,
		IsActiveSend:  v.IsActiveSend,
	}
	if v.IsClose {
		_stream.Body = string(ConnectDisconnect)
	} else {
		_stream.Body = getHexWithSpaces(v.Body)
	}
	return _stream
}
