package rules

import (
	"sort"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/golang/protobuf/ptypes/timestamp"

	"github.com/alittlebrighter/saltshakers/messages"
	"github.com/alittlebrighter/saltshakers/models"
	"github.com/alittlebrighter/saltshakers/persistence"
	"github.com/alittlebrighter/saltshakers/utils"
)

func GroupRulesProducer() actor.Actor {
	return &GroupRulesActor{BaseActor: utils.NewBaseActor("rules.groups")}
}

type GroupRulesActor struct {
	*utils.BaseActor
	persistence *actor.PID
}

func (state *GroupRulesActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case messages.GenerateGroups:
		groupsFuture := context.RequestFuture(state.persistence, persistence.Query{
			EntityType: messages.GroupEntity.String(),
			Model:      func() persistence.HasId { return &models.GroupImpl{new(models.Group)} },
		}, timeout)
		householdsFuture := context.RequestFuture(state.persistence, persistence.Query{
			EntityType: messages.HouseholdEntity.String(),
			Model:      func() persistence.HasId { return &models.HouseholdImpl{new(models.Household)} },
		}, timeout)

		householdsResult, _ := ArrayFromQueryFuture(householdsFuture)
		groupsResult, _ := ArrayFromQueryFuture(groupsFuture)

		context.Respond(GenerateGroups(householdsResult, groupsResult, msg.TargetHouseholdCount))
	case messages.SaveGroups:
		response, _ := context.RequestFuture(state.persistence, persistence.CreateMany{
			EntityType: messages.GroupEntity.String(),
			Entities:   GroupsToHasId(msg.Groups),
			Upsert:     true,
		}, timeout).Result()
		context.Respond(response)

	case messages.Query:
		if msg.Entity != messages.GroupEntity {
			return
		}

		response, _ := context.RequestFuture(state.persistence, persistence.Query{
			EntityType: messages.GroupEntity.String(),
			Model:      func() persistence.HasId { return &models.GroupImpl{new(models.Group)} },
		}, timeout).Result()
		context.Respond(response.(persistence.Query).Entities)

	case messages.DeleteGroups:
		ids := make([][]byte, len(msg.Groups))
		for i, group := range msg.Groups {
			ids[i] = (&models.GroupImpl{Group: group}).GetId()
		}

		_, err := context.RequestFuture(state.persistence, persistence.Delete{
			Ids:        ids,
			EntityType: messages.GroupEntity.String(),
		}, timeout).Result()
		context.Respond(err)

	case *messages.PIDEnvelope:
		switch msg.Type() {
		case messages.PersistencePID:
			state.persistence = msg.PID()
		}

	case *actor.Started:
		context.Request(context.Parent(), messages.NewPIDEnvelope(messages.PersistencePID, nil))
	}
}

func GenerateGroups(householdsResult, groupsResult []persistence.HasId, targetHouseholdCount uint8) []*models.Group {
	historicalGroups := make([]*models.Group, len(groupsResult))
	for i, g := range groupsResult {
		historicalGroups[i] = g.(*models.GroupImpl).Group
	}
	households := []*models.Household{}
	for _, hh := range householdsResult {
		household := hh.(*models.HouseholdImpl).Household
		if household.GetActive() {
			households = append(households, household)
		}
	}

	// low scoring hosts should be picked for new groups
	hostScores := ScoreHosts(GetHostsFromHouseholds(households), historicalGroups)
	By(scoreAsc).Sort(hostScores)

	groupCount := len(households) / int(targetHouseholdCount)
	if float32(len(households)%int(targetHouseholdCount)) > .5*float32(targetHouseholdCount) {
		groupCount++
	}

	date := &timestamp.Timestamp{Seconds: utils.SecondOfTheMonth(time.Now()).Unix()}
	groups := make([]*models.Group, groupCount)
	for i := range groups {
		if i >= len(hostScores) {
			break
		}
		groups[i] = &models.Group{
			HostId:       hostScores[i].Id,
			DateAssigned: date,
			HouseholdIds: [][]byte{hostScores[i].Id},
		}

		households = FilterHouseholds(households, []models.HouseholdFilter{
			func(hh *models.Household) bool {
				return string(hh.GetId()) != string(groups[i].GetHostId())
			},
		})
	}

	for i := range groups {
		// score other households against hosts
		otherHHScores := ScoreGroup(groups[i].GetHostId(), households, historicalGroups)
		By(scoreAsc).Sort(otherHHScores)

		if len(households) >= int(targetHouseholdCount)-1 {
			for _, score := range otherHHScores[:targetHouseholdCount-1] {
				groups[i].HouseholdIds = append(groups[i].HouseholdIds, score.Id)
			}
		} else if float32(len(households)) >= .5*float32(targetHouseholdCount) {
			for _, score := range otherHHScores {
				groups[i].HouseholdIds = append(groups[i].HouseholdIds, score.Id)
			}
		}

		households = FilterHouseholds(households, []models.HouseholdFilter{
			func(hh *models.Household) bool {
				notFound := true
				for _, hhId := range groups[i].GetHouseholdIds() {
					if string(hh.GetId()) == string(hhId) {
						notFound = false
						break
					}
				}
				return notFound
			},
		})
	}

	for i := 0; len(households) > 0; i = (i + 1) % len(groups) {
		groups[i].HouseholdIds = append(groups[i].HouseholdIds, households[0].GetId())
		if len(households) > 1 {
			households = households[1:]
		} else {
			households = []*models.Household{}
		}
	}

	return groups
}

// ScoreHosts returns a list of hosts with a score attached low scores should be selected to host next
func ScoreHosts(hosts [][]byte, groups []*models.Group) []Score {
	// initialize map
	scoreMap := map[string]int{}
	for _, host := range hosts {
		scoreMap[string(host)] = 0
	}

	// score hosts
	const maxScore int = 120
	now := time.Now().Unix()
	for _, group := range groups {
		daysDiff := time.Duration(now-group.GetDateAssigned().GetSeconds()) * time.Second % (time.Hour * 24)
		instanceScore := maxScore - int(daysDiff)

		if instanceScore <= 0 {
			instanceScore = 1
		}

		scoreMap[string(group.GetHostId())] += instanceScore
	}

	// convert to array of scores
	scores := []Score{}
	for hostId, score := range scoreMap {
		scores = append(scores, Score{Id: []byte(hostId), Score: score})
	}

	return scores
}

// ScoreGroup returns a list of householdIds with a score attached low scores should be selected to be grouped with targetHousehold
// targetHousehold should not exist in households list
func ScoreGroup(targetHousehold []byte, households []*models.Household, groups []*models.Group) []Score {
	// initialize map
	scoreMap := map[string]int{}
	for _, hh := range households {
		scoreMap[string(hh.GetId())] = 0
	}

	// score hosts
	const maxScore int = 120
	now := time.Now().Unix()
	for _, group := range groups {
		found := false
		for _, hh := range group.GetHouseholdIds() {
			if string(targetHousehold) == string(hh) {
				found = true
				break
			}
		}
		if !found {
			continue
		}

		daysDiff := time.Duration(now-group.GetDateAssigned().GetSeconds()) * time.Second % (time.Hour * 24)
		instanceScore := maxScore - int(daysDiff)

		if instanceScore <= 0 {
			instanceScore = 1
		}

		for _, hh := range group.GetHouseholdIds() {
			if _, exists := scoreMap[string(hh)]; exists {
				scoreMap[string(hh)] += instanceScore
			}
		}
	}

	// convert to array of scores
	scores := []Score{}
	for hhId, score := range scoreMap {
		scores = append(scores, Score{Id: []byte(hhId), Score: score})
	}

	return scores
}

func FilterHouseholds(households []*models.Household, filters []models.HouseholdFilter) []*models.Household {
	filtered := []*models.Household{}
hhLoop:
	for _, hh := range households {
		for _, filter := range filters {
			if !filter(hh) {
				continue hhLoop
			}
		}
		filtered = append(filtered, hh)
	}
	return filtered
}

func GetHostsFromHouseholds(households []*models.Household) [][]byte {
	hosts := [][]byte{}
	for _, hh := range households {
		if hh.GetActive() && hh.GetHost() {
			hosts = append(hosts, hh.GetId())
		}
	}

	return hosts
}

func ArrayFromQueryFuture(future *actor.Future) ([]persistence.HasId, error) {
	result, err := future.Result()
	if err != nil {
		return nil, err
	}

	return result.(persistence.Query).Entities, nil
}

type Score struct {
	Id    []byte
	Score int
}

func scoreAsc(a, b *Score) bool {
	return a.Score < b.Score
}

type By func(p1, p2 *Score) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by By) Sort(scores []Score) {
	ps := &scoreSorter{
		scores: scores,
		by:     by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(ps)
}

type scoreSorter struct {
	scores []Score
	by     func(p1, p2 *Score) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *scoreSorter) Len() int {
	return len(s.scores)
}

// Swap is part of sort.Interface.
func (s *scoreSorter) Swap(i, j int) {
	s.scores[i], s.scores[j] = s.scores[j], s.scores[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *scoreSorter) Less(i, j int) bool {
	return s.by(&s.scores[i], &s.scores[j])
}

func GroupsToHasId(groups []*models.Group) []persistence.HasId {
	hasIds := make([]persistence.HasId, len(groups))
	for i, group := range groups {
		hasIds[i] = &models.GroupImpl{group}
	}

	return hasIds
}
