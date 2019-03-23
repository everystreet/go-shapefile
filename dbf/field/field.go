package field

// Field base type.
type Field struct {
	name string
}

// Name returns the field name.
func (f Field) Name() string {
	return f.name
}
