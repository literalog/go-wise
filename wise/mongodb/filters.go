package mongodb

import (
	"reflect"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

type indexedFields map[string]struct{}

func (f indexedFields) toBson(input map[string][]any) (bson.M, error) {
	if len(input) == 0 {
		return bson.M{}, nil
	}

	filters := make(bson.M)

	for key, values := range input {
		keyPrefix := strings.Split(key, ".")[0]
		_, ok := f[keyPrefix]
		if !ok {
			continue
		}

		switch len(values) {
		case 0:
			continue
		case 1:
			filters[key] = values[0]
		default:
			filters[key] = bson.M{"$in": values}
		}
	}

	return filters, nil
}

func newIndexedFields(d any) indexedFields {
	typeOfD := reflect.TypeOf(d)

	mapIndexes := make(map[string]struct{})

	for i := 0; i < typeOfD.NumField(); i++ {
		field := typeOfD.Field(i)
		indexed := field.Tag.Get("indexed")
		indexedBool, _ := strconv.ParseBool(indexed)

		if indexedBool {
			raw := field.Tag.Get("bson")
			names := strings.SplitN(raw, ",", -1)
			for _, name := range names {
				if name == "omitempty" {
					continue
				}

				mapIndexes[name] = struct{}{}
			}
		}
	}

	return mapIndexes
}
