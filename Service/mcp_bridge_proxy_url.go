package Service

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

func proxyWayURLConvention() map[string]any {
	return map[string]any{
		"supportedTypes": []string{"http", "https", "socks5"},
		"supportedSchemes": []string{"http", "https", "socks5", "socket"},
		"auth":             "无账号：scheme://host:port；有账号：scheme://user:pass@host:port",
		"note":             "仅支持 HTTP/HTTPS/SOCKS5 代理；socket:// 为 SOCKS5 的引擎写法（与 socks5:// 同类）",
		"examples": []string{
			"http://192.168.1.1:8080",
			"http://admin:123456@192.168.1.1:8080",
			"https://proxy.example.com:443",
			"socks5://127.0.0.1:1080",
			"socks5://user:pass@127.0.0.1:1080",
			"socket://192.168.31.1:4321",
			"socket://admin:123456@192.168.31.1:4321",
		},
	}
}

// validateProxyWayURL 校验上游代理 URL（http/https/socks5；socket 视为 socks5）。
func validateProxyWayURL(raw string) (normalized string, proxyType string, err error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", "", errors.New("url 不能为空")
	}
	if !strings.Contains(raw, "://") {
		return "", "", errors.New("url 须带协议前缀：http://、https://、socks5:// 或 socket://（SOCKS5）")
	}
	u, err := url.Parse(raw)
	if err != nil {
		return "", "", fmt.Errorf("url 解析失败: %w", err)
	}
	scheme := strings.ToLower(strings.TrimSpace(u.Scheme))
	switch scheme {
	case "http":
		proxyType = "http"
	case "https":
		proxyType = "https"
	case "socks5", "socket":
		proxyType = "socks5"
	default:
		return "", "", fmt.Errorf("不支持的代理协议 %q，仅支持 http、https、socks5（或 socket:// 表示 SOCKS5）", u.Scheme)
	}
	if strings.TrimSpace(u.Host) == "" {
		return "", "", errors.New("url 缺少 host:port")
	}
	return raw, proxyType, nil
}

// resolveProxyWayURL 支持 url 或 scheme+host+port+username+password 组装。
func resolveProxyWayURL(m map[string]any) (string, string, error) {
	raw := strings.TrimSpace(argString(m, "url"))
	if raw == "" {
		scheme := strings.ToLower(strings.TrimSpace(argString(m, "scheme")))
		if scheme == "" {
			scheme = strings.ToLower(strings.TrimSpace(argString(m, "type")))
		}
		host := strings.TrimSpace(argString(m, "host"))
		port := strings.TrimSpace(argString(m, "port"))
		if host != "" {
			switch scheme {
			case "", "socks5", "socket", "s5":
				scheme = "socks5"
			case "http", "https":
			default:
				return "", "", fmt.Errorf("scheme/type %q 不支持，须为 http、https、socks5", scheme)
			}
			if port == "" {
				return "", "", errors.New("使用 host 时须同时提供 port")
			}
			raw = scheme + "://" + host + ":" + port
		}
	}
	user := strings.TrimSpace(argString(m, "username"))
	if user == "" {
		user = strings.TrimSpace(argString(m, "user"))
	}
	pass := argString(m, "password")
	if pass == "" {
		pass = argString(m, "pass")
	}
	if raw != "" && (user != "" || pass != "") {
		u, err := url.Parse(raw)
		if err != nil {
			return "", "", err
		}
		if u.User != nil {
			return "", "", errors.New("url 已含账号密码时不要再传 username/password")
		}
		u.User = url.UserPassword(user, pass)
		raw = u.String()
	}
	return validateProxyWayURL(raw)
}
