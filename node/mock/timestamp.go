package mock

import "time"

// MockTime implements Time interface for testing
type MockTime struct {
	NowFunc   func() time.Time
	UnixFunc  func(sec int64, nsec int64) time.Time
	SinceFunc func(t time.Time) time.Duration
}

// Now returns the mocked current time
func (mt *MockTime) Now() time.Time {
	if mt.NowFunc != nil {
		return mt.NowFunc()
	}
	return time.Time{}
}

// Unix returns the mocked Unix time
func (mt *MockTime) Unix(sec int64, nsec int64) time.Time {
	if mt.UnixFunc != nil {
		return mt.UnixFunc(sec, nsec)
	}
	return time.Unix(sec, nsec)
}

// Since returns the mocked duration since t
func (mt *MockTime) Since(t time.Time) time.Duration {
	if mt.SinceFunc != nil {
		return mt.SinceFunc(t)
	}
	return 0
}
