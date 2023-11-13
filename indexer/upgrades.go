package indexer

type UpgradeFork string

// UpgradeForkWatcher is a component that is used to scan a list of headers for an upgrade. Future upgrades may be based on a condition; past upgrades should have a block number configuration provided.
type UpgradeForkWatcher interface {

	// DetectUpgrade takes in a list of headers and sets the CurrentFork and IsUpgrade fields
	DetectUpgrade(headers Headers) Headers

	GetLatestUpgrade(header *Header) uint64
}
