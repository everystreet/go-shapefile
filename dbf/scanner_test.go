package dbf_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/everystreet/go-shapefile/dbf"
	"github.com/everystreet/go-shapefile/dbf/dbase5"
	"github.com/stretchr/testify/require"
)

func TestScanner(t *testing.T) {
	r, err := os.Open(filepath.Join("../testdata", "ne_110m_admin_0_sovereignty.dbf"))
	require.NoError(t, err)

	s := dbf.NewScanner(r)

	v, err := s.Version()
	require.NoError(t, err)
	require.Equal(t, dbf.DBaseLevel5, v)

	h, err := s.Header()
	require.NoError(t, err)
	require.IsType(t, &dbase5.Header{}, h)
	require.Equal(t, uint16(1869), h.RecordLen())
	require.Equal(t, uint32(171), h.NumRecords())
	require.Len(t, h.(*dbase5.Header).Fields, 94)

	err = s.Scan()
	require.NoError(t, err)

	var num uint32
	for {
		rec := s.Record()
		if rec == nil {
			break
		}
		num++

		require.False(t, rec.Deleted())
	}

	require.NoError(t, s.Err())
	require.Equal(t, h.NumRecords(), num)

	require.NoError(t, r.Close())
}
