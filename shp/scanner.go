package shp

import (
	"fmt"
	"io"
	"sync"
)

type Scanner struct {
	in io.Reader

	headerOnce sync.Once
	header     Header

	scanOnce sync.Once
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		in: r,
	}
}

func (s *Scanner) Header() (*Header, error) {
	var err error
	s.headerOnce.Do(func() {
		b := make([]byte, 100)
		var n int
		if n, err = s.in.Read(b); err != nil {
			return
		} else if n != 100 {
			err = fmt.Errorf("expecting to read 100 bytes but only read %d", n)
			return
		}

		var h *Header
		if h, err = DecodeHeader(b); err != nil {
			return
		}
		s.header = *h
	})
	return &s.header, err
}
