package shapefile

import (
	geojson "github.com/everystreet/go-geojson/v3"
	"github.com/everystreet/go-shapefile/dbf"
	"github.com/everystreet/go-shapefile/shp"
)

// Record consists of a shape (read from the .shp file) and attributes (from the .dbf file).
type Record struct {
	shp.Shape
	*dbf.Record
}

// GeoJSONFeature creates a GeoJSON Feature for the Shapefile Record.
func (r Record) GeoJSONFeature(opts ...GeoJSONOption) *geojson.Feature[geojson.Geometry] {
	conf := geoJSONConfig{}
	for _, opt := range opts {
		opt(&conf)
	}

	feature := r.Shape.GeoJSONFeature()
	if r.Record == nil {
		return feature
	}

	feature.Properties = make(geojson.PropertyList, len(r.Record.Fields))

	var i int
	for _, f := range r.Record.Fields {
		name := f.Name()
		if newName, ok := conf.oldNewPropNames[name]; ok {
			name = newName
		}

		feature.Properties[i] = geojson.Property{
			Name:  name,
			Value: f.Value(),
		}
		i++
	}
	return feature
}

// GeoJSONOption funcs can be passed to Record.GeoJSONFeature().
type GeoJSONOption func(*geoJSONConfig)

// RenameProperties allows shapefile field names to be mapped to user-defined GeoJSON property names.
func RenameProperties(oldToNew map[string]string) GeoJSONOption {
	return func(c *geoJSONConfig) {
		c.oldNewPropNames = oldToNew
	}
}

type geoJSONConfig struct {
	oldNewPropNames map[string]string
}
