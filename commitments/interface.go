package commitments

type Commitment interface {
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}
