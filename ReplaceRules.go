package main

import (
	"encoding/base64"
	"github.com/qtgolang/SunnyNet/src/encoding/hex"
	"github.com/qtgolang/SunnyNet/src/protobuf/JSON"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const ReplaceRulesType_Bytes = uint8(1)
const ReplaceRulesType_File = uint8(2)

type ReplaceRules struct {
	Type   uint8
	source []byte //源内容
	target []byte //目标内容
}

var _ReplaceRules []ReplaceRules

func ReplaceRulesEvent(command string, args *JSON.SyJson) any {
	switch command {
	case "保存替换规则":
		_TmpLock.Lock()
		defer _TmpLock.Unlock()
		var failHash []string
		var _Rules []ReplaceRules
		var _CRules []ConfigReplaceRules
		for i := 0; i < args.GetNum("Data"); i++ {
			_Hash := args.GetData("Data[" + strconv.Itoa(i) + "].Hash")
			_Type := args.GetData("Data[" + strconv.Itoa(i) + "].替换类型")
			_source := args.GetData("Data[" + strconv.Itoa(i) + "].源内容")
			_target := args.GetData("Data[" + strconv.Itoa(i) + "].替换内容")
			_source = strings.ReplaceAll(_source, "\\\\", "\\")
			_source = strings.ReplaceAll(_source, "\\\"", "\"")
			_target = strings.ReplaceAll(_target, "\\\\", "\\")
			_target = strings.ReplaceAll(_target, "\\\"", "\"")
			if _source == "" {
				failHash = append(failHash, _Hash)
				continue
			}
			if _Type == "Base64" {
				bs1, e := base64.StdEncoding.DecodeString(_source)
				if e != nil {
					failHash = append(failHash, _Hash)
					continue
				}
				bs2, e := base64.StdEncoding.DecodeString(_target)
				if e != nil {
					failHash = append(failHash, _Hash)
					continue
				}
				_Rules = append(_Rules, ReplaceRules{Type: ReplaceRulesType_Bytes, source: bs1, target: bs2})
				_CRules = append(_CRules, ConfigReplaceRules{Type: _Type, Hash: _Hash, Src: _source, Dest: _target})
			} else if _Type == "HEX" {
				bs1, e := hex.DecodeString(_source)
				if e != nil {
					failHash = append(failHash, _Hash)
					continue
				}
				bs2, e := hex.DecodeString(_target)
				if e != nil {
					failHash = append(failHash, _Hash)
					continue
				}
				_Rules = append(_Rules, ReplaceRules{Type: ReplaceRulesType_Bytes, source: bs1, target: bs2})
				_CRules = append(_CRules, ConfigReplaceRules{Type: _Type, Hash: _Hash, Src: _source, Dest: _target})

			} else if _Type == "String(UTF8)" {
				_Rules = append(_Rules, ReplaceRules{Type: ReplaceRulesType_Bytes, source: []byte(_source), target: []byte(_target)})
				_CRules = append(_CRules, ConfigReplaceRules{Type: _Type, Hash: _Hash, Src: _source, Dest: _target})
			} else if _Type == "String(GBK)" {
				_Rules = append(_Rules, ReplaceRules{Type: ReplaceRulesType_Bytes, source: Utf8ToGBK([]byte(_source)), target: Utf8ToGBK([]byte(_target))})
				_CRules = append(_CRules, ConfigReplaceRules{Type: _Type, Hash: _Hash, Src: _source, Dest: _target})
			} else if _Type == "响应文件" {
				bs1, e := os.ReadFile(_target)
				if e != nil {
					failHash = append(failHash, _Hash)
					continue
				}
				_Rules = append(_Rules, ReplaceRules{Type: ReplaceRulesType_File, source: []byte(_source), target: bs1})
				_CRules = append(_CRules, ConfigReplaceRules{Type: _Type, Hash: _Hash, Src: _source, Dest: _target})
			} else {
				failHash = append(failHash, _Hash)
			}
		}
		GlobalConfig.ReplaceRules = _CRules
		_ = GlobalConfig.saveToFile()
		_ReplaceRules = _Rules
		return failHash
	default:
		return HostsRulesEvent(command, args)
	}
}
func ReplaceURL(u *url.URL) (*url.URL, []byte) {
	if u == nil {
		return u, nil
	}
	ur := u.String()
	_TmpLock.Lock()
	defer _TmpLock.Unlock()
	ok := false
	res := make([]byte, 0)
	for i := 0; i < len(_ReplaceRules); i++ {
		if _ReplaceRules[i].Type == ReplaceRulesType_Bytes {
			if strings.Contains(ur, string(_ReplaceRules[i].source)) {
				ur = strings.ReplaceAll(ur, string(_ReplaceRules[i].source), string(_ReplaceRules[i].target))
				ok = true
			}
		} else if _ReplaceRules[i].Type == ReplaceRulesType_File {
			if strings.Contains(ur, string(_ReplaceRules[i].source)) {
				res = _ReplaceRules[i].target
				ok = true
				break
			}
		}
	}
	if !ok {
		return u, nil
	}
	um, e := url.Parse(ur)
	if e != nil {
		return u, nil
	}
	return um, res
}
func ReplaceHeader(header http.Header) {
	if header == nil {
		return
	}
	for key, values := range header {
		for i, value := range values {
			header[key][i] = string(ReplaceBody([]byte(value)))
		}
	}
}

func ReplaceBody(b []byte) []byte {
	ur := string(b)
	_TmpLock.Lock()
	defer _TmpLock.Unlock()
	for i := 0; i < len(_ReplaceRules); i++ {
		if _ReplaceRules[i].Type == ReplaceRulesType_Bytes {
			if strings.Contains(ur, string(_ReplaceRules[i].source)) {
				ur = strings.ReplaceAll(ur, string(_ReplaceRules[i].source), string(_ReplaceRules[i].target))
			}
		}
	}
	return []byte(ur)
}
