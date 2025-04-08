package services

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/anyisalin/mcp-openapi-to-mcp-adapter/models"
	"github.com/anyisalin/mcp-openapi-to-mcp-adapter/repositories"
)

// SSEConfigService handles SSE configuration operations
type SSEConfigService struct {
	repo *repositories.SSEConfigRepository
}

// NewSSEConfigService creates a new SSE configuration service
func NewSSEConfigService(repo *repositories.SSEConfigRepository) *SSEConfigService {
	return &SSEConfigService{
		repo: repo,
	}
}

// Create creates a new SSE configuration in the database
func (s *SSEConfigService) Create(ctx context.Context, schemaURL, baseURL string, headers map[string]string, filters []string) (string, error) {
	// Validate required fields
	if schemaURL == "" {
		return "", errors.New("schema URL is required")
	}

	if baseURL == "" {
		return "", errors.New("base URL is required")
	}

	// Create the configuration
	config := models.NewSSEConfig(schemaURL, baseURL, headers, filters)

	// Save to database
	id, err := s.repo.Create(ctx, config)
	if err != nil {
		return "", err
	}

	return id, nil
}

// GetByID retrieves an SSE configuration by its ID
func (s *SSEConfigService) GetByID(ctx context.Context, id string) (*models.SSEConfig, error) {
	return s.repo.FindByID(ctx, id)
}

// Update updates an existing SSE configuration
func (s *SSEConfigService) Update(ctx context.Context, id string, schemaURL, baseURL string, headers map[string]string, filters []string) error {
	// Retrieve the existing configuration
	config, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Update fields
	if schemaURL != "" {
		config.SchemaURL = schemaURL
	}
	if baseURL != "" {
		config.BaseURL = baseURL
	}
	if headers != nil {
		config.Headers = headers
	}
	if filters != nil {
		config.Filters = filters
	}

	// Save to database
	return s.repo.Update(ctx, id, config)
}

// Delete removes an SSE configuration
func (s *SSEConfigService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// GetSchemaBytes retrieves the schema content as bytes
func (s *SSEConfigService) GetSchemaBytes(schemaURL string) ([]byte, error) {
	// Check if schemaURL is a local file or a URL
	if strings.HasPrefix(schemaURL, "http://") || strings.HasPrefix(schemaURL, "https://") {
		// Create a custom HTTP client with timeout
		client := &http.Client{
			Timeout: 30 * time.Second,
		}
		
		// Use the client to make the request
		resp, err := client.Get(schemaURL)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		
		return io.ReadAll(resp.Body)
	}
	
	// Local file
	return os.ReadFile(schemaURL)
} 