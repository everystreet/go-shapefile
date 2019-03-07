package shp

import (
	"encoding/binary"
	"fmt"
)

type Header struct {
	FileLength  uint32
	Version     uint32
	ShapeType   ShapeType
	BoundingBox BoundingBox
}

func DecodeHeader(buf []byte) (*Header, error) {
	if len(buf) != 100 {
		return nil, fmt.Errorf("have %d bytes, expecting >= 100", len(buf))
	}

	code := binary.BigEndian.Uint32(buf[0:4])
	if code != 0x0000270a {
		return nil, fmt.Errorf("bad file code")
	}

	shape := binary.LittleEndian.Uint32(buf[32:36])
	if !validShapeType(shape) {
		return nil, fmt.Errorf("invalid shape type %d", shape)
	}

	box, err := DecodeBoundingBox(buf[36:])
	if err != nil {
		return nil, err
	}

	return &Header{
		// file length is in 16-bit words - but bytes is more useful
		FileLength:  binary.BigEndian.Uint32(buf[24:28]) * 2,
		Version:     binary.LittleEndian.Uint32(buf[28:32]),
		ShapeType:   ShapeType(shape),
		BoundingBox: *box,
	}, nil
}
