package Service

import (
	"bytes"
	"changeme/Service/Config"
	"changeme/Service/Session"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"github.com/qtgolang/SunnyNet/src/encoding/hex"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"math"
	"strconv"
	"strings"
)

type FindInfo struct {
	Theology  int
	Options   string
	Value     string
	Type      string
	Range     string
	Color     string
	ProtoSkip string
	Action    string
	Progress  int `json:"-"` //搜索进度
}

func (info *FindInfo) GetSearchValue() [][]byte {
	var result [][]byte
	value := info.Value
	if info.IsDeleteSpacesSearch() {
		value = strings.ReplaceAll(strings.TrimSpace(value), " ", "")
	}
	auto := info.IsAutoSearchType()

	if auto {
		//UTF8
		result = append(result, []byte(value))
		//GBK
		{
			reader := transform.NewReader(strings.NewReader(value), simplifiedchinese.GBK.NewEncoder())
			__bs, e := io.ReadAll(reader)
			if e == nil {
				if !bytes.Equal(__bs, []byte(value)) {
					result = append(result, __bs)
				}
			}
		}
		//hex
		{
			if info.IsHexSearchType() {
				__bs, e := hex.DecodeString(value)
				if e == nil {
					result = append(result, __bs)
				}
			}
		}
		//Base64
		{
			if info.IsBase64SearchType() {
				__bs, e := base64.StdEncoding.DecodeString(value)
				if e == nil {
					result = append(result, __bs)
				}
			}
		}
		//int32
		{
			__bs, err := strconv.Atoi(value)
			if err == nil {
				a, b := int32ToBytes(int32(__bs))
				result = append(result, a)
				result = append(result, b)
			}
		}
		//int64
		{
			__bs, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				a, b := int64ToBytes(__bs)
				result = append(result, a)
				result = append(result, b)
			}
		}
		//float32
		{
			f, err := strconv.ParseFloat(value, 32)
			if err == nil {
				n := math.Float32bits(float32(f))
				a, b := uint32ToBytes(n)
				result = append(result, a)
				result = append(result, b)
			}
		}
		//float64
		{
			f, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return result
			}
			n := math.Float64bits(f)
			a, b := uint64ToBytes(n)
			result = append(result, a)
			result = append(result, b)
		}
		return result
	}
	if info.IsUTF8SearchType() {
		result = append(result, []byte(value))
		return result
	}
	if info.IsGBKSearchType() {
		reader := transform.NewReader(strings.NewReader(value), simplifiedchinese.GBK.NewEncoder())
		__bs, e := io.ReadAll(reader)
		if e != nil {
			return result
		}
		result = append(result, __bs)
		return result
	}
	if info.IsHexSearchType() {
		__bs, e := hex.DecodeString(value)
		if e != nil {
			return result
		}
		result = append(result, __bs)
		return result
	}
	if info.IsBase64SearchType() {
		__bs, e := base64.StdEncoding.DecodeString(value)
		if e != nil {
			return result
		}
		result = append(result, __bs)
		return result
	}
	if info.IsInt32SearchType() {
		__bs, err := strconv.Atoi(value)
		if err != nil {
			return result
		}
		a, b := int32ToBytes(int32(__bs))
		result = append(result, a)
		result = append(result, b)
		return result
	}
	if info.IsInt64SearchType() {
		__bs, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return result
		}
		a, b := int64ToBytes(__bs)
		result = append(result, a)
		result = append(result, b)
		return result
	}
	if info.IsFloat32SearchType() {
		f, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return result
		}
		n := math.Float32bits(float32(f))
		a, b := uint32ToBytes(n)
		result = append(result, a)
		result = append(result, b)
		return result
	}
	if info.IsInt64SearchType() {
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return result
		}
		n := math.Float64bits(f)
		a, b := uint64ToBytes(n)
		result = append(result, a)
		result = append(result, b)
		return result
	}
	if info.IsProtoBufSearchType() {
		result = append(result, []byte(value))
		reader := transform.NewReader(strings.NewReader(value), simplifiedchinese.GBK.NewEncoder())
		__bs, e := io.ReadAll(reader)
		if e == nil {
			if !bytes.Equal(__bs, []byte(value)) {
				result = append(result, __bs)
			}
		}
		return result
	}
	return result
}
func uint32ToBytes(n uint32) ([]byte, []byte) {
	bigEndian := new(bytes.Buffer)
	_ = binary.Write(bigEndian, binary.BigEndian, n)

	littleEndian := new(bytes.Buffer)
	_ = binary.Write(littleEndian, binary.LittleEndian, n)

	return bigEndian.Bytes(), littleEndian.Bytes()
}

// uint64 转 []byte
func uint64ToBytes(n uint64) ([]byte, []byte) {
	bigEndian := new(bytes.Buffer)
	_ = binary.Write(bigEndian, binary.BigEndian, n)

	littleEndian := new(bytes.Buffer)
	_ = binary.Write(littleEndian, binary.LittleEndian, n)

	return bigEndian.Bytes(), littleEndian.Bytes()
}

func int64ToBytes(n int64) ([]byte, []byte) {
	// 大端
	bigEndian := new(bytes.Buffer)
	_ = binary.Write(bigEndian, binary.BigEndian, n)

	// 小端
	littleEndian := new(bytes.Buffer)
	_ = binary.Write(littleEndian, binary.LittleEndian, n)

	return bigEndian.Bytes(), littleEndian.Bytes()
}
func int32ToBytes(n int32) ([]byte, []byte) {
	// 大端
	bigEndian := new(bytes.Buffer)
	_ = binary.Write(bigEndian, binary.BigEndian, n)

	// 小端
	littleEndian := new(bytes.Buffer)
	_ = binary.Write(littleEndian, binary.LittleEndian, n)
	return bigEndian.Bytes(), littleEndian.Bytes()
}

// IsRangeHTTPRequest 是否在HTTP请求查找
func (info *FindInfo) IsRangeHTTPRequest() bool {
	return info.Range == "HTTP请求" || info.Range == "在全部范围寻找"
}

// IsRangeHTTPResponse 是否在HTTP响应查找
func (info *FindInfo) IsRangeHTTPResponse() bool {
	return info.Range == "HTTP响应" || info.Range == "在全部范围寻找"
}

// IsRangeSocketSend 是否在Socket发送 中查找
func (info *FindInfo) IsRangeSocketSend() bool {
	return info.Range == "socketSend" || info.Range == "socketAll" || info.Range == "在全部范围寻找"
}

// IsRangeSocketReceive 是否在Socket接收 中查找
func (info *FindInfo) IsRangeSocketReceive() bool {
	return info.Range == "socketRec" || info.Range == "socketAll" || info.Range == "在全部范围寻找"
}

// GetProtoBufSkipSize 获取 ProtoBuf 跳过多少字节
func (info *FindInfo) GetProtoBufSkipSize() int {
	v, _ := strconv.Atoi(info.ProtoSkip)
	return v
}

// IsColorTag 是否为颜色标记模式
func (info *FindInfo) IsColorTag() bool {
	return !info.IsHideTag()
}

// IsHideTag 是否为隐藏不相关模式
func (info *FindInfo) IsHideTag() bool {
	return info.Action == "hide"
}

// IsGlobal 是否全局搜索
func (info *FindInfo) IsGlobal() bool {
	return info.Theology < 1
}

// IsCancelColor 是否取消颜色标记
func (info *FindInfo) IsCancelColor() bool {
	return strings.Contains(info.Options, "取消之前的颜色标记")
}

// IsNoUpperLowerCase 是否不区分大小写
func (info *FindInfo) IsNoUpperLowerCase() bool {
	return strings.Contains(info.Options, "不区分大小写")
}

// IsDeleteSpacesSearch 是否删除空格后搜索
func (info *FindInfo) IsDeleteSpacesSearch() bool {
	return strings.Contains(info.Options, "删除空格后搜索")
}

// IsAutoSearchType 是否模糊搜索类型
func (info *FindInfo) IsAutoSearchType() bool {
	return info.Type == "auto"
}

// IsUTF8SearchType 是否按照 UTF8 搜索
func (info *FindInfo) IsUTF8SearchType() bool {
	return info.Type == "UTF8"
}

// IsGBKSearchType 是否按照 GBK 搜索
func (info *FindInfo) IsGBKSearchType() bool {
	return info.Type == "GBK"
}

// IsHexSearchType 是否按照 Hex 搜索
func (info *FindInfo) IsHexSearchType() bool {
	return info.Type == "Hex"
}

// IsBase64SearchType 是否按照 Base64 搜索
func (info *FindInfo) IsBase64SearchType() bool {
	return info.Type == "Base64"
}

// IsInt32SearchType 是否按照 Int32 搜索
func (info *FindInfo) IsInt32SearchType() bool {
	return info.Type == "整数4"
}

// IsInt64SearchType 是否按照 Int64 搜索
func (info *FindInfo) IsInt64SearchType() bool {
	return info.Type == "整数8"
}

// IsFloat32SearchType 是否按照 Float32 搜索
func (info *FindInfo) IsFloat32SearchType() bool {
	return info.Type == "浮点数4"
}

// IsFloat64SearchType 是否按照 Float64 搜索
func (info *FindInfo) IsFloat64SearchType() bool {
	return info.Type == "浮点数8"
}

// IsProtoBufSearchType 是否按照 ProtoBuf 搜索
func (info *FindInfo) IsProtoBufSearchType() bool {
	return info.Type == "pb"
}
func (g *AppMain) FindSession(info *FindInfo) []int {
	Config.AppList["Main"].EmitEvent("FindSearchProgress", 1)
	res := g.find(info)
	Config.AppList["Main"].EmitEvent("FindSearchProgress", 100)
	return res
}
func (g *AppMain) find(i *FindInfo) []int {
	var res []int
	if i == nil {
		return res
	}
	lock.Lock()
	defer lock.Unlock()
	searchArray := i.GetSearchValue()
	if len(searchArray) == 0 {
		return res
	}
	if !i.IsGlobal() {
		session := Session.GetAppSession(i.Theology)
		if session != nil {
			session.RangeStream(func(stream Session.AppStream) bool {
				if stream.GetIsSend() && i.IsRangeSocketSend() || !stream.GetIsSend() && i.IsRangeSocketReceive() {
					for _, v := range searchArray {
						if bytes.Contains(stream.GetBody(), v) {
							res = append(res, stream.GetMessageId())
						}
					}
				}
				return true
			})
		}
		return res
	}

	Session.Session.Range(func(key, value any) bool {
		theology, ok := key.(int)
		if !ok {
			return true
		}
		if httpObj, ok := value.(*Session.HttpSession); ok {
			if i.IsRangeHTTPRequest() && containsInHTTP(httpObj.Request.Method, httpObj.Request.Url, httpObj.Request.Proto, httpObj.Request.Header, httpObj.Request.Body, searchArray) {
				res = append(res, theology)
				return true
			}
			if i.IsRangeHTTPResponse() && containsInHTTP(httpObj.Response.Proto, httpObj.Response.Code, httpObj.Response.State, httpObj.Response.Header, httpObj.Response.Body, searchArray) {
				res = append(res, theology)
				return true
			}
		}

		if appSession, ok := value.(Session.AppSession); ok && (i.IsRangeSocketSend() || i.IsRangeSocketReceive()) {
			if searchInStream(appSession, searchArray, i) {
				res = append(res, theology)
				return true
			}
		}
		return true
	})

	return res
}

// 在 HTTP 请求或响应中查找目标数据
func containsInHTTP(proto, code, state string, header map[string][]string, body []byte, searchArray [][]byte) bool {
	var bs bytes.Buffer
	bs.WriteString(fmt.Sprintf("%s %s %s\r\n", proto, code, state))
	for nv, vc := range header {
		for _, vv := range vc {
			bs.WriteString(fmt.Sprintf("%s: %s\r\n", nv, vv))
		}
	}
	bs.WriteString("\r\n")
	bs.Write(body)
	for _, v := range searchArray {
		if bytes.Contains(bs.Bytes(), v) {
			return true
		}
	}
	return false
}

// 在 AppStream 中查找数据
func searchInStream(appSession Session.AppSession, searchArray [][]byte, i *FindInfo) bool {
	found := false
	appSession.RangeStream(func(stream Session.AppStream) bool {
		if stream.GetIsSend() && i.IsRangeSocketSend() || !stream.GetIsSend() && i.IsRangeSocketReceive() {
			for _, v := range searchArray {
				if bytes.Contains(stream.GetBody(), v) {
					found = true
					return false
				}
			}
		}
		return true
	})
	return found
}
