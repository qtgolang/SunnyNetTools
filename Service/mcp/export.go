package mcp

import "context"

// CallWebview 向前端发送 MCP 消息并等待回包（与工具层 callWebviewMsg 相同）。
func CallWebview(ctx context.Context, obj McpMsg) string {
	return callWebviewMsg(ctx, obj)
}

// HomeStartCapture 开始抓包（前端）。
func HomeStartCapture() string {
	return McpMsg{}.homeStartCapture()
}

// HomeStopCapture 停止抓包（前端）。
func HomeStopCapture() string {
	return McpMsg{}.homeStopCapture()
}

// HomeCaptureState 抓包状态（前端）。
func HomeCaptureState() string {
	return McpMsg{}.homeCaptureState()
}

// HomeClearAll 清空列表（前端）。
func HomeClearAll() string {
	return McpMsg{}.homeClearAll()
}

// ProxyEnable 开启系统代理（前端）。
func ProxyEnable() string {
	return McpMsg{}.proxyEnable()
}

// ProxyDisable 关闭系统代理（前端）。
func ProxyDisable() string {
	return McpMsg{}.proxyDisable()
}

// ProxyState 系统代理状态（前端）。
func ProxyState() string {
	return McpMsg{}.proxyState()
}

// ThemeSetDark 暗色主题。
func ThemeSetDark() string {
	return McpMsg{}.themeSetDark()
}

// ThemeSetLight 亮色主题。
func ThemeSetLight() string {
	return McpMsg{}.themeSetLight()
}

// ThemeState 返回 Dark / Light 等。
func ThemeState() string {
	return McpMsg{}.themeState()
}

// SetMainFilter 设置主列表 ag-grid 过滤器（前端）。
func SetMainFilter(filterJSON string) string {
	return McpMsg{}.SetFilter(filterJSON)
}

// PbConvertViaWebview 由前端将 protobuf base64 转为 JSON。
func PbConvertViaWebview(dataB64 string) string {
	return callWebviewMsg(context.Background(), McpMsg{Page: "Main", Tag: "pbConvert", Msg: dataB64})
}
