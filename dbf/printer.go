package dbf

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/everystreet/go-shapefile/dbf/dbase5"
	"github.com/olekukonko/tablewriter"
)

// TablePrinter implements a tabulated view of a shapefile.
type TablePrinter struct {
	scanner *Scanner
	fields  []string
}

// NewTablePrinter creates a TablePrinter for the supplied dbf file,
// optionally displaying only the specified field names.
func NewTablePrinter(s *Scanner, fields ...string) (*TablePrinter, error) {
	header, err := s.Header()
	if err != nil {
		return nil, err
	}

	for _, f := range fields {
		if !header.FieldExists(f) {
			return nil, fmt.Errorf("field with name '%s' does not exist", f)
		}
	}

	return &TablePrinter{
		scanner: s,
		fields:  fields,
	}, nil
}

// Print writes a tab-delimited table to the supplied destination.
func (p TablePrinter) Print(out io.Writer) error {
	if err := p.scanner.Scan(FilterFields(p.fields...)); err != nil {
		return err
	}

	w := tabwriter.NewWriter(out, 0, 0, 1, ' ', 0)

	header, err := p.header()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, strings.Join(header, "\t"))

	for {
		rec := p.scanner.Record()
		if rec == nil {
			break
		}

		row, err := p.row(rec)
		if err != nil {
			return err
		}
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}

	if err := p.scanner.Err(); err != nil {
		return err
	}

	return w.Flush()
}

// PrettyPrint writes a pretty ASCII table to the supplied destination.
func (p TablePrinter) PrettyPrint(out io.Writer) error {
	if err := p.scanner.Scan(); err != nil {
		return err
	}

	table := tablewriter.NewWriter(out)
	table.SetAutoWrapText(false)

	header, err := p.header()
	if err != nil {
		return err
	}
	table.SetHeader(header)

	for {
		rec := p.scanner.Record()
		if rec == nil {
			break
		}

		row, err := p.row(rec)
		if err != nil {
			return err
		}
		table.Append(row)
	}

	if err := p.scanner.Err(); err != nil {
		return err
	}

	table.Render()
	return nil
}

func (p TablePrinter) header() ([]string, error) {
	info, err := p.scanner.Header()
	if err != nil {
		return nil, err
	}

	switch info := info.(type) {
	case *dbase5.Header:
		if len(p.fields) == 0 {
			header := make([]string, len(info.Fields))
			for _, f := range info.Fields {
				header = append(header, f.Name())
			}
			return header, nil
		} else {
			return p.fields, nil
		}
	default:
		return []string{}, fmt.Errorf("unsupported dBase version")
	}
}

func (p TablePrinter) row(rec *Record) ([]string, error) {
	info, err := p.scanner.Header()
	if err != nil {
		return nil, err
	}

	switch info := info.(type) {
	case *dbase5.Header:
		// Add all fields if none specified
		if len(p.fields) == 0 {
			row := make([]string, len(info.Fields))
			for i, field := range info.Fields {
				if f, ok := rec.Field(field.Name()); ok {
					row[i] = fmt.Sprintf("%v", f.Value())
				}
			}
			return row, nil
		}

		// ...or just the specified fields
		row := make([]string, len(p.fields))
		for i, name := range p.fields {
			if f, ok := rec.Field(name); ok {
				row[i] = fmt.Sprintf("%v", f.Value())
			}
		}
		return row, nil
	default:
		return []string{}, fmt.Errorf("unsupported dBase version")
	}
}
