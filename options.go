package shapefile

import (
	"github.com/everystreet/go-shapefile/cpg"
	"github.com/everystreet/go-shapefile/dbf"
	"github.com/everystreet/go-shapefile/shp"
)

// Option funcs can be passed to NewScanner().
type Option func(*options)

// PointPrecision sets shp.PointPrecision.
func PointPrecision(p uint) Option {
	return func(o *options) {
		o.shp = append(o.shp, shp.PointPrecision(p))
	}
}

// CharacterEncoding sets dbf.CharacterEncoding.
func CharacterEncoding(enc cpg.CharacterEncoding) Option {
	return func(o *options) {
		o.dbf = append(o.dbf, dbf.CharacterEncoding(enc))
	}
}

// FilterFields sets dbf.FilterFields.
func FilterFields(names ...string) Option {
	return func(o *options) {
		o.dbf = append(o.dbf, dbf.FilterFields(names...))
	}
}

// Options for shp and dbf parsing.
type options struct {
	shp []shp.Option
	dbf []dbf.Option
}
