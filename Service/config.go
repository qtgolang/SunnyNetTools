package Service

import (
	. "changeme/Service/Config"
	"changeme/Service/IsInstallCert"
	"changeme/Service/Session"
	Theme2 "changeme/Service/Theme"
	"net"
	"strings"
	"sync/atomic"
)

func getDev(theme string) string {
	arr := strings.Split(theme, "\n")
	s := strings.ReplaceAll(strings.Join(arr, ""), " ", "")
	start := strings.Index(s, ",Dev:\"")
	end := strings.LastIndex(s, "\"")
	if start != -1 && end != -1 && start < end {
		if start > 6 {
			return s[start+6 : end]
		}
	}
	return ""
}
func (g *AppMain) AppGetTheme(isDark bool) {
	if isDark {
		AppList["Main"].EmitEvent("setTheme", Config.AgGridDarkTheme, isDark)
	} else {
		AppList["Main"].EmitEvent("setTheme", Config.AgGridLightTheme, isDark)
	}
}

// Theme 获取或应用主题
func (g *AppMain) Theme(isDark bool, _type, Theme string) string {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	switch _type {
	case "默认暗黑配色1":
		return Theme2.ThemeDark1
	case "默认暗黑配色2":
		return Theme2.ThemeDark2
	case "默认暗黑配色3":
		return Theme2.ThemeDark3
	case "默认明亮配色1":
		return Theme2.ThemeLight1
	case "默认明亮配色2":
		return Theme2.ThemeLight2
	case "默认明亮配色3":
		return Theme2.ThemeLight3
	case "明亮", "暗黑":
		g.CallTools("主题设计", true, "")
		Dev := getDev(Theme)
		AppList["主题设计"].ExecJS(`window.SetLocalStorage("` + Dev + `")`)
		AppList["主题设计"].ExecJS(`location.reload();`)
		return ""
	case "应用暗黑配色", "应用明亮配色":
		if isDark {
			Config.AgGridDarkTheme = Theme
		} else {
			Config.AgGridLightTheme = Theme
		}
		Config.Save()
		AppList["Main"].EmitEvent("setTheme", Theme, isDark)
	}
	return Theme
}

// IsDark 设置当前是否选择暗色主题
func (g *AppMain) IsDark() bool {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	return Config.IsDark
}

// GetIPV4InterfaceAdders 获取所有IPV4网卡IP
func (g *AppMain) GetIPV4InterfaceAdders() []string {
	var array []string
	adders, err := net.InterfaceAddrs()
	if err != nil {
		return array
	}
	for _, addr := range adders {
		// 类型断言，排除非 IPNet 类型
		ipNet, ok := addr.(*net.IPNet)
		if !ok || ipNet.IP.IsLoopback() {
			continue
		}
		// 只要 IPv4（IPv6 会返回 nil）
		ip := ipNet.IP.To4()
		if ip == nil {
			continue
		}
		array = append(array, ip.String())
	}
	return array
}

// GetInterfaceOutRouterAdders 获取出口IP
func (g *AppMain) GetInterfaceOutRouterAdders() string {
	return Config.OutRouter
}

// SetInterfaceOutRouterAdders 设置出口IP
func (g *AppMain) SetInterfaceOutRouterAdders(ip string) {
	Config.OutRouter = ip
	g.app.SetOutRouterIP(ip)
}

// GetHomeTextMark 获取列表颜色标记
func (g *AppMain) GetHomeTextMark() string {
	return Config.HomeTextMark
}

// SetHomeTextMark 设置列表颜色标记
func (g *AppMain) SetHomeTextMark(TextMark string) {
	Config.HomeTextMark = TextMark
	Config.Save()
}

// AppDisconnectTCPRequest 断开选中的TCP请求
func (g *AppMain) AppDisconnectTCPRequest(id []int) {
	for _, v := range id {
		obj := Session.GetTCPSession(v)
		if obj != nil {
			_ = obj.Conn.Close()
		}
	}
}

// AppResendRequest 重发请求
func (g *AppMain) AppResendRequest(id int, count, BreakMode int) bool {
	obj := Session.GetAppSession(id)
	if obj == nil {
		return false
	}
	go obj.ResendRequest(count, BreakMode, g.app.Port(), Config.OutRouter)
	return true
}

// SetIsDark 设置当前选择的是否暗色主题
func (g *AppMain) SetIsDark(IsDark bool) {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	Config.IsDark = IsDark
	Config.Save()
	//任意窗口发送一次,所有窗口都会收到该事件
	AppList["Main"].EmitEvent("SetIsDark", IsDark)

	//Debug-SunnyNet

	//AppList["主题设计"].EmitEvent("SetIsDark", IsDark)
}

// GetAgGridLightTheme 获取输入的亮色主题
func (g *AppMain) GetAgGridLightTheme() string {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	return Config.AgGridLightTheme
}

// GetAgGridDarkTheme 获取输入的暗色主题
func (g *AppMain) GetAgGridDarkTheme() string {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	return Config.AgGridDarkTheme
}

// SetListColor 配置列表项颜色
func (g *AppMain) SetListColor(IsDark bool, ColorID string, Color string) {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	lock.Lock()
	if IsDark {
		Config.ListColor["d"+ColorID] = Color
	} else {
		Config.ListColor["l"+ColorID] = Color
	}
	lock.Unlock()
	AppList["Main"].EmitEvent("ListColor", ColorID, Color)
}
func (g *AppMain) GetListColor(IsDark bool) map[string]string {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	res := make(map[string]string)
	lock.Lock()
	for k, Color := range Config.ListColor {
		ColorID := strings.ReplaceAll(strings.ReplaceAll(k, "l", ""), "d", "")
		if IsDark {
			if strings.Contains(k, "d") {
				res[ColorID] = Color
			}
		} else {
			if !strings.Contains(k, "d") {
				res[ColorID] = Color
			}
		}
	}
	lock.Unlock()
	return res
}

// DefaultColor 恢复默认的配色方案
func (g *AppMain) DefaultColor(IsDark bool) {
	lock.Lock()
	for k, Color := range DefaultColor {
		ColorID := strings.ReplaceAll(strings.ReplaceAll(k, "l", ""), "d", "")
		if IsDark {
			if strings.Contains(k, "d") {
				Config.ListColor[k] = Color
				AppList["Main"].EmitEvent("ListColor", ColorID, Color)
			}
		} else {
			if !strings.Contains(k, "d") {
				Config.ListColor[k] = Color
				AppList["Main"].EmitEvent("ListColor", ColorID, Color)
			}
		}
	}
	lock.Unlock()
	AppList["Main"].EmitEvent("RestColor", IsDark)
}

// GetSendIsHTTP1 获取是否强制发送HTTP1.1请求
func (g *AppMain) GetSendIsHTTP1() bool {
	return Config.SendIsHTTP1
}

// SetSendIsHTTP1 获取是否强制发送HTTP1.1请求
func (g *AppMain) SetSendIsHTTP1(value bool) {
	lock.Lock()
	defer lock.Unlock()
	Config.SendIsHTTP1 = value
	Config.Save()
}
func (g *AppMain) AppCheckSunnyNet() bool {
	lock.Lock()
	defer lock.Unlock()
	a, _ := IsInstallCert.CheckSunnyNet()
	if !a {
		g.app.InstallCert()
		a, _ = IsInstallCert.CheckSunnyNet()
	}
	return a
}
func SetCurrentTheology(Theology int) {
	atomic.StoreInt64(&Config.CurrentTheology, int64(Theology))
}
