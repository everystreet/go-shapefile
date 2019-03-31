package shapefile

import (
	"fmt"
	"os"

	"github.com/mercatormaps/go-shapefile/dbf/dbase5"
	"github.com/olekukonko/tablewriter"
)

type TablePrinter struct {
	scanner Scannable
	fields  []string
}

type Scannable interface {
	AddOptions(...Option)
	Info() (*Info, error)
	Scan() error
	Record() *Record
	Err() error
}

func NewTablePrinter(s Scannable, fields ...string) (*TablePrinter, error) {
	info, err := s.Info()
	if err != nil {
		return nil, err
	}

	for _, f := range fields {
		if !info.Fields.Exists(f) {
			return nil, fmt.Errorf("field with name '%s' does not exist", f)
		}
	}

	s.AddOptions(FilterFields(fields...))
	return &TablePrinter{
		scanner: s,
		fields:  fields,
	}, nil
}

func (p *TablePrinter) Print() error {
	if err := p.scanner.Scan(); err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	headers := []string{"Number"}
	table.SetHeader(append(headers, p.fields...))

	for {
		rec := p.scanner.Record()
		if rec == nil {
			break
		}

		switch r := rec.Attributes.(type) {
		case *dbase5.Record:
			row := make([]string, len(p.fields)+1)
			row[0] = fmt.Sprintf("%d", rec.Shape.RecordNumber())

			for i, name := range p.fields {
				if f, ok := r.Fields[name]; ok {
					row[i+1] = fmt.Sprintf("%v", f.Value())
				}
			}

			table.Append(row)
		default:
			return fmt.Errorf("unrecognized record")
		}
	}

	if err := p.scanner.Err(); err != nil {
		return err
	}

	table.Render()
	return nil
}
