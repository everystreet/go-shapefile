package dbf

import (
	"fmt"
	"io"
	"sync"
)

// Scanner parses a dbf file.
type Scanner struct {
	in io.Reader

	headerOnce sync.Once
	header     *Header

	scanOnce  sync.Once
	recordsCh chan *Record
	num       uint32

	errOnce sync.Once
	err     error
}

// NewScanner creates a new scanner for the supplied source.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		in:        r,
		recordsCh: make(chan *Record),
	}
}

// Header parses and returns information found in the dbf file header.
func (s *Scanner) Header() (*Header, error) {
	var err error
	s.headerOnce.Do(func() {
		s.header, err = DecodeHeader(s.in)
	})
	return s.header, err
}

// Scan starts reading the dbf file. Records can be accessed from the Record method.
// An error is returned if there's a problem parsing the header.
// Errors that are encountered when parsing records must be checked with the Err method.
func (s *Scanner) Scan(opts ...Option) error {
	conf := defaultConfig()
	for _, opt := range opts {
		opt(&conf)
	}

	if _, err := s.Header(); err != nil {
		return fmt.Errorf("failed to parse header: %w", err)
	}

	s.scanOnce.Do(func() {
		go func() {
			defer close(s.recordsCh)

			for s.num < s.header.NumRecords() {
				buf := make([]byte, s.header.RecordLen())
				if n, err := io.ReadFull(s.in, buf); err != nil {
					s.setErr(NewError(fmt.Errorf("read %d bytes but expecting %d: %w", n, len(buf), err), s.num))
					return
				}

				rec, err := DecodeRecord(buf, s.header, conf)
				if err != nil {
					s.setErr(NewError(err, s.num))
					return
				}

				s.recordsCh <- rec
				s.num++
			}

			buf := make([]byte, 1)
			if n, err := io.ReadFull(s.in, buf); err == io.EOF {
				return
			} else if err != nil {
				s.setErr(fmt.Errorf("read %d bytes but expecting %d: %w", n, len(buf), err))
				return
			}

			if buf[0] != 0x1A {
				s.setErr(fmt.Errorf("missing file terminator"))
			}
		}()
	})
	return nil
}

// Record returns each record found in the dbf file.
// nil is returned once the last record has been read, or an error occurs -
// the Err method should be used to check for an error at this point.
func (s *Scanner) Record() *Record {
	rec, ok := <-s.recordsCh
	if !ok {
		return nil
	}
	return rec
}

// Err returns the first error encountered when parsing records.
// It should be called after calling the Record method for the last time.
func (s *Scanner) Err() error {
	return s.err
}

func (s *Scanner) setErr(err error) {
	s.errOnce.Do(func() {
		s.err = err
	})
}
