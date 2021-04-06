package shapefile

import (
	"github.com/everystreet/go-geojson/v2"
	"github.com/everystreet/go-shapefile/dbf"
	"github.com/everystreet/go-shapefile/shp"
)

// Record consists of a shape (read from the .shp file) and attributes (from the .dbf file).
type Record struct {
	shp.Shape
	Attributes
}

// Attributes provides access to the dbf record.
type Attributes interface {
	Fields() []dbf.Field
	Field(string) (dbf.Field, bool)
	Deleted() bool
}

// GeoJSONFeature creates a GeoJSON Feature for the Shapefile Record.
func (r Record) GeoJSONFeature(opts ...GeoJSONOption) *geojson.Feature {
	conf := geoJSONConfig{}
	for _, opt := range opts {
		opt(&conf)
	}

	feat := r.Shape.GeoJSONFeature()
	if r.Attributes == nil {
		return feat
	}

	feat.Properties = make(geojson.PropertyList, len(r.Attributes.Fields()))
	for i, f := range r.Attributes.Fields() {
		name := f.Name()
		if newName, ok := conf.oldNewPropNames[name]; ok {
			name = newName
		}

		feat.Properties[i] = geojson.Property{
			Name:  name,
			Value: f.Value(),
		}
	}
	return feat
}

// GeoJSONOption funcs can be passed to Record.GeoJSONFeature().
type GeoJSONOption func(*geoJSONConfig)

// RenameProperties allows shapefile field names to be mapped to user-defined GeoJSON property names.
func RenameProperties(oldNews map[string]string) GeoJSONOption {
	return func(c *geoJSONConfig) {
		c.oldNewPropNames = oldNews
	}
}

type geoJSONConfig struct {
	oldNewPropNames map[string]string
}
