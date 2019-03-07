package dbf

import (
	"fmt"
	"io"

	"github.com/mercatormaps/go-shapefile/dbf/dbase5"
)

const (
	DBaseLevel5 = 3
	DBaseLevel7 = 4
)

type Header interface {
	NumRecords() uint32
}

func DecodeHeader(r io.Reader) (Header, error) {
	buf := make([]byte, 1)
	if n, err := r.Read(buf); err != nil {
		return nil, err
	} else if n != len(buf) {
		return nil, fmt.Errorf("read %d bytes but expecting %d", n, len(buf))
	}

	// dBase version number is first 3 bits
	version := ((buf[0]>>0)&1)<<0 | ((buf[0]>>1)&1)<<1 | ((buf[0]>>2)&1)<<2

	switch version {
	case DBaseLevel5:
		return dbase5.DecodeHeader(r)
	case DBaseLevel7:
		return nil, fmt.Errorf("dBase Level 7 is not supported")
	default:
		return nil, fmt.Errorf("unsupported version")
	}
}
