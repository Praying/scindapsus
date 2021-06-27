package main

import (
	"scindapsus/app"
)

func main() {
	ap := app.App{}
	builder := ap.Builder()
	builder.Init()
	builder.Start()
	builder.Run()
}
