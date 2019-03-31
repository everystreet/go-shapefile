package shapefile

import (
	"github.com/mercatormaps/go-shapefile/cpg"
	"github.com/mercatormaps/go-shapefile/dbf"
	"github.com/mercatormaps/go-shapefile/shp"
)

// Option funcs can be passed to NewScanner().
type Option func(*Options)

// PointPrecision sets shp.PointPrecision.
func PointPrecision(p uint) Option {
	return func(o *Options) {
		o.shp = append(o.shp, shp.PointPrecision(p))
	}
}

// CharacterEncoding sets dbf.CharacterEncoding.
func CharacterEncoding(enc cpg.CharacterEncoding) Option {
	return func(o *Options) {
		o.dbf = append(o.dbf, dbf.CharacterEncoding(enc))
	}
}

// FilterFields sets dbf.FilterFields.
func FilterFields(names ...string) Option {
	return func(o *Options) {
		o.dbf = append(o.dbf, dbf.FilterFields(names...))
	}
}

// Options for shp and dbf parsing.
type Options struct {
	shp []shp.Option
	dbf []dbf.Option
}
