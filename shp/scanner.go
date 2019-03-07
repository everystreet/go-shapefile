package shp

import (
	"encoding/binary"
	"fmt"
	"io"
	"sync"

	"github.com/pkg/errors"
)

type Scanner struct {
	in io.Reader

	headerOnce sync.Once
	header     Header

	scanOnce sync.Once
	shapesCh chan Shape

	errOnce sync.Once
	err     error
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		in:       r,
		shapesCh: make(chan Shape),
	}
}

func (s *Scanner) Header() (*Header, error) {
	var err error
	s.headerOnce.Do(func() {
		buf := make([]byte, 100)
		var n int
		if n, err = s.in.Read(buf); err != nil {
			return
		} else if n != len(buf) {
			err = fmt.Errorf("expecting to read %d bytes but only read %d", len(buf), n)
			return
		}

		var h *Header
		if h, err = DecodeHeader(buf); err != nil {
			return
		}
		s.header = *h
	})
	return &s.header, err
}

func (s *Scanner) Scan() error {
	if _, err := s.Header(); err != nil {
		return errors.Wrap(err, "failed to parse header")
	}

	s.scanOnce.Do(func() {
		go func() {
			defer close(s.shapesCh)

			for {
				rec, err := s.record()
				if err == io.EOF {
					return
				} else if err != nil {
					s.setErr(err)
					return
				}
				s.decodeRecord(rec)
			}
		}()
	})
	return nil
}

func (s *Scanner) Shape() Shape {
	shape, ok := <-s.shapesCh
	if !ok {
		return nil
	}
	return shape
}

func (s *Scanner) Err() error {
	return s.err
}

func (s *Scanner) decodeRecord(rec *record) {
	switch rec.shapeType {
	case PolygonType:
		p, err := DecodePolygon(rec.shape, rec.number)
		if err != nil {
			s.setErr(NewError(err, rec.number))
			return
		}
		s.shapesCh <- p
	default:
		s.setErr(NewError(fmt.Errorf("unknown shape type %d", rec.shapeType), rec.number))
	}
}

func (s *Scanner) record() (*record, error) {
	b := make([]byte, 12)
	if n, err := s.in.Read(b); err != nil {
		return nil, err
	} else if n != 12 {
		return nil, fmt.Errorf("expecting to read 12 bytes but only read %d", n)
	}

	num := binary.BigEndian.Uint32(b[0:4])

	shapeType := ShapeType(binary.LittleEndian.Uint32(b[8:12]))
	if shapeType != s.header.ShapeType {
		return nil, NewError(fmt.Errorf("unexpected shape type; expecting %d, got %d", s.header.ShapeType, shapeType), num)
	}

	length := binary.BigEndian.Uint32(b[4:8]) * 2 // length is in 16-byte words, so multiply by 2 to get bytes

	// length is the length of the record, which consists of the shape type and shape data
	// we've already read the shape type (4 bytes), so the shape data is the next length-4 bytes
	b = make([]byte, length-4)
	if n, err := s.in.Read(b); err != nil {
		return nil, err
	} else if n != len(b) {
		return nil, NewError(fmt.Errorf("expecting to read %d bytes but only read %d", len(b), n), num)
	}

	return &record{
		number:    num,
		length:    length,
		shapeType: ShapeType(shapeType),
		shape:     b,
	}, nil
}

func (s *Scanner) setErr(err error) {
	s.errOnce.Do(func() {
		s.err = err
	})
}

type record struct {
	number    uint32
	length    uint32
	shapeType ShapeType
	shape     []byte
}
