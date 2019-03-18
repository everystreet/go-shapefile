package main

import (
	"fmt"
	"os"

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
	dataShpPath := dataCommand.Flag("shp", "Shape file (.shp) path.").Short('s').Required().String()
	dataDbfPath := dataCommand.Flag("dbf", "Attribute file (.dbf) path.").Short('d').Required().String()
	dataFields := dataCommand.Flag("fields", "Only the specified field names.").Short('f').Strings()

	var err error
	switch kingpin.Parse() {
	case fieldsCommand.FullCommand():
		err = fields(fieldsDbfPath)
	case dataCommand.FullCommand():
		err = data(dataShpPath, dataDbfPath, dataFields)
	default:
		fmt.Fprintf(os.Stderr, "Invalid command\n")
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}

func fields(dbfPath *string) error {
	dbfFile, err := os.Open(*dbfPath)
	if err != nil {
		return errors.Wrapf(err, "failed to open attributes file '%s'", *dbfPath)
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

func data(shpPath, dbfPath *string, fields *[]string) error {
	shpFile, err := os.Open(*shpPath)
	if err != nil {
		return errors.Wrapf(err, "failed to open shape file '%s'", *shpPath)
	}

	dbfFile, err := os.Open(*dbfPath)
	if err != nil {
		return errors.Wrapf(err, "failed to open attributes file '%s'", *dbfPath)
	}

	p, err := shapefile.NewTablePrinter(shapefile.NewScanner(shpFile, dbfFile), *fields...)
	if err != nil {
		return err
	}
	return p.Print()
}
