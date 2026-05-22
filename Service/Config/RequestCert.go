package Config

import (
	"github.com/qtgolang/SunnyNet/SunnyNet"
	"github.com/qtgolang/SunnyNet/src/Certificate"
	"strings"
)

type CertInfo struct {
	Cert   *Certificate.CertManager `json:"-"`
	DoMain string
	Path   string
	Pass   string
	Note   string
	Role   string
	ID     int
	LoadOk bool
}

func (f *config) initRequestCert() {
	//恢复请求证书管理器
	{
		m := f.RequestCert
		f.RequestCert = make(map[int]*CertInfo)
		for _, v := range m {
			if v == nil {
				continue
			}
			if v.LoadOk {
				Cert := SunnyNet.NewCertManager()
				fs := strings.ToLower(v.Path)
				if strings.HasSuffix(fs, ".p12") || strings.HasSuffix(fs, ".pkcs12") {
					if Cert.LoadP12Certificate(v.Path, v.Pass) {
						CertID++
						id := CertID
						f.RequestCert[id] = &CertInfo{
							Cert:   SunnyNet.NewCertManager(),
							DoMain: v.DoMain,
							Path:   v.Path,
							Pass:   v.Pass,
							Note:   v.Note,
							Role:   v.Role,
							ID:     id,
							LoadOk: false,
						}
					}
				} else {
					if Cert.AddCertPoolPath(v.Path) {
						CertID++
						id := CertID
						f.RequestCert[id] = &CertInfo{
							Cert:   SunnyNet.NewCertManager(),
							DoMain: v.DoMain,
							Path:   v.Path,
							Pass:   v.Pass,
							Note:   v.Note,
							Role:   v.Role,
							ID:     id,
							LoadOk: false,
						}
					}
				}
			}
		}
	}
}
