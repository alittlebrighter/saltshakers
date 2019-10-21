package persistence

type HasId interface {
	GetId() []byte
	SetId([]byte)
}

type PersistenceEnvelope struct {
	EntityType string
}

type Create struct {
	EntityType string
	Entity     HasId
	Upsert     bool
}

type CreateMany struct {
	EntityType string
	Entities   []HasId
	Upsert     bool
}

type GetOne struct {
	EntityType string
	Id         []byte
	Entity     HasId
}

// Query with no Filters gets all values
type Query struct {
	EntityType string
	Model      func() HasId // should be a type literal
	Entities   []HasId
	Filters    []Filter
}

type Filter struct {
	Key   string
	Value interface{}
	Op    CompareFunc
}

type CompareFunc func(a interface{}, b interface{}) bool

type Update struct {
	EntityType string
	Entity     HasId
}

type Delete struct {
	EntityType string
	Ids        [][]byte
}
