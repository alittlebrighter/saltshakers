package models

import (
	"encoding/binary"
	"encoding/json"
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
