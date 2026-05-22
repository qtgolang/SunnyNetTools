package Tools

import (
	"changeme/Service/Config"
	"fmt"
	"github.com/qtgolang/SunnyNet/SunnyNet"
	"sort"
	"strings"
	"sync"
)

type RequestCert struct {
	App *SunnyNet.Sunny
}

var mapCertLock sync.Mutex

// CreateRequestCert 创建请求证书
func (g *RequestCert) CreateRequestCert() int {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	mapCertLock.Lock()
	defer mapCertLock.Unlock()
	Config.CertID++
	id := Config.CertID
	Config.Config.RequestCert[id] = &Config.CertInfo{
		Cert:   SunnyNet.NewCertManager(),
		DoMain: "",
		Path:   "",
		Pass:   "",
		Note:   "",
		ID:     id,
		LoadOk: false,
	}
	return id
}

// RequestCertSetFile 设置载入证书文件
func (g *RequestCert) RequestCertSetFile(id int, role, DoMain, file, pass, note string) (res string) {
	defer func() {
		if err := recover(); err != nil {
			res = fmt.Sprintf("%v", err)
		}
	}()
	mapCertLock.Lock()
	defer mapCertLock.Unlock()
	obj := Config.Config.RequestCert[id]
	if obj == nil {
		return "RequestCert not found"
	}
	_role := 3
	switch role {
	case "解析及发送":
		_role = 2
		break
	case "仅解析":
		_role = 3
		break
	default:
		_role = 1
		break
	}
	if obj.DoMain != "" {
		g.App.DelHttpCertificate(obj.DoMain)
	}
	f := strings.ToLower(file)
	if strings.HasSuffix(f, ".p12") || strings.HasSuffix(f, ".pkcs12") {
		if obj.Cert.LoadP12Certificate(file, pass) {
			obj.Path = file
			obj.DoMain = DoMain
			obj.Pass = pass
			obj.Note = note
			obj.Role = role
			obj.LoadOk = true
			g.App.AddHttpCertificate(DoMain, obj.Cert, uint8(_role))
			Config.Config.Save()
			return "ok"
		}
		return "P12 载入失败"
	}
	if obj.Cert.AddCertPoolPath(file) {
		obj.Path = file
		obj.DoMain = DoMain
		obj.Pass = pass
		obj.Note = note
		obj.Role = role
		obj.LoadOk = true
		g.App.AddHttpCertificate(DoMain, obj.Cert, uint8(_role))
		Config.Config.Save()
		return "ok"
	}
	return "Load CertFile failed"
}

// RequestCertGetCommonName 获取证书名称
func (g *RequestCert) RequestCertGetCommonName(id int) string {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	mapCertLock.Lock()
	defer mapCertLock.Unlock()
	obj := Config.Config.RequestCert[id]
	if obj == nil {
		return ""
	}
	return obj.Cert.GetCommonName()
}

// RequestCertRemove 删除证书
func (g *RequestCert) RequestCertRemove(id int) {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	mapCertLock.Lock()
	defer mapCertLock.Unlock()
	obj := Config.Config.RequestCert[id]
	if obj == nil {
		return
	}
	g.App.DelHttpCertificate(obj.DoMain)
	delete(Config.Config.RequestCert, id)
}

// RequestList 获取全部列表
func (g *RequestCert) RequestList() []Config.CertInfo {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	mapCertLock.Lock()
	defer mapCertLock.Unlock()
	var keys []int
	for i, v := range Config.Config.RequestCert {
		if v.Path != "" {
			keys = append(keys, i)
		}
	}
	sort.Ints(keys)
	var array []Config.CertInfo
	for _, k := range keys {
		obj := Config.Config.RequestCert[k]
		array = append(array, *obj)
	}
	return array
}
