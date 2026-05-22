//go:build windows
// +build windows

package IsInstallCert

import (
	"crypto/x509"
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

func CheckSunnyNet() (bool, error) {
	// 打开“根证书颁发机构”存储区
	store, err := windows.CertOpenSystemStore(0, syscall.StringToUTF16Ptr("ROOT"))
	if err != nil {
		return false, err
	}
	defer windows.CertCloseStore(store, 0)
	var cert *windows.CertContext
	for {
		cert, _ = windows.CertFindCertificateInStore(store,
			windows.X509_ASN_ENCODING|windows.PKCS_7_ASN_ENCODING,
			0,
			windows.CERT_FIND_ANY,
			nil,
			cert)
		if cert == nil {
			break
		}

		raw := (*[1 << 20]byte)(unsafe.Pointer(cert.EncodedCert))[:cert.Length:cert.Length]
		parsed, err := x509.ParseCertificate(raw)
		if err != nil {
			continue
		}
		if parsed.Subject.CommonName == "SunnyNet" || parsed.Issuer.CommonName == "SunnyNet" {
			return true, nil
		}
	}
	return false, nil
}
