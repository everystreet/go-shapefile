package shp

import (
	"encoding/binary"
	"fmt"
)

// Header represents a shp header.
type Header struct {
	fileLen uint32
	version uint32
	typ     ShapeType
	box     BoundingBox
}

// DecodeHeader decodes a shp header.
func DecodeHeader(buf []byte, precision *uint) (Header, error) {
	if len(buf) != 100 {
		return Header{}, fmt.Errorf("have %d bytes, expecting >= 100", len(buf))
	}

	code := binary.BigEndian.Uint32(buf[0:4])
	if code != 0x0000270a {
		return Header{}, fmt.Errorf("bad file code")
	}

	typ := ShapeType(binary.LittleEndian.Uint32(buf[32:36]))
	if err := validateShapeType(typ); err != nil {
		return Header{}, err
	}

	out := Header{
		// File length is in 16-bit words, but bytes is more useful.
		fileLen: binary.BigEndian.Uint32(buf[24:28]) * 2,
		version: binary.LittleEndian.Uint32(buf[28:32]),
		typ:     typ,
	}

	var err error
	if precision == nil {
		if out.box, err = DecodeBoundingBox(buf[36:]); err != nil {
			return Header{}, err
		}
	} else {
		if out.box, err = DecodeBoundingBoxP(buf[36:], *precision); err != nil {
			return Header{}, err
		}
	}

	return out, nil
}

func (h Header) FileLength() uint32 {
	return h.fileLen
}

func (h Header) Version() uint32 {
	return h.version
}

func (h Header) ShapeType() ShapeType {
	return h.typ
}

func (h Header) BoundingBox() BoundingBox {
	return h.box
}

func validateShapeType(typ ShapeType) error {
	switch typ {
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
		return nil
	default:
		return fmt.Errorf("unrecognized shape type '%d'", typ)
	}
}
