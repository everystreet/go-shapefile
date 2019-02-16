package shp

type ShapeType uint

const (
	Null        ShapeType = 0
	Point       ShapeType = 1
	Polyline    ShapeType = 3
	Polygon     ShapeType = 5
	MultiPoint  ShapeType = 8
	PointZ      ShapeType = 11
	PolylineZ   ShapeType = 13
	PolygonZ    ShapeType = 15
	MultiPointZ ShapeType = 18
	PointM      ShapeType = 21
	PolylineM   ShapeType = 23
	PolygonM    ShapeType = 25
	MultiPointM ShapeType = 28
	MultiPatch  ShapeType = 31
)

func validShapeType(u uint32) bool {
	switch ShapeType(u) {
	case
		Point,
		Polyline,
		Polygon,
		MultiPoint,
		PointZ,
		PolylineZ,
		PolygonZ,
		MultiPointZ,
		PointM,
		PolylineM,
		PolygonM,
		MultiPointM,
		MultiPatch:
		return true
	default:
		return false
	}
}
