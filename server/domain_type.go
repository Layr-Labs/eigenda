package server

import (
	"fmt"
	"net/http"
)

var (
	ErrInvalidDomainType = fmt.Errorf("invalid domain type")
)

// DomainType is a enumeration type for the different data domains for which a
// blob can exist between
type DomainType uint8

const (
	BinaryDomain DomainType = iota
	PolyDomain
	UnknownDomain
)

func (d DomainType) String() string {
	switch d {
	case BinaryDomain:
		return "binary"
	case PolyDomain:
		return "polynomial"
	default:
		return "unknown"
	}
}

func StrToDomainType(s string) DomainType {
	switch s {
	case "binary":
		return BinaryDomain
	case "polynomial":
		return PolyDomain
	default:
		return UnknownDomain
	}
}

func ReadDomainFilter(r *http.Request) (DomainType, error) {
	query := r.URL.Query()
	key := query.Get(DomainFilterKey)
	if key == "" { // default
		return BinaryDomain, nil
	}
	dt := StrToDomainType(key)
	if dt == UnknownDomain {
		return UnknownDomain, ErrInvalidDomainType
	}

	return dt, nil
}
