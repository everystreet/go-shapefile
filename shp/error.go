package shp

import "fmt"

// Error describes an error that occured when parsing a shape.
type Error struct {
	recordNum uint32
	err       error
}

// NewError returns an attached to record number.
func NewError(err error, recordNum uint32) Error {
	return Error{
		err:       err,
		recordNum: recordNum,
	}
}

func (e Error) Error() string {
	return fmt.Sprintf("error reading record %d: %v", e.recordNum, e.err)
}
