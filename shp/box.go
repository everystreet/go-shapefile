package shp

import "fmt"

// BoundingBox of the shape file.
type BoundingBox struct {
	MinX float64
	MinY float64
	MaxX float64
	MaxY float64
}

// DecodeBoundingBox decodes the bounding box coordinates.
func DecodeBoundingBox(buf []byte) (BoundingBox, error) {
	return decodeBoundingBox(buf, nil)
}

// DecodeBoundingBoxP decodes the bounding box coordinates with a specified precision.
func DecodeBoundingBoxP(buf []byte, precision uint) (BoundingBox, error) {
	return decodeBoundingBox(buf, &precision)
}

func (b BoundingBox) String() string {
	return fmt.Sprintf("(%G,%G), (%G,%G)", b.MaxX, b.MinY, b.MinX, b.MaxY)
}

func decodeBoundingBox(buf []byte, precision *uint) (BoundingBox, error) {
	if len(buf) < 32 {
		return BoundingBox{}, fmt.Errorf("have %d bytes, expecting >= 32", len(buf))
	}

	float := bytesToFloat64Wrapper(precision)
	return BoundingBox{
		MinX: float(buf[0:8]),
		MinY: float(buf[8:16]),
		MaxX: float(buf[16:24]),
		MaxY: float(buf[24:32]),
	}, nil
}
