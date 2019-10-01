package utils

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/router"
	"github.com/alittlebrighter/saltshakers/messages"
)

type BaseActor struct {
	name        string
	knownPIDS   map[messages.PIDType]*actor.PID
	childProps  *actor.Props
	childrenPID *actor.PID
}

func NewBaseActor(name string) *BaseActor {
	return &BaseActor{name: name, knownPIDS: make(map[messages.PIDType]*actor.PID)}
}

func (a *BaseActor) Name() string {
	return a.name
}

func (a *BaseActor) SetChildren(context actor.Context, childProps ...*actor.Props) {
	impls := make([]*actor.PID, len(childProps))
	for i, props := range childProps {
		impls[i] = context.Spawn(props)
	}

	a.childProps = router.NewBroadcastGroup(impls...)
	a.childrenPID = context.Spawn(a.childProps)
}

func (a *BaseActor) Children() *actor.PID {
	return a.childrenPID
}

func (state *BaseActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.PIDEnvelope:
		localPID, localExists := state.knownPIDS[msg.Type()]

		if msg.PID() == nil && localExists {
			context.Respond(messages.NewPIDEnvelope(msg.Type(), localPID))
		} else if msg.PID() == nil {
			context.Forward(context.Parent())
		} else {
			state.knownPIDS[msg.Type()] = msg.PID()
			context.Forward(state.Children())
		}
	}
}
