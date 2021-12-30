package shp

import geojson "github.com/everystreet/go-geojson/v3"

// GeoJSONBoundingBox creates a GeoJSON bounding box from the shapefile bounding box.
func (b BoundingBox) GeoJSONBoundingBox() *geojson.BoundingBox {
	return &geojson.BoundingBox{
		BottomLeft: geojson.MakePosition(b.minY, b.minX),
		TopRight:   geojson.MakePosition(b.maxY, b.maxX),
	}
}

// GeoJSONFeature creates a GeoJSON Point from the shapefile Point.
func (p Point) GeoJSONFeature() *geojson.Feature[geojson.Geometry] {
	return &geojson.Feature[geojson.Geometry]{
		Geometry: geojson.NewPoint(p.Point.X, p.Point.Y),
	}
}

// GeoJSONFeature creates a GeoJSON MultiLineString from the shapefile Polyline.
func (p Polyline) GeoJSONFeature() *geojson.Feature[geojson.Geometry] {
	return &geojson.Feature[geojson.Geometry]{
		Geometry: geojson.NewMultiLineString(sliceOfPositionSlices(p.parts)...),
		BBox:     p.box.GeoJSONBoundingBox(),
	}
}

// GeoJSONFeature creates a GeoJSON Polygon from the shapefile Polygon.
func (p Polygon) GeoJSONFeature() *geojson.Feature[geojson.Geometry] {
	return &geojson.Feature[geojson.Geometry]{
		Geometry: geojson.NewPolygon(sliceOfPositionSlices(p.parts)...),
		BBox:     p.box.GeoJSONBoundingBox(),
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
