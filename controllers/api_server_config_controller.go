package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/anyisalin/mcp-openapi-to-mcp-adapter/services"
)

// APIServerConfigController 处理 API 服务器配置的 HTTP 请求
type APIServerConfigController struct {
	service services.APIServerConfigService
}

// CreateAPIServerConfigRequest 表示创建配置的请求结构
type CreateAPIServerConfigRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	SchemaURL   string `json:"schemaUrl"`
	BaseURL     string `json:"baseUrl"`
}

// APIServerConfigResponse 表示配置操作的响应结构
type APIServerConfigResponse struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	SchemaURL   string `json:"schemaUrl,omitempty"`
	BaseURL     string `json:"baseUrl,omitempty"`
	Message     string `json:"message,omitempty"`
	Error       string `json:"error,omitempty"`
	Status      bool   `json:"status"`
}

// NewAPIServerConfigController 创建一个新的 API 服务器配置控制器
func NewAPIServerConfigController(service services.APIServerConfigService) *APIServerConfigController {
	return &APIServerConfigController{
		service: service,
	}
}

// CreateAPIServerConfig 处理创建新的 API 服务器配置
func (c *APIServerConfigController) CreateAPIServerConfig(w http.ResponseWriter, r *http.Request) {
	// 只接受 POST 请求
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 解析请求体
	var req CreateAPIServerConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 在数据库中创建配置
	config, err := c.service.CreateAPIServerConfig(r.Context(), req.Name, req.Description, req.SchemaURL, req.BaseURL)
	if err != nil {
		c.writeErrorResponse(w, "Failed to create API server configuration: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回成功响应
	c.writeSuccessResponse(w, "API server configuration created successfully", config.ID.Hex(), config.Name, config.Description, config.SchemaURL, config.BaseURL)
}

// GetAPIServerConfig 获取 API 服务器配置
func (c *APIServerConfigController) GetAPIServerConfig(w http.ResponseWriter, r *http.Request) {
	// 只接受 GET 请求
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 提取 URL 路径中的 ID
	// 假设路径是 /api/server-config/{id}
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		c.writeErrorResponse(w, "Invalid request path", http.StatusBadRequest)
		return
	}
	id := pathParts[len(pathParts)-1]

	// 获取所有配置
	if id == "all" {
		configs, err := c.service.GetAllAPIServerConfigs(r.Context())
		if err != nil {
			c.writeErrorResponse(w, "Failed to get API server configurations: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// 返回配置列表
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(configs)
		return
	}

	// 从数据库获取配置
	config, err := c.service.GetAPIServerConfigByID(r.Context(), id)
	if err != nil {
		c.writeErrorResponse(w, "Failed to get API server configuration: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if config == nil {
		c.writeErrorResponse(w, "API server configuration not found", http.StatusNotFound)
		return
	}

	// 返回配置
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}

// UpdateAPIServerConfig 处理更新现有的 API 服务器配置
func (c *APIServerConfigController) UpdateAPIServerConfig(w http.ResponseWriter, r *http.Request) {
	// 只接受 PUT 请求
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 提取 URL 路径中的 ID
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		c.writeErrorResponse(w, "Invalid request path", http.StatusBadRequest)
		return
	}
	id := pathParts[len(pathParts)-1]

	// 解析请求体
	var req CreateAPIServerConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 更新数据库中的配置
	config, err := c.service.UpdateAPIServerConfig(r.Context(), id, req.Name, req.Description, req.SchemaURL, req.BaseURL)
	if err != nil {
		c.writeErrorResponse(w, "Failed to update API server configuration: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回成功响应
	c.writeSuccessResponse(w, "API server configuration updated successfully", config.ID.Hex(), config.Name, config.Description, config.SchemaURL, config.BaseURL)
}

// DeleteAPIServerConfig 处理删除 API 服务器配置
func (c *APIServerConfigController) DeleteAPIServerConfig(w http.ResponseWriter, r *http.Request) {
	// 只接受 DELETE 请求
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 提取 URL 路径中的 ID
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		c.writeErrorResponse(w, "Invalid request path", http.StatusBadRequest)
		return
	}
	id := pathParts[len(pathParts)-1]

	// 从数据库删除配置
	err := c.service.DeleteAPIServerConfig(r.Context(), id)
	if err != nil {
		c.writeErrorResponse(w, "Failed to delete API server configuration: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回成功响应
	c.writeSuccessResponse(w, "API server configuration deleted successfully", id, "", "", "", "")
}

// writeErrorResponse 向客户端写入错误响应
func (c *APIServerConfigController) writeErrorResponse(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := APIServerConfigResponse{
		Error:  message,
		Status: false,
	}
	json.NewEncoder(w).Encode(response)
}

// writeSuccessResponse 向客户端写入成功响应
func (c *APIServerConfigController) writeSuccessResponse(w http.ResponseWriter, message, id, name, description, schemaURL, baseURL string) {
	w.Header().Set("Content-Type", "application/json")
	response := APIServerConfigResponse{
		ID:          id,
		Name:        name,
		Description: description,
		SchemaURL:   schemaURL,
		BaseURL:     baseURL,
		Message:     message,
		Status:      true,
	}
	json.NewEncoder(w).Encode(response)
} 