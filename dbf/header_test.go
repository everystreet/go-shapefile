package dbf_test

import (
	"bytes"
	"testing"

	"github.com/everystreet/go-shapefile/dbf"
	"github.com/stretchr/testify/require"
)

func TestHeader(t *testing.T) {
	header, err := dbf.MakeHeader(
		[]dbf.FieldDesc{
			dbf.MakeFieldDesc(dbf.CharacterType, "field-1", 8),
			dbf.MakeFieldDesc(dbf.FloatingPointType, "field-2", 32),
		},
		dbf.DBaseLevel5, 1024, 128,
	)
	require.NoError(t, err)

	var buf bytes.Buffer
	err = header.Encode(&buf)
	require.NoError(t, err)

	decoded, err := dbf.DecodeHeader(&buf)
	require.NoError(t, err)

	require.Equal(t, header.Fields(), decoded.Fields())
	require.Equal(t, header.Version(), decoded.Version())
	require.Equal(t, header.Len(), decoded.Len())
	require.Equal(t, header.RecordLen(), decoded.RecordLen())
	require.Equal(t, header.NumRecords(), decoded.NumRecords())
}
