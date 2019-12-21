package shapefile_test

import (
	"testing"

	"github.com/everystreet/go-geojson/v2"
	"github.com/everystreet/go-shapefile"
	"github.com/everystreet/go-shapefile/dbf"
	"github.com/everystreet/go-shapefile/shp"
	"github.com/stretchr/testify/require"
)

func TestRecordToGeoJSON(t *testing.T) {
	rec := shapefile.Record{
		Shape: shp.MakePoint(0, 0),
		Attributes: &fakeAttrs{
			fields: []dbf.Field{
				&fakeField{"prop1", "value1"},
				&fakeField{"prop2", 2},
				&fakeField{"prop3", "value3"},
			},
		},
	}

	t.Run("simple", func(t *testing.T) {
		require.Equal(t, geojson.NewPoint(0, 0).WithProperties(
			geojson.Property{Name: "prop1", Value: "value1"},
			geojson.Property{Name: "prop2", Value: 2},
			geojson.Property{Name: "prop3", Value: "value3"},
		), rec.GeoJSONFeature())
	})

	t.Run("renamed properties", func(t *testing.T) {
		require.Equal(t, geojson.NewPoint(0, 0).WithProperties(
			geojson.Property{Name: "new-prop-name", Value: "value1"},
			geojson.Property{Name: "new-prop-name", Value: 2},
			geojson.Property{Name: "prop3", Value: "value3"},
		), rec.GeoJSONFeature(
			shapefile.RenameProperties(map[string]string{
				"prop1": "new-prop-name",
				"prop2": "new-prop-name",
			}),
		))
	})
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
