package Service

import (
	"changeme/Service/Config"
	"changeme/Service/HookKeys"
	"changeme/Service/clipboard"
	"changeme/Service/mcp"
	"changeme/Service/mcpbridge"
	"embed"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

var app *application.App

func getWebviewBrowserPath() string {
	if runtime.GOOS == "windows" {
		// 获取当前可执行文件的路径
		executablePath, _ := os.Executable()
		// 获取可执行文件所在目录，并拼接 WebView2 路径
		executableDir := filepath.Join(filepath.Dir(executablePath), "WebView2")
		webView2ExecutablePath := filepath.Join(executableDir, "msedgewebview2.exe")
		// 检查 WebView2 可执行文件是否存在
		if _, err := os.Stat(webView2ExecutablePath); err == nil {
			// 文件存在，设置 WebView2 浏览器路径
			fmt.Println("当前使用自定义 WebView2 ")
			return executableDir
		}
	}
	//fmt.Println("当前使用系统中的 WebView2")
	return ""
}
func hideWindow(window *Config.AppWindow) {
	if runtime.GOOS != "windows" {
		var wg sync.WaitGroup
		var ok bool
		var l sync.Mutex
		wg.Add(1)
		window.OnWindowEvent(events.Common.WindowRuntimeReady, func(event *application.WindowEvent) {
			l.Lock()
			defer l.Unlock()
			if ok {
				return
			}
			ok = true
			window.Hide()
			wg.Done()
		})
		wg.Wait()
	} else {
		window.Hide()
	}
}

var dragDropWindowLock sync.Mutex

func registerWindowFilesDropped(window string) {
	Config.AppList[window].OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
		files := event.Context().DroppedFiles()
		dragDropWindowLock.Lock()
		defer dragDropWindowLock.Unlock()
		Config.AppList["Main"].EmitEvent("DropFilesEvent", window, files)
	})

}
func CreateMainWindow(assets embed.FS) *application.App {
	Server := NewAppServer()
	SetMCPServer(Server)
	mcp.InitBridge(MCPBridgeInvoke)
	mcpbridge.EmitMCPBridgeChanged = func() {
		if Config.AppList["Main"] != nil {
			Config.AppList["Main"].EmitEvent("mcpBridgeChanged", mcp.StatusJSON())
		}
	}
	CloseWindow := func() {
		_ = mcp.Disable()
		Config.Config.Save()
		Server.app.Close()
		Server.CancelIEProxy()
	}
	app = application.New(application.Options{
		Name:        "SunnyNetTools",
		Description: "SunnyNetTools-网络抓包工具",
		Services: []application.Service{
			application.NewService(Server),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		Logger:  slog.New(slog.NewTextHandler(io.Discard, nil)),
		Windows: application.WindowsOptions{WebviewBrowserPath: getWebviewBrowserPath()},
	})

	Config.AppList["Main"] = Config.NewAppWindow(app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "SunnyNetTools 网络抓包工具",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
		Width:            1470,
		Height:           780,
		MinWidth:         1470,
		MinHeight:        780,
		EnableFileDrop:   true,
		Frameless:        true,
		DevToolsEnabled:  true,
	}))
	go func() {
		time.Sleep(400 * time.Millisecond)
		mcpbridge.EmitMCPBridgeChanged()
	}()
	go func() {
		time.Sleep(1 * time.Second)
		Config.AppList["Main"].OpenDevTools()
	}()
	Config.AppList["Main"].RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
		CloseWindow()
		Config.AppList["Main"].Hide()
		go func() {
			time.Sleep(time.Second)
			os.Exit(0)
		}()
	})
	clipboard.AppWindow = Config.AppList["Main"].WebviewWindow
	registerWindowFilesDropped("Main")
	if runtime.GOOS == "windows" {
		Config.AppList["Main"].Hide()
		HookKeys.RegisterKeys(Config.Config.Keys, Server.CallKeys)
		go func() {
			time.Sleep(time.Second * 3)
			CreateCertWindow()
			CreateReplaceWindow()
			CreateThemeDesignWindow()
			CreateThemeWindow()
			CreateOtherWindow()
			CreateDebugWindow()
		}()
	}

	return app
}
func CreateCertWindow() {
	Config.AppList["Cert"] = Config.NewAppWindow(app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "请求证书设置",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/Cert.html",
		Width:            1000,
		Height:           400,
		MinWidth:         1000,
		MinHeight:        400,
		EnableFileDrop:   true,
		Frameless:        true,
	}))
	Config.AppList["Cert"].RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
		Config.AppList["Cert"].Hide()
		event.Cancel()
	})
	registerWindowFilesDropped("Cert")
	hideWindow(Config.AppList["Cert"])
}
func CreateReplaceWindow() {
	Config.AppList["ReplaceBody"] = Config.NewAppWindow(app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "请求拦截/数据替换设置",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/ReplaceBody.html",
		Width:            1000,
		Height:           600,
		MinWidth:         1000,
		MinHeight:        600,
		EnableFileDrop:   true,
		Frameless:        true,
	}))
	Config.AppList["ReplaceBody"].RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
		Config.AppList["ReplaceBody"].Hide()
		event.Cancel()
	})
	registerWindowFilesDropped("ReplaceBody")
	hideWindow(Config.AppList["ReplaceBody"])
}
func CreateThemeDesignWindow() {
	Config.AppList["主题设计"] = Config.NewAppWindow(app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "自定义主题配置",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/ThemeDesign.html",
		Width:            1400,
		Height:           660,
		MinWidth:         1400,
		MinHeight:        660,
		EnableFileDrop:   true,
		Frameless:        true,
	}))
	Config.AppList["主题设计"].RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
		Config.AppList["主题设计"].Hide()
		event.Cancel()
	})
	registerWindowFilesDropped("主题设计")
	hideWindow(Config.AppList["主题设计"])
}
func CreateThemeWindow() {
	Config.AppList["主题调色"] = Config.NewAppWindow(app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "主题配置",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/Theme.html",
		Width:            570,
		MinWidth:         570,
		Height:           710,
		MinHeight:        710,
		EnableFileDrop:   true,
		Frameless:        true,
	}))
	Config.AppList["主题调色"].RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
		Config.AppList["主题调色"].Hide()
		event.Cancel()
	})
	registerWindowFilesDropped("主题调色")
	hideWindow(Config.AppList["主题调色"])
}
func CreateOtherWindow() {
	Config.AppList["其他窗口"] = Config.NewAppWindow(app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "其他窗口",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/Other.html",
		Width:            1000,
		MinWidth:         1000,
		Height:           600,
		MinHeight:        600,
		EnableFileDrop:   true,
		Frameless:        true,
	}))
	Config.AppList["其他窗口"].RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
		Config.AppList["其他窗口"].Hide()
		event.Cancel()
	})
	registerWindowFilesDropped("其他窗口")
	hideWindow(Config.AppList["其他窗口"])
}
func CreateDebugWindow() {
	Config.AppList["调试工具"] = Config.NewAppWindow(app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "调试工具",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/debugTools.html",
		Width:            1300,
		MinWidth:         1300,
		Height:           700,
		MinHeight:        700,
		EnableFileDrop:   true,
		Frameless:        true,
	}))
	Config.AppList["调试工具"].RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
		Config.AppList["调试工具"].Hide()
		event.Cancel()
	})
	registerWindowFilesDropped("调试工具")
	hideWindow(Config.AppList["调试工具"])
}
