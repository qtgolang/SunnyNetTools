package Service

import "changeme/Service/mcp"

// MCPStatusJSON 返回 MCP 桥状态（JSON 字符串）。
func (g *AppMain) MCPStatusJSON() string {
	return mcp.StatusJSON()
}

// MCPEnable 启动 MCP HTTP 桥；port<=0 使用默认端口；成功返回空字符串。
func (g *AppMain) MCPEnable(port int) string {
	return mcp.Enable(port)
}

// MCPDisable 关闭 MCP HTTP 桥。
func (g *AppMain) MCPDisable() string {
	return mcp.Disable()
}

// MCPListOpsJSON 返回 MCP 桥能力目录 JSON（op 列表与 capabilities 说明）。
func (g *AppMain) MCPListOpsJSON() string {
	return mcp.ListOpsJSON()
}

// MCPDocURL 返回 MCP 文档页 URL；服务未启用时为空。
func (g *AppMain) MCPDocURL() string {
	return mcp.DocURL()
}
