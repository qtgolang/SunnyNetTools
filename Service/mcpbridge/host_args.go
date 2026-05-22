package mcpbridge

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// argIntSlice 从参数 map 读取整数数组（JSON 解码后多为 []any）。
func argIntSlice(m map[string]any, k string) []int {
	v, ok := m[k]
	if !ok {
		return nil
	}
	arr, ok := v.([]any)
	if !ok {
		return nil
	}
	out := make([]int, 0, len(arr))
	for _, x := range arr {
		switch t := x.(type) {
		case float64:
			out = append(out, int(t))
		case int:
			out = append(out, t)
		case int64:
			out = append(out, int(t))
		case json.Number:
			i, _ := t.Int64()
			out = append(out, int(i))
		case string:
			n, _ := strconv.Atoi(strings.TrimSpace(t))
			out = append(out, n)
		}
	}
	return out
}

// argRowIDs 解析主列表 rowId：支持 rowId（单个）、rowIds 或 ids（字符串数组）；去重保序。
func argRowIDs(m map[string]any) ([]string, error) {
	var out []string
	seen := map[string]struct{}{}
	add := func(id string) {
		id = strings.TrimSpace(id)
		if id == "" {
			return
		}
		if _, ok := seen[id]; ok {
			return
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	add(argString(m, "rowId"))
	for _, id := range argStringSlice(m, "rowIds") {
		add(id)
	}
	for _, id := range argStringSlice(m, "ids") {
		add(id)
	}
	if len(out) == 0 {
		return nil, errors.New("rowId、rowIds 或 ids 至少填一项")
	}
	return out, nil
}

// argTheologyList 解析 theology：支持 theology、theologies、rowId、rowIds、ids；去重保序。
func argTheologyList(m map[string]any) ([]int, error) {
	var out []int
	seen := map[int]struct{}{}
	add := func(th int) {
		if th <= 0 {
			return
		}
		if _, ok := seen[th]; ok {
			return
		}
		seen[th] = struct{}{}
		out = append(out, th)
	}
	if th := argInt(m, "theology", 0); th > 0 {
		add(th)
	}
	for _, th := range argIntSlice(m, "theologies") {
		add(th)
	}
	if rid := strings.TrimSpace(argString(m, "rowId")); rid != "" {
		th := RowIDToTheology(rid)
		if th <= 0 {
			return nil, fmt.Errorf("rowId 无效: %s", rid)
		}
		add(th)
	}
	for _, key := range []string{"rowIds", "ids"} {
		for _, rid := range argStringSlice(m, key) {
			rid = strings.TrimSpace(rid)
			if rid == "" {
				continue
			}
			th := RowIDToTheology(rid)
			if th <= 0 {
				return nil, fmt.Errorf("rowId 无效: %s", rid)
			}
			add(th)
		}
	}
	if len(out) == 0 {
		return nil, errors.New("theology、theologies、rowId、rowIds 或 ids 至少填一项")
	}
	return out, nil
}

func mcpBatchOK(count int) map[string]any {
	if count <= 1 {
		return map[string]any{"ok": true}
	}
	return map[string]any{"ok": true, "count": count}
}
