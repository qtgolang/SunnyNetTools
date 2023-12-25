package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS
var app = NewApp()

func main() {
	// Create an instance of the app structure
	Size := GlobalConfig.Size
	// Create application with options
	err := wails.Run(&options.App{
		Frameless: true,
		Title:     "SunnyNet",
		Width:     Size.Width,
		Height:    Size.Height,
		MinWidth:  823,
		MinHeight: 388,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
