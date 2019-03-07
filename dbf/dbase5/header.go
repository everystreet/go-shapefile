package dbase5

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/pkg/errors"
)

type Header struct {
	Fields []Field

	numRecs uint32
}

func DecodeHeader(r io.Reader) (*Header, error) {
	// Read first 31 bytes after first byte
	buf := make([]byte, 31)
	if n, err := r.Read(buf); err != nil {
		return nil, err
	} else if n != len(buf) {
		return nil, fmt.Errorf("read %d bytes but expecting %d", n, len(buf))
	}

	out := &Header{
		numRecs: binary.LittleEndian.Uint32(buf[3:7]),
	}

	// Read remainder of header
	headerLen := binary.LittleEndian.Uint16(buf[7:9])
	buf = make([]byte, int(headerLen)-len(buf)-1)
	if n, err := r.Read(buf); err != nil {
		return nil, err
	} else if n != len(buf) {
		return nil, fmt.Errorf("read %d bytes but expecting %d", n, len(buf))
	}

	if (len(buf)-1)%32 != 0 {
		return nil, fmt.Errorf("invalid header size %d bytes", headerLen)
	}

	if buf[len(buf)-1] != 0x0D {
		return nil, fmt.Errorf("missing field descriptor terminator")
	}

	numFields := (len(buf) - 1) / 32
	out.Fields = make([]Field, numFields)
	for i := 0; i < numFields; i++ {
		f, err := DecodeField(buf[i*32 : (i*32)+32])
		if err != nil {
			return nil, errors.Wrapf(err, "failed to decode field %d", i)
		}
		out.Fields[i] = *f
	}

	return out, nil
}

func (h *Header) NumRecords() uint32 {
	return h.numRecs
}
