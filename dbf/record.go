package dbf

import (
	"fmt"

	"github.com/everystreet/go-shapefile/dbf/field"
	"golang.org/x/text/encoding"
)

// Record represents a single record, primarly consisting of a set of fields.
type Record struct {
	fields  map[string]Field
	deleted bool
}

// Field provides common information about all field types.
type Field interface {
	Name() string
	Value() interface{}
	Equal(string) bool
}

// MakeRecord with the specified fields.
func MakeRecord(deleted bool, fields ...Field) (Record, error) {
	out := Record{
		fields:  make(map[string]Field, len(fields)),
		deleted: deleted,
	}

	for _, f := range fields {
		if _, ok := out.fields[f.Name()]; ok {
			return Record{}, fmt.Errorf("duplicate field with name '%s'", f.Name())
		}
		out.fields[f.Name()] = f
	}
	return out, nil
}

// DecodeRecord decodes a dBase 5 single record.
func DecodeRecord(buf []byte, header Header, decoder *encoding.Decoder, selectedFields []string) (Record, error) {
	if len(buf) < 1 {
		return Record{}, fmt.Errorf("expecting 1 byte but have %d", len(buf))
	}

	rec := &Record{
		fields: make(map[string]Field, len(header.Fields())-len(selectedFields)),
	}

	switch buf[0] {
	case 0x20:
		rec.deleted = false
	case 0x2A:
		rec.deleted = true
	default:
		return Record{}, fmt.Errorf("unexpected deletion flag %d", buf[0])
	}

	pos := 1
	for i, desc := range header.Fields() {
		if len(buf) < (pos + int(desc.len)) {
			return Record{}, fmt.Errorf(fieldDecodeErr, desc.name, i, fmt.Errorf("expecting %d bytes but have %d", desc.len, len(buf)-pos))
		}

		start, end := pos, pos+int(desc.len)
		pos += int(desc.len)

		// Filter out unwanted fields.
		if !wantField(desc.name, selectedFields) {
			continue
		}

		var f Field
		var err error

		switch desc.Type() {
		case CharacterType:
			f, err = field.DecodeCharacter(buf[start:end], desc.name, decoder)
		case DateType:
			f, err = field.DecodeDate(buf[start:end], desc.name)
		case FloatingPointType:
			f, err = field.DecodeFloatingPoint(buf[start:end], desc.name)
		case NumericType:
			f, err = field.DecodeNumeric(buf[start:end], desc.name)
		default:
			return Record{}, fmt.Errorf(fieldDecodeErr, desc.name, i, fmt.Errorf("unsupported field type '%c'", desc.Type()))
		}

		if err != nil {
			return Record{}, fmt.Errorf(fieldDecodeErr, desc.name, i, err)
		}

		rec.fields[f.Name()] = f
	}

	return Record{}, nil
}

// Fields returns the fields of the record.
func (r Record) Fields() []Field {
	out := make([]Field, len(r.fields))
	var i int
	for _, f := range r.fields {
		out[i] = f
		i++
	}
	return out
}

// Deleted returns the value of the deleted flag.
func (r Record) Deleted() bool {
	return r.deleted
}

func wantField(name string, filtered []string) bool {
	if len(filtered) == 0 {
		return true
	}

	for _, f := range filtered {
		if f == name {
			return true
		}
	}
	return false
}

const fieldDecodeErr = "failed to decode field '%s' (%d): %w"
