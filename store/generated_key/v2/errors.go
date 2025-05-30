package eigenda

import "fmt"

// Returned when `cert.L1InclusionBlock > batch.RBN + rbnRecencyWindowSize`,
// to indicate that the cert should be discarded from rollups' derivation pipeline.
// This should get converted by the proxy to an HTTP 418 TEAPOT error code.
type RBNRecencyCheckFailedError struct {
	certRBN              uint32
	certL1IBN            uint64
	rbnRecencyWindowSize uint64
}

func (e RBNRecencyCheckFailedError) Error() string {
	return fmt.Sprintf(
		"Invalid cert (rbn recency check failed): "+
			"certL1InclusionBlockNumber (%d) > cert.RBN (%d) + RBNRecencyWindowSize (%d)",
		e.certL1IBN, e.certRBN, e.rbnRecencyWindowSize,
	)
}
