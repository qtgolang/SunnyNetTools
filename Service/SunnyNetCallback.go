package Service

import (
	. "changeme/Service/Config"
	"changeme/Service/Session"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/qtgolang/SunnyNet/SunnyNet"
	"github.com/qtgolang/SunnyNet/src/http"
	"github.com/qtgolang/SunnyNet/src/public"
	"github.com/sqweek/dialog"
)

func (g *AppMain) SaveSunnyFile(prompt string) (string, error) {
	if prompt == "" {
		prompt = "请选择要保存到的文件"
	}
	d := dialog.File().Title(prompt).Filter("SunnyNetV4记录文件", "sy4").SetStartFile("未命名.sy4")
	p, err := d.Save()
	if err != nil {
		// 用户取消：库会返回特定错误，直接当作空路径处理
		if err == dialog.ErrCancelled {
			return "", nil
		}
		return "", err
	}
	if p != "" && !strings.HasSuffix(strings.ToLower(p), ".sy4") {
		p += ".sy4"
	}
	return p, nil
}

func (g *AppMain) OpenSunnyFile(prompt string) (string, error) {
	if prompt == "" {
		prompt = "请选择要打开的文件"
	}
	p, err := dialog.File().Title(prompt).Filter("SunnyNetV4记录文件", "sy4").SetStartFile("未命名.sy4").Load()
	if err != nil {
		return "", err
	}
	return p, nil
}

// AppExport 导出记录 list 为空数组导出全部
func (g *AppMain) AppExport(list []int, savePath string) string {
	if len(list) < 0 {
		return "未选择任何请求"
	}
	_ = os.Remove(savePath)
	file, err := os.Create(savePath)
	if err != nil {
		return "打开文件或创建文件失败"
	}
	defer file.Close()
	brWriter := brotli.NewWriterLevel(file, brotli.DefaultCompression)
	defer brWriter.Close()
	Session.Export(brWriter, list, g.progress)
	return ""
}
func (g *AppMain) progress(i int) {
	//i=进度百分比
	AppList["Main"].EmitEvent("ExportProgress", i)
}

// AppImport 导入记录
func (g *AppMain) AppImport(filePath string) (err string, list []Session.Insert) {
	defer func() {
		if er := recover(); er != nil {
			err = "导入失败"
		}
	}()
	// 打开 br 文件
	file, e := os.Open(filePath)
	if e != nil {
		return "打开文件失败", nil
	}
	defer file.Close()
	brReader := brotli.NewReader(file)
	importUpdateSocketList = make([]updateSocket, 0)
	g.AppStartInsert()
	go func() {
		time.Sleep(time.Millisecond * 200)
		for _, v := range importUpdateSocketList {
			updateSocketMessage(v)
		}
	}()
	return "", Session.Import(brReader, InsertFuck, g.progress)
}

var importUpdateSocketList []updateSocket

func InsertFuck(obj Session.AppSession) Session.Insert {
	var Filter *Session.Filter
	if _Filter := Config.Filter; _Filter != nil {
		Filter = Session.ParseFilter(_Filter.String)
	}
	var newInsert Session.Insert
	if obj.IsWebsocket() {
		h, ok := obj.(*Session.HttpSession)
		if !ok {
			panic("type error. ")
		}
		h.Lock()
		defer h.Unlock()
		h.SetListFilter(Filter)
		h.Ico = "websocket_close"
		h.WebsocketDisconnect = true
		h.AnewInsert = false
		newInsert = Session.Insert{IsHTTP: true}
		newInsert.Method = h.Request.Method
		newInsert.URL = h.Request.Url
		newInsert.ClientIP = h.Request.ClientIP
		newInsert.ProcessName = h.Request.ProcessName
		newInsert.Theology = h.Theology
		newInsert.Time = h.Request.Time
		newInsert.Ico = h.Ico
		newInsert.UserName = h.UserName
		newInsert.Note = h.Note
		newInsert.Filter = h.ListMatch()
		//InsertList = append(InsertList, newInsert)
		obj2 := updateHTTPDone{}
		obj2.Theology = h.Theology
		obj2.Length = h.Response.Length
		obj2.Type = h.Response.Type
		obj2.Code = h.Response.Code
		obj2.Time = h.Response.Time
		obj2.Ico = h.Ico
		obj2.IP = h.Response.ServerIP
		obj2.Filter = h.ListMatch()
		updateDoneList = append(updateDoneList, obj2)
		obj3 := updateSocket{}
		obj3.Ico = h.Ico
		obj3.Code = "已断开"
		obj3.Theology = h.Theology
		obj3.RecLength = h.RecLength
		obj3.SenLength = h.SenLength
		obj3.Method = "Websocket"

		isDisconnectMessage := false
		h.RangeStream(func(stream Session.AppStream) bool {
			if stream.GetIsClose() {
				isDisconnectMessage = true
				return false
			}
			return true
		})
		if !isDisconnectMessage {
			stream := &Session.Stream{MessageId: -1, IsSend: false, IsActiveSend: false, Time: getTime(), Body: Session.ConnectDisconnect, IsClose: true, ListFilter: h.ListFilter}
			h.AddStream(stream, 8)
		}
		importUpdateSocketList = append(importUpdateSocketList, obj3)
		return newInsert
	}
	if obj.IsHTTP() {
		h := obj.(*Session.HttpSession)
		h.SetListFilter(Filter)
		h.AnewInsert = false
		newInsert = Session.Insert{IsHTTP: true}
		newInsert.Method = h.Request.Method
		newInsert.URL = h.Request.Url
		newInsert.ClientIP = h.Request.ClientIP
		newInsert.ProcessName = h.Request.ProcessName
		newInsert.Theology = obj.GetTheology()
		newInsert.Time = h.Request.Time
		newInsert.Ico = h.Ico
		newInsert.UserName = h.UserName
		newInsert.Note = h.Note
		newInsert.Filter = h.ListMatch()
		newInsert.BreakMode = BreakNone
		//InsertList = append(InsertList, newInsert)
		obj2 := updateHTTPDone{}
		obj2.Theology = h.Theology
		obj2.Length = h.Response.Length
		obj2.Type = h.Response.Type
		obj2.Code = h.Response.Code
		obj2.Time = h.Response.Time
		obj2.Ico = h.Ico
		obj2.IP = h.Response.ServerIP
		obj2.Filter = h.ListMatch()
		updateDoneList = append(updateDoneList, obj2)
		return newInsert
	}
	if obj.IsTCP() {
		t := obj.(*Session.TCPSession)
		t.SetListFilter(Filter)
		t.AnewInsert = false
		t.Ico = "websocket_close"
		t.Disconnect = true
		newInsert = Session.Insert{Theology: t.Theology, Method: t.Method, ClientIP: t.ClientIP, ProcessName: t.ProcessName, Time: getTime(), Ico: "websocket_close"}
		newInsert.Host = t.Host
		newInsert.RemoteAddress = t.RemoteAddress
		newInsert.UserName = t.UserName
		newInsert.Note = t.Note
		newInsert.Filter = obj.ListMatch()
		//InsertList = append(InsertList, _obj)
		res := updateSocket{Code: "已断开"}
		res.Theology = t.Theology
		res.RecLength = t.RecLength
		res.SenLength = t.SenLength
		res.Ico = "websocket_close"
		res.Code = "已断开"
		res.Method = t.Method
		isDisconnectMessage := false
		t.RangeStream(func(stream Session.AppStream) bool {
			if stream.GetIsClose() {
				isDisconnectMessage = true
				return false
			}
			return true
		})
		if !isDisconnectMessage {
			stream := &Session.Stream{MessageId: -1, IsSend: false, IsActiveSend: false, Time: getTime(), Body: Session.ConnectDisconnect, IsClose: true, ListFilter: t.ListFilter}
			t.AddStream(stream)
		}
		updateSocketMessage(res)
		return newInsert
	}
	if obj.IsUDP() {
		t := obj.(*Session.UDPSession)
		t.SetListFilter(Filter)
		t.AnewInsert = false
		t.Ico = "websocket_close"
		t.Disconnect = true
		newInsert = Session.Insert{Theology: t.Theology, Method: t.Method, ClientIP: t.ClientIP, ProcessName: t.ProcessName, Time: getTime(), Ico: "websocket_close"}
		newInsert.Host = t.RemoteAddress
		newInsert.RemoteAddress = newInsert.Host
		newInsert.ClientIP = t.ClientIP
		newInsert.Time = t.Time
		newInsert.ProcessName = t.ProcessName
		newInsert.Method = t.Method
		newInsert.UserName = t.UserName
		newInsert.Note = t.Note
		newInsert.Filter = obj.ListMatch()
		//InsertList = append(InsertList, _obj)
		res := updateSocket{Code: "已断开"}
		res.Theology = t.Theology
		res.RecLength = t.RecLength
		res.SenLength = t.SenLength
		res.Ico = "websocket_close"
		res.Code = "已断开"
		res.Method = t.Method
		isDisconnectMessage := false
		t.RangeStream(func(stream Session.AppStream) bool {
			if stream.GetIsClose() {
				isDisconnectMessage = true
				return false
			}
			return true
		})
		if !isDisconnectMessage {
			stream := &Session.Stream{MessageId: -1, IsSend: false, IsActiveSend: false, Time: getTime(), Body: Session.ConnectDisconnect, IsClose: true, ListFilter: t.ListFilter}
			t.AddStream(stream)
		}
		updateSocketMessage(res)
		return newInsert
	}
	panic("未知类型")
}
func tcpInsert(Conn SunnyNet.ConnTCP) *Session.TCPSession {
	Theology := Conn.Theology()
	lock.Lock()
	obj := Session.GetTCPSession(Theology)
	if obj == nil {
		lock.Unlock()
		return nil
	}

	lock.Unlock()
	return obj
}
func udpInsert(Conn SunnyNet.ConnUDP) *Session.UDPSession {
	Theology := Conn.Theology()
	lock.Lock()
	obj := Session.GetUDPSession(Theology)
	if obj == nil {
		obj = &Session.UDPSession{
			Theology:      Conn.Theology(),
			AnewInsert:    true,
			Stream:        make(map[int]*Session.Stream),
			Conn:          Conn,
			Method:        "UDP",
			RemoteAddress: Conn.RemoteAddress(),
			ClientIP:      Conn.LocalAddress(),
			Time:          getTime(),
			ProcessName:   fmt.Sprintf("(%d)", Conn.PID()) + Conn.GetProcessName(),
			UserName:      Conn.GetSocket5User(),
			Note:          Conn.GetNote(),
		}
		obj.Disconnect = true
		Session.Session.Store(Conn.Theology(), obj)
	}
	lock.Unlock()
	return obj
}
func updateSocketMessage(data updateSocket) {
	lock.Lock()
	ok := false
	for i, v := range updateSocket_List {
		if v.Theology == data.Theology {
			if strings.Contains(updateSocket_List[i].Code, "断开") && !strings.Contains(data.Code, "断开") {
				data.Code = updateSocket_List[i].Code
				updateSocket_List[i] = data
				fmt.Println("请检查代码：此处发现疑似异常：先触发了断开，但后续还有数据来....")
				fmt.Println("请检查代码：此处发现疑似异常：先触发了断开，但后续还有数据来....")
				fmt.Println("请检查代码：此处发现疑似异常：先触发了断开，但后续还有数据来....")
				fmt.Println("请检查代码：此处发现疑似异常：先触发了断开，但后续还有数据来....")
			} else {
				updateSocket_List[i] = data
			}
			ok = true
			break
		}
	}
	if !ok {
		updateSocket_List = append(updateSocket_List, data)
	}
	lock.Unlock()
}
func isBreakSend(req *Session.HttpSession) bool {
	if atomic.LoadUint32(&Config.BreakMode) == BreakSend {
		return true
	}
	return false
}
func isBreakRece(req *Session.HttpSession) bool {
	if atomic.LoadUint32(&Config.BreakMode) == BreakRece {
		return true
	}
	if atomic.LoadUint32(&req.NextBreakMode) == BreakRece {
		return true
	}
	return false
}

func (g *AppMain) httpCallback(Conn SunnyNet.ConnHTTP) {
	if strings.Contains(Conn.URL(), "http://"+LocalServer) {
		return
	}
	isBreak := Config.ReplaceHttp(Conn)
	if Conn.Type() == public.HttpSendRequest {
		{
			req := Session.NewHttpSession()
			req.Lock()
			req.State = public.HttpSendRequest
			req.UserName = Conn.GetSocket5User()
			req.Theology = Conn.Theology()
			req.Request.Body = Conn.GetRequestBody()
			req.Request.Url = Conn.URL()
			req.Request.Proto = Conn.Proto()
			req.Request.Method = Conn.Method()
			req.Request.Time = getTime()
			req.Request.ClientIP = Conn.ClientIP()
			req.Request.ProcessName = fmt.Sprintf("(%d)", Conn.PID()) + Conn.GetProcessName()
			req.Request.Header = Conn.GetRequestHeader()
			req.ListFilter = Config.Filter.Clone()
			Session.Session.Store(Conn.Theology(), req)
			ClientBreakMode := req.Request.Header.Get(public.HTTPClientTags + "_BreakMode")
			req.Request.Header.Del(public.HTTPClientTags + "_BreakMode")
			_GuaranteeDisplay := req.Request.Header.Get(public.HTTPClientTags+"_GuaranteeDisplay") == "true"
			req.Request.Header.Del(public.HTTPClientTags + "_GuaranteeDisplay")
			if ClientBreakMode == "Send" || Conn.Error() == "Debug" {
				isBreak = true
			} else if ClientBreakMode == "Rec" {
				req.NextBreakMode = BreakRece
			}
			if !isBreak {
				isBreak = isBreakSend(req)
			}
			req.Ico = Session.Or(isBreak, "拦截上行", "上行")
			lock.Lock()
			obj := Session.Insert{GuaranteeDisplay: _GuaranteeDisplay, IsHTTP: true}
			obj.Method = req.Request.Method
			obj.URL = req.Request.Url
			obj.ClientIP = req.Request.ClientIP
			obj.ProcessName = req.Request.ProcessName
			obj.Theology = Conn.Theology()
			obj.Time = req.Request.Time
			obj.Ico = req.Ico
			obj.UserName = req.UserName
			if isBreak {
				obj.BreakMode = BreakSend
				obj.Filter = true
			} else {
				obj.Filter = req.ListMatch()
				obj.BreakMode = BreakNone
			}
			note := Conn.GetNote()
			if note != "" {
				obj.Note = note
			}
			InsertList = append(InsertList, obj)
			lock.Unlock()
			if isBreak {
				req.Unlock()
				req.Wg.Add(1)
				req.Wg.Wait()
				req.Lock()
				//恢复为发送中的状态
				{
					lock.Lock()
					{
						if req.Request.Url != Conn.URL() {
							Conn.UpdateURL(req.Request.Url)
						}
						_t := Conn.GetRequestHeader()
						for k, _ := range _t {
							if req.Request.Header[k] == nil {
								_t.Del(k)
							}
						}
						for k, v := range req.Request.Header {
							_t.SetArray(k, v)
						}
						if req.Request.IsUpdateBody {
							Conn.SetRequestBody(req.Request.Body)
							req.Request.IsUpdateBody = false
						}
					}
					obj2 := httpUpdateSend{}
					obj2.Theology = Conn.Theology()
					obj2.Ico = "上行"
					obj2.URL = req.Request.Url
					obj2.Method = req.Request.Method
					obj2.BreakMode = BreakNone
					obj2.Filter = req.ListMatch()
					note = Conn.GetNote()
					if note != "" {
						obj2.Note = note
					}
					//避免重复发送更新请求,但这里实际可能不会同时触发多个更新,但是还是写上把
					{
						ok := false
						for i, v := range updateSendList {
							if v.Theology == obj2.Theology {
								updateSendList[i].Theology = obj2.Theology
								updateSendList[i].Ico = obj2.Ico
								updateSendList[i].Filter = obj2.Filter
								updateSendList[i].BreakMode = obj2.BreakMode
								ok = true
								break
							}
						}
						if !ok {
							updateSendList = append(updateSendList, obj2)
						}
					}
					lock.Unlock()
				}
			}
			if Conn.GetResponseCode() != 0 {
				time.Sleep(time.Millisecond * 200)
				{
					ResponseType := ""
					{
						array := strings.Split(Conn.GetResponseHeader().Get("Content-Type")+";", ";")
						if len(array) > 0 {
							ResponseType = array[0]
						}
					}
					bs := Conn.GetResponseBody()
					req.State = public.HttpResponseOK
					req.Ico = UpdateIco(Conn, ResponseType)
					req.Response.ResponseType = ResponseType
					req.Response.Code = fmt.Sprintf("%d", Conn.GetResponseCode())
					req.Response.State = http.StatusText(Conn.GetResponseCode())
					req.Response.Time = getTime()
					req.Response.Type = ResponseType
					req.Response.Length = len(bs)
					req.Response.Body = bs
					req.Response.Proto = Conn.Proto()
					req.Response.Header = Conn.GetResponseHeader()
					req.Response.ServerIP = "本地响应"
					req.ListFilter = Config.Filter.Clone()
					isBreak = isBreakRece(req)
					lock.Lock()
					obj2 := updateHTTPDone{}
					obj2.Theology = Conn.Theology()
					obj2.Length = req.Response.Length
					obj2.Type = req.Response.Type
					obj2.Code = req.Response.Code
					obj2.Time = req.Response.Time
					if isBreak {
						obj.BreakMode = BreakRece
						obj2.Ico = "拦截下行"
						obj2.Filter = true
					} else {
						obj.BreakMode = BreakNone
						obj2.Ico = req.Ico
						obj2.Filter = req.ListMatch()
					}
					note = Conn.GetNote()
					if note != "" {
						obj2.Note = note
					}
					updateDoneList = append(updateDoneList, obj2)
					lock.Unlock()
					if isBreak {
						req.Unlock()
						req.Wg.Add(1)
						req.Wg.Wait()
						req.Lock()
						lock.Lock()
						{
							Conn.SetResponseBody(req.Response.Body)
							p, _ := strconv.Atoi(req.Response.Code)
							if p != 0 {
								Conn.SetResponseCode(p)
							}
							_t := Conn.GetResponseHeader()
							for k, _ := range _t {
								if req.Response.Header[k] == nil {
									_t.Del(k)
								}
							}
							for k, v := range req.Response.Header {
								_t.SetArray(k, v)
							}
							if req.Response.IsUpdateBody {
								Conn.SetResponseBody(req.Response.Body)
								req.Response.IsUpdateBody = false
							}
						}
						obj2 = updateHTTPDone{}
						obj2.Theology = Conn.Theology()
						obj2.Length = req.Response.Length
						obj2.Type = req.Response.Type
						obj2.Code = req.Response.Code
						obj2.Time = req.Response.Time
						obj2.IP = req.Response.ServerIP
						obj2.Filter = req.ListMatch()
						obj.BreakMode = BreakNone
						obj2.Ico = req.Ico
						updateDoneList = append(updateDoneList, obj2)
						lock.Unlock()
					}
				}
			}
			req.Unlock()
		}
		return
	}
	if Conn.Type() == public.HttpResponseOK {
		{
			req := Session.GetHttpSession(Conn.Theology())
			if req == nil {
				return
			}
			req.Lock()
			req.State = public.HttpResponseOK
			ResponseType := ""
			ResponseHeader := Conn.GetResponseHeader()
			bs := Conn.GetResponseBody()
			{
				array := strings.Split(ResponseHeader.Get("Content-Type")+";", ";")
				if len(array) > 0 {
					ResponseType = array[0]
				}
			}
			req.Ico = UpdateIco(Conn, ResponseType)
			req.Response.ResponseType = ResponseType
			req.Response.Code = fmt.Sprintf("%d", Conn.GetResponseCode())
			req.Response.State = http.StatusText(Conn.GetResponseCode())
			req.Response.Time = getTime()
			req.Response.Type = ResponseType
			req.Response.Length = len(bs)
			req.Response.Body = bs
			req.Response.Proto = Conn.GetResponseProto()
			req.Response.Header = ResponseHeader
			req.Response.ServerIP = Conn.ServerAddress()
			req.ListFilter = Config.Filter.Clone()
			if !isBreak {
				isBreak = isBreakRece(req)
			}
			if !isBreak {
				if Conn.Error() == "Debug" {
					isBreak = true
				}
			}
			lock.Lock()
			obj := updateHTTPDone{}
			obj.Theology = Conn.Theology()
			obj.Length = req.Response.Length
			obj.Type = req.Response.Type
			obj.Code = req.Response.Code
			obj.Time = req.Response.Time
			obj.Ico = req.Ico
			obj.IP = req.Response.ServerIP
			if isBreak {
				obj.BreakMode = BreakRece
				obj.Ico = "拦截下行"
				obj.Filter = true
			} else {
				obj.BreakMode = BreakNone
				obj.Ico = req.Ico
				obj.Filter = req.ListMatch()
			}
			note := Conn.GetNote()
			if note != "" {
				obj.Note = note
			}
			updateDoneList = append(updateDoneList, obj)
			lock.Unlock()
			if isBreak {
				req.Unlock()
				req.Wg.Add(1)
				req.Wg.Wait()
				req.Lock()
				lock.Lock()
				{
					p, _ := strconv.Atoi(req.Response.Code)
					if p != 0 {
						Conn.SetResponseCode(p)
					}
					_t := Conn.GetResponseHeader()
					for k, _ := range _t {
						if req.Response.Header[k] == nil {
							_t.Del(k)
						}
					}
					for k, v := range req.Response.Header {
						_t.SetArray(k, v)
					}
					if req.Response.IsUpdateBody {
						Conn.SetResponseBody(req.Response.Body)
						req.Response.IsUpdateBody = false
					}
				}
				obj = updateHTTPDone{}
				obj.Theology = Conn.Theology()
				obj.Length = req.Response.Length
				obj.Type = req.Response.Type
				obj.Code = req.Response.Code
				obj.Time = req.Response.Time
				obj.Ico = req.Ico
				obj.Filter = req.ListMatch()
				obj.BreakMode = BreakNone
				obj.Ico = req.Ico
				note = Conn.GetNote()
				if note != "" {
					obj.Note = note
				}
				updateDoneList = append(updateDoneList, obj)
				lock.Unlock()
			}
			if Conn.GetResponseCode() == 101 {
				for len(updateDoneList) != 0 {
					time.Sleep(time.Millisecond * 100)
				}
			}
			req.Unlock()
		}
		return
	}
	if Conn.Type() == public.HttpRequestFail {
		{
			_store := false
			req := Session.GetHttpSession(Conn.Theology())
			if req == nil {
				req = Session.NewHttpSession()
				_store = true
				req.UserName = Conn.GetSocket5User()
				req.Theology = Conn.Theology()
				req.Request.Body = Conn.GetRequestBody()
				req.Request.Url = Conn.URL()
				req.Request.Proto = Conn.Proto()
				req.Request.Method = Conn.Method()
				req.Request.Time = getTime()
				req.Request.ClientIP = Conn.ClientIP()
				req.Request.ProcessName = fmt.Sprintf("(%d)", Conn.PID()) + Conn.GetProcessName()
				req.Request.Header = Conn.GetRequestHeader()
				req.ListFilter = Config.Filter.Clone()
				req.Response.Body = []byte(Conn.Error())
				Session.Session.Store(Conn.Theology(), req)
			}
			req.Lock()
			defer req.Unlock()
			req.State = public.HttpRequestFail
			req.Error = Conn.Error()
			if !_store {
				req.ListFilter = Config.Filter.Clone()
			}
			req.Ico = "error"
			lock.Lock()
			defer lock.Unlock()
			if _store {
				insertObj := Session.Insert{IsHTTP: true, State: "错误"}
				insertObj.Method = req.Request.Method
				insertObj.URL = req.Request.Url
				insertObj.ClientIP = req.Request.ClientIP
				insertObj.ProcessName = req.Request.ProcessName
				insertObj.Theology = Conn.Theology()
				insertObj.Time = req.Request.Time
				insertObj.Ico = req.Ico
				insertObj.Method = "错误"
				insertObj.UserName = req.UserName
				insertObj.Filter = req.ListMatch()
				insertObj.BreakMode = BreakNone
				note := Conn.GetNote()
				if note != "" {
					insertObj.Note = note
				}
				InsertList = append(InsertList, insertObj)
			} else {
				obj := updateHTTPError{}
				obj.Theology = Conn.Theology()
				obj.Length = 0
				obj.Code = "错误"
				obj.Time = getTime()
				obj.Ico = "error"
				obj.Filter = req.ListMatch()
				note := Conn.GetNote()
				if note != "" {
					obj.Note = note

				}
				updateErrorList = append(updateErrorList, obj)
			}
		}
		return
	}
}
func (g *AppMain) wsCallback(Conn SunnyNet.ConnWebSocket) {
	Config.ReplaceWebsocket(Conn)
	req := Session.GetHttpSession(Conn.Theology())
	if req == nil {
		return
	}
	req.Lock()
	defer req.Unlock()
	obj := updateSocket{}
	obj.Code = "已连接"
	obj.Theology = Conn.Theology()
	MessageId := Conn.MessageId()
	ConnType := Conn.Type()
	noteNew := Conn.GetNote()
	if noteNew != "" {
		req.SetNote(noteNew)
	}
	switch ConnType {
	case public.WebsocketConnectionOK:
		req.Ico = "websocket_connect"
		req.WebsocketStream = make(map[int]*Session.WebsocketStream)
		req.WebsocketStreamKeys = make([]int, 0)
		req.IsWebsocketRequest = true
		req.WebsocketConn = Conn
		if _Filter := req.GetListFilter(); _Filter != nil {
			req.SetStreamFilter(Session.ParseFilter(_Filter.String))
			obj.Filter = req.ListMatch()
		} else {
			obj.Filter = true
		}
		break
	case public.WebsocketUserSend, public.WebsocketServerSend:
		bs := Conn.Body()
		if ConnType == public.WebsocketUserSend {
			req.SenLength += len(bs)
		} else {
			req.RecLength += len(bs)
		}
		stream := &Session.Stream{MessageId: MessageId, Time: getTime(), Body: bs, IsSend: ConnType == public.WebsocketUserSend, ListFilter: req.ListFilter, WebsocketType: Conn.MessageType()}
		req.AddStream(stream, Conn.MessageType())
		if req.WebsocketConn != nil {
			req.WebsocketConn = Conn
		}
		obj.Stream = stream
		if atomic.LoadInt64(&Config.CurrentTheology) == int64(obj.Theology) {
			_stream := stream.ToUpdateStream(obj.Theology, req.GetStreamFilter())
			lock.Lock()
			updateSocketStreamList = append(updateSocketStreamList, _stream)
			lock.Unlock()
		}
		lock.Lock()
		if req.AnewInsert {
			if req.RecLength != 0 && req.SenLength != 0 {
				req.ListFilter = Config.Filter.Clone()
				req.AnewInsert = false
				newInsert := Session.Insert{IsHTTP: true}
				newInsert.Method = req.Request.Method
				newInsert.URL = req.Request.Url
				newInsert.ClientIP = req.Request.ClientIP
				newInsert.ProcessName = req.Request.ProcessName
				newInsert.Theology = req.Theology
				newInsert.Time = req.Request.Time
				newInsert.Ico = req.Ico
				newInsert.UserName = req.UserName
				newInsert.Filter = req.ListMatch()
				note := Conn.GetNote()
				if note != "" {
					newInsert.Note = note
				}
				InsertList = append(InsertList, newInsert)
			}
		}
		lock.Unlock()
		break
	case public.WebsocketDisconnect:
		req.Ico = "websocket_close"
		obj.Code = "已断开"
		stream := &Session.Stream{MessageId: MessageId, Time: getTime(), Body: Session.ConnectDisconnect, IsSend: false, IsClose: true, ListFilter: req.ListFilter, WebsocketType: 8}
		req.AddStream(stream, Conn.MessageType())
		req.WebsocketDisconnect = true
		req.WebsocketConn = nil
		obj.Stream = stream
		if atomic.LoadInt64(&Config.CurrentTheology) == int64(obj.Theology) {
			_stream := stream.ToUpdateStream(obj.Theology, req.GetStreamFilter())
			lock.Lock()
			updateSocketStreamList = append(updateSocketStreamList, _stream)
			lock.Unlock()
		}
		break
	}
	obj.Ico = req.Ico
	obj.RecLength = req.RecLength
	obj.SenLength = req.SenLength
	obj.Method = "Websocket"
	obj.Note = req.GetNote()
	if obj.Stream != nil {
		obj.Filter = obj.Stream.ListIsMatch()
	}
	updateSocketMessage(obj)
}
func (g *AppMain) tcpCallback(Conn SunnyNet.ConnTCP) {
	RemoteAddress := Conn.RemoteAddress()
	if strings.Contains(RemoteAddress, LocalServer) {
		return
	}
	Config.ReplaceTCP(Conn)
	res := updateSocket{Code: "已连接"}
	res.Theology = Conn.Theology()
	_Type := Conn.Type()
	noteNew := Conn.GetNote()
	switch _Type {
	case public.SunnyNetMsgTypeTCPAboutToConnect:
		{
			obj := &Session.TCPSession{Theology: Conn.Theology(), AnewInsert: true, Stream: make(map[int]*Session.Stream), Conn: Conn}
			obj.Lock()
			defer obj.Unlock()
			obj.UserName = Conn.GetSocket5User()
			obj.Method = strings.ToUpper(string(Conn.Body()))
			obj.ClientIP = Conn.LocalAddress()
			obj.ProcessName = fmt.Sprintf("(%d)", Conn.PID()) + Conn.GetProcessName()
			obj.RemoteAddress = RemoteAddress
			obj.Disconnect = true
			obj.Time = getTime()
			if noteNew != "" {
				obj.SetNote(noteNew)
			}
			arr := strings.Split(RemoteAddress, " -> ")
			if len(arr) >= 2 {
				obj.RemoteAddress = arr[1]
				obj.Host = arr[0]
			} else {
				obj.RemoteAddress = RemoteAddress
				obj.Host = RemoteAddress
			}
			Session.Session.Store(Conn.Theology(), obj)
			//即将开始连接
		}
		break
	case public.SunnyNetMsgTypeTCPConnectOK:
		{
			//连接成功
			obj := Session.GetTCPSession(Conn.Theology())
			if obj != nil {
				obj.Lock()
				defer obj.Unlock()
				obj.Conn = Conn
				if noteNew != "" {
					obj.SetNote(noteNew)
				}
			}
		}
		break
	case public.SunnyNetMsgTypeTCPClientSend:
		{
			//发送数据
			obj := tcpInsert(Conn)
			if obj == nil {
				return
			}
			obj.Lock()
			defer obj.Unlock()
			if noteNew != "" {
				obj.SetNote(noteNew)
			}
			{
				stream := &Session.Stream{MessageId: Conn.MessageId(), IsSend: true, IsActiveSend: false, Time: getTime(), Body: Conn.Body(), ListFilter: obj.ListFilter}
				res.Stream = stream
				obj.AddStream(stream)
				obj.SenLength += Conn.BodyLen()
				res.RecLength = obj.RecLength
				res.SenLength = obj.SenLength
				res.Method = obj.Method
				res.Note = obj.GetNote()
				if atomic.LoadInt64(&Config.CurrentTheology) == int64(obj.Theology) {
					_stream := stream.ToUpdateStream(obj.Theology, obj.GetStreamFilter())
					lock.Lock()
					updateSocketStreamList = append(updateSocketStreamList, _stream)
					lock.Unlock()
				}
				if obj.AnewInsert {
					if obj.RecLength != 0 || obj.SenLength != 0 {
						obj.ListFilter = Config.Filter.Clone()
						obj.AnewInsert = false
						_obj := Session.Insert{Theology: res.Theology, Method: obj.Method, ClientIP: Conn.LocalAddress(), ProcessName: Conn.GetProcessName(), Time: getTime(), Ico: "websocket_connect"}
						_obj.Host = obj.Host
						_obj.RemoteAddress = obj.RemoteAddress
						_obj.UserName = obj.UserName
						_obj.Filter = obj.ListMatch()
						_obj.Note = obj.GetNote()
						InsertList = append(InsertList, _obj)
					}
				}
			}
		}
		break
	case public.SunnyNetMsgTypeTCPClientReceive:
		{
			//接收
			obj := tcpInsert(Conn)
			if obj == nil {
				return
			}
			obj.Lock()
			defer obj.Unlock()
			if noteNew != "" {
				obj.SetNote(noteNew)
			}
			{
				stream := &Session.Stream{MessageId: Conn.MessageId(), IsSend: false, IsActiveSend: false, Time: getTime(), Body: Conn.Body(), ListFilter: obj.ListFilter}
				res.Stream = stream
				obj.AddStream(stream)
				obj.RecLength += Conn.BodyLen()
				res.RecLength = obj.RecLength
				res.SenLength = obj.SenLength
				res.Method = obj.Method
				res.Note = obj.GetNote()
				if atomic.LoadInt64(&Config.CurrentTheology) == int64(obj.Theology) {
					_stream := stream.ToUpdateStream(obj.Theology, obj.GetStreamFilter())
					lock.Lock()
					updateSocketStreamList = append(updateSocketStreamList, _stream)
					lock.Unlock()
				}
				if obj.AnewInsert {
					if obj.RecLength != 0 || obj.SenLength != 0 {
						obj.ListFilter = Config.Filter.Clone()
						obj.AnewInsert = false
						_obj := Session.Insert{Theology: res.Theology, Method: obj.Method, ClientIP: Conn.LocalAddress(), ProcessName: Conn.GetProcessName(), Time: getTime(), Ico: "websocket_connect"}
						_obj.Host = obj.Host
						_obj.RemoteAddress = obj.RemoteAddress
						_obj.UserName = obj.UserName
						_obj.Filter = obj.ListMatch()
						_obj.Note = obj.GetNote()
						InsertList = append(InsertList, _obj)
					}
				}
			}
		}
		break
	case public.SunnyNetMsgTypeTCPClose:
		{
			//连接断开
			obj := tcpInsert(Conn)
			if obj == nil {
				return
			}
			obj.Lock()
			defer obj.Unlock()
			if noteNew != "" {
				obj.SetNote(noteNew)
			}
			{
				obj.Disconnect = false
				res.RecLength = obj.RecLength
				res.SenLength = obj.SenLength
				res.Ico = "websocket_close"
				res.Code = "已断开"
				res.Method = obj.Method
				stream := &Session.Stream{MessageId: Conn.MessageId(), IsSend: false, IsActiveSend: false, Time: getTime(), Body: Session.ConnectDisconnect, IsClose: true, ListFilter: obj.ListFilter}
				res.Stream = stream
				res.Note = obj.GetNote()
				obj.AddStream(stream)
				if atomic.LoadInt64(&Config.CurrentTheology) == int64(obj.Theology) {
					_stream := stream.ToUpdateStream(obj.Theology, obj.GetStreamFilter())
					lock.Lock()
					updateSocketStreamList = append(updateSocketStreamList, _stream)
					lock.Unlock()
				}
				if obj.AnewInsert {
					if obj.RecLength != 0 || obj.SenLength != 0 {
						obj.ListFilter = Config.Filter.Clone()
						obj.AnewInsert = false
						_obj := Session.Insert{Theology: res.Theology, Method: obj.Method, ClientIP: Conn.LocalAddress(), ProcessName: Conn.GetProcessName(), Time: getTime(), Ico: "websocket_connect"}
						_obj.Host = obj.Host
						_obj.RemoteAddress = obj.RemoteAddress
						_obj.UserName = obj.UserName
						_obj.Filter = obj.ListMatch()
						_obj.Note = obj.GetNote()
						InsertList = append(InsertList, _obj)
					}
				}
			}
		}
		break
	}
	if res.Stream != nil {
		res.Filter = res.Stream.ListIsMatch()
		updateSocketMessage(res)
	}
}
func (g *AppMain) udpCallback(Conn SunnyNet.ConnUDP) {
	Config.ReplaceUDP(Conn)
	res := updateSocket{Code: "已连接"}
	res.Theology = Conn.Theology()
	noteNew := Conn.GetNote()
	var req Session.AppSession
	defer func() {
		if res.Stream != nil {
			res.Filter = res.Stream.ListIsMatch()
			if req != nil {
				if noteNew != "" {
					req.SetNote(noteNew)
				}
				res.Note = req.GetNote()
			}
			updateSocketMessage(res)
		}
	}()
	if Conn.Type() == public.SunnyNetUDPTypeClosed {
		{
			//连接断开
			{
				__obj, ok := Session.Session.Load(Conn.Theology())
				if !ok {
					return
				}
				obj, ok := __obj.(*Session.UDPSession)
				if !ok {
					return
				}
				req = obj
				obj.Lock()
				defer obj.Unlock()
				obj.Disconnect = false
				res.RecLength = obj.RecLength
				res.SenLength = obj.SenLength
				res.Ico = "websocket_close"
				res.Code = "已断开"
				res.Method = obj.Method
				stream := &Session.Stream{MessageId: Conn.MessageId(), IsSend: false, IsActiveSend: false, Time: getTime(), Body: Session.ConnectDisconnect, IsClose: true, ListFilter: obj.ListFilter}
				res.Stream = stream
				obj.AddStream(stream)
				if atomic.LoadInt64(&Config.CurrentTheology) == int64(obj.Theology) {
					_stream := stream.ToUpdateStream(obj.Theology, obj.GetStreamFilter())
					lock.Lock()
					updateSocketStreamList = append(updateSocketStreamList, _stream)
					lock.Unlock()
				}
				if obj.AnewInsert {
					if obj.RecLength != 0 || obj.SenLength != 0 {
						obj.AnewInsert = false
						obj.ListFilter = Config.Filter.Clone()
						_obj := Session.Insert{Theology: res.Theology, Method: obj.Method, ClientIP: Conn.LocalAddress(), ProcessName: Conn.GetProcessName(), Time: getTime(), Ico: "websocket_connect"}
						_obj.Host = obj.RemoteAddress
						_obj.RemoteAddress = _obj.Host
						_obj.ClientIP = obj.ClientIP
						_obj.Time = obj.Time
						_obj.ProcessName = obj.ProcessName
						_obj.Method = obj.Method
						_obj.UserName = obj.UserName
						_obj.Filter = obj.ListMatch()
						InsertList = append(InsertList, _obj)
					}
				}
			}
		}
		return
	}
	if Conn.Type() == public.SunnyNetUDPTypeSend {
		{
			//发送数据
			obj := udpInsert(Conn)
			obj.Lock()
			defer obj.Unlock()
			req = obj
			{
				stream := &Session.Stream{MessageId: Conn.MessageId(), IsSend: true, IsActiveSend: false, Time: getTime(), Body: Conn.Body(), ListFilter: obj.ListFilter}
				res.Stream = stream
				obj.AddStream(stream)
				obj.SenLength += Conn.BodyLen()
				res.RecLength = obj.RecLength
				res.SenLength = obj.SenLength
				res.Method = obj.Method
				if atomic.LoadInt64(&Config.CurrentTheology) == int64(obj.Theology) {
					_stream := stream.ToUpdateStream(obj.Theology, obj.GetStreamFilter())
					lock.Lock()
					updateSocketStreamList = append(updateSocketStreamList, _stream)
					lock.Unlock()
				}

				if obj.AnewInsert {
					if obj.RecLength != 0 || obj.SenLength != 0 {
						obj.AnewInsert = false
						obj.ListFilter = Config.Filter.Clone()
						_obj := Session.Insert{Theology: res.Theology, Method: obj.Method, ClientIP: Conn.LocalAddress(), ProcessName: Conn.GetProcessName(), Time: getTime(), Ico: "websocket_connect"}
						_obj.Host = obj.RemoteAddress
						_obj.RemoteAddress = _obj.Host
						_obj.ClientIP = obj.ClientIP
						_obj.Time = obj.Time
						_obj.ProcessName = obj.ProcessName
						_obj.Method = obj.Method
						_obj.UserName = obj.UserName
						_obj.Filter = obj.ListMatch()
						InsertList = append(InsertList, _obj)
					}
				}
			}
		}
		return
	}
	if Conn.Type() == public.SunnyNetUDPTypeReceive {
		{
			//接收
			obj := udpInsert(Conn)
			obj.Lock()
			defer obj.Unlock()
			req = obj
			{
				stream := &Session.Stream{MessageId: Conn.MessageId(), IsSend: false, IsActiveSend: false, Time: getTime(), Body: Conn.Body(), ListFilter: obj.ListFilter}
				res.Stream = stream
				obj.AddStream(stream)
				obj.RecLength += Conn.BodyLen()
				res.RecLength = obj.RecLength
				res.SenLength = obj.SenLength
				res.Method = obj.Method
				if atomic.LoadInt64(&Config.CurrentTheology) == int64(obj.Theology) {
					_stream := stream.ToUpdateStream(obj.Theology, obj.GetStreamFilter())
					lock.Lock()
					updateSocketStreamList = append(updateSocketStreamList, _stream)
					lock.Unlock()
				}
				if obj.AnewInsert {
					if obj.RecLength != 0 || obj.SenLength != 0 {
						obj.AnewInsert = false
						obj.ListFilter = Config.Filter.Clone()
						_obj := Session.Insert{Theology: res.Theology, Method: obj.Method, ClientIP: Conn.LocalAddress(), ProcessName: Conn.GetProcessName(), Time: getTime(), Ico: "websocket_connect"}
						_obj.Host = obj.RemoteAddress
						_obj.RemoteAddress = _obj.Host
						_obj.ClientIP = obj.ClientIP
						_obj.Time = obj.Time
						_obj.ProcessName = obj.ProcessName
						_obj.Method = obj.Method
						_obj.UserName = obj.UserName
						_obj.Filter = obj.ListMatch()
						InsertList = append(InsertList, _obj)
					}
				}
			}
		}
		return
	}
}

func (g *AppMain) log(SunnyNetContext int, info ...any) {

}
func (g *AppMain) save(SunnyNetContext int, code []byte) {

}
