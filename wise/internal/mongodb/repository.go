package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (r *repository[M, D]) InsertOne(ctx context.Context, m M, opts ...*options.InsertOneOptions) error {
	d, err := r.serializer.Serialize(m)
	if err != nil {
		return err
	}

	return r.Repository.InsertOne(ctx, d, opts...)
}

func (r *repository[M, D]) InsertMany(ctx context.Context, mm []M, opts ...*options.InsertManyOptions) error {
	dd := make([]D, len(mm))

	for i, m := range mm {
		d, err := r.serializer.Serialize(m)
		if err != nil {
			return err
		}

		dd[i] = d
	}

	return r.Repository.InsertMany(ctx, dd, opts...)
}

func (r *repository[M, D]) UpdateOne(ctx context.Context, filters map[string][]any, m M, opts ...*options.UpdateOptions) error {
	d, err := r.serializer.Serialize(m)
	if err != nil {
		return err
	}

	return r.Repository.UpdateOne(ctx, filters, d, opts...)
}

func (r *repository[M, D]) UpdateMany(ctx context.Context, filters map[string][]any, mm []M, opts ...*options.UpdateOptions) error {
	dd := make([]D, len(mm))

	for i, m := range mm {
		d, err := r.serializer.Serialize(m)
		if err != nil {
			return err
		}

		dd[i] = d
	}

	return r.Repository.UpdateMany(ctx, filters, dd, opts...)
}

func (r *repository[M, D]) DeleteOne(ctx context.Context, filters map[string][]any) (M, error) {
	d, err := r.Repository.DeleteOne(ctx, filters)
	if err != nil {
		return *new(M), err
	}

	return r.serializer.Deserialize(d)
}
