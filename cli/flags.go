package cli

import (
	"fmt"
	"io"

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
		return OpenZip(f.Zip, fields)
	} else if f.Shp != "" || f.Dbf != "" {
		if f.Zip != "" {
			return nil, nil, fmt.Errorf("--shp and --dbf cannot be used with --zip")
		}
		return OpenExtracted(f.Shp, f.Dbf, fields)
	}
	return nil, nil, fmt.Errorf("missing --zip or --shp & --dbf combination")
}
