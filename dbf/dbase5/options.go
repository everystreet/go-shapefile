package dbase5

import "github.com/mercatormaps/go-shapefile/cpg"

type Option func(Config)

type Config interface {
	CharacterEncoding() cpg.CharacterEncoding
	FilteredFields() []string
}
