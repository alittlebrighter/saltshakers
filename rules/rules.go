package rules

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/alittlebrighter/saltshakers/messages"
	"github.com/alittlebrighter/saltshakers/utils"
)

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
		)
	default:
		context.Forward(state.Children())
	}
}
