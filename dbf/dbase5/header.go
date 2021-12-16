package dbase5

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Header represents a dBase 5 file header.
type Header struct {
	Fields []*FieldDesc

	recLen  uint16
	numRecs uint32
}

// DecodeHeader parses a dBase 5 file header.
func DecodeHeader(r io.Reader) (*Header, error) {
	// Read first 31 bytes after first byte
	buf := make([]byte, 31)
	if n, err := io.ReadFull(r, buf); err != nil {
		return nil, fmt.Errorf("read %d bytes but expecting %d: %w", n, len(buf), err)
	}

	out := &Header{
		recLen:  binary.LittleEndian.Uint16(buf[9:11]),
		numRecs: binary.LittleEndian.Uint32(buf[3:7]),
	}

	// Read remainder of header
	headerLen := binary.LittleEndian.Uint16(buf[7:9])
	buf = make([]byte, int(headerLen)-len(buf)-1)
	if n, err := io.ReadFull(r, buf); err != nil {
		return nil, fmt.Errorf("read %d bytes but expecting %d: %w", n, len(buf), err)
	}

	if (len(buf)-1)%32 != 0 {
		return nil, fmt.Errorf("invalid header size %d bytes", headerLen)
	}

	if buf[len(buf)-1] != 0x0D {
		return nil, fmt.Errorf("missing field descriptor terminator")
	}

	numFields := (len(buf) - 1) / 32
	out.Fields = make([]*FieldDesc, numFields)
	for i := 0; i < numFields; i++ {
		f, err := DecodeFieldDesc(buf[i*32 : (i*32)+32])
		if err != nil {
			return nil, fmt.Errorf("failed to decode field %d: %w", i, err)
		}
		out.Fields[i] = f
	}

	return out, nil
}

// RecordLen returns the size in bytes of each record in the file.
func (h Header) RecordLen() uint16 {
	return h.recLen
}

// NumRecords returns the number of records in the file.
func (h Header) NumRecords() uint32 {
	return h.numRecs
}

func (h Header) FieldExists(name string) bool {
	for _, field := range h.Fields {
		if field.name == name {
			return true
		}
	}
	return false
}
