package mcpbridge

import (
	"encoding/json"
	"strings"
)

// argSyncUi 为 true 时 MCP 在改写 liststore 后同步主界面筛选控件；默认 false（仅后端过滤）。
func argSyncUi(m map[string]any) bool {
	return argBool(m, "syncUi", false)
}

func argFilterMode(m map[string]any) string {
	mode := strings.TrimSpace(strings.ToLower(argString(m, "mode")))
	if mode == "exclude" {
		return "exclude"
	}
	return "match"
}

// normalizeQuickFiltersJSON 合并 jsonKeys 与顶层 mode；jsonKeys 可为数组或 {keys,mode}。
func normalizeQuickFiltersJSON(jsonKeys, modeArg string) string {
	payload := strings.TrimSpace(jsonKeys)
	mode := strings.TrimSpace(strings.ToLower(modeArg))
	if mode != "exclude" {
		mode = "match"
	}
	if payload == "" {
		if mode == "match" {
			return ""
		}
		b, _ := json.Marshal(map[string]any{"keys": []string{}, "mode": mode})
		return string(b)
	}
	var keys []string
	if err := json.Unmarshal([]byte(payload), &keys); err == nil {
		b, _ := json.Marshal(map[string]any{"keys": keys, "mode": mode})
		return string(b)
	}
	var obj struct {
		Keys []string `json:"keys"`
		Mode string   `json:"mode"`
	}
	if err := json.Unmarshal([]byte(payload), &obj); err == nil {
		if modeArg != "" {
			obj.Mode = mode
		} else if strings.TrimSpace(strings.ToLower(obj.Mode)) != "exclude" {
			obj.Mode = "match"
		} else {
			obj.Mode = "exclude"
		}
		b, _ := json.Marshal(obj)
		return string(b)
	}
	return payload
}

func quickFiltersFromJSON(jsonKeys string) (keys []string, mode string) {
	mode = "match"
	payload := strings.TrimSpace(jsonKeys)
	if payload == "" {
		return nil, mode
	}
	if err := json.Unmarshal([]byte(payload), &keys); err == nil {
		return keys, mode
	}
	var obj struct {
		Keys []string `json:"keys"`
		Mode string   `json:"mode"`
	}
	if err := json.Unmarshal([]byte(payload), &obj); err == nil {
		if strings.TrimSpace(strings.ToLower(obj.Mode)) == "exclude" {
			mode = "exclude"
		}
		return obj.Keys, mode
	}
	return nil, mode
}

func toolbarNeedleFromArgs(a map[string]any) (text, needle string) {
	text = argString(a, "text")
	needle = strings.TrimSpace(argString(a, "needle"))
	if needle == "" {
		needle = strings.TrimSpace(text)
	}
	return text, needle
}

func advancedNeedleFromArgs(a map[string]any) (queryText, needle string) {
	queryText = argString(a, "queryText")
	needle = strings.TrimSpace(argString(a, "needle"))
	if needle == "" {
		needle = strings.TrimSpace(queryText)
	}
	return queryText, needle
}

func emitMainListFiltersUISync(patch map[string]any) {
	if len(patch) == 0 || EmitMainListFiltersMcpSynced == nil {
		return
	}
	b, err := json.Marshal(patch)
	if err != nil {
		return
	}
	EmitMainListFiltersMcpSynced(string(b))
}

// mcpFilterInvokeResult syncUi 时落库并可选同步界面；否则仅查询并返回命中行 rowId（=strconv(theology)）。
func mcpFilterInvokeResult(syncUi bool, ids []string) map[string]any {
	if syncUi {
		return map[string]any{"ok": true}
	}
	if ids == nil {
		ids = []string{}
	}
	return map[string]any{
		"ok":    true,
		"ids":   ids,
		"total": len(ids),
	}
}
