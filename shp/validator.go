package shp

import (
	"fmt"

	"github.com/golang/geo/r1"
	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
)

// Validator is used to validate shapes inside a shp file.
type Validator struct {
	box s2.Rect
}

// MakeValidator creates a new Validator based on the constraints of a particular shp file.
func MakeValidator(box BoundingBox) (Validator, error) {
	rect, err := boxToRect(box)
	if err != nil {
		return Validator{}, err
	}

	return Validator{
		box: rect,
	}, nil
}

// Validate the Point by checking that it is within the shp file bounding box.
func (p Point) Validate(v Validator) error {
	ll := pointToLatLng(p)

	if p.box != nil {
		box, err := boxToRect(*p.box)
		if err != nil {
			return err
		}

		if !box.ContainsLatLng(ll) {
			return fmt.Errorf("point %s is not in own bounding box '%s'", ll.String(), box.String())
		}
	}

	if !v.box.ContainsLatLng(ll) {
		return fmt.Errorf("point '%s' is not in file bounding box '%s'", ll.String(), v.box.String())
	}
	return nil
}

// Validate the Polyline.
func (p Polyline) Validate(v Validator) error {
	if len(p.Parts) < 1 {
		return fmt.Errorf("must contain at least 1 part")
	}

	for _, part := range p.Parts {
		latlngs := make([]s2.LatLng, len(part))
		for i, point := range part {
			if err := point.Validate(v); err != nil {
				return err
			}
			latlngs[i] = pointToLatLng(point)
		}

		line := s2.PolylineFromLatLngs(latlngs)
		if line.NumEdges() < 1 {
			return fmt.Errorf("part must have at least 1 edge")
		}
	}
	return nil
}

// Validate the Polygon.
func (p Polygon) Validate(v Validator) error {
	return (Polyline)(p).Validate(v)
}

func boxToRect(box BoundingBox) (s2.Rect, error) {
	tl := s2.LatLngFromDegrees(box.MaxY, box.MinX)
	br := s2.LatLngFromDegrees(box.MinY, box.MaxX)

	rect := s2.Rect{
		Lat: r1.Interval{Lo: br.Lat.Radians(), Hi: tl.Lat.Radians()},
		Lng: s1.Interval{Lo: tl.Lng.Radians(), Hi: br.Lng.Radians()},
	}

	if !rect.IsValid() {
		return s2.Rect{}, fmt.Errorf("invalid box")
	}
	return rect, nil
}

func pointToLatLng(p Point) s2.LatLng {
	return s2.LatLngFromDegrees(p.Y, p.X)
}
