package shp

import (
	"encoding/binary"
	"fmt"
	"math"
)

// Point is a single pair of X and Y coordinates.
type Point struct {
	number uint32
	box    *BoundingBox

	X float64
	Y float64
}

// DecodePoint decodes a single point shape.
func DecodePoint(buf []byte, num uint32) (*Point, error) {
	if len(buf) < 16 {
		return nil, fmt.Errorf("expecting 16 bytes buf only have %d", len(buf))
	}

	return &Point{
		X:      bytesToFloat64(buf[0:8]),
		Y:      bytesToFloat64(buf[8:16]),
		number: num,
	}, nil
}

// RecordNumber returns the position in the shape file.
func (p *Point) RecordNumber() uint32 {
	return p.number
}

func (p *Point) String() string {
	return fmt.Sprintf("(%f,%f)", p.X, p.Y)
}

func bytesToFloat64(b []byte) float64 {
	u := binary.LittleEndian.Uint64(b)
	return math.Float64frombits(u)
}
