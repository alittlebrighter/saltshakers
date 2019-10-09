<template>
  <div>
    <div v-for="group in groups" class="group-box" :key="group.host_id">
    <draggable v-model="group.household_ids">
    <transition-group>
        <div v-for="element in group.household_ids" :key="element">
            {{element}}
        </div>
    </transition-group>
  </draggable>
    </div>
  </div>
</template>

<script>
import draggable from "vuedraggable";

export default {
  name: "EditGroups",
  props: {
    id: null
  },
  components: {
      draggable
  },
  data() {
    self = this;

    backend.WailsActor.Request(JSON.stringify({type: "GenerateGroups", payload: {targetHouseholdCount: 4}}))
      .then(result => {
        self.groups = JSON.parse(result);
        console.log(self.groups);
      })
      .catch(err => {
        console.log("error generating groups:", err);
      });

    const data = {
      groups: [],
      editable: true,
      isDragging: false,
      delayedDragging: false
    }

    return data;
  },
  methods: {
    save(hh) {
      var self = this;
      backend.WailsActor.Request(JSON.stringify({type: "SaveGroups", payload: hh}))
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