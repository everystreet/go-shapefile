package dbase5

import (
	"fmt"
)

type Record struct {
	deleted bool
}

func DecodeRecord(buf []byte) (*Record, error) {
	if len(buf) < 1 {
		return nil, fmt.Errorf("expecting 1 byte but have %d", len(buf))
	}

	rec := &Record{}
	switch buf[0] {
	case 0x20:
		rec.deleted = false
	case 0x2A:
		rec.deleted = true
	default:
		return nil, fmt.Errorf("missing deletion flag")
	}

	return rec, nil
}

func (r *Record) Deleted() bool {
	return r.deleted
}
