package mcp

import (
	"changeme/Service/mcpbridge"
	"changeme/Service/mcpcatalog"
)

// StatusJSON 返回 MCP HTTP 桥状态 JSON。
func StatusJSON() string {
	if mcpControl == nil {
		return `{"enabled":false}`
	}
	return mcpControl.MCPStatusJSON()
}

// Enable 启动 MCP HTTP 桥；成功返回空字符串。
func Enable(port int) string {
	if mcpControl == nil {
		mcpControl = mcpbridge.NewControl()
	}
	return mcpControl.MCPEnable(port)
}

// Disable 关闭 MCP HTTP 桥。
func Disable() string {
	if mcpControl == nil {
		return ""
	}
	return mcpControl.MCPDisable()
}

// DefaultPort 默认监听端口。
func DefaultPort() int {
	return mcpbridge.DefaultPort()
}

// ListOpsJSON 返回 MCP 桥支持的全部 op 能力目录（JSON）。
func ListOpsJSON() string {
	return mcpcatalog.SupportedOpsJSON()
}

// DocURL 返回内置文档页地址；MCP 未启用时返回空字符串。
func DocURL() string {
	if mcpControl == nil {
		return ""
	}
	return mcpControl.MCPDocURL()
}
