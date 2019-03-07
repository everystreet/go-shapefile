package dbase5

import (
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

type Field struct {
	Name string
	Type FieldType
}

func DecodeField(buf []byte) (*Field, error) {
	if len(buf) < 32 {
		return nil, fmt.Errorf("expecting 32 bytes but have %d", len(buf))
	}

	return &Field{
		Name: string(buf[0:11]),
		Type: FieldType(buf[11]),
	}, nil
}
