package enforce

import "fmt"

// If convenient, it's ok at add additional assertions to this collection, as long as those assertions are
// general purpose and not specific to a particular domain or use case. For example, don't import custom
// types or packages that are not part of the standard library or common Go ecosystem.

// Asserts a condition is true and panics with a message if the condition is false.
func True(condition bool, message string) {
	if !condition {
		panic("Expected condition to be true: " + message)
	}
}

// Asserts a condition is false and panics with an error message if the condition is false.
func False(condition bool, message string) {
	if condition {
		panic("Expected condition to be false: " + message)
	}
}

// Asserts that two values are equal and panics with an error if they are not.
func Equals[T comparable](expected T, actual T, message string) {
	if expected != actual {
		panic(fmt.Sprintf("Expected equality, %v != %v: %s", expected, actual, message))
	}
}

// Asserts that two values are not equal and panics with an error if they are equal.
func NotEquals[T comparable](notExpected T, actual T, message string) {
	if notExpected == actual {
		panic(fmt.Sprintf("Expected inequality, %v == %v: %s", notExpected, actual, message))
	}
}

// Asserts that a value is not nil and panics with an error message if it is nil.
func NotNil[T any](value *T, message string) {
	if value == nil {
		panic("Expected value to be not nil: " + message)
	}
}

// Asserts that a value is nil and panics with an error message if it is not nil.
func Nil[T any](value *T, message string) {
	if value != nil {
		panic("Expected value to be nil: " + message)
	}
}

// Asserts that a slice is not empty and panics with an error message if it is empty.
func NotEmptyList[T any](list []T, message string) {
	if len(list) == 0 {
		panic("Expected list to be not empty: " + message)
	}
}

// Asserts that a string is not the empty string and panics with an error message if it is.
func NotEmptyString(value string, message string) {
	if value == "" {
		panic("Expected string to be not empty: " + message)
	}
}

// Asserts that a map is not empty and panics with an error message if it is empty.
func NotEmptyMap[K comparable, V any](m map[K]V, message string) {
	if len(m) == 0 {
		panic("Expected map to be not empty: " + message)
	}
}
