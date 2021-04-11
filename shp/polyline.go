package shp

import (
	"encoding/binary"
	"fmt"

	"github.com/golang/geo/r2"
)

// Polyline is an ordered set of verticies that consists of one or more parts, where a part is one or more Point.
type Polyline struct {
	BoundingBox BoundingBox
	Parts       []Part

	number uint32
}

// Part is a sequence of Points.
type Part []Point

// DecodePolyline parses a single polyline shape, but does not validate its complicance with the spec.
func DecodePolyline(buf []byte, num uint32) (Polyline, error) {
	return decodePolyline(buf, num, nil)
}

// DecodePolylineP parses a single polyline shape with the specified precision,
// but does not validate its complicance with the spec.
func DecodePolylineP(buf []byte, num uint32, precision uint) (Polyline, error) {
	return decodePolyline(buf, num, &precision)
}

// Type is PolylineType.
func (p Polyline) Type() ShapeType {
	return PolylineType
}

// RecordNumber returns the position in the shape file.
func (p Polyline) RecordNumber() uint32 {
	return p.number
}

func (p Polyline) points() []r2.Point {
	var out []r2.Point
	for _, part := range p.Parts {
		for _, point := range part {
			out = append(out, point.Point)
		}
	}
	return out
}

// Polygon has the same syntax as a Polyline, but the parts should be unbroken rings.
type Polygon Polyline

// DecodePolygon decodes a single polygon shape, but does not validate its complicance with the spec.
func DecodePolygon(buf []byte, num uint32) (Polygon, error) {
	p, err := DecodePolyline(buf, num)
	if err != nil {
		return Polygon{}, err
	}
	return Polygon(p), nil
}

// DecodePolygonP decodes a single polygon shape with the specified precision,
// but does not validate its complicance with the spec.
func DecodePolygonP(buf []byte, num uint32, precision uint) (Polygon, error) {
	p, err := DecodePolylineP(buf, num, precision)
	if err != nil {
		return Polygon{}, err
	}
	return Polygon(p), nil
}

// Type is PolygonType.
func (p Polygon) Type() ShapeType {
	return PolygonType
}

// RecordNumber returns the position in the shape file.
func (p Polygon) RecordNumber() uint32 {
	return p.number
}

func (p Polygon) points() []r2.Point {
	return Polyline(p).points()
}

func decodePolyline(buf []byte, num uint32, precision *uint) (Polyline, error) {
	var box BoundingBox
	var err error
	if precision == nil {
		if box, err = DecodeBoundingBox(buf[0:]); err != nil {
			return Polyline{}, err
		}
	} else {
		if box, err = DecodeBoundingBoxP(buf[0:], *precision); err != nil {
			return Polyline{}, err
		}
	}

	const minBytes = 40
	if len(buf) < minBytes {
		return Polyline{}, fmt.Errorf("expecting %d bytes but only have %d", minBytes, len(buf))
	}

	numParts := binary.LittleEndian.Uint32(buf[32:36])
	numPoints := binary.LittleEndian.Uint32(buf[36:40])
	numBytes := minBytes + (numParts * 4) + (numPoints * 16)
	if len(buf) < int(numBytes) {
		return Polyline{}, fmt.Errorf("expecting %d bytes but only have %d", numBytes, len(buf))
	}

	out := Polyline{
		BoundingBox: box,
		Parts:       make([]Part, numParts),
		number:      num,
	}

	parts := make([]uint32, numParts)
	for i := range parts {
		n := minBytes + (i * 4)
		parts[i] = binary.LittleEndian.Uint32(buf[n : n+4])
	}

	var point func([]byte, uint32) (Point, error)
	if precision == nil {
		point = DecodePoint
	} else {
		point = func(buf []byte, num uint32) (Point, error) {
			return DecodePointP(buf, num, *precision)
		}
	}

	pointsOffset := int(minBytes + (numParts * 4))
	for i, start := range parts {
		var end uint32
		if i == len(parts)-1 {
			end = numPoints
		} else {
			end = parts[i+1]
		}

		out.Parts[i] = make(Part, end-start)
		for j := 0; j < len(out.Parts[i]); j++ {
			x := int(start) + j
			p, err := point(buf[pointsOffset+(x*16):pointsOffset+(x*16)+16], num)
			if err != nil {
				return Polyline{}, fmt.Errorf("failed to decode point: %w", err)
			}
			p.box = &box
			out.Parts[i][j] = p
		}
	}
	return out, nil
}
