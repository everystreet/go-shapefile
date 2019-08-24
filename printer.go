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
	for {
		rec := p.scanner.Record()
		if rec == nil {
			break
		}

		row := p.row(rec)
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
	headers := []string{"Number"}
	table.SetHeader(append(headers, p.fields...))

	for {
		rec := p.scanner.Record()
		if rec == nil {
			break
		}
		table.Append(p.row(rec))
	}

	if err := p.scanner.Err(); err != nil {
		return err
	}

	table.Render()
	return nil
}

func (p *TablePrinter) row(rec *Record) []string {
	row := make([]string, len(p.fields)+1)
	row[0] = fmt.Sprintf("%d", rec.Shape.RecordNumber())
	for i, name := range p.fields {
		if f, ok := rec.Attributes.Field(name); ok {
			row[i+1] = fmt.Sprintf("%v", f.Value())
		}
	}
	return row
}
