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

// CreateHousehold upserts a household
type CreateHousehold struct {
	models.Household
}

type GetHousehold struct {
	Id []byte `json:"id"`
}

type QueryHouseholds struct {
	Filters []struct{ Key, Value, Op string } `json:"filters"`
}

type DeleteHousehold struct {
	Id []byte `json:"id"`
}

type SaveGroups struct {
	Groups []models.Group `json:"groups"`
}

type GenerateGroups struct {
	TargetHouseholdCount uint8          `json:"targetHouseholdCount"`
	Groups               []models.Group `json:"groups"`
}
