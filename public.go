package main

import (
	"bytes"
	"changeme/CommAnd"
	"changeme/MapHash"
	"changeme/Resource"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/qtgolang/SunnyNet/Api"
	"github.com/qtgolang/SunnyNet/SunnyNet"
	"github.com/qtgolang/SunnyNet/public"
	"github.com/qtgolang/SunnyNet/src/protobuf/JSON"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var startWork = true
var workLock sync.RWMutex

func SetWorkingState(open bool) {
	workLock.Lock()
	defer workLock.Unlock()
	startWork = open
}
func GetWorkingState() bool {
	workLock.RLock()
	defer workLock.RUnlock()
	return startWork
}

func event(command string, args *JSON.SyJson) any {
	switch command {
	case "设置断点模式":
		breakpoint = getInt(args.GetData("break"))
		return true
	case "清空":
		HashMap.Empty()
		return true
	case "关闭请求会话":
		var TheologyArray []int
		for i := 0; i < args.GetNum("Data"); i++ {
			TheologyArray = append(TheologyArray, getInt(args.GetData("Data["+strconv.Itoa(i)+"]")))
		}
		HashMap.CloseSession(TheologyArray)
		return true
	case "删除请求会话":
		var TheologyArray []int
		for i := 0; i < args.GetNum("Data"); i++ {
			TheologyArray = append(TheologyArray, getInt(args.GetData("Data["+strconv.Itoa(i)+"]")))
		}
		HashMap.Delete(TheologyArray)
		return true
	case "全部放行":
		HashMap.ReleaseAll()
		return true
	case "创建请求代码":
		var TheologyArray []int
		Lang := args.GetData("Lang")
		Module := args.GetData("Module")
		for i := 0; i < args.GetNum("Data"); i++ {
			TheologyArray = append(TheologyArray, getInt(args.GetData("Data["+strconv.Itoa(i)+"]")))
		}
		Code := CreateRequestCode(TheologyArray, Lang, Module)
		_ = runtime.ClipboardSetText(app.ctx, Code)
		CallJs("弹出成功信息", "请求代码已生成到剪辑版")
		return true
	case "保存文件":
		var TheologyArray []int
		Path := strings.ReplaceAll(args.GetData("Path"), "\\\\", "\\")
		if Path == "" {
			return false
		}
		if !strings.HasSuffix(strings.ToLower(Path), ".syn") {
			Path += ".syn"
		}
		ALL := args.GetData("ALL") == "true"
		for i := 0; i < args.GetNum("Data"); i++ {
			TheologyArray = append(TheologyArray, getInt(args.GetData("Data["+strconv.Itoa(i)+"]")))
		}
		return saveToFile(Path, ALL, TheologyArray)
	case "打开记录文件":
		var OpenData []*MapHash.Request
		{
			Path := strings.ReplaceAll(args.GetData("Path"), "\\\\", "\\")
			if Path == "" {
				return false
			}
			SetStatusText("正在读取文件:" + Path)
			bs, e := os.ReadFile(Path)
			if e != nil {
				SetStatusText("读取文件失败:" + e.Error())
				return false
			}
			DATA := MapHash.BrUnCompress(bs)
			e = json.Unmarshal(DATA, &OpenData)
			if e != nil {
				SetStatusText("解密文件失败:" + e.Error())
				return false
			}
		}
		max := strconv.Itoa(len(OpenData))
		xh := OpenData[len(OpenData)-1]
		fmt.Println(xh)
		SetStatusText("正在导入记录:0/" + max)
		var OpenFileListInfo []ListInfo
		for index, v := range OpenData {
			if v != nil {
				SetStatusText("正在导入记录:" + strconv.Itoa(index) + "/" + max)
				Theology := HashMap.CreateUniqueID()
				HashMap.SetRequest(Theology, v)
				State := strconv.Itoa(v.Response.StateCode)
				if !strings.Contains(strings.ToUpper(v.URL), "HTTP") {
					State = "已断开"
				}
				ResponseType := ""
				Method := v.Method
				Ico := "websocket_close"
				ResponseLen := ""
				if v.Way == "Websocket" {
					Method = "Websocket"
					ResponseType = "Websocket"
					ResponseLen = strconv.Itoa(v.SendNum) + "/" + strconv.Itoa(v.RecNum)
				} else if v.Way == "UDP" {
					Method = "UDP"
					ResponseType = "UDP"
					ResponseLen = strconv.Itoa(v.SendNum) + "/" + strconv.Itoa(v.RecNum)
				} else if strings.Contains(strings.ToUpper(Method), "TCP") {
					ResponseLen = strconv.Itoa(v.SendNum) + "/" + strconv.Itoa(v.RecNum)
					ResponseType = Method
				} else {
					if v.Response.Header != nil {
						_a := v.Response.Header["Content-Type"]
						if len(_a) > 0 {
							ResponseType = _a[0]
						} else {
							_a = v.Response.Header["content-type"]
							if len(_a) > 0 {
								ResponseType = _a[0]
							}
						}
						if ResponseType != "" {
							array := strings.Split(ResponseType+";", ";")
							if len(array) > 0 {
								ResponseType = array[0]
							}
						}
					}
					Ico = UpdateIco(v, ResponseType)
					ResponseLen = strconv.Itoa(len(v.Response.Body))
				}
				tmp := ListInfo{
					MessageId: -1,
					Theology:  Theology,
					State:     State,
					URL:       v.URL,
					ClientIP:  v.ClientIP,
					PID:       v.PID,
					Method:    Method,
					Ico:       Ico,
					Len:       ResponseLen,
					Type:      ResponseType,
					SendTime:  v.SendTime,
					RecTime:   v.RecTime,
					Notes:     v.Notes,
				}
				tmp.Color.TagColor = v.Color.TagColor
				tmp.Color.Search = v.Color.Search
				OpenFileListInfo = append(OpenFileListInfo, tmp)
			}
		}
		SetStatusText("导入完成: " + max + " 条记录")
		if len(OpenFileListInfo) > 0 {
			CallJs("插入列表", OpenFileListInfo)
		}
		return true
	case "更新注释":
		Theology := getInt(args.GetData("Theology"))
		Data := args.GetData("Data")
		h := HashMap.GetRequest(Theology)
		if h != nil {
			h.Notes = Data
			return true
		}
		return false
	case "更新主题":
		Dark := args.GetData("Dark") == "true"
		_TmpLock.Lock()
		if Dark {
			GlobalConfig.DarkTheme = 1
		} else {
			GlobalConfig.DarkTheme = 2
		}

		_ = GlobalConfig.saveToFile()
		_TmpLock.Unlock()
		return false
	case "标记颜色":
		empty := args.GetData("empty") == "true"
		if empty {
			for i := 0; i < args.GetNum("Data"); i++ {
				Theology := getInt(args.GetData("Data[" + strconv.Itoa(i) + "]"))
				h := HashMap.GetRequest(Theology)
				if h != nil {
					h.Color.TagColor = ""
				}
			}
			return true
		}
		for i := 0; i < args.GetNum("Data"); i++ {
			Theology := getInt(args.GetData("Data[" + strconv.Itoa(i) + "].id"))
			TagColor := args.GetData("Data[" + strconv.Itoa(i) + "].color")
			h := HashMap.GetRequest(Theology)
			if h != nil {
				h.Color.TagColor = TagColor
			}
		}
		return false
	case "重发请求":
		var TheologyArray []int
		mode := getInt(args.GetData("Mode"))
		for i := 0; i < args.GetNum("Data"); i++ {
			TheologyArray = append(TheologyArray, getInt(args.GetData("Data["+strconv.Itoa(i)+"]")))
		}
		if GlobalConfig.Authentication {
			CallJs("弹出错误信息", "请在设置中关闭身份验证模式后再试！")
			return false
		}
		HashMap.Resend(TheologyArray, mode, app.App.Port())
		CallJs("弹出成功信息", "重发请求已提交")
		return true
	case "工作状态":
		SetWorkingState(args.GetData("State") == "true")
		return true
	case "设置IE代理":
		ok := app.App.SetIeProxy(args.GetData("Set") != "true")
		if !ok {
			CallJs("弹出错误提示", "设置IE代理失败")
		}
		return ok
	case "选择文件":
		DefaultDirectory, _ := CommAnd.GetDesktopPath()
		_Title := args.GetData("Title")
		_Filters := make([]runtime.FileFilter, 0)
		for i := 0; i < args.GetNum("Filters"); i++ {
			_Filters = append(_Filters, runtime.FileFilter{
				DisplayName: args.GetData("Filters[" + strconv.Itoa(i) + "].Name"),
				Pattern:     args.GetData("Filters[" + strconv.Itoa(i) + "].Pattern"),
			})
		}
		res, _ := runtime.OpenFileDialog(app.ctx, runtime.OpenDialogOptions{
			DefaultDirectory:           DefaultDirectory,
			Title:                      _Title,
			Filters:                    _Filters,
			ShowHiddenFiles:            true,
			CanCreateDirectories:       false,
			ResolvesAliases:            true,
			TreatPackagesAsDirectories: false,
		})
		return res
	case "保存文件对话框":
		DefaultDirectory, _ := CommAnd.GetDesktopPath()
		_Title := args.GetData("Title")
		_Filters := make([]runtime.FileFilter, 0)
		for i := 0; i < args.GetNum("Filters"); i++ {
			_Filters = append(_Filters, runtime.FileFilter{
				DisplayName: args.GetData("Filters[" + strconv.Itoa(i) + "].Name"),
				Pattern:     args.GetData("Filters[" + strconv.Itoa(i) + "].Pattern"),
			})
		}
		res, _ := runtime.SaveFileDialog(app.ctx, runtime.SaveDialogOptions{
			DefaultDirectory:           DefaultDirectory,
			Title:                      _Title,
			Filters:                    _Filters,
			ShowHiddenFiles:            true,
			CanCreateDirectories:       true,
			TreatPackagesAsDirectories: false,
		})
		return res
	case "取消搜索颜色标记":
		return CancelSearch()
	case "查找":
		obj := &FindValue{
			Value:   args.GetData("Value"),
			Options: args.GetData("Options"),
			Type:    args.GetData("Type"),
			Range:   args.GetData("Range"),
			Color:   args.GetData("Color"),
			PbSkip:  getInt(args.GetData("ProtoSkip")),
		}
		return obj.Find()
	//主界面的关闭按钮点击
	case "CloseWindow":
		w, h := runtime.WindowGetSize(app.ctx)
		runtime.Hide(app.ctx)
		code, _ := base64.StdEncoding.DecodeString(strings.ReplaceAll(args.GetData("Filter"), "\\\\", "\\"))
		KeysStrings, _ := base64.StdEncoding.DecodeString(strings.ReplaceAll(args.GetData("KeysStrings"), "\\\\", "\\"))
		_TmpLock.Lock()
		GlobalConfig.Filter = string(code)
		GlobalConfig.Size.Width = w
		GlobalConfig.Size.Height = h
		GlobalConfig.KeysStrings = string(KeysStrings)
		_ = GlobalConfig.saveToFile()
		_TmpLock.Unlock()
		app.App.SetIeProxy(true)
		os.Exit(0)
		return nil
	case "获取内网IP":
		return CommAnd.GetWayArray()
	case "保存配置":
		return SaveData(args)
	case "获取脚本模板列表":
		return Resource.ScriptTemplate
	case "加载配置":
		CallJs("加载配置", GlobalConfig)
		return true
	//UI加载完成
	case "init":
		CallJs("加载配置", GlobalConfig)
		if app.App != nil {
			errStr := ""
			if app.App.Error != nil {
				errStr = app.App.Error.Error()
			}
			CallJs("启动状态", base64.StdEncoding.EncodeToString([]byte(errStr)))
			return nil
		}
		go lanZouUpdate()
		app.App = SunnyNet.NewSunny()
		app.App.SetPort(GlobalConfig.Port)
		//app.App.MustTcp(true)
		app.App.SetGoCallback(HttpCallback, TcpCallback, WSCallback, UdpCallback)
		err := app.App.Start().Error
		errStr := ""
		if err != nil {
			errStr = err.Error()
		}
		app.App.Socket5VerifyUser(GlobalConfig.Authentication)
		for k, v := range GlobalConfig.AuthenticationUserInfo {
			app.App.Socket5AddUser(k, v)
			_TmpLock.Lock()
			SocketAuthentication = append(SocketAuthentication, k)
			_TmpLock.Unlock()
		}
		//强制走TCP
		{
			app.App.MustTcp(GlobalConfig.MustTcp.Open)
			_ = app.App.SetMustTcpRegexp(GlobalConfig.MustTcp.Rules)
		}
		//证书选择
		{
			if !GlobalConfig.Cert.Default {
				id := Api.CreateCertificate()
				defer Api.RemoveCertificate(id)
				CaFilePath := GlobalConfig.Cert.CaPath
				KeyFilePath := GlobalConfig.Cert.KeyPath
				if !Api.LoadX509KeyPair(id, CaFilePath, KeyFilePath) {
					bs1, e := os.ReadFile(CaFilePath)
					if e != nil {
						GlobalConfig.Cert.Default = true
						_ = GlobalConfig.saveToFile()
					} else {
						bs2, e1 := os.ReadFile(KeyFilePath)
						if e1 != nil {
							GlobalConfig.Cert.Default = true
							_ = GlobalConfig.saveToFile()
						} else {
							if !Api.LoadX509Certificate(id, "", string(bs1), string(bs2)) {
								GlobalConfig.Cert.Default = true
								_ = GlobalConfig.saveToFile()
							} else {
								app.App.Error = nil
								e2 := app.App.SetCert(id).Error
								if e2 != nil {
									GlobalConfig.Cert.Default = true
									_ = GlobalConfig.saveToFile()
								}
							}
						}
					}
				} else {
					app.App.Error = nil
					e2 := app.App.SetCert(id).Error
					if e2 != nil {
						GlobalConfig.Cert.Default = true
						_ = GlobalConfig.saveToFile()
					}
				}
			}
		}
		//替换规则
		{
			var _Rules []ReplaceRules
			for i := 0; i < len(GlobalConfig.ReplaceRules); i++ {
				v := GlobalConfig.ReplaceRules[i]
				_Type := v.Type
				_source := v.Src
				_target := v.Dest
				if _source == "" {
					continue
				}
				if _Type == "Base64" {
					bs1, e := base64.StdEncoding.DecodeString(_source)
					if e != nil {
						continue
					}
					bs2, e := base64.StdEncoding.DecodeString(_target)
					if e != nil {
						continue
					}
					_Rules = append(_Rules, ReplaceRules{Type: ReplaceRulesType_Bytes, source: bs1, target: bs2})
				} else if _Type == "HEX" {
					bs1, e := hex.DecodeString(_source)
					if e != nil {
						continue
					}
					bs2, e := hex.DecodeString(_target)
					if e != nil {
						continue
					}
					_Rules = append(_Rules, ReplaceRules{Type: ReplaceRulesType_Bytes, source: bs1, target: bs2})
				} else if _Type == "String(UTF8)" {
					_Rules = append(_Rules, ReplaceRules{Type: ReplaceRulesType_Bytes, source: []byte(_source), target: []byte(_target)})
				} else if _Type == "String(GBK)" {
					_Rules = append(_Rules, ReplaceRules{Type: ReplaceRulesType_Bytes, source: Utf8ToGBK([]byte(_source)), target: Utf8ToGBK([]byte(_target))})
				} else if _Type == "响应文件" {
					bs1, e := os.ReadFile(_target)
					if e != nil {
						continue
					}
					_Rules = append(_Rules, ReplaceRules{Type: ReplaceRulesType_File, source: []byte(_source), target: bs1})
				}
			}
			_ReplaceRules = _Rules
		}
		//Hosts规则
		{
			var _Rules []HostsRules
			for i := 0; i < len(GlobalConfig.HostsRules); i++ {
				v := GlobalConfig.HostsRules[i]
				_source := v.Src
				_target := v.Dest
				regex, e := regexp.Compile(_source)
				if _source == "" || _target == "" || e != nil {
					continue
				}
				_Rules = append(_Rules, HostsRules{regex: regex, target: _target})
			}
			_HostsRules = _Rules
		}
		//请求证书规则初始化
		{
			var _Rules []ConfigRequestCertManager
			for _, v := range GlobalConfig.RequestCertManager {
				_Rules = append(_Rules, v)
			}
			GlobalConfig.RequestCertManager = make(map[int]ConfigRequestCertManager)
			for i := 0; i < len(_Rules); i++ {
				GlobalConfig.RequestCertManager[Api.CreateCertificate()] = _Rules[i]
			}
		}
		//11111111111111111
		CallJs("启动状态", base64.StdEncoding.EncodeToString([]byte(errStr)))
		return nil
	//UI获取运行端口
	case "获取运行端口":
		return app.App.Port()
	case "重置所有配置":
		s := GlobalConfig.ResetAll()
		if s != "" {
			CallJs("弹出错误提示", s)
			return false
		}
		_GlobalConfig := &UserConfig{}
		_GlobalConfig.LoadLocalFile()
		configLock.Lock()
		GlobalConfig = _GlobalConfig
		configLock.Unlock()
		CallJs("加载配置", GlobalConfig)
		CallJs("弹出成功提示", "所有配置已重置！")
		return true
	case "获取脚本代码":
		return []byte(FormatCode(string(GlobalConfig.GoScriptCode)))
	case "获取默认Go脚本代码":
		return []byte(FormatCode(string(Resource.GoCode) + Resource.ScriptAnnotation))
	case "格式化Go脚本代码":
		r1 := args.GetData("code")
		r2 := Resource.Bs64ToBs(r1)
		r := FormatCode(string(r2))
		return []byte(r)
	case "保存Go脚本代码":
		r1 := args.GetData("code")
		r2 := Resource.Bs64ToBs(r1)
		GlobalConfig.GoScriptCode = r2
		return RunCode()
	case "获取脚本日志":
		return RunCodeLog()
	case "HTTP请求获取":
		Theology := getInt(args.GetData("Theology"))
		Insert.Lock()
		currentlySelected = Theology
		Insert.Unlock()
		h := HashMap.GetRequestWeb(Theology)
		return h
	case "socket请求获取":
		Theology := getInt(args.GetData("Theology"))
		Index := getInt(args.GetData("Index")) - 1
		//HTTPHeader := args.GetData("HTTPHeader") == "true"
		if Index < 0 {
			return nil
		}
		h := HashMap.GetRequest(Theology)
		Insert.Lock()
		defer Insert.Unlock()
		if len(h.SocketData) < Index {
			return nil
		}
		return base64.StdEncoding.EncodeToString(h.SocketData[Index].Body)
	case "socket点击右键菜单":
		Theology := getInt(args.GetData("Theology"))
		//IsWs := args.GetData("IsWs") == "true"
		//IsTCP := args.GetData("IsTCP") == "true"
		Type := args.GetData("Type")
		switch Type {
		//复制选中HEX到剪辑版
		case "Selected":
			SelectedID := getInt(args.GetData("SelectedID")) - 1
			if SelectedID < 0 {
				CallJs("弹出错误提示", "复制失败:请求可能失效")
				return false
			}
			h := HashMap.GetRequest(Theology)
			if h == nil {
				CallJs("弹出错误提示", "复制失败:请求可能失效")
				return false
			}
			if len(h.SocketData) < SelectedID+1 {
				CallJs("弹出错误提示", "复制失败:选中的ID错误")
				return false
			}
			a := h.SocketData[SelectedID].Body
			err := runtime.ClipboardSetText(app.ctx, fmt.Sprintf("% X", a))
			if err == nil {
				CallJs("弹出成功提示", "当前选择的HEX数据,已复制到剪辑版")
			} else {
				CallJs("弹出错误提示", "复制失败:"+err.Error())
				return false
			}
			return true
		//复制所有HEX到剪辑版
		case "AllHEX":
			h := HashMap.GetRequest(Theology)
			if h == nil {
				CallJs("弹出错误提示", "复制失败:请求可能失效")
				return false
			}
			str := ""
			for i := 0; i < len(h.SocketData); i++ {
				a := h.SocketData[i].Body
				if h.SocketData[i].Info.Ico == "上行" {
					str += "发送的数据:" + fmt.Sprintf("% X", a) + "\r\n"
				} else if h.SocketData[i].Info.Ico == "下行" {
					str += "接收的数据:" + fmt.Sprintf("% X", a) + "\r\n"
				}
			}
			if str == "" {
				CallJs("弹出错误提示", "复制失败:无数据")
				return false
			}
			err := runtime.ClipboardSetText(app.ctx, str)
			if err == nil {
				CallJs("弹出成功提示", "当前请求的全部数据,已复制到剪辑版")
			} else {
				CallJs("弹出错误提示", "复制失败:"+err.Error())
				return false
			}
			return true
		//复制所有发送数据HEX到剪辑版
		case "sendHEX":
			h := HashMap.GetRequest(Theology)
			if h == nil {
				CallJs("弹出错误提示", "复制失败:请求可能失效")
				return false
			}
			str := ""
			for i := 0; i < len(h.SocketData); i++ {
				a := h.SocketData[i].Body
				if h.SocketData[i].Info.Ico == "上行" {
					str += "发送的数据:" + fmt.Sprintf("% X", a) + "\r\n"
				}
			}
			if str == "" {
				CallJs("弹出错误提示", "复制失败:无数据")
				return false
			}
			err := runtime.ClipboardSetText(app.ctx, str)
			if err == nil {
				CallJs("弹出成功提示", "当前请求的所有发送数据,已复制到剪辑版")
			} else {
				CallJs("弹出错误提示", "复制失败:"+err.Error())
				return false
			}
			return true
		//复制所有接收数据HEX到剪辑版
		case "recHEX":
			h := HashMap.GetRequest(Theology)
			if h == nil {
				CallJs("弹出错误提示", "复制失败:请求可能失效")
				return false
			}
			str := ""
			for i := 0; i < len(h.SocketData); i++ {
				a := h.SocketData[i].Body
				if h.SocketData[i].Info.Ico == "下行" {
					str += "接收的数据:" + fmt.Sprintf("% X", a) + "\r\n"
				}
			}
			if str == "" {
				CallJs("弹出错误提示", "复制失败:无数据")
				return false
			}
			err := runtime.ClipboardSetText(app.ctx, str)
			if err == nil {
				CallJs("弹出成功提示", "当前请求的所有接收数据,已复制到剪辑版")
			} else {
				CallJs("弹出错误提示", "复制失败:"+err.Error())
				return false
			}
			return true
		case "empty":
			return HashMap.SetSocketDataEmpty(Theology)
		}
		return false
	case "设置右键菜单配置":
		Theology := getInt(args.GetData("Theology"))
		//IsWs := args.GetData("IsWs") == "true"
		//IsTCP := args.GetData("IsTCP") == "true"
		StopSend := args.GetData("StopSend") == "true"
		StopRec := args.GetData("StopRec") == "true"
		StopALL := args.GetData("StopALL") == "true"
		return HashMap.SetOptions(Theology, StopSend, StopRec, StopALL)
	case "主动发送":
		Theology := getInt(args.GetData("Theology"))
		IsWs := args.GetData("IsWs") == "true"
		IsTCP := args.GetData("IsTCP") == "true"
		wsType := args.GetData("wsType")
		SendType := args.GetData("SendType")
		direction := args.GetData("direction")
		Data := args.GetData("Data")
		var _Bytes []byte
		{
			_tmp1, _ := base64.StdEncoding.DecodeString(Data)
			_tmp := string(_tmp1)
			if SendType == "HEX" {
				c := strings.ReplaceAll(_tmp, " ", "")
				a, e := hex.DecodeString(c)
				if e != nil {
					CallJs("弹出错误提示", "发送失败:您的 十六进制 数据有问题,请检查")
					return false
				}
				_Bytes = a
			} else if SendType == "Base64" {
				a, e := base64.StdEncoding.DecodeString(_tmp)
				if e != nil {
					CallJs("弹出错误提示", "发送失败:您的 Base64 数据有问题,请检查")
					return false
				}
				_Bytes = a
			} else if SendType == "UTF8" {
				_Bytes = _tmp1
			} else if SendType == "GBK" {
				_Bytes = Utf8ToGBK(_tmp1)
			}
			if len(_Bytes) < 1 {
				CallJs("弹出错误提示", "发送失败:无有效的发送的数据")
				return false
			}
		}
		h := HashMap.GetRequest(Theology)
		if h == nil {
			CallJs("弹出错误提示", "发送失败:请求可能失效")
			return false
		}
		BodyHash := ""
		Ico := "上行"
		{
			if direction == "Server" {
				BodyHash = "[手动发送] "
			} else {
				BodyHash = "[手动接收] "
				Ico = "下行"
			}
			if len(_Bytes) > 64 {
				BodyHash += fmt.Sprintf("% X", _Bytes[:64]) + "..."
			} else {
				BodyHash += fmt.Sprintf("% X", _Bytes)
			}
		}
		if IsWs {
			if h.WsConn == nil {
				CallJs("弹出错误提示", "发送失败:请求已断开")
				return false
			}
			MessageType := 2
			if wsType == "Text" {
				MessageType = 1
			}
			SendBool := false
			if direction == "Server" {
				SendBool = h.WsConn.SendToServer(MessageType, _Bytes)
			} else {
				SendBool = h.WsConn.SendToClient(MessageType, _Bytes)
			}
			if SendBool {
				CallJs("弹出成功提示", "发送成功")
			} else {
				CallJs("弹出错误提示", "主动发送WebSocket消息失败")
				return false
			}
			_update := &MapHash.UpdateSocketData{
				Body: _Bytes,
				Info: &MapHash.UpdateSocketList{
					Theology: Theology,
					Ico:      Ico,
					BodyHash: BodyHash,
					Length:   len(_Bytes),
					Time:     time.Now().Format("15:04:05.000"),
					WsType:   wsType,
				},
			}
			HashMap.SetSocketData(Theology, _update, Ico == "上行", len(_Bytes))
			Insert.Lock()
			_update.Info.Index = len(h.SocketData)
			isUpdateRequestInfo := currentlySelected == Theology
			if isUpdateRequestInfo {
				SocketData = append(SocketData, _update.Info)
			}
			Insert.Unlock()
			return true
		}
		if IsTCP {
			if h.TcpConn == nil {
				CallJs("弹出错误提示", "发送失败:请求已断开")
				return false
			}
			SendBool := false
			if direction == "Server" {
				SendBool = h.TcpConn.SendToServer(_Bytes) != 0
			} else {
				SendBool = h.TcpConn.SendToClient(_Bytes) != 0
			}
			if SendBool {
				CallJs("弹出成功提示", "发送成功")
			} else {
				CallJs("弹出错误提示", "主动发送 TCP 消息失败")
				return false
			}
			_update := &MapHash.UpdateSocketData{
				Body: _Bytes,
				Info: &MapHash.UpdateSocketList{
					Theology: Theology,
					Ico:      Ico,
					BodyHash: BodyHash,
					Length:   len(_Bytes),
					Time:     time.Now().Format("15:04:05.000"),
					WsType:   wsType,
				},
			}
			HashMap.SetSocketData(Theology, _update, Ico == "上行", len(_Bytes))
			Insert.Lock()
			_update.Info.Index = len(h.SocketData)
			isUpdateRequestInfo := currentlySelected == Theology
			if isUpdateRequestInfo {
				SocketData = append(SocketData, _update.Info)
			}
			Insert.Unlock()
			return true
		}
		fmt.Println(Theology, IsWs, IsTCP, wsType, SendType, direction, Data)
		return nil
	case "保存响应图片":
		Theology := getInt(args.GetData("Theology"))
		h := HashMap.GetRequest(Theology)
		if h != nil {
			_path, err := CommAnd.GetDesktopPath()
			if _path == "" {
				CallJs("弹出错误提示", "获取桌面路径失败:"+err.Error())
				return false
			}
			timestamp10 := strconv.FormatInt(time.Now().Unix(), 10)
			_path = strings.ReplaceAll(_path+"/"+timestamp10+"."+args.GetData("type"), "\\", "/")
			err = os.WriteFile(_path, h.Response.Body, 777)
			if err != nil {
				CallJs("弹出错误提示", "写入文件时出错:"+err.Error())
				return false
			}
			CallJs("弹出成功提示", "保存文件成功："+_path)
		} else {
			CallJs("弹出错误提示", "请求不存在")
		}
		return true
	case "保存修改数据":
		Theology := getInt(args.GetData("Theology"))
		Type := args.GetData("Type")
		Tabs := args.GetData("Tabs")
		Coding := args.GetData("UTF8") == "true"
		h := HashMap.GetRequest(Theology)
		if h == nil {
			CallJs("弹出错误提示", "修改数据失败:请求可能失效")
			return false
		}
		if Type == "Request" {
			if h.Conn == nil {
				CallJs("弹出错误提示", "修改数据失败:请求可能失效")
				return false
			}
			switch Tabs {
			case "Raw", "Hex":
				Data, _ := base64.StdEncoding.DecodeString(args.GetData("Data"))
				if !Coding {
					Data = Utf8ToGBK(Data)
				}
				//修改原始数据
				{
					DataStr := string(Data)
					array := strings.Split(DataStr, "\r\n\r\n")
					if len(array) < 1 {
						CallJs("弹出错误提示", "修改数据失败:无法正确识别修改后的数据")
						return false
					}
					array1 := strings.Split(array[0], "\r\n")
					if len(array1) < 1 {
						CallJs("弹出错误提示", "修改数据失败:无法正确识别修改后的数据")
						return false
					}
					array2 := strings.Split(array1[0], " ")
					if len(array2) < 3 {
						CallJs("弹出错误提示", "修改数据失败:无法正确识别修改后的数据")
						return false
					}
					Method := strings.ToUpper(array2[0])
					URL, err := url.Parse(array2[1])
					if err != nil {
						CallJs("弹出错误提示", "修改数据失败:无法正确识别修改后的URL error:"+err.Error())
						return false
					}
					header := make(http.Header)
					for i := 1; i < len(array1)-1; i++ {
						array2 = strings.Split(array1[i], ":")
						name := ""
						value := ""
						if len(array2) >= 1 {
							name = array2[0]
						}
						value = strings.TrimSpace(strings.Replace(array1[i], name+":", "", 1))
						header[name] = []string{value}
					}
					h.Body = []byte(strings.Join(array[1:], "\r\n\r\n"))
					if h.Conn.Request.Body != nil {
						_ = h.Conn.Request.Body.Close()
					}
					h.Conn.Request.Body = io.NopCloser(bytes.NewBuffer(h.Body))
					h.Conn.Request.Header = header
					h.Header = header
					h.URL = URL.String()
					h.Method = Method
					h.Conn.Request.Method = Method
					h.Conn.Request.URL = URL
					CallJs("弹出成功提示", "修改请求提交的包体成功")
				}
				return true
			case "Headers":
				//修改协议头
				{
					header := make(http.Header)
					num := args.GetNum("Data")
					for i := 0; i < num; i++ {
						name := args.GetData("Data[" + strconv.Itoa(i) + "].名称")
						value := args.GetData("Data[" + strconv.Itoa(i) + "].值")
						if !Coding {
							name1 := string(Utf8ToGBK([]byte(name)))
							value1 := string(Utf8ToGBK([]byte(value)))
							header[name1] = []string{value1}
						} else {
							header[name] = []string{value}
						}

					}
					h.Conn.Request.Header = header
					h.Header = header
					CallJs("弹出成功提示", "修改请求提交的协议头成功")
				}
				return true
			case "Args":
				//修改URL中的参数
				{
					RawQuery := ""

					num := args.GetNum("Data")
					for i := 0; i < num; i++ {
						name := args.GetData("Data[" + strconv.Itoa(i) + "].名称")
						value := args.GetData("Data[" + strconv.Itoa(i) + "].值")
						_utf8 := args.GetData("Data["+strconv.Itoa(i)+"].编码") == "true"
						if !_utf8 {
							name = string(Utf8ToGBK([]byte(name)))
							value = string(Utf8ToGBK([]byte(value)))
						}
						if RawQuery == "" {
							RawQuery = name + "=" + url.QueryEscape(value)
						} else {
							RawQuery += "&" + name + "=" + url.QueryEscape(value)
						}
					}
					h.Conn.Request.URL.RawQuery = RawQuery
					h.URL = h.Conn.Request.URL.String()
					CallJs("弹出成功提示", "修改请求URL中的参数成功")
				}
				return true
			case "Cookies":
				//修改Cookies
				{
					Cookie := ""
					num := args.GetNum("Data")
					for i := 0; i < num; i++ {
						name := args.GetData("Data[" + strconv.Itoa(i) + "].名称")
						value := args.GetData("Data[" + strconv.Itoa(i) + "].值")
						if !Coding {
							name = string(Utf8ToGBK([]byte(name)))
							value = string(Utf8ToGBK([]byte(value)))
						}
						if Cookie == "" {
							Cookie = name + "=" + value
						} else {
							Cookie += "; " + name + "=" + value
						}
					}
					h.Conn.Request.Header["Cookie"] = []string{Cookie}
					CallJs("弹出成功提示", "修改请求Cookies成功")
				}
				return true
			case "BodyArgs":
				//修改Post中的参数
				{
					PostData := ""
					num := args.GetNum("Data")
					for i := 0; i < num; i++ {
						name := args.GetData("Data[" + strconv.Itoa(i) + "].名称")
						value := args.GetData("Data[" + strconv.Itoa(i) + "].值")
						_utf8 := args.GetData("Data["+strconv.Itoa(i)+"].编码") == "true"
						if !_utf8 {
							name = string(Utf8ToGBK([]byte(name)))
							value = string(Utf8ToGBK([]byte(value)))
						}
						if PostData == "" {
							PostData = name + "=" + url.QueryEscape(value)
						} else {
							PostData += "&" + name + "=" + url.QueryEscape(value)
						}
					}
					{
						h.Body = []byte(PostData)
						if h.Conn.Request.Body != nil {
							_ = h.Conn.Request.Body.Close()
						}
						h.Conn.Request.Body = io.NopCloser(bytes.NewBuffer(h.Body))
					}
					CallJs("弹出成功提示", "修改请求POST数据成功")
					return PostData
				}
			case "Body", "Json":
				Data, _ := base64.StdEncoding.DecodeString(args.GetData("Data"))
				if !Coding {
					Data = Utf8ToGBK(Data)
				}
				{
					h.Body = Data
					if h.Conn.Request.Body != nil {
						_ = h.Conn.Request.Body.Close()
					}
					h.Conn.Request.Body = io.NopCloser(bytes.NewBuffer(h.Body))
					CallJs("弹出成功提示", "修改请求POST数据成功")
				}
				return true
			}
		} else if Type == "Response" {
			if h.Response.Conn == nil {
				CallJs("弹出错误提示", "修改数据失败:请求可能失效")
				return false
			}
			defer func() {
				CallJs("更新响应", &UpdateCurrentResponse{
					Theology:  Theology,
					Header:    h.Response.Header,
					Body:      h.Response.Body,
					StateText: http.StatusText(h.Response.StateCode),
					StateCode: h.Response.StateCode,
					Break:     true,
				})
			}()
			switch Tabs {
			case "Raw", "Hex":
				{
					Data, _ := base64.StdEncoding.DecodeString(args.GetData("Data"))
					if !Coding {
						Data = Utf8ToGBK(Data)
					}
					DataStr := string(Data)
					DataArray := strings.Split(DataStr, "\r\n\r\n")
					if len(DataArray) < 1 {
						CallJs("弹出错误提示", "保存响应数据失败,响应体不正确！！")
						return false
					}
					_header := DataArray[0]
					array2 := strings.Split(_header, "\r\n")
					if len(array2) < 1 {
						CallJs("弹出错误提示", "保存响应数据失败,响应体不正确！！")
						return false
					}
					DataBody := []byte(strings.Join(DataArray[1:], "\r\n\r\n"))
					array1 := strings.Split(array2[0], " ")
					if len(array1) < 3 {
						CallJs("弹出错误提示", "保存响应数据失败,响应体不正确！！")
						return false
					}
					h.Response.StateCode, _ = strconv.Atoi(array1[1])

					h.Response.Header = make(http.Header)
					for i := 1; i < len(array2)-1; i++ {
						array1 = strings.Split(array2[i], ":")
						if len(array1) >= 1 {
							name := array1[0]
							value := strings.TrimSpace(strings.Replace(array2[i], name+":", "", 1))
							if len(h.Response.Header[name]) > 0 {
								h.Response.Header[name] = append(h.Response.Header[name], value)
							} else {
								h.Response.Header[name] = []string{value}
							}
						}
					}
					h.Response.Body = DataBody
					if h.Response.Conn.Response.Body != nil {
						_ = h.Response.Conn.Response.Body.Close()
					}
					h.Response.Conn.Response.Body = io.NopCloser(bytes.NewBuffer(h.Response.Body))
					h.Response.Conn.Response.Header = h.Response.Header
					h.Response.Conn.Response.StatusCode = h.Response.StateCode
					h.Response.Conn.Response.Status = http.StatusText(h.Response.StateCode)

				}
				CallJs("弹出成功提示", "修改响应数据保存成功！！")
				return true
			case "Headers":
				//修改协议头
				{
					header := make(http.Header)
					num := args.GetNum("Data")
					for i := 0; i < num; i++ {
						name := args.GetData("Data[" + strconv.Itoa(i) + "].名称")
						value := args.GetData("Data[" + strconv.Itoa(i) + "].值")
						if !Coding {
							name = string(Utf8ToGBK([]byte(name)))
							value = string(Utf8ToGBK([]byte(value)))
						}
						if len(header[name]) > 0 {
							header[name] = append(header[name], value)
						} else {
							header[name] = []string{value}
						}
					}
					h.Response.Header = header
					h.Response.Conn.Response.Header = header
					CallJs("弹出成功提示", "修改响应协议头成功！！")
				}
				return true
			case "Body", "Json":
				Data, _ := base64.StdEncoding.DecodeString(args.GetData("Data"))
				if !Coding {
					Data = Utf8ToGBK(Data)
				}
				//修改URL中的参数
				{
					h.Response.Body = Data
					if h.Response.Conn.Response.Body != nil {
						_ = h.Response.Conn.Response.Body.Close()
					}
					h.Response.Conn.Response.Body = io.NopCloser(bytes.NewBuffer(h.Response.Body))
					CallJs("弹出成功提示", "修改响应Body成功！！")
				}
				return true
			case "Cookies":
				//修改Cookies
				{
					var Cookie []string
					num := args.GetNum("Data")
					for i := 0; i < num; i++ {
						name := args.GetData("Data[" + strconv.Itoa(i) + "].名称")
						value := args.GetData("Data[" + strconv.Itoa(i) + "].值")
						value2 := args.GetData("Data[" + strconv.Itoa(i) + "].其他值")
						if !Coding {
							name = string(Utf8ToGBK([]byte(name)))
							value = string(Utf8ToGBK([]byte(value)))
							value2 = string(Utf8ToGBK([]byte(value2)))
						}
						Cookie = append(Cookie, name+"="+value+"; "+value2)
					}
					h.Response.Conn.Response.Header["Set-Cookie"] = Cookie
					//h.Response.Header = h.Response.Conn.Response.Header
					CallJs("弹出成功提示", "修改响应Cookies成功！！")
				}
				return true
			}
		} else {
			CallJs("弹出错误提示", "保存修改数据失败,出现未知错误！！")
			return false
		}
		fmt.Println(Type, Tabs, Theology, Coding)
		return true
	case "protobufToJson":
		skip := getInt(args.GetData("skip"))
		pb := Resource.Bs64ToBs(args.GetData("Data"))
		return _PbToJson(pb, skip)
	case "断点点击":
		Theology := getInt(args.GetData("Theology"))
		NextBreak := getInt(args.GetData("NextBreak"))
		h := HashMap.GetRequest(Theology)
		if h != nil {
			h.Break = uint8(NextBreak)
			h.Wait.Done()
		}
		return true
	case "重置颜色列表":
		dark := args.GetData("dark") == "true"
		_TmpLock.Lock()
		GlobalConfig.ResetColorConfig(dark)
		_ = GlobalConfig.saveToFile()
		_TmpLock.Unlock()
		return GlobalConfig.ColorConfig
	case "保存上游代理使用规则":
		code := strings.ReplaceAll(args.GetData("Data"), "\\\\", "\\")
		_TmpLock.Lock()
		GlobalConfig.GlobalProxyRules = code
		_ = GlobalConfig.saveToFile()
		_TmpLock.Unlock()
		return app.App.CompileProxyRegexp(code) == nil
	case "设置上游代理":
		code := strings.ReplaceAll(args.GetData("Data"), "\\\\", "\\")
		_TmpLock.Lock()
		GlobalConfig.GlobalProxy = code
		_ = GlobalConfig.saveToFile()
		_TmpLock.Unlock()
		Set := args.GetData("Set") == "true"
		if Set {
			return app.App.SetGlobalProxy(code)
		}
		app.App.SetGlobalProxy("socket5://:@")
		return true
	case "保存强制TCP使用规则":
		code := strings.ReplaceAll(args.GetData("Data"), "\\\\", "\\")
		if code == "MustTcp" {
			_TmpLock.Lock()
			GlobalConfig.MustTcp.Open = true
			_ = GlobalConfig.saveToFile()
			_TmpLock.Unlock()
			app.App.MustTcp(true)
			return true
		} else if code == "CancelMustTcp" {
			_TmpLock.Lock()
			GlobalConfig.MustTcp.Open = false
			_ = GlobalConfig.saveToFile()
			_TmpLock.Unlock()
			_ = app.App.SetMustTcpRegexp("ALL.ALL.ALL.ALL")
			app.App.MustTcp(false)
			return true
		}
		_TmpLock.Lock()
		GlobalConfig.MustTcp.Rules = code
		_ = GlobalConfig.saveToFile()
		_TmpLock.Unlock()
		return app.App.SetMustTcpRegexp(code) == nil
	case "创建证书":
		domain := args.GetData("domain")
		country := args.GetData("country")
		company := args.GetData("company")
		department := args.GetData("department")
		province := args.GetData("province")
		city := args.GetData("city")
		outTime := getInt(args.GetData("outTime"))
		if domain == "" {
			CallJs("弹出错误提示", "请输入 '颁发给' 不可为空,一般为域名")
			return false
		}
		{
			if country == "" {
				country = "CN"
			}
			if company == "" {
				company = "Sunny"
			}
			if department == "" {
				department = "Sunny"
			}
			if province == "" {
				province = "Beijing"
			}
			if city == "" {
				city = "Beijing"
			}
			if outTime <= 0 {
				outTime = 1
			}
		}
		id := Api.CreateCertificate()
		defer Api.RemoveCertificate(id)
		ok := Api.CreateCA(id, country, company, department, province, domain, city, 2048, outTime)
		if !ok {
			CallJs("弹出错误提示", "创建证书失败")
			return false
		}
		path, e := CommAnd.GetDesktopPath()
		fileName := strconv.FormatInt(time.Now().Unix(), 10)
		path1 := path + "/" + fileName + ".cer"
		path2 := path + "/" + fileName + ".key"
		if e != nil {
			CallJs("弹出错误提示", "创建证书失败:没有获取到桌面路径")
			return false
		}
		_ = os.Remove(path1)
		_ = os.Remove(path2)
		ca := Api.ExportCA(id)
		key := Api.ExportKEY(id)
		e = os.WriteFile(path1, []byte(ca), 777)
		if e != nil {
			CallJs("弹出错误提示", "创建证书失败:出错CA文件到桌面失败")
			return false
		}
		e = os.WriteFile(path2, []byte(key), 777)
		if e != nil {
			CallJs("弹出错误提示", "创建证书失败:出错KEY文件到桌面失败")
			return false
		}
		CallJsAlert("创建证书成功：已储存在桌面", "储存文件路径:\n\tCA文件: "+path1+"\n\n\tKEY文件: "+path2)
		return true
	case "安装CA证书":
		CaFilePath := args.GetData("CaFilePath")
		CaFilePath = strings.ReplaceAll(CaFilePath, "\\\\", "\\")
		bs, e := os.ReadFile(CaFilePath)
		if e != nil {
			CallJs("弹出错误提示", "读取CA文件失败,请检查文件是否存在")
			return false
		}
		CallJsAlert("安装结果：", CommAnd.InstallCert(bs))
		return ""
	case "导入证书":
		CaFilePath := args.GetData("CaFilePath")
		KeyFilePath := args.GetData("KeyFilePath")
		CaFilePath = strings.ReplaceAll(CaFilePath, "\\\\", "\\")
		KeyFilePath = strings.ReplaceAll(KeyFilePath, "\\\\", "\\")
		_TmpLock.Lock()
		GlobalConfig.Cert.CaPath = CaFilePath
		GlobalConfig.Cert.KeyPath = KeyFilePath
		_ = GlobalConfig.saveToFile()
		_TmpLock.Unlock()
		id := Api.CreateCertificate()
		defer Api.RemoveCertificate(id)
		if !Api.LoadX509KeyPair(id, CaFilePath, KeyFilePath) {
			bs1, e := os.ReadFile(CaFilePath)
			if e != nil {
				CallJs("弹出错误提示", "读取CA文件失败,请检查文件是否存在")
				return false
			}
			bs2, e := os.ReadFile(KeyFilePath)
			if e != nil {
				CallJs("弹出错误提示", "读取KEY文件失败,请检查文件是否存在")
				return false
			}
			if !Api.LoadX509Certificate(id, "", string(bs1), string(bs2)) {
				CallJs("弹出错误提示", "加载自定义证书失败！！")
				return false
			}
		}
		app.App.Error = nil
		e := app.App.SetCert(id).Error
		if e == nil {
			_TmpLock.Lock()
			GlobalConfig.Cert.Default = false
			_ = GlobalConfig.saveToFile()
			_TmpLock.Unlock()
			CallJs("弹出成功提示", "已应用自定义证书")
			return true
		}
		CallJs("弹出错误提示", "应用自定义证书失败:"+e.Error())
		return false
	case "安装默认证书":
		CallJsAlert("安装结果：", CommAnd.InstallCert([]byte(public.RootCa)))
		return ""
	case "保存默认证书到桌面":
		path, e := CommAnd.GetDesktopPath()
		if e != nil {
			CallJs("弹出错误提示", "未能获取到桌面路径！！")
			return false
		}
		path1 := path + "/SunnyNet.cer"
		_ = os.Remove(path1)
		if os.WriteFile(path1, []byte(public.RootCa), 777) == nil {
			CallJs("弹出成功提示", "保存默认证书文件成功")
			return true
		}
		return false
	case "应用默认证书":
		id := Api.CreateCertificate()
		defer Api.RemoveCertificate(id)
		if !Api.LoadX509Certificate(id, public.NULL, public.RootCa, public.RootKey) {
			CallJs("弹出错误提示", "加载默认自定义证书！！")
			return false
		}
		_TmpLock.Lock()
		GlobalConfig.Cert.Default = true
		_ = GlobalConfig.saveToFile()
		_TmpLock.Unlock()
		return true
	case "修改端口号":
		Port := getInt(args.GetData("Port"))
		a := &StartMsg{}
		if Port < 1 || Port > 65535 {
			a.Err = "Port error"
			return a
		}
		app.App.Close()
		app.App.SetPort(Port)
		app.App.Start()
		if app.App.Error == nil {
			a.Ok = true
			CallJs("启动状态", "")
		} else {
			a.Err = base64.StdEncoding.EncodeToString([]byte(app.App.Error.Error()))
			CallJs("启动状态", a.Err)
		}
		GlobalConfig.Port = Port
		_ = GlobalConfig.saveToFile()
		return a
	case "禁止UDP":
		_TmpLock.Lock()
		DisableUDP = args.GetData("DisableUDP") == "true"
		GlobalConfig.DisableUDP = DisableUDP
		_ = GlobalConfig.saveToFile()
		_TmpLock.Unlock()
		return true
	case "禁止缓存":
		_TmpLock.Lock()
		DisableCache = args.GetData("DisableBrowserCache") == "true"
		GlobalConfig.DisableCache = DisableCache
		_ = GlobalConfig.saveToFile()
		_TmpLock.Unlock()
		return true
	case "身份验证模式":
		_TmpLock.Lock()
		authentication := args.GetData("authentication") == "true"
		GlobalConfig.Authentication = authentication
		_ = GlobalConfig.saveToFile()
		_TmpLock.Unlock()
		app.App.Socket5VerifyUser(authentication)
		return true
	case "更新身份验证账号信息":
		_TmpLock.Lock()
		for _, user := range SocketAuthentication {
			app.App.Socket5DelUser(user)
		}
		SocketAuthentication = make([]string, 0)
		_TmpLock.Unlock()
		num := args.GetNum("Data")
		UserInfo := make(map[string]string)
		for i := 0; i < num; i++ {
			name := args.GetData("Data[" + strconv.Itoa(i) + "].账号")
			value := args.GetData("Data[" + strconv.Itoa(i) + "].密码")
			UserInfo[name] = value
			app.App.Socket5AddUser(name, value)
			_TmpLock.Lock()
			SocketAuthentication = append(SocketAuthentication, name)
			_TmpLock.Unlock()
		}
		_TmpLock.Lock()
		GlobalConfig.AuthenticationUserInfo = UserInfo
		_ = GlobalConfig.saveToFile()
		_TmpLock.Unlock()
		return true
	case "":
		return ""
	default:
		return RequestCertificate(command, args)
	}
}

var DisableUDP = false
var DisableCache = false
var SocketAuthentication []string

type StartMsg struct {
	Ok  bool   `json:"ok"`
	Err string `json:"err"`
}

func saveToFile(Path string, All bool, TheologyArray []int) bool {
	return HashMap.SaveToFile(Path, All, TheologyArray, SetStatusText)
}
