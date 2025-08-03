package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

// main 应用入口点
func main() {
	// 创建应用实例
	app := NewApp()

	// 创建应用选项
	err := wails.Run(&options.App{
		Title:            "Clash SpeedTest GUI",
		Width:            1024,
		Height:           768,
		MinWidth:         800,
		MinHeight:        600,
		MaxWidth:         1400,
		MaxHeight:        1000,
		DisableResize:    false,
		Fullscreen:       false,
		Frameless:        false,
		StartHidden:      false,
		AlwaysOnTop:      false,
		BackgroundColour: &options.RGBA{R: 45, G: 49, B: 66, A: 1}, // #2D3142
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:        app.OnStartup,
		OnDomReady:       app.OnDomReady,
		OnBeforeClose:    app.OnBeforeClose,
		OnShutdown:       app.OnShutdown,
		WindowStartState: options.Normal,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("启动应用失败:", err.Error())
	}
}