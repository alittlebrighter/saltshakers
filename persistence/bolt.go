package persistence

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/alittlebrighter/saltshakers/messages"
	"github.com/alittlebrighter/saltshakers/utils"
)

func BoltDBProducer() actor.Actor {
	return &BoltDBActor{BaseActor: utils.NewBaseActor("persistence.boltdb")}
}

type BoltDBActor struct {
	*utils.BaseActor
}

func (state *BoltDBActor) Receive(context actor.Context) {
	switch context.Message().(type) {
	case messages.PIDEnvelope:

	case *actor.Started:
	case []configuration.Persistence:
		for _, config := range msg {
			if config.Kind() != configuration.Bolt {
				continue
			}
		}
	default:
		context.Forward(state.Children())
	}
}
