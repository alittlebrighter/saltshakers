export const SET_HOUSEHOLDS = 'SET_HOUSEHOLDS';
export const setHouseholds = households => ({
    type: SET_HOUSEHOLDS,
    payload: households
});

export const UPSERT_HOUSEHOLD = 'UPSERT_HOUSEHOLD'
export const upsertHousehold = household => ({
    type: UPSERT_HOUSEHOLD,
    payload: household
})

export const SET_GROUPS = 'SET_GROUPS';
export const setGroups = groups => ({
    type: SET_GROUPS,
    payload: groups
});

export const SET_ACTIVE_GROUPS = 'SET_ACTIVE_GROUPS';
export const setActiveGroups = groups => ({
    type: SET_ACTIVE_GROUPS,
    payload: groups
});

export const getHouseholds = () => backend.WailsActor.Request(JSON.stringify({type: "QueryHouseholds", payload: {filters: []}}))
    .then(JSON.parse)
    .then(households => setHouseholds(households));

export const getHousehold = hhId => backend.WailsActor.Request(JSON.stringify({type: "GetHousehold", payload: {id: hhId}}))
    .then(JSON.parse)
    .then(hh => upsertHousehold(hh));

export const saveHousehold = hh => backend.WailsActor.Request(JSON.stringify({type: "CreateHousehold", payload: hh}))
    .then(JSON.parse)
    .then(hh => upsertHousehold(hh));

export const generateGroups = hhCount => backend.WailsActor.Request(JSON.stringify({type: "GenerateGroups", payload: {targetHouseholdCount: hhCount}}))
    .then(JSON.parse)
    .then(groups => setActiveGroups(groups));