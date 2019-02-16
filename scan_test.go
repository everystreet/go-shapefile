package shapefile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mercatormaps/go-shapefile/shp"
	"github.com/stretchr/testify/require"
)

func TestScan(t *testing.T) {
	r, err := os.Open(filepath.Join("testdata", "ne_110m_admin_0_sovereignty.shp"))
	require.NoError(t, err)

	s := shp.NewScanner(r)

	h, err := s.Header()
	require.NoError(t, err)

	require.Equal(t, &shp.Header{
		FileLength: 180400,
		Version:    1000,
		ShapeType:  5,
		MinX:       -180,
		MinY:       -90,
		MaxX:       180.00000000000006,
		MaxY:       83.64513000000001,
	}, h)
}
