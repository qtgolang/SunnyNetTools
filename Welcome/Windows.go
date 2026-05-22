//go:build windows
// +build windows

package Welcome

import "C"
import (
	"changeme/Welcome/Windows"
	"golang.org/x/sys/windows"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

func init() {
	if !isRunAsAdmin() {
		runAsAdmin()
	}
	Windows.DisableEdgeTelemetry()
	go Windows.Start()
}

func Stop() {
	Windows.Stop()
}

var shell32 = syscall.NewLazyDLL("shell32.dll")
var procShellExecuteExW = shell32.NewProc("ShellExecuteExW")

type ShellExecuteInfo struct {
	CbSize       uint32
	FMask        uint32
	Hwnd         uintptr
	LpVerb       *uint16
	LpFile       *uint16
	LpParameters *uint16
	LpDirectory  *uint16
	NShow        int32
	HInstApp     uintptr
	LpIDList     unsafe.Pointer
	LpClass      *uint16
	HkeyClass    uintptr
	DwHotKey     uint32
	HIcon        uintptr
	HProcess     windows.Handle
}

func shellExecuteEx(sei *ShellExecuteInfo) error {
	r1, _, err := syscall.SyscallN(
		procShellExecuteExW.Addr(),
		uintptr(unsafe.Pointer(sei)),
	)
	if r1 == 0 {
		return err
	}
	return nil
}
func isRunAsAdmin() bool {
	adminSID, err := windows.CreateWellKnownSid(windows.WinBuiltinAdministratorsSid)
	if err != nil {
		return false
	}
	token := windows.Token(0)
	isMember, err := token.IsMember(adminSID)
	return err == nil && isMember
}
func runAsAdmin() {
	exePath, _ := os.Executable()
	args := strings.Join(os.Args[1:], " ")
	verbPtr, _ := syscall.UTF16PtrFromString("runas")
	exePtr, _ := syscall.UTF16PtrFromString(exePath)
	argsPtr, _ := syscall.UTF16PtrFromString(args)
	sei := &ShellExecuteInfo{
		CbSize:       uint32(unsafe.Sizeof(ShellExecuteInfo{})),
		FMask:        0x00000040,
		Hwnd:         0,
		LpVerb:       verbPtr,
		LpFile:       exePtr,
		LpParameters: argsPtr,
		NShow:        1,
	}
	_ = shellExecuteEx(sei)
	os.Exit(0)
}
