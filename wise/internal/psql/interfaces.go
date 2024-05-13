package psql

import "context"

type Repository[M any] interface {
	Select(ctx context.Context, filters map[string][]any) ([]M, error)

	// Upsert(ctx context.Context, id string, m M) error

	// Delete(ctx context.Context, id string) (M, error)
	// DeleteMany(ctx context.Context, filters map[string][]any) error
}
