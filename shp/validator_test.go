package shp_test

import (
	"testing"

	"github.com/everystreet/go-shapefile/shp"
	"github.com/stretchr/testify/require"
)

func TestValidatePolyline(t *testing.T) {
	v, err := shp.MakeValidator(
		shp.BoundingBox{
			MinY: -90,
			MaxY: 90,
			MinX: -180,
			MaxX: 180,
		},
	)
	require.NoError(t, err)

	tests := []struct {
		name     string
		polyline shp.Polyline
		err      string
	}{
		{
			"no parts",
			shp.Polyline{},
			"must contain at least 1 part",
		},
		{
			"no edges",
			shp.Polyline{
				Parts: []shp.Part{
					{
						shp.MakePoint(0, 0),
					},
				},
			},
			"part must have at least 1 edge",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.polyline.Validate(v)
			if tt.err == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.err)
			}
		})
	}
}
