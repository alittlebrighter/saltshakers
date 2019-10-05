package models

import "encoding/json"

type HouseholdImpl struct {
	*Household
}

func (hh *HouseholdImpl) SetId(id uint64) {
	hh.Id = id
}

func (hh *HouseholdImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(hh.Household)
}

func (hh *HouseholdImpl) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, hh.Household)
}
