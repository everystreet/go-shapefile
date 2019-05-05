package shp_test

import (
	"testing"

	"github.com/mercatormaps/go-geojson"
	"github.com/mercatormaps/go-shapefile/shp"
	"github.com/stretchr/testify/require"
)

func TestPointToGeoJSON(t *testing.T) {
	p := shp.Point{X: 12.34, Y: 56.78}
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
				shp.Point{X: 12.34, Y: 56.78},
				shp.Point{X: 23.45, Y: 67.89},
			},
		},
	}

	require.Equal(t,
		geojson.NewMultiLineString(
			[]geojson.Position{
				geojson.NewPosition(12.34, 56.78),
				geojson.NewPosition(23.45, 67.89),
			}).
			WithBoundingBox(
				geojson.NewPosition(1, 1),
				geojson.NewPosition(100, 100),
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
				shp.Point{X: 12.34, Y: 56.78},
				shp.Point{X: 23.45, Y: 67.89},
			},
		},
	}

	require.Equal(t,
		geojson.NewPolygon(
			[]geojson.Position{
				geojson.NewPosition(12.34, 56.78),
				geojson.NewPosition(23.45, 67.89),
			}).
			WithBoundingBox(
				geojson.NewPosition(1, 1),
				geojson.NewPosition(100, 100),
			),
		p.GeoJSONFeature())
}
