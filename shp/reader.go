package shp

import (
	"encoding/binary"
	"fmt"
	"io"
	"sync"
)

// Reader parses a shp file.
type Reader struct {
	in         io.Reader
	conf       readerConfig
	headerOnce sync.Once
	header     Header
}

// NewReader creates a new reader for the supplied source.
func NewReader(r io.Reader, opts ...ReaderOption) *Reader {
	out := Reader{
		in:   r,
		conf: defaultConfig(),
	}

	for _, opt := range opts {
		opt(&out.conf)
	}
	return &out
}

// Header parses the shp file header.
func (r *Reader) Header() (Header, error) {
	var err error
	r.headerOnce.Do(func() {
		buf := make([]byte, 100)
		if n, readErr := io.ReadFull(r.in, buf); readErr != nil {
			err = fmt.Errorf("expecting to read %d bytes but only read %d: %w", len(buf), n, readErr)
			return
		}

		r.header, err = DecodeHeader(buf, r.conf.precision)
	})
	return r.header, err
}

// Validator returns a Validator that can be passed to Shape.Validate().
func (r *Reader) Validator() (Validator, error) {
	h, err := r.Header()
	if err != nil {
		return Validator{}, fmt.Errorf("failed to decode header: %w", err)
	}

	return MakeValidator(h.BoundingBox)
}

// Shape reads the next shape from the shp file.
// After the last shape is read, further calls to this function result in io.EOF.
func (r *Reader) Shape() (Shape, error) {
	record, err := r.readRecord()
	if err != nil {
		return nil, err
	}

	return r.decodeRecord(record)
}

func (r *Reader) decodeRecord(rec record) (Shape, error) {
	h, err := r.Header()
	if err != nil {
		return nil, fmt.Errorf("failed to decode header: %w", err)
	}

	if rec.shapeType == NullType {
		return nil, nil
	} else if rec.shapeType != h.ShapeType {
		return nil, fmt.Errorf("shape type %d differs from expected type %d", rec.shapeType, h.ShapeType)
	}

	switch h.ShapeType {
	case PointType:
		if r.conf.precision == nil {
			return DecodePoint(rec.shape, rec.number)
		}
		return DecodePointP(rec.shape, rec.number, *r.conf.precision)
	case PolylineType:
		if r.conf.precision == nil {
			return DecodePolyline(rec.shape, rec.number)
		}
		return DecodePolylineP(rec.shape, rec.number, *r.conf.precision)
	case PolygonType:
		if r.conf.precision == nil {
			return DecodePolygon(rec.shape, rec.number)
		}
		return DecodePolygonP(rec.shape, rec.number, *r.conf.precision)
	default:
		return nil, fmt.Errorf("unknown shape type %d", h.ShapeType)
	}
}

func (r *Reader) readRecord() (record, error) {
	buf := make([]byte, 12)
	if _, err := io.ReadFull(r.in, buf); err != nil {
		return record{}, io.EOF
	}

	num := binary.BigEndian.Uint32(buf[0:4])

	shapeType := ShapeType(binary.LittleEndian.Uint32(buf[8:12]))

	// length is in 16-byte words, so multiply by 2 to get bytes.
	length := binary.BigEndian.Uint32(buf[4:8]) * 2

	// length is the length of the record, which consists of the shape type and shape data.
	// We've already read the shape type (4 bytes), so the shape data is the next `length-4` bytes.
	buf = make([]byte, length-4)
	if _, err := io.ReadFull(r.in, buf); err != nil {
		return record{}, io.EOF
	}

	return record{
		number:    num,
		length:    length,
		shapeType: ShapeType(shapeType),
		shape:     buf,
	}, nil
}

type record struct {
	number    uint32
	length    uint32
	shapeType ShapeType
	shape     []byte
}

// PointPrecision sets the precision of coordinates.
func PointPrecision(p uint) ReaderOption {
	return func(c *readerConfig) {
		c.precision = &p
	}
}

// ReaderOption funcs can be passed to reading operations.
type ReaderOption func(*readerConfig)

// Config for shp parsing.
type readerConfig struct {
	precision *uint
}

func defaultConfig() readerConfig {
	return readerConfig{}
}
