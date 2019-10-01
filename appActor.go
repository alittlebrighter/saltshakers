package main

import (
	"log"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/alittlebrighter/saltshakers/configuration"
	"github.com/alittlebrighter/saltshakers/io"
	"github.com/alittlebrighter/saltshakers/messages"
	"github.com/alittlebrighter/saltshakers/persistence"
	"github.com/alittlebrighter/saltshakers/rules"
)

// AppActor controls the application by setting up all of the layers
type AppActor struct {
	config, io, rules, persistence *actor.PID
}

func (state *AppActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.PIDEnvelope:
		switch msg.Type() {
		case messages.ConfigurationPID:
			context.Respond(messages.NewPIDEnvelope(messages.ConfigurationPID, state.config))
		case messages.RulesPID:
			context.Respond(messages.NewPIDEnvelope(messages.RulesPID, state.rules))
		case messages.PersistencePID:
			context.Respond(messages.NewPIDEnvelope(messages.PersistencePID, state.persistence))
		}
	case *actor.Started:
		log.Println("AppActor started")
		// start in dependency order, config has no dependencies (at the moment)
		state.config = context.Spawn(actor.PropsFromProducer(configuration.Producer))
		state.persistence = context.Spawn(actor.PropsFromProducer(persistence.Producer))
		state.rules = context.Spawn(actor.PropsFromProducer(rules.Producer))
		state.io = context.Spawn(actor.PropsFromProducer(io.Producer))
	case *actor.Stopping:
		/*
			state.io.Stop()
			state.persistence.Stop()
			state.logic.Stop()
			state.config.Stop()
		*/
	case *actor.Stopped:

	case *actor.Restarting:

	}
}
