package mongodb

import "errors"

var (
	ErrNotStruct     = errors.New("not a struct")
	ErrNilCollection = errors.New("nil collection")
	ErrNilSerializer = errors.New("nil serializer")
)
