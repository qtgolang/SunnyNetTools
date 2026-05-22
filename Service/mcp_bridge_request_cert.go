package Service

import (
	"changeme/Service/Config"
	"errors"
	"fmt"
	"os"
	"strings"
)

var requestCertRoles = map[string]bool{
	"解析及发送": true,
	"仅解析":   true,
	"仅发送":   true,
}

func isRequestCertP12(path string) bool {
	f := strings.ToLower(strings.TrimSpace(path))
	return strings.HasSuffix(f, ".p12") || strings.HasSuffix(f, ".pkcs12")
}

func isRequestCertPEM(path string) bool {
	f := strings.ToLower(strings.TrimSpace(path))
	return strings.HasSuffix(f, ".pem") || strings.HasSuffix(f, ".cer")
}

func normalizeRequestCertRole(role string) (string, error) {
	role = strings.TrimSpace(role)
	if role == "" {
		return "解析及发送", nil
	}
	if !requestCertRoles[role] {
		return "", fmt.Errorf("role 须为 解析及发送、仅解析 或 仅发送")
	}
	return role, nil
}

func normalizeRequestCertPassword(pass, path string) (string, error) {
	pass = strings.TrimSpace(pass)
	if pass == "" || pass == "该文件类型无需密码" || strings.HasPrefix(pass, "双击") {
		if isRequestCertP12(path) {
			return "", errors.New("P12/PKCS12 证书须填写 password")
		}
		return "", nil
	}
	return pass, nil
}

func certInfoToMCP(c Config.CertInfo) map[string]any {
	status := "未载入"
	if c.LoadOk {
		status = "已载入"
	}
	pass := c.Pass
	if pass == "" && isRequestCertPEM(c.Path) {
		pass = "该文件类型无需密码"
	}
	certType := "PEM/CER"
	if isRequestCertP12(c.Path) {
		certType = "P12/PKCS12"
	}
	return map[string]any{
		"id":       c.ID,
		"domain":   c.DoMain,
		"certPath": c.Path,
		"password": pass,
		"note":     c.Note,
		"role":     c.Role,
		"loadOk":   c.LoadOk,
		"status":   status,
		"certType": certType,
	}
}

func requestCertListForMCP(app *AppMain) map[string]any {
	list := app.RequestCert.RequestList()
	certs := make([]map[string]any, 0, len(list))
	for _, c := range list {
		certs = append(certs, certInfoToMCP(c))
	}
	return map[string]any{
		"certs": certs,
		"total": len(certs),
		"convention": map[string]any{
			"roles":     []string{"解析及发送", "仅解析", "仅发送"},
			"certTypes": map[string]any{
				"P12/PKCS12": map[string]any{"extensions": []string{".p12", ".pkcs12"}, "passwordRequired": true},
				"PEM/CER":    map[string]any{"extensions": []string{".pem", ".cer"}, "passwordRequired": false},
			},
		},
	}
}

func bridgeRequestCertList(app *AppMain) (any, error) {
	return requestCertListForMCP(app), nil
}

func bridgeRequestCertAdd(app *AppMain, m map[string]any) (any, error) {
	path := strings.TrimSpace(argString(m, "certPath"))
	if path == "" {
		path = strings.TrimSpace(argString(m, "path"))
	}
	if path == "" {
		return nil, errors.New("certPath 必填（本地 .p12/.pkcs12/.pem/.cer 路径）")
	}
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("证书文件不存在: %s", path)
	}
	if !isRequestCertP12(path) && !isRequestCertPEM(path) {
		return nil, errors.New("certPath 扩展名须为 .p12/.pkcs12/.pem/.cer")
	}
	domain := strings.TrimSpace(argString(m, "domain"))
	if domain == "" {
		domain = strings.TrimSpace(argString(m, "DoMain"))
	}
	if domain == "" {
		return nil, errors.New("domain 必填（绑定 Host，如 api.example.com）")
	}
	role, err := normalizeRequestCertRole(argString(m, "role"))
	if err != nil {
		return nil, err
	}
	pass, err := normalizeRequestCertPassword(argString(m, "password"), path)
	if err != nil {
		return nil, err
	}
	note := strings.TrimSpace(argString(m, "note"))

	id := app.RequestCert.CreateRequestCert()
	res := app.RequestCert.RequestCertSetFile(id, role, domain, path, pass, note)
	if res != "ok" {
		app.RequestCert.RequestCertRemove(id)
		if res == "P12 载入失败" {
			return nil, errors.New("P12 载入失败（请检查 password 与文件）")
		}
		return nil, fmt.Errorf("证书载入失败: %s", res)
	}
	emitMCPRequestCertReload()
	return map[string]any{"ok": true, "id": id, "status": "已载入", "cert": certInfoToMCP(*Config.Config.RequestCert[id])}, nil
}

func bridgeRequestCertDelete(app *AppMain, m map[string]any) (any, error) {
	id := argInt(m, "id", 0)
	if id <= 0 {
		return nil, errors.New("id 必填")
	}
	obj := Config.Config.RequestCert[id]
	if obj == nil {
		return nil, fmt.Errorf("证书 id %d 不存在", id)
	}
	app.RequestCert.RequestCertRemove(id)
	Config.Config.Save()
	emitMCPRequestCertReload()
	return map[string]any{"ok": true, "id": id}, nil
}

func bridgeRequestCertUpdate(app *AppMain, m map[string]any) (any, error) {
	id := argInt(m, "id", 0)
	if id <= 0 {
		return nil, errors.New("id 必填")
	}
	obj := Config.Config.RequestCert[id]
	if obj == nil {
		return nil, fmt.Errorf("证书 id %d 不存在", id)
	}

	domain := obj.DoMain
	if v, ok := m["domain"]; ok {
		domain = strings.TrimSpace(fmt.Sprint(v))
		if domain == "" {
			return nil, errors.New("domain 不能为空")
		}
	} else if v, ok := m["DoMain"]; ok {
		domain = strings.TrimSpace(fmt.Sprint(v))
	}

	note := obj.Note
	if v, ok := m["note"]; ok {
		note = strings.TrimSpace(fmt.Sprint(v))
	}

	role := obj.Role
	if v := strings.TrimSpace(argString(m, "role")); v != "" {
		var err error
		role, err = normalizeRequestCertRole(v)
		if err != nil {
			return nil, err
		}
	}

	newPath := strings.TrimSpace(argString(m, "certPath"))
	if newPath == "" {
		newPath = strings.TrimSpace(argString(m, "path"))
	}
	path := obj.Path
	if newPath != "" {
		if _, err := os.Stat(newPath); err != nil {
			return nil, fmt.Errorf("证书文件不存在: %s", newPath)
		}
		path = newPath
	}

	pass := obj.Pass
	if v, ok := m["password"]; ok {
		var err error
		pass, err = normalizeRequestCertPassword(fmt.Sprint(v), path)
		if err != nil {
			return nil, err
		}
	} else if isRequestCertPEM(path) {
		pass = ""
	}

	if path == "" {
		obj.DoMain = domain
		obj.Note = note
		obj.Role = role
		Config.Config.Save()
		emitMCPRequestCertReload()
		return map[string]any{"ok": true, "id": id, "status": "未载入", "cert": certInfoToMCP(*obj)}, nil
	}

	res := app.RequestCert.RequestCertSetFile(id, role, domain, path, pass, note)
	status := "已载入"
	if res != "ok" {
		status = "载入失败"
		if res == "P12 载入失败" {
			return nil, errors.New("P12 载入失败（请检查 password 与文件）")
		}
		return nil, fmt.Errorf("证书更新失败: %s", res)
	}
	emitMCPRequestCertReload()
	return map[string]any{"ok": true, "id": id, "status": status, "cert": certInfoToMCP(*Config.Config.RequestCert[id])}, nil
}
