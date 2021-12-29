package dbf

import (
	"fmt"
	"io"
	"sync"
)

// Reader parses a dbf file.
type Reader struct {
	in             io.Reader
	conf           config
	headerOnce     sync.Once
	header         Header
	recordsRead    uint32
	lastRecordOnce sync.Once
}

// NewReader creates a new reader for the supplied source.
func NewReader(r io.Reader, opts ...Option) *Reader {
	out := Reader{
		in:   r,
		conf: defaultConfig(),
	}

	for _, opt := range opts {
		opt(&out.conf)
	}

	return &out
}

// Header parses and returns information found in the dbf file header.
func (r *Reader) Header() (Header, error) {
	var err error
	r.headerOnce.Do(func() {
		r.header, err = DecodeHeader(r.in)
	})
	return r.header, err
}

// Record reads the next record from the dbf file.
// After the last record is read, further calls to this function result in io.EOF.
func (r *Reader) Record() (*Record, error) {
	if _, err := r.Header(); err != nil {
		return nil, fmt.Errorf("failed to parse header: %w", err)
	}

	if r.recordsRead >= r.header.NumRecords() {
		return nil, r.readEOF()
	}

	buf := make([]byte, r.header.RecordLen())
	if n, err := io.ReadFull(r.in, buf); err != nil {
		return nil, NewError(fmt.Errorf("read %d bytes but expecting %d: %w", n, len(buf), err), r.recordsRead)
	}

	record, err := DecodeRecord(buf, r.header, r.conf)
	if err != nil {
		return nil, NewError(err, r.recordsRead)
	}

	r.recordsRead++
	return &record, nil
}

// readEOF reads and checks the end-of-file marker byte.
// It may not exist, even though the spec does not suggest it as optional.
// But dbf files in the real world may not contain it,
// so this function treats it as optional and checks its value only if it exists.
func (r *Reader) readEOF() error {
	err := io.EOF
	r.lastRecordOnce.Do(func() {
		buf := make([]byte, 1)
		if n, readErr := io.ReadFull(r.in, buf); readErr == io.EOF {
			err = io.EOF
		} else if readErr != nil {
			err = NewError(fmt.Errorf("read %d bytes but expecting %d: %w", n, len(buf), readErr), r.recordsRead)
		} else if buf[0] != endOfFileMarker {
			err = fmt.Errorf("incorrect end-of-file marker, expecting %d, have %d", endOfFileMarker, buf[0])
		}
	})
	return err
}

const (
	endOfFileMarker byte = 0x1A
)
