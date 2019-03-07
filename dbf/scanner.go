package dbf

import (
	"io"
	"sync"
)

type Scanner struct {
	in io.Reader

	headerOnce sync.Once
	header     Header

	scanOnce  sync.Once
	recordsCh chan *Record

	errOnce sync.Once
	err     error
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		in:        r,
		recordsCh: make(chan *Record),
	}
}

func (s *Scanner) Header() (Header, error) {
	var err error
	s.headerOnce.Do(func() {
		var h Header
		if h, err = DecodeHeader(s.in); err != nil {
			return
		}
		s.header = h
	})
	return s.header, err
}

func (s *Scanner) Scan() error {
	if _, err := s.Header(); err != nil {
		return err
	}

	s.scanOnce.Do(func() {
		go func() {
			defer close(s.recordsCh)

			/*for {
				rec, err := s.record()
				if err == io.EOF {
					return
				} else if err != nil {
					s.setErr(err)
					return
				}
				s.decodeRecord(rec)
			}*/
		}()
	})
	return nil
}

func (s *Scanner) Record() *Record {
	rec, ok := <-s.recordsCh
	if !ok {
		return nil
	}
	return rec
}

func (s *Scanner) Err() error {
	return s.err
}
