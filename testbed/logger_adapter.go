package testbed

import (
	"fmt"

	"github.com/Layr-Labs/eigensdk-go/logging"
	tclog "github.com/testcontainers/testcontainers-go/log"
)

// loggerAdapter adapts eigensdk-go/logging.Logger to testcontainers log.Logger interface
type loggerAdapter struct {
	logger logging.Logger
}

// Printf implements the testcontainers log.Logger interface
func (la *loggerAdapter) Printf(format string, v ...any) {
	la.logger.Debug(fmt.Sprintf(format, v...))
}

// newTestcontainersLogger creates a testcontainers logger from an eigensdk logger
func newTestcontainersLogger(logger logging.Logger) tclog.Logger {
	return &loggerAdapter{logger: logger}
}