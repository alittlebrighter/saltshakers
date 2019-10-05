<template>
  <div>
    <h1>{{ title }}</h1>
    <form>
      <input v-model="household.surname" type="text" placeholder="Surname" /><br>
      <input v-model="household.email" type="text" placeholder="Email" />
      <button v-on:click="save()">Save</button>
    </form>
  </div>
</template>

<script>
export default {
  data() {
    return {
      title: "Add/Edit Household",
      household: {
        surname: "",
        email: "",
      }
    };
  },
  methods: {
    save() {
      backend.WailsActor.Request(JSON.stringify({type: "CreateHousehold", payload: this.household}))
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