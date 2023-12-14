//go:build windows
// +build windows

package CommAnd

import (
	"bufio"
	"bytes"
	"github.com/Trisia/gosysproxy"
	gops "github.com/mitchellh/go-ps"
	"github.com/qtgolang/SunnyNet/public"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
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
	process, err := gops.FindProcess(pid)
	if err != nil {
		return strconv.Itoa(pid)
	}
	if process == nil {
		return strconv.Itoa(pid)
	}
	return strconv.Itoa(pid) + ":" + process.Executable()
}

func SetIEProxy(Set bool, Port int) bool {

	if !Set {
		_ = gosysproxy.Off()
		return true
	}
	ies := "127.0.0.1:" + strconv.Itoa(Port)
	_ = gosysproxy.SetGlobalProxy("http="+ies+";https="+ies, "")
	return true
}

func EnumerateProcesses() map[int]string {
	res := make(map[int]string)
	processes, err := gops.Processes()
	if err != nil {
		return res
	}

	// 遍历每个进程并输出 PID 和进程名
	for _, process := range processes {
		res[process.Pid()] = process.Executable()
	}
	return res
}

// InstallCert 安装证书 将证书安装到Windows系统内
func InstallCert(certificates []byte) string {
	path, err := os.Getwd()
	if err != nil {
		return err.Error()
	}
	err = public.WriteBytesToFile(certificates, path+"\\ca.crt")
	if err != nil {
		return err.Error()
	}
	var args []string
	args = append(args, "-addstore")
	args = append(args, "root")
	args = append(args, path+"\\ca.crt")
	defer func() { _ = public.RemoveFile(path + "\\ca.crt") }()
	cmd := exec.Command("certutil", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err.Error()
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	_ = cmd.Start()
	var Buff bytes.Buffer
	reader := bufio.NewReader(stdout)
	for {
		line, err2 := reader.ReadBytes('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		Buff.Write(line)
	}
	utf8Bytes, err := simplifiedchinese.GBK.NewDecoder().Bytes(Buff.Bytes())
	if err == nil {
		return string(utf8Bytes)
	}
	return Buff.String()
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
