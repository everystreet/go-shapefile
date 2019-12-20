package shapefile

import (
	"archive/zip"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

// ZipScanner wraps Scanner, providing a simple method of reading a zipped shapefile.
type ZipScanner struct {
	opts []Option

	in   *zip.Reader
	name string

	initOnce sync.Once
	scanner  *Scanner
}

// NewZipScanner creates a ZipScanner for the supplied zip file.
// The filename parameter should be the zip file's name (as stored on disk),
// and MUST match the names of the contained shp and dbf files.
func NewZipScanner(r io.ReaderAt, size int64, filename string, opts ...Option) (*ZipScanner, error) {
	in, err := zip.NewReader(r, size)
	if err != nil {
		return nil, err
	}

	if !strings.HasSuffix(filename, ".zip") {
		return nil, fmt.Errorf("expecting name to be *.zip")
	}

	return &ZipScanner{
		opts: opts,
		in:   in,
		name: strings.TrimSuffix(filename, ".zip"),
	}, nil
}

// AddOptions allows additional options to be set after the scanner has already been created.
func (s *ZipScanner) AddOptions(opts ...Option) {
	s.opts = append(s.opts, opts...)
	if s.scanner != nil {
		s.scanner.AddOptions(s.opts...)
	}
}

// Info calls Scanner.Info().
func (s *ZipScanner) Info() (*Info, error) {
	if err := s.init(); err != nil {
		return nil, err
	}
	return s.scanner.Info()
}

// Scan calls Scanner.Scan().
func (s *ZipScanner) Scan() error {
	if err := s.init(); err != nil {
		return err
	}
	return s.scanner.Scan()
}

// Record calls Scanner.Record().
func (s *ZipScanner) Record() *Record {
	if s.scanner == nil {
		return nil
	}
	return s.scanner.Record()
}

// Err returns the first error encountered when parsing records.
// It should be called after calling the Shape method for the last time.
func (s *ZipScanner) Err() error {
	if s.scanner == nil {
		return nil
	}
	return s.scanner.Err()
}

func (s *ZipScanner) init() error {
	var err error

	s.initOnce.Do(func() {
		var shpFile, dbfFile, cfgFile *zip.File
		shpFile, dbfFile, cfgFile, err = s.files()
		if err != nil {
			return
		}
		_ = cfgFile

		var shpR, dbfR io.ReadCloser
		shpR, err = shpFile.Open()
		if err != nil {
			err = errors.Wrapf(err, "failed to open %s", shpFile.Name)
			return
		}

		dbfR, err = dbfFile.Open()
		if err != nil {
			err = errors.Wrapf(err, "failed to open %s", dbfFile.Name)
			return
		}

		s.scanner = NewScanner(shpR, dbfR, s.opts...)
	})

	return err
}

func (s *ZipScanner) files() (shpFile, dbfFile, cpgFile *zip.File, err error) {
	if s.name != "" {
		for _, f := range s.in.File {
			switch f.Name {
			case s.name + ".shp":
				shpFile = f
			case s.name + ".dbf":
				dbfFile = f
			case s.name + ".cpg":
				cpgFile = f
			}
		}
	} else {
		for _, f := range s.in.File {
			switch {
			case strings.HasSuffix(f.Name, ".shp"):
				if shpFile != nil {
					err = fmt.Errorf("found multiple .shp files")
					return
				}
				shpFile = f
			case strings.HasSuffix(f.Name, ".dbf"):
				if dbfFile != nil {
					err = fmt.Errorf("found multiple .dbf files")
					return
				}
				dbfFile = f
			case strings.HasSuffix(f.Name, "cpg"):
				if cpgFile != nil {
					err = fmt.Errorf("found multiple .cpg files")
					return
				}
				cpgFile = f
			}
		}
	}

	if shpFile == nil {
		err = fmt.Errorf("missing .shp file")
	} else if dbfFile == nil {
		err = fmt.Errorf("missing .dbf file")
	}
	return
}
