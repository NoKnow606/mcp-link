package mongo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
)

var (
	defaultClient *Client
	clientMu      sync.RWMutex
)

// InitDefaultClient initializes the default MongoDB client with the provided configuration
func InitDefaultClient(ctx context.Context, config *Config) error {
	clientMu.Lock()
	defer clientMu.Unlock()

	if defaultClient != nil {
		// Disconnect the existing client before creating a new one
		if err := defaultClient.Disconnect(ctx); err != nil {
			return fmt.Errorf("failed to disconnect existing MongoDB client: %w", err)
		}
	}

	client := NewClient(config)
	if err := client.Connect(ctx); err != nil {
		return err
	}

	defaultClient = client
	return nil
}

// GetDefaultClient returns the default MongoDB client
func GetDefaultClient() (*Client, error) {
	clientMu.RLock()
	defer clientMu.RUnlock()

	if defaultClient == nil {
		return nil, errors.New("default MongoDB client is not initialized")
	}
	return defaultClient, nil
}

// CloseDefaultClient closes the default MongoDB client
func CloseDefaultClient(ctx context.Context) error {
	clientMu.Lock()
	defer clientMu.Unlock()

	if defaultClient == nil {
		return nil
	}

	if err := defaultClient.Disconnect(ctx); err != nil {
		return err
	}

	defaultClient = nil
	return nil
}

// LoadConfigFromFile loads MongoDB configuration from a JSON file
func LoadConfigFromFile(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// LoadConfigFromEnv loads MongoDB configuration from environment variables
func LoadConfigFromEnv() *Config {
	config := DefaultConfig()

	if uri := os.Getenv("MONGODB_URI"); uri != "" {
		config.URI = uri
	}

	if db := os.Getenv("MONGODB_DATABASE"); db != "" {
		config.Database = db
	}

	if username := os.Getenv("MONGODB_USERNAME"); username != "" {
		config.Username = username
	}

	if password := os.Getenv("MONGODB_PASSWORD"); password != "" {
		config.Password = password
	}

	if authDB := os.Getenv("MONGODB_AUTH_DATABASE"); authDB != "" {
		config.AuthDatabase = authDB
	}

	if replicaSet := os.Getenv("MONGODB_REPLICA_SET"); replicaSet != "" {
		config.ReplicaSet = replicaSet
	}

	if heartbeatStr := os.Getenv("MONGODB_HEARTBEAT_INTERVAL"); heartbeatStr != "" {
		if heartbeat, err := strconv.Atoi(heartbeatStr); err == nil {
			if heartbeat > 0 {
				config.HeartbeatInterval = heartbeat
			} else {
				// Default to 10 seconds if value is not positive
				config.HeartbeatInterval = 10
			}
		}
	}

	return config
} 