<template>
  <div>
    <h1 class="title">{{ title }}</h1>
    <router-link to="/households">Back to households</router-link>
    <form class="grid-container">
      <div class="grid-x grid-padding-x">
        <label class="medium-10 cell">
          Last Name
          <input v-model="household.surname" type="text" placeholder="Surname" />
        </label>
        <label v-for="(firstName, i) in household.members" :key="i" class="medium-3 cell">
          First Name
          <input
            v-model="household.members[i].given_name"
            type="text"
            placeholder="Surname"
          />
        </label>
        <button @click="addHouseholdMember()" class="button cell">Add Household Member</button>
        <div class="cell"></div>
        <label class="medium-6 cell">
          Email
          <input v-model="household.email" type="text" placeholder="Email" />
        </label>
        <label class="medium-6 cell">
          Phone Number
          <input v-model="household.phone" type="tel" placeholder="###-###-####" />
        </label>
        <fieldset class="medium-6 cell">
          <input v-model="household.active" type="checkbox" id="active" />
          <label for="active">Active?</label>
          <input v-model="household.host" type="checkbox" id="can-host" />
          <label for="can-host">Can Host?</label>
        </fieldset>
        <fieldset class="medium-6 cell">
          <legend>Communication preference:</legend>
          <input
            v-if="household.email.length > 0"
            v-model="household.preferred_contact"
            type="radio"
            value="0"
            id="email"
          />
          <label v-if="household.email.length > 0" for="email">Email</label>
          <input
            v-if="household.phone.length > 0"
            v-model="household.preferred_contact"
            type="radio"
            value="1"
            id="text"
          />
          <label v-if="household.phone.length > 0" for="text">Text (SMS)</label>
          <input
            v-if="household.phone.length > 0"
            v-model="household.preferred_contact"
            type="radio"
            value="2"
            id="call"
          />
          <label v-if="household.phone.length > 0" for="call">Call</label>
        </fieldset>
        <div class="button-group">
          <button @click="save(household)" class="button">
            <i class="fad fa-save"></i>&nbsp;Save
          </button>
          <button @click="remove($route.params.id)" class="alert button">
            <i class="fad fa-trash-alt"></i>&nbsp;Delete
          </button>
        </div>
      </div>
    </form>
  </div>
</template>

<script>
import store from "../store/store";
import { select } from "../store/store";
import { saveHousehold, deleteHousehold, getHousehold } from "../store/actions";

export default {
  props: {
    id: null
  },
  data() {
    const emptyHH = {
      id: "",
      surname: "",
      members: [{ given_name: "" }],
      email: "",
      phone: "",
      host: false,
      active: true,
      preferred_contact: "0"
    };

    if (this.$route.params.id !== "null") {
      const self = this;
      store.subscribe(() => {
        self.household = Object.assign(
          emptyHH,
          select("households", self.$route.params.id)
        );
      });

      store.dispatch(getHousehold(this.$route.params.id));
    }

    const data = {
      title: "Add/Edit Household",
      household: emptyHH
    };
    return data;
  },
  methods: {
    addHouseholdMember() {
      this.household.members.push({ first_name: "" });
    },
    save(hh) {
      hh.preferred_contact = parseInt(hh.preferred_contact);
      hh.phone = parseInt(hh.phone.replace(/[-()\s\t]/g, ""));

      const actualMembers = [];
      for (var i = 0; i < hh.members.length; i++) {
        if (hh.members[i].given_name.length) {
          actualMembers.push(hh.members[i]);
        }
      }
      hh.members = actualMembers;

      store.dispatch(saveHousehold(hh));
    },
    remove(id) {
      store.dispatch(deleteHousehold(id));
    }
  },
  components: {}
};
</script>

<style>
.logo {
  width: 16em;
}

.title {
  text-align: center;
}
</style>