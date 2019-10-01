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

type CreateHousehold struct {
	models.Household
}

type ReadHousehold struct {
}

type UpdateHousehold struct {
	models.Household
}

type DeleteHousehold struct {
	Id string
}
