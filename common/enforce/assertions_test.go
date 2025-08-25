package enforce

import "testing"

func TestTrue(t *testing.T) {
	True(true, "This should not panic")
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but did not panic")
		}
	}()

	True(false, "This should panic, and the deferred function will catch it")
}

func TestFalse(t *testing.T) {
	False(false, "This should not panic")
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but did not panic")
		}
	}()

	False(true, "This should panic, and the deferred function will catch it")
}

func TestEquals(t *testing.T) {
	Equals(1, 1, "This should not panic")
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but did not panic")
		}
	}()

	Equals(1, 2, "This should panic, and the deferred function will catch it")
}

func TestNotEquals(t *testing.T) {
	NotEquals(1, 2, "This should not panic")
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but did not panic")
		}
	}()

	NotEquals(1, 1, "This should panic, and the deferred function will catch it")
}

func TestNotNil(t *testing.T) {
	notNilValue := "not nil"
	NotNil(&notNilValue, "This should not panic")
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but did not panic")
		}
	}()

	var nilValue *string
	NotNil(nilValue, "This should panic, and the deferred function will catch it")
}

func TestNil(t *testing.T) {
	nilValue := (*string)(nil)
	Nil(nilValue, "This should not panic")
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but did not panic")
		}
	}()

	notNilValue := "not nil"
	Nil(&notNilValue, "This should panic, and the deferred function will catch it")
}

func TestNotEmptyList(t *testing.T) {
	notEmptyList := []int{1, 2, 3}
	NotEmptyList(notEmptyList, "This should not panic")
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but did not panic")
		}
	}()

	emptyList := []int{}
	NotEmptyList(emptyList, "This should panic, and the deferred function will catch it")
}

func TestNotEmptyMap(t *testing.T) {
	notEmptyMap := map[string]int{"key": 1}
	NotEmptyMap(notEmptyMap, "This should not panic")
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but did not panic")
		}
	}()

	emptyMap := map[string]int{}
	NotEmptyMap(emptyMap, "This should panic, and the deferred function will catch it")
}

func TestNotEmptyString(t *testing.T) {
	notEmptyString := "not empty"
	NotEmptyString(notEmptyString, "This should not panic")
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but did not panic")
		}
	}()

	emptyString := ""
	NotEmptyString(emptyString, "This should panic, and the deferred function will catch it")
}

func TestMapContainsKey(t *testing.T) {
	data := map[string]int{"key": 1}
	MapContainsKey(data, "key", "This should not panic")
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but did not panic")
		}
	}()

	MapContainsKey(data, "missing", "This should panic, and the deferred function will catch it")
}

func TestMapDoesNotContainKey(t *testing.T) {
	data := map[string]int{"key": 1}
	MapDoesNotContainKey(data, "missing", "This should not panic")
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but did not panic")
		}
	}()

	MapDoesNotContainKey(data, "key", "This should panic, and the deferred function will catch it")
}
