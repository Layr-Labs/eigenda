package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommaOMatic(t *testing.T) {
	require.Equal(t, "0", CommaOMatic(0))
	require.Equal(t, "1", CommaOMatic(1))
	require.Equal(t, "12", CommaOMatic(12))
	require.Equal(t, "123", CommaOMatic(123))
	require.Equal(t, "1,234", CommaOMatic(1234))
	require.Equal(t, "12,345", CommaOMatic(12345))
	require.Equal(t, "123,456", CommaOMatic(123456))
	require.Equal(t, "1,234,567", CommaOMatic(1234567))
	require.Equal(t, "12,345,678", CommaOMatic(12345678))
	require.Equal(t, "123,456,789", CommaOMatic(123456789))
	require.Equal(t, "1,234,567,890", CommaOMatic(1234567890))

	require.Equal(t, "-1", CommaOMatic(-1))
	require.Equal(t, "-12", CommaOMatic(-12))
	require.Equal(t, "-123", CommaOMatic(-123))
	require.Equal(t, "-1,234", CommaOMatic(-1234))
	require.Equal(t, "-12,345", CommaOMatic(-12345))
	require.Equal(t, "-123,456", CommaOMatic(-123456))
	require.Equal(t, "-1,234,567", CommaOMatic(-1234567))
	require.Equal(t, "-12,345,678", CommaOMatic(-12345678))
	require.Equal(t, "-123,456,789", CommaOMatic(-123456789))
	require.Equal(t, "-1,234,567,890", CommaOMatic(-1234567890))
}
