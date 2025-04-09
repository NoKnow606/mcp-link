package models

import (
	"time"

	"github.com/anyisalin/mcp-openapi-to-mcp-adapter/db/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SSEConfig represents the SSE configuration stored in MongoDB
type SSEConfig struct {
	mongo.BaseModel   `bson:",inline"`
	APIServerConfigId string            `json:"apiServerConfigId" bson:"api_server_config_id"`
	SchemaURL         string            `json:"schemaURL" bson:"schema_url"`      // URL or path to the OpenAPI schema
	BaseURL           string            `json:"baseURL" bson:"base_url"`          // Base URL for API requests
	Headers           map[string]string `json:"headers" bson:"headers"`           // Headers to send with API requests
	Filters           []string          `json:"filters" bson:"filters,omitempty"` // Filter expressions for API paths
	CreatedAt         time.Time         `json:"createdAt" bson:"created_at"`
	UpdatedAt         time.Time         `json:"updatedAt" bson:"updated_at,omitempty"`
}

// GetID returns the ID of the model
func (c *SSEConfig) GetID() primitive.ObjectID {
	return c.BaseModel.ID
}

// SetID sets the ID of the model
func (c *SSEConfig) SetID(id primitive.ObjectID) {
	c.BaseModel.ID = id
}

// SetCreatedAt sets the created timestamp
func (c *SSEConfig) SetCreatedAt(t time.Time) {
	c.CreatedAt = t
}

// SetUpdatedAt sets the updated timestamp
func (c *SSEConfig) SetUpdatedAt(t time.Time) {
	c.UpdatedAt = t
}

// NewSSEConfig creates a new SSE configuration
func NewSSEConfig(apiServerConfigId string, schemaURL string, baseURL string, headers map[string]string, filters []string) *SSEConfig {
	return &SSEConfig{
		APIServerConfigId: apiServerConfigId,
		SchemaURL:         schemaURL,
		BaseURL:           baseURL,
		Headers:           headers,
		Filters:           filters,
		CreatedAt:         time.Now(),
	}
}

// BeforeInsert is called before inserting the document
func (c *SSEConfig) BeforeInsert() {
	c.CreatedAt = time.Now()
}

// BeforeUpdate is called before updating the document
func (c *SSEConfig) BeforeUpdate() {
	c.UpdatedAt = time.Now()
}
