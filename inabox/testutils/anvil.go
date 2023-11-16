package testutils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/Layr-Labs/eigenda/inabox/config"
	ethinitscript "github.com/Layr-Labs/eigenda/inabox/eth-init-script"
	"github.com/Layr-Labs/eigenda/inabox/utils"
)

type AnvilContainer struct {
	anvil *exec.Cmd
	lock  *config.ConfigLock
}

func NewAnvilContainer(lock *config.ConfigLock) *AnvilContainer {
	return &AnvilContainer{
		lock: lock,
	}
}

func (a *AnvilContainer) MustStart() {
	rootDir := utils.MustGetModuleRootPath()

	a.anvil = ethinitscript.SetupEigenDA(a.lock, filepath.Join(rootDir, "contracts"))
}

func (a *AnvilContainer) MustStop() {
	if a.anvil == nil {
		panic(fmt.Errorf("anvil cannot be stopped, it was never started by us (maybe started externally)"))
	}
	err := a.anvil.Process.Signal(os.Interrupt)
	if err != nil {
		panic(err)
	}

	// Create a channel to signal when the process exits
	done := make(chan error)
	go func() {
		// Wait for the process to exit
		done <- a.anvil.Wait()
	}()

	// Select on the done channel and a 2 second timeout
	select {
	case <-time.After(2 * time.Second):
		// If the process does not exit in 2 seconds, kill it
		if err := a.anvil.Process.Kill(); err != nil {
			panic(fmt.Errorf("failed to kill anvil: %v", err))
		}
	case <-done:
		// Process exited gracefully
	}
}
