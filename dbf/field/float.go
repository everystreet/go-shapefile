package field

// FloatingPoint field.
type FloatingPoint Numeric

// DecodeFloatingPoint decodes a single floating point field.
func DecodeFloatingPoint(buf []byte, name string) (*FloatingPoint, error) {
	n, err := DecodeNumeric(buf, name)
	if err != nil {
		return nil, err
	}
	return (*FloatingPoint)(n), nil
}

// Value returns the field value.
func (f FloatingPoint) Value() interface{} {
	return f.Number
}

// Equal returns true if v contains the same value as f.
func (f FloatingPoint) Equal(v string) bool {
	return Numeric(f).Equal(v)
}
