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

// ZipReader wraps a reader, providing a simple method of reading a zipped shapefile.
type ZipReader struct {
	in       *zip.Reader
	name     string
	opts     []ReaderOption
	initOnce sync.Once
	reader   *Reader
}

// NewZipReader creates a reader for the supplied zip file.
// The filename parameter should be the zip file's name (as stored on disk),
// and MUST match the names of the contained files.
func NewZipReader(r io.ReaderAt, size int64, filename string, opts ...ReaderOption) (*ZipReader, error) {
	in, err := zip.NewReader(r, size)
	if err != nil {
		return nil, err
	}

	if !strings.HasSuffix(filename, ".zip") {
		return nil, fmt.Errorf("expecting name to be *.zip")
	}

	return &ZipReader{
		in:   in,
		name: strings.TrimSuffix(filename, ".zip"),
		opts: opts,
	}, nil
}

// Info returns combined information about the shp and dbf pair.
func (r *ZipReader) Info() (Info, error) {
	if err := r.init(); err != nil {
		return Info{}, err
	}
	return r.reader.Info()
}

// Record reads the next record from the input files.
// After the last record is read, further calls to this function result in io.EOF.
func (r *ZipReader) Record() (*Record, error) {
	if _, err := r.Info(); err != nil {
		return nil, err
	}
	return r.reader.Record()
}

func (r *ZipReader) init() error {
	var err error
	r.initOnce.Do(func() {
		shpFile, dbfFile, cpgFile, filesErr := r.files()
		if filesErr != nil {
			err = filesErr
			return
		}

		shpR, openErr := shpFile.Open()
		if openErr != nil {
			err = fmt.Errorf("failed to open %s: %w", shpFile.Name, openErr)
			return
		}

		dbfR, openErr := dbfFile.Open()
		if openErr != nil {
			err = fmt.Errorf("failed to open %s: %w", dbfFile.Name, openErr)
			return
		}

		opts := make([]ReaderOption, len(r.opts))
		copy(opts, r.opts)

		if cpgFile != nil {
			decoder, readErr := readCpg(cpgFile)
			if readErr != nil {
				err = readErr
				return
			}

			opts = append(opts, CharacterDecoder(decoder))
		}

		r.reader = NewReader(shpR, dbfR, opts...)
	})
	return err
}

func (s *ZipReader) files() (shpFile, dbfFile, cpgFile *zip.File, err error) {
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
