package codecs

// BlobForm is an enum that represents the different ways that a blob may be represented
type BlobForm uint

const (
	// Eval is short for "evaluation form". The field elements represent the evaluation at the polynomial's expanded
	// roots of unity
	Eval BlobForm = iota
	// Coeff is short for "coefficient form". The field elements represent the coefficients of the polynomial
	Coeff
)
