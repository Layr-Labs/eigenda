package secret

// Convert a slice of Secret pointers to a slice of strings. If the input slice is nil, returns an empty slice.
func SecretSliceToStringSlice(secrets []*Secret) []string {
	if secrets == nil {
		return make([]string, 0)
	}

	result := make([]string, len(secrets))
	for i, secret := range secrets {
		result[i] = secret.Get()
	}
	return result
}
