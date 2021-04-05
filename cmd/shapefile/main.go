package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/everystreet/go-shapefile"
	"github.com/everystreet/go-shapefile/dbf"
	"github.com/everystreet/go-shapefile/dbf/dbase5"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	readCommand := kingpin.Command("read", "Display shapefile data.")
	readZipPath := readCommand.Flag("zip",
		"Zipped shape file (.zip). Cannot be used with --shp or --dbf.").Short('z').String()
	readShpPath := readCommand.Flag("shp",
		"Shape file (.shp) path. Must be used in combination with --dbf.").Short('s').String()
	readDbfPath := readCommand.Flag("dbf",
		"Attribute file (.dbf) path. Must be used in combination with --shp.").Short('d').String()
	readFields := readCommand.Flag("fields", "Only the specified field names.").Short('f').Strings()
	readListFields := readCommand.Flag("list-fields", "List fields only - no data.").Bool()
	pretty := readCommand.Flag("pretty", "Enable pretty-printing.").Short('p').Bool()

	var err error
	switch kingpin.Parse() {
	case readCommand.FullCommand():
		switch {
		case *readZipPath == "" && *readShpPath == "" && *readDbfPath == "":
			err = fmt.Errorf("no data source specified")
		case *readZipPath != "" && (*readShpPath != "" || *readDbfPath != ""):
			err = fmt.Errorf("--zip cannot be used with --shp or --dbf")
		case *readZipPath != "":
			err = dataFromZip(*readZipPath, readFields, *readListFields, *pretty)
		case *readListFields == false && (*readShpPath == "" || *readDbfPath == ""):
			err = fmt.Errorf("--shp and --dbf must be used together")
		case *readListFields == true && *readDbfPath != "":
			err = fieldsFromExtracted(*readDbfPath, *pretty)
		case *readShpPath != "" && *readDbfPath != "":
			err = dataFromExtracted(*readShpPath, *readDbfPath, readFields, *pretty)
		default:
			err = fmt.Errorf("unspecified source")
		}
	default:
		fmt.Fprintf(os.Stderr, "Invalid command\n")
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}

func dataFromZip(path string, fields *[]string, meta, pretty bool) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open zip file '%s': %w", path, err)
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat zip file '%s': %w", path, err)
	}

	_, name := filepath.Split(path)

	var s *shapefile.ZipScanner
	if fields == nil || len(*fields) == 0 {
		if s, err = shapefile.NewZipScanner(f, stat.Size(), name); err != nil {
			return err
		}
	} else {
		if s, err = shapefile.NewZipScanner(f, stat.Size(), name, shapefile.FilterFields(*fields...)); err != nil {
			return err
		}
	}

	if meta {
		info, err := s.Info()
		if err != nil {
			return err
		}

		if pretty {
			return fieldsPrettyTable(info.Fields)
		}
		return fieldsTable(info.Fields)
	}
	return dataTable(s, fields, pretty)
}

func dataFromExtracted(shpPath, dbfPath string, fields *[]string, pretty bool) error {
	shpFile, err := os.Open(shpPath)
	if err != nil {
		return fmt.Errorf("failed to open shape file '%s': %w", shpPath, err)
	}
	defer shpFile.Close()

	dbfFile, err := os.Open(dbfPath)
	if err != nil {
		return fmt.Errorf("failed to open attributes file '%s': %w", dbfPath, err)
	}
	defer dbfFile.Close()

	var s *shapefile.Scanner
	if fields == nil || len(*fields) == 0 {
		s = shapefile.NewScanner(shpFile, dbfFile)
	} else {
		s = shapefile.NewScanner(shpFile, dbfFile, shapefile.FilterFields(*fields...))
	}
	return dataTable(s, fields, pretty)
}

func fieldsFromExtracted(dbfPath string, pretty bool) error {
	dbfFile, err := os.Open(dbfPath)
	if err != nil {
		return fmt.Errorf("failed to open attributes file '%s': %w", dbfPath, err)
	}
	defer dbfFile.Close()

	s := dbf.NewScanner(dbfFile)
	header, err := s.Header()
	if err != nil {
		return fmt.Errorf("failed to parse dbf header: %w", err)
	}

	switch h := header.(type) {
	case *dbase5.Header:
		fields := make([]shapefile.FieldDesc, len(h.Fields))
		for i, f := range h.Fields {
			fields[i] = f
		}

		if pretty {
			return fieldsPrettyTable(fields)
		}
		return fieldsTable(fields)
	default:
		return fmt.Errorf("unrecognized file type")
	}
}

func dataTable(s shapefile.Scannable, fields *[]string, pretty bool) error {
	var p *shapefile.TablePrinter
	var err error
	if fields != nil {
		p, err = shapefile.NewTablePrinter(s, *fields...)
	} else {
		p, err = shapefile.NewTablePrinter(s)
	}

	if err != nil {
		return err
	}

	if pretty {
		return p.PrettyPrint(os.Stdout)
	}
	return p.Print(os.Stdout)
}

func fieldsTable(fields shapefile.FieldDescList) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	for _, field := range fields {
		f, err := fieldRow(field)
		if err != nil {
			return err
		}

		fmt.Fprintln(w, strings.Join(f, "\t"))
	}
	return w.Flush()
}

func fieldsPrettyTable(fields shapefile.FieldDescList) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetHeader([]string{"Name", "Type"})

	for _, field := range fields {
		f, err := fieldRow(field)
		if err != nil {
			return err
		}
		table.Append(f)
	}

	table.Render()
	return nil
}

func fieldRow(field shapefile.FieldDesc) ([]string, error) {
	switch f := field.(type) {
	case *dbase5.FieldDesc:
		typ := ""
		switch f.Type {
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
			typ = fmt.Sprintf("%c", f.Type)
		}

		return []string{field.Name(), typ}, nil
	default:
		return nil, fmt.Errorf("unrecognized file type")
	}
}
