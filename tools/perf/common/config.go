package common

import "time"

type MQ struct {
	QueueName string
}

type EigendaClient struct {
	Hostname          string
	Port              string
	UseSecureGrpcFlag bool
	Timeout           time.Duration
}

type Utils struct {
	L1URL           string
	RetrieveTimeout time.Duration
}

type Config struct {
	MQ
	EigendaClient
	Utils
}
