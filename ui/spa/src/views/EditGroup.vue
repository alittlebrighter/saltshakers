<template>
  <div class="main grid-x grid-margin-x">
    <div class="sidebar cell medium-2">
      <ul>
        <li @click="generate(targetSize)">
          Generate Groups
          <label>
            Size
            <input v-model="targetSize" type="number" class="small" id="size" />
          </label>
        </li>
        <li>
          <i class="fas fa-ellipsis-h"></i>
        </li>
        <li
          v-for="(groups, date_assigned) in savedGroups"
          :key="date_assigned.seconds"
          @click="showGroups(groups)"
        >{{ secondsToMonthYear(date_assigned) }}</li>
      </ul>
    </div>
    <div class="cell medium-10 grid-x grid-margin-x">
      <h1 class="cell">{{title}}</h1>

      <h3 v-if="!canSave" class="cell">{{ secondsToMonthYear(groupsTimestamp.seconds) }}</h3>

      <div
        v-for="group in groups"
        :key="group.host_id"
        class="card cell medium-3"
        style="width: 300px;"
      >
        <div class="card-divider">
          <b>Host:</b>
          &nbsp;
          {{households[group.host_id] ? households[group.host_id].surname : group.hose_id}}&nbsp;
          <button
            @click="deleteGroup(group)"
            v-if="!canSave"
            class="button hollow warn float-right"
          >
            <i class="fad fa-trash"></i>
          </button>
        </div>
        <div class="card-section">
          <draggable
            tag="ul"
            class="group-box"
            v-bind="dragOptions"
            :move="onMove"
            @start="isDragging=true"
            @end="isDragging=false"
          >
            <transition-group>
              <li
                :class="{grabbable: canSave}"
                v-for="hhId in getNonHosts(group)"
                :key="hhId"
              >{{households[hhId] ? households[hhId].surname : hhId}}</li>
            </transition-group>
          </draggable>
        </div>
      </div>
      <button v-if="canSave" @click="save(groups)" class="cell button">Save</button>
    </div>
  </div>
</template>

<script>
import draggable from "vuedraggable";

import store from "../store/store";
import {
  generateGroups,
  getHouseholds,
  saveGroups,
  getGroups,
  deleteGroup
} from "../store/actions";

export default {
  name: "EditGroups",
  components: {
    draggable
  },
  data() {
    const self = this;
    store.subscribe(() => {
      const state = store.getState();

      self.households = state.households;
      self.groups = state.activeGroups;
      self.canSave = true;

      self.savedGroups = {};
      for (var i = 0; i < state.groups.length; i++) {
        if (self.savedGroups[state.groups[i].date_assigned.seconds]) {
          self.savedGroups[state.groups[i].date_assigned.seconds].push(
            state.groups[i]
          );
        } else {
          self.savedGroups[state.groups[i].date_assigned.seconds] = [
            state.groups[i]
          ];
        }
      }
    });

    if (Object.keys(store.getState().households).length === 0) {
      store.dispatch(getHouseholds());
      store.dispatch(getGroups());
    }

    const data = {
      title: "Groups",
      targetSize: 4,
      groups: [],
      savedGroups: {},
      households: {},
      editable: true,
      isDragging: false,
      delayedDragging: false,
      canSave: true,

      months: [
        "",
        "January",
        "February",
        "March",
        "April",
        "May",
        "June",
        "July",
        "August",
        "September",
        "October",
        "November",
        "December"
      ]
    };
    this.generate(data.targetSize);

    return data;
  },
  methods: {
    showGroups(groups) {
      this.groups = groups;
      this.canSave = false;
    },
    secondsToDate(seconds) {
      return new Date(seconds * 1000);
    },
    generate(hhCount) {
      store.dispatch(generateGroups(hhCount));
    },
    getNonHosts(group) {
      return group.household_ids.filter(hh => hh !== group.host_id);
    },
    secondsToMonthYear(seconds) {
      const date = this.secondsToDate(seconds);
      return this.months[date.getMonth()] + " " + date.getFullYear();
    },
    save(groups) {
      store.dispatch(saveGroups(groups));
    },
    deleteGroup(group) {
      store.dispatch(deleteGroup(group));
    },
    onMove({ relatedContext, draggedContext }) {
      const relatedElement = relatedContext.element;
      const draggedElement = draggedContext.element;
      return (
        (!relatedElement || !relatedElement.fixed) && !draggedElement.fixed
      );
    }
  },
  computed: {
    groupsTimestamp() {
      return this.groups && this.groups.length > 0
        ? this.groups[0].date_assigned
        : Math.floor(Date.now() / 1000);
    },
    dragOptions() {
      return {
        animation: 0,
        group: "description",
        disabled: !this.editable,
        ghostClass: "ghost"
      };
    }
  },
  watch: {
    isDragging(newValue) {
      if (newValue) {
        this.delayedDragging = true;
        return;
      }
      this.$nextTick(() => {
        this.delayedDragging = false;
      });
    }
  }
};
</script>

<style>
.main {
  height: 100%;
}

.sidebar {
  overflow: auto;
  border-right: solid 1px #000;
  text-align: center;
}

.sidebar > ul > li {
  list-style: none;
}

.grabbable {
  list-style: none;
  text-align: center;
  border: solid 1px gray;
  border-radius: 5px;
  cursor: grab;
}

.small {
  display: inline;
  width: 50px;
}
</style>