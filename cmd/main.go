package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mercatormaps/go-shapefile"
	"github.com/mercatormaps/go-shapefile/dbf"
	"github.com/mercatormaps/go-shapefile/dbf/dbase5"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	fieldsCommand := kingpin.Command("fields", "Display fields in the attributes file.")
	fieldsDbfPath := fieldsCommand.Flag("dbf", "Attribute file (.dbf) path.").Short('d').Required().String()

	dataCommand := kingpin.Command("data", "Display shape attribute data.")
	dataZipPath := dataCommand.Flag("zip",
		"Zipped shape file (.zip). Cannot be used with --shp or --dbf.").Short('z').String()
	dataShpPath := dataCommand.Flag("shp",
		"Shape file (.shp) path. Must be used in combination with --dbf.").Short('s').String()
	dataDbfPath := dataCommand.Flag("dbf",
		"Attribute file (.dbf) path. Must be used in combination with --shp.").Short('d').String()
	dataFields := dataCommand.Flag("fields", "Only the specified field names.").Short('f').Strings()

	var err error
	switch kingpin.Parse() {
	case fieldsCommand.FullCommand():
		err = fields(*fieldsDbfPath)
	case dataCommand.FullCommand():
		switch {
		case *dataZipPath != "" && (*dataShpPath != "" || *dataDbfPath != ""):
			err = fmt.Errorf("--zip cannot be used with --shp or --dbf")
		case *dataZipPath != "":
			err = dataFromZip(*dataZipPath, dataFields)
		case *dataShpPath != "" && *dataDbfPath != "":
			err = dataFromExtracted(*dataShpPath, *dataDbfPath, dataFields)
		default:
			err = fmt.Errorf("--shp and --dbf must be used together")
		}
	default:
		fmt.Fprintf(os.Stderr, "Invalid command\n")
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}

func fields(dbfPath string) error {
	dbfFile, err := os.Open(dbfPath)
	if err != nil {
		return errors.Wrapf(err, "failed to open attributes file '%s'", dbfPath)
	}

	s := dbf.NewScanner(dbfFile)
	header, err := s.Header()
	if err != nil {
		return errors.Wrap(err, "failed to parse header")
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetHeader([]string{"Name", "Type"})

	switch h := header.(type) {
	case *dbase5.Header:
		for _, field := range h.Fields {
			typ := ""
			switch field.Type {
			case dbase5.CharacterType:
				typ = "Character"
			case dbase5.DateType:
				typ = "Date"
			case dbase5.FloatingPointType:
				typ = "Float"
			case dbase5.LogicalType:
				typ = "Logical"
			case dbase5.MemoType:
				typ = "Memo"
			case dbase5.NumericType:
				typ = "Numeric"
			default:
				typ = fmt.Sprintf("%c", field.Type)
			}

			table.Append([]string{field.Name(), typ})
		}
	default:
		return fmt.Errorf("unrecognized file type")
	}

	table.Render()
	return nil
}

func dataFromZip(path string, fields *[]string) error {
	f, err := os.Open(path)
	if err != nil {
		return errors.Wrapf(err, "failed to open zip file '%s'", path)
	}

	stat, err := f.Stat()
	if err != nil {
		return errors.Wrapf(err, "failed to stat zip file '%s'", path)
	}

	_, name := filepath.Split(path)

	s, err := shapefile.NewZipScanner(f, stat.Size(), name)
	if err != nil {
		return err
	}

	var p *shapefile.TablePrinter
	if fields != nil {
		p, err = shapefile.NewTablePrinter(s, *fields...)
	} else {
		p, err = shapefile.NewTablePrinter(s)
	}

	if err != nil {
		return err
	}
	return p.Print()
}

func dataFromExtracted(shpPath, dbfPath string, fields *[]string) error {
	shpFile, err := os.Open(shpPath)
	if err != nil {
		return errors.Wrapf(err, "failed to open shape file '%s'", shpPath)
	}

	dbfFile, err := os.Open(dbfPath)
	if err != nil {
		return errors.Wrapf(err, "failed to open attributes file '%s'", dbfPath)
	}

	s := shapefile.NewScanner(shpFile, dbfFile)

	var p *shapefile.TablePrinter
	if fields != nil {
		p, err = shapefile.NewTablePrinter(s, *fields...)
	} else {
		p, err = shapefile.NewTablePrinter(s)
	}

	if err != nil {
		return err
	}
	return p.Print()
}
