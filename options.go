package shapefile

import (
	"github.com/everystreet/go-shapefile/dbf"
	"github.com/everystreet/go-shapefile/shp"
	"golang.org/x/text/encoding"
)

// PointPrecision sets shp.PointPrecision.
func PointPrecision(p uint) Option {
	return func(c *config) {
		c.shp = append(c.shp, shp.PointPrecision(p))
	}
}

// CharacterDecoder sets dbf.CharacterDecoder.
func CharacterDecoder(dec *encoding.Decoder) Option {
	return func(c *config) {
		c.dbf = append(c.dbf, dbf.CharacterDecoder(dec))
	}
}

// FilterFields sets dbf.FilterFields.
func FilterFields(names ...string) Option {
	return func(c *config) {
		c.dbf = append(c.dbf, dbf.FilterFields(names...))
	}
}

// Option funcs modify reading of shapefiles.
type Option func(*config)

type config struct {
	shp []shp.Option
	dbf []dbf.Option
}

func defaultConfig() config {
	return config{}
}
