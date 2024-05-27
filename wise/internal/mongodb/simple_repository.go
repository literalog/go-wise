package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type simpleRepository[M any] struct {
	collection      *mongo.Collection
	indexedFields   indexedFields
	opts            *repositoryOptions
	MaxPageSize     int
	DefaultPageSize int
}

func NewSimpleRepository[M any](col *mongo.Collection, opts ...RepositoryOptions) (Repository[M], error) {
	if col == nil {
		return nil, ErrNilCollection
	}

	repo := &simpleRepository[M]{
		collection:    col,
		indexedFields: newIndexedFields(*new(M)),
		opts:          NewRepositoryOptions(opts...),
	}

	return repo, nil
}

func (r *simpleRepository[M]) Find(ctx context.Context, id string) (M, error) {
	m := new(M)

	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&m)
	if err != nil {
		return *new(M), err
	}

	return *m, nil
}

func (r *simpleRepository[M]) FindAll(ctx context.Context) ([]M, error) {
	return r.search(ctx, bson.M{})
}

func (r *simpleRepository[M]) FindMany(ctx context.Context, ids []string) ([]M, error) {
	return r.search(ctx, bson.M{"_id": bson.M{"$in": ids}})
}

func (r *simpleRepository[M]) Search(ctx context.Context, filters map[string][]any, opts ...SearchOptions) ([]M, error) {
	bson, err := r.indexedFields.toBson(filters)
	if err != nil {
		return nil, err
	}

	opt := NewSearchOptions(opts...)

	return r.search(ctx, bson, opt.ToFindOptions(r.MaxPageSize))
}

func (r *simpleRepository[M]) Aggregate(ctx context.Context, pipeline map[string][]any) ([]M, error) {
	mm := make([]M, 0)

	cur, err := r.collection.Aggregate(ctx, pipeline)
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

func (r *simpleRepository[M]) CountDocuments(ctx context.Context, filters map[string][]any) (int64, error) {
	opts := options.Count()

	bson, err := r.indexedFields.toBson(filters)
	if err != nil {
		return 0, err
	}

	return r.collection.CountDocuments(ctx, bson, opts)
}

func (r *simpleRepository[M]) Upsert(ctx context.Context, id string, m M) error {
	opt := options.Update().SetUpsert(true)
	update := bson.M{"$set": m}

	_, err := r.collection.UpdateByID(context.TODO(), id, update, opt)

	return err
}

func (r *simpleRepository[M]) Delete(ctx context.Context, id string) (M, error) {
	m, err := r.Find(ctx, id)
	if err != nil {
		return *new(M), err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": id})
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
	_, err = r.collection.DeleteMany(ctx, mongoFilters)
	if err != nil {
		return err
	}

	return nil
}

func (r *simpleRepository[M]) search(ctx context.Context, filters bson.M, opts ...*options.FindOptions) ([]M, error) {
	mm := make([]M, 0)

	cur, err := r.collection.Find(ctx, filters, opts...)
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
