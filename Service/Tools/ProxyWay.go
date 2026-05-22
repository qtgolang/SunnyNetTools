package Tools

import (
	"changeme/Service/Config"
	"github.com/qtgolang/SunnyNet/SunnyNet"
	"sort"
	"strings"
	"sync"
)

type ProxyWay struct {
	App *SunnyNet.Sunny
}

var mapProxyLock sync.Mutex

func (g *ProxyWay) SetProxyRoles(roles string) {
	mapProxyLock.Lock()
	defer mapProxyLock.Unlock()
	Config.Config.ProxyRoles = roles
	_ = g.App.CompileProxyRegexp(parseProxyRoles(roles))
	Config.Config.Save()
}
func parseProxyRoles(Roles string) string {
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

func (g *ProxyWay) SetProxyDns(dns string) {
	mapProxyLock.Lock()
	defer mapProxyLock.Unlock()
	Config.Config.ProxyDns = dns
	g.App.SetDnsServer(dns)
	Config.Config.Save()
}
func (g *ProxyWay) GetProxyDns() string {
	mapProxyLock.Lock()
	defer mapProxyLock.Unlock()
	return Config.Config.ProxyDns
}
func (g *ProxyWay) GetProxyRoles() string {
	mapProxyLock.Lock()
	defer mapProxyLock.Unlock()
	return Config.Config.ProxyRoles
}
func (g *ProxyWay) CreateProxyWay() int {
	mapProxyLock.Lock()
	defer mapProxyLock.Unlock()
	Config.CertID++
	id := Config.CertID
	Config.Config.ProxyWay[id] = &Config.ProxyWayInfo{
		ID: id,
	}
	return id
}
func (g *ProxyWay) ProxyWayUpdate(id int, URL, State, Note string) bool {
	mapProxyLock.Lock()
	defer mapProxyLock.Unlock()
	obj := Config.Config.ProxyWay[id]
	if obj == nil {
		return false
	}
	for _, v := range Config.Config.ProxyWay {
		v.State = "禁用"
	}
	obj.State = State
	obj.URL = URL
	obj.Note = Note
	var res *Config.ProxyWayInfo
	for _, v := range Config.Config.ProxyWay {
		if v.State == "启用" {
			res = v
			break
		}
	}
	if res != nil {
		g.App.SetGlobalProxy(res.URL, 30*1000) //代理超时30秒
	} else {
		g.App.SetGlobalProxy("", 30*1000) //代理超时30秒
	}
	Config.Config.Save()
	return true
}

func (g *ProxyWay) ProxyWayRemove(id int) {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	mapProxyLock.Lock()
	defer mapProxyLock.Unlock()
	obj := Config.Config.ProxyWay[id]
	if obj == nil {
		return
	}
	delete(Config.Config.ProxyWay, id)
	var res *Config.ProxyWayInfo
	for _, v := range Config.Config.ProxyWay {
		if v.State == "启用" {
			res = v
			break
		}
	}
	if res != nil {
		g.App.SetGlobalProxy(res.URL, 30*1000)
	} else {
		g.App.SetGlobalProxy("", 30*1000)
	}
	Config.Config.Save()
}

func (g *ProxyWay) ProxyWayList() []Config.ProxyWayInfo {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	mapReplaceLock.Lock()
	defer mapReplaceLock.Unlock()
	var keys []int
	for i := range Config.Config.ProxyWay {
		keys = append(keys, i)
	}
	sort.Ints(keys)
	var array []Config.ProxyWayInfo
	for _, k := range keys {
		obj := Config.Config.ProxyWay[k]
		array = append(array, *obj)
	}
	return array
}
