package repositories

import (
	"context"

	"github.com/anyisalin/mcp-openapi-to-mcp-adapter/db/mongo"
	"github.com/anyisalin/mcp-openapi-to-mcp-adapter/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const SSEConfigCollectionName = "sse_configs"

// SSEConfigRepository handles database operations for SSE configurations
type SSEConfigRepository struct {
	repo *mongo.BaseRepository[*models.SSEConfig]
}

// NewSSEConfigRepository creates a new repository for SSE configurations
func NewSSEConfigRepository(client *mongo.Client) (*SSEConfigRepository, error) {
	repo := &SSEConfigRepository{
		repo: mongo.NewRepository[*models.SSEConfig](client, SSEConfigCollectionName),
	}
	return repo, nil
}

// Create saves a new SSE configuration to the database
func (r *SSEConfigRepository) Create(ctx context.Context, config *models.SSEConfig) (string, error) {
	// Call before insert hook
	config.BeforeInsert()
	
	// Insert the document
	if err := r.repo.Create(ctx, config); err != nil {
		return "", err
	}
	
	return config.GetID().Hex(), nil
}

// FindByID retrieves an SSE configuration by its ID
func (r *SSEConfigRepository) FindByID(ctx context.Context, id string) (*models.SSEConfig, error) {
	return r.repo.FindByID(ctx, id)
}

// FindOne finds a single document matching the filter
func (r *SSEConfigRepository) FindOne(ctx context.Context, filter interface{}) (*models.SSEConfig, error) {
	return r.repo.FindOne(ctx, filter)
}

// Update updates an existing SSE configuration
func (r *SSEConfigRepository) Update(ctx context.Context, id string, config *models.SSEConfig) error {
	// Convert string ID to ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	
	// Set the ID on the config
	config.SetID(objID)
	
	// Call before update hook
	config.BeforeUpdate()
	
	// Update the document
	return r.repo.Update(ctx, config)
}

// Delete removes an SSE configuration from the database
func (r *SSEConfigRepository) Delete(ctx context.Context, id string) error {
	return r.repo.Delete(ctx, id)
} 