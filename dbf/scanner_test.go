package dbf_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/mercatormaps/go-shapefile/dbf"
	"github.com/mercatormaps/go-shapefile/dbf/dbase5"
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
	require.Equal(t, uint32(171), h.NumRecords())
	require.Len(t, h.(*dbase5.Header).Fields, 94)

	for i, f := range h.(*dbase5.Header).Fields {
		fmt.Printf("%d) %s: %c\n", i, f.Name, f.Type)
	}
}
