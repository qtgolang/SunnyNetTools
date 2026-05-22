package mcpbridge

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

func argsMap(args map[string]any) map[string]any {
	if args == nil {
		return map[string]any{}
	}
	return args
}

func argInt(m map[string]any, k string, def int) int {
	v, ok := m[k]
	if !ok {
		return def
	}
	switch t := v.(type) {
	case float64:
		return int(t)
	case int:
		return t
	case int64:
		return int(t)
	case json.Number:
		i, _ := t.Int64()
		return int(i)
	case string:
		n, _ := strconv.Atoi(strings.TrimSpace(t))
		return n
	default:
		return def
	}
}

func argBool(m map[string]any, k string, def bool) bool {
	v, ok := m[k]
	if !ok {
		return def
	}
	b, ok := v.(bool)
	if ok {
		return b
	}
	if f, ok := v.(float64); ok {
		return f != 0
	}
	return def
}

func argString(m map[string]any, k string) string {
	v, ok := m[k]
	if !ok {
		return ""
	}
	s, _ := v.(string)
	return s
}

func argStringSlice(m map[string]any, k string) []string {
	v, ok := m[k]
	if !ok {
		return nil
	}
	arr, ok := v.([]any)
	if !ok {
		return nil
	}
	out := make([]string, 0, len(arr))
	for _, x := range arr {
		if s, ok := x.(string); ok && strings.TrimSpace(s) != "" {
			out = append(out, s)
		}
	}
	return out
}

// BackendInvoke 由主程序在启动时注入（避免 mcpbridge 与 Service 包循环依赖）。
var BackendInvoke func(op string, args map[string]any) (any, error)

// Host 将 MCP 调用转发至 BackendInvoke。
type Host struct{}

// NewHost 构造 Host。
func NewHost() *Host {
	return &Host{}
}

// Invoke 执行桥接操作。
func (h *Host) Invoke(op string, args map[string]any) (any, error) {
	if h == nil {
		return nil, errors.New("mcp host nil")
	}
	if BackendInvoke == nil {
		return nil, errors.New("mcp backend 未初始化")
	}
	return BackendInvoke(op, args)
}

// RowIDToTheology 解析 rowId；约定 rowId === strconv.Itoa(theology)，与 ag-grid 行键 Theology 一致。
func RowIDToTheology(rowID string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(rowID))
	return n
}
