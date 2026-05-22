package Service

import (
	"changeme/Service/Config"
	"changeme/Service/Tools"
)

// ReapplyEngineFromConfig 将当前 Config 中的引擎相关设置同步到运行中的 Sunny（MCP/UI 修改后调用）。
func (g *AppMain) ReapplyEngineFromConfig() map[string]any {
	if g.app == nil {
		return map[string]any{"engineReady": false}
	}
	c := Config.Config

	g.app.SetRandomTLS(c.RandomJa3)
	g.app.DisableTCP(c.DisableTCP)
	g.app.DisableUDP(c.DisableUDP)

	limit := c.LimitRequestSize
	if limit < 10240 {
		limit = 1024000
	}
	g.app.SetHTTPRequestMaxUpdateLength(int64(limit))

	g.SetProxyDns(c.ProxyDns)
	g.SetProxyRoles(c.ProxyRoles)

	g.app.MustTcp(c.MustTcp.Type == Config.MustTcpTypeAll)
	_ = g.app.SetMustTcpRegexp(
		Tools.ParseMustTcpRoles(c.MustTcp.Roles),
		c.MustTcp.Type == Config.MustTcpTypeLei,
	)

	enabledProxy := ""
	for _, v := range c.ProxyWay {
		if v != nil && v.State == "启用" && v.URL != "" {
			enabledProxy = v.URL
			break
		}
	}
	g.app.SetGlobalProxy(enabledProxy, 30*1000)

	return g.engineSettingsSnapshot()
}

func (g *AppMain) engineSettingsSnapshot() map[string]any {
	c := Config.Config
	enabledProxy := ""
	for _, v := range c.ProxyWay {
		if v != nil && v.State == "启用" {
			enabledProxy = v.URL
			break
		}
	}
	protocol := "h2"
	if c.SendIsHTTP1 {
		protocol = "http/1.1"
	}
	return map[string]any{
		"engineReady":      g.app != nil,
		"disableTCP":         c.DisableTCP,
		"disableUDP":         c.DisableUDP,
		"disableCache":       c.DisableCache,
		"limitRequestSize":   c.LimitRequestSize,
		"proxyDns":           c.ProxyDns,
		"proxyRolesCompiled": true,
		"globalProxy":        enabledProxy,
		"mustTcpType":        int(c.MustTcp.Type),
		"mustTcpTypeUI":      mustTcpTypeToUI(int(c.MustTcp.Type)),
		"randomJa3":          c.RandomJa3,
		"sendIsHTTP1":        c.SendIsHTTP1,
		"protocol":           protocol,
		"http2WillApply":     !c.SendIsHTTP1 && c.HTTPSProto != "",
	}
}

// settingsApplyResult MCP 写设置后的统一返回（含 applied 与引擎快照）。
func settingsApplyResult(app *AppMain, extra map[string]any) map[string]any {
	out := map[string]any{
		"ok":      true,
		"applied": true,
		"engine":  app.ReapplyEngineFromConfig(),
	}
	for k, v := range extra {
		out[k] = v
	}
	return out
}
