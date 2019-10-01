<template>
  <v-layout justify-center align-center class="px-0">
    <v-row align="center">
      <h1>{{ title }}</h1>
    </v-row>
    <v-form>
      <v-text-field v-model="household.surname" />
    </v-form>
  </v-layout>
</template>

<script>
export default {
  data: () => ({
    title: "Add/Edit Household",
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