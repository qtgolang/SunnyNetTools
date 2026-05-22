package mcp

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

// BuildFilterJSON 将 Rules + Join 转为 Session.ParseFilter 使用的 JSON（与 capture_search 一致）。
func BuildFilterJSON(rulesRaw any, joinOp string, required bool) (string, error) {
	return buildFilterJSONString(rulesRaw, joinOp, required)
}

// BuildQuickFilterJSON 合并 keys 与 mode 为快速过滤 JSON。
func BuildQuickFilterJSON(keys []string, mode string) string {
	if mode != "exclude" {
		mode = "match"
	}
	b, _ := json.Marshal(map[string]any{"keys": keys, "mode": mode})
	return string(b)
}

// QuickFilterKeysMode 解析 jsonKeys 字符串。
func QuickFilterKeysMode(jsonKeys string) (keys []string, mode string) {
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

// BuildRulesFilterForKeys 多关键字 OR 过滤（每个 key 在常用列上 contains）。
func BuildRulesFilterForKeys(keys []string, mode string, join string) (string, error) {
	if len(keys) == 0 {
		return "", nil
	}
	filterType := "contains"
	if mode == "exclude" {
		filterType = "notContains"
	}
	cols := []string{"主机名", "请求地址", "进程", "注释", "方式", "状态码", "HOST"}
	conds := make([]any, 0, len(keys)*len(cols))
	for _, key := range keys {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		for _, col := range cols {
			conds = append(conds, map[string]any{
				"filterType": "text",
				"colId":      col,
				"type":       filterType,
				"filter":     key,
			})
		}
	}
	if len(conds) == 0 {
		return "", errors.New("keys 无有效项")
	}
	if join == "" {
		join = "OR"
	}
	return buildFilterJSONString(conds, join, true)
}

// ParseRulesArray 从 args["Rules"] 或 args["rules"] 读取规则数组。
func ParseRulesArray(m map[string]any) ([]any, error) {
	if v, ok := m["Rules"]; ok {
		if arr, ok := v.([]any); ok {
			return arr, nil
		}
	}
	if v, ok := m["rules"]; ok {
		if arr, ok := v.([]any); ok {
			return arr, nil
		}
	}
	return nil, errors.New("Rules 类型错误")
}

// RuleFromKeyword 由关键词构造单条 ag-grid 文本过滤 JSON。
func RuleFromKeyword(text, mode, scope string) (string, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return "", errors.New("text 为空")
	}
	filterType := "contains"
	if strings.ToLower(mode) == "exclude" {
		filterType = "notContains"
	}
	f := map[string]any{
		"filterType": "text",
		"type":       filterType,
		"filter":     text,
	}
	if scope != "" && scope != "all" {
		f["colId"] = scope
	}
	b, err := json.Marshal(f)
	return string(b), err
}

// EncodeDataPart 按 type 编码二进制片段（auto/hex/base64/str）。
func EncodeDataPart(bs []byte, typ string) (data string, outType string) {
	typ = strings.TrimSpace(strings.ToLower(typ))
	if typ == "" {
		typ = "auto"
	}
	switch typ {
	case "hex":
		return bytesToHexSpaced(bs), "hex"
	case "base64":
		return base64.StdEncoding.EncodeToString(bs), "base64"
	case "str":
		return string(bs), "str"
	case "auto":
		if hasControl(bs) {
			return base64.StdEncoding.EncodeToString(bs), "base64"
		}
		return string(bs), "str"
	default:
		return base64.StdEncoding.EncodeToString(bs), "base64"
	}
}

func hasControl(bs []byte) bool {
	for _, b := range bs {
		if b < 32 && b != '\t' && b != '\n' && b != '\r' {
			return true
		}
		if b == 127 {
			return true
		}
	}
	if !utf8.Valid(bs) {
		return true
	}
	return false
}

// buildFilterJSONString 将 Rules 或已构造的 condition 列表转为 Session.ParseFilter 使用的 JSON。
func buildFilterJSONString(rulesRaw any, joinOp string, required bool) (string, error) {
	rulesArr, ok := rulesRaw.([]any)
	if !ok {
		return "", fmt.Errorf("Rules 类型错误")
	}
	conds := make([]any, 0, len(rulesArr))
	for _, item := range rulesArr {
		ruleMap, ok := item.(map[string]any)
		if !ok {
			continue
		}
		// 已是 ag-grid 条件对象（filterType + filter）
		if ft, _ := ruleMap["filterType"].(string); ft != "" {
			conds = append(conds, ruleMap)
			continue
		}
		colID, _ := ruleMap["colId"].(string)
		text, _ := ruleMap["text"].(string)
		typ, _ := ruleMap["type"].(string)
		if colID == "" {
			continue
		}
		if colID == "响应长度" || colID == "长度" {
			var n float64
			switch v := ruleMap["filter"].(type) {
			case float64:
				n = v
			case int:
				n = float64(v)
			case string:
				parsed, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
				if err != nil {
					return "", fmt.Errorf("解析长度数据错误: %w", err)
				}
				n = parsed
			default:
				parsed, err := strconv.Atoi(strings.TrimSpace(text))
				if err != nil {
					return "", fmt.Errorf("解析长度数据错误: %w", err)
				}
				n = float64(parsed)
			}
			conds = append(conds, map[string]any{
				"filterType": "number",
				"colId":      colID,
				"type":       typ,
				"filter":     n,
			})
			continue
		}
		filterVal := ruleMap["filter"]
		if filterVal == nil {
			filterVal = text
		}
		conds = append(conds, map[string]any{
			"filterType": "text",
			"colId":      colID,
			"type":       typ,
			"filter":     filterVal,
		})
	}
	if required && len(conds) == 0 {
		return "", fmt.Errorf("请至少提供一条过滤条件")
	}
	if len(conds) == 0 {
		return "", nil
	}
	if len(conds) == 1 {
		b, err := json.Marshal(conds[0])
		return string(b), err
	}
	if joinOp == "" {
		joinOp = "AND"
	}
	b, err := json.Marshal(map[string]any{
		"filterType": "join",
		"type":       joinOp,
		"conditions": conds,
	})
	return string(b), err
}

func bytesToHexSpaced(bs []byte) string {
	if len(bs) == 0 {
		return ""
	}
	const hexChars = "0123456789ABCDEF"
	out := make([]byte, len(bs)*3-1)
	j := 0
	for i, b := range bs {
		out[j], out[j+1] = hexChars[b>>4], hexChars[b&0x0F]
		if i < len(bs)-1 {
			out[j+2] = ' '
			j += 3
		} else {
			j += 2
		}
	}
	return string(out)
}
