package Tools

import (
	"changeme/Service/Config"
	"changeme/Service/HookKeys"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/qtgolang/SunnyNet/SunnyNet"
	"github.com/qtgolang/SunnyNet/src/public"
)

type BaseSettings struct {
	App     *SunnyNet.Sunny
	IsStart bool
}

var baseSettings sync.Mutex

func (b *BaseSettings) GetColumnState() string {
	baseSettings.Lock()
	defer baseSettings.Unlock()
	return Config.Config.ColumnState
}
func (b *BaseSettings) SetColumnState(ColumnState string) {
	baseSettings.Lock()
	defer baseSettings.Unlock()
	Config.Config.ColumnState = ColumnState
	Config.Config.Save()
}
func (b *BaseSettings) GetPort() int {
	baseSettings.Lock()
	defer baseSettings.Unlock()
	return Config.Config.Port
}
func (b *BaseSettings) SetPort(Port int, noStart bool) string {
	baseSettings.Lock()
	defer baseSettings.Unlock()
	Config.Config.Port = Port
	Config.Config.Save()
	b.App.Close()
	b.App.SetPort(Port)
	if !noStart {
		b.App.Start()
	}
	if Config.Config.IEProxy {
		b.SetIEProxy()
	}
	e := b.App.Error
	b.IsStart = e == nil
	if e != nil {
		return e.Error()
	}
	return ""
}
func (b *BaseSettings) AppIsSetPort() bool {
	return Config.Config.IEProxy
}

// SetIEProxy 设置IE代理
func (b *BaseSettings) SetIEProxy() bool {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	if b.App.SetIEProxy() {
		Config.Config.IEProxy = true
	} else {
		Config.Config.IEProxy = false
	}
	Config.Config.Save()
	return Config.Config.IEProxy
}

// CancelIEProxy 取消IE代理
func (b *BaseSettings) CancelIEProxy() bool {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	if Config.Config.IEProxy {
		if b.App.CancelIEProxy() {
			Config.Config.IEProxy = false
		} else {
			return false
		}
		Config.Config.Save()
	}
	return true
}

// DisableTCP 禁用TCP
func (b *BaseSettings) SetDisableTCP(Disable bool) {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	baseSettings.Lock()
	defer baseSettings.Unlock()
	Config.Config.DisableTCP = Disable
	b.App.DisableTCP(Disable)
	Config.Config.Save()
	b.App.DisableTCP(Disable)
}

// DisableUDP 禁用UDP
func (b *BaseSettings) SetDisableUDP(Disable bool) {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	baseSettings.Lock()
	defer baseSettings.Unlock()
	Config.Config.DisableUDP = Disable
	b.App.DisableUDP(Disable)
	Config.Config.Save()
	b.App.DisableUDP(Disable)
}

// DisableCache 禁用浏览器缓存
func (b *BaseSettings) SetDisableCache(Disable bool) {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	baseSettings.Lock()
	defer baseSettings.Unlock()
	Config.Config.DisableCache = Disable
	Config.Config.Save()
}

// LimitRequestSize 设置限制请求大小
func (b *BaseSettings) SetLimitRequestSize(size int) {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	baseSettings.Lock()
	defer baseSettings.Unlock()
	Config.Config.LimitRequestSize = size
	b.App.SetHTTPRequestMaxUpdateLength(int64(size))
	Config.Config.Save()
}

// SetAuthMode 使用身份验证模式
func (b *BaseSettings) SetAuthMode(open bool) {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	baseSettings.Lock()
	defer baseSettings.Unlock()
	Config.Config.OpenAuthMode = open
	Config.Config.Save()
}

// AuthModeCreate 添加
func (b *BaseSettings) AuthModeCreate() int {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	baseSettings.Lock()
	defer baseSettings.Unlock()
	Config.CertID++
	id := Config.CertID
	Config.Config.AuthMap[id] = &Config.Auth{
		User: "",
		Pass: "",
	}
	return id
}

// AuthModeSet 设置身份验证账号密码内容
func (b *BaseSettings) AuthModeSet(id int, User, Pass string) {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	baseSettings.Lock()
	defer baseSettings.Unlock()
	obj := Config.Config.AuthMap[id]
	if obj == nil {
		return
	}
	obj.User = User
	obj.Pass = Pass
	Config.Config.Save()
}

// AuthModeRemove 删除身份验证
func (b *BaseSettings) AuthModeRemove(id int) {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	baseSettings.Lock()
	defer baseSettings.Unlock()
	obj := Config.Config.AuthMap[id]
	if obj == nil {
		return
	}
	delete(Config.Config.AuthMap, id)
}
func (b *BaseSettings) AuthModeList() []Config.Auth {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	mapReplaceLock.Lock()
	defer mapReplaceLock.Unlock()
	var keys []int
	for i := range Config.Config.AuthMap {
		keys = append(keys, i)
	}
	sort.Ints(keys)
	var array []Config.Auth
	for _, k := range keys {
		obj := Config.Config.AuthMap[k]
		array = append(array, *obj)
	}
	return array
}
func (b *BaseSettings) AppGetEditorFontSize() int {
	if Config.Config.EditorFontSize < 10 {
		Config.Config.EditorFontSize = 14
		Config.Config.Save()
	}
	return Config.Config.EditorFontSize
}
func (b *BaseSettings) AppSetEditorFontSize(size int) {
	if size < 10 {
		Config.Config.EditorFontSize = 14
	} else {
		Config.Config.EditorFontSize = size
	}
	Config.AppList["Main"].EmitEvent("updateEditorFontSize", Config.Config.EditorFontSize)
	Config.Config.Save()
}
func (b *BaseSettings) GetBaseSettingsValue() (bool, bool, bool, bool, int) {
	return Config.Config.DisableTCP, Config.Config.DisableUDP, Config.Config.DisableCache, Config.Config.OpenAuthMode, Config.Config.LimitRequestSize
}
func (b *BaseSettings) ResetALLConfig() {
	Config.Config.Reset()
	Config.AppList["Main"].EmitEvent("onRest", true)
	Config.Config.Tour[public.SunnyVersion] = true
	Config.Config.Save()
}
func (b *BaseSettings) GetTour(newTour bool) bool {
	lod := Config.Config.Tour[public.SunnyVersion]
	Config.Config.Tour[public.SunnyVersion] = newTour
	Config.Config.Save()
	return lod
}
func (b *BaseSettings) ExportCert(path string) bool {
	_ = os.Remove(path)
	return os.WriteFile(path, []byte(public.RootCa), 0777) == nil
}
func (b *BaseSettings) SetKeys(obj string) {
	HookKeys.RegisterKeys(obj, b.CallKeys)
	Config.Config.Keys = obj
	Config.Config.Save()
}
func (b *BaseSettings) GetKeys() string {
	return Config.Config.Keys
}
func (b *BaseSettings) CallKeys(id string) {
	if id == "Boss" {
		if !Config.AppList["Main"].IsMinimised() {
			Config.AppList["Main"].Minimise()
			time.Sleep(time.Millisecond * 100)
			Config.AppList["Main"].Hide()
		} else {
			Config.AppList["Main"].Show()
			time.Sleep(time.Millisecond * 100)
			Config.AppList["Main"].UnMinimise()
		}
		return
	}
	Config.AppList["Main"].EmitEvent("ExternalKeydownEventListener", id)
}
func (b *BaseSettings) SetIsEditKeyDown(i bool) {
	HookKeys.UpdateEditKeyDown(i)
}
