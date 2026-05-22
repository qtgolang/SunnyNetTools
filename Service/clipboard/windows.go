//go:build !darwin
// +build !darwin

package clipboard

import (
	"github.com/atotto/clipboard"
	"github.com/wailsapp/wails/v3/pkg/application"
)

var AppWindow *application.WebviewWindow

func ClipboardWriteAll(value string) error {
	return clipboard.WriteAll(value)
}
func ClipboardReadAll() string {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	t, e := clipboard.ReadAll()
	if e == nil {
		return t
	}
	return e.Error()
}
