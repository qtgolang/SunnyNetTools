package mcpbridge

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"changeme/Service/mcpcatalog"
	"changeme/Service/mcpwire"
)

// httpBridge 本机 HTTP：REST 调用 Host.Invoke + SSE 推送 liststore 变更（与 TCP 行协议语义一致）。
type httpBridge struct {
	host *Host
	hub  *sseHub
	srv  *http.Server
	addr string
	mu   sync.Mutex
	ln   net.Listener
}

type sseHub struct {
	mu   sync.Mutex
	subs map[chan string]struct{}
}

func newSSEHub() *sseHub {
	return &sseHub{subs: make(map[chan string]struct{})}
}

func (h *sseHub) subscribe(buf int) chan string {
	ch := make(chan string, buf)
	h.mu.Lock()
	h.subs[ch] = struct{}{}
	h.mu.Unlock()
	return ch
}

func (h *sseHub) unsubscribe(ch chan string) {
	h.mu.Lock()
	if _, ok := h.subs[ch]; ok {
		delete(h.subs, ch)
		close(ch)
	}
	h.mu.Unlock()
}

func (h *sseHub) broadcast(msg string) {
	if h == nil || msg == "" {
		return
	}
	h.mu.Lock()
	for ch := range h.subs {
		select {
		case ch <- msg:
		default:
			// 慢消费者丢事件，避免阻塞主路径
		}
	}
	h.mu.Unlock()
}

var (
	httpFanoutMu sync.Mutex
	httpFanout   func(string) // MCP HTTP 启用时由 httpBridge 注册；禁用时清空
)

// RouteHook 在默认 MCP 路由注册后追加自定义路由（如 /doc、旧版 JSON-RPC）。
var RouteHook func(mux *http.ServeMux)

func registerHTTPListstoreFanout(fn func(string)) {
	httpFanoutMu.Lock()
	httpFanout = fn
	httpFanoutMu.Unlock()
}

// NotifyListstoreMCPHTTPFanout 由 main 在 liststore 变更时调用；无 HTTP 桥或未启用时为 no-op。
func NotifyListstoreMCPHTTPFanout(payload string) {
	httpFanoutMu.Lock()
	f := httpFanout
	httpFanoutMu.Unlock()
	if f != nil {
		f(payload)
	}
}

func newHTTPBridge(host *Host) *httpBridge {
	return &httpBridge{host: host, hub: newSSEHub()}
}

func (b *httpBridge) broadcast(payload string) {
	if b != nil && b.hub != nil {
		b.hub.broadcast(payload)
	}
}

func (b *httpBridge) Addr() string {
	if b == nil {
		return ""
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.addr
}

func (b *httpBridge) Start(host string, port int) error {
	if b == nil || b.host == nil {
		return errors.New("mcp http: bridge or host nil")
	}
	if port <= 0 || port > 65535 {
		return errors.New("mcp http: invalid port")
	}
	if host == "" {
		host = "127.0.0.1"
	}
	addr := fmt.Sprintf("%s:%d", host, port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	mux := http.NewServeMux()
	mux.HandleFunc(SunnyNetHTTPAPIPrefix+"/health", b.handleHealth)
	mux.HandleFunc(SunnyNetHTTPAPIPrefix+"/supported-ops", b.handleSupportedOps)
	mux.HandleFunc(SunnyNetHTTPAPIPrefix+"/invoke", b.handleInvoke)
	mux.HandleFunc(SunnyNetHTTPAPIPrefix+"/events", b.handleSSE)
	mux.Handle(MCPStreamablePath, newStreamableMCPHandler(b.host))
	if RouteHook != nil {
		RouteHook(mux)
	}

	srv := &http.Server{
		Handler:      withCORS(mux),
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 0, // SSE 长连接
	}
	b.mu.Lock()
	b.ln = ln
	b.srv = srv
	b.addr = ln.Addr().String()
	b.mu.Unlock()

	go func() {
		if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
		}
	}()
	registerHTTPListstoreFanout(b.broadcast)
	return nil
}

func (b *httpBridge) Stop() {
	registerHTTPListstoreFanout(nil)
	if b == nil {
		return
	}
	b.mu.Lock()
	srv := b.srv
	ln := b.ln
	b.srv = nil
	b.ln = nil
	b.addr = ""
	b.mu.Unlock()
	if srv != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		_ = srv.Shutdown(ctx)
		cancel()
	}
	if ln != nil {
		_ = ln.Close()
	}
	log.Print("【MCP】HTTP 桥已停止")
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, Mcp-Session-Id, Mcp-Protocol-Version, Last-Event-ID")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (b *httpBridge) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"ok":true}`))
}

func (b *httpBridge) handleSupportedOps(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(mcpcatalog.SupportedOpsJSON()))
}

func (b *httpBridge) handleInvoke(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req mcpwire.InvokeLine
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(mcpwire.InvokeResponse{ID: 0, OK: false, Error: "invalid json: " + err.Error()})
		return
	}
	if req.ID == 0 {
		req.ID = 1
	}
	res, err := b.host.Invoke(req.Op, req.Args)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mcpwire.InvokeResponse{ID: req.ID, OK: false, Error: err.Error()})
		return
	}
	var raw json.RawMessage
	switch v := res.(type) {
	case json.RawMessage:
		raw = v
	default:
		raw, _ = json.Marshal(v)
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(mcpwire.InvokeResponse{ID: req.ID, OK: true, Result: raw})
}

func (b *httpBridge) handleSSE(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, ": connected\n\n")
	flusher.Flush()

	ch := b.hub.subscribe(32)
	defer b.hub.unsubscribe(ch)

	// 周期性注释行保活（部分代理会断空闲连接）
	tick := time.NewTicker(25 * time.Second)
	defer tick.Stop()

	for {
		select {
		case msg := <-ch:
			// data 单行：payload 已为紧凑 JSON，一般无换行
			_, _ = fmt.Fprintf(w, "event: liststore\n")
			_, _ = fmt.Fprintf(w, "data: %s\n\n", strings.ReplaceAll(msg, "\n", ""))
			flusher.Flush()
		case <-tick.C:
			_, _ = fmt.Fprintf(w, ": ping %d\n\n", time.Now().Unix())
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}
