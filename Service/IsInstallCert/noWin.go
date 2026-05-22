//go:build !windows
// +build !windows

package IsInstallCert

import (
	"bytes"
	"os/exec"
	"strings"
)

func CheckSunnyNet() (bool, error) {
	keychains := []string{
		"/System/Library/Keychains/SystemRootCertificates.keychain",
		"/Library/Keychains/System.keychain",
	}

	for _, keychain := range keychains {
		cmd := exec.Command("security", "find-certificate", "-a", "-c", "SunnyNet", keychain)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			// 继续下一个 keychain
			continue
		}
		if strings.Contains(out.String(), "SunnyNet") {
			return true, nil
		}
	}

	// 检查用户 Keychain (可选)
	userCmd := exec.Command("security", "find-certificate", "-a", "-c", "SunnyNet")
	var userOut bytes.Buffer
	userCmd.Stdout = &userOut
	err := userCmd.Run()
	if err == nil && strings.Contains(userOut.String(), "SunnyNet") {
		return true, nil
	}

	return false, nil
}
