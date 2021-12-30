package shapefile_test

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	shapefile "github.com/everystreet/go-shapefile"
	"github.com/stretchr/testify/require"
)

func TestZipReader(t *testing.T) {
	const filename = "ne_110m_admin_0_sovereignty.zip"

	f, err := os.Open(filepath.Join("testdata", filename))
	require.NoError(t, err)

	defer func() {
		require.NoError(t, f.Close())
	}()

	stat, err := f.Stat()
	require.NoError(t, err)

	r, err := shapefile.NewZipReader(f, stat.Size(), filename)
	require.NoError(t, err)

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
