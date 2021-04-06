package shapefile

import (
	"fmt"
	"io"
	"sync"

	"github.com/everystreet/go-shapefile/dbf"
	"github.com/everystreet/go-shapefile/dbf/dbase5"
	"github.com/everystreet/go-shapefile/shp"
)

// Scanner parses a pair or shp and dbf files.
type Scanner struct {
	shp  *shp.Scanner
	dbf  *dbf.Scanner
	opts options

	infoOnce sync.Once
	info     Info

	scanOnce  sync.Once
	recordsCh chan *Record

	errOnce sync.Once
	err     error
}

// Info contains combined information from the pair of input files.
type Info struct {
	BoundingBox shp.BoundingBox
	NumRecords  uint32
	ShapeType   shp.ShapeType
	Fields      FieldDescList
}

// FieldDescList is a list of field descriptors.
type FieldDescList []FieldDesc

// Exists reutnrs true if the named field exists, and false otherwise.
func (l FieldDescList) Exists(name string) bool {
	for _, f := range l {
		if f.Name() == name {
			return true
		}
	}
	return false
}

// FieldDesc provides information about an attribute field.
type FieldDesc interface {
	Name() string
}

// NewScanner creates a new Scanner for the provided shp and dbf files.
func NewScanner(shpR, dbfR io.Reader, opts ...Option) *Scanner {
	s := &Scanner{
		dbf:       dbf.NewScanner(dbfR),
		recordsCh: make(chan *Record),
	}

	for _, opt := range opts {
		opt(&s.opts)
	}
	s.shp = shp.NewScanner(shpR, s.opts.shp...)
	return s
}

// AddOptions allows additional options to be set after the scanner has already been created.
func (s *Scanner) AddOptions(opts ...Option) {
	for _, opt := range opts {
		opt(&s.opts)
	}
	s.shp.AddOptions(s.opts.shp...)
}

// Info returns combined information about the shp and dbf pair.
func (s *Scanner) Info() (*Info, error) {
	var err error

	s.infoOnce.Do(func() {
		var shpHeader shp.Header
		if shpHeader, err = s.shp.Header(); err != nil {
			err = fmt.Errorf("failed to parse shp header: %w", err)
			return
		}

		var dbfHeader dbf.Header
		if dbfHeader, err = s.dbf.Header(); err != nil {
			err = fmt.Errorf("failed to parse dbf header: %w", err)
			return
		}

		var fields []FieldDesc
		switch h := dbfHeader.(type) {
		case *dbase5.Header:
			fields = make([]FieldDesc, len(h.Fields))
			for i, f := range h.Fields {
				fields[i] = f
			}
		default:
			err = fmt.Errorf("unrecognized dbf header")
			return
		}

		s.info = Info{
			BoundingBox: shpHeader.BoundingBox,
			NumRecords:  dbfHeader.NumRecords(),
			ShapeType:   shpHeader.ShapeType,
			Fields:      fields,
		}
	})

	return &s.info, err
}

// Scan begins reading the shp and dbf files for records. Records can be accessed from the Record method.
// An error is returned if there's a problem parsing the header of either file.
// Errors that are encountered when parsing records must be checked with the Err method.
func (s *Scanner) Scan() error {
	info, err := s.Info()
	if err != nil {
		return err
	}

	s.scanOnce.Do(func() {
		if err = s.shp.Scan(); err != nil {
			return
		} else if err = s.dbf.Scan(s.opts.dbf...); err != nil {
			return
		}

		go func() {
			defer func() {
				if err := s.shp.Err(); err != nil {
					s.setErr(fmt.Errorf("error in shp file: %w", err))
				} else if err = s.dbf.Err(); err != nil {
					s.setErr(fmt.Errorf("error in dbf file: %w", err))
				}

				close(s.recordsCh)
			}()

			for i := uint32(0); i < info.NumRecords; i++ {
				shape := s.shp.Shape()
				if err := s.shp.Err(); err != nil {
					s.setErr(fmt.Errorf("error in shp file: %w", err))
					return
				} else if shape == nil {
					s.setErr(fmt.Errorf("failed to read shape; expecting %d but have read %d", info.NumRecords, i+1))
					return
				}

				attr := s.dbf.Record()
				if err = s.dbf.Err(); err != nil {
					s.setErr(fmt.Errorf("error in dbf file: %w", err))
					return
				} else if attr == nil {
					s.setErr(fmt.Errorf("failed to read attributes; expecting %d but have read %d", info.NumRecords, i+1))
					return
				}

				s.recordsCh <- &Record{
					Shape:      shape,
					Attributes: attr,
				}
			}
		}()
	})

	return err
}

// Record returns each record found in the shp and dbf files.
// A single record consists of a shape and a set of attributes.
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
