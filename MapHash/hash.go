package MapHash

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"github.com/qtgolang/SunnyNet/SunnyNet"
	"github.com/qtgolang/SunnyNet/public"
	"github.com/qtgolang/SunnyNet/src/GoWinHttp"
	"io"
	"net"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Map struct {
	Request      map[int]*Request
	lock         sync.Mutex
	UpdateLength map[int]*ResponseLength
}

type WaitGroup struct {
	lock      sync.Mutex
	waitGroup sync.WaitGroup
	i         int
}

func (v *WaitGroup) Add(i int) {
	v.lock.Lock()
	v.i += i
	v.waitGroup.Add(i)
	v.lock.Unlock()
}
func (v *WaitGroup) Done() {
	v.lock.Lock()
	v.i--
	if v.i < 0 {
		v.lock.Unlock()
		return
	}
	v.waitGroup.Done()
	v.lock.Unlock()
}
func (v *WaitGroup) Wait() {
	v.waitGroup.Wait()
}

type Request struct {
	PID      string      `json:"PID"`    //进程
	Method   string      `json:"Method"` //方式
	URL      string      `json:"URL"`    //请求地址
	Proto    string      `json:"Proto"`
	Header   http.Header `json:"Header"`
	Body     []byte      `json:"Body"`
	Display  bool        `json:"Display"` //是否需要显示到列表
	Response struct {
		Conn      *SunnyNet.HttpConn `json:"-"`
		Header    http.Header        `json:"Header"`
		Body      []byte             `json:"Body"`
		StateCode int                `json:"StateCode"`
		Error     bool               `json:"Error"`
	} `json:"Response"`
	Break   uint8              `json:"Break"`
	Wait    WaitGroup          `json:"-"`
	Conn    *SunnyNet.HttpConn `json:"-"`
	WsConn  *SunnyNet.WsConn   `json:"-"`
	TcpConn *SunnyNet.TcpConn  `json:"-"`
	UdpConn *SunnyNet.UDPConn  `json:"-"`
	Options struct {
		StopSend bool `json:"StopSend"`
		StopRec  bool `json:"StopRec"`
		StopALL  bool `json:"StopALL"`
	} `json:"StopSend"`
	SocketData []*UpdateSocketData `json:"SocketData"`
	SendTime   string              `json:"SendTime"`
	RecTime    string              `json:"RecTime"`
	SendNum    int                 `json:"SendNum"`
	RecNum     int                 `json:"RecNum"`
	Way        string              `json:"Way"`
	Notes      string              `json:"Notes"`
	ClientIP   string              `json:"ClientIP"`
	Color      struct {
		TagColor string `json:"TagColor"` //标记的文本颜色
		Search   string `json:"search"`   //搜索的背景颜色
	} `json:"color"` //显示图标
}
type RequestWeb struct {
	Method   string
	URL      string
	Proto    string //HTTP/1.1
	Header   http.Header
	Body     []byte
	Response struct {
		Header    http.Header
		Body      []byte
		StateCode int
		StateText string
		Error     bool
	}
	SocketData []*UpdateSocketList
	Options    struct {
		StopSend bool
		StopRec  bool
		StopALL  bool
	}
}
type ResponseLength struct {
	Send int
	Rec  int
}
type UpdateResponseLength struct {
	Send     int
	Rec      int
	Theology int
}

func (m *Map) GetALLResponseLength() []*UpdateResponseLength {
	m.lock.Lock()
	defer m.lock.Unlock()
	res := make([]*UpdateResponseLength, 0)
	for kk, vv := range m.UpdateLength {
		res = append(res, &UpdateResponseLength{Theology: kk, Rec: vv.Rec, Send: vv.Send})
	}
	m.UpdateLength = make(map[int]*ResponseLength)
	return res
}
func (m *Map) addResponseLength(TheologyID int, a, b int) {
	if m.UpdateLength[TheologyID] == nil {
		m.UpdateLength[TheologyID] = &ResponseLength{Send: a, Rec: b}
	} else {
		m.UpdateLength[TheologyID].Send = a
		m.UpdateLength[TheologyID].Rec = b
	}

}
func (m *Map) SetRequest(TheologyID int, h *Request) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.Request[TheologyID] = h
}
func (m *Map) CreateUniqueID() int {
	m.lock.Lock()
	defer m.lock.Unlock()
	tm := atomic.AddInt64(&public.Theology, 1)
	return int(tm)
}
func (m *Map) SetRequestUDP(TheologyID int, Conn *SunnyNet.UDPConn) *Request {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.Request[TheologyID] == nil {
		m.Request[TheologyID] = &Request{UdpConn: Conn}
	} else {
		m.Request[TheologyID].UdpConn = Conn
	}
	m.Request[TheologyID].Display = true
	return m.Request[TheologyID]
}
func (m *Map) SetRequestTCP(TheologyID int, Conn *SunnyNet.TcpConn) *Request {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.Request[TheologyID] == nil {
		m.Request[TheologyID] = &Request{TcpConn: Conn}
	} else {
		m.Request[TheologyID].TcpConn = Conn
	}
	m.Request[TheologyID].Display = true
	return m.Request[TheologyID]
}
func (m *Map) GetRequestWeb(Theology int) *RequestWeb {
	m.lock.Lock()
	defer m.lock.Unlock()
	h := m.Request[Theology]
	r := &RequestWeb{Body: h.Body, URL: h.URL, Proto: h.Proto, Header: h.Header, Method: h.Method}
	r.Response.Header = h.Response.Header
	r.Response.Body = h.Response.Body
	r.Response.StateText = http.StatusText(h.Response.StateCode)
	r.Response.StateCode = h.Response.StateCode
	r.Response.Error = h.Response.Error
	r.SocketData = make([]*UpdateSocketList, 0)
	for i := 0; i < len(h.SocketData); i++ {
		r.SocketData = append(r.SocketData, h.SocketData[i].Info)
	}
	r.Options.StopSend = h.Options.StopSend
	r.Options.StopRec = h.Options.StopRec
	r.Options.StopALL = h.Options.StopALL

	return r
}
func (m *Map) GetRequest(Theology int) *Request {
	m.lock.Lock()
	defer m.lock.Unlock()
	h := m.Request[Theology]
	return h
}
func (m *Map) SetOptions(Theology int, send, rec, all bool) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	h := m.Request[Theology]
	if h != nil {
		h.Options.StopRec = rec
		h.Options.StopALL = all
		h.Options.StopSend = send
	}
	return h != nil
}
func (m *Map) SetSocketData(Theology int, data *UpdateSocketData, up bool, num int) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	h := m.Request[Theology]
	if h != nil {
		if num > 0 {
			if up {
				h.SendNum += num
			} else {
				h.RecNum += num
			}
			m.addResponseLength(Theology, h.SendNum, h.RecNum)
		}
		h.SocketData = append(h.SocketData, data)
	}
	return h != nil
}
func (m *Map) SetSocketDataEmpty(Theology int) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	h := m.Request[Theology]
	if h != nil {
		h.SocketData = make([]*UpdateSocketData, 0)
		h.SendNum = 0
		h.RecNum = 0
		m.addResponseLength(Theology, h.SendNum, h.RecNum)
	}
	return h != nil
}
func NewHashMap() *Map {
	return &Map{Request: make(map[int]*Request)}
}
func (m *Map) Empty() {
	m.lock.Lock()
	defer m.lock.Unlock()
	mz := make(map[int]*Request)
	for k, v := range m.Request {
		if v != nil {
			if v.UdpConn != nil || v.TcpConn != nil || v.WsConn != nil {
				if v.Display {
					v.SocketData = make([]*UpdateSocketData, 0)
					v.RecNum = 0
					v.SendNum = 0
					mz[k] = v
				}
			}
			v.Wait.Done()
		}
	}
	m.Request = mz
	m.UpdateLength = make(map[int]*ResponseLength)
}

// ReleaseAll 全部放行
func (m *Map) ReleaseAll() {
	m.lock.Lock()
	defer m.lock.Unlock()
	for _, v := range m.Request {
		if v != nil {
			v.Wait.Done()
		}
	}
}
func (m *Map) Delete(TheologyArray []int) {
	m.lock.Lock()
	defer m.lock.Unlock()
	for _, k := range TheologyArray {
		v := m.Request[k]
		vv := m.UpdateLength[k]
		if vv != nil {
			vv.Rec = 0
			vv.Send = 0
		}
		if v != nil {
			if v.UdpConn != nil || v.TcpConn != nil || v.WsConn != nil {
				if v.Display {
					v.SocketData = make([]*UpdateSocketData, 0)
					v.RecNum = 0
					v.SendNum = 0
				}
			} else {
				delete(m.Request, k)
				delete(m.UpdateLength, k)
			}
			v.Wait.Done()
		}

	}
}
func (m *Map) Search(callSearch func(int, int, *Request)) {
	m.lock.Lock()
	defer m.lock.Unlock()
	max := float64(len(m.Request))
	i := float64(0)
	for k, v := range m.Request {
		i++
		callSearch(k, int(i/max*100), v)
	}
}
func (m *Map) CloseSession(TheologyArray []int) {
	m.lock.Lock()
	defer m.lock.Unlock()
	for _, k := range TheologyArray {
		v := m.Request[k]
		if v != nil {
			if v.TcpConn != nil {
				v.TcpConn.Close()
			}
			if v.WsConn != nil {
				v.WsConn.Close()
			}
		}
	}

}

func (m *Map) SaveToFile(Path string, All bool, TheologyArray []int, SetStatusText func(string)) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	var SaveData []*Request
	SetStatusText("正在统计需要储存的信息")
	if All {
		var keys []int
		for k := range m.Request {
			keys = append(keys, k)
		}
		sort.Ints(keys)
		for _, _ket := range keys {
			k := m.Request[_ket]
			if k != nil {
				if k.Display {
					SaveData = append(SaveData, k)
				}
			}
		}
	} else {
		for _, k := range TheologyArray {
			v := m.Request[k]
			if v != nil {
				SaveData = append(SaveData, v)
			}
		}
	}
	if len(SaveData) < 1 {
		SetStatusText("需要储存的数量小于1")
		return false
	}
	SetStatusText("有 " + strconv.Itoa(len(SaveData)) + " 条数据正在序列化储存...")
	bs, e := json.Marshal(&SaveData)
	if e != nil {
		SetStatusText("数据序列化储存失败！！")
		return false
	}
	SetStatusText("数据压缩中...")
	bs2 := BrCompress(bs)
	if len(bs2) < 1 {
		SetStatusText("数据失败！")
		return false
	}
	SetStatusText("正在写入文件")
	err := os.WriteFile(Path, bs2, 666)
	if err == nil {
		SetStatusText("保存记录文件成功：" + Path)
	} else {
		SetStatusText("保存记录文件失败：" + err.Error())
	}
	return err == nil
}

// Resend 重发请求
func (m *Map) Resend(TheologyArray []int, mode int, Port int) {
	m.lock.Lock()
	defer m.lock.Unlock()
	for _, k := range TheologyArray {
		v := m.Request[k]
		if v != nil {
			go resend(v, mode, Port)
		}
	}
}

func resend(m *Request, mode, Port int) {
	if m != nil {
		/*
			if m.Way == "Websocket" {
				resendWS(m, mode, Port)
			}
		*/
		if m.Way == "HTTP" {
			resendHttp(m, mode, Port)
		}
		/*
			else if m.Way == "UDP" {
				resendUDP(m, Port)
			} else if strings.Contains(strings.ToUpper(m.Way), "TCP") {
				resendTCP(m, strings.Contains(strings.ToUpper(m.Way), "TLS"), Port)
			}
		*/
	}

}

// 写完但是貌似有问题
func resendTCP(m *Request, _TLS bool, LocalSunnyNetPort int) {
	_t := strings.Split(m.URL, "->")
	if len(_t) != 2 {
		return
	}
	DataInfo := m.SocketData
	if DataInfo == nil {
		return
	}
	RemoteAddr := _t[1]
	uAddr := SunnyNet.TargetInfo{}
	uAddr.Parse(RemoteAddr, 0)
	if uAddr.Port == 0 {
		return
	}
	Conn, err := net.DialTimeout("tcp", "127.0.0.1:"+strconv.Itoa(LocalSunnyNetPort), time.Duration(10000)*time.Millisecond)
	defer func() { _ = Conn.Close() }()
	if err != nil {
		return
	}
	if GoWinHttp.ConnectS5(&Conn, &GoWinHttp.Proxy{}, uAddr.Host, uAddr.Port) == false {
		return
	}
	if _TLS {
		cfg := &tls.Config{ServerName: uAddr.Host}
		cfg.InsecureSkipVerify = true
		tlsConn := tls.Client(Conn, cfg)
		err = tlsConn.Handshake()
		if err != nil {
			return
		}
		Conn = tlsConn
	}
	var t time.Time
	var t2 time.Time
	for _, v := range m.SocketData {
		t2, err = time.Parse("2006-01-02 15:04:05.000", "2024-01-01 "+v.Info.Time)
		if t.Year() == 2024 {
			if err == nil {
				time.Sleep(t2.Sub(t))
			}
		}
		t = t2
		if v.Info.Ico == "上行" {
			_, _ = Conn.Write(v.Body)
			mx := make([]byte, 4096)
			_, _ = Conn.Read(mx)
		}
	}
	time.Sleep(time.Second)
}

// 未写完
func resendUDP(m *Request, LocalSunnyNetPort int) {
	_t := strings.Split(m.URL, "->")
	if len(_t) != 2 {
		return
	}
	DataInfo := m.SocketData
	if DataInfo == nil {
		return
	}
	RemoteAddr := _t[1]
	uAddr := SunnyNet.TargetInfo{}
	uAddr.Parse(RemoteAddr, 0)
	if uAddr.Port == 0 {
		return
	}
	PackMsg := func(data []byte) []byte {
		return nil
	}
	// 创建一个 UDP 地址
	serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:"+strconv.Itoa(LocalSunnyNetPort))
	if err != nil {
		return
	}
	// 连接 UDP 服务器
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		return
	}
	defer func() { _ = conn.Close() }()

	var t time.Time
	var t2 time.Time
	for _, v := range m.SocketData {
		t2, err = time.Parse("2006-01-02 15:04:05.000", "2024-01-01 "+v.Info.Time)
		if t.Year() == 2024 {
			if err == nil {
				time.Sleep(t2.Sub(t))
			}
		}
		t = t2
		if v.Info.Ico == "上行" {
			_, _ = conn.Write(PackMsg(v.Body))
			buffer := make([]byte, 4096)
			_ = conn.SetReadDeadline(time.Now().Add(time.Second))
			_, _, _ = conn.ReadFromUDP(buffer)
		}
	}
	time.Sleep(time.Second)

}

func resendWS(m *Request, mode, Port int) {

}
func resendHttp(m *Request, mode, SunnyNetServerPort int) {
	if m == nil {
		return
	}
	Body := io.NopCloser(bytes.NewBuffer(m.Body))
	defer func() {
		if Body != nil {
			_ = Body.Close()
		}
	}()
	h, e := http.NewRequest(m.Method, m.URL, Body)
	if e != nil {
		return
	}
	h.Header = m.Header.Clone()
	h.Header.Set("SunnyNetMode", strconv.Itoa(mode))
	w := GoWinHttp.NewGoWinHttp()
	w.SetProxyType(true)
	w.SetProxyIP("127.0.0.1:" + strconv.Itoa(SunnyNetServerPort))
	RES, _ := w.Do(h)
	if RES != nil {
		if RES.Body != nil {
			_ = RES.Body.Close()
		}
	}
}

type UpdateSocketList struct {
	Index    int    `json:"#"`
	Theology int    `json:"Theology"`
	BodyHash string `json:"数据"`
	Ico      string `json:"ico"`
	Time     string `json:"时间"`
	Length   int    `json:"长度"`
	WsType   string `json:"类型"`
	Color    string `json:"background"`
}

type UpdateSocketData struct {
	Info *UpdateSocketList `json:"Info"`
	Body []byte            `json:"Body"`
}
type RequestImg struct {
	Body string `json:"Body"`
	Type string `json:"Type"`
}

var multipartTag = []byte("--")
var multipartTag2 = []byte("\r\n\r\n")

func (r *Request) GetRequestImg() *RequestImg {
	if r == nil || r.Body == nil {
		return nil
	}
	ar := bytes.Split(r.Body, []byte("\n"))
	tag := make([]byte, 0)
	if len(ar) > 0 {
		tag = ar[0]
	}
	if !bytes.HasPrefix(tag, multipartTag) {
		return nil
	}
	ar = bytes.Split(r.Body, tag)
	for _, v := range ar {
		m := strings.ToLower(string(v))
		_type := strings.TrimSpace(public.SubString(m, "content-type:", "\n"))
		if strings.Contains(_type, "image/") {
			ar1 := strings.Split(_type, "/")
			if len(ar1) < 2 {
				continue
			}
			_type = ar1[1]
			ar2 := bytes.Split(v, multipartTag2)
			if len(ar2) < 2 {
				continue
			}
			bs := make([]byte, 0)
			for k, vv := range ar2 {
				if k == 0 {
					continue
				}
				if k == 1 {
					bs = append(bs, vv...)
					continue
				}
				bs = append(bs, multipartTag2...)
				bs = append(bs, vv...)
			}
			res := &RequestImg{}
			res.Body = base64.StdEncoding.EncodeToString(bs)
			res.Type = _type
			return res
		}
	}
	return nil
}
