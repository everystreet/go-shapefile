package shapefile_test

import (
	"testing"

	"github.com/mercatormaps/go-geojson"
	"github.com/mercatormaps/go-shapefile"
	"github.com/mercatormaps/go-shapefile/dbf"
	"github.com/mercatormaps/go-shapefile/shp"
	"github.com/stretchr/testify/require"
)

func TestRecordToGeoJSON(t *testing.T) {
	rec := shapefile.Record{
		Shape: &shp.Point{X: 0, Y: 0},
		Attributes: &fakeAttrs{
			fields: []dbf.Field{
				&fakeField{"prop1", "value1"},
				&fakeField{"prop2", 2},
			},
		},
	}

	require.Equal(t, geojson.NewPoint(0, 0).WithProperties(
		geojson.Property{Name: "prop1", Value: "value1"},
		geojson.Property{Name: "prop2", Value: 2},
	), rec.GeoJSONFeature())
}

type fakeAttrs struct {
	fields []dbf.Field
}

func (f *fakeAttrs) Fields() []dbf.Field {
	return f.fields
}

func (f *fakeAttrs) Field(name string) (dbf.Field, bool) {
	return nil, false
}

func (f *fakeAttrs) Deleted() bool {
	return false
}

type fakeField struct {
	name  string
	value interface{}
}

func (f *fakeField) Name() string {
	return f.name
}

func (f *fakeField) Value() interface{} {
	return f.value
}
