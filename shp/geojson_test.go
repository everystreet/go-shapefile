package shp_test

import (
	"testing"

	geojson "github.com/everystreet/go-geojson/v3"
	"github.com/everystreet/go-shapefile/shp"
	"github.com/stretchr/testify/require"
)

func TestPointToGeoJSON(t *testing.T) {
	p := shp.MakePoint(12.34, 56.78)
	require.Equal(t,
		geojson.Feature[geojson.Geometry]{
			Geometry: geojson.NewPoint(12.34, 56.78),
		}, *p.GeoJSONFeature(),
	)
}

func TestPolylineToGeoJSON(t *testing.T) {
	polyline := shp.MakePolyline(
		[]shp.Part{
			{
				shp.MakePoint(12.34, 56.78),
				shp.MakePoint(23.45, 67.89),
			},
		},
		shp.MakeBoundingBox(1, 1, 100, 100),
	)

	require.Equal(t,
		geojson.Feature[geojson.Geometry]{
			Geometry: geojson.NewMultiLineString(
				[]geojson.Position{
					geojson.MakePosition(56.78, 12.34),
					geojson.MakePosition(67.89, 23.45),
				},
			),
			BBox: &geojson.BoundingBox{
				BottomLeft: geojson.MakePosition(1, 1),
				TopRight:   geojson.MakePosition(100, 100),
			},
		}, *polyline.GeoJSONFeature(),
	)
}

func TestPolygonToGeoJSON(t *testing.T) {
	polygon := shp.MakePolygon(
		[]shp.Part{
			{
				shp.MakePoint(12.34, 56.78),
				shp.MakePoint(23.45, 67.89),
			},
		},
		shp.MakeBoundingBox(1, 1, 100, 100),
	)

	require.Equal(t,
		geojson.Feature[geojson.Geometry]{
			Geometry: geojson.NewPolygon(
				[]geojson.Position{
					geojson.MakePosition(56.78, 12.34),
					geojson.MakePosition(67.89, 23.45),
				},
			),
			BBox: &geojson.BoundingBox{
				BottomLeft: geojson.MakePosition(1, 1),
				TopRight:   geojson.MakePosition(100, 100),
			},
		}, *polygon.GeoJSONFeature(),
	)
}
