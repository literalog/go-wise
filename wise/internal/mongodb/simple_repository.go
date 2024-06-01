package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type simpleRepository[M any] struct {
	coll            *mongo.Collection
	indexedFields   indexedFields
	opts            *repositoryOptions
	MaxPageSize     int
	DefaultPageSize int
}

func NewSimpleRepository[M any](coll *mongo.Collection, opts ...RepositoryOptions) (Repository[M], error) {
	if coll == nil {
		return nil, ErrNilCollection
	}

	repo := &simpleRepository[M]{
		coll:          coll,
		indexedFields: newIndexedFields(*new(M)),
		opts:          NewRepositoryOptions(opts...),
	}

	return repo, nil
}

func (r *simpleRepository[M]) FindOne(ctx context.Context, filters map[string][]any) (M, error) {
	bson, err := r.indexedFields.toBson(filters)
	if err != nil {
		return *new(M), err
	}

	return r.searchOne(ctx, bson)
}

func (r *simpleRepository[M]) Find(ctx context.Context, filters map[string][]any, opts ...SearchOptions) ([]M, error) {
	bson, err := r.indexedFields.toBson(filters)
	if err != nil {
		return nil, err
	}

	opt := NewSearchOptions(opts...)

	return r.searchMany(ctx, bson, opt.ToFindOptions(r.MaxPageSize))
}

func (r *simpleRepository[M]) Upsert(ctx context.Context, id string, m M) error {
	opt := options.Update().SetUpsert(true)
	update := bson.M{"$set": m}

	_, err := r.coll.UpdateByID(context.TODO(), id, update, opt)

	return err
}

func (r *simpleRepository[M]) DeleteOne(ctx context.Context, filters map[string][]any) (M, error) {
	bson, err := r.indexedFields.toBson(filters)
	if err != nil {
		return *new(M), err
	}

	m, err := r.searchOne(ctx, bson)
	if err != nil {
		return *new(M), err
	}

	_, err = r.coll.DeleteOne(ctx, bson)
	if err != nil {
		return *new(M), err
	}

	return m, nil
}

func (r *simpleRepository[M]) DeleteMany(ctx context.Context, filters map[string][]any) error {
	mongoFilters, err := r.indexedFields.toBson(filters)
	if err != nil {
		return err
	}

	_, err = r.coll.DeleteMany(ctx, mongoFilters)
	if err != nil {
		return err
	}

	return nil
}

func (r *simpleRepository[M]) CountDocuments(ctx context.Context, filters map[string][]any) (int64, error) {
	opts := options.Count()

	bson, err := r.indexedFields.toBson(filters)
	if err != nil {
		return 0, err
	}

	return r.coll.CountDocuments(ctx, bson, opts)
}

func (r *simpleRepository[M]) searchMany(ctx context.Context, filters bson.M, opts ...*options.FindOptions) ([]M, error) {
	mm := make([]M, 0)

	cur, err := r.coll.Find(ctx, filters, opts...)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	err = cur.All(ctx, &mm)
	if err != nil {
		return nil, err
	}

	return mm, nil
}

func (r *simpleRepository[M]) searchOne(ctx context.Context, filters bson.M, opts ...*options.FindOneOptions) (M, error) {
	m := new(M)

	err := r.coll.FindOne(ctx, filters, opts...).Decode(&m)
	if err != nil {
		return *new(M), err
	}

	return *m, nil
}
