package Config

import (
	"bytes"
	"changeme/Service/Session"
	Theme2 "changeme/Service/Theme"
	_ "embed"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/qtgolang/SunnyNet/SunnyNet"
	"github.com/qtgolang/SunnyNet/src/Compress"
	"github.com/qtgolang/SunnyNet/src/public"
)

//go:embed codeTemplate.txt
var codeTemplate string

const (
	BreakNone = 0 //不拦截
	BreakSend = 1 //拦截上行
	BreakRece = 2 //拦截下行
)

var CertID = 0

type config struct {
	IsDark                bool
	AgGridLightTheme      string
	AgGridDarkTheme       string
	Port                  int
	IEProxy               bool
	DisableTCP            bool
	DisableUDP            bool
	DisableCache          bool
	OpenAuthMode          bool
	LimitRequestSize      int
	SendIsHTTP1           bool
	Tour                  map[string]bool //是否已经引导过了
	Keys                  string
	ScriptCode            string
	ColumnState           string //列状态
	AuthMap               map[int]*Auth
	OutRouter             string `json:"-"`
	HomeTextMark          string
	GenerateCodeInterface string
	EditorFontSize        int
	Filter                *Session.Filter `json:"-"`
	IsHideHook            uint32          `json:"-"` //检查是否隐藏捕获 0 = false, 1 = true
	BreakMode             uint32          `json:"-"` //拦截状态
	CurrentTheology       int64           `json:"-"` //当前选中的唯一ID
	ListColor             map[string]string
	RequestCert           map[int]*CertInfo
	ReplaceRoles          map[int]*ReplaceBodyInfo
	ReplaceHost           map[int]*ReplaceHostInfo
	ProxyWay              map[int]*ProxyWayInfo
	Authentication        map[int]*AuthenticationInfo
	ToolsList             []ToolsInfo
	ProxyRoles            string
	ProxyDns              string
	MustTcp               struct {
		Type  MustTcpType
		Roles string
	}
	HTTPSProto string
	RandomJa3  bool
}
type AuthenticationInfo struct {
	ID   int
	User string
	Pass string
}

type ToolsInfo struct {
	ID   string
	Name string
	Icon string
	File string
	Args string
}

var Config = config{
	ListColor:      make(map[string]string),
	RequestCert:    make(map[int]*CertInfo),
	ReplaceRoles:   make(map[int]*ReplaceBodyInfo),
	ReplaceHost:    make(map[int]*ReplaceHostInfo),
	ProxyWay:       make(map[int]*ProxyWayInfo),
	AuthMap:        make(map[int]*Auth),
	Authentication: make(map[int]*AuthenticationInfo),
}

type Auth struct {
	User string
	Pass string
}

func init() {
	Config = config{
		ListColor:      make(map[string]string),
		RequestCert:    make(map[int]*CertInfo),
		ReplaceRoles:   make(map[int]*ReplaceBodyInfo),
		ReplaceHost:    make(map[int]*ReplaceHostInfo),
		ProxyWay:       make(map[int]*ProxyWayInfo),
		AuthMap:        make(map[int]*Auth),
		Authentication: make(map[int]*AuthenticationInfo),
		Tour:           make(map[string]bool),
	}
	Config.IsDark = true
	Config.AgGridLightTheme = Theme2.ThemeLight1
	Config.AgGridDarkTheme = Theme2.ThemeDark1
	Config.Load()
	initListColor()
}

// DefaultColor 默认列表配色 暗色和亮色
var DefaultColor = make(map[string]string)

func initListColor() {
	_DefaultColor := `
	{
	"d1":	"#14BEF2",
	"d2":	"#14BEF2",
	"d3":	"#14BEF2",
	"d4":	"#E6ABFF",
	"d5":	"#AFF86F",
	"d6":	"#9A9A9A",
	"d7":	"#10D090",
	"d8":	"#FE6363",
	"d9":	"#FFBF00",
	"d10":	"#FF0000",
	"d99":	"#D6D6D6",
	"l1":	"#FC2CD6",
	"l10":	"#AA0000",
	"l2":	"#FC2CD6",
	"l3":	"#FC2CD6",
	"l4":	"#AF0ABE",
	"l5":	"#0F8C62",
	"l6":	"#B3B1B1",
	"l7":	"#0476EF",
	"l8":	"#AA0000",
	"l9":	"#9FB416",
	"l99":	"#0A0A0A"
	}`
	_ = json.Unmarshal([]byte(_DefaultColor), &DefaultColor)
	_ = json.Unmarshal([]byte(_DefaultColor), &Config.ListColor)
	if Config.ListColor["d1"] != "" {
		return
	}
	//暗色
	Config.ListColor["d1"] = `#FF00000`  //TCP请求          如果是 TCP/TLS-TCP/TCP-Must/ 请求
	Config.ListColor["d2"] = `#FF00001`  //UDP请求          如果是UDP请求
	Config.ListColor["d3"] = `#FF00002`  //Websocket       如果是WS/WSS请求
	Config.ListColor["d4"] = `#FF00003`  //CSS             如果响应是 CSS 文件
	Config.ListColor["d5"] = `#FF00004`  //javaScript      如果响应是 js 文件
	Config.ListColor["d6"] = `#FF00005`  //图片             如果响应是 图片 文件
	Config.ListColor["d7"] = `#FF00006`  //文档             如果响应是 TXT/HTML/...文档类型
	Config.ListColor["d8"] = `#FF00007`  //错误请求          如果响应状态码是 -1 的请求
	Config.ListColor["d9"] = `#FF00008`  //重定向请求        如果响应状态码是 301/302 的重定向请求
	Config.ListColor["d10"] = `#FF00009` //40x             如果响应状态码是 401/403/404/500 的请求
	Config.ListColor["d99"] = `#FF000aa` //普通项           如果不是上列的任何一个
	//===============================================
	//亮色
	Config.ListColor["l1"] = `#AA00000`  //TCP请求          如果是 TCP/TLS-TCP/TCP-Must/ 请求
	Config.ListColor["l2"] = `#AA00001`  //UDP请求          如果是UDP请求
	Config.ListColor["l3"] = `#AA00002`  //Websocket       如果是WS/WSS请求
	Config.ListColor["l4"] = `#AA00003`  //CSS             如果响应是 CSS 文件
	Config.ListColor["l5"] = `#AA00004`  //javaScript      如果响应是 js 文件
	Config.ListColor["l6"] = `#AA00005`  //图片             如果响应是 图片 文件
	Config.ListColor["l7"] = `#AA00006`  //文档             如果响应是 TXT/HTML/...文档类型
	Config.ListColor["l8"] = `#AA00007`  //错误请求          如果响应状态码是 -1 的请求
	Config.ListColor["l9"] = `#AA00008`  //重定向请求        如果响应状态码是 301/302 的重定向请求
	Config.ListColor["l10"] = `#AA00009` //40x            如果响应状态码是 401/403/404/500 的请求
	Config.ListColor["l99"] = `#AA000aa` //普通项           如果不是上列的任何一个

}
func (f *config) Save() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}
	t := f.Tour[public.SunnyVersion]
	f.Tour = make(map[string]bool)
	f.Tour[public.SunnyVersion] = t
	_config, _ := json.Marshal(f)
	_ = os.WriteFile(filepath.Join(homeDir, "SunnyNet.json"), _config, 0777)
}
func (f *config) Reset() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	r := filepath.Join(homeDir, "SunnyNet.json")
	_ = os.Remove(r)
	_ = os.WriteFile(filepath.Join(homeDir, "SunnyNet.json"), []byte("{}"), 0777)
	Config = config{
		ListColor:    make(map[string]string),
		RequestCert:  make(map[int]*CertInfo),
		ReplaceRoles: make(map[int]*ReplaceBodyInfo),
		ReplaceHost:  make(map[int]*ReplaceHostInfo),
		ProxyWay:     make(map[int]*ProxyWayInfo),
		AuthMap:      make(map[int]*Auth),
	}
	Config.IsDark = true
	Config.AgGridLightTheme = Theme2.ThemeLight1
	Config.AgGridDarkTheme = Theme2.ThemeDark1
	Config.Load()
	initListColor()
}

func (f *config) Load() {
	{
		homeDir, err := os.UserHomeDir()
		if err == nil {
			bs, _ := os.ReadFile(filepath.Join(homeDir, "SunnyNet.json"))
			_ = json.Unmarshal(bs, &Config)
		}
	}
	f.initAuthentication()
	f.initBaseSettings()
	f.initRequestCert()
	f.initReplaceBody()
	f.initReplaceHost()
	f.initHTTPSProto()
	f.initMustTcp()
	f.initProxy()
	f.InitCodeTemplate()
	f.IEProxy = false
}
func (f *config) InitCodeTemplate() {
	if len(f.GenerateCodeInterface) < 5 {
		f.GenerateCodeInterface = codeTemplate
	}
}

var decompressors = map[string]func([]byte) []byte{
	"gzip":    Compress.GzipUnCompress,
	"br":      Compress.BrUnCompress,
	"deflate": Compress.DeflateUnCompress,
	"zstd":    Compress.ZSTDDecompress,
	"zlib":    Compress.ZlibUnCompress,
}

func (f *config) ReplaceHttp(Conn SunnyNet.ConnHTTP) bool {
	isBreak := false
	if Conn.Type() == public.HttpSendRequest {
		if f.SendIsHTTP1 {
			Conn.SetHTTP2Config("http/1.1")
		} else {
			if f.HTTPSProto != "" {
				Conn.SetHTTP2Config(f.HTTPSProto)
			}
		}
		{
			u, e := url.Parse(Conn.URL())
			if e == nil && u != nil {
				for _, v := range f.ReplaceHost {
					if v.IsInvalid() {
						continue
					}
					if u.Host == v.LodInfo.Host {
						if v.NewInfo.Port != 0 {
							u.Host = v.NewInfo.Host + ":" + strconv.Itoa(int(v.NewInfo.Port))
						} else {
							u.Host = v.NewInfo.Host
						}
						Conn.UpdateURL(u.String())
						continue
					}
					h, p, _ := net.SplitHostPort(u.Host)
					if h != "" && h == v.LodInfo.Host && v.NewInfo.Port == 0 {
						u.Host = v.NewInfo.Host + ":" + p
						Conn.UpdateURL(u.String())
						continue
					}
				}
			}
		}
		_url := []byte(Conn.URL())
		_header := Conn.GetRequestHeader()
		ok := false
		bs := Conn.GetRequestBody()
		ok2 := false
		for _, v := range f.ReplaceRoles {
			if !v.Ok || v.State == "已禁用" {
				continue
			}
			if v.IsBreak() && !isBreak {
				if v.SourceType_IsHTTP_URL() {
					if bytes.Contains(_url, v.BytesLod) {
						isBreak = true
					}
				}
				if v.SourceType_IsHTTP_RequestHeader() {
					for k, v1 := range _header {
						for _, v11 := range v1 {
							if bytes.Contains([]byte(v11), v.BytesLod) {
								isBreak = true
							}
						}
						if bytes.Contains([]byte(k), v.BytesLod) {
							isBreak = true
						}
					}
				}
				if v.SourceType_IsHTTP_RequestBody() {
					if bytes.Contains(bs, v.BytesLod) {
						isBreak = true
					}
				}
				continue
			}
			if v.SourceType_IsHTTP_URL() {
				if bytes.Contains(_url, v.BytesLod) {
					_url = bytes.ReplaceAll(_url, v.BytesLod, v.BytesNew)
					ok = true
				}
			}
			if v.SourceType_IsHTTP_RequestHeader() {
				for k, v1 := range _header {
					var v12 []string
					for _, v11 := range v1 {
						if bytes.Contains([]byte(v11), v.BytesLod) {
							v12 = append(v12, strings.ReplaceAll(v11, string(v.BytesLod), string(v.BytesNew)))
						} else {
							v12 = append(v12, v11)
						}
					}
					if bytes.Contains([]byte(k), v.BytesLod) {
						delete(_header, k)
						_header[k] = v12
					} else {
						_header[k] = v12
					}
				}
			}
			if v.SourceType_IsHTTP_RequestBody() {
				if bytes.Contains(bs, v.BytesLod) {
					bs = bytes.ReplaceAll(bs, v.BytesLod, v.BytesNew)
					ok2 = true
				}
			}
		}
		if ok {
			Conn.UpdateURL(string(_url))
		}
		if ok2 {
			Conn.SetResponseBody(bs)
		}
		return isBreak
	}
	if Conn.Type() == public.HttpResponseOK {
		_header := Conn.GetResponseHeader()
		if f.DisableCache {
			_header.Del("Cache-Control")
			_header.Del("Pragma")
			_header.Del("Expires")
		}
		bs := Conn.GetResponseBody()
		encoding := strings.ToLower(_header.Get("Content-Encoding"))
		if decompressor, ok := decompressors[encoding]; ok {
			bb := decompressor(bs)
			if len(bb) > 0 {
				_header.Del("Content-Encoding")
				bs = bb
			}
		}
		for _, v := range f.ReplaceRoles {
			if !v.Ok || v.State == "已禁用" {
				continue
			}
			// 判断是否需要拦截
			if v.IsBreak() && !isBreak {
				if v.SourceType_IsHTTP_ResponseHeader() {
					for k, v1 := range _header {
						for _, v11 := range v1 {
							if bytes.Contains([]byte(v11), v.BytesLod) {
								isBreak = true
							}
						}
						if bytes.Contains([]byte(k), v.BytesLod) {
							isBreak = true
						}
					}
				}
				if v.SourceType_IsHTTP_ResponseBody() {
					if bytes.Contains(bs, v.BytesLod) {
						isBreak = true
					}
				}
				continue
			}
			//替换数据
			if v.SourceType_IsHTTP_ResponseHeader() {
				for k, v1 := range _header {
					var v12 []string
					for _, v11 := range v1 {
						if bytes.Contains([]byte(v11), v.BytesLod) {
							v12 = append(v12, strings.ReplaceAll(v11, string(v.BytesLod), string(v.BytesNew)))
						} else {
							v12 = append(v12, v11)
						}
					}
					if bytes.Contains([]byte(k), v.BytesLod) {
						delete(_header, k)
						_header[k] = v12
					} else {
						_header[k] = v12
					}
				}
			}
			if v.SourceType_IsHTTP_ResponseBody() {
				if bytes.Contains(bs, v.BytesLod) {
					bs = bytes.ReplaceAll(bs, v.BytesLod, v.BytesNew)
				}
			}
		}
		Conn.SetResponseBody(bs)
	}
	return isBreak
}
func (f *config) ReplaceTCP(Conn SunnyNet.ConnTCP) {
	if Conn.Type() == public.SunnyNetMsgTypeTCPAboutToConnect {
		u := Conn.RemoteAddress()
		for _, v := range f.ReplaceHost {
			if v.IsInvalid() {
				continue
			}
			if u == v.LodInfo.Host {
				if v.NewInfo.Port != 0 {
					Conn.SetNewAddress(v.NewInfo.Host + ":" + strconv.Itoa(int(v.NewInfo.Port)))
				} else {
					Conn.SetNewAddress(v.NewInfo.Host)
				}
				continue
			}

			h, p, _ := net.SplitHostPort(u)
			if h != "" && h == v.LodInfo.Host && v.NewInfo.Port == 0 {
				Conn.SetNewAddress(v.NewInfo.Host + ":" + p)
				continue
			}
		}
	}
	if Conn.Type() == public.SunnyNetMsgTypeTCPClientSend {
		bs := Conn.Body()
		for _, v := range f.ReplaceRoles {
			if !v.Ok || v.State == "已禁用" {
				continue
			}
			if v.SourceType_IsTCP_Send() {
				if bytes.Contains(bs, v.BytesLod) {
					bs = bytes.ReplaceAll(bs, v.BytesLod, v.BytesNew)
				}
			}
		}
		Conn.SetBody(bs)
	}
	if Conn.Type() == public.SunnyNetMsgTypeTCPClientReceive {
		bs := Conn.Body()
		for _, v := range f.ReplaceRoles {
			if !v.Ok || v.State == "已禁用" {
				continue
			}
			if v.SourceType_IsTCP_Receive() {
				if bytes.Contains(bs, v.BytesLod) {
					bs = bytes.ReplaceAll(bs, v.BytesLod, v.BytesNew)
				}
			}
		}
		Conn.SetBody(bs)
	}
}
func (f *config) ReplaceUDP(Conn SunnyNet.ConnUDP) {
	if Conn.Type() == public.SunnyNetUDPTypeSend {
		bs := Conn.Body()
		for _, v := range f.ReplaceRoles {
			if !v.Ok || v.State == "已禁用" {
				continue
			}
			if v.SourceType_IsUDP_Send() {
				if bytes.Contains(bs, v.BytesLod) {
					bs = bytes.ReplaceAll(bs, v.BytesLod, v.BytesNew)
				}
			}
		}
		Conn.SetBody(bs)
	}
	if Conn.Type() == public.SunnyNetUDPTypeReceive {
		bs := Conn.Body()
		for _, v := range f.ReplaceRoles {
			if !v.Ok || v.State == "已禁用" {
				continue
			}
			if v.SourceType_IsUDP_Receive() {
				if bytes.Contains(bs, v.BytesLod) {
					bs = bytes.ReplaceAll(bs, v.BytesLod, v.BytesNew)
				}
			}
		}
		Conn.SetBody(bs)
	}
}
func (f *config) ReplaceWebsocket(Conn SunnyNet.ConnWebSocket) {
	if Conn.Type() == public.WebsocketUserSend {
		bs := Conn.Body()
		for _, v := range f.ReplaceRoles {
			if !v.Ok || v.State == "已禁用" {
				continue
			}
			if v.SourceType_IsWebsocket_Send() {
				if bytes.Contains(bs, v.BytesLod) {
					bs = bytes.ReplaceAll(bs, v.BytesLod, v.BytesNew)
				}
			}
		}
		Conn.SetBody(bs)
	}
	if Conn.Type() == public.WebsocketServerSend {
		bs := Conn.Body()
		for _, v := range f.ReplaceRoles {
			if !v.Ok || v.State == "已禁用" {
				continue
			}
			if v.SourceType_IsWebsocket_Receive() {
				if bytes.Contains(bs, v.BytesLod) {
					bs = bytes.ReplaceAll(bs, v.BytesLod, v.BytesNew)
				}
			}
		}
		Conn.SetBody(bs)
	}
}
