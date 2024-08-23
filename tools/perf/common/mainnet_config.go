package common

import "time"

var MainnetConfig = Config{
	MQ{
		QueueName: "eigenda_mainnet",
	},
	EigendaClient{
		Hostname:          "disperser.eigenda.xyz",
		Port:              "443",
		UseSecureGrpcFlag: true,
		Timeout:           60 * time.Second,
	},
	Utils{
		L1URL:           "https://mainnet.infura.io/v3/b75e92bfc8454d3abbcd0160e4b3debb",
		RetrieveTimeout: 10 * time.Second,
	},
}
