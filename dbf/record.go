package dbf

import "github.com/everystreet/go-shapefile/dbf/dbase5"

// Record wraps a dBase level-specific record.
type Record struct {
	rec interface {
		Deleted() bool
	}
}

// Field provides generic access to record fields of any type.
type Field interface {
	Name() string
	Value() interface{}
	Equal(string) bool
}

// Fields returns a list of all fields in the record.
// The order of the fields is nondeterministic.
func (r Record) Fields() []Field {
	switch rec := r.rec.(type) {
	case *dbase5.Record:
		fields := make([]Field, len(rec.Fields))
		i := 0
		for _, f := range rec.Fields {
			fields[i] = f
			i++
		}
		return fields
	default:
		return nil
	}
}

// Field returns a field by name.
func (r Record) Field(name string) (Field, bool) {
	switch rec := r.rec.(type) {
	case *dbase5.Record:
		f, ok := rec.Fields[name]
		if !ok {
			return nil, false
		}
		return f.(Field), true
	default:
		return nil, false
	}
}

// Deleted returns the state of the "deleted" marker.
func (r Record) Deleted() bool {
	return r.rec.Deleted()
}
