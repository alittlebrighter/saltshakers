<template>
  <div id="app">
    <div class="top-bar">
      <div class="top-bar-left">
        <ul id="nav" class="menu" data-dropdown-menu>
          <li class="menu-text">Saltshakers Household and Group Manager</li>
          <li class="menu-option" v-for="action in actions" :key="action.route">
            <router-link class="inline" :to="action.route">
              <i :class="action.icon"></i>
              {{ action.label }}
            </router-link>
          </li>
        </ul>
      </div>
      <div class="top-bar-right"></div>
    </div>
    <router-view />
  </div>
</template>

<script>
export default {
  data: () => ({
    drawer: false,
    actions: [
      {
        label: "Households",
        icon: "fad fa-home",
        route: "/households"
      },
      {
        label: "Groups",
        icon: "fad fa-user-friends",
        route: "/groups"
      }
      /*
      {
        label: "Settings",
        icon: "",
        route: "/settings"
      }
      */
    ]
  }),
  components: {},
  props: {
    source: String
  }
};
</script>

<style>
.inline {
  display: inline;
}
</style>
