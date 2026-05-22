package Service

import (
	"bytes"
	"changeme/Service/Config"
	"changeme/Service/Session"
	"changeme/Service/mcp"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/qtgolang/SunnyNet/src/http"
	"github.com/qtgolang/SunnyNet/src/public"
)

const mcpBodyChunkMax = 4 << 20

func theologyIDsToRowIDs(ids []int) []string {
	out := make([]string, len(ids))
	for i, th := range ids {
		out[i] = strconv.Itoa(th)
	}
	return out
}

func bridgeEngineApplyAdvanced(app *AppMain, m map[string]any) (any, error) {
	payload := strings.TrimSpace(argString(m, "payload"))
	if payload == "" {
		return nil, errors.New("payload 为空")
	}
	var p map[string]any
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		return nil, err
	}
	if v, ok := p["port"]; ok {
		port := intFromAny(v)
		if port > 0 {
			if msg := app.SetPort(port, false); msg != "" {
				return nil, errors.New(msg)
			}
		}
	}
	if v, ok := p["disableTCP"]; ok {
		app.SetDisableTCP(boolFromAny(v))
	}
	if v, ok := p["disableUDP"]; ok {
		app.SetDisableUDP(boolFromAny(v))
	}
	if v, ok := p["disableCache"]; ok {
		app.SetDisableCache(boolFromAny(v))
	}
	if v, ok := p["openAuthMode"]; ok {
		app.SetAuthMode(boolFromAny(v))
	}
	if v, ok := p["limitRequestSize"]; ok {
		app.SetLimitRequestSize(intFromAny(v))
	}
	if v, ok := p["proxyRoles"]; ok {
		if s, ok := v.(string); ok {
			app.SetProxyRoles(s)
		}
	}
	if v, ok := p["proxyDns"]; ok {
		if s, ok := v.(string); ok {
			app.SetProxyDns(s)
		}
	}
	if v, ok := p["mustTcpRoles"]; ok {
		if s, ok := v.(string); ok {
			tp := intFromAny(p["mustTcpType"])
			app.SetMustTcpRoles(tp, s)
		}
	}
	if v, ok := p["outRouter"]; ok {
		if s, ok := v.(string); ok {
			app.SetInterfaceOutRouterAdders(s)
		}
	}
	if v, ok := p["httpsProto"]; ok {
		if s, ok := v.(string); ok {
			app.SetHTTPSProto(s)
		}
	}
	return settingsApplyResult(app, map[string]any{"source": "engine_apply_advanced"}), nil
}

func boolFromAny(v any) bool {
	switch t := v.(type) {
	case bool:
		return t
	case float64:
		return t != 0
	case string:
		return strings.EqualFold(strings.TrimSpace(t), "true") || t == "1"
	default:
		return false
	}
}

func intFromAny(v any) int {
	switch t := v.(type) {
	case float64:
		return int(t)
	case int:
		return t
	case int64:
		return int(t)
	case string:
		n, _ := strconv.Atoi(strings.TrimSpace(t))
		return n
	default:
		return 0
	}
}

func bridgeMainSearch(app *AppMain, m map[string]any) (any, error) {
	needle := strings.TrimSpace(argString(m, "needle"))
	if needle == "" {
		needle = strings.TrimSpace(argString(m, "queryText"))
	}
	if needle == "" {
		return nil, errors.New("needle 或 queryText 必填")
	}
	searchType := strings.TrimSpace(argString(m, "searchType"))
	if searchType == "" {
		searchType = "auto"
	}
	searchRange := strings.TrimSpace(argString(m, "searchRange"))
	if searchRange == "" {
		searchRange = "在全部范围寻找"
	}
	opts := ""
	if argBool(m, "caseInsensitive", true) {
		opts += "不区分大小写 "
	}
	if argBool(m, "stripSpaces", false) {
		opts += "删除空格后搜索 "
	}
	fi := &FindInfo{
		Theology: 0,
		Value:    needle,
		Type:     searchType,
		Range:    searchRange,
		Options:  strings.TrimSpace(opts),
	}
	ids := app.FindSession(fi)
	if len(ids) == 0 {
		return nil, errors.New("未找到匹配会话")
	}
	if len(ids) > 1 {
		return map[string]any{
			"ok":         false,
			"unique":     false,
			"count":      len(ids),
			"theologies": ids,
		}, nil
	}
	th := ids[0]
	listIndex := theologyListIndex(th)
	if listIndex <= 0 {
		return nil, errors.New("匹配会话不在当前主列表可见范围")
	}
	return map[string]any{
		"ok":        true,
		"unique":    true,
		"theology":  th,
		"listIndex": listIndex,
		"rowId":     strconv.Itoa(th),
	}, nil
}

func bridgeStreamCount(m map[string]any) (any, error) {
	th, err := streamKeyTheology(m)
	if err != nil {
		return nil, err
	}
	return map[string]any{"total": sessionStreamTotal(th)}, nil
}

func bridgeStreamSlice(m map[string]any) (any, error) {
	th, err := streamKeyTheology(m)
	if err != nil {
		return nil, err
	}
	offset := argInt(m, "offset", 0)
	limit := argInt(m, "limit", 50)
	if limit <= 0 {
		limit = 50
	}
	if limit > 500 {
		limit = 500
	}
	keys := sessionStreamKeys(th)
	total := len(keys)
	end := offset + limit
	if end > total {
		end = total
	}
	if offset > total {
		offset = total
	}
	page := keys[offset:end]
	rows := make([]map[string]any, len(page))
	for i, mid := range page {
		rows[i] = map[string]any{
			"messageId": mid,
			"cells":     streamMessageCells(th, mid),
		}
	}
	return map[string]any{
		"streamKey": strconv.Itoa(th),
		"rows":      rows,
		"total":     total,
		"offset":    offset,
		"limit":     limit,
	}, nil
}

func streamKeyTheology(m map[string]any) (int, error) {
	key := strings.TrimSpace(argString(m, "streamKey"))
	if key == "" {
		th, err := argTheologyOne(m)
		return th, err
	}
	th, _ := strconv.Atoi(key)
	if th <= 0 {
		return 0, errors.New("streamKey 无效")
	}
	return th, nil
}

func sessionStreamKeys(th int) []int {
	var keys []int
	if hs := Session.GetHttpSession(th); hs != nil {
		keys = append(keys, hs.WebsocketStreamKeys...)
	} else if ts := Session.GetTCPSession(th); ts != nil {
		keys = append(keys, ts.StreamKeys...)
	} else if us := Session.GetUDPSession(th); us != nil {
		keys = append(keys, us.StreamKeys...)
	}
	sort.Ints(keys)
	return keys
}

func sessionStreamTotal(th int) int {
	return len(sessionStreamKeys(th))
}

func streamMessageCells(th, mid int) []string {
	if hs := Session.GetHttpSession(th); hs != nil {
		st := hs.WebsocketStream[mid]
		if st == nil {
			return []string{strconv.Itoa(mid), "", "", "", ""}
		}
		dir := "接收"
		if st.IsSend {
			dir = "发送"
		}
		ft := "Text"
		if st.WebsocketType != 1 && st.WebsocketType != 0 {
			ft = "Binary"
		}
		return []string{strconv.Itoa(mid), dir, ft, st.Time, strconv.Itoa(len(st.Body))}
	}
	if ts := Session.GetTCPSession(th); ts != nil {
		st := ts.Stream[mid]
		if st == nil {
			return []string{strconv.Itoa(mid), "", "", "", ""}
		}
		dir := "接收"
		if st.IsSend {
			dir = "发送"
		}
		return []string{strconv.Itoa(mid), dir, "", st.Time, strconv.Itoa(len(st.Body))}
	}
	if us := Session.GetUDPSession(th); us != nil {
		st := us.Stream[mid]
		if st == nil {
			return []string{strconv.Itoa(mid), "", "", "", ""}
		}
		dir := "接收"
		if st.IsSend {
			dir = "发送"
		}
		return []string{strconv.Itoa(mid), dir, "", st.Time, strconv.Itoa(len(st.Body))}
	}
	return []string{strconv.Itoa(mid), "", "", "", ""}
}

func bridgeUiThemeGet() (any, error) {
	disp := "light"
	if Config.Config.IsDark {
		disp = "dark"
	}
	themeLocal := map[string]any{
		"sunnynet-theme": disp,
		"isDark":         Config.Config.IsDark,
		"listColors":     Config.Config.ListColor,
	}
	return map[string]any{
		"themeLocal":  themeLocal,
		"displayMode": disp,
		"dataTheme":   disp,
		"colorHints": map[string]string{
			"sunnynet-theme": "dark/light，与主窗口明暗一致",
			"listColors":     "Config.ListColor 行配色键",
		},
	}, nil
}

func bridgeUiThemeSet(m map[string]any) (any, error) {
	mode := strings.TrimSpace(strings.ToLower(argString(m, "mode")))
	if mode != "dark" && mode != "light" && mode != "toggle" {
		return nil, errors.New("mode 须为 dark、light 或 toggle")
	}
	if mode == "toggle" {
		if Config.Config.IsDark {
			mode = "light"
		} else {
			mode = "dark"
		}
	}
	if mcpApp == nil {
		return nil, errors.New("mcp: 主程序未初始化")
	}
	var res string
	if mode == "dark" {
		res = mcp.ThemeSetDark()
		mcpApp.SetIsDark(true)
	} else {
		res = mcp.ThemeSetLight()
		mcpApp.SetIsDark(false)
	}
	if res == "timeout" {
		return nil, errors.New("主题切换超时")
	}
	out, err := bridgeUiThemeGet()
	if err != nil {
		return nil, err
	}
	outMap := out.(map[string]any)
	outMap["ok"] = true
	outMap["webview"] = res
	return outMap, nil
}

func bridgeHTTPGetPartMulti(app *AppMain, m map[string]any) (any, error) {
	ids, err := argTheologyList(m)
	if err != nil {
		return nil, err
	}
	part := strings.TrimSpace(argString(m, "part"))
	if part == "" {
		part = "responseBody"
	}
	typ := strings.TrimSpace(argString(m, "type"))
	if typ == "" {
		typ = "auto"
	}
	offset := argInt(m, "offset", 0)
	maxLen := argInt(m, "maxLen", 0)
	if maxLen <= 0 {
		maxLen = mcpBodyChunkMax
	}
	one := func(th int) (map[string]any, error) {
		bs, frameType, e := httpPartBytes(app, th, part, maxLen)
		if e != nil {
			return nil, e
		}
		if bs == nil {
			return map[string]any{"ok": false, "total": 0, "data": "", "theology": th}, nil
		}
		total := len(bs)
		if offset > total {
			offset = total
		}
		chunk := bs[offset:]
		if len(chunk) > maxLen {
			chunk = chunk[:maxLen]
		}
		data, outType := mcp.EncodeDataPart(chunk, typ)
		out := map[string]any{
			"ok": true, "total": total, "type": outType, "data": data, "theology": th,
		}
		if frameType != "" {
			out["frameType"] = frameType
		}
		return out, nil
	}
	if len(ids) == 1 {
		return one(ids[0])
	}
	items := make([]map[string]any, 0, len(ids))
	for _, th := range ids {
		item, e := one(th)
		if e != nil {
			return nil, e
		}
		items = append(items, item)
	}
	return map[string]any{"items": items}, nil
}

func httpPartBytes(app *AppMain, th int, part string, maxLen int) ([]byte, string, error) {
	getAll := maxLen <= 0
	switch part {
	case "requestBody":
		return app.GetHTTPRequestBody(th, getAll), "", nil
	case "responseBody":
		return app.GetHTTPResponseBody(th, getAll), "", nil
	case "rawRequest":
		return buildHTTPRaw(th, true), "", nil
	case "rawResponse":
		return buildHTTPRaw(th, false), "", nil
	default:
		return nil, "", errors.New("part 须为 requestBody | responseBody | rawRequest | rawResponse")
	}
}

func buildHTTPRaw(th int, isRequest bool) []byte {
	hs := Session.GetHttpSession(th)
	if hs == nil {
		return nil
	}
	var buf bytes.Buffer
	hs.ReadLocked(func() {
		if isRequest {
			buf.WriteString(fmt.Sprintf("%s %s %s\r\n", hs.Request.Method, hs.Request.Url, hs.Request.Proto))
			writeHeader(&buf, hs.Request.Header)
			buf.WriteString("\r\n")
			buf.Write(hs.Request.Body)
		} else {
			buf.WriteString(fmt.Sprintf("%s %s %s\r\n", hs.Response.Proto, hs.Response.Code, hs.Response.State))
			writeHeader(&buf, hs.Response.Header)
			buf.WriteString("\r\n")
			buf.Write(hs.Response.Body)
		}
	})
	return buf.Bytes()
}

func writeHeader(buf *bytes.Buffer, hdr http.Header) {
	for k, vals := range hdr {
		for _, v := range vals {
			buf.WriteString(k + ": " + v + "\r\n")
		}
	}
}

func bridgeStreamGetPartMulti(app *AppMain, m map[string]any, defaultType string) (any, error) {
	ids, err := argTheologyList(m)
	if err != nil {
		return nil, err
	}
	msgID := argInt(m, "messageId", 0)
	typ := strings.TrimSpace(argString(m, "type"))
	if typ == "" {
		typ = defaultType
	}
	offset := argInt(m, "offset", 0)
	maxLen := argInt(m, "maxLen", 0)
	if maxLen <= 0 {
		maxLen = mcpBodyChunkMax
	}
	one := func(th int) (map[string]any, error) {
		bs := app.GetSessionMessageBody(th, msgID)
		if bs == nil {
			return map[string]any{"ok": false, "total": 0, "data": "", "theology": th}, nil
		}
		total := len(bs)
		if offset > total {
			offset = total
		}
		chunk := bs[offset:]
		if len(chunk) > maxLen {
			chunk = chunk[:maxLen]
		}
		data, outType := mcp.EncodeDataPart(chunk, typ)
		out := map[string]any{
			"ok": true, "total": total, "type": outType, "data": data, "theology": th,
		}
		if hs := Session.GetHttpSession(th); hs != nil && hs.IsWebsocketRequest {
			if st := hs.WebsocketStream[msgID]; st != nil {
				if st.WebsocketType == 1 {
					out["frameType"] = "Text"
				} else {
					out["frameType"] = "Binary"
				}
			}
		}
		return out, nil
	}
	if len(ids) == 1 {
		return one(ids[0])
	}
	items := make([]map[string]any, 0, len(ids))
	for _, th := range ids {
		item, e := one(th)
		if e != nil {
			return nil, e
		}
		items = append(items, item)
	}
	return map[string]any{"items": items}, nil
}

func sessionJSONDetail(th int) (any, error) {
	if hs := Session.GetHttpSession(th); hs != nil {
		var out map[string]any
		hs.ReadLocked(func() {
			intercept := "非拦截"
			if hs.IsWait() {
				switch hs.State {
				case public.HttpSendRequest:
					intercept = "上行"
				case public.HttpResponseOK:
					intercept = "下行"
				default:
					intercept = "拦截"
				}
			}
			kind := "http"
			if hs.IsWebsocketRequest {
				kind = "ws"
			}
			bm := uint32(0)
			if waiting := hs.IsWait(); waiting {
				switch hs.State {
				case public.HttpSendRequest:
					bm = 1
				case public.HttpResponseOK:
					bm = 2
				default:
					bm = 1
				}
			}
			out = map[string]any{
				"theology":       th,
				"rowId":          strconv.Itoa(th),
				"type":           kind,
				"breakMode":      bm,
				"isWaiting":      hs.IsWait(),
				"interceptState": intercept,
				"state":          hs.State,
				"method":         hs.Request.Method,
				"url":            hs.Request.Url,
				"code":           hs.Response.Code,
				"note":           hs.GetNote(),
				"processName":    hs.Request.ProcessName,
				"recLength":      hs.RecLength,
				"senLength":      hs.SenLength,
				"disconnect":     hs.WebsocketDisconnect,
			}
		})
		b, err := json.Marshal(out)
		return json.RawMessage(b), err
	}
	if ts := Session.GetTCPSession(th); ts != nil {
		out := map[string]any{
			"theology": th, "rowId": strconv.Itoa(th), "type": "tcp",
			"method": ts.GetMethod(), "remote": ts.RemoteAddress,
			"note": ts.GetNote(), "disconnect": ts.Disconnect,
			"recLength": ts.RecLength, "senLength": ts.SenLength,
		}
		b, err := json.Marshal(out)
		return json.RawMessage(b), err
	}
	if us := Session.GetUDPSession(th); us != nil {
		out := map[string]any{
			"theology": th, "rowId": strconv.Itoa(th), "type": "udp",
			"remote": us.RemoteAddress, "note": us.GetNote(),
			"recLength": us.RecLength, "senLength": us.SenLength,
		}
		b, err := json.Marshal(out)
		return json.RawMessage(b), err
	}
	return nil, fmt.Errorf("theology %d 不存在", th)
}
