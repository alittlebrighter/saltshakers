package messages

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/alittlebrighter/saltshakers/models"
)

type PIDEnvelope struct {
	pidType PIDType
	pid     *actor.PID
}

func (env *PIDEnvelope) Type() PIDType {
	return env.pidType
}

func (env *PIDEnvelope) PID() *actor.PID {
	return env.pid
}

func NewPIDEnvelope(pidType PIDType, pid *actor.PID) *PIDEnvelope {
	return &PIDEnvelope{pidType, pid}
}

type PIDType string

const (
	ConfigurationPID PIDType = "Configuration"
	RulesPID                 = "Rules"
	PersistencePID           = "Persistence"
)

type GetIOConfiguration struct{}
type GetRulesConfiguration struct{}
type GetPersistenceConfiguration struct{}

// SaveHousehold upserts a household
type SaveHousehold struct {
	models.Household
}

type GetHousehold struct {
	Id []byte `json:"id"`
}

type DeleteHousehold struct {
	Id []byte `json:"id"`
}

type SaveGroups struct {
	Groups []*models.Group `json:"groups"`
}

type DeleteGroups struct {
	Groups []*models.Group `json:"groups"`
}

type GenerateGroups struct {
	TargetHouseholdCount uint8          `json:"targetHouseholdCount"`
	Groups               []models.Group `json:"groups"`
}

type Query struct {
	Entity  EntityType
	Filters []Filter `json:"filters"`
}

type Filter struct {
	Op, Key string
	Value   interface{}
}

type EntityType string

const (
	HouseholdEntity EntityType = "Household"
	GroupEntity     EntityType = "Group"
)

func (e EntityType) String() string {
	return string(e)
}

var FilterOperations = map[string]func(interface{}) func(interface{}) bool{
	"equals":      FilterEquals,
	"greaterThan": FilterGreaterThan,
}

// FilterEquals works for primitive data types
func FilterEquals(val interface{}) func(interface{}) bool {
	return func(other interface{}) bool {
		return val == other
	}
}

func FilterGreaterThan(val interface{}) func(interface{}) bool {
	return func(other interface{}) bool {
		switch val.(type) {
		case int64:
			return val.(int64) < other.(int64)
		default:
			return false
		}
	}
}
