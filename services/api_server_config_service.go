package services

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/url"

	"github.com/anyisalin/mcp-openapi-to-mcp-adapter/models"
	"github.com/anyisalin/mcp-openapi-to-mcp-adapter/repositories"
)

// APIServerConfigService 是 API 服务器配置的服务接口
type APIServerConfigService interface {
	CreateAPIServerConfig(ctx context.Context, name, description, schemaURL, baseURL string) (*models.APIServerConfig, error)
	GetAPIServerConfigByID(ctx context.Context, id string) (*models.APIServerConfig, error)
	GetAllAPIServerConfigs(ctx context.Context) ([]*models.APIServerConfig, error)
	UpdateAPIServerConfig(ctx context.Context, id, name, description, schemaURL, baseURL string) (*models.APIServerConfig, error)
	DeleteAPIServerConfig(ctx context.Context, id string) error
}

// DefaultAPIServerConfigService 是 API 服务器配置服务的默认实现
type DefaultAPIServerConfigService struct {
	repo repositories.APIServerConfigRepository
}

// NewAPIServerConfigService 创建一个新的 API 服务器配置服务
func NewAPIServerConfigService(repo repositories.APIServerConfigRepository) APIServerConfigService {
	return &DefaultAPIServerConfigService{
		repo: repo,
	}
}

// validateURLs 验证 URL 格式
func validateURLs(schemaURL, baseURL string) error {
	// 验证 schemaURL
	if schemaURL != "" {
		if _, err := url.ParseRequestURI(schemaURL); err != nil {
			return errors.Wrap(err, "invalid schema URL format")
		}
	}

	// 验证 baseURL
	if baseURL != "" {
		if _, err := url.ParseRequestURI(baseURL); err != nil {
			return errors.Wrap(err, "invalid base URL format")
		}
	}

	return nil
}

// CreateAPIServerConfig 创建新的 API 服务器配置
func (s *DefaultAPIServerConfigService) CreateAPIServerConfig(ctx context.Context, name, description, schemaURL, baseURL string) (*models.APIServerConfig, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}

	if schemaURL == "" {
		return nil, errors.New("schema URL is required")
	}

	if baseURL == "" {
		return nil, errors.New("base URL is required")
	}

	if err := validateURLs(schemaURL, baseURL); err != nil {
		return nil, err
	}

	config := models.NewAPIServerConfig(name, description, schemaURL, baseURL)
	// 调用前置钩子
	config.BeforeInsert()

	return s.repo.Create(ctx, config)
}

// GetAPIServerConfigByID 通过 ID 获取 API 服务器配置
func (s *DefaultAPIServerConfigService) GetAPIServerConfigByID(ctx context.Context, id string) (*models.APIServerConfig, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	return s.repo.GetByID(ctx, id)
}

// GetAllAPIServerConfigs 获取所有 API 服务器配置
func (s *DefaultAPIServerConfigService) GetAllAPIServerConfigs(ctx context.Context) ([]*models.APIServerConfig, error) {
	return s.repo.GetAll(ctx)
}

// UpdateAPIServerConfig 更新 API 服务器配置
func (s *DefaultAPIServerConfigService) UpdateAPIServerConfig(ctx context.Context, id, name, description, schemaURL, baseURL string) (*models.APIServerConfig, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	if err := validateURLs(schemaURL, baseURL); err != nil {
		return nil, err
	}

	// 获取当前配置
	config, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if config == nil {
		return nil, errors.New("config not found")
	}

	// 更新配置
	config.Update(name, description, schemaURL, baseURL)
	// 调用前置钩子
	config.BeforeUpdate()

	return s.repo.Update(ctx, config)
}

// DeleteAPIServerConfig 删除 API 服务器配置
func (s *DefaultAPIServerConfigService) DeleteAPIServerConfig(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is required")
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.Wrap(err, "invalid id format")
	}

	return s.repo.Delete(ctx, objID)
}
