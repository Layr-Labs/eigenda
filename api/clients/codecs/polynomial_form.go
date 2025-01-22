package codecs

// PolynomialForm is an enum that describes the different ways a polynomial may be represented.
type PolynomialForm uint

const (
	// PolynomialFormEval is short for polynomial "evaluation form".
	// The field elements represent the evaluation of the polynomial at roots of unity.
	PolynomialFormEval PolynomialForm = iota
	// PolynomialFormCoeff is short for polynomial "coefficient form".
	// The field elements represent the coefficients of the polynomial.
	PolynomialFormCoeff
)
