<script setup lang="ts">
import axios from "axios";
import { ref, computed, onMounted } from "vue";
import { useAuthStore } from "./stores/auth";
import { useAppStore } from "./stores/app";
import { RouterLink, RouterView } from "vue-router";

const isLoggingOut = ref(false);

function logout() {
  isLoggingOut.value = true;
  const auth = useAuthStore();
  auth.logout();
  setTimeout(() => {
    isLoggingOut.value = false;
  }, 1000);
}

const activeIsAdmin = computed(() => {
  const auth = useAuthStore();
  return auth.activeIsAdmin;
});

const activeUsername = computed(() => {
  const auth = useAuthStore();
  return auth.activeUsername;
});

const activeDisplayName = computed(() => {
  const auth = useAuthStore();
  return auth.activeDisplayName;
});

const pendingConfirmation = computed(() => {
  const app = useAppStore();
  return app.pendingConfirmation;
});

function cancelConfirmation() {
  const app = useAppStore();
  app.cancelConfirmation();
}

function confirmConfirmation() {
  const app = useAppStore();
  app.confirmConfirmation();
}

const version = computed(() => {
  return import.meta.env.STABLE_VERSION;
});

onMounted(() => {
  const auth = useAuthStore();
  axios.defaults.headers.common["x-user-token"] = auth.authToken;
});
</script>

<template>
  <div id="app">
    <div>
      <b-navbar toggleable="sm" size="sm" type="dark" variant="dark">
        <RouterLink :to="{ name: 'HomeRoute' }">
          <span class="title">LDAP Manager</span>
        </RouterLink>

        <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>

        <b-collapse class="navbar" id="nav-collapse" is-nav>
          <b-navbar-nav v-if="activeUsername !== null">
            <!--
            <b-nav-item
              :to="{
                name: 'EditAccountRoute',
                params: { username: activeUsername },
              }"
              >My account</b-nav-item
            >
            <b-nav-item v-if="activeIsAdmin" :to="{ name: 'AccountsRoute' }"
              >Accounts</b-nav-item
            >
            <b-nav-item v-if="activeIsAdmin" :to="{ name: 'GroupsRoute' }"
              >Groups</b-nav-item
            >
            -->
          </b-navbar-nav>

          <b-navbar-nav>
            <b-nav-item right href="https://github.com/romnn/ldap-manager">{{
              version
            }}</b-nav-item>

            <b-nav-item-dropdown right v-if="activeUsername !== null">
              <template v-slot:button-content>
                <em>{{ activeDisplayName }} </em>
              </template>
              <!--
              <b-dropdown-item
                :to="{
                  name: 'EditAccountRoute',
                  params: { username: activeUsername },
                }"
                >My account</b-dropdown-item
              >
              -->
              <b-dropdown-item @click="logout">Logout</b-dropdown-item>
            </b-nav-item-dropdown>
          </b-navbar-nav>
        </b-collapse>
      </b-navbar>
    </div>
    <div class="app-content-container">
      <div class="app-content">
        <div v-if="pendingConfirmation !== null">
          <!--
          <confirmation-c
            :message="pendingConfirmation.message"
            :ackMessage="pendingConfirmation.ack"
            v-on:cancel="cancelConfirmation"
            v-on:confirm="confirmConfirmation"
          ></confirmation-c>
          -->
        </div>
        <div class="logout-container" v-if="isLoggingOut">
          <p>You are being logged out...</p>
          <p><b-spinner label="Logging out..."></b-spinner></p>
        </div>
        <div v-else>
          <!--
          <b-breadcrumb v-if="items.length > 0" :items="items"></b-breadcrumb>
          -->
          <RouterView />
        </div>
      </div>
    </div>
    <div class="footer">
      <a href="https://github.com/romnn/ldap-manager"
        >LDAPManager {{ version }}</a
      >
    </div>
  </div>
</template>

<style lang="sass" scoped>
.app-content-container
  display: flex
  .app-content
    padding: 50px 0
    max-width: 1000px
    width: 90%
    margin: 0 auto

.footer
  display: flex
  padding: 20px 0
  justify-content: center
  font-size: 0.7rem
  color: gray

.navbar
  overflow-x: hidden
  .title
    font-weight: 600
    color: white
</style>
