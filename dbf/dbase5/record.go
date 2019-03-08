package dbase5

import (
	"fmt"

	"github.com/mercatormaps/go-shapefile/dbf/field"
	"github.com/pkg/errors"
)

type Record struct {
	deleted bool
	fields  map[string]Field
}

type Field interface {
	Name() string
}

func DecodeRecord(buf []byte, header *Header, conf Config) (*Record, error) {
	if len(buf) < 1 {
		return nil, fmt.Errorf("expecting 1 byte but have %d", len(buf))
	}

	rec := &Record{
		fields: make(map[string]Field, len(header.Fields)-len(conf.FilteredFields())),
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
		if len(buf) < (pos + int(desc.Length)) {
			return nil, errors.Wrapf(fmt.Errorf("expecting %d bytes but have %d", desc.Length, len(buf)-pos),
				fieldDecodeErr, desc.Name, i)
		}
		start, end := pos, pos+int(desc.Length)
		pos += int(desc.Length)

		// filter out unwanted fields
		if !wantField(desc.Name, conf.FilteredFields()) {
			continue
		}

		var f Field
		var err error

		switch desc.Type {
		case CharacterType:
			f, err = field.DecodeCharacter(buf[start:end], desc.Name, conf.CharacterEncoding())
		default:
			continue // TODO remove
			return nil, errors.Wrapf(fmt.Errorf("unsupported field type '%c'", desc.Type),
				fieldDecodeErr, desc.Name, i)
		}

		if err != nil {
			return nil, errors.Wrapf(err, fieldDecodeErr, desc.Name, i)
		}
		rec.fields[f.Name()] = f
	}

	return rec, nil
}

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
