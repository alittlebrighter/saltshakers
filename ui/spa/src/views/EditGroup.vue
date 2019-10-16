<template>
<div :v-if="Object.keys(households).length > 0">
  <ul v-for="group in groups" :key="group.host_id">
    <li>Host: {{households[group.host_id].surname}}</li>
  <draggable class="group-box" v-bind="dragOptions" :move="onMove" @start="isDragging=true" @end="isDragging=false">
    <transition-group>
      <li v-for="hhId in getNonHosts(group)" :key="hhId">
          {{households[hhId].surname}}
      </li>
    </transition-group>
  </draggable>
  </ul>
</div>
</template>

<script>
import draggable from "vuedraggable";

import store from '../store/store';
import { generateGroups, getHouseholds } from '../store/actions';

export default {
  name: "EditGroups",
  props: {
    id: null
  },
  components: {
      draggable
  },
  data() {
    const self = this;
    store.subscribe(() => {
      const state = store.getState();

      self.households = state.households;
      self.groups = state.activeGroups;
    });

    if (Object.keys(store.getState().households).length === 0) {
      store.dispatch(getHouseholds());
    }
    this.generate(4);

    const data = {
      groups: [],
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
      // saveGroups action
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