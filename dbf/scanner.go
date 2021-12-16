package dbf

import (
	"fmt"
	"io"
	"sync"

	"github.com/everystreet/go-shapefile/dbf/dbase5"
)

// Version is the dBase version or "level".
type Version uint

// dBase versions.
const (
	DBaseLevel5 Version = 3
	DBaseLevel7 Version = 4
)

// Scanner parses a dbf file.
type Scanner struct {
	in io.Reader

	versionOnce sync.Once
	version     Version

	headerOnce sync.Once
	header     Header

	scanOnce  sync.Once
	recordsCh chan *Record
	num       uint32

	errOnce sync.Once
	err     error
}

// Header provides common information for all dbf version headers.
type Header interface {
	RecordLen() uint16
	NumRecords() uint32
	FieldExists(string) bool
}

// NewScanner creates a new Scanner for the supplied source.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		in:        r,
		recordsCh: make(chan *Record),
	}
}

// Version reads and returns the dBase version.
// This function does not validate the version number.
func (s *Scanner) Version() (Version, error) {
	var err error
	s.versionOnce.Do(func() {
		buf := make([]byte, 1)
		var n int
		if n, err = io.ReadFull(s.in, buf); err != nil {
			err = fmt.Errorf("read %d bytes but expecting %d: %w", n, len(buf), err)
			return
		}

		// dBase version number is first 3 bits
		s.version = Version(((buf[0]>>0)&1)<<0 | ((buf[0]>>1)&1)<<1 | ((buf[0]>>2)&1)<<2)
	})
	return s.version, err
}

// Header invokes the Header method for the appropriate dbase version scanner.
// A type assertion can be used to access information specific to the version.
func (s *Scanner) Header() (Header, error) {
	var err error
	if _, err = s.Version(); err != nil {
		return nil, fmt.Errorf("failed to parse version number: %w", err)
	}

	s.headerOnce.Do(func() {
		switch s.version {
		case DBaseLevel5:
			s.header, err = dbase5.DecodeHeader(s.in)
		case DBaseLevel7:
			err = fmt.Errorf("dBase Level 7 is not supported")
		default:
			err = fmt.Errorf("unsupported version")
		}
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
				rec, err := s.record()
				if err != nil {
					s.setErr(err)
					return
				}
				s.decodeRecord(rec, conf)
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

func (s *Scanner) decodeRecord(buf []byte, conf config) {
	switch s.version {
	case DBaseLevel5:
		rec, err := dbase5.DecodeRecord(buf, s.header.(*dbase5.Header), conf)
		if err != nil {
			s.setErr(NewError(err, s.num))
			return
		}
		s.recordsCh <- &Record{
			rec: rec,
		}
		s.num++
	case DBaseLevel7:
		s.setErr(fmt.Errorf("dBase Level 7 is not supported"))
	default:
		s.setErr(fmt.Errorf("unsupported version"))
	}
}

func (s *Scanner) record() ([]byte, error) {
	buf := make([]byte, s.header.RecordLen())
	if n, err := io.ReadFull(s.in, buf); err != nil {
		return nil, NewError(fmt.Errorf("read %d bytes but expecting %d: %w", n, len(buf), err), s.num)
	}
	return buf, nil
}

func (s *Scanner) setErr(err error) {
	s.errOnce.Do(func() {
		s.err = err
	})
}
