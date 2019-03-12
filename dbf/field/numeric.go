package field

import (
	"bytes"
	"strconv"

	"github.com/pkg/errors"
)

type Numeric struct {
	Field
	Number float64
}

func DecodeNumeric(buf []byte, name string) (*Numeric, error) {
	val := bytes.Trim(buf, "\x20") // trim spaces
	num, err := strconv.ParseFloat(string(val), 0)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse number '%s'", string(val))
	}

	return &Numeric{
		Field:  Field{name: name},
		Number: num,
	}, nil
}

func (n *Numeric) Value() interface{} {
	return n.Number
}
