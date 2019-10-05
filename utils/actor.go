package utils

import (
	"log"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/router"
	"github.com/alittlebrighter/saltshakers/messages"
)

type BaseActor struct {
	name        string
	knownPIDs   map[messages.PIDType]*actor.PID
	childProps  *actor.Props
	childrenPID *actor.PID

	restarts uint8
}

func NewBaseActor(name string) *BaseActor {
	return &BaseActor{name: name, knownPIDs: make(map[messages.PIDType]*actor.PID)}
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

func (state *BaseActor) ManagePIDs(context actor.Context, msg *messages.PIDEnvelope) {
	localPID, localExists := state.knownPIDs[msg.Type()]

	if msg.PID() == nil && localExists {
		log.Println(state.Name(), "providing known PID", msg.Type())
		context.Respond(messages.NewPIDEnvelope(msg.Type(), localPID))
	} else if msg.PID() == nil {
		log.Println(state.Name(), "requesting unknown PID", msg.Type())
		context.Forward(context.Parent())
	} else {
		log.Println(state.Name(), "setting and forwarding PID", msg.Type())
		state.knownPIDs[msg.Type()] = msg.PID()
		context.Forward(state.Children())
	}
}

func (state *BaseActor) GetPID(pidType messages.PIDType) *actor.PID {
	return state.knownPIDs[pidType]
}

func (state *BaseActor) Restarting(context actor.Context, msg *actor.Restarting) {
	state.restarts++

	if state.restarts > 3 {
		context.Poison(context.Self())
	}
}
