package shapefile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mercatormaps/go-shapefile"
	"github.com/stretchr/testify/require"
)

func TestScanZip(t *testing.T) {
	const filename = "ne_110m_admin_0_sovereignty.zip"
	r, err := os.Open(filepath.Join("testdata", filename))
	require.NoError(t, err)

	stat, err := r.Stat()
	require.NoError(t, err)

	s, err := shapefile.NewZipScanner(r, stat.Size(), filename)
	require.NoError(t, err)

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

	require.NoError(t, r.Close())
}
