package rules

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/alittlebrighter/saltshakers/messages"
	"github.com/alittlebrighter/saltshakers/models"
	"github.com/alittlebrighter/saltshakers/persistence"
	"github.com/alittlebrighter/saltshakers/utils"
)

func HouseholdRulesProducer() actor.Actor {
	return &HouseholdRulesActor{BaseActor: utils.NewBaseActor("rules.households")}
}

type HouseholdRulesActor struct {
	*utils.BaseActor
	persistence *actor.PID
}

func (state *HouseholdRulesActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case messages.CreateHousehold:
		response, _ := context.RequestFuture(state.persistence, persistence.Create{
			EntityType: HouseholdEntity.String(),
			Entity:     &models.HouseholdImpl{&msg.Household},
			Upsert:     true,
		}, timeout).Result()
		context.Respond(response.(persistence.Create).Entity)

	case messages.GetHousehold:
		response, _ := context.RequestFuture(state.persistence, persistence.GetOne{
			EntityType: HouseholdEntity.String(),
			Id:         msg.Id,
			Entity:     &models.HouseholdImpl{&models.Household{}},
		}, timeout).Result()
		context.Respond(response.(persistence.GetOne).Entity)

	case messages.QueryHouseholds:
		response, _ := context.RequestFuture(state.persistence, persistence.Query{
			EntityType: HouseholdEntity.String(),
			Model:      models.Household{},
		}, timeout).Result()
		context.Respond(response.(persistence.Query).Entities)

	case *messages.PIDEnvelope:
		switch msg.Type() {
		case messages.PersistencePID:
			state.persistence = msg.PID()
		}
	case *actor.Started:
		context.Request(context.Parent(), messages.NewPIDEnvelope(messages.PersistencePID, nil))
	case *actor.Stopping:
	case *actor.Stopped:
	case *actor.Restarting:
	}
}

type EntityType string

const (
	HouseholdEntity EntityType = "Household"
	GroupEntity     EntityType = "Group"
)

func (e EntityType) String() string {
	return string(e)
}
