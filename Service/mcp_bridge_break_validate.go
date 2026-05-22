package Service

import (
	"changeme/Service/Session"
	"fmt"
	"strings"

	"github.com/qtgolang/SunnyNet/src/public"
)

const (
	breakPhaseNone       = 0
	breakPhaseUpstream   = 1 // 拦截上行：HttpSendRequest + IsWait
	breakPhaseDownstream = 2 // 拦截下行：HttpResponseOK + IsWait
)

// httpBreakPhase 返回是否处于断点等待及阶段（0=非拦截/未知，1=上行，2=下行）。
func httpBreakPhase(obj *Session.HttpSession) (waiting bool, phase int) {
	if obj == nil || !obj.IsWait() {
		return false, breakPhaseNone
	}
	switch obj.State {
	case public.HttpSendRequest:
		return true, breakPhaseUpstream
	case public.HttpResponseOK:
		return true, breakPhaseDownstream
	default:
		return true, breakPhaseNone
	}
}

func requireBreakUpstream(obj *Session.HttpSession, th int) error {
	if obj == nil {
		return fmt.Errorf("theology %d 不存在", th)
	}
	waiting, phase := httpBreakPhase(obj)
	if !waiting {
		return fmt.Errorf("theology %d 非拦截等待状态，无法改请求", th)
	}
	if phase != breakPhaseUpstream {
		return fmt.Errorf("theology %d 非拦截上行模式，无法改请求", th)
	}
	return nil
}

func requireBreakDownstream(obj *Session.HttpSession, th int) error {
	if obj == nil {
		return fmt.Errorf("theology %d 不存在", th)
	}
	waiting, phase := httpBreakPhase(obj)
	if !waiting {
		return fmt.Errorf("theology %d 非拦截等待状态，无法改响应", th)
	}
	if phase != breakPhaseDownstream {
		return fmt.Errorf("theology %d 非拦截下行模式，无法改响应", th)
	}
	return nil
}

func hasBreakSyncRequestPatch(m map[string]any) bool {
	return strings.TrimSpace(argString(m, "requestURL")) != "" ||
		strings.TrimSpace(argString(m, "headersJSON")) != "" ||
		strings.TrimSpace(argString(m, "bodyB64")) != ""
}

func hasBreakSyncResponsePatch(m map[string]any) bool {
	if _, ok := m["statusCode"]; ok {
		return true
	}
	return strings.TrimSpace(argString(m, "headersJSON")) != "" ||
		strings.TrimSpace(argString(m, "bodyB64")) != ""
}

func rejectUpstreamMethodChange(m map[string]any, current string) error {
	method := strings.TrimSpace(argString(m, "method"))
	if method == "" {
		method = strings.TrimSpace(argString(m, "requestMethod"))
	}
	if method != "" && !strings.EqualFold(method, current) {
		return fmt.Errorf("拦截上行模式下不可修改 Method（当前为 %s）", current)
	}
	return nil
}
