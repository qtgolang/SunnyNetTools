package Config

import (
	"encoding/base64"
	"github.com/qtgolang/SunnyNet/src/encoding/hex"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"strings"
)

type ReplaceBodyInfo struct {
	Type     string
	Lod      string
	New      string
	Source   string
	Note     string
	ID       int
	State    string
	BytesLod []byte `json:"-"`
	BytesNew []byte `json:"-"`
	Ok       bool
}

func (v *ReplaceBodyInfo) IsBreak() bool {
	return v.Type == "拦截请求"
}
func (v *ReplaceBodyInfo) SourceType_IsHTTP_URL() bool {
	return v.Source == "任意" || v.Source == "HTTP请求/响应" || v.Source == "HTTP请求-全部" || v.Source == "HTTP请求-URL"
}
func (v *ReplaceBodyInfo) SourceType_IsHTTP_RequestHeader() bool {
	return v.Source == "任意" || v.Source == "HTTP请求/响应" || v.Source == "HTTP请求-全部" || v.Source == "HTTP请求-协议头"
}
func (v *ReplaceBodyInfo) SourceType_IsHTTP_RequestBody() bool {
	return v.Source == "任意" || v.Source == "HTTP请求/响应" || v.Source == "HTTP请求-全部" || v.Source == "HTTP请求-提交数据"
}
func (v *ReplaceBodyInfo) SourceType_IsHTTP_ResponseHeader() bool {
	return v.Source == "任意" || v.Source == "HTTP请求/响应" || v.Source == "HTTP响应-全部" || v.Source == "HTTP响应-协议头"
}
func (v *ReplaceBodyInfo) SourceType_IsHTTP_ResponseBody() bool {
	return v.Source == "任意" || v.Source == "HTTP请求/响应" || v.Source == "HTTP响应-全部" || v.Source == "HTTP响应-响应数据"
}
func (v *ReplaceBodyInfo) SourceType_IsTCP_Send() bool {
	return v.Source == "任意" || v.Source == "Socket-任意" || v.Source == "TCP-全部" || v.Source == "TCP-发送"
}
func (v *ReplaceBodyInfo) SourceType_IsTCP_Receive() bool {
	return v.Source == "任意" || v.Source == "Socket-任意" || v.Source == "TCP-全部" || v.Source == "TCP-接收"
}
func (v *ReplaceBodyInfo) SourceType_IsUDP_Send() bool {
	return v.Source == "任意" || v.Source == "Socket-任意" || v.Source == "UDP-全部" || v.Source == "UDP-发送"
}
func (v *ReplaceBodyInfo) SourceType_IsUDP_Receive() bool {
	return v.Source == "任意" || v.Source == "Socket-任意" || v.Source == "UDP-全部" || v.Source == "UDP-接收"
}
func (v *ReplaceBodyInfo) SourceType_IsWebsocket_Send() bool {
	return v.Source == "任意" || v.Source == "Socket-任意" || v.Source == "UDP-全部" || v.Source == "Websocket-发送"
}
func (v *ReplaceBodyInfo) SourceType_IsWebsocket_Receive() bool {
	return v.Source == "任意" || v.Source == "Socket-任意" || v.Source == "Websocket-全部" || v.Source == "Websocket-接收"
}
func (v *ReplaceBodyInfo) Parse() bool {
	switch v.Type {
	case "拦截请求":
		v.BytesLod = []byte(v.Lod)
		return len(v.BytesLod) > 0
	case "字符串(UTF8)":
		v.BytesLod = []byte(v.Lod)
		v.BytesNew = []byte(v.New)
		return true
	case "字符串(GBK)":
		reader := transform.NewReader(strings.NewReader(v.Lod), simplifiedchinese.GBK.NewEncoder())
		__bs1, e1 := ioutil.ReadAll(reader)
		reader = transform.NewReader(strings.NewReader(v.New), simplifiedchinese.GBK.NewEncoder())
		__bs2, e2 := ioutil.ReadAll(reader)
		if e1 != nil || e2 != nil {
			return false
		}
		v.BytesLod = __bs1
		v.BytesNew = __bs2
		return true
	case "Base64":
		bs1, e1 := base64.StdEncoding.DecodeString(strings.ReplaceAll(strings.TrimSpace(v.Lod), " ", ""))
		bs2, e2 := base64.StdEncoding.DecodeString(strings.ReplaceAll(strings.TrimSpace(v.New), " ", ""))
		if e1 != nil || e2 != nil {
			return false
		}
		v.BytesLod = bs1
		v.BytesNew = bs2
		return true
	case "十六进制":
		bs1, e1 := hex.DecodeString(strings.ReplaceAll(strings.TrimSpace(v.Lod), " ", ""))
		bs2, e2 := hex.DecodeString(strings.ReplaceAll(strings.TrimSpace(v.New), " ", ""))
		if e1 != nil || e2 != nil {
			return false
		}
		v.BytesLod = bs1
		v.BytesNew = bs2
		return true
	}
	return false
}
func (f *config) initReplaceBody() {

	//恢复替换列表
	{
		m := f.ReplaceRoles
		f.ReplaceRoles = make(map[int]*ReplaceBodyInfo)
		for _, v := range m {
			if v == nil {
				continue
			}
			CertID++
			id := CertID
			if v.ID > 0 {
				id = v.ID
				if id > CertID {
					CertID = id
				}
			}
			f.ReplaceRoles[id] = &ReplaceBodyInfo{
				Type:   v.Type,
				Lod:    v.Lod,
				New:    v.New,
				Source: v.Source,
				Note:   v.Note,
				ID:     id,
			}
			f.ReplaceRoles[id].Ok = f.ReplaceRoles[id].Parse()
		}
	}
}
