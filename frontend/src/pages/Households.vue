<template>
  <v-layout justify-center align-center class="px-0">
    <h1>{{ title }}</h1>
    <router-link to="/add-household">Add Household</router-link>
  </v-layout>
</template>

<script>
export default {
  data: () => ({
    title: "Households",
    household: {
      surname: "",
    }
  }),
  components: {
  }
};
</script>

<style>
.logo {
  width: 16em;
}
</style>