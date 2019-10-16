package rules

import (
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/alittlebrighter/saltshakers/messages"
	"github.com/alittlebrighter/saltshakers/utils"
)

const timeout = 2 * time.Second

func Producer() actor.Actor {
	return &RulesActor{BaseActor: utils.NewBaseActor("rules")}
}

type RulesActor struct {
	*utils.BaseActor
}

func (state *RulesActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.PIDEnvelope:
		state.ManagePIDs(context, msg)
	case *actor.Started:
		state.SetChildren(context,
			actor.PropsFromProducer(HouseholdRulesProducer),
			actor.PropsFromProducer(GroupRulesProducer),
		)
	case *actor.Stopping:
		state.Stopping(context)
	case *actor.Stopped:
	default:
		context.Forward(state.Children())
	}
}
