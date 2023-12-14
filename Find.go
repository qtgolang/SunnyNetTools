package main

import (
	"changeme/MapHash"
	"encoding/base64"
	"encoding/binary"
	"github.com/qtgolang/SunnyNet/src/encoding/hex"
	"math"
	"strconv"
	"strings"
)

type FindValue struct {
	Options       string
	Value         string
	Type          string
	Range         string
	Color         string
	Bytes         []byte
	BytesReverse  []byte //倒叙后的Bytes 查找 Int Float 需要使用
	SearchResult  map[int]bool
	CaseSensitive bool //是否区分大小写
	PbSkip        int
}

func (c *FindValue) getValue() string {
	DelSpace := strings.Contains(c.Options, "删除空格后搜索")
	v := ""
	if DelSpace {
		v = strings.ReplaceAll(strings.ReplaceAll(c.Value, "\t", ""), " ", "")
	} else {
		v = c.Value
	}
	if c.CaseSensitive {
		return v
	}
	return strings.ToLower(v)
}
func (c *FindValue) Find() any {
	c.CaseSensitive = !strings.Contains(c.Options, "不区分大小写")
	defer func() {
		Insert.Lock()
		SearchPercentage = -1
		Insert.Unlock()
	}()
	c.SearchResult = make(map[int]bool)
	if c.Value == "" {
		CallJs("弹出错误提示", "查找失败：请输入要搜索的内容")
		return nil
	}
	switch c.Type {
	case "UTF8":
		return c.FindUTF8()
	case "GBK":
		return c.FindGBK()
	case "Hex":
		c.CaseSensitive = true
		return c.FindHex()
	case "pb":
		return c.FindProtoBuf()
	case "Base64":
		c.CaseSensitive = true
		return c.FindBase64()
	case "整数4":
		c.CaseSensitive = true
		return c.FindInt32()
	case "整数8":
		c.CaseSensitive = true
		return c.FindInt64()
	case "浮点数4":
		c.CaseSensitive = true
		return c.FindFloat32()
	case "浮点数8":
		c.CaseSensitive = true
		return c.FindFloat64()
	default:
		return nil
	}
}
func (c *FindValue) FindUTF8() any {
	c.Bytes = []byte(c.getValue())
	return c.FindStart()
}
func (c *FindValue) FindGBK() any {
	c.Bytes = Utf8ToGBK([]byte(c.getValue()))
	return c.FindStart()
}
func (c *FindValue) FindHex() any {
	bs, e := hex.DecodeString(c.getValue())
	if e != nil {
		CallJs("弹出错误提示", "查找失败：输入的 HEX 不正确,请检查！！")
		return nil
	}
	c.Bytes = bs
	return c.FindStart()
}
func (c *FindValue) FindBase64() any {
	c.CaseSensitive = false
	b, e := base64.StdEncoding.DecodeString(c.getValue())
	if e != nil {
		CallJs("弹出错误提示", "查找失败：输入的 Base64 不正确,请检查！！")
		return nil
	}
	c.Bytes = b
	return c.FindStart()
}
func (c *FindValue) FindProtoBuf() any {
	c.Bytes = []byte(c.getValue())
	return c.FindStart()
}
func (c *FindValue) FindInt32() any {
	num, e := strconv.Atoi(c.getValue())
	if e != nil {
		CallJs("弹出错误提示", "查找失败：输入的数值不正确,或类型选择错误")
		return nil
	}
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, uint32(num))
	c.Bytes = bs
	bs = make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(num))
	c.BytesReverse = bs
	return c.FindStart()
}
func (c *FindValue) FindInt64() any {
	num, e := strconv.Atoi(c.getValue())
	if e != nil {
		CallJs("弹出错误提示", "查找失败：输入的数值不正确,或查找类型选择错误")
		return nil
	}
	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, uint64(num))
	c.Bytes = bs
	bs = make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, uint64(num))
	c.BytesReverse = bs
	return c.FindStart()
}
func (c *FindValue) FindFloat32() any {
	// 将字符串转换为float32
	f32, err := strconv.ParseFloat(c.getValue(), 32)
	if err != nil {
		CallJs("弹出错误提示", "查找失败：输入的数值不正确,或查找类型选择错误")
		return nil
	}
	bytes32 := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes32, math.Float32bits(float32(f32)))
	c.Bytes = bytes32
	binary.LittleEndian.PutUint32(bytes32, math.Float32bits(float32(f32)))
	c.BytesReverse = bytes32
	return c.FindStart()
}
func (c *FindValue) FindFloat64() any {
	// 将字符串转换为 float64
	f64, err := strconv.ParseFloat(c.getValue(), 64)
	if err != nil {
		CallJs("弹出错误提示", "查找失败：输入的数值不正确,或查找类型选择错误")
		return nil
	}
	bytes32 := make([]byte, 4)
	binary.BigEndian.PutUint64(bytes32, math.Float64bits(f64))
	c.Bytes = bytes32
	binary.LittleEndian.PutUint64(bytes32, math.Float64bits(f64))
	c.BytesReverse = bytes32
	return c.FindStart()
}

func CancelSearch() []int {
	Insert.Lock()
	i := LastSearch
	for _, v := range LastSearch {
		h := HashMap.GetRequest(v)
		if h != nil {
			h.Color.Search = ""
		}
	}
	LastSearch = make([]int, 0)
	for n := 0; n < len(LastSearchSocket); n++ {
		LastSearchSocket[n].Color = ""
	}
	LastSearchSocket = make([]*MapHash.UpdateSocketList, 0)
	Insert.Unlock()
	return i
}

func (c *FindValue) FindStart() any {
	if strings.Contains(c.Options, "取消之前的颜色标记") {
		Insert.Lock()
		for i := 0; i < len(LastSearchSocket); i++ {
			LastSearchSocket[i].Color = ""
		}
		LastSearchSocket = make([]*MapHash.UpdateSocketList, 0)
		Insert.Unlock()
	}
	HashMap.Search(c.Search)
	Insert.Lock()
	var _SearchResult []int
	for o, v := range c.SearchResult {
		if v {
			_SearchResult = append(_SearchResult, o)
		}
	}
	f := &SearchResult{SearchResult: _SearchResult, Color: c.Color}
	if strings.Contains(c.Options, "取消之前的颜色标记") {
		for _, v := range LastSearch {
			h := HashMap.GetRequest(v)
			if h != nil {
				h.Color.Search = ""
			}
		}
		f.LastSearchResult = LastSearch
		LastSearch = _SearchResult
	} else {
		for i := 0; i < len(_SearchResult); i++ {
			LastSearch = append(LastSearch, _SearchResult[i])
		}
	}
	Insert.Unlock()
	return f
}

var LastSearch []int
var LastSearchSocket []*MapHash.UpdateSocketList

type SearchResult struct {
	LastSearchResult []int
	SearchResult     []int
	Color            string
}

func (c *FindValue) caseSensitiveSearch(Theology int, s string) bool {
	res := ""
	if !c.CaseSensitive {
		res = strings.ToLower(s)
	} else {
		res = s
	}
	if strings.Contains(res, string(c.Bytes)) {
		c.SearchResult[Theology] = true
		return true
	}
	if len(c.BytesReverse) > 0 {
		if strings.Contains(res, string(c.BytesReverse)) {
			c.SearchResult[Theology] = true
			return true
		}
	}
	return false
}

func (c *FindValue) Search(Theology, percentage int, request *MapHash.Request) {
	Insert.Lock()
	SearchPercentage = percentage
	Insert.Unlock()
	if request == nil {
		return
	}
	if !request.Display {
		return
	}
	if c.Type == "pb" {
		if c.Range == "全部" || c.Range == "HTTP请求" {
			//在请求Body中搜索
			{
				if c.caseSensitiveSearch(Theology, _PbToJson(request.Body, c.PbSkip)) {
					request.Color.Search = c.Color
					return
				}
			}
		}
		if c.Range == "全部" || c.Range == "HTTP响应" {
			//在请求响应Body中搜索
			{
				if c.caseSensitiveSearch(Theology, _PbToJson(request.Response.Body, c.PbSkip)) {
					request.Color.Search = c.Color
					return
				}
			}
		}
		if c.Range == "全部" || c.Range == "socketSend" || c.Range == "socketRec" || c.Range == "socketAll" {
			sv := "上行"
			if c.Range == "socketRec" {
				sv = "下行"
			}
			if c.Range == "socketAll" || c.Range == "全部" {
				sv = "上行下行"
			}
			if request.SocketData != nil {
				for i := 0; i < len(request.SocketData); i++ {
					v := request.SocketData[i]
					if v != nil {
						if strings.Contains(sv, v.Info.Ico) {
							if c.caseSensitiveSearch(Theology, _PbToJson(v.Body, c.PbSkip)) {
								v.Info.Color = c.Color
								request.Color.Search = c.Color
								Insert.Lock()
								LastSearchSocket = append(LastSearchSocket, v.Info)
								Insert.Unlock()
							}
						}
					}
				}
			}
		}

		return
	}

	if c.Range == "全部" || c.Range == "HTTP请求" {
		//在URL中搜索
		{
			if c.caseSensitiveSearch(Theology, request.URL) {
				request.Color.Search = c.Color
				return
			}
		}
		//在请求协议头中搜索
		{
			if request.Header != nil {
				_t := ""
				for k, v := range request.Header {
					if len(v) == 0 {
						_t += k + ": \r\n"
						continue
					}
					for i := 0; i < len(v); i++ {
						_t += k + ": " + v[i] + "\r\n"
					}
				}
				if c.caseSensitiveSearch(Theology, _t) {
					request.Color.Search = c.Color
					return
				}
			}
		}
		//在请求Body中搜索
		{
			if c.caseSensitiveSearch(Theology, string(request.Body)) {
				request.Color.Search = c.Color
				return
			}
		}
	}
	if c.Range == "全部" || c.Range == "HTTP响应" {
		//在请求响应协议头中搜索
		{
			if request.Response.Header != nil {
				_t := ""
				for k, v := range request.Response.Header {
					if len(v) == 0 {
						_t += k + ": \r\n"
						continue
					}
					for i := 0; i < len(v); i++ {
						_t += k + ": " + v[i] + "\r\n"
					}
				}
				if c.caseSensitiveSearch(Theology, _t) {
					request.Color.Search = c.Color
					return
				}
			}
		}
		//在请求响应Body中搜索
		{
			if c.caseSensitiveSearch(Theology, string(request.Response.Body)) {
				request.Color.Search = c.Color
				return
			}
		}
	}
	if c.Range == "全部" || c.Range == "socketSend" || c.Range == "socketRec" || c.Range == "socketAll" {
		sv := "上行"
		if c.Range == "socketRec" {
			sv = "下行"
		}
		if c.Range == "socketAll" || c.Range == "全部" {
			sv = "上行下行"
		}
		if request.SocketData != nil {
			for i := 0; i < len(request.SocketData); i++ {
				v := request.SocketData[i]
				if v != nil {
					if strings.Contains(sv, v.Info.Ico) {
						if c.caseSensitiveSearch(Theology, string(v.Body)) {
							v.Info.Color = c.Color
							request.Color.Search = c.Color
							Insert.Lock()
							LastSearchSocket = append(LastSearchSocket, v.Info)
							Insert.Unlock()
						}
					}
				}
			}
		}
	}
}
