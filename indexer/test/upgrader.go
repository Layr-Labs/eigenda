package weth_test

import "github.com/Layr-Labs/eigenda/indexer"

type Upgrader struct {
}

// DetectUpgrade takes in a list of headers and sets the CurrentFork and IsUpgrade fields
func (u *Upgrader) DetectUpgrade(headers indexer.Headers) indexer.Headers {
	for i := 0; i < len(headers); i++ {
		headers[i].CurrentFork = "genesis"
	}
	return headers
}

func (u *Upgrader) GetLatestUpgrade(header *indexer.Header) uint64 {
	return header.Number
}
