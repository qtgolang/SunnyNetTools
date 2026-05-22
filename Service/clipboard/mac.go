//go:build darwin
// +build darwin

package clipboard

import (
	"github.com/wailsapp/wails/v3/pkg/application"
	"golang.design/x/clipboard"
)

var AppWindow *application.WebviewWindow

func init() {
	clipboard.Init()
}
func ClipboardWriteAll(value string) error {
	clipboard.Write(clipboard.FmtText, []byte(value))
	return nil
}
func ClipboardReadAll() string {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	t := clipboard.Read(clipboard.FmtText)
	return string(t)
}
