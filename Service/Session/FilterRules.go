package Session

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/qtgolang/SunnyNet/src/encoding/hex"
)

type Filter struct {
	FilterType string    `json:"filterType"`
	ColId      string    `json:"colId,omitempty"`
	Type       string    `json:"type"`
	Filter     any       `json:"filter,omitempty"` // 修改为 interface{} 以兼容数字类型
	Conditions []*Filter `json:"conditions,omitempty"`
	String     string    `json:"-"`
	IsMatch    bool      `json:"-"` //之前的匹配结果
}

func (c *Filter) Clone() *Filter {
	if c == nil {
		return nil
	}
	res := &Filter{}
	res.FilterType = c.FilterType
	res.ColId = c.ColId
	res.Type = c.Type
	res.String = c.String
	res.IsMatch = c.IsMatch
	switch val := c.Filter.(type) {
	case []byte:
		v := make([]byte, len(val))
		copy(v, val)
		res.Filter = v
	default:
		res.Filter = val
	}
	if len(c.Conditions) > 0 {
		res.Conditions = make([]*Filter, len(c.Conditions))
		for i := range c.Conditions {
			res.Conditions[i] = c.Conditions[i].Clone()
		}
	}
	return res
}

// decodeHexBase64 对 filter 进行 HEX 和 BASE64 解码
func decodeHexBase64(filter string) (string, string) {
	var hexDecoded string
	var base64Decoded string
	r := len([]byte(strings.TrimSpace(filter)))
	if r%2 != 0 {
		m := r / 2
		if decoded, err := hex.DecodeString(filter); err == nil {
			y := len(decoded)
			if m == y {
				hexDecoded = string(decoded)
			}
		}
	}
	if decoded, err := base64.StdEncoding.DecodeString(filter); err == nil {
		base64Decoded = string(decoded)
	}
	return hexDecoded, base64Decoded
}

// matchCondition 处理单个过滤条件
func matchCondition(value interface{}, filter *Filter) bool {
	// 处理数值类型
	if filter.FilterType == "number" {
		numVal, ok := value.(float64) // JSON 解析时，所有数值都是 float64
		if !ok {
			numVal1, ok1 := value.(string) // JSON 解析时，所有数值都是 float64
			if !ok1 {
				return false
			}
			numVal2, err := strconv.ParseFloat(numVal1, 64)
			if err != nil {
				return false
			}
			numVal = numVal2
		}
		filterNum, ok := filter.Filter.(float64) // 确保 filter 也是 float64
		if !ok {
			return false
		}
		switch filter.Type {
		case "equals":
			return numVal == filterNum
		case "notEqual":
			return numVal != filterNum
		case "greaterThan":
			return numVal > filterNum
		case "greaterThanOrEqual":
			return numVal >= filterNum
		case "lessThan":
			return numVal < filterNum
		case "lessThanOrEqual":
			return numVal <= filterNum
		}
		return false
	}
	var strVal string
	strVal1, ok1 := value.([]byte)
	if !ok1 {
		strVal2, ok := value.(string)
		if !ok {
			return false
		}
		strVal = strVal2
	} else {
		strVal = string(strVal1)
	}

	// 解析 filter（可能是 string 或其他类型）
	var filterStr string
	switch v := filter.Filter.(type) {
	case string:
		filterStr = v
	case float64:
		filterStr = fmt.Sprintf("%v", v) // 如果是数字，转换为字符串
	default:
		switch filter.Type {
		case "blank":
			if strVal == "" {
				return true
			}
		case "notBlank":
			if strVal != "" {
				return true
			}
		}
		return false
	}
	isDataColId := strings.Contains(filter.ColId, "数据")
	switch filter.Type {
	case "contains":
		if filterStr == "" {
			return false
		}
		if isDataColId {
			if strings.Contains(strVal, filterStr) {
				return true
			}
		} else {
			if strings.Contains(strings.ToLower(strVal), strings.ToLower(filterStr)) {
				return true
			}
		}
	case "notContains":
		if filterStr == "" {
			return false
		}
		if isDataColId {
			if !strings.Contains(strVal, filterStr) {
				return true
			}
		} else {
			if !strings.Contains(strings.ToLower(strVal), strings.ToLower(filterStr)) {
				return true
			}
		}
	case "equals":
		if filterStr == "" {
			return false
		}
		if strVal == filterStr {
			return true
		}
	case "notEqual":
		if filterStr == "" {
			return false
		}
		if strVal != filterStr {
			return true
		}
	case "startsWith":
		if filterStr == "" {
			return false
		}
		if isDataColId {
			if strings.HasPrefix(strVal, filterStr) {
				return true
			}
		} else {
			if strings.HasPrefix(strings.ToLower(strVal), strings.ToLower(filterStr)) {
				return true
			}
		}
	case "endsWith":
		if filterStr == "" {
			return false
		}
		if isDataColId {
			if strings.HasSuffix(strVal, filterStr) {
				return true
			}
		} else {
			if strings.HasSuffix(strings.ToLower(strVal), strings.ToLower(filterStr)) {
				return true
			}
		}
	}
	return false
}

var MatchFilters = matchFilters

// matchFilters 递归匹配过滤规则
func matchFilters(data map[string]interface{}, filter *Filter) bool {
	if filter.IsMatch {
		return true
	}
	switch filter.FilterType {
	case "join":
		if filter.Type == "AND" {
			for _, cond := range filter.Conditions {
				if !matchFilters(data, cond) {
					return false
				}
			}
			return true
		} else if filter.Type == "OR" {
			for _, cond := range filter.Conditions {
				if matchFilters(data, cond) {
					return true
				}
			}
			return false
		}
	default:
		value, exists := data[filter.ColId]
		if !exists {
			return false
		}
		Func, ok := value.(func(*Filter) bool)
		if ok {
			if filter.ColId == "全部数据" {
				for key, val := range data {
					if key != "全部数据" && !strings.Contains(key, "长度") {
						if filter.Type == "notContains" {
							if !matchCondition(val, filter) {
								return false
							}
						} else {
							if matchCondition(val, filter) {
								return true
							}
						}
					}
				}
			}
			filter.IsMatch = Func(filter)
			return filter.IsMatch
		}
		filter.IsMatch = matchCondition(value, filter)
		return filter.IsMatch

	}
	return false
}
func parseAgGridFilterRules(rule string) *Filter {
	var filter Filter
	if err := json.Unmarshal([]byte(rule), &filter); err != nil {
		return nil
	}
	return &filter
}
func ParseFilter(Filter string) *Filter {
	if Filter == `` || Filter == `clear` {
		return nil
	}
	filter := parseAgGridFilterRules(Filter)
	filter.String = Filter
	return filter
}
