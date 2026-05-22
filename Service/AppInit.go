package Service

import (
	"changeme/Service/Config"
	"changeme/Service/Tools"
	"changeme/Service/Tools/DebugTools"
	"encoding/base64"
	"runtime"

	"github.com/qtgolang/SunnyNet/SunnyNet"
)

func newAppMain(sApp *SunnyNet.Sunny) *AppMain {
	A := &AppMain{app: sApp}
	{
		//初始化
		sApp.SetRandomTLS(Config.Config.RandomJa3)
		sApp.DisableTCP(Config.Config.DisableTCP)
		sApp.DisableUDP(Config.Config.DisableUDP)
		{
			for _, v := range Config.Config.Authentication {
				sApp.Socket5AddUser(v.User, v.Pass)
			}
		}
		A.Device = Tools.Device{App: sApp}
		A.BaseSettings = Tools.BaseSettings{App: sApp}
		A.AddCustomTools = Tools.AddCustomTools{App: sApp}
		A.ReplaceBody = Tools.ReplaceBody{App: sApp}
		A.RequestCert = Tools.RequestCert{App: sApp}
		A.HTTPSProto = Tools.HTTPSProto{App: sApp}
		A.ProxyWay = Tools.ProxyWay{App: sApp}
		A.MustTcp = Tools.MustTcp{App: sApp}
		A.ScriptLog = Tools.ScriptLog{App: sApp}
		A.DebugTools = DebugTools.DebugTools{App: sApp}
	}
	A.localServerInit()
	{
		bs, _ := base64.StdEncoding.DecodeString(Config.Config.ScriptCode)
		A.app.SetScriptCall(A.PrintScriptLog, A.SaveScriptCode)
		A.app.SetScriptCode(string(bs))
		A.ScriptLogInit()
	}
	_ = A.app.SetMustTcpRegexp(Tools.ParseMustTcpRoles(Config.Config.MustTcp.Roles), Config.Config.MustTcp.Type == Config.MustTcpTypeLei)
	if runtime.GOOS == "windows" {
		go A.statusProxy()
	}
	return A
}
