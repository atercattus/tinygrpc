package http2

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBufReader_readFromBuf(t *testing.T) {
	src := []byte(`hello world`)

	buf := &BufReader{
		prepend: src,
	}

	dst := make([]byte, 10)

	n, err := buf.readFromBuf(dst[:4])
	require.NoError(t, err)
	require.Equal(t, 4, n)
	require.Equal(t, []byte(`hell`), dst[:n])

	n, err = buf.readFromBuf(dst[:6])
	require.NoError(t, err)
	require.Equal(t, 6, n)
	require.Equal(t, []byte(`o worl`), dst[:n])

	n, err = buf.readFromBuf(dst[:4])
	require.NoError(t, err)
	require.Equal(t, 1, n)
	require.Equal(t, []byte(`d`), dst[:n])

	n, err = buf.readFromBuf(dst[:4])
	require.NoError(t, err)
	require.Equal(t, 0, n)
}
