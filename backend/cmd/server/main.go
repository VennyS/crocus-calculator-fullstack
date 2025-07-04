package main

import "crocus-calculator/internal/setting"

func main() {
	app := setting.App{}
	app.LoadConfig()
	r := app.MountRouter()

	app.RunServer(r)
}
