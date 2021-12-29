package shp

import "fmt"

// BoundingBox of the shape file.
type BoundingBox struct {
	minX float64
	minY float64
	maxX float64
	maxY float64
}

// MakeBoundingBox at the specified coordinates.
func MakeBoundingBox(minX, minY, maxX, maxY float64) BoundingBox {
	return BoundingBox{
		minX: minX,
		minY: minY,
		maxX: maxX,
		maxY: maxY,
	}
}

// DecodeBoundingBox decodes the bounding box coordinates.
func DecodeBoundingBox(buf []byte) (BoundingBox, error) {
	return decodeBoundingBox(buf, nil)
}

// DecodeBoundingBoxP decodes the bounding box coordinates with a specified precision.
func DecodeBoundingBoxP(buf []byte, precision uint) (BoundingBox, error) {
	return decodeBoundingBox(buf, &precision)
}

func (b BoundingBox) BottomLeft() (x, y float64) {
	return b.minX, b.minY
}

func (b BoundingBox) TopRight() (x, y float64) {
	return b.maxX, b.maxY
}

func (b BoundingBox) String() string {
	return fmt.Sprintf("(%G,%G), (%G,%G)", b.maxX, b.minY, b.minX, b.maxY)
}

func decodeBoundingBox(buf []byte, precision *uint) (BoundingBox, error) {
	if len(buf) < 32 {
		return BoundingBox{}, fmt.Errorf("have %d bytes, expecting >= 32", len(buf))
	}

	float := bytesToFloat64Wrapper(precision)
	return BoundingBox{
		minX: float(buf[0:8]),
		minY: float(buf[8:16]),
		maxX: float(buf[16:24]),
		maxY: float(buf[24:32]),
	}, nil
}
