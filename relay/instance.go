package relay

import (
	"net"

	"github.com/Layr-Labs/eigensdk-go/logging"
)

// Instance holds the state for a single relay server instance.
type Instance struct {
	Server   *Server
	Listener net.Listener
	Port     string
	URL      string
	Logger   logging.Logger
}

// Stop gracefully stops the relay instance.
// It performs a graceful shutdown of the gRPC server (which also closes the listener).
func (i *Instance) Stop() error {
	i.Logger.Info("Stopping relay instance", "url", i.URL)

	// Gracefully stop the gRPC server (this also closes the listener)
	if i.Server != nil {
		if err := i.Server.Stop(); err != nil {
			i.Logger.Warn("Error during server graceful stop", "error", err)
			return err
		}
	}

	return nil
}
