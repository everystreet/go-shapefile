package field

import (
	"bytes"
	"fmt"
	"strconv"
)

// Numeric field.
type Numeric struct {
	Field
	Number float64
}

// DecodeNumeric decodes a single numeric field.
func DecodeNumeric(buf []byte, name string) (*Numeric, error) {
	val := bytes.Trim(buf, "\x20") // trim spaces
	num, err := strconv.ParseFloat(string(val), 0)
	if err != nil {
		return nil, fmt.Errorf("failed to parse number '%s': %w", string(val), err)
	}

	return &Numeric{
		Field:  Field{name: name},
		Number: num,
	}, nil
}

// Value returns the field value.
func (n *Numeric) Value() interface{} {
	return n.Number
}
