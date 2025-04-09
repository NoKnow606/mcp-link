package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// APIServerConfig 代表一个 API 服务器的配置信息
type APIServerConfig struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	SchemaURL   string             `bson:"schema_url" json:"schemaUrl"`
	BaseURL     string             `bson:"base_url" json:"baseUrl"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updatedAt"`
}

// NewAPIServerConfig 创建一个新的 API 服务器配置
func NewAPIServerConfig(name, description, schemaURL, baseURL string) *APIServerConfig {
	now := time.Now()
	return &APIServerConfig{
		Name:        name,
		Description: description,
		SchemaURL:   schemaURL,
		BaseURL:     baseURL,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// BeforeInsert 在插入到数据库前执行
func (a *APIServerConfig) BeforeInsert() {
	if a.ID.IsZero() {
		a.ID = primitive.NewObjectID()
	}
	now := time.Now()
	a.CreatedAt = now
	a.UpdatedAt = now
}

// BeforeUpdate 在更新前执行
func (a *APIServerConfig) BeforeUpdate() {
	a.UpdatedAt = time.Now()
}

// GetID 获取记录ID
func (a *APIServerConfig) GetID() primitive.ObjectID {
	return a.ID
}

// SetID 设置记录ID
func (a *APIServerConfig) SetID(id primitive.ObjectID) {
	a.ID = id
}

// SetCreatedAt 设置创建时间
func (a *APIServerConfig) SetCreatedAt(t time.Time) {
	a.CreatedAt = t
}

// SetUpdatedAt 设置更新时间
func (a *APIServerConfig) SetUpdatedAt(t time.Time) {
	a.UpdatedAt = t
}

// Update 更新配置信息
func (a *APIServerConfig) Update(name, description, schemaURL, baseURL string) {
	if name != "" {
		a.Name = name
	}
	a.Description = description
	if schemaURL != "" {
		a.SchemaURL = schemaURL
	}
	if baseURL != "" {
		a.BaseURL = baseURL
	}
	a.UpdatedAt = time.Now()
}
