package mcpbridge

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"

	"changeme/Service/mcpwire"
)

const maxMCPWireLine = 32 << 20 // 32 MiB

// TCPServer 在 127.0.0.1 上监听，供 stdio 子进程通过 TCP 调用 Host（无令牌鉴权；首行即为 InvokeLine）。
type TCPServer struct {
	mu       sync.Mutex
	ln       net.Listener
	token    string
	host     *Host
	closeCh  chan struct{}
	closedWG sync.WaitGroup
}

// NewTCPServer 创建未监听的实例；调用 Start 后生效。
func NewTCPServer(host *Host, token string) *TCPServer {
	return &TCPServer{host: host, token: strings.TrimSpace(token)}
}

// Addr 返回监听地址（未启动时为空串）。
func (s *TCPServer) Addr() string {
	if s == nil {
		return ""
	}
	s.mu.Lock()
	ln := s.ln
	s.mu.Unlock()
	if ln == nil {
		return ""
	}
	return ln.Addr().String()
}

// Start 在 host:port 上监听（host 为空则用 127.0.0.1）；仅本机回环。
func (s *TCPServer) Start(host string, port int) (string, error) {
	if s == nil || s.host == nil {
		return "", errors.New("mcp tcp: server or host nil")
	}
	if port <= 0 || port > 65535 {
		return "", errors.New("mcp tcp: invalid port")
	}
	if host == "" {
		host = "127.0.0.1"
	}
	addr := fmt.Sprintf("%s:%d", host, port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return "", err
	}
	s.mu.Lock()
	if s.ln != nil {
		s.mu.Unlock()
		_ = ln.Close()
		return "", errors.New("mcp tcp: already started")
	}
	s.ln = ln
	s.closeCh = make(chan struct{})
	s.mu.Unlock()

	s.closedWG.Add(1)
	go s.acceptLoop(ln)
	log.Printf("【MCP】TCP 桥已监听 %s", ln.Addr().String())
	return ln.Addr().String(), nil
}

func (s *TCPServer) acceptLoop(ln net.Listener) {
	defer s.closedWG.Done()
	for {
		conn, err := ln.Accept()
		if err != nil {
			select {
			case <-s.closeCh:
			default:
				if !strings.Contains(err.Error(), "use of closed") {
					log.Printf("【MCP】accept: %v", err)
				}
			}
			return
		}
		go s.handleConn(conn)
	}
}

func (s *TCPServer) handleConn(conn net.Conn) {
	defer conn.Close()
	sc := bufio.NewScanner(conn)
	buf := make([]byte, 0, 64*1024)
	sc.Buffer(buf, maxMCPWireLine)
	enc := json.NewEncoder(conn)

	tokenSet := strings.TrimSpace(s.token) != ""
	authDone := !tokenSet

	for sc.Scan() {
		line := sc.Bytes()
		if !authDone {
			var auth mcpwire.AuthLine
			if err := json.Unmarshal(line, &auth); err != nil {
				return
			}
			if !subtleConstantTimeEq(auth.Auth, s.token) {
				log.Print("【MCP】鉴权失败，已断开")
				return
			}
			authDone = true
			continue
		}

		var req mcpwire.InvokeLine
		if err := json.Unmarshal(line, &req); err != nil {
			_ = enc.Encode(mcpwire.InvokeResponse{ID: 0, OK: false, Error: "invalid json: " + err.Error()})
			continue
		}
		if req.ID == 0 {
			req.ID = 1
		}
		res, err := s.host.Invoke(req.Op, req.Args)
		if err != nil {
			_ = enc.Encode(mcpwire.InvokeResponse{ID: req.ID, OK: false, Error: err.Error()})
			continue
		}
		var raw json.RawMessage
		switch v := res.(type) {
		case json.RawMessage:
			raw = v
		default:
			raw, _ = json.Marshal(v)
		}
		_ = enc.Encode(mcpwire.InvokeResponse{ID: req.ID, OK: true, Result: raw})
	}
}

func subtleConstantTimeEq(a, b string) bool {
	a = strings.TrimSpace(a)
	b = strings.TrimSpace(b)
	if len(a) != len(b) {
		return false
	}
	var v byte
	for i := 0; i < len(a); i++ {
		v |= a[i] ^ b[i]
	}
	return v == 0
}

// Stop 关闭监听；可重复调用。
func (s *TCPServer) Stop() {
	s.mu.Lock()
	ln := s.ln
	ch := s.closeCh
	s.ln = nil
	s.closeCh = nil
	s.mu.Unlock()
	if ch != nil {
		close(ch)
	}
	if ln != nil {
		_ = ln.Close()
	}
	s.closedWG.Wait()
	log.Print("【MCP】TCP 桥已停止")
}
