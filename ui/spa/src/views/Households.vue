<template>
  <div>
    <h1>{{ title }}</h1>
    <router-link to="/add-household/null">Add Household</router-link>
    <ul>
      <li v-for="hh in households" :key="hh.id">
        <router-link :to="'/add-household/' + hh.id">
          {{hh.surname || "unknown"}}
        </router-link>
      </li>
    </ul>
  </div>
</template>

<script>
import store from '../store/store';
import { select } from '../store/store';
import { getHouseholds } from '../store/actions';

export default {
  data() {
    const self = this;
    store.subscribe(() => {
      // TODO: figure out why we need to filter for undefined
      self.households = Object.values(select("households")).filter(hh => hh.id !== undefined);
    });

    this.getHouseholds();

    return {
      title: "Households",
      households: [],
    };
  },
  methods: {
    getHouseholds: () => store.dispatch(getHouseholds())
  },
  components: {}
};
</script>

<style>
.logo {
  width: 16em;
}
</style>