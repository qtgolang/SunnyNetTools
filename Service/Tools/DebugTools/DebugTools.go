package DebugTools

import (
	"bytes"
	"changeme/Service/Config"
	"changeme/Service/Session"
	"encoding/base64"
	"encoding/json"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/qtgolang/SunnyNet/SunnyNet"
	"github.com/qtgolang/SunnyNet/src/Compress"
	"github.com/qtgolang/SunnyNet/src/SunnyProxy"
	"github.com/qtgolang/SunnyNet/src/encoding/hex"
	"github.com/qtgolang/SunnyNet/src/http"
	"github.com/qtgolang/SunnyNet/src/httpClient"
	"github.com/qtgolang/SunnyNet/src/loop"
	"github.com/qtgolang/SunnyNet/src/public"
)

type DebugTools struct {
	App *SunnyNet.Sunny
	Encryption
	Encoding
	HttpDebug
}
type HttpDebug struct {
	isSelect bool
	App      *SunnyNet.Sunny
}

func (a *HttpDebug) AppStartSelectRequest(isSelect bool) {
	a.isSelect = isSelect
}
func (a *HttpDebug) IsGetSelectRequest() bool {
	return a.isSelect
}

func (a *HttpDebug) SetSelectRequest(req *Session.HttpSession) {
	obj := Config.AppList["Main"]
	if obj != nil {
		for _, v := range req.Request.Body {
			if v < 9 {
				obj.EmitEvent("SetSelectRequest", req.Request.Url, req.Request.Method, req.Request.Header, hex.EncodeToString(req.Request.Body), "hex")
				return
			}
		}
		obj.EmitEvent("SetSelectRequest", req.Request.Url, req.Request.Method, req.Request.Header, req.Request.Body, "text")
	}
}

var decompressors = map[string]func([]byte) []byte{
	"gzip":    Compress.GzipUnCompress,
	"br":      Compress.BrUnCompress,
	"deflate": Compress.DeflateUnCompress,
	"zstd":    Compress.ZSTDDecompress,
	"zlib":    Compress.ZlibUnCompress,
}

func (a *HttpDebug) DoHTTPRequest(req Config.DoHTTPRequestInfo) string {
	var resp Config.DoHTTPResponseInfo
	body, err, errorLevel := req.GetBody()
	if err != nil {
		resp.ErrorLevel = errorLevel
		resp.Error = err.Error()
		i, er := json.Marshal(resp)
		if er == nil {
			return string(i)
		}
		return "未知的内部错误0x002"
	}

	h, _ := http.NewRequest(req.Method, req.Url, bytes.NewReader(body))
	h.Header = req.Header
	RequestProxy, _ := SunnyProxy.ParseProxy(req.ProxyIP, req.OutTime*1000)
	ok, ip := public.IsLocalIP(Config.Config.OutRouter)
	if ok {
		if ip.To4() != nil {
			localAddr, er := net.ResolveTCPAddr("tcp", Config.Config.OutRouter+":0")
			if er == nil {
				h.SetContext(public.OutRouterIPKey, localAddr)
			}

		} else {
			localAddr, er := net.ResolveTCPAddr("tcp", "["+Config.Config.OutRouter+"]:0")
			if er == nil {
				h.SetContext(public.OutRouterIPKey, localAddr)
			}
		}
	}
	errorStr := `

1.请检查网址域名是否正确？

2.请检查网络是否通常？

3.请用浏览器检查是否能访问？

4.请检查 ‘设置->强制TCP设置’ 该域名是否在强制走TCP的规则范围内？
`
	if RequestProxy != nil && RequestProxy.URL != nil {
		if RequestProxy.Port() == strconv.Itoa(Config.Config.Port) {
			bl := false
			if RequestProxy.Host == "127.0.0.1" || RequestProxy.Host == "127.0.0.1:"+strconv.Itoa(Config.Config.Port) {
				bl = true
			} else if RequestProxy.Host == "localhost" || RequestProxy.Host == "localhost:"+strconv.Itoa(Config.Config.Port) {
				bl = true
			} else {
				m1 := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(RequestProxy.Host, "[", ""), "]", ""))
				adders, err := net.InterfaceAddrs()
				if err == nil {
					for _, addr := range adders {
						// 类型断言，排除非 IPNet 类型
						ipNet, ok := addr.(*net.IPNet)
						if !ok || ipNet.IP.IsLoopback() {
							continue
						}
						m2 := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(ipNet.IP.String(), "[", ""), "]", ""))
						m3 := m2 + ":" + strconv.Itoa(Config.Config.Port)
						if m1 == m2 || m1 == m3 {
							bl = true
							break
						}
					}
				}
			}
			if bl && Config.Config.MustTcp.Type == Config.MustTcpTypeAll {
				return "你设置了代理,这个代理是本程序当前的代理端口\r\n\r\n但你开启 全部强制走TCP 选项\r\n\r\n请取消代理设置 或 在设置中 关闭全部强制走TCP"
			}
			if bl && Config.Config.MustTcp.Type != 0 {
				h.Header.Set(public.HTTPClientTags, "true")
			}
		}
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
		RequestProxy:  RequestProxy,
		CheckRedirect: !req.Redirect,
		TLSConfig:     nil,
		OutTime:       time.Duration(req.OutTime) * time.Second,
		GetTLSValues:  public.GetTLSValues,
		MConn:         nil,
		Event: httpClient.Event{
			Connection: _connection,
			Close:      _close,
		},
	}
	r := httpClient.DoOptions(h, op)
	defer func() {
		if r.Close != nil {
			r.Close()
		}
	}()
	if r.Err != nil {
		resp.ErrorLevel = 4
		resp.Error = r.Err.Error() + errorStr
		i, er := json.Marshal(resp)
		if er == nil {
			return string(i)
		}
		return "未知的内部错误0x002" + errorStr
	}
	bs, _ := io.ReadAll(r.Response.Body)
	_bs := bs
	encoding := r.Response.Header.Get("content-encoding")
	if decompressor, ok := decompressors[encoding]; ok {
		if uncompressed := decompressor(bs); len(uncompressed) > 0 {
			_bs = uncompressed
			r.Response.Header.Del("Content-Encoding")
		}
	}
	resp.Body = base64.StdEncoding.EncodeToString(_bs)
	resp.ErrorLevel = 0
	resp.Header = r.Response.Header
	resp.Code = r.Response.StatusCode
	resp.Proto = r.Response.Proto
	resp.Status = r.Response.Status
	i, er := json.Marshal(resp)
	if er != nil {
		return "未知的错误0x003"
	}
	return string(i)
}
func (a *HttpDebug) DoHTTPGenerateCode(req Config.DoHTTPRequestInfo) string {
	var rs Session.HttpSession
	rs.Request.Body, _, _ = req.GetBody()
	rs.Request.Header = req.Header
	rs.Request.Method = req.Method
	rs.Request.Url = req.Url
	if Session.CreateRequestCode(&rs, req.Language, req.Type) != nil {
		return ""
	}
	return "ok"
}
func extractPorts(conn net.Conn) (uint16, uint16) {
	tcpLocal, ok1 := conn.LocalAddr().(*net.TCPAddr)   // 本地 TCP 地址
	tcpRemote, ok2 := conn.RemoteAddr().(*net.TCPAddr) // 对端 TCP 地址
	if !ok1 || !ok2 {
		return 0, 0
	}
	return uint16(tcpLocal.Port), uint16(tcpRemote.Port)
}
