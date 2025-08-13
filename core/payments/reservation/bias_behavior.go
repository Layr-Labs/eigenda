package reservation

// In the leaky bucket implementation, there are different points where we need to decide whether we should err on the
// side of permitting *more* or *less* throughput.
//
// Consider the different users of the leaky bucket:
//   - Validator nodes should err on the side of permitting *more* throughput. Processing a little extra data isn't
//     a big deal, but denying usage that a user is entitled to is something to be avoided at all costs.
//   - Clients should err on the side of utilizing *less* throughput. They should do their best to use the
//     full capacity of the reservation they're entitled to, but should prefer slight under-use.
type BiasBehavior string

const (
	// When in doubt, permit *more* throughput instead of less.
	//
	// This is what a validator node should use.
	BiasPermitMore BiasBehavior = "permitMore"
	// When in doubt, permit *less* throughput instead of more.
	//
	// This is what a client should use.
	BiasPermitLess BiasBehavior = "permitLess"
)
