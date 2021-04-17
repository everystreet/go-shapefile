package field

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/text/encoding"
)

// Character field is a string of characters.
type Character struct {
	Field
	String string
}

// DecodeCharacter decodes a single character field with the specified encoding.
func DecodeCharacter(buf []byte, name string, decoder *encoding.Decoder) (*Character, error) {
	val := bytes.Trim(buf, "\x00")

	decVal, err := decoder.Bytes(val)
	if err != nil {
		return nil, fmt.Errorf("failed to decode value: %w", err)
	}

	return &Character{
		Field:  Field{name: name},
		String: strings.TrimSpace(string(decVal)),
	}, nil
}

// Value returns the field value.
func (c Character) Value() interface{} {
	return c.String
}

// Equal returns true if v contains the same value as c.
func (c Character) Equal(v string) bool {
	return v == c.String
}
