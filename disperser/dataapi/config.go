package dataapi

type Config struct {
	SocketAddr         string
	ServerMode         string
	AllowOrigins       []string
	DisperserHostname  string
	ChurnerHostname    string
	BatcherHealthEndpt string
}
