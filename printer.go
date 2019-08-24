package shapefile

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

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

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

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

func (p *TablePrinter) PrettyPrint() error {
	if err := p.scanner.Scan(); err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
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

func (p *TablePrinter) header() ([]string, error) {
	info, err := p.scanner.Info()
	if err != nil {
		return nil, err
	}

	header := []string{"Number"}
	if len(p.fields) == 0 {
		for _, f := range info.Fields {
			header = append(header, f.Name())
		}
	} else {
		header = append(header, p.fields...)
	}
	return header, nil
}

func (p *TablePrinter) row(rec *Record) ([]string, error) {
	row := []string{fmt.Sprintf("%d", rec.Shape.RecordNumber())}

	// Add all fields if none specified
	if len(p.fields) == 0 {
		info, err := p.scanner.Info()
		if err != nil {
			return nil, err
		}

		row = append(row, make([]string, len(info.Fields))...)
		for i, field := range info.Fields {
			if f, ok := rec.Attributes.Field(field.Name()); ok {
				row[i+1] = fmt.Sprintf("%v", f.Value())
			}
		}
		return row, nil
	}

	// ...or just the specified fields
	for i, name := range p.fields {
		row = append(row, make([]string, len(p.fields))...)

		if f, ok := rec.Attributes.Field(name); ok {
			row[i+1] = fmt.Sprintf("%v", f.Value())
		}
	}
	return row, nil
}
