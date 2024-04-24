package mongodb

import (
	"encoding/json"
	"reflect"
)

type Serializer[M, D any] interface {
	Serialize(m M) (D, error)
	Deserialize(d D) (M, error)
}

type DefaultSerializer[M, D any] struct{}

func NewDefaultSerializer[M, D any]() (Serializer[M, D], error) {
	if reflect.TypeOf(*new(M)).Kind() != reflect.Struct ||
		reflect.TypeOf(*new(D)).Kind() != reflect.Struct {
		return nil, ErrNotStruct
	}

	return DefaultSerializer[M, D]{}, nil
}

func (ds DefaultSerializer[M, D]) Serialize(m M) (D, error) {
	d := new(D)

	b, err := json.Marshal(m)
	if err != nil {
		return *d, err
	}

	err = json.Unmarshal(b, &d)
	if err != nil {
		return *d, err
	}

	return *d, nil
}

func (ds DefaultSerializer[M, D]) Deserialize(d D) (M, error) {
	m := new(M)

	b, err := json.Marshal(d)
	if err != nil {
		return *m, err
	}

	err = json.Unmarshal(b, &m)
	if err != nil {
		return *m, err
	}

	return *m, nil
}
