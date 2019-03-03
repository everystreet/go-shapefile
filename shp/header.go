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

func DecodeHeader(b []byte) (*Header, error) {
	if len(b) != 100 {
		return nil, fmt.Errorf("have %d bytes, expecting >= 100", len(b))
	}

	code := binary.BigEndian.Uint32(b[0:4])
	if code != 0x0000270a {
		return nil, fmt.Errorf("bad file code")
	}

	shape := binary.LittleEndian.Uint32(b[32:36])
	if !validShapeType(shape) {
		return nil, fmt.Errorf("invalid shape type %d", shape)
	}

	box, err := DecodeBoundingBox(b[36:])
	if err != nil {
		return nil, err
	}

	return &Header{
		// file length is in 16-bit words - but bytes is more useful
		FileLength:  binary.BigEndian.Uint32(b[24:28]) * 2,
		Version:     binary.LittleEndian.Uint32(b[28:32]),
		ShapeType:   ShapeType(shape),
		BoundingBox: *box,
	}, nil
}
