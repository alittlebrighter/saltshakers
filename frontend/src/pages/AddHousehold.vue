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

    console.log("hh id:", this.$route.params.id)
    if (this.$route.params.id) {
      const self = this;
      backend.WailsActor.Request(JSON.stringify({type: "GetHousehold", payload: {id: this.$route.params.id}}))
        .then(result => {
          console.log("result:", result);
          self.household = JSON.parse(result) || emptyHH;
        })
    }

    const data = {
      title: "Add/Edit Household",
      household: emptyHH
    }

    return data;
  },
  methods: {
    save(hh) {
      backend.WailsActor.Request(JSON.stringify({type: "CreateHousehold", payload: hh}))
        .then(result => {
          console.log("households:", result);
          self.households = result;
        });
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