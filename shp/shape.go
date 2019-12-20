package shp

import (
	"github.com/everystreet/go-geojson/v2"
)

// ShapeType represents a shape type in the shp file.
type ShapeType uint

// Valid shape types. All shapes in a single shp file must be of the same type.
const (
	// Null shapes are allowed in any shp file, regardless of the type specified in the header.
	Null ShapeType = 0

	PointType       ShapeType = 1
	PolylineType    ShapeType = 3
	PolygonType     ShapeType = 5
	MultiPointType  ShapeType = 8
	PointZType      ShapeType = 11
	PolylineZType   ShapeType = 13
	PolygonZType    ShapeType = 15
	MultiPointZType ShapeType = 18
	PointMType      ShapeType = 21
	PolylineMType   ShapeType = 23
	PolygonMType    ShapeType = 25
	MultiPointMType ShapeType = 28
	MultiPatchType  ShapeType = 31
)

// Shape provides common information for all shapes of any type.
type Shape interface {
	Type() ShapeType
	RecordNumber() uint32
	Validate(*Validator) error
	GeoJSONFeature() *geojson.Feature
}

func validShapeType(u uint32) bool {
	switch ShapeType(u) {
	case
		PointType,
		PolylineType,
		PolygonType,
		MultiPointType,
		PointZType,
		PolylineZType,
		PolygonZType,
		MultiPointZType,
		PointMType,
		PolylineMType,
		PolygonMType,
		MultiPointMType,
		MultiPatchType:
		return true
	default:
		return false
	}
}
