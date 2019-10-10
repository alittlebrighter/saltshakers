<template>
  <div>
    <h1>{{ title }}</h1>
    <form>
      <input v-model="household.surname" type="text" placeholder="Surname" /><br>
      <input v-model="household.email" type="text" placeholder="Email" /><br>
      <span>host</span><input v-model="household.host" type="checkbox" />&nbsp;
      <span>active</span><input v-model="household.active" type="checkbox" /><br>
      <button v-on:click="save(household)">Save</button>
    </form>
  </div>
</template>

<script>
import store from '../store/store';
import { select } from '../store/store';
import { saveHousehold, getHousehold } from '../store/actions';

export default {
  props: {
    id: null
  },
  data() {
    const emptyHH = {
        surname: "",
        email: "",
        host: false,
        active: true
      };

    if (this.$route.params.id) {
      const self = this;
      store.subscribe(() => {
        self.household = select("households", self.$route.params.id) || emptyHH;
      })

      store.dispatch(getHousehold(this.$route.params.id));
    }

    const data = {
      title: "Add/Edit Household",
      household: emptyHH
    };
    return data;
  },
  methods: {
    save(hh) {
      store.dispatch(saveHousehold(hh));
    }
  },
  components: {}
};
</script>

<style>
.logo {
  width: 16em;
}
</style>