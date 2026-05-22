package Service

import (
	"changeme/Service/Config"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const interceptRuleType = "拦截请求"

type replaceRuleMCP struct {
	ID     int    `json:"id"`
	Type   string `json:"type"`
	Source string `json:"source"`
	Lod    string `json:"lod"`
	New    string `json:"new"`
	Note   string `json:"note"`
	State  string `json:"state"`
	Ok     bool   `json:"ok"`
}

// isInterceptRuleType 与 UI ReplaceType「拦截请求」及 Config.ReplaceBodyInfo.IsBreak 一致。
func isInterceptRuleType(typ string) bool {
	return strings.TrimSpace(typ) == interceptRuleType
}

func isReplaceRuleType(typ string) bool {
	return strings.TrimSpace(typ) != "" && !isInterceptRuleType(typ)
}

func replaceRulesConvention() map[string]any {
	return map[string]any{
		"page": "请求拦截/数据替换设置（ReplaceBody）",
		"classification": map[string]any{
			"intercept": "type 必须为「拦截请求」",
			"replace":   "type 为 字符串(UTF8)/字符串(GBK)/Base64/十六进制 之一",
		},
		"intercept": map[string]any{
			"type":          interceptRuleType,
			"lodColumn":     "旧数据 = 要匹配的数据；在 source 范围内命中则进入 HTTP 断点",
			"newColumn":     "新数据 = 不填或占位「拦截请求[此项不用填写]」",
			"sourceColumn":  "查找范围 = 在哪些位置检测 lod（见 interceptSourceOptions）",
			"sourceOptions": interceptSourceOptions(),
		},
		"replace": map[string]any{
			"types":         []string{"字符串(UTF8)", "字符串(GBK)", "Base64", "十六进制"},
			"lodColumn":     "旧数据 = 被替换内容",
			"newColumn":     "新数据 = 替换结果",
			"sourceColumn":  "查找范围 = 在哪些位置检测旧数据并替换（见 replaceSourceOptions）",
			"sourceOptions": replaceSourceOptions(),
		},
	}
}

func interceptSourceOptions() []string {
	return []string{
		"HTTP请求/响应", "HTTP请求-全部", "HTTP响应-全部",
		"HTTP请求-URL", "HTTP请求-协议头", "HTTP请求-提交数据",
		"HTTP响应-协议头", "HTTP响应-响应数据",
	}
}

func replaceSourceOptions() []string {
	return []string{
		"任意", "Socket-任意",
		"TCP-全部", "UDP-全部", "Websocket-全部",
		"HTTP请求-全部", "HTTP请求-URL", "HTTP请求-协议头", "HTTP请求-提交数据",
		"HTTP响应-全部", "HTTP响应-协议头", "HTTP响应-响应数据",
		"TCP-发送", "TCP-接收", "UDP-发送", "UDP-接收",
		"Websocket-发送", "Websocket-接收",
	}
}

func ruleMCPEnriched(r replaceRuleMCP) map[string]any {
	intercept := isInterceptRuleType(r.Type)
	kind := "replace"
	lodHint := "旧数据：在 source 范围内出现则替换为新数据"
	newHint := "新数据：替换结果"
	if intercept {
		kind = "intercept"
		lodHint = "匹配数据：在 source 范围内出现则添加 HTTP 拦截断点"
		newHint = "拦截规则无需填写新数据（可留空或 UI 占位）"
	}
	return map[string]any{
		"id": r.ID, "type": r.Type, "source": r.Source,
		"lod": r.Lod, "new": r.New, "note": r.Note, "state": r.State, "ok": r.Ok,
		"ruleKind": kind,
		"lodHint":  lodHint,
		"newHint":  newHint,
	}
}

func replaceRulesSnapshot(app *AppMain, interceptOnly, replaceOnly bool) map[string]any {
	list := app.ReplaceBody.ReplaceBodyList()
	rules := make([]map[string]any, 0, len(list))
	for _, item := range list {
		if interceptOnly && !isInterceptRuleType(item.Type) {
			continue
		}
		if replaceOnly && !isReplaceRuleType(item.Type) {
			continue
		}
		rules = append(rules, ruleMCPEnriched(replaceRuleMCPFromInfo(item)))
	}
	return map[string]any{
		"rules":      rules,
		"total":      len(rules),
		"convention": replaceRulesConvention(),
	}
}

func parseReplaceRulesJSON(rulesJSON string) ([]replaceRuleMCP, error) {
	rulesJSON = strings.TrimSpace(rulesJSON)
	if rulesJSON == "" || rulesJSON == "[]" || rulesJSON == "{}" {
		return nil, nil
	}
	var arr []replaceRuleMCP
	if err := json.Unmarshal([]byte(rulesJSON), &arr); err == nil {
		return arr, nil
	}
	var obj map[string]replaceRuleMCP
	if err := json.Unmarshal([]byte(rulesJSON), &obj); err == nil {
		keys := make([]int, 0, len(obj))
		for k := range obj {
			if id, e := strconv.Atoi(k); e == nil {
				keys = append(keys, id)
			}
		}
		sort.Ints(keys)
		out := make([]replaceRuleMCP, 0, len(keys))
		for _, id := range keys {
			r := obj[strconv.Itoa(id)]
			if r.ID == 0 {
				r.ID = id
			}
			out = append(out, r)
		}
		return out, nil
	}
	return nil, errors.New("rulesJSON 须为规则数组 [{id,type,source,lod,new,note,state},...]")
}

func validateIncomingReplaceRules(rules []replaceRuleMCP, wantIntercept bool) error {
	for i, r := range rules {
		if wantIntercept {
			if strings.TrimSpace(r.Type) != "" && !isInterceptRuleType(r.Type) {
				return fmt.Errorf("rules[%d].type=%q：config_set_intercept 仅接受「拦截请求」", i, r.Type)
			}
		} else {
			if isInterceptRuleType(r.Type) {
				return fmt.Errorf("rules[%d].type=拦截请求：config_set_rewrite 不接受拦截类型，请用 config_set_intercept", i)
			}
		}
	}
	return nil
}

func applyReplaceRulesToConfig(app *AppMain, rules []replaceRuleMCP) {
	Config.Config.ReplaceRoles = make(map[int]*Config.ReplaceBodyInfo)
	for _, r := range rules {
		id := r.ID
		if id <= 0 {
			id = app.ReplaceBody.CreateReplaceBody()
		} else if Config.Config.ReplaceRoles[id] == nil {
			Config.Config.ReplaceRoles[id] = &Config.ReplaceBodyInfo{ID: id}
		}
		state := r.State
		if state == "" {
			state = "已启用"
		}
		typ := r.Type
		if isInterceptRuleType(typ) {
			if strings.TrimSpace(r.New) == "" {
				r.New = "拦截请求[此项不用填写]"
			}
		}
		app.ReplaceBody.ReplaceBodyUpdate(id, typ, r.Source, r.Lod, r.New, r.Note, state)
	}
	Config.Config.Save()
}

func replaceRuleMCPFromInfo(item Config.ReplaceBodyInfo) replaceRuleMCP {
	state := item.State
	if state == "" {
		if item.Ok {
			state = "已启用"
		} else {
			state = "已禁用"
		}
	}
	return replaceRuleMCP{
		ID: item.ID, Type: item.Type, Source: item.Source,
		Lod: item.Lod, New: item.New, Note: item.Note, State: state, Ok: item.Ok,
	}
}

func upsertReplaceRulesSubset(app *AppMain, incoming []replaceRuleMCP, interceptSubset bool, replaceAll bool) {
	current := app.ReplaceBody.ReplaceBodyList()
	others := make([]replaceRuleMCP, 0)
	subset := make(map[int]replaceRuleMCP)
	for _, item := range current {
		r := replaceRuleMCPFromInfo(item)
		if interceptSubset == isInterceptRuleType(item.Type) {
			subset[item.ID] = r
		} else {
			others = append(others, r)
		}
	}
	if replaceAll {
		subset = make(map[int]replaceRuleMCP)
	}
	var incomingNoID []replaceRuleMCP
	for _, r := range incoming {
		if interceptSubset {
			if strings.TrimSpace(r.Type) == "" {
				r.Type = interceptRuleType
			}
		}
		if r.State == "" {
			r.State = "已启用"
		}
		if r.ID <= 0 {
			incomingNoID = append(incomingNoID, r)
			continue
		}
		subset[r.ID] = r
	}
	merged := make([]replaceRuleMCP, 0, len(others)+len(subset)+len(incomingNoID))
	merged = append(merged, others...)
	for _, r := range subset {
		merged = append(merged, r)
	}
	merged = append(merged, incomingNoID...)
	sort.Slice(merged, func(i, j int) bool { return merged[i].ID < merged[j].ID })
	applyReplaceRulesToConfig(app, merged)
}

func bridgeConfigRuleSetState(app *AppMain, m map[string]any) (any, error) {
	id := argInt(m, "id", 0)
	if id <= 0 {
		return nil, errors.New("id 必填（规则 id，与 config_get 返回的 rules[].id 一致）")
	}
	state := strings.TrimSpace(argString(m, "state"))
	if state == "" {
		if argBool(m, "enabled", true) {
			state = "已启用"
		} else {
			state = "已禁用"
		}
	}
	if state != "已启用" && state != "已禁用" {
		return nil, errors.New("state 须为 已启用 或 已禁用")
	}
	obj := Config.Config.ReplaceRoles[id]
	if obj == nil {
		if Config.Config.ReplaceHost[id] != nil {
			return nil, fmt.Errorf("规则 id %d 属于 Host 映射，请用 config_host_delete 或 config_host_update", id)
		}
		return nil, fmt.Errorf("规则 id %d 不存在", id)
	}
	app.ReplaceBody.ReplaceBodyUpdate(id, obj.Type, obj.Source, obj.Lod, obj.New, obj.Note, state)
	emitMCPRulesPageReload()
	return map[string]any{"ok": true, "id": id, "state": state}, nil
}
