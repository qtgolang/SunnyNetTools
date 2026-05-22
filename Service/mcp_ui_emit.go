package Service

import (
	"encoding/json"

	"changeme/Service/Config"
)

// emitMCPMain 向主窗口 WebView 发送 MCP 同步事件（page=main）。
func emitMCPMain(tag, msg string) {
	if Config.AppList["Main"] == nil {
		return
	}
	Config.AppList["Main"].EmitEvent("mcp", map[string]any{
		"page": "main",
		"tag":  tag,
		"msg":  msg,
	})
}

func emitMCPMainJSON(tag string, payload any) {
	b, err := json.Marshal(payload)
	if err != nil {
		return
	}
	emitMCPMain(tag, string(b))
}

func emitMCPConfigReload(rule string) {
	emitMCPMainJSON("configreload", map[string]any{
		"rule":   rule,
		"action": "reload",
	})
}

// emitMCPRulesPageReload 通知规则页刷新（ReplaceBody 等监听 rule=all，只发一次避免重复追加行）。
func emitMCPRulesPageReload() {
	emitMCPConfigReload("all")
}

func emitMCPDevice(action string, payload map[string]any) {
	if payload == nil {
		payload = map[string]any{}
	}
	payload["action"] = action
	emitMCPMainJSON("device", payload)
}

func emitMCPRowRefresh(theology int) {
	emitMCPMainJSON("rowrefresh", map[string]any{"theology": theology})
}

// emitMCPRequestCertReload 通知请求证书窗口刷新表格（Cert 独立 WebView）。
func emitMCPRequestCertReload() {
	payload, _ := json.Marshal(map[string]any{"action": "reload"})
	msg := string(payload)
	ev := map[string]any{
		"page": "cert",
		"tag":  "requestcert",
		"msg":  msg,
	}
	if win := Config.AppList["Cert"]; win != nil {
		win.EmitEvent("mcp", ev)
	}
}
