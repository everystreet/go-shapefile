package dbf

import (
	"bytes"
	"fmt"
)

// FieldType is the type of a field.
type FieldType uint8

// Field types for dBase Level 5.
const (
	CharacterType     FieldType = 'C'
	DateType          FieldType = 'D'
	FloatingPointType FieldType = 'F'
	LogicalType       FieldType = 'L'
	MemoType          FieldType = 'M'
	NumericType       FieldType = 'N'
)

// FieldDesc represents a field descriptor consisting of a type, name and size in bytes.
type FieldDesc struct {
	typ  FieldType
	name string
	len  uint8
}

func MakeFieldDesc(typ FieldType, name string, len uint8) FieldDesc {
	return FieldDesc{
		typ:  typ,
		name: name,
		len:  len,
	}
}

// DecodeFieldDesc parses a single field descriptor.
func DecodeFieldDesc(buf []byte) (FieldDesc, error) {
	if len(buf) < 32 {
		return FieldDesc{}, fmt.Errorf("expecting 32 bytes but have %d", len(buf))
	}

	typ := FieldType(buf[11])
	if err := validateFieldType(typ); err != nil {
		return FieldDesc{}, err
	}

	return FieldDesc{
		typ:  typ,
		name: string(bytes.Trim(buf[0:11], "\x00")),
		len:  buf[16],
	}, nil
}

// Type of the field.
func (f FieldDesc) Type() FieldType {
	return f.typ
}

// Name of the field.
func (f FieldDesc) Name() string {
	return f.name
}

// Length of the field.
func (f FieldDesc) Length() uint8 {
	return f.len
}

func (f FieldDesc) Encode() ([]byte, error) {
	out := make([]byte, 32)

	if err := validateFieldType(f.typ); err != nil {
		return nil, err
	}
	out[11] = byte(f.typ)

	name := []byte(f.name)
	if len(name) > 11 {
		return nil, fmt.Errorf("field name exceeds maximum length of 11 bytes (%d bytes)", len(name))
	}
	copy(out, name)

	out[16] = f.len

	return out, nil
}

func validateFieldType(typ FieldType) error {
	switch typ {
	case
		CharacterType,
		DateType,
		FloatingPointType,
		LogicalType,
		MemoType,
		NumericType:
		return nil
	default:
		return fmt.Errorf("unrecognized field type '%c'", typ)
	}
}
