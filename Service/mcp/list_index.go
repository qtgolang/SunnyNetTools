package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
)

// ListIndexTheologyMap 以 listIndex 的十进制字符串为键，theology 为值。
// theology 可为负数；仅 0 表示解析失败或该行不存在。
type ListIndexTheologyMap map[string]int

// ListIndexToTheology 将主列表界面「序号」(listIndex，从 1 开始) 转为内部 theology。
func ListIndexToTheology(listIndex int) (int, error) {
	m, err := ListIndexesToTheologies(listIndex)
	if err != nil {
		return 0, err
	}
	th := m[strconv.Itoa(listIndex)]
	if th == 0 {
		return 0, fmt.Errorf("listIndex %d 不存在", listIndex)
	}
	return th, nil
}

// ListIndexesToTheologies 批量解析 listIndex → theology。
// 返回 map 键为 listIndex 字符串（如 "1"、"66"），值为 theology（可为负）；失败或不存在时为 0。
func ListIndexesToTheologies(listIndexes ...int) (ListIndexTheologyMap, error) {
	if len(listIndexes) == 0 {
		return nil, errors.New("没有传入 listIndex")
	}
	want := make(map[int]struct{}, len(listIndexes))
	for _, idx := range listIndexes {
		if idx <= 0 {
			return nil, fmt.Errorf("listIndex 须为正整数，收到 %d", idx)
		}
		want[idx] = struct{}{}
	}
	m, err := fetchListIndexTheologyMap(want)
	if err != nil {
		return nil, err
	}
	out := make(ListIndexTheologyMap, len(listIndexes))
	for _, idx := range listIndexes {
		key := strconv.Itoa(idx)
		if th, ok := m[idx]; ok {
			out[key] = th
		} else {
			out[key] = 0
		}
	}
	return out, nil
}

// fetchListIndexTheologyMap 向主窗口 ag-grid 查询 listIndex→theology 映射。
// want 非空时仅返回指定序号；为空时返回当前可见视图全部映射。
func fetchListIndexTheologyMap(want map[int]struct{}) (map[int]int, error) {
	msg := "[]"
	if len(want) > 0 {
		ids := make([]int, 0, len(want))
		for id := range want {
			ids = append(ids, id)
		}
		sort.Ints(ids)
		b, err := json.Marshal(ids)
		if err != nil {
			return nil, err
		}
		msg = string(b)
	}
	res := callWebviewMsg(context.Background(), McpMsg{Page: "main", Tag: "listIndexToTheology", Msg: msg})
	if res == "" || res == "timeout" || res == "处理失败" {
		return nil, fmt.Errorf("listIndex 映射失败: %s", res)
	}
	return parseListIndexTheologyJSON(res)
}

func parseListIndexTheologyJSON(raw string) (map[int]int, error) {
	var anyMap map[string]json.RawMessage
	if err := json.Unmarshal([]byte(raw), &anyMap); err != nil {
		return nil, err
	}
	out := make(map[int]int, len(anyMap))
	for k, v := range anyMap {
		idx, err := strconv.Atoi(k)
		if err != nil || idx <= 0 {
			continue
		}
		var th int
		if err := json.Unmarshal(v, &th); err != nil {
			var s string
			if err2 := json.Unmarshal(v, &s); err2 == nil {
				th, _ = strconv.Atoi(s)
			}
		}
		out[idx] = th
	}
	return out, nil
}
