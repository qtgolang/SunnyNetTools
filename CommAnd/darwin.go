//go:build darwin
// +build darwin

package CommAnd

import (
	"github.com/lwch/rdesktop/clipboard"
	"net"
	"os"
	"path/filepath"
	"strconv"
)

var UserSelectPath = ""

func GetDesktopPath() (string, error) {
	if UserSelectPath != "" {
		return UserSelectPath, nil
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, "Desktop"), nil
}

func GetPidName(pid int) string {
	if pid < 1 {
		return "代理"
	}
	return strconv.Itoa(pid)
}

func SetIEProxy(Set bool, Port int) bool {

	return false
}

func EnumerateProcesses() map[int]string {
	return nil
}

func InstallCert(certificates []byte) string {
	return ""
}
func GetWayArray() []string {
	var ipArray []string
	interfaces, err := net.Interfaces()
	if err != nil {
		return ipArray
	}
	for _, face := range interfaces {
		adders, err1 := face.Addrs()
		if err1 != nil {
			continue
		}
		for _, addr := range adders {
			ipNet, ok := addr.(*net.IPNet)
			if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
				ipArray = append(ipArray, ipNet.IP.String())
			}
		}
	}
	return ipArray
}
func ClipboardText(text string) error {
	return clipboard.Set(text)
}
