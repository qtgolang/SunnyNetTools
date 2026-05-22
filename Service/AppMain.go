package Service

import (
	"changeme/Service/Config"
	"changeme/Service/Session"
	"changeme/Service/Tools"
	"changeme/Service/Tools/DebugTools"
	"changeme/Service/clipboard"
	"changeme/Service/mcp"
	"changeme/Welcome"
	"fmt"
	"runtime"
	"time"

	"github.com/qtgolang/SunnyNet/Api"
	"github.com/qtgolang/SunnyNet/SunnyNet"
	"github.com/qtgolang/SunnyNet/src/http"
)

type AppMain struct {
	app *SunnyNet.Sunny
	Tools.ReplaceBody
	Tools.RequestCert
	Tools.ProxyWay
	Tools.HTTPSProto
	Tools.MustTcp
	Tools.BaseSettings
	Tools.Device
	Tools.ScriptLog
	Tools.AddCustomTools
	DebugTools.DebugTools
}

func NewAppServer() *AppMain {
	A := newAppMain(SunnyNet.NewSunny())
	A.app.SetGoCallback(A.httpCallback, A.tcpCallback, A.wsCallback, A.udpCallback)
	//A.app.SetGoCallback(nil, nil, nil, A.udpCallback)
	//A.app.SetMustTcpRegexp("124.222.224.186:8800;*.qq.com;*.baidu.com", false)
	//A.app.MustTcp(true)
	//fmt.Println(A.app.OpenDrive(true))
	return A
}

// McpFuncRes 获取 mcp 结果
func (g *AppMain) McpFuncRes(res mcp.McpMsg) {
	if mcp.MsgCallback != nil {
		mcp.MsgCallback(res)
	}
}

// GOOS 当前系统是否Windows
func (g *AppMain) GOOS() bool {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	return runtime.GOOS == "windows"
}

// 检测系统代理
func (g *AppMain) statusProxy() {
	return
	for {
		time.Sleep(time.Second)
		//fmt.Println(gosysproxy.Status())
	}
}

// Start 启动
func (g *AppMain) Start() {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	g.app.Start()
	if runtime.GOOS == "windows" {
		go func() {
			Config.AppList["Main"].Show()
			Welcome.Stop()
		}()
	}
	g.IsStart = g.app.Error == nil
}

// GetError 获取错误信息
func (g *AppMain) GetError() string {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	if g.app.Error == nil {
		return ""
	}
	return g.app.Error.Error()
}

func (g *AppMain) CallTools(name string, open bool, args string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("CallTools panic:", err)
		}
	}()

	lock.Lock()
	defer lock.Unlock()

	getOrCreateWindow := func(windowName string, createFunc func()) *Config.AppWindow {
		obj := Config.AppList[windowName]
		if obj == nil && open && createFunc != nil {
			createFunc()
			obj = Config.AppList[windowName]
		}
		return obj
	}

	var obj *Config.AppWindow

	switch name {
	case "Cert":
		obj = getOrCreateWindow(name, CreateCertWindow)
	case "ReplaceBody":
		obj = getOrCreateWindow(name, CreateReplaceWindow)
	case "主题调色":
		obj = getOrCreateWindow(name, CreateThemeWindow)
	case "主题设计":
		obj = getOrCreateWindow(name, CreateThemeDesignWindow)
	case "调试工具":
		obj = getOrCreateWindow(name, CreateDebugWindow)
	case "证书安装", "脚本代码", "代码生成", "文本对比", "MCP能力描述":
		// 特殊处理，复用 "其他窗口"
		obj = getOrCreateWindow("其他窗口", CreateOtherWindow)
	default:
		obj = Config.AppList[name]
		if obj == nil {
			panic("未找到指定窗口名称: " + name)
		}
	}

	if obj == nil {
		fmt.Println("未找到窗口:", name)
		return
	}

	if !open {
		obj.Hide()
		return
	}

	if !obj.IsVisible() {
		obj.Hide()
	}
	obj.Center()
	obj.Show()
	obj.SetAlwaysOnTop(true)
	go func() {
		time.Sleep(time.Second)
		obj.SetAlwaysOnTop(false)
	}()
	if name == "代码生成" || name == "证书安装" || name == "脚本代码" || name == "文本对比" || name == "MCP能力描述" {
		obj.SetTitle(name)
		obj.EmitEvent("LoadUrl", name, args)
	}
}

func (g *AppMain) GoGetHex(data []byte) string {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	return Session.GetHexAllSpaces(data)
}

// GetAllStream 获取ws,tcp,udp,全部Stream
func (g *AppMain) GetAllStream(Theology int) []Session.UpdateSocketStream {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	g.AppStartInsert()
	var array []Session.UpdateSocketStream
	{
		stream := Session.GetAppSession(Theology)
		if stream != nil {
			lock.Lock()
			stream.RangeStream(func(val Session.AppStream) bool {
				array = append(array, val.ToUpdateStream(Theology, stream.GetStreamFilter()))
				return true
			})
			SetCurrentTheology(Theology)
			lock.Unlock()
		}
	}
	return array
}

// ClipboardReadAll 获取剪辑版内容,成功返回空字符串
func (g *AppMain) ClipboardReadAll() string {
	return clipboard.ClipboardReadAll()
}
func (g *AppMain) ProtobufToJson(aa []byte, skip int) string {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	if skip > len(aa) {
		return ""
	}
	return Api.PbToJson(aa[skip:])
}

// AppGenerateCode 生成代码
func (g *AppMain) AppGenerateCode(Theology int, Language, Module string) string {
	obj := Session.GetAppSession(Theology)
	if obj == nil {
		return "没有找到这个请求"
	}
	e := Session.CreateRequestCode(obj, Language, Module)
	if e != nil {
		return e.Error()
	}
	return ""
}

// AppSaveGenerateCodeInterface 保存生成代码接口
func (g *AppMain) AppSaveGenerateCodeInterface(code string) {
	Config.Config.GenerateCodeInterface = code
	Config.Config.InitCodeTemplate()
	Config.AppList["Main"].EmitEvent("GenerateCodeInterface", Config.Config.GenerateCodeInterface)
	Config.Config.Save()
}

// AppGetGenerateCodeInterface 保存生成代码接口
func (g *AppMain) AppGetGenerateCodeInterface() string {
	Config.Config.InitCodeTemplate()
	return Config.Config.GenerateCodeInterface
}

// AppGenerateCodeInterface 生成代码接口
func (g *AppMain) AppGenerateCodeInterface(Theology int) GenerateCodeInterface {
	obj := Session.GetAppSession(Theology)
	var e GenerateCodeInterface
	if obj == nil {
		return e
	}
	if obj.IsWebsocket() || obj.IsTCP() {
		path, _ := Session.ExportMessage(obj)
		e.Path = path
	}
	if obj.IsHTTP() {
		H, OK := obj.(*Session.HttpSession)
		if OK {
			e.URL = H.Request.Url
			e.Header = H.Request.Header
			e.IP = H.Response.ServerIP
			e.Method = H.Request.Method
			e.Body = H.Request.Body
		}
	}
	if obj.IsTCP() {
		T, OK := obj.(*Session.TCPSession)
		if OK {
			e.IP = T.RemoteAddress
			e.URL = T.Host
		}
	}
	return e
}

type GenerateCodeInterface struct {
	Path   string
	URL    string
	Header http.Header
	IP     string
	Method string
	Body   []byte
}
