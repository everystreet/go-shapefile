package dbase5

import (
	"fmt"

	"github.com/everystreet/go-shapefile/dbf/field"
	"golang.org/x/text/encoding"
)

// Record represents a single record, primarly consisting of a set of fields.
type Record struct {
	Fields map[string]Field

	deleted bool
}

// Field provides common information about all field types.
type Field interface {
	Name() string
	Value() interface{}
	Equal(string) bool
}

// Config provides config for record parsing.
type Config interface {
	CharacterDecoder() *encoding.Decoder
	FilteredFields() []string
}

// DecodeRecord decodes a dBase 5 single record.
func DecodeRecord(buf []byte, header *Header, conf Config) (*Record, error) {
	if len(buf) < 1 {
		return nil, fmt.Errorf("expecting 1 byte but have %d", len(buf))
	}

	rec := &Record{
		Fields: make(map[string]Field, len(header.Fields)-len(conf.FilteredFields())),
	}

	switch buf[0] {
	case 0x20:
		rec.deleted = false
	case 0x2A:
		rec.deleted = true
	default:
		return nil, fmt.Errorf("missing deletion flag")
	}

	pos := 1
	for i, desc := range header.Fields {
		if len(buf) < (pos + int(desc.len)) {
			return nil, fmt.Errorf(fieldDecodeErr, desc.name, i,
				fmt.Errorf("expecting %d bytes but have %d", desc.len, len(buf)-pos))
		}
		start, end := pos, pos+int(desc.len)
		pos += int(desc.len)

		// filter out unwanted fields
		if !wantField(desc.name, conf.FilteredFields()) {
			continue
		}

		var f Field
		var err error

		switch desc.Type {
		case CharacterType:
			f, err = field.DecodeCharacter(buf[start:end], desc.name, conf.CharacterDecoder())
		case DateType:
			f, err = field.DecodeDate(buf[start:end], desc.name)
		case FloatingPointType:
			f, err = field.DecodeFloatingPoint(buf[start:end], desc.name)
		case NumericType:
			f, err = field.DecodeNumeric(buf[start:end], desc.name)
		default:
			return nil, fmt.Errorf(fieldDecodeErr, desc.name, i,
				fmt.Errorf("unsupported field type '%c'", desc.Type))
		}

		if err != nil {
			return nil, fmt.Errorf(fieldDecodeErr, desc.name, i, err)
		}
		rec.Fields[f.Name()] = f
	}

	return rec, nil
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
