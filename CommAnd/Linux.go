//go:build !windows && !darwin
// +build !windows,!darwin

package CommAnd

import (
	"os"
	"path/filepath"
	"strconv"
)

func GetDesktopPath() (string, error) {
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

	return nil
}
