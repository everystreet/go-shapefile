package dbf

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Version is the dBase version or "level".
type Version uint8

// dBase versions.
const (
	DBaseLevel5 Version = 3
	DBaseLevel7 Version = 4
)

// Header represents a dBase 5 file header.
type Header struct {
	fields     map[string]FieldDesc
	fieldOrder []string
	version    Version
	headerLen  uint16
	recLen     uint16
	numRecs    uint32
}

func MakeHeader(fields []FieldDesc, version Version, recordLen uint16, numRecords uint32) (Header, error) {
	out := Header{
		fields:     make(map[string]FieldDesc, len(fields)),
		fieldOrder: make([]string, len(fields)),
		version:    version,
		headerLen:  uint16(32+len(fields)*32) + 1,
		recLen:     recordLen,
		numRecs:    numRecords,
	}

	for i, f := range fields {
		if _, ok := out.fields[f.name]; ok {
			return Header{}, fmt.Errorf("duplicate field with name '%s'", f.name)
		}

		out.fields[f.name] = f
		out.fieldOrder[i] = f.name
	}

	return out, nil
}

// DecodeHeader parses a dBase 5 file header.
func DecodeHeader(r io.Reader) (Header, error) {
	buf := make([]byte, 32)
	if n, err := io.ReadFull(r, buf); err != nil {
		return Header{}, fmt.Errorf("read %d bytes but expecting %d: %w", n, len(buf), err)
	}

	out := Header{
		numRecs:   binary.LittleEndian.Uint32(buf[4:8]),
		headerLen: binary.LittleEndian.Uint16(buf[8:10]),
		recLen:    binary.LittleEndian.Uint16(buf[10:12]),
	}

	out.version = Version(((buf[0]>>0)&1)<<0 | ((buf[0]>>1)&1)<<1 | ((buf[0]>>2)&1)<<2)
	if out.version != DBaseLevel5 {
		return Header{}, fmt.Errorf("unsupported bBase version '%d'", out.version)
	}

	if (out.headerLen-1)%32 != 0 {
		return Header{}, fmt.Errorf("invalid header size %d bytes", out.headerLen)
	}

	numFields := (out.headerLen - 32 - 1) / 32
	out.fields = make(map[string]FieldDesc, numFields)
	out.fieldOrder = make([]string, numFields)

	for i := 0; i < int(numFields); i++ {
		buf := make([]byte, 32)
		if n, err := io.ReadFull(r, buf); err != nil {
			return Header{}, fmt.Errorf("read %d bytes but expecting %d: %w", n, len(buf), err)
		}

		f, err := DecodeFieldDesc(buf)
		if err != nil {
			return Header{}, fmt.Errorf("failed to decode field %d: %w", i, err)
		}

		if _, ok := out.fields[f.name]; ok {
			return Header{}, fmt.Errorf("duplicate field with name '%s'", f.name)
		}

		out.fields[f.name] = f
		out.fieldOrder[i] = f.name
	}

	buf = make([]byte, 1)
	if n, err := io.ReadFull(r, buf); err != nil {
		return Header{}, fmt.Errorf("read %d bytes but expecting %d: %w", n, len(buf), err)
	}

	if buf[0] != 0x0D {
		return Header{}, fmt.Errorf("missing field descriptor terminator")
	}

	return out, nil
}

// Fields returns the list of field descriptors.
func (h Header) Fields() []FieldDesc {
	out := make([]FieldDesc, len(h.fields))
	for i, name := range h.fieldOrder {
		out[i] = h.fields[name]
	}
	return out
}

// Version returns the dBase version.
func (h Header) Version() Version {
	return h.version
}

// Length of the header.
func (h Header) Len() uint16 {
	return h.headerLen
}

// RecordLen returns the size in bytes of each record in the file.
func (h Header) RecordLen() uint16 {
	return h.recLen
}

// NumRecords returns the number of records in the file.
func (h Header) NumRecords() uint32 {
	return h.numRecs
}

func (h Header) Encode(w io.Writer) error {
	buf := make([]byte, 32)

	if h.version != DBaseLevel5 {
		return fmt.Errorf("unsupported bBase version '%d'", h.version)
	}
	buf[0] = ((byte(h.version)>>0)&1)<<0 | ((byte(h.version)>>1)&1)<<1 | ((byte(h.version)>>2)&1)<<2

	binary.LittleEndian.PutUint32(buf[4:], h.numRecs)
	binary.LittleEndian.PutUint16(buf[8:], h.headerLen)
	binary.LittleEndian.PutUint16(buf[10:], h.recLen)

	if _, err := w.Write(buf); err != nil {
		return err
	}

	for _, f := range h.fields {
		buf, err := f.Encode()
		if err != nil {
			return err
		}

		if _, err := w.Write(buf); err != nil {
			return err
		}
	}

	if _, err := w.Write([]byte{0x0D}); err != nil {
		return err
	}
	return nil
}
