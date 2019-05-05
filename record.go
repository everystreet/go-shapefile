package shapefile

import (
	"github.com/mercatormaps/go-geojson"
	"github.com/mercatormaps/go-shapefile/dbf"
	"github.com/mercatormaps/go-shapefile/shp"
)

// Record consists of a shape (read from the .shp file) and attributes (from the .dbf file).
type Record struct {
	Shape      shp.Shape
	Attributes Attributes
}

// Attributes provides access to the dbf record.
type Attributes interface {
	Fields() []dbf.Field
	Field(string) (dbf.Field, bool)
	Deleted() bool
}

// GeoJSONFeature creates a GeoJSON Feature for the Shapefile Record.
func (r *Record) GeoJSONFeature() *geojson.Feature {
	feat := r.Shape.GeoJSONFeature()
	if r.Attributes == nil {
		return feat
	}

	feat.Properties = make(geojson.PropertyList, len(r.Attributes.Fields()))
	for i, f := range r.Attributes.Fields() {
		feat.Properties[i] = geojson.Property{
			Name:  f.Name(),
			Value: f.Value(),
		}
	}
	return feat
}
