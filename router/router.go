package router

import (
	"net/http"
	"strings"

	"github.com/anyisalin/mcp-openapi-to-mcp-adapter/controllers"
)

// Router handles HTTP requests routing
type Router struct {
	sseConfigController *controllers.SSEConfigController
}

// NewRouter creates a new router instance
func NewRouter(sseConfigController *controllers.SSEConfigController) *Router {
	return &Router{
		sseConfigController: sseConfigController,
	}
}

// ServeHTTP implements the http.Handler interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	// Set CORS headers for all responses
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Handle OPTIONS (preflight) requests
	if req.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Routes for SSE configuration API
	if path == "/api/config" {
		switch req.Method {
		case http.MethodPost:
			r.sseConfigController.CreateConfig(w, req)
			return
		case http.MethodGet:
			// List configs - could be added later
			http.Error(w, "Not implemented", http.StatusNotImplemented)
			return
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
	}

	// Routes for specific configuration
	if strings.HasPrefix(path, "/api/config/") && len(path) > 12 {
		switch req.Method {
		case http.MethodGet:
			r.sseConfigController.GetConfig(w, req)
			return
		case http.MethodPut:
			r.sseConfigController.UpdateConfig(w, req)
			return
		case http.MethodDelete:
			r.sseConfigController.DeleteConfig(w, req)
			return
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
	}

	// Route for SSE connections with configuration ID
	if path == "/sse/config" {
		r.sseConfigController.SSEHandler(w, req)
		return
	}

	// If no routes match, return 404
	http.NotFound(w, req)
} 