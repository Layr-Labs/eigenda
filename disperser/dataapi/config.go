package dataapi

type Config struct {
	SocketAddr         string
	ServerMode         string
	AllowOrigins       []string
	DisperserHostname  string
	ChurnerHostname    string
	BatcherHealthEndpt string
}

type DataApiVersion uint

const (
	V1 DataApiVersion = 1
	V2 DataApiVersion = 2
)
