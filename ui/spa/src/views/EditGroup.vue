<template>
<div class="grid-x grid-margin-x">
  <h1 class="cell">{{title}}</h1>
  <h3 class="cell">{{months[jsDate.getMonth()]}} - {{jsDate.getFullYear()}}</h3>
  
  <div class="cell medium-1"></div>
  <div v-for="group in groups" :key="group.host_id" class="card cell medium-3" style="width: 300px;">
    <div class="card-divider">
      <b>Host:</b> {{households[group.host_id] ? households[group.host_id].surname : group.hose_id}}
    </div>
    <div class="card-section">
      <draggable tag="ul" class="group-box" v-bind="dragOptions" :move="onMove" @start="isDragging=true" @end="isDragging=false">
        <transition-group>
          <li v-for="hhId in getNonHosts(group)" :key="hhId">
              {{households[hhId] ? households[hhId].surname : hhId}}
          </li>
        </transition-group>
      </draggable>
    </div>
  </div>

  <button @click="save(groups)">Save</button>
</div>
</template>

<script>
import draggable from "vuedraggable";

import store from '../store/store';
import { generateGroups, getHouseholds, saveGroups, getGroups } from '../store/actions';

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
      self.savedGroups = state.groups;
    });

    if (Object.keys(store.getState().households).length === 0) {
      store.dispatch(getHouseholds());
      store.dispatch(getGroups());
    }
    this.generate(4);

    const data = {
      title: "Groups",
      groups: [],
      savedGroups: [],
      households: {},
      editable: true,
      isDragging: false,
      delayedDragging: false,

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
        "December",
      ],
    };
    return data;
  },
  methods: {
    generate(hhCount) {
      store.dispatch(generateGroups(hhCount));
    },
    getNonHosts(group) {
      return group.household_ids.filter(hh => hh !== group.host_id);
    },
    save(groups) {
      store.dispatch(saveGroups(groups));
    },
    onMove({ relatedContext, draggedContext }) {
      const relatedElement = relatedContext.element;
      const draggedElement = draggedContext.element;
      return (
        (!relatedElement || !relatedElement.fixed) && !draggedElement.fixed
      );
    },
  },
  computed: {
    jsDate() {
      return this.groups && this.groups.length > 0 ? new Date(this.groups[0].date_assigned.seconds * 1000) : new Date();
    },
    dragOptions() {
      return {
        animation: 0,
        group: "description",
        disabled: !this.editable,
        ghostClass: "ghost"
      };
    },
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
</style>