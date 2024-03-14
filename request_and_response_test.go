package http2

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_splitHeaderValue(t *testing.T) {
	tests := [...]struct {
		In           string
		ExpectHeader string
		ExpectValue  string
	}{
		{"Header: value", "Header", "value"},
		{"Header:value", "Header", "value"},
		{"Header:    value", "Header", "value"},
		{"Header: value   ", "Header", "value"},
		{"Header: :value :  ", "Header", ":value :"},
		{"Header:", "Header", ""},
		{"Header : value", "Header ", "value"}, // !
		{"Header:     ", "Header", ""},
		{"Header", "", ""},
		{":header", "", ""},
		{":", "", ""},
		{"::", "", ""},
		{"", "", ""},
	}

	for _, test := range tests {
		gotHdr, gotValue := splitHeaderValue([]byte(test.In))
		require.Equal(t, test.ExpectHeader, string(gotHdr), test.In)
		require.Equal(t, test.ExpectValue, string(gotValue), test.In)
	}
}
