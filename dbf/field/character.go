package field

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/everystreet/go-shapefile/cpg"
)

// Character field is a string of characters.
type Character struct {
	Field
	String string
}

// DecodeCharacter decodes a single character field with the specified encoding.
func DecodeCharacter(buf []byte, name string, encoding cpg.CharacterEncoding) (*Character, error) {
	val := bytes.Trim(buf, "\x00")

	switch encoding {
	case cpg.EncodingASCII:
		fallthrough
	case cpg.EncodingUTF8:
		return &Character{
			Field:  Field{name: name},
			String: strings.TrimSpace(string(val)),
		}, nil
	default:
		return nil, fmt.Errorf("unsupported character encoding")
	}
}

// Value returns the field value.
func (c *Character) Value() interface{} {
	return c.String
}
