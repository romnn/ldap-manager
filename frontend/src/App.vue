<template>
  <div id="app">
    <div>
      <b-navbar toggleable="sm" size="sm" type="dark" variant="dark">
        <router-link :to="{ name: 'HomeRoute' }">
          <b-navbar-brand>LDAP Manager</b-navbar-brand>
        </router-link>

        <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>

        <b-collapse id="nav-collapse" is-nav>
          <b-navbar-nav>
            <b-nav-item
              :to="{
                name: 'EditAccountRoute',
                params: { username: activeUser }
              }"
              >My account</b-nav-item
            >
            <b-nav-item :to="{ name: 'AccountsRoute' }">Accounts</b-nav-item>
            <b-nav-item :to="{ name: 'GroupsRoute' }">Groups</b-nav-item>
          </b-navbar-nav>

          <b-navbar-nav class="ml-auto">
            <!--
            <b-nav-form>
              <b-form-input size="sm" class="mr-sm-2" placeholder="Search"></b-form-input>
              <b-button size="sm" class="my-2 my-sm-0" type="submit">Search</b-button>
            </b-nav-form>
            -->
            <b-nav-item right href="https://github.com/romnnn/ldap-manager">{{
              version
            }}</b-nav-item>
            <!-- TODO: Add internationalization
            <b-nav-item-dropdown text="Lang" right>
              <b-dropdown-item href="#">EN</b-dropdown-item>
              <b-dropdown-item href="#">ES</b-dropdown-item>
              <b-dropdown-item href="#">RU</b-dropdown-item>
              <b-dropdown-item href="#">FA</b-dropdown-item>
            </b-nav-item-dropdown>
            -->

            <b-nav-item-dropdown right>
              <template v-slot:button-content>
                <em>{{ activeUser }} </em>
              </template>
              <b-dropdown-item
                :to="{
                  name: 'EditAccountRoute',
                  params: { username: activeUser }
                }"
                >My account</b-dropdown-item
              >
              <b-dropdown-item :to="{ name: 'LogoutRoute' }"
                >Logout</b-dropdown-item
              >
            </b-nav-item-dropdown>
          </b-navbar-nav>
        </b-collapse>
      </b-navbar>
    </div>
    <div class="app-content">
      <b-breadcrumb :items="items"></b-breadcrumb>
      <router-view />
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import "bootstrap/dist/css/bootstrap.css";
import "bootstrap-vue/dist/bootstrap-vue.css";

@Component({
  components: {}
})
export default class App extends Vue {
  activeUser = "TODO";
  version = "v0.0.1";
  items = [
    {
      text: "Admin",
      href: "#"
    },
    {
      text: "Manage",
      href: "#"
    },
    {
      text: "Library",
      active: true
    }
  ];
}
</script>

<style lang="sass" scoped>
#app
  font-family: Avenir, Helvetica, Arial, sans-serif
  -webkit-font-smoothing: antialiased
  -moz-osx-font-smoothing: grayscale
  text-align: center
  color: #2c3e50

.app-content
  position: relative
  top: 70px
  min-width: 600px
  width: 70%
  margin: 0 auto
</style>
