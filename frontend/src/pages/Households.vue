<template>
  <div>
    <h1>{{ title }}</h1>
    <router-link to="/add-household/null">Add Household</router-link>
    <ul>
      <li v-for="hh in households" :key="hh.id">
        <router-link :to="'/add-household/' + hh.id">
          {{hh.surname}}
        </router-link>
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  data() {
    const self = this;

    backend.WailsActor.Request(JSON.stringify({type: "QueryHouseholds", payload: {filters: []}}))
      .then(function(toParse) {
        self.households = JSON.parse(toParse);
      });

    return {
      title: "Households",
      households: [],
    };
  },
  methods: {
    
  },
  components: {}
};
</script>

<style>
.logo {
  width: 16em;
}
</style>