package shp

import (
	"encoding/binary"
	"fmt"
	"math"
)

// Point is a single pair of X and Y coordinates.
type Point struct {
	X float64
	Y float64

	number uint32
	box    *BoundingBox
}

// DecodePoint decodes a single point shape.
func DecodePoint(buf []byte, num uint32) (*Point, error) {
	return decodePoint(buf, num, nil)
}

// DecodePointP decodes a single point shape with specified precision.
func DecodePointP(buf []byte, num uint32, precision uint) (*Point, error) {
	return decodePoint(buf, num, &precision)
}

// Type is PointType.
func (p *Point) Type() ShapeType {
	return PointType
}

// RecordNumber returns the position in the shape file.
func (p *Point) RecordNumber() uint32 {
	return p.number
}

func (p *Point) String() string {
	return fmt.Sprintf("(%G,%G)", p.X, p.Y)
}

func decodePoint(buf []byte, num uint32, precision *uint) (*Point, error) {
	if len(buf) < 16 {
		return nil, fmt.Errorf("expecting 16 bytes buf only have %d", len(buf))
	}

	float := bytesToFloat64Wrapper(precision)
	return &Point{
		X:      float(buf[0:8]),
		Y:      float(buf[8:16]),
		number: num,
	}, nil
}

func bytesToFloat64(buf []byte) float64 {
	u := binary.LittleEndian.Uint64(buf)
	return math.Float64frombits(u)
}

func bytesToFloat64P(buf []byte, precision uint) float64 {
	f := bytesToFloat64(buf)
	s := math.Pow(10, float64(precision))
	return math.Round(f*s) / s
}

func bytesToFloat64Wrapper(precision *uint) func([]byte) float64 {
	if precision == nil {
		return bytesToFloat64
	}

	return func(buf []byte) float64 {
		return bytesToFloat64P(buf, *precision)
	}
}
