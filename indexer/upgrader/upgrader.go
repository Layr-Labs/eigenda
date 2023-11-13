package eigendaUpgrader

import "github.com/Layr-Labs/eigenda/indexer"

type Upgrader struct {
}

// DetectUpgrade takes in a list of headers and sets the CurrentFork and IsUpgrade fields
func (u *Upgrader) DetectUpgrade(headers []indexer.Header) []indexer.Header {
	return nil
}

func (u *Upgrader) GetLatestUpgrade(header indexer.Header) uint64 {
	return 0
}
