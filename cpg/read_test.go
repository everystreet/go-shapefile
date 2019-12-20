package cpg_test

import (
	"bytes"
	"testing"

	"github.com/everystreet/go-shapefile/cpg"
	"github.com/stretchr/testify/require"
)

func TestRead(t *testing.T) {
	tests := []struct {
		in  string
		out cpg.CharacterEncoding
	}{
		{"ASCII", cpg.EncodingASCII},
		{"UTF-8", cpg.EncodingUTF8},
		{"UTF8", cpg.EncodingUTF8},
		{"iicsa", cpg.EncodingUnknown},
	}

	for _, tt := range tests {
		in := "\n\t   " + tt.in + "   \t\n"
		out, err := cpg.Read(bytes.NewBuffer([]byte(in)))
		require.NoError(t, err)
		require.Equal(t, tt.out, out)
	}
}
