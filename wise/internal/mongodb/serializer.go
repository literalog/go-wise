package mongodb

type Serializer[M, D any] interface {
	Serialize(m M) (D, error)
	Deserialize(d D) (M, error)
}
