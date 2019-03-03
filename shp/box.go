package shp

import (
	"fmt"
)

type BoundingBox struct {
	MinX float64
	MinY float64
	MaxX float64
	MaxY float64
}

func DecodeBoundingBox(buf []byte) (*BoundingBox, error) {
	if len(buf) < 32 {
		return nil, fmt.Errorf("have %d bytes, expecting >= 32", len(buf))
	}

	return &BoundingBox{
		MinX: bytesToFloat64(buf[0:8]),
		MinY: bytesToFloat64(buf[8:16]),
		MaxX: bytesToFloat64(buf[16:24]),
		MaxY: bytesToFloat64(buf[24:32]),
	}, nil
}
