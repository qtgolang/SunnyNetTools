package mcp

import (
	"changeme/Service/mcpbridge"
	"log"
	"net/http"
)

var mcpControl *mcpbridge.Control

// InitBridge 注册 MCP 业务回调并挂载 /doc；invoke 由 Service.MCPBridgeInvoke 提供。
func InitBridge(invoke func(op string, args map[string]any) (any, error)) {
	mcpbridge.BackendInvoke = invoke
	mcpbridge.RouteHook = func(mux *http.ServeMux) {
		mux.HandleFunc("/doc", handleDoc)
	}
	mcpControl = mcpbridge.NewControl()
}

// Start 启动 MCP HTTP（默认 6987）。
func Start() error {
	if mcpbridge.BackendInvoke == nil {
		log.Println("【MCP】BackendInvoke 未设置，跳过启动")
		return nil
	}
	if mcpControl == nil {
		mcpControl = mcpbridge.NewControl()
	}
	if msg := mcpControl.MCPEnable(mcpbridge.DefaultPort()); msg != "" {
		return &mcpStartError{msg}
	}
	return nil
}

type mcpStartError struct{ s string }

func (e *mcpStartError) Error() string { return e.s }
