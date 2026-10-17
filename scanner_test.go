package shapefile_test

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/everystreet/go-shapefile"
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

	require.NoError(t, shp.Close())
	require.NoError(t, dbf.Close())
}

func TestScannerReadErrors(t *testing.T) {
	shpData, err := ioutil.ReadFile(filepath.Join("testdata", "ne_110m_admin_0_sovereignty.shp"))
	require.NoError(t, err)

	dbfData, err := ioutil.ReadFile(filepath.Join("testdata", "ne_110m_admin_0_sovereignty.dbf"))
	require.NoError(t, err)

	// number of valid records that precede the error
	const valid = 10

	// offsets of the first record after the valid ones
	shpPos := 100
	for i := 0; i < valid; i++ {
		shpPos += 8 + int(binary.BigEndian.Uint32(shpData[shpPos+4:shpPos+8]))*2
	}
	dbfPos := int(binary.LittleEndian.Uint16(dbfData[8:10])) + valid*int(binary.LittleEndian.Uint16(dbfData[10:12]))

	t.Run("shp record", func(t *testing.T) {
		// give the record a polyline shape type (3), which differs from the file's polygon type (5)
		shp := append([]byte{}, shpData...)
		binary.LittleEndian.PutUint32(shp[shpPos+8:shpPos+12], 3)

		num, err := scanRecords(shapefile.NewScanner(bytes.NewReader(shp), bytes.NewReader(dbfData)))
		require.Error(t, err)
		require.Contains(t, err.Error(), "error in shp file")
		require.Equal(t, uint32(valid), num)
	})

	t.Run("shp shape", func(t *testing.T) {
		// give the shape more points than its record has room for
		shp := append([]byte{}, shpData...)
		binary.LittleEndian.PutUint32(shp[shpPos+48:shpPos+52], 1<<20)

		num, err := scanRecords(shapefile.NewScanner(bytes.NewReader(shp), bytes.NewReader(dbfData)))
		require.Error(t, err)
		require.Contains(t, err.Error(), "error in shp file")
		require.Equal(t, uint32(valid), num)
	})

	t.Run("shp trailing", func(t *testing.T) {
		// append a record header with a polyline shape type (3) after the last record
		shp := append(append([]byte{}, shpData...), 0, 0, 0, 172, 0, 0, 0, 2, 3, 0, 0, 0)

		num, err := scanRecords(shapefile.NewScanner(bytes.NewReader(shp), bytes.NewReader(dbfData)))
		require.Error(t, err)
		require.Contains(t, err.Error(), "error in shp file")
		require.Equal(t, uint32(171), num)
	})

	t.Run("dbf record", func(t *testing.T) {
		// replace the deletion flag with an invalid value
		dbf := append([]byte{}, dbfData...)
		dbf[dbfPos] = 0

		num, err := scanRecords(shapefile.NewScanner(bytes.NewReader(shpData), bytes.NewReader(dbf)))
		require.Error(t, err)
		require.Contains(t, err.Error(), "error in dbf file")
		require.Equal(t, uint32(valid), num)
	})

	t.Run("dbf terminator", func(t *testing.T) {
		// replace the file terminator that follows the last record
		dbf := append([]byte{}, dbfData...)
		dbf[len(dbf)-1] = 0

		num, err := scanRecords(shapefile.NewScanner(bytes.NewReader(shpData), bytes.NewReader(dbf)))
		require.Error(t, err)
		require.Contains(t, err.Error(), "error in dbf file")
		require.Equal(t, uint32(171), num)
	})
}

func scanRecords(s *shapefile.Scanner) (uint32, error) {
	if err := s.Scan(); err != nil {
		return 0, err
	}

	var num uint32
	for s.Record() != nil {
		num++
	}
	return num, s.Err()
}
