package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type simpleRepository[D any] struct {
	collection      *mongo.Collection
	indexedFields   indexedFields
	opts            *repositoryOptions
	MaxPageSize     int
	DefaultPageSize int
}

func NewSimpleRepository[D any](col *mongo.Collection, opts ...RepositoryOptions) (Repository[D], error) {
	if col == nil {
		return nil, ErrNilCollection
	}

	repo := &simpleRepository[D]{
		collection:    col,
		indexedFields: newIndexedFields(*new(D)),
		opts:          NewRepositoryOptions(opts...),
	}

	return repo, nil
}

func (r *simpleRepository[D]) Find(ctx context.Context, id string) (D, error) {
	d := new(D)

	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&d)
	if err != nil {
		return *new(D), err
	}

	return *d, nil
}

func (r *simpleRepository[D]) FindAll(ctx context.Context) ([]D, error) {
	return r.search(ctx, bson.M{})
}

func (r *simpleRepository[D]) FindMany(ctx context.Context, ids []string) ([]D, error) {
	return r.search(ctx, bson.M{"_id": bson.M{"$in": ids}})
}

func (r *simpleRepository[D]) Search(ctx context.Context, filters map[string][]any, opts ...SearchOptions) ([]D, error) {
	bson, err := r.indexedFields.toBson(filters)
	if err != nil {
		return nil, err
	}

	opt := NewSearchOptions(opts...)

	return r.search(ctx, bson, opt.ToFindOptions(r.MaxPageSize))
}

func (r *simpleRepository[D]) CountDocuments(ctx context.Context, filters map[string][]any) (int64, error) {
	opts := options.Count()

	bson, err := r.indexedFields.toBson(filters)
	if err != nil {
		return 0, err
	}

	return r.collection.CountDocuments(ctx, bson, opts)
}

func (r *simpleRepository[D]) Upsert(ctx context.Context, id string, d D) error {
	opt := options.Update().SetUpsert(true)
	update := bson.M{"$set": d}

	_, err := r.collection.UpdateByID(context.TODO(), id, update, opt)

	return err
}

func (r *simpleRepository[D]) Delete(ctx context.Context, id string) (D, error) {
	d, err := r.Find(ctx, id)
	if err != nil {
		return *new(D), err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return *new(D), err
	}

	return d, nil
}

func (r *simpleRepository[D]) DeleteMany(ctx context.Context, filters map[string][]any) error {
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

func (r *simpleRepository[D]) search(ctx context.Context, filters bson.M, opts ...*options.FindOptions) ([]D, error) {
	dd := make([]D, 0)

	cur, err := r.collection.Find(ctx, filters, opts...)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	err = cur.All(ctx, &dd)
	if err != nil {
		return nil, err
	}

	return dd, nil
}
