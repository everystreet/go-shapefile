package shapefile_test

import (
	"testing"

	geojson "github.com/everystreet/go-geojson/v2"
	shapefile "github.com/everystreet/go-shapefile"
	"github.com/everystreet/go-shapefile/dbf"
	"github.com/everystreet/go-shapefile/shp"
	"github.com/stretchr/testify/require"
)

func TestRecordToGeoJSON(t *testing.T) {
	rec := shapefile.Record{
		Shape: shp.MakePoint(0, 0),
		Record: &dbf.Record{
			Fields: map[string]dbf.Field{
				"prop1": Field{"prop1", "value1"},
				"prop2": Field{"prop2", 2},
				"prop3": Field{"prop3", "value3"},
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

type Field struct {
	name  string
	value any
}

func (f Field) Name() string {
	return f.name
}

func (f Field) Value() any {
	return f.value
}

func (f Field) Equal(_ string) bool {
	return false
}
