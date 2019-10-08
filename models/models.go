package models

import (
	"encoding/binary"
	"encoding/json"
	"sort"
)

type HouseholdImpl struct {
	*Household
}

func (hh *HouseholdImpl) SetId(id []byte) {
	hh.Id = id
}

func (hh *HouseholdImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(hh.Household)
}

func (hh *HouseholdImpl) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, hh.Household)
}

type HouseholdFilter func(hh *Household) bool

type GroupImpl struct {
	*Group
}

func (g *GroupImpl) GetId() []byte {
	seconds := make([]byte, 8)
	binary.BigEndian.PutUint64(seconds, uint64(g.DateAssigned.GetSeconds()))
	id := append(seconds, g.GetHostId()...)
	for _, hh := range g.GetHouseholdIds() {
		id = append(id, hh...)
	}
	return id
}

func (g *GroupImpl) SetId(id []byte) {
	return // no op
}

func (g *GroupImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(g.Group)
}

func (g *GroupImpl) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, g.Group)
}

type By func(p1, p2 *Group) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by By) Sort(groups []Group) {
	ps := &groupSorter{
		groups: groups,
		by:     by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(ps)
}

type groupSorter struct {
	groups []Group
	by     func(p1, p2 *Group) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *groupSorter) Len() int {
	return len(s.groups)
}

// Swap is part of sort.Interface.
func (s *groupSorter) Swap(i, j int) {
	s.groups[i], s.groups[j] = s.groups[j], s.groups[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *groupSorter) Less(i, j int) bool {
	return s.by(&s.groups[i], &s.groups[j])
}
