package field

import (
	"bytes"
	"fmt"

	"github.com/mercatormaps/go-shapefile/cpg"
)

type Character struct {
	Field
	value string
}

func DecodeCharacter(buf []byte, name string, encoding cpg.CharacterEncoding) (*Character, error) {
	val := bytes.Trim(buf, "\x00")

	switch encoding {
	case cpg.ASCII:
		fallthrough
	case cpg.UTF8:
		return &Character{
			Field: Field{name: name},
			value: string(val),
		}, nil
	default:
		return nil, fmt.Errorf("unsupported character encoding")
	}
}

func (c *Character) Value() string {
	return c.value
}
