package Service

import (
	. "changeme/Service/Config"
	"changeme/Service/Session"
	"changeme/Service/clipboard"
	"changeme/Service/update"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"sync/atomic"

	"github.com/qtgolang/SunnyNet/SunnyNet"
	"github.com/qtgolang/SunnyNet/src/encoding/hex"
	"github.com/qtgolang/SunnyNet/src/public"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// ClipboardWriteAll 设置剪辑版内容,成功返回空字符串
func (g *AppMain) ClipboardWriteAll(value string) string {
	_ = clipboard.ClipboardWriteAll(value)
	return ""
}

// SetWorking 设置是否隐藏捕获数据
func (g *AppMain) SetWorking(working bool) {
	if working {
		Config.IsHideHook = 1
	} else {
		Config.IsHideHook = 0
	}
}

// SetBreakMode 设置断点模式
func (g *AppMain) SetBreakMode(working uint32) {
	Config.BreakMode = working
}

const MaxBodyLength = 10240 * 4

// GetHTTPSession 根据唯一ID 获取HTTP消息内容
func (g *AppMain) GetHTTPSession(Theology int) *Session.HttpSession {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	req := Session.GetHttpSession(Theology)
	lock.Lock()
	if req != nil {
		req.Request.BodyLength = len(req.Request.Body)
		req.Request.IsMaxLength = req.Request.BodyLength > MaxBodyLength

		req.Response.BodyLength = len(req.Response.Body)
		req.Response.IsMaxLength = req.Response.BodyLength > MaxBodyLength
		if g.IsGetSelectRequest() {
			go g.SetSelectRequest(req)
		}
	}
	SetCurrentTheology(Theology)
	lock.Unlock()

	return req
}
func (g *AppMain) GetHTTPResponseBody(Theology int, getAllBody bool) []byte {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	req := Session.GetHttpSession(Theology)
	lock.Lock()
	if req != nil {
		lock.Unlock()
		if !getAllBody {
			if len(req.Response.Body) > MaxBodyLength {
				return req.Response.Body[:MaxBodyLength]
			}
		}
		return req.Response.Body
	}
	lock.Unlock()
	return nil
}
func (g *AppMain) GetHTTPRequestBody(Theology int, getAllBody bool) []byte {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	req := Session.GetHttpSession(Theology)
	lock.Lock()
	if req != nil {
		lock.Unlock()
		if !getAllBody {
			if len(req.Request.Body) > MaxBodyLength {
				return req.Request.Body[:MaxBodyLength]
			}
		}
		return req.Request.Body
	}
	lock.Unlock()
	return nil
}

// SessionActiveSend 主动发送
func (g *AppMain) SessionActiveSend(Theology int, isSendServer bool, SendType string, wsType int, _bs []byte) string {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()

	var bs []byte
	_SendType := strings.ToLower(SendType)
	switch _SendType {
	case "hex":
		__bs, e := hex.DecodeString(strings.ReplaceAll(strings.TrimSpace(string(_bs)), " ", ""))
		if e != nil {
			return "尝试十六进制解码失败：" + e.Error()
		}
		bs = __bs
		break
	case "base64":
		__bs, e := base64.StdEncoding.DecodeString(strings.ReplaceAll(strings.TrimSpace(string(_bs)), " ", ""))
		if e != nil {
			return "尝试 Base64 解码失败：" + e.Error()
		}
		bs = __bs
		break
	case "gbk":
		reader := transform.NewReader(strings.NewReader(string(_bs)), simplifiedchinese.GBK.NewEncoder())
		__bs, e := ioutil.ReadAll(reader) // 读取转换后的 GBK 字节流
		if e != nil {
			return "尝试 将字符串转到 GBK 编码失败：" + e.Error()
		}
		bs = __bs
		break
	default:
		bs = _bs
		break

	}
	obj := Session.GetAppSession(Theology)
	if obj != nil {
		obj.Lock()
		defer obj.Unlock()
		err := obj.ActiveSend(isSendServer, wsType, bs)
		if err != "" {
			return err
		}
		stream := &Session.Stream{MessageId: SunnyNet.NewMessageId(), Time: getTime(), Body: bs, IsSend: isSendServer, IsActiveSend: true}
		obj.AddStream(stream, wsType)
		if atomic.LoadInt64(&Config.CurrentTheology) == int64(obj.GetTheology()) {
			_stream := stream.ToUpdateStream(obj.GetTheology(), obj.GetStreamFilter())
			lock.Lock()
			updateSocketStreamList = append(updateSocketStreamList, _stream)
			lock.Unlock()
		}
	}
	return ""
}

// URLQueryUnescape URL解码,失败返回空
func (g *AppMain) URLQueryUnescape(value string) []byte {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	a, _ := url.QueryUnescape(value)
	return []byte(a)
}

// URLQueryEscape URL编码,失败返回空
func (g *AppMain) URLQueryEscape(value string) string {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	a := url.QueryEscape(value)
	return a
}

// SetRequestNextBreakMode 设置断点模式
func (g *AppMain) SetRequestNextBreakMode(Theology int, working uint32) {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	obj := Session.GetHttpSession(Theology)
	if obj != nil {
		obj.NextBreakMode = working
		obj.Wg.Done()
	}
}

// FreeAllRequest 释放全部
func (g *AppMain) FreeAllRequest() {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	Session.Session.Range(func(_, value any) bool {
		obj, ok := value.(*Session.HttpSession)
		if ok {
			obj.NextBreakMode = 0
			obj.Wg.Done()
		}
		return true
	})
}

// UpdateHttpRequest 更新 HTTP 请求（MCP / UI 共用）。
func (g *AppMain) UpdateHttpRequest(Theology int, req *Session.HttpSessionRequest) {
	g.updateHttpRequest(Theology, req)
}

// UpdateHttpRequest 更新HTTP请求信息
func (g *AppMain) updateHttpRequest(Theology int, req *Session.HttpSessionRequest) {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	obj := Session.GetHttpSession(Theology)
	if obj == nil {
		return
	}
	waiting, phase := httpBreakPhase(obj)
	if !waiting || phase != breakPhaseUpstream {
		return
	}
	lock.Lock()
	obj.Request.Url = req.Url
	// 拦截上行：Method 不可修改，不应用 req.Method
	obj.Request.Header = req.Header
	obj.Request.Body = req.Body
	obj.Request.IsUpdateBody = true
	lock.Unlock()
	g.queueHttpRequestRowUpdate(obj)
}

// queueHttpRequestRowUpdate 断点改请求后刷新主列表（URL 等），不触发放行。
func (g *AppMain) queueHttpRequestRowUpdate(obj *Session.HttpSession) {
	if obj == nil {
		return
	}
	lock.Lock()
	defer lock.Unlock()
	bm := uint32(BreakNone)
	ico := obj.Ico
	if obj.IsWait() {
		switch obj.State {
		case public.HttpSendRequest:
			bm = BreakSend
			ico = "拦截上行"
		case public.HttpResponseOK:
			bm = BreakRece
			ico = "拦截下行"
		default:
			bm = BreakSend
			ico = "拦截上行"
		}
	}
	filter := obj.ListMatch()
	if obj.IsWait() {
		filter = true
	}
	row := httpUpdateSend{
		Theology:  obj.Theology,
		Ico:       ico,
		Filter:    filter,
		URL:       obj.Request.Url,
		Method:    obj.Request.Method,
		Note:      obj.GetNote(),
		BreakMode: bm,
	}
	for i, v := range updateSendList {
		if v.Theology == row.Theology {
			updateSendList[i] = row
			g.emitHttpRequestRowUpdate(row)
			return
		}
	}
	updateSendList = append(updateSendList, row)
	g.emitHttpRequestRowUpdate(row)
}

// emitHttpRequestRowUpdate 立即通知主窗口刷新请求行（不依赖批处理轮询）。
func (g *AppMain) emitHttpRequestRowUpdate(row httpUpdateSend) {
	if AppList["Main"] != nil {
		AppList["Main"].EmitEvent("updateSendHTTP", []httpUpdateSend{row})
	}
	emitMCPMainJSON("rowupdatesend", map[string]any{
		"theology":  row.Theology,
		"url":       row.URL,
		"method":    row.Method,
		"ico":       row.Ico,
		"note":      row.Note,
		"breakMode": row.BreakMode,
		"filter":    row.Filter,
	})
}

// UpdateHttpResponse 更新 HTTP 响应（MCP / UI 共用）。
func (g *AppMain) UpdateHttpResponse(Theology int, req *Session.HttpSessionResponse) {
	g.updateHttpResponse(Theology, req)
}

// UpdateHttpResponse 更新HTTP响应信息
func (g *AppMain) updateHttpResponse(Theology int, req *Session.HttpSessionResponse) {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	obj := Session.GetHttpSession(Theology)
	if obj == nil {
		return
	}
	waiting, phase := httpBreakPhase(obj)
	if !waiting || phase != breakPhaseDownstream {
		return
	}
	lock.Lock()
	obj.Response.Code = req.Code
	obj.Response.State = req.State
	obj.Response.Header = req.Header
	obj.Response.Body = req.Body
	obj.Response.IsUpdateBody = true
	lock.Unlock()
	g.queueHttpResponseRowUpdate(obj)
}

// queueHttpResponseRowUpdate 断点改响应后刷新主列表状态列等，不触发放行。
func (g *AppMain) queueHttpResponseRowUpdate(obj *Session.HttpSession) {
	if obj == nil {
		return
	}
	lock.Lock()
	defer lock.Unlock()
	bm := uint32(BreakNone)
	ico := obj.Ico
	if obj.IsWait() {
		switch obj.State {
		case public.HttpResponseOK:
			bm = BreakRece
			ico = "拦截下行"
		case public.HttpSendRequest:
			bm = BreakSend
			ico = "拦截上行"
		default:
			bm = BreakRece
			ico = "拦截下行"
		}
	}
	filter := obj.ListMatch()
	if obj.IsWait() {
		filter = true
	}
	row := updateHTTPDone{
		Theology:  obj.Theology,
		Code:      obj.Response.Code,
		Length:    obj.Response.Length,
		Type:      obj.Response.Type,
		Time:      obj.Response.Time,
		Ico:       ico,
		BreakMode: bm,
		Filter:    filter,
		IP:        obj.Response.ServerIP,
		Note:      obj.GetNote(),
	}
	for i, v := range updateDoneList {
		if v.Theology == row.Theology {
			updateDoneList[i] = row
			return
		}
	}
	updateDoneList = append(updateDoneList, row)
}

// AppDeleteSession 删除指定会话
func (g *AppMain) AppDeleteSession(tid []int) {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	lock.Lock()
	var remove []int
	for _, v := range tid {
		obj := Session.GetAppSession(v)
		if obj != nil {
			if obj.IsHTTP() && !obj.IsWebsocket() {
				remove = append(remove, v)
				continue
			}
			obj.ClearAllSession()
			if obj.IsDisconnect() {
				remove = append(remove, v)
			} else {
				obj.SetAnewInsert(true)
			}
		}
	}
	for _, v := range remove {
		Session.Session.Delete(v)
	}
	lock.Unlock()
	return
}

// ClearAllSession 删除全部
func (g *AppMain) ClearAllSession() []int {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	lock.Lock()
	var array []int
	var remove []int
	Session.Session.Range(func(key, value any) bool {
		theology, ok := key.(int)
		if !ok {
			return true
		}
		array = append(array, theology)
		{
			obj, ok := value.(Session.AppSession)
			if ok {
				if obj.IsHTTP() && !obj.IsWebsocket() {
					remove = append(remove, theology)
					return true
				}
				obj.ClearAllSession()
				if obj.IsDisconnect() {
					remove = append(remove, theology)
				} else {
					obj.SetAnewInsert(true)
				}
			}
		}
		return true
	})
	for _, v := range remove {
		Session.Session.Delete(v)
	}
	lock.Unlock()
	return array
}

// GetSessionMessageBody 获取指定 TCP/UDP/WSS 字节数组 数据
func (g *AppMain) GetSessionMessageBody(Theology, MessageId int) []byte {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	lock.Lock()
	bs := make([]byte, 0)
	obj := Session.GetAppSession(Theology)
	if obj != nil {
		stream := obj.GetStream(MessageId)
		if stream != nil {
			bs = stream.GetBody()
		}
	}
	lock.Unlock()
	return bs
}

// DelSessionMessageIdArray 删除指定 TCP/UDP/WSS 数据
func (g *AppMain) DelSessionMessageIdArray(theology int, MessageIdArray []int) []int {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	lock.Lock()
	var array []int
	obj := Session.GetAppSession(theology)
	if obj != nil {
		array = obj.DeleteMessageIdArray(MessageIdArray)
		objs := updateSocket{}
		objs.Theology = theology
		objs.Ico = "updateLen"
		objs.RecLength = obj.GetRecLength()
		objs.SenLength = obj.GetSenLength()
		objs.Method = obj.GetMethod()
		updateSocket_List = append(updateSocket_List, objs)
	}
	lock.Unlock()
	return array
}

// ClearAllSessionMessageIdArray 删除指定唯一ID下的全部数据流
func (g *AppMain) ClearAllSessionMessageIdArray(theology int) []int {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	lock.Lock()
	var array []int
	obj := Session.GetAppSession(theology)
	if obj != nil {
		array = obj.ClearAllSession()
		objs := updateSocket{}
		objs.Theology = theology
		objs.Ico = "updateLen"
		objs.RecLength = obj.GetRecLength()
		objs.SenLength = obj.GetSenLength()
		objs.Method = obj.GetMethod()
		updateSocket_List = append(updateSocket_List, objs)
	}
	lock.Unlock()
	return array
}

// CopySessionMessageIdArray 复制数据 CopyType 复制类型 MessageIdArray 要复制是消息ID
func (g *AppMain) CopySessionMessageIdArray(theology int, CopyType string, MessageIdArray []int) string {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	lock.Lock()
	var sb strings.Builder
	obj := Session.GetAppSession(theology)
	if obj != nil {
		switch CopyType {
		case "copy":
			for _, v := range MessageIdArray {
				m := obj.GetStream(v)
				if m.GetIsSend() {
					sb.WriteString("发送数据:" + Session.GetHexAllSpaces(m.GetBody()) + "\r\n")
				} else {
					sb.WriteString("接收数据:" + Session.GetHexAllSpaces(m.GetBody()) + "\r\n")
				}
			}
			break
		case "all", "send", "rec":
			{
				obj.RangeStream(func(stream Session.AppStream) bool {
					isSend := stream.GetIsSend()
					body := Session.GetHexAllSpaces(stream.GetBody())
					switch CopyType {
					case "all": // 发送和接收都记录
						if isSend {
							sb.WriteString("发送数据:" + body + "\r\n")
						} else {
							sb.WriteString("接收数据:" + body + "\r\n")
						}
					case "send": // 只记录发送
						if isSend {
							sb.WriteString("发送数据:" + body + "\r\n")
						}
					case "rec": // 只记录接收
						if !isSend {
							sb.WriteString("接收数据:" + body + "\r\n")
						}
					}
					return true
				})
			}
			break
		}
	}
	lock.Unlock()
	return g.ClipboardWriteAll(sb.String())
}

// StreamSearch 设置 WS/UDP/TCP 过滤器
func (g *AppMain) StreamSearch(Theology int, FilterJson string) []int {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	res := make([]int, 0)
	stream := Session.GetAppSession(Theology)
	if stream == nil {
		return res
	}
	stream.SetStreamFilter(Session.ParseFilter(FilterJson))
	if stream.GetStreamFilter() == nil {
		return res
	}
	lock.Lock()
	stream.RangeStream(func(val Session.AppStream) bool {
		if val.MatchFilters(stream.GetStreamFilter()) {
			res = append(res, val.GetMessageId())
		}
		return true
	})
	lock.Unlock()
	SetCurrentTheology(Theology)
	return res
}

// GetSocketFilter 获取  WS/UDP/TCP 过滤器
func (g *AppMain) GetSocketFilter(theology int) string {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	if value, ok := Session.Session.Load(theology); ok {
		if obj, ok1 := value.(Session.AppSession); ok1 {
			if f := obj.GetStreamFilter(); f != nil {
				return f.String
			}
		}
	}
	return ""
}

// ListSearch 设置主列表过滤器
func (g *AppMain) ListSearch(FilterJson string) []int {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	lock.Lock()
	defer lock.Unlock()
	Filter := Session.ParseFilter(FilterJson)
	Config.Filter = Filter
	var array []int
	if Filter != nil {
		Session.Session.Range(func(key, value any) bool {
			theology, ok := key.(int)
			if !ok {
				return true
			}
			obj := value.(Session.AppSession)
			if obj != nil {
				{
					//当前是否正在断点拦截,仅HTTP请求有效， tcp/dup 始终返回false
					//如果正在拦截 那么始终显示到列表区,不执行过滤操作
					if obj.IsWait() {
						array = append(array, theology)
						return true
					}
				}
				obj.SetListFilter(Filter.Clone())
				if obj.ListMatch() {
					array = append(array, theology)
				}
			}
			return true
		})
	} else {
		Session.Session.Range(func(key, value any) bool {
			obj := value.(Session.AppSession)
			if obj != nil {
				obj.SetListFilter(nil)
			}
			return true
		})
	}
	return array
}

// UpdateNote 更新注释信息
func (g *AppMain) UpdateNote(Theology int, Note string) bool {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	obj := Session.GetAppSession(Theology)
	if obj == nil {
		return false
	}
	obj.SetNote(Note)
	if _Filter := obj.GetListFilter(); _Filter != nil {
		obj.SetStreamFilter(Session.ParseFilter(_Filter.String))
		return obj.ListMatch()
	}
	return true
}

// AppSaveRequestImg 保存请求图片
func (g *AppMain) AppSaveRequestImg(Theology int, ImgType string, IsRequest bool, path string) string {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	obj := Session.GetHttpSession(Theology)
	if obj == nil {
		return "没有在缓存中找到这个请求"
	}
	_ = os.Remove(path)
	if IsRequest {
		if os.WriteFile(path, obj.Request.Body, 0777) != nil {
			return "请检查文件目录是否可写入"
		}
	} else {
		if os.WriteFile(path, obj.Response.Body, 0777) != nil {
			return "请检查文件目录是否可写入"
		}
	}
	return ""
}

// AppVersion 获取SDK版本信息
func (g *AppMain) AppVersion() string {
	return public.SunnyVersion + " " + os.Getenv("minor_version")
}

func init() {
	update.SetNewVerCallback(func(newUrl string) {
		fmt.Println("发现新版本", newUrl)
	})
}
