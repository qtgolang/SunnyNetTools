package mcpbridge

// EmitAppConfigRulesChanged 由 main 注入：MCP config_set_* 改规则后通知前端重新读盘。
var EmitAppConfigRulesChanged func()
