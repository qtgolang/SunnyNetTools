package mcp

import (
	"changeme/Service/Config"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"hash/fnv"
	"sync"
	"time"
)

const (
	defaultRPCTime = 30 * time.Second // WebView 往返默认超时
)

// MCPMsg 与前端 WebView 通信的消息结构
type McpMsg struct {
	Page string `json:"page"` // 页面名（例如 Main）
	Tag  string `json:"tag"`  // 消息类型（例如 home/SystemProxy）
	Msg  string `json:"msg"`  // 请求内容
	Res  string `json:"res"`  // 回包内容
	Id   uint32 `json:"id"`   // 请求ID（用于匹配请求/回包）
}

// MsgCallback 由 AppMain.McpFuncRes 调用，将前端回包投递到等待中的 callWebviewMsg。
var MsgCallback func(McpMsg)

var (
	waitMu  sync.Mutex
	waitMap = make(map[uint32]chan McpMsg) // id -> 等待通道
)

func init() {
	MsgCallback = mcpCall
}

// callWebviewMsg 发消息给前端并等待回包（带超时）
//
// 核心思路：
// 1) 根据 Page/Tag/Msg + 时间生成稳定的 requestId
// 2) 把 requestId -> chan 放进 waitMap
// 3) EmitEvent 通知前端
// 4) 等待回包 / 超时 / ctx 取消
func callWebviewMsg(ctx context.Context, obj McpMsg) string {
	seed := fmt.Sprintf("%d|%s|%s|%s", time.Now().UnixNano(), obj.Tag, obj.Page, obj.Msg)

	hash := fnv.New32a()
	_, _ = hash.Write([]byte(seed))
	obj.Id = hash.Sum32()

	ch := make(chan McpMsg, 1)

	waitMu.Lock()
	waitMap[obj.Id] = ch
	waitMu.Unlock()

	// 退出时一定要把等待项清理掉，避免内存泄漏
	defer func() {
		waitMu.Lock()
		delete(waitMap, obj.Id)
		waitMu.Unlock()
	}()

	// 这里依赖你的项目：向前端派发事件
	Config.AppList["Main"].EmitEvent("mcp", obj)

	// ctx 有 deadline 就用 deadline；否则走默认超时
	timeout := defaultRPCTime
	if dl, ok := ctx.Deadline(); ok {
		if remain := time.Until(dl); remain > 0 {
			timeout = remain
		}
	}

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case res := <-ch:
		return res.Res
	case <-timer.C:
		return "timeout"
	case <-ctx.Done():
		return "timeout"
	}
}

// mcpCall 前端回包入口：把回包投递到对应等待通道
func mcpCall(res McpMsg) {
	waitMu.Lock()
	ch, ok := waitMap[res.Id]
	waitMu.Unlock()
	if !ok {
		return
	}
	select {
	case ch <- res:
	default:
		// 如果对方已经不等了，就直接丢弃
	}
}

// SetFilter 设置主列表 ag-grid 过滤器（filter JSON 字符串）。
func (m McpMsg) SetFilter(filterJSON string) string {
	return callWebviewMsg(context.Background(), McpMsg{Page: "main", Tag: "filter", Msg: filterJSON})
}

// DeleteID 按主列表「序号」(listIndex) 删除：先映射为 theology，再下发给前端 delreq。
func (m McpMsg) DeleteID(ids ...int) string {
	if len(ids) == 0 {
		return "没有传入id"
	}
	idxMap, err := ListIndexesToTheologies(ids...)
	if err != nil {
		return callWebviewMsg(context.Background(), McpMsg{Page: "main", Tag: "delreq", Msg: "[]"})
	}
	theologies := make([]int, 0, len(ids))
	for _, id := range ids {
		if th := idxMap[strconv.Itoa(id)]; th != 0 {
			theologies = append(theologies, th)
		}
	}
	b, err := json.Marshal(theologies)
	if err != nil {
		return callWebviewMsg(context.Background(), McpMsg{Page: "main", Tag: "delreq", Msg: "[]"})
	}
	return callWebviewMsg(context.Background(), McpMsg{Page: "main", Tag: "delreq", Msg: string(b)})
}

// ------------------- 前端 UI 同步（可选，与桥接 op 配合） -------------------

func (m McpMsg) homeStartCapture() string {
	return callWebviewMsg(context.Background(), McpMsg{Page: "Main", Tag: "home", Msg: "start"})
}
func (m McpMsg) homeStopCapture() string {
	return callWebviewMsg(context.Background(), McpMsg{Page: "Main", Tag: "home", Msg: "stop"})
}
func (m McpMsg) homeCaptureState() string {
	return callWebviewMsg(context.Background(), McpMsg{Page: "Main", Tag: "home", Msg: "state"})
}
func (m McpMsg) homeClearAll() string {
	return callWebviewMsg(context.Background(), McpMsg{Page: "Main", Tag: "home", Msg: "clear"})
}
func (m McpMsg) proxyEnable() string {
	return callWebviewMsg(context.Background(), McpMsg{Page: "Main", Tag: "SystemProxy", Msg: "Set"})
}
func (m McpMsg) proxyDisable() string {
	return callWebviewMsg(context.Background(), McpMsg{Page: "Main", Tag: "SystemProxy", Msg: "Cancel"})
}
func (m McpMsg) proxyState() string {
	return callWebviewMsg(context.Background(), McpMsg{Page: "Main", Tag: "SystemProxy", Msg: "State"})
}
func (m McpMsg) themeSetDark() string {
	return callWebviewMsg(context.Background(), McpMsg{Page: "theme", Tag: "theme", Msg: "setDark"})
}
func (m McpMsg) themeSetLight() string {
	return callWebviewMsg(context.Background(), McpMsg{Page: "theme", Tag: "theme", Msg: "setLight"})
}
func (m McpMsg) themeState() string {
	return callWebviewMsg(context.Background(), McpMsg{Page: "theme", Tag: "theme", Msg: "state"})
}
