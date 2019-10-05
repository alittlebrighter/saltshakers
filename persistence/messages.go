package persistence

type HasID interface {
	GetId() uint64
	SetId(uint64)
}

type PersistenceEnvelope struct {
	EntityType string
}

type Create struct {
	EntityType string
	Entity     HasID
	Upsert     bool
}

type GetOne struct {
	EntityType string
	Id         uint64
	Entity     HasID
}

// Query with no Filters gets all values
type Query struct {
	EntityType string
	Model      interface{} // should be a type literal
	Entities   []interface{}
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
	Entity     HasID
}

type Delete struct {
	EntityType string
	Id         uint64
}
