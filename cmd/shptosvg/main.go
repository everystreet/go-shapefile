package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/everystreet/go-shapefile/cli"
)

// App defines the command line interface.
type App struct {
	Export ExportCmd `kong:"cmd,default=1"`
}

func main() {
	var app App
	ctx := kong.Parse(&app)
	ctx.FatalIfErrorf(ctx.Run())
}

type ExportCmd struct {
	Flags cli.Flags `kong:"embed"`
}

func (c ExportCmd) Run(_ *kong.Context) (err error) {
	scanner, closer, err := c.Flags.OpenAllFields()
	if err != nil {
		return err
	}

	defer func() {
		if closeErr := closer.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	if err := scanner.Scan(); err != nil {
		return err
	}

	for {
		rec := scanner.Record()
		if rec == nil {
			break
		}

		field, ok := rec.Attributes.Field("NAME_EN")
		if !ok {
			break
		}

		fmt.Println(rec.Shape.RecordNumber(), field.Value())
	}

	return scanner.Err()
}
