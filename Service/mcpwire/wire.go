// Package mcpwire 定义主进程 MCP TCP 桥与 stdio 子进程之间的 JSON 行协议（与实现解耦，供 mcpbridge / cmd 共用）。
package mcpwire

import "encoding/json"

// AuthLine 旧版首行鉴权 JSON；当前主程序不校验令牌，stdio 子进程不应再发送此行。
type AuthLine struct {
	Auth string `json:"auth"`
}

// InvokeLine 业务调用一行 JSON。
type InvokeLine struct {
	ID   int            `json:"id"`
	Op   string         `json:"op"`
	Args map[string]any `json:"args,omitempty"`
}

// InvokeResponse 对 InvokeLine 的应答一行 JSON。
type InvokeResponse struct {
	ID     int             `json:"id"`
	OK     bool            `json:"ok"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  string          `json:"error,omitempty"`
}
