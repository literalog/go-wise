package mongodb

import "context"

type Repository[M any] interface {
	FindOne(ctx context.Context, filters map[string][]any) (M, error)
	Find(ctx context.Context, filters map[string][]any, opts ...SearchOptions) ([]M, error)

	Upsert(ctx context.Context, id string, m M) error

	DeleteOne(ctx context.Context, filters map[string][]any) (M, error)
	DeleteMany(ctx context.Context, filters map[string][]any) error

	CountDocuments(ctx context.Context, filters map[string][]any) (int64, error)
}
