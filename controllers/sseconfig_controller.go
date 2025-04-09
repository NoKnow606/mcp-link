package controllers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/anyisalin/mcp-openapi-to-mcp-adapter/services"
	"github.com/anyisalin/mcp-openapi-to-mcp-adapter/utils"
)

// SSEConfigController handles HTTP requests for SSE configurations
type SSEConfigController struct {
	service    *services.SSEConfigService
	sseServer  *utils.SSEServer
	sseBaseURL string
}

// CreateConfigRequest represents the request structure for creating a configuration
type CreateConfigRequest struct {
	ApiConfigId string            `json:"apiConfigId"`
	SchemaURL   string            `json:"schemaURL"`
	BaseURL     string            `json:"baseURL"`
	Headers     map[string]string `json:"headers"`
	Filters     []string          `json:"filters"`
}

// ConfigResponse represents the response structure for configuration operations
type ConfigResponse struct {
	ID      string `json:"id,omitempty"`
	SSEUrl  string `json:"sseUrl,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	Status  bool   `json:"status"`
}

// schemaBytesContextKey is used to store schema bytes in the request context
type schemaBytesContextKey struct{}

// NewSSEConfigController creates a new SSE configuration controller
func NewSSEConfigController(service *services.SSEConfigService, sseServer *utils.SSEServer, baseURL string) *SSEConfigController {
	return &SSEConfigController{
		service:    service,
		sseServer:  sseServer,
		sseBaseURL: baseURL,
	}
}

// CreateConfig handles the creation of a new SSE configuration
func (c *SSEConfigController) CreateConfig(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req CreateConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create configuration in database
	id, err := c.service.Create(r.Context(), req.ApiConfigId, req.SchemaURL, req.BaseURL, req.Headers, req.Filters)
	if err != nil {
		c.writeErrorResponse(w, "Failed to create configuration: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Build SSE URL with the configuration ID
	sseURL := c.buildSSEURL(id)

	// Return success response with the SSE URL
	c.writeSuccessResponse(w, "Configuration created successfully", id, sseURL)
}

// GetConfig retrieves an SSE configuration
func (c *SSEConfigController) GetConfig(w http.ResponseWriter, r *http.Request) {
	// Only accept GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL path
	// Assuming path is like /api/config/{id}
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		c.writeErrorResponse(w, "Invalid request path", http.StatusBadRequest)
		return
	}
	id := pathParts[len(pathParts)-1]

	// Get configuration from database
	config, err := c.service.GetByID(r.Context(), id)
	if err != nil {
		c.writeErrorResponse(w, "Failed to get configuration: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the configuration
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}

// UpdateConfig handles updating an existing SSE configuration
func (c *SSEConfigController) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	// Only accept PUT requests
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		c.writeErrorResponse(w, "Invalid request path", http.StatusBadRequest)
		return
	}
	id := pathParts[len(pathParts)-1]

	// Parse request body
	var req CreateConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update configuration in database
	err := c.service.Update(r.Context(), id, req.SchemaURL, req.BaseURL, req.Headers, req.Filters)
	if err != nil {
		c.writeErrorResponse(w, "Failed to update configuration: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Build SSE URL with the configuration ID
	sseURL := c.buildSSEURL(id)

	// Return success response
	c.writeSuccessResponse(w, "Configuration updated successfully", id, sseURL)
}

// DeleteConfig handles deleting an SSE configuration
func (c *SSEConfigController) DeleteConfig(w http.ResponseWriter, r *http.Request) {
	// Only accept DELETE requests
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		c.writeErrorResponse(w, "Invalid request path", http.StatusBadRequest)
		return
	}
	id := pathParts[len(pathParts)-1]

	// Delete configuration from database
	err := c.service.Delete(r.Context(), id)
	if err != nil {
		c.writeErrorResponse(w, "Failed to delete configuration: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	c.writeSuccessResponse(w, "Configuration deleted successfully", id, "")
}

// SSEHandler handles SSE connection requests with configuration ID
func (c *SSEConfigController) SSEHandler(w http.ResponseWriter, r *http.Request) {
	// Only accept GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract configID from query parameters
	configID := r.URL.Query().Get("configId")
	if configID == "" {
		http.Error(w, "Missing configId parameter", http.StatusBadRequest)
		return
	}

	// Get configuration from database
	config, err := c.service.GetByID(r.Context(), configID)
	if err != nil {
		http.Error(w, "Failed to get configuration: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get schema content
	schemaBytes, err := c.service.GetSchemaBytes(config.SchemaURL)
	if err != nil {
		http.Error(w, "Failed to get schema content: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a new context with the schema bytes
	ctx := context.WithValue(r.Context(), schemaBytesContextKey{}, schemaBytes)
	r = r.WithContext(ctx)

	// Set parameters directly in the request's URL query instead of using 'code' parameter
	q := r.URL.Query()
	q.Set("u", config.BaseURL)

	// Convert headers to JSON
	headersJSON, err := json.Marshal(config.Headers)
	if err == nil {
		q.Set("h", string(headersJSON))
	}

	// Add filters if any
	for _, filter := range config.Filters {
		q.Add("f", filter)
	}

	// Create an encoded parameters object to match what the SSE server expects
	paramsObj := map[string]interface{}{
		"s": config.SchemaURL,
		"u": config.BaseURL,
		"h": config.Headers,
	}

	if len(config.Filters) > 0 {
		paramsObj["f"] = strings.Join(config.Filters, ";")
	}

	// Encode the params as JSON and then base64
	paramsJSON, _ := json.Marshal(paramsObj)
	encodedParams := base64.StdEncoding.EncodeToString(paramsJSON)

	// Add the encoded params as 'code' param
	q.Set("code", encodedParams)

	// Set the modified query
	r.URL.RawQuery = q.Encode()

	// Handle the SSE connection with the original SSE server handler
	c.sseServer.ServeHTTP(w, r)
}

// buildSSEURL constructs the SSE URL with the configuration ID
func (c *SSEConfigController) buildSSEURL(id string) string {
	return c.sseBaseURL + "/sse?configId=" + id
}

// writeErrorResponse writes an error response to the client
func (c *SSEConfigController) writeErrorResponse(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := ConfigResponse{
		Error:  message,
		Status: false,
	}
	json.NewEncoder(w).Encode(response)
}

// writeSuccessResponse writes a success response to the client
func (c *SSEConfigController) writeSuccessResponse(w http.ResponseWriter, message, id, sseURL string) {
	w.Header().Set("Content-Type", "application/json")
	response := ConfigResponse{
		ID:      id,
		SSEUrl:  sseURL,
		Message: message,
		Status:  true,
	}
	json.NewEncoder(w).Encode(response)
}
