package common

import "time"

var DevConfig = Config{
	MQ{
		QueueName: "eigenda_dev",
	},
	EigendaClient{
		Hostname:          "localhost",
		Port:              "32003",
		UseSecureGrpcFlag: false,
		Timeout:           60 * time.Second,
	},
	Utils{
		L1URL:           "http://127.0.0.1:8545",
		RetrieveTimeout: 10 * time.Second,
	},
}
