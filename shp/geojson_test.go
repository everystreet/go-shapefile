package shp_test

import (
	"testing"

	geojson "github.com/everystreet/go-geojson/v2"
	"github.com/everystreet/go-shapefile/shp"
	"github.com/stretchr/testify/require"
)

func TestPointToGeoJSON(t *testing.T) {
	p := shp.MakePoint(12.34, 56.78)
	require.Equal(t, geojson.NewPoint(12.34, 56.78), p.GeoJSONFeature())
}

func TestPolylineToGeoJSON(t *testing.T) {
	p := shp.Polyline{
		BoundingBox: shp.BoundingBox{
			MinX: 1,
			MinY: 1,
			MaxX: 100,
			MaxY: 100,
		},
		Parts: []shp.Part{
			{
				shp.MakePoint(12.34, 56.78),
				shp.MakePoint(23.45, 67.89),
			},
		},
	}

	require.Equal(t,
		geojson.NewMultiLineString(
			[]geojson.Position{
				geojson.MakePosition(56.78, 12.34),
				geojson.MakePosition(67.89, 23.45),
			}).
			WithBoundingBox(
				geojson.MakePosition(1, 1),
				geojson.MakePosition(100, 100),
			),
		p.GeoJSONFeature())
}

func TestPolygonToGeoJSON(t *testing.T) {
	p := shp.Polygon{
		BoundingBox: shp.BoundingBox{
			MinX: 1,
			MinY: 1,
			MaxX: 100,
			MaxY: 100,
		},
		Parts: []shp.Part{
			{
				shp.MakePoint(12.34, 56.78),
				shp.MakePoint(23.45, 67.89),
			},
		},
	}

	require.Equal(t,
		geojson.NewPolygon(
			[]geojson.Position{
				geojson.MakePosition(56.78, 12.34),
				geojson.MakePosition(67.89, 23.45),
			}).
			WithBoundingBox(
				geojson.MakePosition(1, 1),
				geojson.MakePosition(100, 100),
			),
		p.GeoJSONFeature())
}
