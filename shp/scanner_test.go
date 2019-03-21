package shp_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mercatormaps/go-shapefile/shp"
	"github.com/stretchr/testify/require"
)

func TestScan(t *testing.T) {
	r, err := os.Open(filepath.Join("../testdata", "ne_110m_admin_0_sovereignty.shp"))
	require.NoError(t, err)

	s := shp.NewScanner(r)

	h, err := s.Header()
	require.NoError(t, err)

	require.Equal(t, &shp.Header{
		FileLength: 180400,
		Version:    1000,
		ShapeType:  5,
		BoundingBox: shp.BoundingBox{
			MinX: -180,
			MinY: -90,
			MaxX: 180.00000000000006,
			MaxY: 83.64513000000001,
		},
	}, h)

	err = s.Scan()
	require.NoError(t, err)

	num := 0
	points := 0
	for {
		shape := s.Shape()
		if shape == nil {
			break
		}
		num++

		switch s := shape.(type) {
		case *shp.Polygon:
			for _, p := range s.Parts {
				points += len(p)
			}
		}
	}

	require.Equal(t, 10641, points)

	require.NoError(t, s.Err())
	require.Equal(t, 171, num)

	require.NoError(t, r.Close())
}
