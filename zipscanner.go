package shapefile

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"strings"
	"sync"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
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
		var shpFile, dbfFile, cpgFile *zip.File
		shpFile, dbfFile, cpgFile, err = s.files()
		if err != nil {
			return
		}

		var shpR, dbfR io.Reader
		shpR, err = shpFile.Open()
		if err != nil {
			err = fmt.Errorf("failed to open %s: %w", shpFile.Name, err)
			return
		}

		dbfR, err = dbfFile.Open()
		if err != nil {
			err = fmt.Errorf("failed to open %s: %w", dbfFile.Name, err)
			return
		}

		opts := make([]Option, len(s.opts))
		copy(opts, s.opts)

		if cpgFile != nil {
			var dec *encoding.Decoder
			dec, err = readCpg(cpgFile)
			if err != nil {
				return
			}
			opts = append(opts, CharacterDecoder(dec))
		}

		s.scanner = NewScanner(shpR, dbfR, opts...)
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

func readCpg(f *zip.File) (*encoding.Decoder, error) {
	r, err := f.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open cpg file: %w", err)
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		str := strings.TrimSpace(scanner.Text())
		if len(str) == 0 {
			continue
		}

		enc, _ := charset.Lookup(str)
		if enc == nil {
			return nil, fmt.Errorf("unknown charset '%s'", str)
		}
		return enc.NewDecoder(), nil
	}
	return nil, fmt.Errorf("missing charset")
}
