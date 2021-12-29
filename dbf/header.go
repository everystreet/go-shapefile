package dbf

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Header represents a dBase 5 file header.
type Header struct {
	fields  []FieldDesc
	version Version
	recLen  uint16
	numRecs uint32
}

// Version is the dBase version or "level".
type Version uint

// dBase versions.
const (
	DBaseLevel5 Version = 3
	DBaseLevel7 Version = 4
)

// DecodeHeader parses a dBase 5 file header.
func DecodeHeader(r io.Reader) (Header, error) {
	buf := make([]byte, 32)
	if n, err := io.ReadFull(r, buf); err != nil {
		return Header{}, fmt.Errorf("read %d bytes but expecting %d: %w", n, len(buf), err)
	}

	out := Header{
		numRecs: binary.LittleEndian.Uint32(buf[4:8]),
		recLen:  binary.LittleEndian.Uint16(buf[10:12]),
	}

	out.version = Version(((buf[0]>>0)&1)<<0 | ((buf[0]>>1)&1)<<1 | ((buf[0]>>2)&1)<<2)
	if out.version != DBaseLevel5 {
		return Header{}, fmt.Errorf("unsupported bBase version '%d'", out.version)
	}

	// Read remainder of header.
	headerLen := binary.LittleEndian.Uint16(buf[8:10])
	buf = make([]byte, int(headerLen)-len(buf))
	if n, err := io.ReadFull(r, buf); err != nil {
		return Header{}, fmt.Errorf("read %d bytes but expecting %d: %w", n, len(buf), err)
	}

	if (len(buf)-1)%32 != 0 {
		return Header{}, fmt.Errorf("invalid header size %d bytes", headerLen)
	}

	if buf[len(buf)-1] != 0x0D {
		return Header{}, fmt.Errorf("missing field descriptor terminator")
	}

	numFields := (len(buf) - 1) / 32
	out.fields = make([]FieldDesc, numFields)
	for i := 0; i < numFields; i++ {
		f, err := DecodeFieldDesc(buf[i*32 : (i*32)+32])
		if err != nil {
			return Header{}, fmt.Errorf("failed to decode field %d: %w", i, err)
		}
		out.fields[i] = f
	}

	return out, nil
}

// Fields returns the list of field descriptors.
func (h Header) Fields() []FieldDesc {
	return h.fields
}

// Version returns the dBase version.
func (h Header) Version() Version {
	return h.version
}

// RecordLen returns the size in bytes of each record in the file.
func (h Header) RecordLen() uint16 {
	return h.recLen
}

// NumRecords returns the number of records in the file.
func (h Header) NumRecords() uint32 {
	return h.numRecs
}
