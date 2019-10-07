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
  props: {
    id: null
  },
  data() {
    const emptyHH = {
      surname: "",
      email: "",
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