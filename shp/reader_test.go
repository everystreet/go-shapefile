package shp_test

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/everystreet/go-shapefile/shp"
	"github.com/stretchr/testify/require"
)

func TestReader(t *testing.T) {
	f, err := os.Open(filepath.Join("../testdata", "ne_110m_admin_0_sovereignty.shp"))
	require.NoError(t, err)

	defer func() {
		require.NoError(t, f.Close())
	}()

	r := shp.NewReader(f, shp.PointPrecision(6))

	h, err := r.Header()
	require.NoError(t, err)

	require.Equal(t, 180400, int(h.FileLength()))
	require.Equal(t, 1000, int(h.Version()))
	require.Equal(t, shp.PolygonType, h.ShapeType())
	require.Equal(t, shp.MakeBoundingBox(-180.000000, -90.000000, 180.000000, 83.645130), h.BoundingBox())

	v, err := r.Validator()
	require.NoError(t, err)

	shapes := 0
	points := 0
	for {
		shape, err := r.Shape()
		if shape == nil {
			require.ErrorIs(t, err, io.EOF)
			break
		}

		require.NoError(t, err)

		require.Equal(t, h.ShapeType(), shape.Type())
		require.NoError(t, shape.Validate(v))

		switch s := shape.(type) {
		case shp.Polygon:
			for _, p := range s.Parts() {
				points += len(p)
			}
		}

		shapes++
	}

	require.Equal(t, 10641, points)

	record, err := r.Shape()
	require.Nil(t, record)
	require.ErrorIs(t, err, io.EOF)
}
