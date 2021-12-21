package field

import (
	"bytes"
	"time"
)

// Date field is a date with no time component.
type Date struct {
	Field
	Date *time.Time
}

// DecodeDate decodes a single date field.
func DecodeDate(buf []byte, name string) (*Date, error) {
	val := bytes.Trim(buf, "\x00\x20")

	out := &Date{
		Field: Field{name: name},
	}

	if len(val) == 0 {
		return out, nil
	}

	date, err := time.Parse("01/02/2006", string(val))
	if err != nil {
		return nil, err
	}

	return &Date{
		Field: Field{name: name},
		Date:  &date,
	}, nil
}

// Value returns the field value.
func (d Date) Value() interface{} {
	return d.Date
}

// Equal returns true if v contains the same value as c.
func (d Date) Equal(v string) bool {
	d2, err := time.Parse("01/02/2006", v)
	if err != nil {
		return false
	}

	return d.Date.Year() == d2.Year() && d.Date.Month() == d2.Month() && d.Date.Day() == d2.Day()
}
