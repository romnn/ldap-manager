<template>
  <div id="app">
    <div>
      <b-navbar toggleable="sm" size="sm" type="dark" variant="dark">
        <router-link :to="{ name: 'HomeRoute' }">
          <b-navbar-brand>LDAP Manager</b-navbar-brand>
        </router-link>

        <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>

        <b-collapse id="nav-collapse" is-nav>
          <b-navbar-nav v-if="activeUsername !== null">
            <b-nav-item
              :to="{
                name: 'EditAccountRoute',
                params: { username: activeUsername }
              }"
              >My account</b-nav-item
            >
            <b-nav-item :to="{ name: 'AccountsRoute' }">Accounts</b-nav-item>
            <b-nav-item :to="{ name: 'GroupsRoute' }">Groups</b-nav-item>
          </b-navbar-nav>

          <b-navbar-nav class="ml-auto">
            <b-nav-item right href="https://github.com/romnnn/ldap-manager">{{
              version
            }}</b-nav-item>

            <b-nav-item-dropdown right v-if="activeUsername !== null">
              <template v-slot:button-content>
                <em>{{ activeDisplayName }} </em>
              </template>
              <b-dropdown-item
                :to="{
                  name: 'EditAccountRoute',
                  params: { username: activeUsername }
                }"
                >My account</b-dropdown-item
              >
              <b-dropdown-item @click="logout">Logout</b-dropdown-item>
            </b-nav-item-dropdown>
          </b-navbar-nav>
        </b-collapse>
      </b-navbar>
    </div>
    <div class="app-content">
      <div v-if="pendingConfirmation !== null">
        <confirmation-c
          :message="pendingConfirmation.message"
          :ackMessage="pendingConfirmation.ack"
          v-on:cancel="cancelConfirmation"
          v-on:confirm="confirmConfirmation"
        ></confirmation-c>
      </div>
      <div class="logout-container" v-if="isLoggingOut">
        <p>You are being logged out...</p>
        <p><b-spinner label="Logging out..."></b-spinner></p>
      </div>
      <div v-else>
        <b-breadcrumb v-if="items.length > 0" :items="items"></b-breadcrumb>
        <router-view />
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import "bootstrap/dist/css/bootstrap.css";
import "bootstrap-vue/dist/bootstrap-vue.css";
import { Dictionary } from "vue-router/types/router";
import { AppModule } from "./store/modules/app";
import ConfirmationC from "./components/Confirmation.vue";
import { AuthModule } from "./store/modules/auth";

export interface BreadcrumbItem {
  text: string;
  active?: boolean;
  to?: { name?: string; params?: Dictionary<string> };
}

@Component({
  components: { ConfirmationC }
})
export default class App extends Vue {
  protected isLoggingOut = false;

  logout() {
    this.isLoggingOut = true;
    setTimeout(() => {
      AuthModule.logout();
      this.isLoggingOut = false;
    }, 1000);
  }

  get activeUsername() {
    return AuthModule.activeUsername;
  }

  get activeDisplayName() {
    return AuthModule.activeDisplayName;
  }

  get pendingConfirmation() {
    return AppModule.pendingConfirmation;
  }

  cancelConfirmation() {
    AppModule.cancelConfirmation();
  }

  confirmConfirmation() {
    AppModule.confirmConfirmation();
  }

  get version() {
    return process.env.STABLE_VERSION;
  }

  get items(): BreadcrumbItem[] {
    if (!(this.$route.meta?.showBreadcrumb ?? true)) return [];
    const base = this.$route?.meta?.base ?? [];
    const params = this.$route?.params ?? {};
    const paramsItems = Object.values(params).reduce((acc, param) => {
      const name = this.$route.name;
      if (name)
        acc.push({
          text: param,
          to: { name: name, params: this.$route.params },
          active: true
        });
      return acc;
    }, [] as BreadcrumbItem[]);
    return base.concat(paramsItems);
  }

  mounted() {
    Vue.axios.defaults.headers.common["x-user-token"] = AuthModule.authToken;
  }
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
  padding-bottom: 70px
  min-width: 600px
  width: 70%
  margin: 0 auto
</style>
