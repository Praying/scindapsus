package app

import "testing"

func TestApp_Builder(t *testing.T) {
	ap := App{}
	builder := ap.Builder()
	builder.Init()
	builder.Start()

	builder.Run()
}
