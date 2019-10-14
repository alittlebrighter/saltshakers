package io

import (
	c "context"
	"log"
	"net/http"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/go-chi/chi"

	"github.com/alittlebrighter/saltshakers/configuration"
	"github.com/alittlebrighter/saltshakers/messages"
	"github.com/alittlebrighter/saltshakers/utils"
)

func HttpRestProducer() actor.Actor {
	return &HttpRestActor{BaseActor: utils.NewBaseActor("io.httpRest")}
}

type HttpRestActor struct {
	*utils.BaseActor
	ctx      actor.Context
	rulesPID *actor.PID

	server  *http.Server
}

func (state *HttpRestActor) Receive(context actor.Context) {
	state.ctx = context

	switch msg := context.Message().(type) {
	case []configuration.IOConfig:
		var serveAt string
		for _, config := range msg {
			if config.Kind() != configuration.HttpRest {
				continue
			}

			serveAt = config.Params["serveAt"].(string)
		}
		state.startServer(serveAt)
	case *messages.PIDEnvelope:
		switch msg.Type() {
		case messages.RulesPID:
			state.rulesPID = msg.PID()
		}
	case *actor.Started:
		context.Request(context.Parent(), messages.NewPIDEnvelope(messages.RulesPID, nil))
	case *actor.Stopping:
		ctx, cancel := c.WithTimeout(c.Background(), 5 * time.Second)
		defer cancel()
		state.server.Shutdown(ctx)
	case *actor.Stopped:
	case *actor.Restarting:
		state.startServer("")
	}
}

func (state *HttpRestActor) startServer(serveAt string) {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	state.server = &http.Server{Addr: serveAt}
	go func() {
		log.Println(state.Name(), "server exited with:", state.server.ListenAndServe())
	}()
}
