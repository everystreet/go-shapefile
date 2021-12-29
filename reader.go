package shapefile

import (
	"fmt"
	"io"
	"sync"

	"github.com/everystreet/go-shapefile/dbf"
	"github.com/everystreet/go-shapefile/shp"
)

// Reader parses a pair or shp and dbf files.
type Reader struct {
	shp         *shp.Reader
	dbf         *dbf.Reader
	infoOnce    sync.Once
	info        Info
	recordsRead uint32
}

// Info contains combined information from the pair of input files.
type Info struct {
	BoundingBox shp.BoundingBox
	NumRecords  uint32
	ShapeType   shp.ShapeType
	Fields      []dbf.FieldDesc
}

// NewReader creates a new reader for the provided shp and dbf files.
func NewReader(shpR, dbfR io.Reader, opts ...Option) *Reader {
	conf := defaultConfig()
	for _, opt := range opts {
		opt(&conf)
	}

	return &Reader{
		shp: shp.NewReader(shpR, conf.shp...),
		dbf: dbf.NewReader(dbfR, conf.dbf...),
	}
}

// Info returns combined information about the shp and dbf pair.
func (r *Reader) Info() (Info, error) {
	var err error
	r.infoOnce.Do(func() {
		shpHeader, headerErr := r.shp.Header()
		if headerErr != nil {
			err = fmt.Errorf("failed to parse shp header: %w", headerErr)
			return
		}

		dbfHeader, headerErr := r.dbf.Header()
		if headerErr != nil {
			err = fmt.Errorf("failed to parse dbf header: %w", headerErr)
			return
		}

		r.info = Info{
			BoundingBox: shpHeader.BoundingBox,
			NumRecords:  dbfHeader.NumRecords(),
			ShapeType:   shpHeader.ShapeType,
			Fields:      dbfHeader.Fields,
		}
	})
	return r.info, err
}

// Record reads the next record from the input files.
// After the last record is read, further calls to this function result in io.EOF.
func (r *Reader) Record() (*Record, error) {
	info, err := r.Info()
	if err != nil {
		return nil, err
	}

	if r.recordsRead >= info.NumRecords {
		return nil, io.EOF
	}

	shape, err := r.shp.Shape()
	if err != nil {
		return nil, fmt.Errorf("failed to read shape; expecting %d but have read %d: %w", info.NumRecords, r.info.NumRecords+1, err)
	}

	record, err := r.dbf.Record()
	if err != nil {
		return nil, fmt.Errorf("failed to read attributes; expecting %d but have read %d: %w", info.NumRecords, r.recordsRead+1, err)
	}

	r.recordsRead++

	return &Record{
		Shape:  shape,
		Record: record,
	}, nil
}
