package persistence

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/alittlebrighter/saltshakers/messages"
	"github.com/alittlebrighter/saltshakers/utils"
)

func Producer() actor.Actor {
	return &PersistenceActor{BaseActor: utils.NewBaseActor("persistence")}
}

type PersistenceActor struct {
	*utils.BaseActor
	configPID
}

func (state *PersistenceActor) Receive(context actor.Context) {
main:
	switch msg := context.Message().(type) {
	case *messages.PIDEnvelope:
		if msg.PID() == nil {
			context.Forward(context.Parent())
			break main
		}

		switch msg.Type() {
		case messages.ConfigurationPID:
			state.configPID = msg.PID()
			context.Request(state.configPID, messages.GetIOConfiguration{})
		default:
			context.Forward(state.Children())
		}
	case *actor.Started:
		state.SetChildren(context,
			actor.PropsFromProducer(BoltDBProducer),
		)
	default:
		context.Forward(state.Children())
	}
}

func (state *PersistenceActor) Receive(context actor.Context) {
main:
	switch msg := context.Message().(type) {
	case *messages.PIDEnvelope:
		if msg.PID() == nil {
			context.Forward(context.Parent())
			break main
		}

		switch msg.Type() {
		case messages.ConfigurationPID:
			state.configPID = msg.PID()
			context.Request(state.configPID, messages.GetIOConfiguration{})
		default:
			context.Forward(state.Children())
		}

	case *actor.Started:
		state.SetChildren(context,
			actor.PropsFromProducer(WailsProducer),
		)

		context.Request(context.Parent(), messages.NewPIDEnvelope(messages.ConfigurationPID, nil))
	default:
		context.Forward(state.Children())
	}
}

type PersistenceType string

const (
	Queryable PersistenceType = "persistence.Queryable"
)
