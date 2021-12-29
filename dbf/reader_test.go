package dbf_test

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/everystreet/go-shapefile/dbf"
	"github.com/stretchr/testify/require"
)

func TestReader(t *testing.T) {
	f, err := os.Open(filepath.Join("../testdata", "ne_110m_admin_0_sovereignty.dbf"))
	require.NoError(t, err)

	defer func() {
		require.NoError(t, f.Close())
	}()

	r := dbf.NewReader(f)

	h, err := r.Header()
	require.NoError(t, err)
	require.Equal(t, dbf.DBaseLevel5, h.Version())
	require.Equal(t, 1869, int(h.RecordLen()))
	require.Equal(t, 171, int(h.NumRecords()))
	require.Len(t, h.Fields(), 94)

	var num uint32
	for {
		record, err := r.Record()
		if record == nil {
			require.ErrorIs(t, err, io.EOF)
			break
		}

		require.NoError(t, err)
		require.False(t, record.Deleted())
		num++
	}

	require.Equal(t, h.NumRecords(), num)

	record, err := r.Record()
	require.Nil(t, record)
	require.ErrorIs(t, err, io.EOF)
}

func TestReader2(t *testing.T) {
	f, err := os.Open(filepath.Join("../testdata", "water_main_dist.dbf"))
	require.NoError(t, err)

	defer func() {
		require.NoError(t, f.Close())
	}()

	r := dbf.NewReader(f)

	h, err := r.Header()
	require.NoError(t, err)
	require.Equal(t, dbf.DBaseLevel5, h.Version())
	require.Equal(t, 170, int(h.RecordLen()))
	require.Equal(t, 2274, int(h.NumRecords()))
	require.Len(t, h.Fields(), 17)

	var num uint32
	for {
		record, err := r.Record()
		if record == nil {
			require.ErrorIs(t, err, io.EOF)
			break
		}

		require.NoError(t, err)
		require.False(t, record.Deleted())
		num++
	}

	require.Equal(t, h.NumRecords(), num)

	record, err := r.Record()
	require.Nil(t, record)
	require.ErrorIs(t, err, io.EOF)
}
