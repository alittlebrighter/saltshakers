<template>
<div>
  <draggable v-for="group in groups" tag="div" class="group-box" :key="group.host_id" v-bind="dragOptions" :move="onMove" @start="isDragging=true" @end="isDragging=false">
    <transition-group>
      <div v-for="hhId in group.household_ids" :key="hhId">
          {{households[hhId].surname}}
      </div>
    </transition-group>
  </draggable>
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
    save(groups) {
      var self = this;
      backend.WailsActor.Request(JSON.stringify({type: "SaveGroups", payload: groups}))
        .then(result => {
          self.groups = JSON.parse(result);
        });
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