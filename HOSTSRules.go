package main

import (
	"github.com/qtgolang/SunnyNet/src/protobuf/JSON"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type HostsRules struct {
	regex  *regexp.Regexp //源内容
	target string         //目标内容
}

var _HostsRules []HostsRules

func HostsRulesEvent(command string, args *JSON.SyJson) any {
	switch command {
	case "保存HOSTS规则":
		_TmpLock.Lock()
		defer _TmpLock.Unlock()
		var failHash []string
		var _Rules []HostsRules
		var _CRules []ConfigReplaceRules
		for i := 0; i < args.GetNum("Data"); i++ {
			_Hash := args.GetData("Data[" + strconv.Itoa(i) + "].Hash")
			_source := args.GetData("Data[" + strconv.Itoa(i) + "].源地址")
			_target := args.GetData("Data[" + strconv.Itoa(i) + "].新地址")
			_source = strings.ReplaceAll(_source, "\\\\", "\\")
			_target = strings.ReplaceAll(_target, "\\\\", "\\")
			_source = strings.ReplaceAll(_source, ".*", "{点星}")
			_source = strings.ReplaceAll(_source, "*", ".*")
			_source = strings.ReplaceAll(_source, "{点星}", ".*")
			regex, e := regexp.Compile(_source)
			if _source == "" || _target == "" || e != nil {
				failHash = append(failHash, _Hash)
				continue
			}
			_Rules = append(_Rules, HostsRules{regex: regex, target: _target})
			_CRules = append(_CRules, ConfigReplaceRules{Hash: _Hash, Src: _source, Dest: _target})
		}
		_HostsRules = _Rules
		GlobalConfig.HostsRules = _CRules
		_ = GlobalConfig.saveToFile()
		return failHash
	}
	return false
}

func HostsRulesUrl(u *url.URL) {
	if u == nil {
		return
	}
	_TmpLock.Lock()
	defer _TmpLock.Unlock()
	for i := 0; i < len(_HostsRules); i++ {
		regexReplace := false
		rx := _HostsRules[i].regex.ReplaceAllStringFunc(u.Host, func(match string) string {
			regexReplace = true
			return _HostsRules[i].target
		})
		if regexReplace {
			u.Host = rx
		}
	}
}
func HostsRulesAddress(u string) string {
	if u == "" {
		return ""
	}
	um := u
	_TmpLock.Lock()
	defer _TmpLock.Unlock()
	for i := 0; i < len(_HostsRules); i++ {
		regexReplace := false
		rx := _HostsRules[i].regex.ReplaceAllStringFunc(um, func(match string) string {
			regexReplace = true
			return _HostsRules[i].target
		})
		if regexReplace {
			um = rx
		}
	}
	return um
}
