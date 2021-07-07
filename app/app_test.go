package app

import (
	"scindapsus/exchanges"
	"testing"
)

func TestApp_Builder(t *testing.T) {
	ap := App{}
	builder := ap.Builder()
	builder.Exchange = exchanges.NewOKExchange()
	builder.Init()
	builder.Start()

	builder.Run()
}

func TestOther(t *testing.T) {

}
