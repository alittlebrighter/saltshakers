<template>
<div>
  <h1>{{title}}</h1>
  <h3>{{groups.length > 0 ? new Date(groups[0].date_assigned.seconds * 1000).toString() : "unknown"}}</h3>
  <ul :v-if="Object.keys(households).length > 0" v-for="group in groups" :key="group.host_id">
    <li>Host: {{households[group.host_id] ? households[group.host_id].surname : group.hose_id}}</li>
  <draggable class="group-box" v-bind="dragOptions" :move="onMove" @start="isDragging=true" @end="isDragging=false">
    <transition-group>
      <li v-for="hhId in getNonHosts(group)" :key="hhId">
          {{households[hhId] ? households[hhId].surname : hhId}}
      </li>
    </transition-group>
  </draggable>
  </ul>
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

      if (self.savedGroups) {
        console.log("saved:", self.savedGroups);
      }
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
      delayedDragging: false
    }
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
    }
  },
  computed: {
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
.group-box {
  padding: 20px;
}
</style>