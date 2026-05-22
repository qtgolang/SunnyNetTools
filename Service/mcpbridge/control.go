package mcpbridge

import (
	"encoding/json"
	"fmt"
	"sync"

	"changeme/Service/mcpcatalog"
)

const defaultMCPPort = 6987

const bridgeModeHTTP = "http"

// Control 启用/禁用本机 MCP HTTP 桥。
type Control struct {
	mu             sync.Mutex
	host           *Host
	httpBr         *httpBridge
	httpListenAddr string
	lastPort       int
}

// NewControl 构造 MCP 控制服务。
func NewControl() *Control {
	return &Control{host: NewHost()}
}

// MCPEnable 在 127.0.0.1:port 启动 HTTP MCP。
func (c *Control) MCPEnable(port int) string {
	return c.MCPEnableMode(port, bridgeModeHTTP)
}

// MCPEnableMode 仅支持 http。
func (c *Control) MCPEnableMode(port int, mode string) string {
	_ = mode
	return c.enableHTTP(port)
}

func (c *Control) enableHTTP(port int) string {
	if c == nil || c.host == nil {
		return "MCP 未初始化"
	}
	if port <= 0 || port > 65535 {
		port = defaultMCPPort
	}
	c.mu.Lock()
	if c.httpBr != nil && c.httpBr.Addr() != "" {
		c.mu.Unlock()
		return "MCP 已在运行"
	}
	c.mu.Unlock()

	host := "127.0.0.1"
	hb := newHTTPBridge(c.host)
	if err := hb.Start(host, port); err != nil {
		return err.Error()
	}
	c.mu.Lock()
	c.httpBr = hb
	c.httpListenAddr = hb.Addr()
	c.lastPort = port
	c.mu.Unlock()
	notifyMCPBridgeChanged()
	return ""
}

// MCPDisable 关闭 MCP HTTP 桥。
func (c *Control) MCPDisable() string {
	if c == nil {
		return ""
	}
	c.mu.Lock()
	hb := c.httpBr
	c.httpBr = nil
	c.httpListenAddr = ""
	c.mu.Unlock()
	if hb != nil {
		hb.Stop()
	}
	notifyMCPBridgeChanged()
	return ""
}

// MCPStatusJSON 返回桥状态 JSON。
func (c *Control) MCPStatusJSON() string {
	if c == nil {
		return `{"enabled":false}`
	}
	c.mu.Lock()
	httpAddr := c.httpListenAddr
	lp := c.lastPort
	httpOn := c.httpBr != nil && c.httpBr.Addr() != ""
	c.mu.Unlock()
	if !httpOn {
		lp = defaultMCPPort
	}
	mcpURL := ""
	if httpOn && httpAddr != "" {
		mcpURL = fmt.Sprintf("http://%s%s", httpAddr, MCPStreamablePath)
	}
	out := map[string]any{
		"enabled":           httpOn,
		"httpEnabled":       httpOn,
		"bridgeMode":        bridgeModeHTTP,
		"httpListenAddr":    httpAddr,
		"defaultPort":       defaultMCPPort,
		"lastPort":          lp,
		"mcpStreamablePath": MCPStreamablePath,
		"mcpStreamableURL":  mcpURL,
		"httpInvokePath":    SunnyNetHTTPAPIPrefix + "/invoke",
		"httpEventsPath":    SunnyNetHTTPAPIPrefix + "/events",
		"httpHealthPath":    SunnyNetHTTPAPIPrefix + "/health",
		"httpSupportedOps":  SunnyNetHTTPAPIPrefix + "/supported-ops",
		"httpApiPrefix":     SunnyNetHTTPAPIPrefix,
	}
	b, _ := json.Marshal(out)
	return string(b)
}

// MCPListOpsJSON 能力目录。
func (c *Control) MCPListOpsJSON() string {
	return mcpcatalog.SupportedOpsJSON()
}

// MCPDocURL 返回 /doc 文档页完整 URL；未启用时为空。
func (c *Control) MCPDocURL() string {
	if c == nil {
		return ""
	}
	c.mu.Lock()
	httpAddr := c.httpListenAddr
	httpOn := c.httpBr != nil && c.httpBr.Addr() != ""
	c.mu.Unlock()
	if !httpOn || httpAddr == "" {
		return ""
	}
	return fmt.Sprintf("http://%s/doc", httpAddr)
}

// DefaultPort 默认监听端口。
func DefaultPort() int {
	return defaultMCPPort
}
