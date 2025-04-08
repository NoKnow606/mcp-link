package mongo

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Model is an interface for models with ID field
type Model interface {
	GetID() primitive.ObjectID
	SetID(id primitive.ObjectID)
	SetCreatedAt(time time.Time)
	SetUpdatedAt(time time.Time)
}

// BaseModel provides basic fields for MongoDB documents
type BaseModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// GetID returns the ID of the model
func (m *BaseModel) GetID() primitive.ObjectID {
	return m.ID
}

// SetID sets the ID of the model
func (m *BaseModel) SetID(id primitive.ObjectID) {
	m.ID = id
}

// SetCreatedAt sets the created timestamp
func (m *BaseModel) SetCreatedAt(t time.Time) {
	m.CreatedAt = t
}

// SetUpdatedAt sets the updated timestamp
func (m *BaseModel) SetUpdatedAt(t time.Time) {
	m.UpdatedAt = t
}

// Repository defines the interface for MongoDB repositories
type Repository[T Model] interface {
	FindByID(ctx context.Context, id string) (T, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) (T, error)
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) ([]T, error)
	Create(ctx context.Context, model T) error
	Update(ctx context.Context, model T) error
	Delete(ctx context.Context, id string) error
	DeleteMany(ctx context.Context, filter interface{}) (int64, error)
	Count(ctx context.Context, filter interface{}) (int64, error)
	Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) ([]T, error)
}

// BaseRepository is the base implementation of Repository
type BaseRepository[T Model] struct {
	client     *Client
	collection string
}

// NewRepository creates a new repository for the given model type
func NewRepository[T Model](client *Client, collection string) *BaseRepository[T] {
	return &BaseRepository[T]{
		client:     client,
		collection: collection,
	}
}

// getCollection returns the MongoDB collection
func (r *BaseRepository[T]) getCollection() (*mongo.Collection, error) {
	return r.client.Collection(r.collection)
}

// FindByID finds a document by its ID
func (r *BaseRepository[T]) FindByID(ctx context.Context, id string) (T, error) {
	var result T
	
	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}
	
	collection, err := r.getCollection()
	if err != nil {
		return result, err
	}
	
	filter := bson.M{"_id": objectID}
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return result, nil
		}
		return result, err
	}
	
	return result, nil
}

// FindOne finds a single document matching the filter
func (r *BaseRepository[T]) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) (T, error) {
	var result T
	
	collection, err := r.getCollection()
	if err != nil {
		return result, err
	}
	
	err = collection.FindOne(ctx, filter, opts...).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return result, nil
		}
		return result, err
	}
	
	return result, nil
}

// Find finds all documents matching the filter
func (r *BaseRepository[T]) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) ([]T, error) {
	collection, err := r.getCollection()
	if err != nil {
		return nil, err
	}
	
	cursor, err := collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var results []T
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	
	return results, nil
}

// Create inserts a new document
func (r *BaseRepository[T]) Create(ctx context.Context, model T) error {
	now := time.Now()
	
	// If ID is not set, generate a new one
	if model.GetID().IsZero() {
		model.SetID(primitive.NewObjectID())
	}
	
	// Set timestamps
	model.SetCreatedAt(now)
	model.SetUpdatedAt(now)
	
	collection, err := r.getCollection()
	if err != nil {
		return err
	}
	
	_, err = collection.InsertOne(ctx, model)
	return err
}

// Update updates an existing document
func (r *BaseRepository[T]) Update(ctx context.Context, model T) error {
	// Make sure we have an ID
	if model.GetID().IsZero() {
		return errors.New("model ID is required for update")
	}
	
	// Update timestamp
	model.SetUpdatedAt(time.Now())
	
	collection, err := r.getCollection()
	if err != nil {
		return err
	}
	
	filter := bson.M{"_id": model.GetID()}
	_, err = collection.ReplaceOne(ctx, filter, model)
	return err
}

// Delete removes a document by its ID
func (r *BaseRepository[T]) Delete(ctx context.Context, id string) error {
	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	
	collection, err := r.getCollection()
	if err != nil {
		return err
	}
	
	filter := bson.M{"_id": objectID}
	_, err = collection.DeleteOne(ctx, filter)
	return err
}

// DeleteMany removes multiple documents matching the filter
func (r *BaseRepository[T]) DeleteMany(ctx context.Context, filter interface{}) (int64, error) {
	collection, err := r.getCollection()
	if err != nil {
		return 0, err
	}
	
	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	
	return result.DeletedCount, nil
}

// Count returns the number of documents matching the filter
func (r *BaseRepository[T]) Count(ctx context.Context, filter interface{}) (int64, error) {
	collection, err := r.getCollection()
	if err != nil {
		return 0, err
	}
	
	return collection.CountDocuments(ctx, filter)
}

// Aggregate performs an aggregation pipeline
func (r *BaseRepository[T]) Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) ([]T, error) {
	collection, err := r.getCollection()
	if err != nil {
		return nil, err
	}
	
	cursor, err := collection.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var results []T
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	
	return results, nil
} 