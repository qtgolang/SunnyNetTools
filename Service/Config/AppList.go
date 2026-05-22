package Config

import (
	"context"
	"github.com/wailsapp/wails/v3/pkg/application"
)

var AppList = make(map[string]*AppWindow)

type AppWindow struct {
	*application.WebviewWindow
	context.Context
}

func NewAppWindow(window *application.WebviewWindow) *AppWindow {
	return &AppWindow{window, context.Background()}
}

// ContextWithValue 储存的信息
func (w *AppWindow) ContextWithValue(key, val any) {
	w.Context = context.WithValue(w.Context, key, val)
}

// ContextNewValue 丢弃之前储存的信息
func (w *AppWindow) ContextNewValue(key, val any) {
	w.Context = context.WithValue(context.Background(), key, val)
}
