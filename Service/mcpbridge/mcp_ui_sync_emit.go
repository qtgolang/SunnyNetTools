package mcpbridge

// 以下回调由 main 注入，在 MCP 改写「与主界面展示相关」的运行时状态后通知前端同步。

// EmitEngineAdvancedMcpApplied MCP engine_apply_advanced 成功后，前端应调用 SunnyEngineAdvancedLastApplied 对齐高级设置表单。
var EmitEngineAdvancedMcpApplied func()

// EmitSystemProxyMcpToggled MCP system_proxy_enable / system_proxy_disable 成功后，前端应重新读取系统代理开关。
var EmitSystemProxyMcpToggled func()

// EmitMainListFiltersMcpSynced MCP 主列表筛选 syncUi=true 时，payload 为 JSON（quickFilterKeys / filterText / advancedSearch 等局部字段）。
var EmitMainListFiltersMcpSynced func(payload string)

// EmitMCPBridgeChanged MCPEnable / MCPDisable / MCPRegenerateToken 后通知所有 WebView 刷新 MCP 状态（如主窗状态栏）。
var EmitMCPBridgeChanged func()

func notifyMCPBridgeChanged() {
	if EmitMCPBridgeChanged != nil {
		EmitMCPBridgeChanged()
	}
}
