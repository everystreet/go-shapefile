package shp

import geojson "github.com/everystreet/go-geojson/v3"

// GeoJSONFeature creates a GeoJSON Point from a Shapefile Point.
func (p Point) GeoJSONFeature() *geojson.Feature[geojson.Geometry] {
	return &geojson.Feature[geojson.Geometry]{
		Geometry: geojson.NewPoint(p.Point.X, p.Point.Y),
	}
}

// GeoJSONFeature creates a GeoJSON MultiLineString from a Shapefile Polyline.
func (p Polyline) GeoJSONFeature() *geojson.Feature[geojson.Geometry] {
	return &geojson.Feature[geojson.Geometry]{
		Geometry: geojson.NewMultiLineString(sliceOfPositionSlices(p.Parts)...),
		BBox: &geojson.BoundingBox{
			BottomLeft: geojson.MakePosition(p.BoundingBox.MinY, p.BoundingBox.MinX),
			TopRight:   geojson.MakePosition(p.BoundingBox.MaxY, p.BoundingBox.MaxX),
		},
	}
}

// GeoJSONFeature creates a GeoJSON Polygon from a Shapefile Polygon.
func (p Polygon) GeoJSONFeature() *geojson.Feature[geojson.Geometry] {
	return &geojson.Feature[geojson.Geometry]{
		Geometry: geojson.NewPolygon(sliceOfPositionSlices(p.Parts)...),
		BBox: &geojson.BoundingBox{
			BottomLeft: geojson.MakePosition(p.BoundingBox.MinY, p.BoundingBox.MinX),
			TopRight:   geojson.MakePosition(p.BoundingBox.MaxY, p.BoundingBox.MaxX),
		},
	}
}

func sliceOfPositionSlices(parts []Part) [][]geojson.Position {
	strings := make([][]geojson.Position, len(parts))
	for i, part := range parts {
		strings[i] = make([]geojson.Position, len(part))
		for j, point := range part {
			strings[i][j] = geojson.MakePosition(point.Point.Y, point.Point.X)
		}
	}
	return strings
}
