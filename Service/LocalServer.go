package Service

import (
	"bytes"
	"changeme/Service/Config"
	"changeme/Service/Session"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"

	http2 "github.com/qtgolang/SunnyNet/src/http"
	"github.com/rs/cors"
)

var LocalServer = ""

func (g *AppMain) GetLocalServerPATH() string {
	return "http://" + LocalServer
}

// 启动一个 本地 服务器
func (g *AppMain) localServerInit() {
	// 创建 CORS 处理器
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           0, // 设置为 0 禁用缓存
	})

	// 将 CORS 处理器包装到 http.Handler 中
	handler := c.Handler(http.DefaultServeMux)
	http.HandleFunc("/UpdateHttpResponse", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		w.Header().Set("Server", "local")

		// 读取 POST 参数
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		params := r.URL.Query()
		Theology, _ := strconv.Atoi(params.Get("Theology"))
		var req HttpSessionResponse
		_ = json.Unmarshal(body, &req)
		var bs bytes.Buffer
		bs.Write(req.Body)
		if req.IsMaxLength {
			Body := g.GetHTTPResponseBody(Theology, true)
			if len(Body) > MaxBodyLength {
				bs.Write(Body[MaxBodyLength:])
			}
		}
		_res := &Session.HttpSessionResponse{Body: bs.Bytes(), Header: req.Header, Code: req.Code, State: req.State}
		g.updateHttpResponse(Theology, _res)
	})
	http.HandleFunc("/UpdateHttpRequest", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		w.Header().Set("Server", "local")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		params := r.URL.Query()
		Theology, _ := strconv.Atoi(params.Get("Theology"))
		var req HttpSessionRequest
		_ = json.Unmarshal(body, &req)
		var bs bytes.Buffer
		bs.Write(req.Body)
		if req.IsMaxLength {
			Body := g.GetHTTPRequestBody(Theology, true)
			if len(Body) > MaxBodyLength {
				bs.Write(Body[MaxBodyLength:])
			}
		}
		_res := &Session.HttpSessionRequest{Body: bs.Bytes(), Header: req.Header, Url: req.Url, Method: req.Method}
		g.updateHttpRequest(Theology, _res)
	})
	http.HandleFunc("/CopyRequest", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		w.Header().Set("Server", "local")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusOK)
			return
		}

		params := r.URL.Query()
		theology, _ := strconv.Atoi(params.Get("Theology"))
		name := params.Get("Name")
		rng := params.Get("Range")
		copyType := params.Get("Type")

		obj := Session.GetHttpSession(theology)
		if obj == nil {
			http.Error(w, "请求不存在了", http.StatusOK)
			return
		}

		var bs bytes.Buffer
		if len(body) < 5 {
			switch rng {
			case "all":
				if name == "RequestHex" {
					writeFullRequest(&bs, obj.Request)
				} else {
					writeFullRequest(&bs, obj.Response)
				}
			case "body":
				if name == "RequestHex" {
					bs.Write(obj.Request.Body)
				} else {
					bs.Write(obj.Response.Body)
				}
				if bs.Len() == 0 {
					http.Error(w, "没有Body内容", http.StatusOK)
					return
				}
			default:
				http.Error(w, "复制类型错误", http.StatusOK)
				return
			}
		} else {
			if rng != "all" && rng != "body" {
				http.Error(w, "复制类型错误", http.StatusOK)
				return
			}
			if name == "RequestHex" {
				var req HttpSessionRequest
				_ = json.Unmarshal(body, &req)
				if rng == "all" {
					writeFullRequest(&bs, req)
				} else {
					bs.Write(req.Body)
				}
			} else {
				var resp HttpSessionResponse
				_ = json.Unmarshal(body, &resp)
				if rng == "all" {
					writeFullRequest(&bs, resp)
				} else {
					bs.Write(resp.Body)
				}
			}
		}
		writeClipboardResponse(g, w, bs.Bytes(), copyType)
		http.Error(w, "", http.StatusOK)
	})
	http.HandleFunc("/DoHTTPRequest", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		w.Header().Set("Server", "local")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusOK)
			return
		}

		var req Config.DoHTTPRequestInfo
		_ = json.Unmarshal(body, &req)

		http.Error(w, g.DoHTTPRequest(req), http.StatusOK)
	})
	http.HandleFunc("/GenerateCode", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		w.Header().Set("Server", "local")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusOK)
			return
		}

		var req Config.DoHTTPRequestInfo
		_ = json.Unmarshal(body, &req)
		http.Error(w, g.DoHTTPGenerateCode(req), http.StatusOK)
	})

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	LocalServer = "127.0.0.1:" + fmt.Sprint(listener.Addr().(*net.TCPAddr).Port)
	go func() {
		_ = http.Serve(listener, handler)
		_ = listener.Close()
	}()
}
func writeFullRequest(buf *bytes.Buffer, msg Session.HttpMessage) {
	buf.WriteString(msg.FirstLine() + "\r\n")
	for k, vs := range msg.GetHeader() {
		for _, v := range vs {
			buf.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
		}
	}
	buf.Write([]byte("\r\n"))
	buf.Write(msg.GetBody())
	return

}

func writeClipboardResponse(g *AppMain, w http.ResponseWriter, data []byte, typ string) {
	var err string
	switch typ {
	case "Ansi":
		err = g.ClipboardWriteAll(string(data))
	case "Base64":
		err = g.ClipboardWriteAll(base64.StdEncoding.EncodeToString(data))
	case "十六进制":
		err = g.ClipboardWriteAll(Session.GetHexAllSpaces(data))
	default:
		http.Error(w, "复制类型不正确", http.StatusBadRequest)
		return
	}

	if err != "" {
		_, _ = w.Write([]byte(err))
	} else {
		_, _ = w.Write([]byte(""))
	}
}

type HttpSessionRequest struct {
	Method      string
	Url         string
	Proto       string
	Header      http2.Header
	Body        []byte
	BodyLength  int
	IsMaxLength bool
}
type HttpSessionResponse struct {
	Code        string //响应状态码
	State       string //响应状态文本
	Type        string
	Proto       string
	Header      http2.Header
	Body        []byte
	BodyLength  int
	IsMaxLength bool
}

func (h HttpSessionRequest) FirstLine() string {
	return h.Method + " " + h.Url + " " + h.Proto
}

func (h HttpSessionRequest) GetHeader() map[string][]string {
	return h.Header
}

func (h HttpSessionRequest) GetBody() []byte {
	return h.Body
}
func (h HttpSessionResponse) FirstLine() string {
	return h.Proto + " " + h.Code + " " + h.State
}

func (h HttpSessionResponse) GetHeader() map[string][]string {
	return h.Header
}

func (h HttpSessionResponse) GetBody() []byte {
	return h.Body
}
