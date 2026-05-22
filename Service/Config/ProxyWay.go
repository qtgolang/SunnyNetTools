package Config

type ProxyWayInfo struct {
	State string
	Note  string
	URL   string
	ID    int
}

func (f *config) initProxy() {
	//恢复上游代理网关列表
	{
		m := f.ProxyWay
		f.ProxyWay = make(map[int]*ProxyWayInfo)
		for _, v := range m {
			if v == nil {
				continue
			}
			CertID++
			id := CertID
			if v.ID > 0 {
				id = v.ID
				if id > CertID {
					CertID = id
				}
			}
			f.ProxyWay[id] = &ProxyWayInfo{
				State: "禁用",
				Note:  v.Note,
				URL:   v.URL,
				ID:    id,
			}
		}
		if f.ProxyDns == "" {
			f.ProxyDns = "localhost"
		}
		if f.ProxyRoles == "" {
			f.ProxyRoles = "//使用换行 或 " + `";" 分号分割` + "\n*.test.com;qqqqq.com;\ndome.com\n*.abc.mmm.cn"
		}
		if len(f.ProxyWay) < 1 {
			CertID++
			id := CertID
			f.ProxyWay[id] = &ProxyWayInfo{
				State: "禁用",
				Note:  "示例:S5代理:无账号密码",
				URL:   "socket://192.168.31.1:4321",
				ID:    id,
			}
			CertID++
			id = CertID
			f.ProxyWay[id] = &ProxyWayInfo{
				State: "禁用",
				Note:  "示例:S5代理:有账号密码",
				URL:   "socket://admin:123456@192.168.31.1:4321",
				ID:    id,
			}
			CertID++
			id = CertID
			f.ProxyWay[id] = &ProxyWayInfo{
				State: "禁用",
				Note:  "示例:http代理:无账号密码",
				URL:   "http://192.168.31.1:4321",
				ID:    id,
			}
			CertID++
			id = CertID
			f.ProxyWay[id] = &ProxyWayInfo{
				State: "禁用",
				Note:  "示例:http代理:有账号密码",
				URL:   "http://admin:123456@192.168.31.1:4321",
				ID:    id,
			}
		}
	}
}
