<template>
  <div id="app">
    <ul id="nav">
      <li v-for="action in actions" :key="action.route">
        <router-link :to="action.route">{{ action.label }}</router-link>
      </li>
    </ul>
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
#app {
  font-family: "Avenir", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}
#nav {
  padding: 30px;
}

#nav a {
  font-weight: bold;
  color: #2c3e50;
}

#nav a.router-link-exact-active {
  color: #42b983;
}
</style>
