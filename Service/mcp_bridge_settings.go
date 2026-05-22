package Service

import (
	"changeme/Service/Config"
	"errors"
	"fmt"
	"strings"
)

func emitMCPSettingsReload(scopes ...string) {
	for _, scope := range scopes {
		emitMCPMainJSON("settings", map[string]any{
			"scope":  scope,
			"action": "reload",
		})
	}
}

func proxyDnsSnapshot(app *AppMain) map[string]any {
	raw := app.GetProxyDns()
	mode := "remotes"
	remoteServer := raw
	switch raw {
	case "localhost", "local":
		mode = "local"
		remoteServer = ""
	case "remote":
		mode = "remote"
		remoteServer = ""
	default:
		if raw == "" {
			mode = "local"
			remoteServer = ""
		}
	}
	return map[string]any{
		"mode":         mode,
		"remoteServer": remoteServer,
		"raw":          raw,
		"convention": map[string]any{
			"modes": []map[string]any{
				{"mode": "local", "storedValue": "localhost", "label": "本地解析"},
				{"mode": "remote", "storedValue": "remote", "label": "远程解析"},
				{"mode": "remotes", "label": "远程服务器解析", "note": "须为 TLS DNS（端口 853），如 223.5.5.5:853"},
			},
		},
	}
}

func bridgeProxyDnsSet(app *AppMain, m map[string]any) (any, error) {
	mode := strings.TrimSpace(strings.ToLower(argString(m, "mode")))
	remote := strings.TrimSpace(argString(m, "remoteServer"))
	if remote == "" {
		remote = strings.TrimSpace(argString(m, "dns"))
	}
	dns := ""
	switch mode {
	case "local", "localhost":
		dns = "localhost"
	case "remote":
		dns = "remote"
	case "remotes", "remote_server", "server":
		if remote == "" {
			return nil, errors.New("mode=remotes 时 remoteServer 必填（如 223.5.5.5:853）")
		}
		if !strings.Contains(remote, ":853") {
			return nil, fmt.Errorf("远程 DNS 须使用 853 端口（TLS），当前: %s", remote)
		}
		dns = remote
	case "":
		if remote != "" {
			dns = remote
		} else {
			return nil, errors.New("mode 必填：local | remote | remotes")
		}
	default:
		return nil, errors.New("mode 须为 local、remote 或 remotes")
	}
	app.SetProxyDns(dns)
	emitMCPSettingsReload("proxy_dns")
	snap := proxyDnsSnapshot(app)
	return settingsApplyResult(app, map[string]any{
		"dns": snap["raw"], "mode": snap["mode"], "remoteServer": snap["remoteServer"],
	}), nil
}

func proxyWaySnapshot(app *AppMain) map[string]any {
	list := app.ProxyWayList()
	proxies := make([]map[string]any, 0, len(list))
	for _, p := range list {
		proxies = append(proxies, proxyWayItemFromInfo(p))
	}
	conv := proxyWayURLConvention()
	conv["state"] = []string{"启用", "禁用"}
	conv["note"] = "同时仅一条可为启用；启用后成为全局上游代理"
	return map[string]any{
		"proxies":    proxies,
		"total":      len(proxies),
		"convention": conv,
	}
}

func normalizeProxyWayState(m map[string]any, defaultState string) (string, error) {
	state := strings.TrimSpace(defaultState)
	if state == "" {
		state = "禁用"
	}
	if v := strings.TrimSpace(argString(m, "state")); v != "" {
		state = v
	} else if _, ok := m["enabled"]; ok {
		if argBool(m, "enabled", false) {
			state = "启用"
		} else {
			state = "禁用"
		}
	}
	if state != "启用" && state != "禁用" {
		return "", errors.New("state 须为 启用 或 禁用")
	}
	return state, nil
}

func proxyWayItemFromInfo(p Config.ProxyWayInfo) map[string]any {
	_, proxyType, _ := validateProxyWayURL(p.URL)
	return map[string]any{
		"id": p.ID, "url": p.URL, "state": p.State, "note": p.Note,
		"enabled": p.State == "启用", "type": proxyType,
	}
}

func bridgeProxyWayAdd(app *AppMain, m map[string]any) (any, error) {
	proxyURL, _, err := resolveProxyWayURL(m)
	if err != nil {
		return nil, err
	}
	state, err := normalizeProxyWayState(m, "禁用")
	if err != nil {
		return nil, err
	}
	note := strings.TrimSpace(argString(m, "note"))

	id := app.CreateProxyWay()
	if !app.ProxyWayUpdate(id, proxyURL, state, note) {
		app.ProxyWayRemove(id)
		return nil, errors.New("上游代理添加失败，请检查 url 是否可达")
	}
	obj := Config.Config.ProxyWay[id]
	emitMCPSettingsReload("proxy_way")
	out := settingsApplyResult(app, map[string]any{
		"proxy": proxyWayItemFromInfo(*obj),
	})
	return out, nil
}

func bridgeProxyWayUpdateNote(app *AppMain, m map[string]any) (any, error) {
	id := argInt(m, "id", 0)
	if id <= 0 {
		return nil, errors.New("id 必填")
	}
	obj := Config.Config.ProxyWay[id]
	if obj == nil {
		return nil, fmt.Errorf("上游代理 id %d 不存在", id)
	}
	note := obj.Note
	if v, ok := m["note"]; ok {
		note = strings.TrimSpace(fmt.Sprint(v))
	}
	if !app.ProxyWayUpdate(id, obj.URL, obj.State, note) {
		return nil, fmt.Errorf("上游代理 id %d 注释更新失败", id)
	}
	emitMCPSettingsReload("proxy_way")
	return settingsApplyResult(app, map[string]any{
		"id": id, "note": note, "proxy": proxyWayItemFromInfo(*Config.Config.ProxyWay[id]),
	}), nil
}

func bridgeProxyWayDelete(app *AppMain, m map[string]any) (any, error) {
	id := argInt(m, "id", 0)
	if id <= 0 {
		return nil, errors.New("id 必填")
	}
	if Config.Config.ProxyWay[id] == nil {
		return nil, fmt.Errorf("上游代理 id %d 不存在", id)
	}
	app.ProxyWayRemove(id)
	emitMCPSettingsReload("proxy_way")
	return settingsApplyResult(app, map[string]any{"id": id}), nil
}

func bridgeProxyWaySetState(app *AppMain, m map[string]any) (any, error) {
	id := argInt(m, "id", 0)
	if id <= 0 {
		return nil, errors.New("id 必填")
	}
	obj := Config.Config.ProxyWay[id]
	if obj == nil {
		return nil, fmt.Errorf("上游代理 id %d 不存在", id)
	}
	state, err := normalizeProxyWayState(m, obj.State)
	if err != nil {
		return nil, err
	}
	if !app.ProxyWayUpdate(id, obj.URL, state, obj.Note) {
		return nil, fmt.Errorf("上游代理 id %d 状态更新失败", id)
	}
	emitMCPSettingsReload("proxy_way")
	return settingsApplyResult(app, map[string]any{"id": id, "state": state}), nil
}

func bridgeProxyWayUpdate(app *AppMain, m map[string]any) (any, error) {
	id := argInt(m, "id", 0)
	if id <= 0 {
		return nil, errors.New("id 必填")
	}
	obj := Config.Config.ProxyWay[id]
	if obj == nil {
		return nil, fmt.Errorf("上游代理 id %d 不存在", id)
	}
	url := obj.URL
	changesURL := m["url"] != nil || m["scheme"] != nil || m["host"] != nil || m["port"] != nil ||
		m["username"] != nil || m["user"] != nil || m["password"] != nil || m["pass"] != nil
	if changesURL {
		mm := m
		if argString(m, "url") == "" {
			mm = map[string]any{
				"url": obj.URL, "scheme": m["scheme"], "host": m["host"], "port": m["port"],
				"username": m["username"], "user": m["user"], "password": m["password"], "pass": m["pass"],
			}
		}
		merged, _, err := resolveProxyWayURL(mm)
		if err != nil {
			return nil, err
		}
		url = merged
	}
	if _, _, err := validateProxyWayURL(url); err != nil {
		return nil, err
	}
	state := obj.State
	if _, hasState := m["state"]; hasState || m["enabled"] != nil {
		var err error
		state, err = normalizeProxyWayState(m, obj.State)
		if err != nil {
			return nil, err
		}
	}
	note := obj.Note
	if v, ok := m["note"]; ok {
		note = strings.TrimSpace(fmt.Sprint(v))
	}
	if url == "" {
		return nil, errors.New("url 不能为空")
	}
	if !app.ProxyWayUpdate(id, url, state, note) {
		return nil, fmt.Errorf("上游代理 id %d 更新失败", id)
	}
	emitMCPSettingsReload("proxy_way")
	return settingsApplyResult(app, map[string]any{
		"id": id, "url": url, "state": state, "note": note,
		"proxy": proxyWayItemFromInfo(*Config.Config.ProxyWay[id]),
	}), nil
}

func bridgeProxyRolesGet(app *AppMain) (any, error) {
	return map[string]any{
		"roles": app.GetProxyRoles(),
		"convention": map[string]any{
			"format": "换行或 ; 分号分割；// 开头为注释行",
			"example": "*.test.com;qqqqq.com;\ndome.com\n*.abc.mmm.cn",
		},
	}, nil
}

func bridgeProxyRolesSet(app *AppMain, m map[string]any) (any, error) {
	roles := argString(m, "roles")
	if roles == "" {
		roles = argString(m, "rolesJSON")
	}
	app.SetProxyRoles(roles)
	emitMCPSettingsReload("proxy_roles")
	return settingsApplyResult(app, map[string]any{"rolesLength": len(roles)}), nil
}

func mustTcpTypeToUI(t int) string {
	switch Config.MustTcpType(t) {
	case Config.MustTcpTypeAll:
		return "ALLMustTcp"
	case Config.MustTcpTypeWai:
		return "MustTcpWai"
	default:
		return "MustTcpLei"
	}
}

func mustTcpTypeFromUI(s string) (int, error) {
	switch strings.TrimSpace(s) {
	case "ALLMustTcp", "all", "0":
		return int(Config.MustTcpTypeAll), nil
	case "MustTcpWai", "wai", "2":
		return int(Config.MustTcpTypeWai), nil
	case "MustTcpLei", "lei", "1", "":
		return int(Config.MustTcpTypeLei), nil
	default:
		if n := argInt(map[string]any{"type": s}, "type", -1); n >= 0 && n <= 2 {
			return n, nil
		}
		return 0, errors.New("type 须为 ALLMustTcp | MustTcpLei | MustTcpWai 或 0|1|2")
	}
}

func bridgeMustTcpGet(app *AppMain) (any, error) {
	return map[string]any{
		"type":     app.GetMustTcpType(),
		"typeUI":   mustTcpTypeToUI(app.GetMustTcpType()),
		"roles":    app.GetMustTcpRoles(),
		"convention": map[string]any{
			"types": map[string]any{
				"ALLMustTcp":  "全部走 TCP (0)",
				"MustTcpLei":  "规则内走 TCP (1)",
				"MustTcpWai":  "规则外走 TCP (2)",
			},
			"rolesFormat": "多个规则用 ; 分号分割",
		},
	}, nil
}

func bridgeMustTcpSet(app *AppMain, m map[string]any) (any, error) {
	tp := -1
	if v := strings.TrimSpace(argString(m, "typeUI")); v != "" {
		var err error
		tp, err = mustTcpTypeFromUI(v)
		if err != nil {
			return nil, err
		}
	} else if _, ok := m["type"]; ok {
		tp = argInt(m, "type", int(Config.MustTcpTypeLei))
	} else {
		tp = int(Config.MustTcpTypeLei)
	}
	roles := app.GetMustTcpRoles()
	if _, ok := m["roles"]; ok {
		roles = argString(m, "roles")
	}
	app.SetMustTcpRoles(tp, roles)
	emitMCPSettingsReload("musttcp")
	return map[string]any{"ok": true, "type": tp, "typeUI": mustTcpTypeToUI(tp)}, nil
}

func engineDisableStatus(app *AppMain) map[string]any {
	tcp, udp, cache, _, limit := app.GetBaseSettingsValue()
	return map[string]any{
		"disableTCP":   tcp,
		"disableUDP":   udp,
		"disableCache": cache,
		"limitRequestSize": limit,
		"items": []map[string]any{
			{"key": "disableTCP", "label": "禁用TCP", "disabled": tcp, "enabled": !tcp},
			{"key": "disableUDP", "label": "禁用UDP", "disabled": udp, "enabled": !udp},
			{"key": "disableCache", "label": "禁用浏览器缓存", "disabled": cache, "enabled": !cache},
		},
	}
}

func bridgeEngineTogglesGet(app *AppMain) (any, error) {
	return engineDisableStatus(app), nil
}

func bridgeDisableTCPGet(app *AppMain) (any, error) {
	tcp, _, _, _, _ := app.GetBaseSettingsValue()
	return map[string]any{
		"disableTCP": tcp,
		"disabled":   tcp,
		"enabled":    !tcp,
		"label":      "禁用TCP",
	}, nil
}

func bridgeDisableUDPGet(app *AppMain) (any, error) {
	_, udp, _, _, _ := app.GetBaseSettingsValue()
	return map[string]any{
		"disableUDP": udp,
		"disabled":   udp,
		"enabled":    !udp,
		"label":      "禁用UDP",
	}, nil
}

func bridgeDisableCacheGet(app *AppMain) (any, error) {
	_, _, cache, _, _ := app.GetBaseSettingsValue()
	return map[string]any{
		"disableCache": cache,
		"disabled":     cache,
		"enabled":      !cache,
		"label":        "禁用浏览器缓存",
	}, nil
}

func bridgeDisableTCPSet(app *AppMain, m map[string]any) (any, error) {
	v, ok := m["disableTCP"]
	if !ok {
		v = m["disabled"]
	}
	if !ok {
		return nil, errors.New("disableTCP 或 disabled 必填")
	}
	disabled := boolFromAny(v)
	app.SetDisableTCP(disabled)
	emitMCPSettingsReload("base")
	return settingsApplyResult(app, map[string]any{"disableTCP": disabled}), nil
}

func bridgeDisableUDPSet(app *AppMain, m map[string]any) (any, error) {
	v, ok := m["disableUDP"]
	if !ok {
		v = m["disabled"]
	}
	if !ok {
		return nil, errors.New("disableUDP 或 disabled 必填")
	}
	disabled := boolFromAny(v)
	app.SetDisableUDP(disabled)
	emitMCPSettingsReload("base")
	return settingsApplyResult(app, map[string]any{"disableUDP": disabled}), nil
}

func bridgeDisableCacheSet(app *AppMain, m map[string]any) (any, error) {
	v, ok := m["disableCache"]
	if !ok {
		v = m["disabled"]
	}
	if !ok {
		return nil, errors.New("disableCache 或 disabled 必填")
	}
	disabled := boolFromAny(v)
	app.SetDisableCache(disabled)
	emitMCPSettingsReload("base")
	return settingsApplyResult(app, map[string]any{"disableCache": disabled}), nil
}

func bridgeLimitRequestSizeSet(app *AppMain, m map[string]any) (any, error) {
	size := argInt(m, "limitRequestSize", 0)
	if size <= 0 {
		size = argInt(m, "size", 0)
	}
	if size <= 0 {
		return nil, errors.New("limitRequestSize 必填（默认 1024000）")
	}
	app.SetLimitRequestSize(size)
	emitMCPSettingsReload("base")
	return map[string]any{"ok": true, "limitRequestSize": size}, nil
}

func bridgeHTTPSProtocolGet(app *AppMain) (any, error) {
	return map[string]any{
		"sendIsHTTP1":  app.GetSendIsHTTP1(),
		"protocol":     ternaryStr(app.GetSendIsHTTP1(), "http/1.1", "h2"),
		"label":        ternaryStr(app.GetSendIsHTTP1(), "仅使用 HTTP/1.1 发送", "HTTP/2.0 优先"),
		"randomJa3":    app.GetRandomJa3(),
		"http2Fingerprint": app.GetHTTPSProto(),
	}, nil
}

func ternaryStr(cond bool, a, b string) string {
	if cond {
		return a
	}
	return b
}

func bridgeHTTPSProtocolSet(app *AppMain, m map[string]any) (any, error) {
	var sendIsHTTP1 *bool
	if v, ok := m["sendIsHTTP1"]; ok {
		b := boolFromAny(v)
		sendIsHTTP1 = &b
	} else if p := strings.TrimSpace(argString(m, "protocol")); p != "" {
		b := p == "http/1.1" || strings.EqualFold(p, "h1")
		sendIsHTTP1 = &b
	}

	protoJSON := strings.TrimSpace(argString(m, "http2Fingerprint"))
	if protoJSON == "" {
		protoJSON = strings.TrimSpace(argString(m, "fingerprint"))
	}
	if tpl := strings.TrimSpace(argString(m, "template")); tpl != "" {
		t, ok := findHTTP2Template(tpl)
		if !ok {
			return nil, fmt.Errorf("未找到 HTTP2 模板 %q", tpl)
		}
		protoJSON = t["config"].(string)
	}

	var randomJa3 *bool
	if v, ok := m["randomJa3"]; ok {
		b := boolFromAny(v)
		randomJa3 = &b
	} else if v, ok := m["ja3"]; ok {
		b := boolFromAny(v)
		randomJa3 = &b
	}

	if sendIsHTTP1 == nil && protoJSON == "" && randomJa3 == nil {
		return nil, errors.New("请提供 protocol/sendIsHTTP1、http2Fingerprint 或 template、randomJa3 至少一项")
	}

	out, err := app.HTTPSProto.ApplyHTTPSProtocol(sendIsHTTP1, protoJSON, randomJa3)
	if err != nil {
		return nil, err
	}
	emitMCPSettingsReload("https")
	out["engine"] = app.ReapplyEngineFromConfig()
	out["applied"] = true
	return out, nil
}

func bridgeRandomJa3Set(app *AppMain, m map[string]any) (any, error) {
	v, ok := m["randomJa3"]
	if !ok {
		v = m["enabled"]
	}
	if !ok {
		return nil, errors.New("randomJa3 或 enabled 必填")
	}
	ja3 := boolFromAny(v)
	out, err := app.HTTPSProto.ApplyHTTPSProtocol(nil, "", &ja3)
	if err != nil {
		return nil, err
	}
	emitMCPSettingsReload("https")
	out["engine"] = app.ReapplyEngineFromConfig()
	return out, nil
}

func bridgeHTTP2FingerprintSet(app *AppMain, m map[string]any) (any, error) {
	proto := argString(m, "http2Fingerprint")
	if proto == "" {
		proto = argString(m, "fingerprint")
	}
	if proto == "" {
		proto = argString(m, "proto")
	}
	if tpl := strings.TrimSpace(argString(m, "template")); tpl != "" {
		t, ok := findHTTP2Template(tpl)
		if !ok {
			return nil, fmt.Errorf("未找到 HTTP2 模板 %q", tpl)
		}
		proto = t["config"].(string)
	}
	if strings.TrimSpace(proto) == "" {
		return nil, errors.New("http2Fingerprint 或 template 必填")
	}
	out, err := app.HTTPSProto.ApplyHTTPSProtocol(nil, proto, nil)
	if err != nil {
		return nil, err
	}
	emitMCPSettingsReload("https")
	out["engine"] = app.ReapplyEngineFromConfig()
	return out, nil
}

func bridgeHTTP2TemplateApply(app *AppMain, m map[string]any) (any, error) {
	name := strings.TrimSpace(argString(m, "name"))
	if name == "" {
		name = strings.TrimSpace(argString(m, "template"))
	}
	if name == "" {
		return nil, errors.New("name 或 template 必填")
	}
	t, ok := findHTTP2Template(name)
	if !ok {
		return nil, fmt.Errorf("未找到 HTTP2 模板 %q", name)
	}
	proto := t["config"].(string)
	sendH1 := false
	if v, ok := m["sendIsHTTP1"]; ok {
		b := boolFromAny(v)
		sendH1 = b
	} else if p := strings.TrimSpace(argString(m, "protocol")); p != "" {
		sendH1 = p == "http/1.1" || strings.EqualFold(p, "h1")
	}
	out, err := app.HTTPSProto.ApplyHTTPSProtocol(&sendH1, proto, nil)
	if err != nil {
		return nil, err
	}
	out["template"] = name
	emitMCPSettingsReload("https")
	out["engine"] = app.ReapplyEngineFromConfig()
	return out, nil
}

func bridgeReapplyEngine(app *AppMain) (any, error) {
	return map[string]any{
		"ok":      true,
		"applied": true,
		"engine":  app.ReapplyEngineFromConfig(),
	}, nil
}

func bridgeHTTP2TemplateList() (any, error) {
	names := make([]string, 0, len(http2FingerprintTemplates))
	for _, t := range http2FingerprintTemplates {
		names = append(names, t["name"].(string))
	}
	return map[string]any{
		"templates": http2FingerprintTemplates,
		"names":     names,
		"total":     len(names),
	}, nil
}

func bridgeHTTP2TemplateGet(m map[string]any) (any, error) {
	name := strings.TrimSpace(argString(m, "name"))
	if name == "" {
		return nil, errors.New("name 必填（如 Firefox、Chrome_117_120_124）")
	}
	t, ok := findHTTP2Template(name)
	if !ok {
		return nil, fmt.Errorf("未找到 HTTP2 模板 %q", name)
	}
	return map[string]any{"ok": true, "template": t}, nil
}
