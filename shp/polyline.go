package shp

import (
	"encoding/binary"
	"fmt"

	"github.com/pkg/errors"
)

type Polyline struct {
	number uint32

	BoundingBox BoundingBox
	Parts       []Part
}

func DecodePolyline(buf []byte, num uint32) (*Polyline, error) {
	box, err := DecodeBoundingBox(buf[0:])
	if err != nil {
		return nil, err
	}

	const minBytes = 40
	if len(buf) < minBytes {
		return nil, fmt.Errorf("expecting %d bytes but only have %d", minBytes, len(buf))
	}

	numParts := binary.LittleEndian.Uint32(buf[32:36])
	numPoints := binary.LittleEndian.Uint32(buf[36:40])
	numBytes := minBytes + (numParts * 4) + (numPoints * 16)
	if len(buf) < int(numBytes) {
		return nil, fmt.Errorf("expecting %d bytes but only have %d", numBytes, len(buf))
	}

	out := &Polyline{
		BoundingBox: *box,
		Parts:       make([]Part, numParts),
		number:      num,
	}

	parts := make([]uint32, numParts)
	for i := range parts {
		n := minBytes + (i * 4)
		parts[i] = binary.LittleEndian.Uint32(buf[n : n+4])
	}

	pointsOffset := int(minBytes + (numParts * 4))

	for i, start := range parts {
		var end uint32
		if i == len(parts)-1 {
			end = numPoints
		} else {
			end = parts[i+1]
		}

		out.Parts[i] = make(Part, end-start)
		for j := 0; j < len(out.Parts[i]); j++ {
			x := int(start) + j
			p, err := DecodePoint(buf[pointsOffset+(x*16):pointsOffset+(x*16)+16], num)
			if err != nil {
				return nil, errors.Wrap(err, "failed to decode point")
			}
			out.Parts[i][j] = *p
		}
	}

	return out, nil
}

func (p *Polyline) RecordNumber() uint32 {
	return p.number
}

type Polygon Polyline

func DecodePolygon(buf []byte, num uint32) (*Polygon, error) {
	p, err := DecodePolyline(buf, num)
	if err != nil {
		return nil, err
	}
	return (*Polygon)(p), nil
}

func (p *Polygon) RecordNumber() uint32 {
	return p.number
}

type Part []Point
