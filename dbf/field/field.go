package field

type Field struct {
	name string
}

func (f Field) Name() string {
	return f.name
}
