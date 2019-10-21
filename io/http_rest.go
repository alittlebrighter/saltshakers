package io

import (
	ctxLib "context"
	"encoding/base64"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/alittlebrighter/saltshakers/configuration"
	"github.com/alittlebrighter/saltshakers/messages"
	"github.com/alittlebrighter/saltshakers/models"
	"github.com/alittlebrighter/saltshakers/utils"
)

func HttpRestProducer() actor.Actor {
	return &HttpRestActor{BaseActor: utils.NewBaseActor("io.httpRest")}
}

type HttpRestActor struct {
	*utils.BaseActor
	ctx      actor.Context
	rulesPID *actor.PID

	serveAt string
	server  *echo.Echo
}

func (state *HttpRestActor) Receive(context actor.Context) {
	state.ctx = context

	switch msg := context.Message().(type) {
	case []configuration.IOConfig:
		for _, config := range msg {
			if config.Kind() != configuration.HttpRest {
				continue
			}

			state.serveAt = config.Params["serveAt"].(string)
		}
		state.startServer(state.serveAt)
	case *messages.PIDEnvelope:
		switch msg.Type() {
		case messages.RulesPID:
			state.rulesPID = msg.PID()
		}
	case *actor.Started:
		context.Request(context.Parent(), messages.NewPIDEnvelope(messages.RulesPID, nil))
	case *actor.Stopping:
		ctx, cancel := ctxLib.WithTimeout(ctxLib.Background(), 10*time.Second)
		defer cancel()
		state.server.Shutdown(ctx)
	case *actor.Stopped:
	case *actor.Restarting:
		state.startServer(state.serveAt)
	}
}

func (state *HttpRestActor) startServer(serveAt string) {
	state.server = echo.New()
	state.server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:8080"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	const prefix = "/api"

	state.server.POST(prefix+"/households", state.SaveHousehold)
	state.server.GET(prefix+"/households", state.GetHouseholds)
	state.server.GET(prefix+"/households/:id", state.GetHousehold)
	state.server.DELETE(prefix+"/households/:id", state.DeleteHousehold)

	state.server.GET(prefix+"/groups/generate", state.GenerateGroups)
	state.server.POST(prefix+"/groups", state.SaveGroups)
	state.server.GET(prefix+"/groups", state.GetGroups)

	go func() {
		log.Println(state.Name(), "server exited with:", state.server.Start(serveAt))
		state.ctx.Poison(state.ctx.Self())
	}()
}

func (state *HttpRestActor) SaveHousehold(c echo.Context) error {
	var payload messages.SaveHousehold
	if err := c.Bind(&payload); err != nil {
		c.String(http.StatusBadRequest, "could not parse request: "+err.Error())
	}

	return state.sendRequest(c, payload)
}

func (state *HttpRestActor) SaveGroups(c echo.Context) error {
	var payload []*models.Group
	if err := c.Bind(&payload); err != nil {
		c.String(http.StatusBadRequest, "could not parse request: "+err.Error())
	}

	return state.sendRequest(c, messages.SaveGroups{Groups: payload})
}

func (state *HttpRestActor) GetHouseholds(c echo.Context) error {
	return state.sendRequest(c, messages.Query{Entity: messages.HouseholdEntity})
}

func (state *HttpRestActor) GetGroups(c echo.Context) error {
	return state.sendRequest(c, messages.Query{Entity: messages.GroupEntity})
}

func (state *HttpRestActor) GetHousehold(c echo.Context) error {
	id, _ := ParseBytesFromURL(c.Param("id"))
	return state.sendRequest(c, messages.GetHousehold{Id: id})
}

func (state *HttpRestActor) DeleteHousehold(c echo.Context) error {
	id, _ := ParseBytesFromURL(c.Param("id"))
	return state.sendRequest(c, messages.DeleteHousehold{Id: id})
}

func (state *HttpRestActor) DeleteGroup(c echo.Context) error {
	var payload []*models.Group
	if err := c.Bind(&payload); err != nil {
		c.String(http.StatusBadRequest, "could not parse request: "+err.Error())
	}

	return state.sendRequest(c, messages.DeleteGroups{Groups: payload})
}

func (state *HttpRestActor) GenerateGroups(c echo.Context) error {
	hhCount, _ := strconv.Atoi(c.QueryParam("targetHouseholdCount"))
	return state.sendRequest(c, messages.GenerateGroups{TargetHouseholdCount: uint8(hhCount)})
}

func (state *HttpRestActor) sendRequest(c echo.Context, payload interface{}) error {
	result, err := state.ctx.RequestFuture(state.rulesPID, payload, 5*time.Second).Result()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, result)
}

func ParseBytesFromURL(param string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(param)
}
