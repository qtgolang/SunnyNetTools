package Config

import "github.com/qtgolang/SunnyNet/SunnyNet"

type ReplaceHostInfo struct {
	Lod     string
	New     string
	Note    string
	ID      int
	LodInfo SunnyNet.TargetInfo `json:"-"`
	NewInfo SunnyNet.TargetInfo `json:"-"`
}

func (v *ReplaceHostInfo) GetLodPort() bool {
	return v.LodInfo.Port != 0
}
func (v *ReplaceHostInfo) GetNewPort() bool {
	return v.NewInfo.Port != 0
}

const InvalidHost = "testSunny.test111222.com"

func (v *ReplaceHostInfo) IsInvalid() bool {
	if v.LodInfo.Host == InvalidHost || v.NewInfo.Host == InvalidHost {
		return true
	}
	if v.LodInfo.Host == "" || v.NewInfo.Host == "" {
		return true
	}
	return false
}
func (f *config) initReplaceHost() {
	//恢复替换Host列表
	{
		m := f.ReplaceHost
		f.ReplaceHost = make(map[int]*ReplaceHostInfo)
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
			f.ReplaceHost[id] = &ReplaceHostInfo{
				Lod:  v.Lod,
				New:  v.New,
				Note: v.Note,
				ID:   id,
			}
			f.ReplaceHost[id].LodInfo.Parse(v.Lod, 0)
			f.ReplaceHost[id].NewInfo.Parse(v.New, 0)
		}
		if len(f.ReplaceHost) == 0 {
			CertID++
			id := CertID
			f.ReplaceHost[id] = &ReplaceHostInfo{
				Note: "示例,不指定域名端口,端口不变",
				ID:   id,
			}
			f.ReplaceHost[id].LodInfo.Parse("www.test.com", 0)
			f.ReplaceHost[id].NewInfo.Parse("www.dome.com", 0)
			f.ReplaceHost[id].Lod = "www.test.com"
			f.ReplaceHost[id].New = "www.dome.com"
			CertID++
			id = CertID
			f.ReplaceHost[id] = &ReplaceHostInfo{
				Note: "示例,旧Host不指定域名端口,任意端口变成新端口",
				ID:   id,
			}
			f.ReplaceHost[id].LodInfo.Parse("www.Sunny.com", 0)
			f.ReplaceHost[id].NewInfo.Parse("www.dome.com:8443", 0)
			f.ReplaceHost[id].Lod = "www.Sunny.com"
			f.ReplaceHost[id].New = "www.dome.com:8443"
			CertID++
			id = CertID
			f.ReplaceHost[id] = &ReplaceHostInfo{
				Note: "示例,新的HOST不指定域名端口,默认80",
				ID:   id,
			}
			f.ReplaceHost[id].LodInfo.Parse("www.test.com:8443", 0)
			f.ReplaceHost[id].NewInfo.Parse("www.dome.com", 0)
			f.ReplaceHost[id].Lod = "www.test.com:8443"
			f.ReplaceHost[id].New = "www.dome.com"
			CertID++
			id = CertID
			f.ReplaceHost[id] = &ReplaceHostInfo{
				Note: "示例 新旧都指定端口",
				ID:   id,
			}
			f.ReplaceHost[id].LodInfo.Parse("6.6.9.9:8443", 0)
			f.ReplaceHost[id].NewInfo.Parse("9.9.6.6:8443", 0)
			f.ReplaceHost[id].Lod = "6.6.9.9:8443"
			f.ReplaceHost[id].New = "9.9.6.6:8443"
			CertID++
			id = CertID
			f.ReplaceHost[id] = &ReplaceHostInfo{
				Note: "示例 新旧都不指定端口,端口不变",
				ID:   id,
			}
			f.ReplaceHost[id].LodInfo.Parse("6.6.9.1", 0)
			f.ReplaceHost[id].NewInfo.Parse("1.9.6.6", 0)
			f.ReplaceHost[id].Lod = "6.6.9.1"
			f.ReplaceHost[id].New = "1.9.6.6"
		}
	}
}
