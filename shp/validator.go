package shp

import (
	"fmt"
)

type Validator struct {
	box *BoundingBox
}

func NewValidator(box *BoundingBox) *Validator {
	return &Validator{
		box: box,
	}
}

func (p *Point) Validate(v *Validator) error {
	if p.box != nil && (p.X > p.box.MaxX || p.X < p.box.MinX || p.Y > p.box.MaxY || p.Y < p.box.MinY) {
		return fmt.Errorf("shape (point %s) is not in own bounding box '%s'", p.String(), p.box.String())
	} else if p.X > v.box.MaxX || p.X < v.box.MinX || p.Y > v.box.MaxY || p.Y < v.box.MinY {
		return fmt.Errorf("shape (point %s) is not in file bounding box '%s'", p.String(), v.box.String())
	}
	return nil
}

func (p *Polyline) Validate(v *Validator) error {
	for _, part := range p.Parts {
		for _, point := range part {
			if err := point.Validate(v); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *Polygon) Validate(v *Validator) error {
	return (*Polyline)(p).Validate(v)
}
