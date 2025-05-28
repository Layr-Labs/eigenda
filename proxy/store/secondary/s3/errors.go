package s3

// Keccak256KeyValueMismatchError is an error that indicates a mismatch between the key and the keccaked value.
// KeccakCommitments should always respect the invariant that key=keccak(value).
// Before writing to S3 (in the POST route), or after reading the value from S3 (in the GET route),
// we check this invariant and return this error if it is violated.
// We only store the keccakedValue directly and not the value because the value is a full payload,
// which could be large (e.g. 1MB).
//
// TODO: this doesn't belong in the s3 package, but currently the Verify function returns
// this error and is on S3. That also should be moved elsewhere.
type Keccak256KeyValueMismatchError struct {
	Key           string
	KeccakedValue string
}

func NewKeccak256KeyValueMismatchErr(key, keccakedValue string) Keccak256KeyValueMismatchError {
	return Keccak256KeyValueMismatchError{
		Key:           key,
		KeccakedValue: keccakedValue,
	}
}

func (e Keccak256KeyValueMismatchError) Error() string {
	return "key!=keccak(value): key=" + e.Key + " keccak(value)=" + e.KeccakedValue
}
