package mcpbridge

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"changeme/Service/mcpcatalog"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// MCPStreamablePath 为 Cursor 等通过 Streamable HTTP 接入 MCP 的路径（与 REST /invoke、/events 不同）。
const MCPStreamablePath = SunnyNetHTTPAPIPrefix + "/mcp"

func newStreamableMCPHandler(host *Host) http.Handler {
	srv := mcp.NewServer(&mcp.Implementation{Name: "sunnynet", Version: "1.0.0"}, &mcp.ServerOptions{
		Instructions: mcpcatalog.MCPStreamableInstructions,
	})
	for _, tt := range mcpcatalog.BridgeMCPTools() {
		op := tt.Op
		tool := &mcp.Tool{
			Name:        tt.MCPName,
			Description: tt.Description,
		}
		if schema := mcpcatalog.ToolInputSchema(op); schema != nil {
			tool.InputSchema = schema
		}
		mcp.AddTool(srv, tool, func(ctx context.Context, req *mcp.CallToolRequest, in map[string]any) (*mcp.CallToolResult, any, error) {
			args := in
			if args == nil {
				args = map[string]any{}
			}
			rawAny, err := host.Invoke(op, args)
			if err != nil {
				var errRes mcp.CallToolResult
				errRes.SetError(err)
				return &errRes, nil, nil
			}
			// 与 REST /invoke 的 result 字段一致：直接返回业务 JSON 对象，避免多包一层 {"result":...}
			switch v := rawAny.(type) {
			case map[string]any:
				return nil, v, nil
			case json.RawMessage:
				var out map[string]any
				if len(v) > 0 {
					_ = json.Unmarshal(v, &out)
				}
				if out == nil {
					out = map[string]any{}
				}
				return nil, out, nil
			case string:
				s := strings.TrimSpace(v)
				if s == "" {
					return nil, map[string]any{}, nil
				}
				var out map[string]any
				if json.Unmarshal([]byte(s), &out) == nil && out != nil {
					return nil, out, nil
				}
				var arr []any
				if json.Unmarshal([]byte(s), &arr) == nil {
					return nil, map[string]any{"rules": arr, "total": len(arr)}, nil
				}
				return nil, map[string]any{"text": v}, nil
			default:
				raw, _ := json.Marshal(v)
				var out map[string]any
				if len(raw) > 0 && json.Unmarshal(raw, &out) == nil && out != nil {
					return nil, out, nil
				}
				var arr []any
				if len(raw) > 0 && json.Unmarshal(raw, &arr) == nil {
					return nil, map[string]any{"data": arr}, nil
				}
				if out == nil {
					out = map[string]any{}
				}
				return nil, out, nil
			}
		})
	}
	return mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		return srv
	}, &mcp.StreamableHTTPOptions{
		Stateless: true,
	})
}
