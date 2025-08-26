package enforce

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTrue(t *testing.T) {
	True(true, "This should not panic")

	require.Panics(t, func() {
		True(false, "This should panic")
	})
}

func TestFalse(t *testing.T) {
	False(false, "This should not panic")

	require.Panics(t, func() {
		False(true, "This should panic")
	})
}

func TestEquals(t *testing.T) {
	Equals(1, 1, "This should not panic")

	require.Panics(t, func() {
		Equals(1, 2, "This should panic")
	})
}

func TestNotEquals(t *testing.T) {
	NotEquals(1, 2, "This should not panic")

	require.Panics(t, func() {
		NotEquals(1, 1, "This should panic")
	})
}

func TestGreaterThan(t *testing.T) {
	GreaterThan(2, 1, "This should not panic")

	require.Panics(t, func() {
		GreaterThan(1, 2, "This should panic")
	})
	require.Panics(t, func() {
		GreaterThan(2, 2, "This should panic")
	})
}

func TestGreaterThanOrEqual(t *testing.T) {
	GreaterThanOrEqual(2, 1, "This should not panic")
	GreaterThanOrEqual(2, 2, "This should not panic")

	require.Panics(t, func() {
		GreaterThanOrEqual(1, 2, "This should panic")
	})
}

func TestLessThan(t *testing.T) {
	LessThan(1, 2, "This should not panic")

	require.Panics(t, func() {
		LessThan(2, 1, "This should panic")
	})
	require.Panics(t, func() {
		LessThan(2, 2, "This should panic")
	})
}

func TestLessThanOrEqual(t *testing.T) {
	LessThanOrEqual(1, 2, "This should not panic")
	LessThanOrEqual(2, 2, "This should not panic")

	require.Panics(t, func() {
		LessThanOrEqual(2, 1, "This should panic")
	})
}

func TestNotNil(t *testing.T) {
	notNilValue := "not nil"
	NotNil(&notNilValue, "This should not panic")

	require.Panics(t, func() {
		var nilValue *string
		NotNil(nilValue, "This should panic")
	})
}

func TestNil(t *testing.T) {
	nilValue := (*string)(nil)
	Nil(nilValue, "This should not panic")

	require.Panics(t, func() {
		notNilValue := "not nil"
		Nil(&notNilValue, "This should panic")
	})
}

func TestNotEmptyList(t *testing.T) {
	notEmptyList := []int{1, 2, 3}
	NotEmptyList(notEmptyList, "This should not panic")

	require.Panics(t, func() {
		emptyList := []int{}
		NotEmptyList(emptyList, "This should panic")
	})
}

func TestNotEmptyMap(t *testing.T) {
	notEmptyMap := map[string]int{"key": 1}
	NotEmptyMap(notEmptyMap, "This should not panic")

	require.Panics(t, func() {
		emptyMap := map[string]int{}
		NotEmptyMap(emptyMap, "This should panic")
	})
}

func TestNotEmptyString(t *testing.T) {
	notEmptyString := "not empty"
	NotEmptyString(notEmptyString, "This should not panic")

	require.Panics(t, func() {
		emptyString := ""
		NotEmptyString(emptyString, "This should panic")
	})
}

func TestMapContainsKey(t *testing.T) {
	data := map[string]int{"key": 1}
	MapContainsKey(data, "key", "This should not panic")

	require.Panics(t, func() {
		MapContainsKey(data, "missing", "This should panic")
	})
}

func TestMapDoesNotContainKey(t *testing.T) {
	data := map[string]int{"key": 1}
	MapDoesNotContainKey(data, "missing", "This should not panic")

	require.Panics(t, func() {
		MapDoesNotContainKey(data, "key", "This should panic")
	})
}

func TestNilError(t *testing.T) {
	NilError(nil, "This should not panic")

	require.Panics(t, func() {
		NilError(fmt.Errorf("test error"), "This should panic")
	})
}
