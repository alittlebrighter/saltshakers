import { createStore, combineReducers, applyMiddleware } from "redux";
import {
  SET_HOUSEHOLDS,
  SET_GROUPS,
  UPSERT_HOUSEHOLD,
  SET_ACTIVE_GROUPS,
  DELETE_HOUSEHOLD,
  DELETE_GROUP,
  ADD_GROUPS
} from "./actions";

const hhHandlers = {};
hhHandlers[SET_HOUSEHOLDS] = (state, action) => {
  for (let i = 0; i < action.payload.length; i++) {
    const hh = action.payload[i];
    state[hh.id] = hh;
  }
  return { ...state };
};

hhHandlers[UPSERT_HOUSEHOLD] = (state, action) => {
  state[action.payload.id] = action.payload;
  return { ...state };
};

hhHandlers[DELETE_HOUSEHOLD] = (state, action) => {
  delete state[action.payload];
  return { ...state };
};

function householdsReducer(state = {}, action) {
  if (hhHandlers[action.type]) {
    return hhHandlers[action.type](state, action);
  }

  return state;
}

const groupsHandlers = {};
groupsHandlers[SET_GROUPS] = (state, action) => {
  return action.payload;
};

groupsHandlers[ADD_GROUPS] = (state, action) => {
  console.log(state);
  state.push(...action.payload);
  return state;
};

groupsHandlers[DELETE_GROUP] = (state, action) => {
  for (var i = 0; i < state.length; i++) {
    if (
      state[i].date_assigned.seconds === action.payload.date_assigned.seconds &&
      state[i].host_id === action.payload.host_id
    ) {
      const first = state.slice(0, i + 1);
      if (i === 0) {
        console.log("new group count", state.slice(1).length);
        return state.slice(1);
      } else if (i + 1 <= state.length) {
        console.log("new group count", first.length);
        return first;
      } else {
        first.push(...state.slice(i + 1));
        console.log("new group count", first.length);
        return first;
      }
    }
  }
  return state;
};

function groupsReducer(state = [], action) {
  if (groupsHandlers[action.type]) {
    return groupsHandlers[action.type](state, action);
  }

  return state;
}

const activeGroupsHandlers = {};
activeGroupsHandlers[SET_ACTIVE_GROUPS] = (state, action) => {
  return action.payload;
};

function activeGroupsReducer(state = [], action) {
  if (activeGroupsHandlers[action.type]) {
    return activeGroupsHandlers[action.type](state, action);
  }

  return state;
}

/**
 * Lets you dispatch promises in addition to actions.
 * If the promise is resolved, its result will be dispatched as an action.
 * The promise is returned from `dispatch` so the caller may handle rejection.
 */
const vanillaPromise = store => next => action => {
  if (typeof action.then !== "function") {
    return next(action);
  }

  return Promise.resolve(action).then(store.dispatch);
};

const store = createStore(
  combineReducers({
    households: householdsReducer,
    groups: groupsReducer,
    activeGroups: activeGroupsReducer
  }),
  applyMiddleware(vanillaPromise)
);

export default store;

export const select = (...args) => {
  if (args.length === 0) {
    return null;
  }

  let val = store.getState()[args[0]];
  args.slice(1).forEach(arg => {
    val = val[arg];
  });
  return val;
};
