package router

import (
	"net/http"
	"strings"

	"github.com/anyisalin/mcp-openapi-to-mcp-adapter/controllers"
)

// Router handles HTTP requests routing
type Router struct {
	sseConfigController       *controllers.SSEConfigController
	apiServerConfigController *controllers.APIServerConfigController
}

// NewRouter creates a new router instance
func NewRouter(sseConfigController *controllers.SSEConfigController, apiServerConfigController *controllers.APIServerConfigController) *Router {
	return &Router{
		sseConfigController:       sseConfigController,
		apiServerConfigController: apiServerConfigController,
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
	if path == "/api/v1/config" {
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
	if strings.HasPrefix(path, "/api/v1/config/") && len(path) > 12 {
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

	// Routes for API server configuration
	if path == "/api/v1/api-server/config" {
		switch req.Method {
		case http.MethodPost:
			r.apiServerConfigController.CreateAPIServerConfig(w, req)
			return
		case http.MethodGet:
			// Get all configurations
			// 假设URL路径为 /api/server-config?all=true
			if req.URL.Query().Get("all") == "true" {
				r.apiServerConfigController.GetAPIServerConfig(w, req)
				return
			}
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
	}

	// Routes for specific API server configuration
	if strings.HasPrefix(path, "/api/v1/api-server/config/") && len(path) > 19 {
		switch req.Method {
		case http.MethodGet:
			r.apiServerConfigController.GetAPIServerConfig(w, req)
			return
		case http.MethodPut:
			r.apiServerConfigController.UpdateAPIServerConfig(w, req)
			return
		case http.MethodDelete:
			r.apiServerConfigController.DeleteAPIServerConfig(w, req)
			return
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
	}

	// If no routes match, return 404
	http.NotFound(w, req)
}
