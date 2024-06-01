package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository[M any] interface {
	FindOne(ctx context.Context, filters map[string][]any) (M, error)
	Find(ctx context.Context, filters map[string][]any, opts ...SearchOptions) ([]M, error)

	InsertOne(ctx context.Context, m M, opts ...*options.InsertOneOptions) error
	InsertMany(ctx context.Context, mm []M, opts ...*options.InsertManyOptions) error

	UpdateOne(ctx context.Context, filters map[string][]any, m M, opts ...*options.UpdateOptions) error
	UpdateMany(ctx context.Context, filters map[string][]any, mm []M, opts ...*options.UpdateOptions) error

	DeleteOne(ctx context.Context, filters map[string][]any) (M, error)
	DeleteMany(ctx context.Context, filters map[string][]any) error

	CountDocuments(ctx context.Context, filters map[string][]any) (int64, error)
}
