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
	case messages.SaveHousehold:
		if msg.Household.GetSurname() == "" {
			context.Respond("cannot save household with no surname")
			return
		}
		response, _ := context.RequestFuture(state.persistence, persistence.Create{
			EntityType: messages.HouseholdEntity.String(),
			Entity:     &models.HouseholdImpl{&msg.Household},
			Upsert:     true,
		}, timeout).Result()
		context.Respond(response.(persistence.Create).Entity)

	case messages.GetHousehold:
		response, _ := context.RequestFuture(state.persistence, persistence.GetOne{
			EntityType: messages.HouseholdEntity.String(),
			Id:         msg.Id,
			Entity:     &models.HouseholdImpl{new(models.Household)},
		}, timeout).Result()
		context.Respond(response.(persistence.GetOne).Entity)

	case messages.Query:
		if msg.Entity != messages.HouseholdEntity {
			return
		}

		response, _ := context.RequestFuture(state.persistence, persistence.Query{
			EntityType: messages.HouseholdEntity.String(),
			Model:      func() persistence.HasId { return &models.HouseholdImpl{new(models.Household)} },
		}, timeout).Result()
		context.Respond(response.(persistence.Query).Entities)

	case messages.DeleteHousehold:
		_, err := context.RequestFuture(state.persistence, persistence.Delete{
			Ids:        [][]byte{msg.Id},
			EntityType: messages.HouseholdEntity.String(),
		}, timeout).Result()
		context.Respond(err)

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
