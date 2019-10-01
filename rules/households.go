package rules

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/alittlebrighter/saltshakers/messages"
	"github.com/alittlebrighter/saltshakers/utils"
)

func HouseholdRulesProducer() actor.Actor {
	return &HouseholdRulesActor{BaseActor: utils.NewBaseActor("rules.households")}
}

type HouseholdRulesActor struct {
	*utils.BaseActor
}

func (state *HouseholdRulesActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case messages.CreateHousehold:
		msg.Id = "success"
		context.Respond(msg)
	case *actor.Started:
	case *actor.Stopping:
	case *actor.Stopped:
	case *actor.Restarting:
	}
}
