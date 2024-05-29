package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type repository[M, D any] struct {
	Repository[D]
	serializer Serializer[M, D]
}

func NewRepository[M, D any](coll *mongo.Collection, ser Serializer[M, D]) (Repository[M], error) {
	if coll == nil {
		return nil, ErrNilCollection
	}

	if ser == nil {
		return nil, ErrNilSerializer
	}

	innerRepo, err := NewSimpleRepository[D](coll)
	if err != nil {
		return nil, err
	}

	repo := &repository[M, D]{
		Repository: innerRepo,
		serializer: ser,
	}

	return repo, nil
}

func (r *repository[M, D]) FindOne(ctx context.Context, filters map[string][]any) (M, error) {
	d, err := r.Repository.FindOne(ctx, filters)
	if err != nil {
		return *new(M), err
	}

	return r.serializer.Deserialize(d)
}

func (r *repository[M, D]) Find(ctx context.Context, filters map[string][]any, opts ...SearchOptions) ([]M, error) {
	dd, err := r.Repository.Find(ctx, filters, opts...)
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

func (r *repository[M, D]) FindById(ctx context.Context, id string) (M, error) {
	d, err := r.Repository.FindById(ctx, id)
	if err != nil {
		return *new(M), err
	}

	return r.serializer.Deserialize(d)
}

func (r *repository[M, D]) Upsert(ctx context.Context, id string, m M) error {
	d, err := r.serializer.Serialize(m)
	if err != nil {
		return err
	}

	return r.Repository.Upsert(ctx, id, d)
}

func (r *repository[M, D]) DeleteOne(ctx context.Context, filters map[string][]any) (M, error) {
	d, err := r.Repository.DeleteOne(ctx, filters)
	if err != nil {
		return *new(M), err
	}

	return r.serializer.Deserialize(d)
}

func (r *repository[M, D]) DeleteById(ctx context.Context, id string) (M, error) {
	d, err := r.Repository.DeleteById(ctx, id)
	if err != nil {
		return *new(M), err
	}

	return r.serializer.Deserialize(d)
}
