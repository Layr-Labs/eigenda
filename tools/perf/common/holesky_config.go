package common

import "time"

var HolesSkyConfig = Config{
	MQ{
		QueueName: "eigenda_holesky",
	},
	EigendaClient{
		Hostname:          "disperser-holesky.eigenda.xyz",
		Port:              "443",
		UseSecureGrpcFlag: true,
		Timeout:           60 * time.Second,
	},
	Utils{
		L1URL:           "https://holesky.infura.io/v3/b75e92bfc8454d3abbcd0160e4b3debb",
		RetrieveTimeout: 10 * time.Second,
	},
}
