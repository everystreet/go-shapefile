package dbf_test

import (
	"fmt"
	"testing"

	"github.com/everystreet/go-shapefile/dbf"
	"github.com/stretchr/testify/require"
)

func TestFieldDesc(t *testing.T) {
	for _, tt := range []dbf.FieldType{
		dbf.CharacterType,
		dbf.DateType,
		dbf.FloatingPointType,
		dbf.LogicalType,
		dbf.MemoType,
		dbf.NumericType,
	} {
		t.Run(fmt.Sprintf("type %c", tt), func(t *testing.T) {
			field := dbf.MakeFieldDesc(tt, "name", 12)

			buf, err := field.Encode()
			require.NoError(t, err)

			decoded, err := dbf.DecodeFieldDesc(buf)
			require.NoError(t, err)

			require.Equal(t, field.Type(), decoded.Type())
			require.Equal(t, field.Name(), decoded.Name())
			require.Equal(t, field.Length(), decoded.Length())
		})
	}
}
