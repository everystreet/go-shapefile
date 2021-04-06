package shp_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/everystreet/go-shapefile/shp"
	"github.com/stretchr/testify/require"
)

func TestScan(t *testing.T) {
	r, err := os.Open(filepath.Join("../testdata", "ne_110m_admin_0_sovereignty.shp"))
	require.NoError(t, err)

	s := shp.NewScanner(r, shp.PointPrecision(6))

	h, err := s.Header()
	require.NoError(t, err)

	require.Equal(t, shp.Header{
		FileLength: 180400,
		Version:    1000,
		ShapeType:  shp.PolygonType,
		BoundingBox: shp.BoundingBox{
			MinX: -180.000000,
			MinY: -90.000000,
			MaxX: 180.000000,
			MaxY: 83.645130,
		},
	}, h)

	err = s.Scan()
	require.NoError(t, err)

	v, err := s.Validator()
	require.NoError(t, err)

	shapes := 0
	points := 0
	for {
		shape := s.Shape()
		if shape == nil {
			break
		}
		shapes++

		require.Equal(t, h.ShapeType, shape.Type())
		require.NoError(t, shape.Validate(v))

		switch s := shape.(type) {
		case shp.Polygon:
			for _, p := range s.Parts {
				points += len(p)
			}
		}
	}

	require.Equal(t, 10641, points)

	require.NoError(t, s.Err())
	require.Equal(t, 171, shapes)

	require.NoError(t, r.Close())
}
