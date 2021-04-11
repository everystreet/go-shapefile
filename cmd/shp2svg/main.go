package main

import (
	"github.com/alecthomas/kong"
	"github.com/everystreet/go-shapefile/cli"
)

// App defines the command line interface.
type App struct {
	Export cli.ExportSVGCmd `kong:"cmd,default=1"`
}

func main() {
	var app App
	ctx := kong.Parse(&app)
	ctx.FatalIfErrorf(ctx.Run())
}
