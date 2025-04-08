package mongo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Client represents a MongoDB client manager
type Client struct {
	config *Config
	client *mongo.Client
	db     *mongo.Database
	mu     sync.RWMutex
}

// NewClient creates a new MongoDB client with the provided configuration
func NewClient(config *Config) *Client {
	if config == nil {
		config = DefaultConfig()
	}
	return &Client{
		config: config,
	}
}

// Connect establishes a connection to the MongoDB server
func (c *Client) Connect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.client != nil {
		return nil // Already connected
	}

	// Ensure HeartbeatInterval is at least 1 second to prevent RTT monitor panic
	heartbeatInterval := c.config.HeartbeatInterval
	if heartbeatInterval <= 0 {
		heartbeatInterval = 10 // Default to 10 seconds if not set correctly
	}

	// Create client options based on the configuration
	clientOptions := options.Client().
		ApplyURI(c.config.URI).
		SetMaxPoolSize(c.config.MaxPoolSize).
		SetMinPoolSize(c.config.MinPoolSize).
		SetConnectTimeout(time.Duration(c.config.ConnectTimeout) * time.Second).
		SetSocketTimeout(time.Duration(c.config.SocketTimeout) * time.Second).
		SetServerSelectionTimeout(time.Duration(c.config.ServerSelectionTimeout) * time.Second).
		SetHeartbeatInterval(time.Duration(heartbeatInterval) * time.Second).
		SetLocalThreshold(time.Duration(c.config.LocalThreshold) * time.Millisecond)

	// Apply optional configuration
	if c.config.Username != "" && c.config.Password != "" {
		credential := options.Credential{
			Username: c.config.Username,
			Password: c.config.Password,
		}
		if c.config.AuthDatabase != "" {
			credential.AuthSource = c.config.AuthDatabase
		}
		clientOptions.SetAuth(credential)
	}

	if c.config.ReplicaSet != "" {
		clientOptions.SetReplicaSet(c.config.ReplicaSet)
	}

	// Connect to the MongoDB server
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Verify the connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		closeErr := client.Disconnect(ctx)
		if closeErr != nil {
			log.Printf("Failed to disconnect from MongoDB after ping failure: %v", closeErr)
		}
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	c.client = client

	// Set the default database if specified
	if c.config.Database != "" {
		c.db = client.Database(c.config.Database)
	}

	return nil
}

// Disconnect closes the connection to the MongoDB server
func (c *Client) Disconnect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.client == nil {
		return nil // Already disconnected
	}

	if err := c.client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
	}

	c.client = nil
	c.db = nil
	return nil
}

// Client returns the underlying MongoDB client
func (c *Client) Client() (*mongo.Client, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.client == nil {
		return nil, errors.New("MongoDB client is not connected")
	}
	return c.client, nil
}

// Database returns the MongoDB database
func (c *Client) Database() (*mongo.Database, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.db == nil {
		return nil, errors.New("MongoDB client is not connected or no database is specified")
	}
	return c.db, nil
}

// Collection returns the specified MongoDB collection from the default database
func (c *Client) Collection(name string) (*mongo.Collection, error) {
	db, err := c.Database()
	if err != nil {
		return nil, err
	}
	return db.Collection(name), nil
}

// Ping checks the connection to the MongoDB server
func (c *Client) Ping(ctx context.Context) error {
	client, err := c.Client()
	if err != nil {
		return err
	}
	return client.Ping(ctx, readpref.Primary())
}

// WithDatabase creates a new instance of Client with a different database
func (c *Client) WithDatabase(dbName string) (*Client, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.client == nil {
		return nil, errors.New("MongoDB client is not connected")
	}

	newClient := &Client{
		config: c.config,
		client: c.client,
		db:     c.client.Database(dbName),
	}
	return newClient, nil
} 