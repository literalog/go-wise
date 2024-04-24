package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type repository[M, D any] struct {
	Repository[D]
	serializer Serializer[M, D]
}

func NewRepository[M, D any](col *mongo.Collection, ser Serializer[M, D]) (Repository[M], error) {
	if col == nil {
		return nil, ErrNilCollection
	}

	if ser == nil {
		return nil, ErrNilSerializer
	}

	innerRepo, err := NewSimpleRepository[D](col)
	if err != nil {
		return nil, err
	}

	repo := &repository[M, D]{
		Repository: innerRepo,
		serializer: ser,
	}

	return repo, nil
}

func (r *repository[M, D]) Find(ctx context.Context, id string) (M, error) {
	d, err := r.Repository.Find(ctx, id)
	if err != nil {
		return *new(M), err
	}

	return r.serializer.Deserialize(d)
}

func (r *repository[M, D]) FindAll(ctx context.Context) ([]M, error) {
	dd, err := r.Repository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	mm := make([]M, len(dd))

	for i, d := range dd {
		m, err := r.serializer.Deserialize(d)
		if err != nil {
			return nil, err
		}

		mm[i] = m
	}

	return mm, nil
}

func (r *repository[M, D]) FindMany(ctx context.Context, ids []string) ([]M, error) {
	dd, err := r.Repository.FindMany(ctx, ids)
	if err != nil {
		return nil, err
	}

	mm := make([]M, len(dd))

	for i, d := range dd {
		m, err := r.serializer.Deserialize(d)
		if err != nil {
			return nil, err
		}

		mm[i] = m
	}

	return mm, nil
}

func (r *repository[M, D]) Search(ctx context.Context, filters map[string][]any, opts ...SearchOptions) ([]M, error) {
	dd, err := r.Repository.Search(ctx, filters, opts...)
	if err != nil {
		return nil, err
	}

	mm := make([]M, len(dd))

	for i, d := range dd {
		m, err := r.serializer.Deserialize(d)
		if err != nil {
			return nil, err
		}

		mm[i] = m
	}

	return mm, nil
}

func (r *repository[M, D]) Upsert(ctx context.Context, id string, m M) error {
	d, err := r.serializer.Serialize(m)
	if err != nil {
		return err
	}

	return r.Repository.Upsert(ctx, id, d)
}

func (r *repository[M, D]) Delete(ctx context.Context, id string) (M, error) {
	d, err := r.Repository.Delete(ctx, id)
	if err != nil {
		return *new(M), err
	}

	return r.serializer.Deserialize(d)
}
