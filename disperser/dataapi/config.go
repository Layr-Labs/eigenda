package dataapi

type Config struct {
	SocketAddr         string
	ServerMode         string
	AllowOrigins       []string
	AvailabilityCheck  bool
	DisperserHostname  string
	ChurnerHostname    string
	BatcherHealthEndpt string
}
