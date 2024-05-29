# Go wise with MongoDB

Wise has two types of Mongo repositories:

1. Simple Repository: Utilized when a document perfectly aligns with a model.

2. (Custom) Repository: Utilized when a document's structure diverges from the model. In such instances, you'll need to develop a serializer.

If you're exclusively interested in the simple repository, you can skip ahead to [here](#).

## Getting started

```go
type User struct {
    ID string `json:"ID"`
    Name string `json:"name"`
}

type UserDoc struct {
    ID string `bson:"_id"`
    Name string `bson:"name"`
    Surname string `bson:"surname"`
}
```

## Serializer

The `serializer` provides functionality to serialize a model into a document and deserialize a document into a model.

```go
// Define a serializer struct
type serializer struct{}

// Create a new instance of the serializer
func NewSerializer() wise.Serializer[User, UserDoc] {
    return &serializer{}
}

// Convert Model to Document
func (s *serializer) Serialize(v User) (UserDoc, error) {    
    name, surname, _ := strings.Cut(v.Name, " ")

    return UserDoc{
        ID: v.ID,
        Name: name,
        Surname: surname
    }
}

// Convert Document to Model
func (s *serializer) Deserialize(v UserDoc) (User, error) {
    name := strings.Join([]string{v.Name, v.Surname}, " ")

    return User{
        ID: v.ID,
        Name: name,
    }
}
```