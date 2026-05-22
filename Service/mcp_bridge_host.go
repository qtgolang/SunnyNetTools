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

type hostRuleMCP struct {
	ID   int    `json:"id"`
	Lod  string `json:"lod"`
	New  string `json:"new"`
	Note string `json:"note"`
}

func hostRuleLod(r hostRuleMCP) string {
	if s := strings.TrimSpace(r.Lod); s != "" {
		return s
	}
	return ""
}

func hostRuleNew(r hostRuleMCP) string {
	return strings.TrimSpace(r.New)
}

func hostRuleNote(r hostRuleMCP) string {
	return strings.TrimSpace(r.Note)
}

func parseHostRulesJSON(rulesJSON string) ([]hostRuleMCP, error) {
	rulesJSON = strings.TrimSpace(rulesJSON)
	if rulesJSON == "" || rulesJSON == "[]" || rulesJSON == "{}" {
		return nil, nil
	}
	var arr []hostRuleMCP
	if err := json.Unmarshal([]byte(rulesJSON), &arr); err == nil {
		return normalizeHostRulesArray(arr), nil
	}
	// 兼容 UI/配置：{"3":{...},"8":{...}}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal([]byte(rulesJSON), &raw); err != nil {
		return nil, errors.New("rulesJSON 须为 Host 规则数组 [{id,lod,new,note},...] 或对象 {\"id\":{...}}")
	}
	keys := make([]int, 0, len(raw))
	for k := range raw {
		if id, e := strconv.Atoi(k); e == nil {
			keys = append(keys, id)
		}
	}
	sort.Ints(keys)
	out := make([]hostRuleMCP, 0, len(keys))
	for _, id := range keys {
		var r hostRuleMCP
		if err := json.Unmarshal(raw[strconv.Itoa(id)], &r); err != nil {
			return nil, err
		}
		if r.ID == 0 {
			r.ID = id
		}
		out = append(out, r)
	}
	return normalizeHostRulesArray(out), nil
}

// normalizeHostRulesArray 兼容 Lod/New/Note 大写字段名。
func normalizeHostRulesArray(arr []hostRuleMCP) []hostRuleMCP {
	var raw []map[string]any
	b, _ := json.Marshal(arr)
	_ = json.Unmarshal(b, &raw)
	out := make([]hostRuleMCP, 0, len(arr))
	for i, r := range arr {
		if i < len(raw) {
			if v, ok := raw[i]["Lod"].(string); ok && r.Lod == "" {
				r.Lod = v
			}
			if v, ok := raw[i]["New"].(string); ok && r.New == "" {
				r.New = v
			}
			if v, ok := raw[i]["Note"].(string); ok && r.Note == "" {
				r.Note = v
			}
			if v, ok := raw[i]["ID"].(float64); ok && r.ID == 0 {
				r.ID = int(v)
			}
		}
		out = append(out, r)
	}
	return out
}

func hostRulesSnapshot(app *AppMain) map[string]any {
	list := app.ReplaceHostList()
	rules := make([]map[string]any, 0, len(list))
	for _, h := range list {
		rules = append(rules, map[string]any{
			"id": h.ID, "lod": h.Lod, "new": h.New, "note": h.Note,
		})
	}
	return map[string]any{
		"rules": rules,
		"total": len(rules),
		"convention": map[string]any{
			"lod":  "旧 Host（匹配来源）",
			"new":  "新 Host（替换目标，可含端口如 host:8443）",
			"note": "注释",
			"deleteOp": "config_host_delete（勿用 config_rule_set_state，该 op 仅用于数据替换/拦截规则）",
		},
	}
}

func applyHostRulesFull(app *AppMain, rules []hostRuleMCP) error {
	Config.Config.ReplaceHost = make(map[int]*Config.ReplaceHostInfo)
	for _, r := range rules {
		lod := hostRuleLod(r)
		newH := hostRuleNew(r)
		if lod == "" || newH == "" {
			return fmt.Errorf("规则 id=%d 缺少 lod 或 new", r.ID)
		}
		id := r.ID
		if id <= 0 {
			id = app.CreateReplaceHost()
		} else {
			Config.Config.ReplaceHost[id] = &Config.ReplaceHostInfo{ID: id}
			if id > Config.CertID {
				Config.CertID = id
			}
		}
		if !app.ReplaceHostUpdate(id, lod, newH, hostRuleNote(r)) {
			return fmt.Errorf("规则 id=%d 更新失败", id)
		}
	}
	Config.Config.Save()
	return nil
}

func findHostRuleIDByLod(lod string) (int, bool) {
	lod = strings.TrimSpace(lod)
	if lod == "" {
		return 0, false
	}
	for id, v := range Config.Config.ReplaceHost {
		if v == nil {
			continue
		}
		if v.Lod == lod || v.LodInfo.Host == lod {
			return id, true
		}
		h, _, _ := strings.Cut(lod, ":")
		if h != "" && v.LodInfo.Host == h {
			return id, true
		}
	}
	return 0, false
}

func bridgeHostAdd(app *AppMain, m map[string]any) (any, error) {
	lod := strings.TrimSpace(argString(m, "lod"))
	if lod == "" {
		lod = strings.TrimSpace(argString(m, "old"))
	}
	newH := strings.TrimSpace(argString(m, "new"))
	if lod == "" || newH == "" {
		return nil, errors.New("lod 与 new 必填")
	}
	note := strings.TrimSpace(argString(m, "note"))
	id := argInt(m, "id", 0)
	if id <= 0 {
		id = app.CreateReplaceHost()
	} else if Config.Config.ReplaceHost[id] == nil {
		Config.Config.ReplaceHost[id] = &Config.ReplaceHostInfo{ID: id}
		if id > Config.CertID {
			Config.CertID = id
		}
	}
	if !app.ReplaceHostUpdate(id, lod, newH, note) {
		return nil, fmt.Errorf("Host 规则 id=%d 写入失败", id)
	}
	emitMCPConfigReload("host")
	return map[string]any{"ok": true, "id": id, "lod": lod, "new": newH, "note": note}, nil
}

func bridgeHostDelete(app *AppMain, m map[string]any) (any, error) {
	id := argInt(m, "id", 0)
	if id <= 0 {
		lod := strings.TrimSpace(argString(m, "lod"))
		if lod == "" {
			lod = strings.TrimSpace(argString(m, "old"))
		}
		if lod == "" {
			return nil, errors.New("请提供 id 或 lod（旧 Host）")
		}
		var ok bool
		id, ok = findHostRuleIDByLod(lod)
		if !ok {
			return nil, fmt.Errorf("未找到 lod=%q 的 Host 规则", lod)
		}
	}
	if Config.Config.ReplaceHost[id] == nil {
		return nil, fmt.Errorf("Host 规则 id %d 不存在（config_rule_set_state 仅用于数据替换/拦截，删除 Host 请用 config_host_delete）", id)
	}
	app.ReplaceHostRemove(id)
	emitMCPConfigReload("host")
	return map[string]any{"ok": true, "id": id}, nil
}

func bridgeHostUpdate(app *AppMain, m map[string]any) (any, error) {
	id := argInt(m, "id", 0)
	if id <= 0 {
		return nil, errors.New("id 必填")
	}
	obj := Config.Config.ReplaceHost[id]
	if obj == nil {
		return nil, fmt.Errorf("Host 规则 id %d 不存在", id)
	}
	lod := obj.Lod
	if v, ok := m["lod"]; ok {
		lod = strings.TrimSpace(fmt.Sprint(v))
	} else if v, ok := m["old"]; ok {
		lod = strings.TrimSpace(fmt.Sprint(v))
	}
	newH := obj.New
	if v, ok := m["new"]; ok {
		newH = strings.TrimSpace(fmt.Sprint(v))
	}
	note := obj.Note
	if v, ok := m["note"]; ok {
		note = strings.TrimSpace(fmt.Sprint(v))
	}
	if lod == "" || newH == "" {
		return nil, errors.New("lod 与 new 不能为空")
	}
	if !app.ReplaceHostUpdate(id, lod, newH, note) {
		return nil, fmt.Errorf("Host 规则 id=%d 更新失败", id)
	}
	emitMCPConfigReload("host")
	return map[string]any{"ok": true, "id": id, "lod": lod, "new": newH, "note": note}, nil
}
