package shapefile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mercatormaps/go-shapefile"
	"github.com/stretchr/testify/require"
)

func TestScanner(t *testing.T) {
	shp, err := os.Open(filepath.Join("testdata", "ne_110m_admin_0_sovereignty.shp"))
	require.NoError(t, err)

	dbf, err := os.Open(filepath.Join("testdata", "ne_110m_admin_0_sovereignty.dbf"))
	require.NoError(t, err)

	s := shapefile.NewScanner(shp, dbf)

	info, err := s.Info()
	require.NoError(t, err)

	err = s.Scan()
	require.NoError(t, err)

	var num uint32
	for {
		rec := s.Record()
		if rec == nil {
			break
		}
		num++
	}

	require.NoError(t, s.Err())
	require.Equal(t, info.NumRecords, num)
}
