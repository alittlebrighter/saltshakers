package persistence

type HasID interface {
	ID() string
}

type Create struct {
	EntityType string
	Entity     interface{}
	Upsert     bool
}

type GetOne struct {
	EntityType string
	ID         string
}

type Query struct {
	Filters []Filter
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
}
