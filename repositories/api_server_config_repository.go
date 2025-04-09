package repositories

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/anyisalin/mcp-openapi-to-mcp-adapter/db/mongo"
	"github.com/anyisalin/mcp-openapi-to-mcp-adapter/models"
)

const APIServerConfigCollectionName = "api_server_configs"

// APIServerConfigRepository 处理 API 服务器配置的数据库操作
type APIServerConfigRepository interface {
	Create(ctx context.Context, config *models.APIServerConfig) (*models.APIServerConfig, error)
	GetByID(ctx context.Context, id string) (*models.APIServerConfig, error)
	GetAll(ctx context.Context) ([]*models.APIServerConfig, error)
	Update(ctx context.Context, config *models.APIServerConfig) (*models.APIServerConfig, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
}

// MongoAPIServerConfigRepository 是 APIServerConfigRepository 的MongoDB实现
type MongoAPIServerConfigRepository struct {
	repo *mongo.BaseRepository[*models.APIServerConfig]
}

// NewAPIServerConfigRepository 创建一个新的 API 服务器配置仓库
func NewAPIServerConfigRepository(client *mongo.Client) (APIServerConfigRepository, error) {
	if client == nil {
		return nil, errors.New("MongoDB client is nil")
	}
	
	repo := &MongoAPIServerConfigRepository{
		repo: mongo.NewRepository[*models.APIServerConfig](client, APIServerConfigCollectionName),
	}
	return repo, nil
}

// Create 将新的 API 服务器配置保存到数据库
func (r *MongoAPIServerConfigRepository) Create(ctx context.Context, config *models.APIServerConfig) (*models.APIServerConfig, error) {
	if err := r.repo.Create(ctx, config); err != nil {
		return nil, errors.Wrap(err, "failed to create API server config")
	}
	
	return config, nil
}

// GetByID 通过 ID 查找 API 服务器配置
func (r *MongoAPIServerConfigRepository) GetByID(ctx context.Context, id string) (*models.APIServerConfig, error) {
	var objID primitive.ObjectID
	var err error

	// 检查ID是否为有效的ObjectID
	if len(id) == 24 {
		objID, err = primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, errors.Wrap(err, "invalid id format")
		}
	} else {
		// 如果不是有效的ObjectID，使用自定义查询方法
		// 这里需要根据你的需求实现适当的查询逻辑
		return nil, errors.New("id must be a valid MongoDB ObjectID")
	}

	config, err := r.repo.FindOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return nil, errors.Wrap(err, "failed to find API server config by ID")
	}
	
	return config, nil
}

// GetAll 获取所有 API 服务器配置
func (r *MongoAPIServerConfigRepository) GetAll(ctx context.Context) ([]*models.APIServerConfig, error) {
	// 添加排序选项，按创建时间降序排序
	opts := options.Find().SetSort(bson.M{"createdAt": -1})
	
	configs, err := r.repo.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find all API server configs")
	}
	
	return configs, nil
}

// Update 更新现有的 API 服务器配置
func (r *MongoAPIServerConfigRepository) Update(ctx context.Context, config *models.APIServerConfig) (*models.APIServerConfig, error) {
	if err := r.repo.Update(ctx, config); err != nil {
		return nil, errors.Wrap(err, "failed to update API server config")
	}
	
	return config, nil
}

// Delete 从数据库中删除 API 服务器配置
func (r *MongoAPIServerConfigRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	if err := r.repo.Delete(ctx, id.Hex()); err != nil {
		return errors.Wrap(err, "failed to delete API server config")
	}
	
	return nil
} 