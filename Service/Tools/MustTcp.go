package Tools

import (
	"changeme/Service/Config"
	"github.com/qtgolang/SunnyNet/SunnyNet"
	"strings"
)

type MustTcp struct {
	App *SunnyNet.Sunny
}

func (g *MustTcp) SetMustTcpRoles(Type int, Roles string) {
	mapProxyLock.Lock()
	defer mapProxyLock.Unlock()
	Config.Config.MustTcp.Type = Config.MustTcpType(Type)
	g.App.MustTcp(Config.Config.MustTcp.Type == Config.MustTcpTypeAll)
	_ = g.App.SetMustTcpRegexp(ParseMustTcpRoles(Roles), Config.Config.MustTcp.Type == Config.MustTcpTypeLei)
	Config.Config.MustTcp.Roles = Roles
	Config.Config.Save()
}
func (g *MustTcp) GetMustTcpRoles() string {
	mapProxyLock.Lock()
	defer mapProxyLock.Unlock()
	return Config.Config.MustTcp.Roles
}
func (g *MustTcp) GetMustTcpType() int {
	mapProxyLock.Lock()
	defer mapProxyLock.Unlock()
	return int(Config.Config.MustTcp.Type)
}

func ParseMustTcpRoles(Roles string) string {
	array := strings.Split(strings.ReplaceAll(Roles, "\r", ""), "\n")
	str := ""
	for _, v := range array {
		o := strings.ReplaceAll(strings.TrimSpace(v), " ", "")
		if strings.HasPrefix(o, "//") {
			continue
		}
		if o == "" {
			continue
		}
		str += o + ";"
	}
	return str
}
