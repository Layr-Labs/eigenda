package codecs

// PolynomialForm is an enum that represents the different ways that a blob polynomial may be represented
type PolynomialForm uint

const (
	// Eval is short for "evaluation form". The field elements represent the evaluation at the polynomial's expanded
	// roots of unity
	Eval PolynomialForm = iota
	// Coeff is short for "coefficient form". The field elements represent the coefficients of the polynomial
	Coeff
)
