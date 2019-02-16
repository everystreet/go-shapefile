package shp

import (
	"encoding/binary"
	"fmt"
	"math"
)

type Header struct {
	FileLength uint32
	Version    uint32
	ShapeType  ShapeType

	MinX float64
	MinY float64
	MaxX float64
	MaxY float64
}

func DecodeHeader(b []byte) (*Header, error) {
	if len(b) != 100 {
		return nil, fmt.Errorf("have %d bytes, expecting 100", len(b))
	}

	code := binary.BigEndian.Uint32(b[0:4])
	if code != 0x0000270a {
		return nil, fmt.Errorf("bad file code")
	}

	shape := binary.LittleEndian.Uint32(b[32:36])
	if !validShapeType(shape) {
		return nil, fmt.Errorf("invalid shape type %d", shape)
	}

	return &Header{
		// file length is in 16-bit words - but bytes is more useful
		FileLength: binary.BigEndian.Uint32(b[24:28]) * 2,
		Version:    binary.LittleEndian.Uint32(b[28:32]),
		ShapeType:  ShapeType(shape),
		MinX:       bytesToFloat64(b[36:44]),
		MinY:       bytesToFloat64(b[44:52]),
		MaxX:       bytesToFloat64(b[52:60]),
		MaxY:       bytesToFloat64(b[60:68]),
	}, nil
}

func bytesToFloat64(b []byte) float64 {
	u := binary.LittleEndian.Uint64(b)
	return math.Float64frombits(u)
}
