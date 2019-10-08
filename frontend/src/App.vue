<template>
  <v-app id="inspire">
    <v-navigation-drawer v-model="drawer" clipped fixed app>
      <v-list dense>
        <router-link v-for="(action, i) in actions" :key="i" :to="action.route"> 
          <v-list-tile>
            <v-list-tile-action>
              <v-icon>{{ action.icon }}</v-icon>
            </v-list-tile-action>
            <v-list-tile-content>
              <v-list-tile-title>{{ action.label }}</v-list-tile-title>
            </v-list-tile-content>
          </v-list-tile>
        </router-link>
      </v-list>
    </v-navigation-drawer>
    <v-toolbar app fixed clipped-left>
      <v-toolbar-side-icon @click.stop="drawer = !drawer"></v-toolbar-side-icon>
      <v-toolbar-title>Saltshakers</v-toolbar-title>
    </v-toolbar>
    <v-content>
      <v-container fluid class="px-0">
        <router-view></router-view>
      </v-container>
    </v-content>
    <v-footer app fixed>
      <span style="margin-left:1em">&copy; </span>
    </v-footer>
  </v-app>
</template>

<script>
//import HelloWorld from "./components/HelloWorld.vue";

export default {
  data: () => {
    backend.WailsActor.Request(JSON.stringify({type: "GenerateGroups", payload: {targetHouseholdCount: 4}}))
      .then(result => {
        var parsed = JSON.parse(result);
        console.log("generated groups", parsed);
      })
      .catch(err => {
        console.log("error generating groups:", err);
      });
    return {
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
  }
  },
  components: {
  },
  props: {
    source: String
  }
};
</script>

<style>
.logo {
  width: 16em;
}

a {
  color: initial;
  
}
</style>