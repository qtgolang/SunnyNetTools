package Config

import (
	"encoding/base64"
	"errors"
	"github.com/qtgolang/SunnyNet/src/encoding/hex"
	"github.com/qtgolang/SunnyNet/src/http"
)

type DoHTTPRequestInfo struct {
	Url      string      `json:"url"`
	Method   string      `json:"method"`
	Header   http.Header `json:"header"`
	Body     string      `json:"body"`
	BodyType string      `json:"bodyType"`
	OutTime  int         `json:"outTime"`
	Redirect bool        `json:"redirect"`
	Disguise bool        `json:"disguise"`
	ProxyIP  string      `json:"proxyIP"`
	Language string      `json:"Language"` //httpDebug调试工具，生成代码时使用
	Type     string      `json:"Type"`     //httpDebug调试工具，生成代码时使用
}
type DoHTTPResponseInfo struct {
	Body       string      `json:"body"`
	Error      string      `json:"Error"`
	Header     http.Header `json:"header"`
	Code       int         `json:"code"`
	ErrorLevel int         `json:"errorLevel"`
	Proto      string      `json:"proto"`
	Status     string      `json:"status"`
}

func (e *DoHTTPRequestInfo) GetBody() ([]byte, error, int) {
	bs, _ := base64.StdEncoding.DecodeString(e.Body)
	if e.BodyType == "text" {
		return bs, nil, 0
	}
	if e.BodyType == "hex" {
		a, b := hex.DecodeString(string(bs))
		return a, b, 1
	}
	if e.BodyType == "base64" {
		a, b := base64.StdEncoding.DecodeString(string(bs))
		return a, b, 2
	}
	return nil, errors.New("未知的内部错误0x001"), 3
}
