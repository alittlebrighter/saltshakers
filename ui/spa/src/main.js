import './../node_modules/foundation-sites/dist/css/foundation.min.css';
import './../node_modules/foundation-sites/dist/js/foundation.min.js';

import Vue from "vue";

import VueRouter from "vue-router";

import Households from "./views/Households.vue";
import AddHousehold from "./views/AddHousehold.vue";
import EditGroup from "./views/EditGroup.vue";

Vue.use(VueRouter);

const router = new VueRouter({
  routes: [
    { path: '/households', component: Households },
    { path: '/add-household/:id', component: AddHousehold, props: { id: null } },
    { path: '/create-groups', component: EditGroup},
  ]
});

import App from "./App.vue";

Vue.config.productionTip = false;
Vue.config.devtools = true;

new Vue({
  router,
  render: h => h(App)
}).$mount("#app");
