package cli

import (
	"fmt"
	"io"
	"math"
	"os"
	"strings"

	svg "github.com/ajstarks/svgo"
	"github.com/alecthomas/kong"
	"github.com/everystreet/go-shapefile"
	"github.com/everystreet/go-shapefile/shp"
)

type ExportSVGCmd struct {
	Shapefiles  []string `kong:"required,name=shapefiles,short=z,help='Path to zipped shapefiles.'"`
	Destination string   `kong:"required,type=path,name=destination,short=d,help='Path to destination SVG.'"`
	Filters     []string `kong:"optional,name=filter,short=f,sep=';',help='Filter expressions.'"`
	Scale       float64  `kong:"optional,default=1,name=scale-factor,short=s,help='Scale factor.'"`
}

func (c ExportSVGCmd) Run(_ *kong.Context) error {
	filters, err := c.parseFilters()
	if err != nil {
		return err
	}

	fields := make(map[string]string)
	var shapes shp.Shapes

	for _, path := range c.Shapefiles {
		if err := func() (err error) {
			scanner, closer, err := open(path)
			defer func() {
				if closeErr := closer.Close(); closeErr != nil && err == nil {
					err = closeErr
				}
			}()

			info, err := scanner.Info()
			if err != nil {
				return err
			}

			for _, filter := range filters {
				for _, field := range info.Fields {
					if field.Name() != filter.name {
						continue
					}

					if other, ok := fields[field.Name()]; ok {
						return fmt.Errorf("filter field name '%s' is ambiguous - exists in '%s' and '%s'", field.Name(), path, other)
					}
					fields[field.Name()] = path
				}
			}

			if err := scanner.Scan(); err != nil {
				return err
			}

		Record:
			for {
				record := scanner.Record()
				if record == nil {
					break
				}

				for _, field := range record.Fields() {
					for _, filter := range filters {
						if filter.name != field.Name() {
							continue
						}

						for _, value := range filter.values {
							if field.Equal(value) {
								shapes = append(shapes, record.Shape)
								continue Record
							}
						}
					}
				}
			}

			return scanner.Err()
		}(); err != nil {
			return err
		}
	}

	for _, filter := range filters {
		if _, ok := fields[filter.name]; !ok {
			return fmt.Errorf("unrecognized field '%s' not present in any shapefile", filter.name)
		}
	}

	if len(shapes) == 0 {
		return fmt.Errorf("no records selected")
	}

	f, err := os.Create(c.Destination)
	if err != nil {
		return fmt.Errorf("failed to create file '%s': %w", c.Destination, err)
	}

	defer func() {
		if closeErr := f.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("failed to close file: %w", err)
		}
	}()

	box := shapes.BoundingBox()

	canvas := createCanvas(f, box, c.Scale)
	defer canvas.End()

	for _, shape := range shapes {
		switch v := shape.(type) {
		case shp.Polyline:
			renderPolyline(canvas, v, box, c.Scale)
		case shp.Polygon:
			renderPolygon(canvas, v, box, c.Scale)
		}

		switch shape.Type() {
		case shp.PolylineType, shp.PolygonType:
		default:
			return fmt.Errorf("record %d is of unsupported type '%s'", shape.RecordNumber(), shape.Type())
		}
	}
	return nil
}

type filter struct {
	name   string
	values []string
}

func (c ExportSVGCmd) parseFilters() ([]filter, error) {
	filters := make(map[string][]string)
	for _, str := range c.Filters {
		parts := strings.Split(str, "=")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid filter expression '%s'", str)
		}

		name := strings.TrimSpace(parts[0])
		valuesStr := strings.TrimSpace(parts[1])
		if name == "" || valuesStr == "" {
			return nil, fmt.Errorf("missing name or values from '%s'", str)
		}

		if valuesStr[0] == '[' && valuesStr[len(valuesStr)-1] == ']' {
			values := strings.Split(valuesStr[1:len(valuesStr)-1], ",")
			for i := 0; i < len(values); i++ {
				values[i] = strings.TrimSpace(values[i])
			}
			filters[name] = append(filters[name], values...)
		} else {
			filters[name] = append(filters[name], valuesStr)
		}
	}

	out := make([]filter, len(filters))
	var i int
	for name, values := range filters {
		out[i] = filter{name: name, values: values}
		i++
	}
	return out, nil
}

func open(path string) (shapefile.Scannable, io.Closer, error) {
	scanner, closer, err := OpenZip(path, nil)
	if err != nil {
		return nil, nil, err
	}

	info, err := scanner.Info()
	if err != nil {
		return nil, nil, err
	}

	switch info.ShapeType {
	case
		shp.PolylineType,
		shp.PolygonType:
		return scanner, closer, err
	default:
		return nil, nil, fmt.Errorf("unsupported shape type '%s'", info.ShapeType)
	}
}

func renderPolyline(canvas *svg.SVG, polyline shp.Polyline, box shp.BoundingBox, scale float64) {
	for _, part := range polyline.Parts {
		var xs, ys []int
		for _, point := range part {
			x, y := mapPoint(point.X, point.Y, box, scale)
			xs = append(xs, x)
			ys = append(ys, y)
		}
		canvas.Polyline(xs, ys)
	}
}

func renderPolygon(canvas *svg.SVG, polygon shp.Polygon, box shp.BoundingBox, scale float64) {
	for _, part := range polygon.Parts {
		var xs, ys []int
		for _, point := range part {
			x, y := mapPoint(point.X, point.Y, box, scale)
			xs = append(xs, x)
			ys = append(ys, y)
		}
		canvas.Polygon(xs, ys)
	}
}

func mapPoint(x, y float64, box shp.BoundingBox, scale float64) (mappedX, mappedY int) {
	return int(math.Round((x - box.MinX) * scale)),
		int(math.Round(box.MaxY*scale)) - int(math.Round(box.MinY*scale)) - int(math.Round((y-box.MinY)*scale)) - 1
}

func canvasSize(box shp.BoundingBox, scale float64) (width, height int) {
	return int(math.Round(box.MaxX*scale)) - int(math.Round(box.MinX*scale)),
		int(math.Round(box.MaxY*scale)) - int(math.Round(box.MinY*scale))
}

func createCanvas(w io.Writer, box shp.BoundingBox, scale float64) *svg.SVG {
	out := svg.New(w)
	out.Start(canvasSize(box, scale))
	return out
}
