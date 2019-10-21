<template>
  <div id="app">
    <div class="top-bar">
      <div class="top-bar-left">
        <ul id="nav" class="dropdown menu" data-dropdown-menu>
          <li class="menu-text">Saltshakers Household and Group Manager</li>
          <li v-for="action in actions" :key="action.route">
            <router-link :to="action.route">{{ action.label }}</router-link>
          </li>
        </ul>
      </div>
      <div class="top-bar-right"></div>
    </div>
    <router-view/>
  </div>
</template>

<script>
export default {
  data: () => ({
    drawer: false,
    actions: [
      {
        label: "Manage Households",
        icon: "",
        route: "/households"
      },
      {
        label: "Create Groups",
        icon: "",
        route: "/create-groups"
      },
      {
        label: "History",
        icon: "",
        route: "/group-history"
      },
      {
        label: "Settings",
        icon: "",
        route: "/settings"
      }
    ]
  }),
  components: {
  },
  props: {
    source: String
  }
};
</script>

<style>

</style>
