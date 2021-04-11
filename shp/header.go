package shp

import (
	"encoding/binary"
	"fmt"
)

// Header represents a shp header.
type Header struct {
	FileLength  uint32
	Version     uint32
	ShapeType   ShapeType
	BoundingBox BoundingBox
}

// DecodeHeader decodes a shp header.
func DecodeHeader(buf []byte, opts ...Option) (Header, error) {
	var conf config
	for _, opt := range opts {
		opt(&conf)
	}

	if len(buf) != 100 {
		return Header{}, fmt.Errorf("have %d bytes, expecting >= 100", len(buf))
	}

	code := binary.BigEndian.Uint32(buf[0:4])
	if code != 0x0000270a {
		return Header{}, fmt.Errorf("bad file code")
	}

	shape := binary.LittleEndian.Uint32(buf[32:36])
	if !validShapeType(shape) {
		return Header{}, fmt.Errorf("invalid shape type %d", shape)
	}

	out := Header{
		// file length is in 16-bit words - but bytes is more useful
		FileLength: binary.BigEndian.Uint32(buf[24:28]) * 2,
		Version:    binary.LittleEndian.Uint32(buf[28:32]),
		ShapeType:  ShapeType(shape),
	}
	var err error
	if conf.precision == nil {
		if out.BoundingBox, err = DecodeBoundingBox(buf[36:]); err != nil {
			return Header{}, err
		}
	} else {
		if out.BoundingBox, err = DecodeBoundingBoxP(buf[36:], *conf.precision); err != nil {
			return Header{}, err
		}
	}

	return out, nil
}

func validShapeType(u uint32) bool {
	switch ShapeType(u) {
	case
		PointType,
		PolylineType,
		PolygonType,
		MultiPointType,
		PointZType,
		PolylineZType,
		PolygonZType,
		MultiPointZType,
		PointMType,
		PolylineMType,
		PolygonMType,
		MultiPointMType,
		MultiPatchType:
		return true
	default:
		return false
	}
}
