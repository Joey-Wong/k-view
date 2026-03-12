package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "视频预览",
		Width:  400, // 调整窗口大小以适应紧凑的易语言风格配置界面
		Height: 280, // 调整窗口大小以适应紧凑的易语言风格配置界面
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 212, G: 208, B: 200, A: 255}, // 易语言风格灰色背景 #d4d0c8
		OnStartup:        app.startup,
		DisableResize:    true, // 禁止调整窗口大小
		Windows: &windows.Options{
			DisableWindowIcon: false, // 是否禁用窗口图标
		},
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
