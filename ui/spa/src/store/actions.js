import config from "../config";

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

export const DELETE_HOUSEHOLD = 'DELETE_HOUSEHOLD'
export const removeHousehold = id => ({
    type: DELETE_HOUSEHOLD,
    payload: id
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

export const NO_OP = 'NO_OP';
export const noOp = () => ({
    type: NO_OP,
});

export const getHouseholds = () => fetch(config.apiBaseURL + "/households", {method: "GET", mode: "cors"})
    .then((response) => response.json())
    .then(households => setHouseholds(households));

export const getHousehold = hhId => fetch(config.apiBaseURL + "/households/" + hhId)
    .then(response => response.json())
    .then(hh => upsertHousehold(hh));

export const saveHousehold = hh => fetch(config.apiBaseURL + "/households", {
    method: "POST", 
    headers: {
        'Content-Type': 'application/json'
      },
    mode: "cors", 
    body: JSON.stringify(hh),
})
    .then(JSON.parse)
    .then(hh => upsertHousehold(hh));

export const deleteHousehold = id => fetch(config.apiBaseURL + "/households/" + id, {method: "DELETE"})
    .then(() => removeHousehold(id));

export const generateGroups = hhCount => fetch(config.apiBaseURL + "/groups/generate?targetHouseholdCount=" + hhCount)
    .then(response => response.json())
    .then(groups => setActiveGroups(groups));

export const saveGroups = groups => fetch(config.apiBaseURL + "/groups", {
        method: "POST",
        headers: {
            'Content-Type': 'application/json',
        },
        mode: "cors",
        body: JSON.stringify(groups),
    })
    .then(() => noOp())
    .catch(console.log);

export const getGroups = () => fetch(config.apiBaseURL + "/groups")
    .then(response => response.json())
    .then(groups => setGroups(groups))