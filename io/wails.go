package io

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/wailsapp/wails"

	"github.com/alittlebrighter/saltshakers/configuration"
	"github.com/alittlebrighter/saltshakers/messages"
	"github.com/alittlebrighter/saltshakers/utils"
)

func WailsProducer() actor.Actor {
	return &WailsActor{BaseActor: utils.NewBaseActor("io.wails")}
}

type WailsActor struct {
	*utils.BaseActor
	ctx      actor.Context
	app      *wails.App
	runtime  *wails.Runtime
	rulesPID *actor.PID
}

func (state *WailsActor) Receive(context actor.Context) {
	state.ctx = context

	switch msg := context.Message().(type) {
	case []configuration.IOConfig:
		for _, config := range msg {
			if config.Kind() != configuration.Wails {
				continue
			}

			wailsConfig := config.Params["wailsConfig"].(wails.AppConfig)
			state.app = wails.CreateApp(&wailsConfig)
		}
		state.app.Bind(state)
		go state.app.Run()
	case *messages.PIDEnvelope:
		switch msg.Type() {
		case messages.RulesPID:
			state.rulesPID = msg.PID()
		}
	case *actor.Started:
		context.Request(context.Parent(), messages.NewPIDEnvelope(messages.RulesPID, nil))
	case *actor.Stopping:
		fmt.Println("Stopping, wailsActor is about shut down")
	case *actor.Stopped:
		fmt.Println("Stopped, wailsActor and it's children are stopped")
	case *actor.Restarting:
		fmt.Println("Restarting, wailsActor is about restart")
	}
}

func (w *WailsActor) Request(msg string) (string, error) {
	if w.rulesPID == nil {
		return msg, errors.New("no rules available to process this request")
	}

	envelope := ActionEnvelope{}
	if err := json.Unmarshal([]byte(msg), &envelope); err != nil {
		log.Println(w.Name(), ": could not decode ActionEnvelope, message:", msg, "error:", err)
		return "", err
	}

	var payload interface{}
	switch envelope.ActionType {
	case "CreateHousehold":
		p := messages.CreateHousehold{}
		if err := json.Unmarshal(envelope.Payload, &p); err != nil {
			return "", err
		}
		payload = p
	case "GetHousehold":
		p := messages.GetHousehold{}
		if err := json.Unmarshal(envelope.Payload, &p); err != nil {
			return "", err
		}
		payload = p
	case "QueryHouseholds":
		p := messages.QueryHouseholds{}
		if err := json.Unmarshal(envelope.Payload, &p); err != nil {
			return "", err
		}
		payload = p
	}

	log.Println(w.Name(), "making request, payload:", payload)
	result, err := w.ctx.RequestFuture(w.rulesPID, payload, 5*time.Second).Result()
	if err != nil {
		return "", err
	}
	data, err := json.Marshal(&result)
	return string(data), err
}

// https://wails.app/reference/#wails-runtime
func (w *WailsActor) WailsInit(runtime *wails.Runtime) error {
	w.runtime = runtime

	return nil
}

type ActionEnvelope struct {
	ActionType string          `json:"type"`
	Payload    json.RawMessage `json:"payload"`
}
