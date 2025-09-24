package version

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func verifySemver(t *testing.T, str string, major int, minor int, patch int, errata string) {
	semver, err := SemverFromString(str)
	require.NoError(t, err)

	require.Equal(t, major, semver.Major())
	require.Equal(t, minor, semver.Minor())
	require.Equal(t, patch, semver.Patch())
	require.Equal(t, errata, semver.Errata())
	require.Equal(t, str, semver.String())
}

func TestSerialization(t *testing.T) {
	verifySemver(t, "0.0.0", 0, 0, 0, "")
	verifySemver(t, "1.2.3", 1, 2, 3, "")
	verifySemver(t, "10.20.30", 10, 20, 30, "")
	verifySemver(t, "1.2.3-alpha", 1, 2, 3, "alpha")
	verifySemver(t, "1.2.3-beta.1", 1, 2, 3, "beta.1")
	verifySemver(t, "1.2.3-rc.1", 1, 2, 3, "rc.1")
	verifySemver(t, "1.2.3-rc.1+build.123", 1, 2, 3, "rc.1+build.123")
}

func TestInvalidSyntax(t *testing.T) {
	_, err := SemverFromString("1")
	require.Error(t, err)
	_, err = SemverFromString("1.2")
	require.Error(t, err)
	_, err = SemverFromString("1.2-beta")
	require.Error(t, err)
	_, err = SemverFromString("1.2.3.4")
	require.Error(t, err)
	_, err = SemverFromString("1.2.-3")
	require.Error(t, err)
	_, err = SemverFromString("asdfasdf1.2.3")
	require.Error(t, err)
	_, err = SemverFromString("asdfasdf1.2.3-beta")
	require.Error(t, err)
}

func TestEquals(t *testing.T) {
	a := NewSemver(1, 2, 3, "")
	b := NewSemver(1, 2, 3, "alpha")
	c := NewSemver(1, 2, 100, "")
	d := NewSemver(1, 100, 3, "")
	e := NewSemver(100, 2, 3, "")

	require.True(t, a.Equals(a))
	require.True(t, a.Equals(b))
	require.False(t, a.Equals(c))
	require.False(t, a.Equals(d))
	require.False(t, a.Equals(e))
}

func TestStrictEquals(t *testing.T) {
	a := NewSemver(1, 2, 3, "")
	b := NewSemver(1, 2, 3, "alpha")
	c := NewSemver(1, 2, 100, "")
	d := NewSemver(1, 100, 3, "")
	e := NewSemver(100, 2, 3, "")

	require.True(t, a.StrictEquals(a))
	require.False(t, a.StrictEquals(b))
	require.False(t, a.StrictEquals(c))
	require.False(t, a.StrictEquals(d))
	require.False(t, a.StrictEquals(e))
}

func TestLessThan(t *testing.T) {
	a := NewSemver(1, 2, 3, "")
	b := NewSemver(1, 2, 3, "alpha")
	c := NewSemver(1, 2, 100, "")
	d := NewSemver(1, 100, 3, "")
	e := NewSemver(100, 2, 3, "")
	f := NewSemver(0, 2, 3, "")
	g := NewSemver(1, 1, 3, "")
	h := NewSemver(1, 2, 2, "")

	require.False(t, a.LessThan(a))
	require.False(t, a.LessThan(b))
	require.True(t, a.LessThan(c))
	require.True(t, a.LessThan(d))
	require.True(t, a.LessThan(e))
	require.False(t, a.LessThan(f))
	require.False(t, a.LessThan(g))
	require.False(t, a.LessThan(h))
}

func TestGreaterThan(t *testing.T) {
	a := NewSemver(1, 2, 3, "")
	b := NewSemver(1, 2, 3, "alpha")
	c := NewSemver(1, 2, 100, "")
	d := NewSemver(1, 100, 3, "")
	e := NewSemver(100, 2, 3, "")
	f := NewSemver(0, 2, 3, "")
	g := NewSemver(1, 1, 3, "")
	h := NewSemver(1, 2, 2, "")

	require.False(t, a.GreaterThan(a))
	require.False(t, a.GreaterThan(b))
	require.False(t, a.GreaterThan(c))
	require.False(t, a.GreaterThan(d))
	require.False(t, a.GreaterThan(e))
	require.True(t, a.GreaterThan(f))
	require.True(t, a.GreaterThan(g))
	require.True(t, a.GreaterThan(h))
}

func TestLessThanOrEqual(t *testing.T) {
	a := NewSemver(1, 2, 3, "")
	b := NewSemver(1, 2, 3, "alpha")
	c := NewSemver(1, 2, 100, "")
	d := NewSemver(1, 100, 3, "")
	e := NewSemver(100, 2, 3, "")
	f := NewSemver(0, 2, 3, "")
	g := NewSemver(1, 1, 3, "")
	h := NewSemver(1, 2, 2, "")

	require.True(t, a.LessThanOrEqual(a))
	require.True(t, a.LessThanOrEqual(b))
	require.True(t, a.LessThanOrEqual(c))
	require.True(t, a.LessThanOrEqual(d))
	require.True(t, a.LessThanOrEqual(e))
	require.False(t, a.LessThanOrEqual(f))
	require.False(t, a.LessThanOrEqual(g))
	require.False(t, a.LessThanOrEqual(h))
}

func TestGreaterThanOrEqual(t *testing.T) {
	a := NewSemver(1, 2, 3, "")
	b := NewSemver(1, 2, 3, "alpha")
	c := NewSemver(1, 2, 100, "")
	d := NewSemver(1, 100, 3, "")
	e := NewSemver(100, 2, 3, "")
	f := NewSemver(0, 2, 3, "")
	g := NewSemver(1, 1, 3, "")
	h := NewSemver(1, 2, 2, "")

	require.True(t, a.GreaterThanOrEqual(a))
	require.True(t, a.GreaterThanOrEqual(b))
	require.False(t, a.GreaterThanOrEqual(c))
	require.False(t, a.GreaterThanOrEqual(d))
	require.False(t, a.GreaterThanOrEqual(e))
	require.True(t, a.GreaterThanOrEqual(f))
	require.True(t, a.GreaterThanOrEqual(g))
	require.True(t, a.GreaterThanOrEqual(h))
}

func TestComparator(t *testing.T) {
	a := NewSemver(1, 2, 3, "")
	b := NewSemver(1, 2, 3, "alpha")
	c := NewSemver(1, 2, 100, "")
	d := NewSemver(1, 100, 3, "")
	e := NewSemver(100, 2, 3, "")
	f := NewSemver(0, 2, 3, "")
	g := NewSemver(1, 1, 3, "")
	h := NewSemver(1, 2, 2, "")

	require.Equal(t, 0, SemverComparator(a, a))
	require.Equal(t, 0, SemverComparator(a, b))
	require.Equal(t, -1, SemverComparator(a, c))
	require.Equal(t, -1, SemverComparator(a, d))
	require.Equal(t, -1, SemverComparator(a, e))
	require.Equal(t, 1, SemverComparator(a, f))
	require.Equal(t, 1, SemverComparator(a, g))
	require.Equal(t, 1, SemverComparator(a, h))
}
