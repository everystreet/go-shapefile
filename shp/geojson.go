package shp

import "github.com/everystreet/go-geojson/v2"

// GeoJSONFeature creates a GeoJSON Point from a Shapefile Point.
func (p Point) GeoJSONFeature() *geojson.Feature {
	return geojson.NewPoint(p.X, p.Y)
}

// GeoJSONFeature creates a GeoJSON MultiLineString from a Shapefile Polyline.
func (p Polyline) GeoJSONFeature() *geojson.Feature {
	strings := sliceOfPositionSlices(p.Parts)
	return withBox(&p.BoundingBox, geojson.NewMultiLineString(strings...))
}

// GeoJSONFeature creates a GeoJSON Polygon from a Shapefile Polygon.
func (p Polygon) GeoJSONFeature() *geojson.Feature {
	strings := sliceOfPositionSlices(p.Parts)
	return withBox(&p.BoundingBox, geojson.NewPolygon(strings...))
}

func sliceOfPositionSlices(parts []Part) [][]geojson.Position {
	strings := make([][]geojson.Position, len(parts))
	for i, part := range parts {
		strings[i] = make([]geojson.Position, len(part))
		for j, point := range part {
			strings[i][j] = geojson.MakePosition(point.Y, point.X)
		}
	}
	return strings
}

func withBox(b *BoundingBox, f *geojson.Feature) *geojson.Feature {
	return f.WithBoundingBox(
		geojson.MakePosition(b.MinY, b.MinX),
		geojson.MakePosition(b.MaxY, b.MaxX),
	)
}
