package config

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

// RuntimeConfig represents the configuration that can be modified while the traffic generator is running
type RuntimeConfig struct {
	// WriterGroups defines different groups of writers with their own configurations
	WriterGroups []WriterGroupConfig `yaml:"writer_groups"`
}

// WriterGroupConfig represents the configuration for a group of writers with the same settings
type WriterGroupConfig struct {
	// Name of the writer group for identification
	Name string `yaml:"name"`

	// The number of worker threads that generate write traffic.
	NumWriteInstances uint `yaml:"num_write_instances"`

	// The period of the submission rate of new blobs for each write worker thread.
	WriteRequestInterval time.Duration `yaml:"write_request_interval"`

	// The Size of each blob dispersed, in bytes.
	DataSize uint64 `yaml:"data_size"`

	// If true, then each blob will contain unique random data. If false, the same random data
	// will be dispersed for each blob by a particular worker thread.
	RandomizeBlobs bool `yaml:"randomize_blobs"`

	// The amount of time to wait for a blob to be written.
	WriteTimeout time.Duration `yaml:"write_timeout"`

	// Custom quorum numbers to use for the traffic generator.
	CustomQuorums []uint8 `yaml:"custom_quorums"`
}

// RuntimeConfigManager handles loading and watching of runtime configuration
type RuntimeConfigManager struct {
	sync.RWMutex
	currentConfig *RuntimeConfig
	configPath    string
	onChange      func(*RuntimeConfig)
}

// NewRuntimeConfigManager creates a new runtime config manager
func NewRuntimeConfigManager(configPath string, onChange func(*RuntimeConfig)) (*RuntimeConfigManager, error) {
	manager := &RuntimeConfigManager{
		configPath: configPath,
		onChange:   onChange,
	}

	// Load initial config
	if err := manager.loadConfig(); err != nil {
		return nil, err
	}

	return manager, nil
}

// GetConfig returns the current runtime configuration
func (m *RuntimeConfigManager) GetConfig() *RuntimeConfig {
	m.RLock()
	defer m.RUnlock()
	return m.currentConfig
}

// loadConfig loads the configuration from disk
func (m *RuntimeConfigManager) loadConfig() error {
	data, err := os.ReadFile(m.configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var config RuntimeConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate each writer group if any exist
	for i, group := range config.WriterGroups {
		if group.Name == "" {
			return fmt.Errorf("writer group at index %d must have a name", i)
		}
		if group.NumWriteInstances == 0 {
			return fmt.Errorf("writer group '%s' must have at least one writer instance", group.Name)
		}
		if group.WriteRequestInterval == 0 {
			return fmt.Errorf("writer group '%s' must have a non-zero write request interval", group.Name)
		}
		if group.DataSize == 0 {
			return fmt.Errorf("writer group '%s' must have a non-zero data size", group.Name)
		}
		if group.WriteTimeout == 0 {
			return fmt.Errorf("writer group '%s' must have a non-zero write timeout", group.Name)
		}
		if len(group.CustomQuorums) == 0 {
			return fmt.Errorf("writer group '%s' must have at least one custom quorum", group.Name)
		}
	}

	m.Lock()
	defer m.Unlock()

	// Check if config has actually changed
	if m.currentConfig != nil {
		// Convert both configs to YAML for comparison
		currentYAML, err := yaml.Marshal(m.currentConfig)
		if err != nil {
			return fmt.Errorf("failed to marshal current config: %w", err)
		}
		newYAML, err := yaml.Marshal(&config)
		if err != nil {
			return fmt.Errorf("failed to marshal new config: %w", err)
		}

		if string(currentYAML) == string(newYAML) {
			// No changes, skip update
			return nil
		}
	}

	m.currentConfig = &config

	if m.onChange != nil {
		m.onChange(&config)
	}

	return nil
}

// StartWatching begins watching the config file for changes
func (m *RuntimeConfigManager) StartWatching(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				if err := m.loadConfig(); err != nil {
					// Just log the error and continue
					fmt.Printf("Error reloading config: %v\n", err)
				}
			}
		}
	}()
}
