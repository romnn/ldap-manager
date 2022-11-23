<script setup lang="ts">
import axios from "axios";
import { ref, computed, onMounted } from "vue";
import { useAuthStore } from "./stores/auth";
import { useAppStore } from "./stores/app";
import { RouterLink, RouterView } from "vue-router";

const authStore = useAuthStore();
const appStore = useAppStore();

const isLoggingOut = ref(false);

function logout() {
  isLoggingOut.value = true;
  authStore.logout();
  setTimeout(() => {
    isLoggingOut.value = false;
  }, 1000);
}

const isAdmin = computed(() => {
  return authStore.isAdmin;
});

const username = computed(() => {
  return authStore.username;
});

const displayName = computed(() => {
  return authStore.displayName;
});

const pendingConfirmation = computed(() => {
  return appStore.pendingConfirmation;
});

function cancelConfirmation() {
  appStore.cancelConfirmation();
}

function confirmConfirmation() {
  appStore.confirmConfirmation();
}

const version = computed(() => {
  return import.meta.env.STABLE_VERSION;
});

onMounted(() => {
  authStore.init();
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
          <b-navbar-nav v-if="username !== null">
            <b-nav-item
              :to="{
                name: 'EditAccountRoute',
                params: { username: username },
              }"
              >My account</b-nav-item
            >
            <b-nav-item v-if="isAdmin" :to="{ name: 'AccountsRoute' }"
              >Accounts</b-nav-item
            >
            <b-nav-item v-if="isAdmin" :to="{ name: 'GroupsRoute' }"
              >Groups</b-nav-item
            >
          </b-navbar-nav>

          <b-navbar-nav>
            <b-nav-item right href="https://github.com/romnn/ldap-manager">{{
              version
            }}</b-nav-item>

            <b-nav-item-dropdown right v-if="username !== null">
              <template v-slot:button-content>
                <em>{{ displayName }} </em>
              </template>
              <!--
              <b-dropdown-item
                :to="{
                  name: 'EditAccountRoute',
                  params: { username: username },
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
        <b-container
          :toast="{ root: true }"
          fluid="sm"
          position="position-fixed"
          style="top: 50px; left: -200px"
        ></b-container>
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
