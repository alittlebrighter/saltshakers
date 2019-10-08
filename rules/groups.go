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
			EntityType: GroupEntity.String(),
			Model:      models.Group{},
		}, timeout)
		householdsFuture := context.RequestFuture(context.Parent(), messages.QueryHouseholds{}, timeout)

		groupsResult, _ := ArrayFromQueryFuture(groupsFuture)
		householdsResult, _ := ArrayFromQueryFuture(householdsFuture)

		historicalGroups := make([]models.Group, len(groupsResult))
		for i, g := range groupsResult {
			historicalGroups[i] = g.(models.Group)
		}
		households := []*models.Household{}
		for _, hh := range householdsResult {
			household := hh.(*models.Household)
			if household.GetActive() {
				households = append(households, household)
			}
		}

		// low scoring hosts should be picked for new groups
		hostScores := ScoreHosts(GetHostsFromHouseholds(households), historicalGroups)
		By(scoreAsc).Sort(hostScores)

		groupCount := len(households) / int(msg.TargetHouseholdCount)
		if float32(len(households)%int(msg.TargetHouseholdCount)) > .5*float32(msg.TargetHouseholdCount) {
			groupCount++
		}

		now := &timestamp.Timestamp{Seconds: time.Now().Unix()}
		groups := make([]*models.Group, groupCount)
		for i, group := range groups {
			group.HostId = hostScores[i].Id
			group.DateAssigned = now
		}
		context.Respond(groups)
	case *messages.PIDEnvelope:
		switch msg.Type() {
		case messages.PersistencePID:
			state.persistence = msg.PID()
		}
	case *actor.Started:
		context.Request(context.Parent(), messages.NewPIDEnvelope(messages.PersistencePID, nil))
	}
}

// ScoreHosts returns a list of hosts with a score attached low scores should be selected to host next
func ScoreHosts(hosts [][]byte, groups []models.Group) []Score {
	// initialize map
	scoreMap := map[string]int{}
	for _, host := range hosts {
		scoreMap[string(host)] = 0
	}

	// score hosts
	const maxScore int = 90
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

func ArrayFromQueryFuture(future *actor.Future) ([]interface{}, error) {
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