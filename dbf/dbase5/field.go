package dbase5

import (
	"bytes"
	"fmt"
)

type FieldType uint8

const (
	CharacterType     FieldType = 'C'
	DateType          FieldType = 'D'
	FloatingPointType FieldType = 'F'
	LogicalType       FieldType = 'L'
	MemoType          FieldType = 'M'
	NumericType       FieldType = 'N'
)

type FieldDesc struct {
	Type FieldType

	name string
	len  uint8
}

func DecodeFieldDesc(buf []byte) (*FieldDesc, error) {
	if len(buf) < 32 {
		return nil, fmt.Errorf("expecting 32 bytes but have %d", len(buf))
	}

	name := bytes.Trim(buf[0:11], "\x00")
	return &FieldDesc{
		Type: FieldType(buf[11]),
		name: string(name),
		len:  buf[16],
	}, nil
}

func (f *FieldDesc) Name() string {
	return f.name
}
