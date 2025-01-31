package main

import (
	"os"

	"godo/deps"
	App "godoos/app"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

func main() {
	// Create an instance of the app structure
	app := App.NewApp()
	os.Setenv("GODOTOPTYPE", "desktop")
	// Create application with options
	err := wails.Run(&options.App{
		Title:  "GodoOS",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: deps.Frontendassets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.Startup,
		OnShutdown:       app.Shutdown,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
