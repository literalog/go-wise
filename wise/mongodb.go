package wise

import (
	"github.com/literalog/go-wise/wise/internal/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepository[M any] mongodb.Repository[M]

func NewMongoSimpleRepository[M any](col *mongo.Collection) (mongodb.Repository[M], error) {
	return mongodb.NewSimpleRepository[M](col)
}

func NewMongoRepository[M, D any](col *mongo.Collection, ser mongodb.Serializer[M, D]) (mongodb.Repository[M], error) {
	return mongodb.NewRepository(col, ser)
}
