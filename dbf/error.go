package dbf

import (
	"fmt"
)

type Error struct {
	recordNum uint32
	err       error
}

func NewError(err error, recordNum uint32) *Error {
	return &Error{
		err:       err,
		recordNum: recordNum,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("error reading record %d: %v", e.recordNum, e.err)
}
