package field

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/mercatormaps/go-shapefile/cpg"
)

type Character struct {
	Field
	String string
}

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

func (c *Character) Value() interface{} {
	return c.String
}
