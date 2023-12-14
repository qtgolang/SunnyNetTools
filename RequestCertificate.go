package main

import (
	"github.com/qtgolang/SunnyNet/Api"
	"github.com/qtgolang/SunnyNet/src/protobuf/JSON"
	"os"
	"strconv"
	"strings"
	"sync"
)

var _TmpLock sync.Mutex
var _RequestCertificateMap = make(map[int]string)

func RequestCertificate(command string, args *JSON.SyJson) any {
	switch command {
	case "创建证书管理器":
		id := Api.CreateCertificate()
		return id
	case "删除证书管理器":
		for i := 0; i < args.GetNum("context"); i++ {
			id := getInt(args.GetData("context[" + strconv.Itoa(i) + "].id"))
			_TmpLock.Lock()
			Api.DelHttpCertificate(_RequestCertificateMap[id])
			delete(GlobalConfig.RequestCertManager, id)
			_ = GlobalConfig.saveToFile()
			_TmpLock.Unlock()
			Api.RemoveCertificate(id)
		}
		return true
	case "查询证书CommonName":
		id := getInt(args.GetData("context"))
		host := Api.GetCommonName(id)
		if host != "" {
			{
				_TmpLock.Lock()
				if _RequestCertificateMap[id] != "" {
					Api.DelHttpCertificate(_RequestCertificateMap[id])
				}
				_TmpLock.Unlock()
			}

			{
				rule := args.GetData("rule")
				_rule := uint8(1)
				if rule == "解析及发送" {
					_rule = 2
				} else if rule == "仅解析" {
					_rule = 3
				}
				if host != "" {
					Api.AddHttpCertificate(host, id, _rule)
				}
			}

			{
				_TmpLock.Lock()
				if host != "" {
					_RequestCertificateMap[id] = host
				}
				_TmpLock.Unlock()
			}
		}
		return host
	case "载入请求证书":
		id := getInt(args.GetData("Data.context"))
		rule := args.GetData("Data.使用规则")
		file := strings.ReplaceAll(args.GetData("Data.证书文件"), "\\\\", "\\")
		file = strings.ToLower(file)
		if file == "null" {
			file = ""
		}
		pass := args.GetData("Data.密码")
		if pass == "null" {
			pass = ""
		}
		host := args.GetData("Data.主机名")
		if host == "null" {
			host = ""
		}
		if file == "" {
			CallJs("弹出错误信息", "载入证书失败:未选择证书文件")
			return false
		}
		if strings.HasSuffix(file, ".p12") || strings.HasSuffix(file, ".pfx") || strings.HasSuffix(file, ".pkcs12") {
			A, E := os.Open(file)
			if E != nil {
				CallJs("弹出错误信息", "载入证书失败:打开证书文件失败:"+E.Error())
				return false
			}
			_ = A.Close()
			if !Api.LoadP12Certificate(id, file, pass) {
				CallJs("弹出错误信息", "载入证书失败:请检查密码是否正确？")
				return false
			}
		} else if strings.HasSuffix(file, ".pem") || strings.HasSuffix(file, ".cer") {
			if !Api.AddCertPoolPath(id, file) {
				CallJs("弹出错误信息", "载入证书失败")
				return false
			}
		} else {
			CallJs("弹出错误信息", "载入证书失败:证书文件格式错误")
			return false
		}
		_rule := uint8(1)
		if rule == "解析及发送" {
			_rule = 2
		} else if rule == "仅解析" {
			_rule = 3
		}
		_TmpLock.Lock()
		GlobalConfig.RequestCertManager[id] = ConfigRequestCertManager{FilePath: file, PassWord: pass, Rule: _rule, Host: host}
		_ = GlobalConfig.saveToFile()
		_TmpLock.Unlock()
		{
			_TmpLock.Lock()
			if _RequestCertificateMap[id] != "" {
				Api.DelHttpCertificate(_RequestCertificateMap[id])
			}
			_TmpLock.Unlock()
		}

		if host != "" {
			Api.AddHttpCertificate(host, id, _rule)
		}
		_TmpLock.Lock()
		if host != "" {
			_RequestCertificateMap[id] = host
		}
		_TmpLock.Unlock()
		return true
	default:
		return processEvent(command, args)
	}
}
