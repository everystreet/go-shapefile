package shp

import (
	"encoding/binary"
	"fmt"
	"math"
)

// Point is a single pair of X and Y coordinates.
type Point struct {
	number uint32

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

func bytesToFloat64(b []byte) float64 {
	u := binary.LittleEndian.Uint64(b)
	return math.Float64frombits(u)
}
