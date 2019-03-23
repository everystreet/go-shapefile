package dbase5

import (
	"fmt"

	"github.com/mercatormaps/go-shapefile/cpg"
	"github.com/mercatormaps/go-shapefile/dbf/field"
	"github.com/pkg/errors"
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
}

// Config provides config for record parsing.
type Config interface {
	CharacterEncoding() cpg.CharacterEncoding
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
			return nil, errors.Wrapf(fmt.Errorf("expecting %d bytes but have %d", desc.len, len(buf)-pos),
				fieldDecodeErr, desc.name, i)
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
			f, err = field.DecodeCharacter(buf[start:end], desc.name, conf.CharacterEncoding())
		case NumericType:
			f, err = field.DecodeNumeric(buf[start:end], desc.name)
		default:
			return nil, errors.Wrapf(fmt.Errorf("unsupported field type '%c'", desc.Type),
				fieldDecodeErr, desc.name, i)
		}

		if err != nil {
			return nil, errors.Wrapf(err, fieldDecodeErr, desc.name, i)
		}
		rec.Fields[f.Name()] = f
	}

	return rec, nil
}

// Deleted returns the value of the deleted flag.
func (r *Record) Deleted() bool {
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

const fieldDecodeErr = "failed to decode field '%s' (%d)"
