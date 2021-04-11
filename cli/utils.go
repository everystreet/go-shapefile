package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/everystreet/go-shapefile"
)

func OpenZip(path string, exclusive []string) (*shapefile.ZipScanner, closer, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open zip file '%s': %w", path, err)
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

	_, name := filepath.Split(path)

	if len(exclusive) == 0 {
		scan, err := shapefile.NewZipScanner(file, stat.Size(), name)
		return scan, close, err
	}

	scan, err := shapefile.NewZipScanner(file, stat.Size(), name, shapefile.FilterFields(exclusive...))
	return scan, close, err
}

func OpenExtracted(shpPath, dbfPath string, exclusive []string) (*shapefile.Scanner, io.Closer, error) {
	shp, err := os.Open(shpPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open shp file '%s': %w", shpPath, err)
	}

	dbf, err := os.Open(dbfPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open dbf file '%s': %w", dbfPath, err)
	}

	var close closer = func() error {
		if err := shp.Close(); err != nil {
			return fmt.Errorf("failed to close shp file: %w", err)
		} else if err := dbf.Close(); err != nil {
			return fmt.Errorf("failed to close dbf file: %w", err)
		}
		return nil
	}

	if len(exclusive) == 0 {
		return shapefile.NewScanner(shp, dbf), close, err
	}
	return shapefile.NewScanner(shp, dbf, shapefile.FilterFields(exclusive...)), close, err
}

type closer func() error

func (c closer) Close() error {
	return c()
}
