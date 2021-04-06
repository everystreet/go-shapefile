package shp

import (
	"encoding/binary"
	"fmt"
	"io"
	"sync"
)

// Scanner parses a shp file.
type Scanner struct {
	in   io.Reader
	opts []Option

	headerOnce sync.Once
	header     Header

	scanOnce sync.Once
	shapesCh chan Shape

	errOnce sync.Once
	err     error
}

// NewScanner creates a new Scanner for the supplied source.
func NewScanner(r io.Reader, opts ...Option) *Scanner {
	return &Scanner{
		in:       r,
		opts:     opts,
		shapesCh: make(chan Shape),
	}
}

// AddOptions to scanner after creation.
func (s *Scanner) AddOptions(opts ...Option) {
	s.opts = append(s.opts, opts...)
}

// Header parses the shp file header.
func (s *Scanner) Header() (Header, error) {
	var err error
	s.headerOnce.Do(func() {
		buf := make([]byte, 100)
		var n int
		if n, err = io.ReadFull(s.in, buf); err != nil {
			err = fmt.Errorf("expecting to read %d bytes but only read %d: %w", len(buf), n, err)
			return
		}

		if s.header, err = DecodeHeader(buf, s.opts...); err != nil {
			return
		}
	})
	return s.header, err
}

// Validator returns a Validator that can be used to validate Shapes using the Validate method.
func (s *Scanner) Validator() (Validator, error) {
	h, err := s.Header()
	if err != nil {
		return Validator{}, fmt.Errorf("failed to decode header: %w", err)
	}
	return MakeValidator(h.BoundingBox)
}

// Scan starts reading the shp file. Shapes can be accessed from the Shape method.
// An error is returned if there's a problem parsing the header.
// Errors that are encountered when parsing records must be checked with the Err method.
func (s *Scanner) Scan() error {
	var conf config
	for _, opt := range s.opts {
		opt(&conf)
	}

	if _, err := s.Header(); err != nil {
		return fmt.Errorf("failed to parse header: %w", err)
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
				s.decodeRecord(*rec, conf)
			}
		}()
	})
	return nil
}

// Shape returns each shape found in the shp file.
// nil is returned once the last record has been read, or an error occurs -
// the Err method should be used to check for an error at this point.
func (s *Scanner) Shape() Shape {
	shape, ok := <-s.shapesCh
	if !ok {
		return nil
	}
	return shape
}

// Err returns the first error encountered when parsing records.
// It should be called after calling the Shape method for the last time.
func (s *Scanner) Err() error {
	return s.err
}

func (s *Scanner) decodeRecord(rec record, conf config) {
	if rec.shapeType == NullType {
		return
	} else if rec.shapeType != s.header.ShapeType {
		err := fmt.Errorf("shape type %d differs to specified type %d", rec.shapeType, s.header.ShapeType)
		s.setErr(NewError(err, rec.number))
		return
	}

	var shape Shape
	var err error
	if conf.precision == nil {
		shape, err = s.decodeShape(rec)
	} else {
		shape, err = s.decodeShapeP(rec, *conf.precision)
	}

	if err != nil {
		s.setErr(NewError(err, rec.number))
		return
	}
	s.shapesCh <- shape
}

func (s *Scanner) decodeShape(rec record) (Shape, error) {
	switch rec.shapeType {
	case PointType:
		return DecodePoint(rec.shape, rec.number)
	case PolylineType:
		return DecodePolyline(rec.shape, rec.number)
	case PolygonType:
		return DecodePolygon(rec.shape, rec.number)
	default:
		return nil, fmt.Errorf("unknown shape type %d", rec.shapeType)
	}
}

func (s *Scanner) decodeShapeP(rec record, precision uint) (Shape, error) {
	switch rec.shapeType {
	case PointType:
		return DecodePointP(rec.shape, rec.number, precision)
	case PolylineType:
		return DecodePolylineP(rec.shape, rec.number, precision)
	case PolygonType:
		return DecodePolygonP(rec.shape, rec.number, precision)
	default:
		return nil, fmt.Errorf("unknown shape type %d", rec.shapeType)
	}
}

func (s *Scanner) record() (*record, error) {
	buf := make([]byte, 12)
	if _, err := io.ReadFull(s.in, buf); err != nil {
		return nil, io.EOF
	}

	num := binary.BigEndian.Uint32(buf[0:4])

	shapeType := ShapeType(binary.LittleEndian.Uint32(buf[8:12]))
	if shapeType != s.header.ShapeType {
		return nil, NewError(fmt.Errorf("unexpected shape type; expecting %d, got %d", s.header.ShapeType, shapeType), num)
	}

	length := binary.BigEndian.Uint32(buf[4:8]) * 2 // length is in 16-byte words, so multiply by 2 to get bytes

	// length is the length of the record, which consists of the shape type and shape data
	// we've already read the shape type (4 bytes), so the shape data is the next length-4 bytes
	buf = make([]byte, length-4)
	if _, err := io.ReadFull(s.in, buf); err != nil {
		return nil, io.EOF
	}

	return &record{
		number:    num,
		length:    length,
		shapeType: ShapeType(shapeType),
		shape:     buf,
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
