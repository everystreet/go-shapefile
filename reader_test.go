package shapefile_test

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	shapefile "github.com/everystreet/go-shapefile"
	"github.com/stretchr/testify/require"
)

func TestReader(t *testing.T) {
	shp, err := os.Open(filepath.Join("testdata", "ne_110m_admin_0_sovereignty.shp"))
	require.NoError(t, err)

	defer func() {
		require.NoError(t, shp.Close())
	}()

	dbf, err := os.Open(filepath.Join("testdata", "ne_110m_admin_0_sovereignty.dbf"))
	require.NoError(t, err)

	defer func() {
		require.NoError(t, dbf.Close())
	}()

	r := shapefile.NewReader(shp, dbf)

	info, err := r.Info()
	require.NoError(t, err)

	var num uint32
	for {
		record, err := r.Record()
		if record == nil {
			require.ErrorIs(t, err, io.EOF)
			break
		}

		require.NoError(t, err)

		num++
	}

	require.Equal(t, info.NumRecords, num)

	record, err := r.Record()
	require.Nil(t, record)
	require.ErrorIs(t, err, io.EOF)
}
