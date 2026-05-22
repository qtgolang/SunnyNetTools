package Session

import (
	"bytes"
	"changeme/Service/Session/Conv"
	"encoding/binary"
	"github.com/qtgolang/SunnyNet/src/http"
	"github.com/qtgolang/SunnyNet/src/public"
	"io"
	"sort"
	"sync/atomic"
)

const ImportExportType_HTTP = "1"
const ImportExportType_TCP = "2"
const ImportExportType_UDP = "3"

type Insert struct {
	IsHTTP           bool
	Method           string
	URL              string
	ProcessName      string
	ClientIP         string
	Theology         int
	Time             string
	Ico              string
	Filter           bool
	UserName         string //身份验证的账号
	BreakMode        uint32 //断点模式
	GuaranteeDisplay bool   //保证显示此消息
	RemoteAddress    string
	Host             string
	Note             string
	State            string //状态 已连接/已断开

}

// 导入记录
func Import(_r io.Reader, InsertFuck func(session AppSession) Insert, Progress func(i int)) []Insert {
	r := &Reader{_r}
	ma := r.readInt()
	var IDArray []Insert
	mi := 0
	for {
		Type, err := r.readMessageType()
		if err != nil {
			if err != io.EOF {
				panic("文件数据错误")
			}
			break
		}
		mi++
		Progress(int((float64(mi) / float64(ma)) * 100))
		Theology := int(atomic.AddInt64(&public.Theology, 1))
		if Type == ImportExportType_HTTP {
			obj := r.readHttpSession()
			obj.Theology = Theology
			Session.Store(Theology, obj)
			IDArray = append(IDArray, InsertFuck(obj))
			continue
		}
		if Type == ImportExportType_TCP {
			obj := r.readTCPSession()
			if obj.SenLength == 0 && obj.RecLength == 0 {
				continue
			}
			obj.Theology = Theology
			Session.Store(Theology, obj)
			IDArray = append(IDArray, InsertFuck(obj))
			continue
		}
		if Type == ImportExportType_UDP {
			obj := r.readUDPSession()
			obj.Theology = Theology
			Session.Store(Theology, obj)
			IDArray = append(IDArray, InsertFuck(obj))
			continue
		}
		panic("文件数据错误")
	}

	return IDArray
}
func (r Reader) readMessageType() (string, error) {
	bs := make([]byte, 1)
	_, e := r.Read(bs)
	return string(bs), e
}

type Reader struct {
	io.Reader
}

func (r Reader) readInt() int {
	buf := make([]byte, 8)
	_, _ = io.ReadFull(r, buf)
	return int(binary.BigEndian.Uint64(buf))
}

func (r Reader) readBool() bool {
	b := []byte{0}
	_, _ = r.Read(b)
	return b[0] != 0
}

func (r Reader) readString() string {
	length := r.readInt()
	if length <= 0 {
		return ""
	}
	buf := make([]byte, length)
	_, _ = io.ReadFull(r, buf)
	return string(buf)
}

func (r Reader) readBody() []byte {
	length := r.readInt()
	if length <= 0 {
		return nil
	}
	buf := make([]byte, length)
	_, _ = io.ReadFull(r, buf)
	return buf
}

func (r Reader) readIntArray() []int {
	count := r.readInt()
	arr := make([]int, count)
	for i := 0; i < count; i++ {
		arr[i] = r.readInt()
	}
	return arr
}

func (r Reader) readHeader() http.Header {
	count := r.readInt()
	header := make(http.Header)
	for i := 0; i < count; i++ {
		key := r.readString()
		valCount := r.readInt()
		vals := make([]string, valCount)
		for j := 0; j < valCount; j++ {
			vals[j] = r.readString()
		}
		header[key] = vals
	}
	return header
}

func (r Reader) readHttpSessionRequest() HttpSessionRequest {
	return HttpSessionRequest{
		Method:      r.readString(),
		Url:         r.readString(),
		ClientIP:    r.readString(),
		Time:        r.readString(),
		ProcessName: r.readString(),
		Proto:       r.readString(),
		Header:      r.readHeader(),
		Body:        r.readBody(),
		BodyLength:  r.readInt(),
		IsMaxLength: r.readBool(),
	}
}

func (r Reader) readHttpSessionResponse() HttpSessionResponse {
	return HttpSessionResponse{
		Time:         r.readString(),
		Code:         r.readString(),
		State:        r.readString(),
		Type:         r.readString(),
		Length:       r.readInt(),
		Proto:        r.readString(),
		Header:       r.readHeader(),
		Body:         r.readBody(),
		ResponseType: r.readString(),
		ServerIP:     r.readString(),
		BodyLength:   r.readInt(),
		IsMaxLength:  r.readBool(),
	}
}

func (r Reader) readWebsocketStream() *WebsocketStream {
	return &WebsocketStream{
		MessageId:     r.readInt(),
		IsSend:        r.readBool(),
		IsActiveSend:  r.readBool(),
		Time:          r.readString(),
		Body:          r.readBody(),
		WebsocketType: r.readInt(),
		IsClose:       r.readBool(),
	}
}

func (r Reader) readWebsocketStreamMap() map[int]*WebsocketStream {
	count := r.readInt()
	m := make(map[int]*WebsocketStream)
	for i := 0; i < count; i++ {
		key := r.readInt()
		m[key] = r.readWebsocketStream()
	}
	return m
}

func (r Reader) readHttpSession() *HttpSession {
	return &HttpSession{
		Theology:            r.readInt(),
		Ico:                 r.readString(),
		Request:             r.readHttpSessionRequest(),
		Response:            r.readHttpSessionResponse(),
		Error:               r.readString(),
		State:               r.readInt(),
		RecLength:           r.readInt(),
		SenLength:           r.readInt(),
		Note:                r.readString(),
		WebsocketStream:     r.readWebsocketStreamMap(),
		WebsocketStreamKeys: r.readIntArray(),
		IsWebsocketRequest:  r.readBool(),
		WebsocketDisconnect: r.readBool(),
		UserName:            r.readString(),
	}
}
func (r Reader) readStream() *Stream {
	return &Stream{
		MessageId:     r.readInt(),
		IsSend:        r.readBool(),
		IsActiveSend:  r.readBool(),
		Time:          r.readString(),
		Body:          r.readBody(),
		IsClose:       r.readBool(),
		WebsocketType: r.readInt(),
	}
}

func (r Reader) readStreamMap() map[int]*Stream {
	count := r.readInt()
	m := make(map[int]*Stream, count)
	for i := 0; i < count; i++ {
		key := r.readInt()
		m[key] = r.readStream()
	}
	return m
}

func (r Reader) readTCPSession() *TCPSession {
	return &TCPSession{
		Theology:      r.readInt(),
		Method:        r.readString(),
		Ico:           r.readString(),
		Note:          r.readString(),
		ClientIP:      r.readString(),
		ProcessName:   r.readString(),
		RemoteAddress: r.readString(),
		Host:          r.readString(),
		Time:          r.readString(),
		RecLength:     r.readInt(),
		SenLength:     r.readInt(),
		Stream:        r.readStreamMap(),
		StreamKeys:    r.readIntArray(),
		UserName:      r.readString(),
	}
}
func (r Reader) readUDPSession() *UDPSession {
	return &UDPSession{
		Theology:      r.readInt(),
		Method:        r.readString(),
		Ico:           r.readString(),
		Note:          r.readString(),
		ClientIP:      r.readString(),
		ProcessName:   r.readString(),
		RemoteAddress: r.readString(),
		Time:          r.readString(),
		RecLength:     r.readInt(),
		SenLength:     r.readInt(),
		Stream:        r.readStreamMap(),
		StreamKeys:    r.readIntArray(),
		UserName:      r.readString(),
	}
}

type Writer struct {
	io.Writer
}

func (w Writer) writeHeader(app http.Header) {
	w.writeInt(len(app))
	for k, v := range app {
		w.writeString(k)
		w.writeInt(len(v))
		for _, vv := range v {
			w.writeString(vv)
		}
	}
}

func (w Writer) writeHttpSessionResponse(app HttpSessionResponse) {
	w.writeString(app.Time)
	w.writeString(app.Code)
	w.writeString(app.State)
	w.writeString(app.Type)
	w.writeInt(app.Length)
	w.writeString(app.Proto)
	w.writeHeader(app.Header)
	w.writeBody(app.Body)
	w.writeString(app.ResponseType)
	w.writeString(app.ServerIP)
	w.writeInt(app.BodyLength)
	w.writeBool(app.IsMaxLength)
}
func (w Writer) writeHttpSessionRequest(app HttpSessionRequest) {
	w.writeString(app.Method)
	w.writeString(app.Url)
	w.writeString(app.ClientIP)
	w.writeString(app.Time)
	w.writeString(app.ProcessName)
	w.writeString(app.Proto)
	w.writeHeader(app.Header)
	w.writeBody(app.Body)
	w.writeInt(app.BodyLength)
	w.writeBool(app.IsMaxLength)
}

func (w Writer) writeWebsocketStreamMap(app map[int]*WebsocketStream) {
	w.writeInt(len(app))
	for k, v := range app {
		w.writeInt(k)
		w.writeWebsocketStream(v)
	}
}

func (w Writer) writeHttpSession(app *HttpSession) {
	w.writeInt(app.Theology)
	w.writeString(app.Ico)
	w.writeHttpSessionRequest(app.Request)
	w.writeHttpSessionResponse(app.Response)
	w.writeString(app.Error)
	w.writeInt(app.State)
	w.writeInt(app.RecLength)
	w.writeInt(app.SenLength)
	w.writeString(app.Note)
	w.writeWebsocketStreamMap(app.WebsocketStream)
	w.writeIntArray(app.WebsocketStreamKeys)
	w.writeBool(app.IsWebsocketRequest)
	w.writeBool(app.WebsocketDisconnect)
	w.writeString(app.UserName)
}

func (w Writer) writeUDPSession(app *UDPSession) {
	w.writeInt(app.Theology)
	w.writeString(app.Method)
	w.writeString(app.Ico)
	w.writeString(app.Note)
	w.writeString(app.ClientIP)
	w.writeString(app.ProcessName)
	w.writeString(app.RemoteAddress)
	w.writeString(app.Time)
	w.writeInt(app.RecLength)
	w.writeInt(app.SenLength)
	w.writeStreamMap(app.Stream)
	w.writeIntArray(app.StreamKeys)
	w.writeString(app.UserName)
}
func (w Writer) writeTCPSession(app *TCPSession) {
	w.writeInt(app.Theology)
	w.writeString(app.Method)
	w.writeString(app.Ico)
	w.writeString(app.Note)
	w.writeString(app.ClientIP)
	w.writeString(app.ProcessName)
	w.writeString(app.RemoteAddress)
	w.writeString(app.Host)
	w.writeString(app.Time)
	w.writeInt(app.RecLength)
	w.writeInt(app.SenLength)
	w.writeStreamMap(app.Stream)
	w.writeIntArray(app.StreamKeys)
	w.writeString(app.UserName)
}
func (w Writer) writeStreamMap(app map[int]*Stream) {
	w.writeInt(len(app))
	for k, v := range app {
		w.writeInt(k)
		w.writeStream(v)
	}
}
func (w Writer) writeStream(app *Stream) {
	w.writeInt(app.MessageId)
	w.writeBool(app.IsSend)
	w.writeBool(app.IsActiveSend)
	w.writeString(app.Time)
	w.writeBody(app.Body)
	w.writeBool(app.IsClose) //IsClose
	w.writeInt(app.WebsocketType)

}

func (w Writer) writeWebsocketStream(app *WebsocketStream) {
	w.writeInt(app.MessageId)
	w.writeBool(app.IsSend)
	w.writeBool(app.IsActiveSend)
	w.writeString(app.Time)
	w.writeBody(app.Body)
	w.writeInt(app.WebsocketType)
	w.writeBool(app.IsClose) //IsClose
}

func (w Writer) writeBool(i bool) {
	if i {
		_, _ = w.Write([]byte{1})
	} else {
		_, _ = w.Write([]byte{0})
	}
}
func (w Writer) writeBody(i []byte) {
	_, _ = w.Write(Conv.IntToBytes(len(i)))
	_, _ = w.Write(i)
}
func (w Writer) writeInt(i int) {
	_, _ = w.Write(Conv.IntToBytes(i))
}
func (w Writer) writeIntArray(i []int) {
	_, _ = w.Write(Conv.IntToBytes(len(i)))
	for _, v := range i {
		_, _ = w.Write(Conv.IntToBytes(v))
	}
}
func (w Writer) writeString(i string) {
	_, _ = w.Write(Conv.IntToBytes(len(i)))
	_, _ = w.Write([]byte(i))
}

// 导出记录
func Export(writer io.Writer, _ExportTheology []int, Progress func(i int)) {
	W := &Writer{writer}
	var array []int
	array = _ExportTheology
	sort.Ints(array)
	ma := len(array)
	W.writeInt(ma)
	for mi, v := range array {
		value, ok := Session.Load(v)
		if ok {
			Progress(int((float64(mi) / float64(ma)) * 100))
			var buff bytes.Buffer
			buf := &Writer{&buff}
			export(buf, value)
			_, _ = W.Write(buff.Bytes())
		}
	}
}
func export(W *Writer, value any) bool {
	obj, ok := value.(AppSession)
	if !ok {
		return true
	}
	if obj.IsHTTP() {
		o, isOK := obj.(*HttpSession)
		if isOK {
			_, _ = W.Write([]byte(ImportExportType_HTTP))
			W.writeHttpSession(o)
		}
	}
	if obj.IsTCP() {
		o, isOK := obj.(*TCPSession)
		if isOK {
			_, _ = W.Write([]byte(ImportExportType_TCP))
			W.writeTCPSession(o)
		}
	}
	if obj.IsUDP() {
		o, isOK := obj.(*UDPSession)
		if isOK {
			_, _ = W.Write([]byte(ImportExportType_UDP))
			W.writeUDPSession(o)
		}
	}
	return true
}
