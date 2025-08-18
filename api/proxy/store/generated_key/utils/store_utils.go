package utils

// ConvertToRetryGoAttempts converts the user-facing PutTries value to retry-go's "attempts" semantic.
// In retry-go:
// - 0 "attempts" means retry forever (corresponds to our negative PutTries)
// - >0 "attempts" means try that many times total (corresponds to our PutTries values)
// Note: This function doesn't handle the PutTries=0 case, since 0 isn't a valid configuration, and this is checked
// at construction time
func ConvertToRetryGoAttempts(putTries int) uint {
	if putTries < 0 {
		return 0
	}
	return uint(putTries)
}
