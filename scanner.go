package shapefile

import (
	"fmt"
	"io"
	"sync"

	"github.com/mercatormaps/go-shapefile/dbf"
	"github.com/mercatormaps/go-shapefile/shp"
	"github.com/pkg/errors"
)

type Scanner struct {
	shp *shp.Scanner
	dbf *dbf.Scanner

	infoOnce sync.Once
	info     Info

	scanOnce  sync.Once
	recordsCh chan *Record

	errOnce sync.Once
	err     error
}

type Info struct {
	BoundingBox shp.BoundingBox
	NumRecords  uint32
	ShapeType   shp.ShapeType
}

type Record struct {
	Shape      shp.Shape
	Attributes dbf.Record
}

func NewScanner(shpR, dbfR io.Reader) *Scanner {
	return &Scanner{
		shp:       shp.NewScanner(shpR),
		dbf:       dbf.NewScanner(dbfR),
		recordsCh: make(chan *Record),
	}
}

func (s *Scanner) Info() (*Info, error) {
	var err error

	s.infoOnce.Do(func() {
		var shpHeader *shp.Header
		if shpHeader, err = s.shp.Header(); err != nil {
			err = errors.Wrap(err, "failed to parse shp header")
			return
		}

		var dbfHeader dbf.Header
		if dbfHeader, err = s.dbf.Header(); err != nil {
			err = errors.Wrap(err, "failed to parse dbf header")
			return
		}

		s.info = Info{
			BoundingBox: shpHeader.BoundingBox,
			NumRecords:  dbfHeader.NumRecords(),
			ShapeType:   shpHeader.ShapeType,
		}
	})

	return &s.info, err
}

func (s *Scanner) Scan() error {
	info, err := s.Info()
	if err != nil {
		return err
	}

	s.scanOnce.Do(func() {
		if err = s.shp.Scan(); err != nil {
			return
		} else if err = s.dbf.Scan(); err != nil {
			return
		}

		go func() {
			defer func() {
				if err := s.shp.Err(); err != nil {
					s.setErr(errors.Wrap(err, "error in shp file"))
				} else if err = s.dbf.Err(); err != nil {
					s.setErr(errors.Wrap(err, "error in dbf file"))
				}

				close(s.recordsCh)
			}()

			for i := uint32(0); i < info.NumRecords; i++ {
				shape := s.shp.Shape()
				if shape == nil {
					s.setErr(fmt.Errorf("failed to read shape; expecting %d but have read %d", info.NumRecords, i+1))
					return
				}

				attr := s.dbf.Record()
				if attr == nil {
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

func (s *Scanner) setErr(err error) {
	s.errOnce.Do(func() {
		s.err = err
	})
}
