package mongodb

import "context"

type Repository[M any] interface {
	FindOne(ctx context.Context, filters map[string][]any) (M, error)
	Find(ctx context.Context, filters map[string][]any, opts ...SearchOptions) ([]M, error)
	FindById(ctx context.Context, id string) (M, error)

	Upsert(ctx context.Context, id string, m M) error

	DeleteOne(ctx context.Context, filters map[string][]any) (M, error)
	DeleteMany(ctx context.Context, filters map[string][]any) error
	DeleteById(ctx context.Context, id string) (M, error)

	CountDocuments(ctx context.Context, filters map[string][]any) (int64, error)
}
