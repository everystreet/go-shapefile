package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/everystreet/go-shapefile"
)

type Flags struct {
	Zip string `kong:"optional,name=zip,short=z,type=existingfile,help='Path to zipped (.zip) shapefile. Not to be used in conjunction with --shp or --dbf.'"`
	Shp string `kong:"optional,name=shp,type=existingfile,help='Path to shape file (.shp). Must be used in conjunction with --dbf.'"`
	Dbf string `kong:"optional,name=dbf,type=existingfile,help='Path to attribute file (.dbf). Must be used in conjunction with --shp.'"`
}

func (f Flags) OpenAllFields() (shapefile.Scannable, io.Closer, error) {
	return f.open(nil)
}

func (f Flags) OpenFilteredFields(fields []string) (shapefile.Scannable, io.Closer, error) {
	return f.open(fields)
}

func (f Flags) open(fields []string) (shapefile.Scannable, io.Closer, error) {
	if f.Zip != "" {
		if f.Shp != "" || f.Dbf != "" {
			return nil, nil, fmt.Errorf("--zip cannot be used with --shp or --dbf")
		}
		return f.openZip(fields)
	} else if f.Shp != "" || f.Dbf != "" {
		if f.Zip != "" {
			return nil, nil, fmt.Errorf("--shp and --dbf cannot be used with --zip")
		}
		return f.openExtracted(fields)
	}
	return nil, nil, fmt.Errorf("missing --zip or --shp & --dbf combination")
}

func (f Flags) openZip(fields []string) (*shapefile.ZipScanner, closer, error) {
	file, err := os.Open(f.Zip)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open zip file '%s': %w", f.Zip, err)
	}

	close := func() error {
		if err := file.Close(); err != nil {
			return fmt.Errorf("failed to close zip file: %w", err)
		}
		return nil
	}

	stat, err := file.Stat()
	if err != nil {
		return nil, close, fmt.Errorf("failed to stat zip file: %w", err)
	}

	_, name := filepath.Split(f.Zip)

	if len(fields) == 0 {
		scan, err := shapefile.NewZipScanner(file, stat.Size(), name)
		return scan, close, err
	}

	scan, err := shapefile.NewZipScanner(file, stat.Size(), name, shapefile.FilterFields(fields...))
	return scan, close, err
}

func (f Flags) openExtracted(fields []string) (*shapefile.Scanner, closer, error) {
	shp, err := os.Open(f.Shp)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open shp file '%s': %w", f.Shp, err)
	}

	dbf, err := os.Open(f.Dbf)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open dbf file '%s': %w", f.Dbf, err)
	}

	close := func() error {
		if err := shp.Close(); err != nil {
			return fmt.Errorf("failed to close shp file: %w", err)
		} else if err := dbf.Close(); err != nil {
			return fmt.Errorf("failed to close dbf file: %w", err)
		}
		return nil
	}

	if len(fields) == 0 {
		return shapefile.NewScanner(shp, dbf), close, err
	}
	return shapefile.NewScanner(shp, dbf, shapefile.FilterFields(fields...)), close, err
}

type closer func() error

func (c closer) Close() error {
	return c()
}
