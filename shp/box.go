package shp

import (
	"fmt"
)

// BoundingBox of the shape file.
type BoundingBox struct {
	MinX float64
	MinY float64
	MaxX float64
	MaxY float64
}

// DecodeBoundingBox decodes the bounding box coordinates.
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

func (b *BoundingBox) String() string {
	return fmt.Sprintf("(%f,%f), (%f,%f)", b.MaxX, b.MinY, b.MinX, b.MaxY)
}
