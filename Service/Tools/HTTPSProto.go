package Tools

import (
	"changeme/Service/Config"
	"encoding/json"
	"errors"
	"strings"

	"github.com/qtgolang/SunnyNet/SunnyNet"
	"sync"
)

type HTTPSProto struct {
	App *SunnyNet.Sunny
}

var httpSProtoLock sync.Mutex

func (g *HTTPSProto) GetHTTPSProto() string {
	httpSProtoLock.Lock()
	defer httpSProtoLock.Unlock()
	return Config.Config.HTTPSProto
}
func (g *HTTPSProto) SetHTTPSProto(proto string) {
	httpSProtoLock.Lock()
	defer httpSProtoLock.Unlock()
	Config.Config.HTTPSProto = proto
	Config.Config.Save()
}

// ApplyHTTPSProtocol 保存并生效 HTTPS 协议/HTTP2 指纹/JA3（指纹在后续 HTTP 请求 ReplaceHttp 中按连接应用）。
func (g *HTTPSProto) ApplyHTTPSProtocol(sendIsHTTP1 *bool, protoJSON string, randomJa3 *bool) (map[string]any, error) {
	httpSProtoLock.Lock()
	defer httpSProtoLock.Unlock()

	if sendIsHTTP1 != nil {
		Config.Config.SendIsHTTP1 = *sendIsHTTP1
	}
	protoJSON = strings.TrimSpace(protoJSON)
	if protoJSON != "" {
		if !json.Valid([]byte(protoJSON)) {
			return nil, errors.New("http2Fingerprint 不是合法 JSON")
		}
		Config.Config.HTTPSProto = protoJSON
	}
	if randomJa3 != nil {
		Config.Config.RandomJa3 = *randomJa3
		g.App.SetRandomTLS(*randomJa3)
	}
	Config.Config.Save()

	protocol := "h2"
	label := "HTTP/2.0 优先"
	if Config.Config.SendIsHTTP1 {
		protocol = "http/1.1"
		label = "仅使用 HTTP/1.1 发送"
	}
	return map[string]any{
		"ok":               true,
		"applied":          true,
		"sendIsHTTP1":      Config.Config.SendIsHTTP1,
		"protocol":         protocol,
		"label":            label,
		"randomJa3":        Config.Config.RandomJa3,
		"http2Fingerprint": Config.Config.HTTPSProto,
		"http2WillApply":   !Config.Config.SendIsHTTP1 && Config.Config.HTTPSProto != "",
		"ja3Applied":       true,
	}, nil
}
func (g *HTTPSProto) SetRandomJa3(open bool) {
	httpSProtoLock.Lock()
	defer httpSProtoLock.Unlock()
	Config.Config.RandomJa3 = open
	g.App.SetRandomTLS(Config.Config.RandomJa3)
	Config.Config.Save()
}
func (g *HTTPSProto) GetRandomJa3() bool {
	httpSProtoLock.Lock()
	defer httpSProtoLock.Unlock()
	return Config.Config.RandomJa3
}
