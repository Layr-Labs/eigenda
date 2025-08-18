package replay

import (
	"time"
)

var _ ReplayGuardian = &noOpReplayGuardian{}

// noOpReplayGuardian is a ReplayGuardian that does nothing, always accepting requests without actually verifying them.
// Useful for unit tests where that want to be able to send duplicate requests without mocking the clock.
type noOpReplayGuardian struct{}

// NewNoOpReplayGuardian creates a new ReplayGuardian that does nothing, always accepting requests without actually
// verifying them. Useful for unit tests where that want to be able to send duplicate requests without mocking the
// clock.
func NewNoOpReplayGuardian() ReplayGuardian {
	return &noOpReplayGuardian{}
}

func (n *noOpReplayGuardian) VerifyRequest(requestHash []byte, requestTimestamp time.Time) error {
	return nil
}
