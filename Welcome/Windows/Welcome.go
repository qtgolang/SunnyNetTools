package Windows

/*
#cgo CXXFLAGS: -std=c++17
#cgo LDFLAGS: -Wl,-Bstatic -lstdc++ -Wl,-Bdynamic -lgdi32 -lgdiplus
#cgo LDFLAGS: -static-libgcc
#include "player.hpp"
#include <stdlib.h>
*/
import "C"
import (
	"bytes"
	_ "embed"
	"os"
	"path/filepath"
	"time"
	"unsafe"

	"github.com/qtgolang/SunnyNet/src/Resource"
)

//go:embed welcome.gif
var embeddedGif []byte // 内嵌 GIF 动画

var initialized bool    // 是否已经初始化
var iconTempPath string // 图标临时文件路径
var gifTempPath string  // GIF 临时文件路径
var startTime time.Time

func init() {
	if initialized {
		return
	}
	//不知为何,从字节数组 载入 icon 图标 失败,只能从文件载入 那就干脆从文件载入算了

	// 设置图标和 GIF 的临时路径
	iconTempPath = filepath.Join(os.TempDir(), "_SunnyNetIcon.tmp")
	gifTempPath = filepath.Join(os.TempDir(), "_SunnyNetGif.tmp")

	// 写入图标文件（如果内容有变化）
	existingIcon, _ := os.ReadFile(iconTempPath)
	if !bytes.Equal(existingIcon, Resource.Icon) {
		_ = os.Remove(iconTempPath)
		_ = os.WriteFile(iconTempPath, Resource.Icon, 0600)
	}
	// 写入 GIF 文件（如果内容有变化）
	existingGif, _ := os.ReadFile(gifTempPath)
	if !bytes.Equal(existingGif, embeddedGif) {
		_ = os.Remove(gifTempPath)
		_ = os.WriteFile(gifTempPath, embeddedGif, 0600)
	}
	// 创建 C 字符串路径（用于跨语言调用）
	cIconPath := C.CString(iconTempPath)
	cGifPath := C.CString(gifTempPath)
	C.SetImgPath(cGifPath, cIconPath)
	C.free(unsafe.Pointer(cGifPath))
	C.free(unsafe.Pointer(cIconPath))
}
func Start() {
	if initialized {
		return
	}
	initialized = true
	startTime = time.Now()
	go func() {
		C.StartGIFWindow()
		initialized = false
	}()
}

// Stop 停止欢迎窗口动画
func Stop() {
	if !initialized {
		return
	}
	C.StopGIFWindow() // 调用 C++ 停止函数
}
