package shp

type ShapeType uint

const (
	Null            ShapeType = 0
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

type Shape interface {
	RecordNumber() uint32
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
