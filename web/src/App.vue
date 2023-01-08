<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { RouterLink, RouterView } from "vue-router";
import { parseBool } from "./utils";
import type { RouteParams, RouteParamValue } from "vue-router";
import type { BreadcrumbItem } from "bootstrap-vue-3";
import ConfirmationComponent from "./components/ConfirmationComponent.vue";

import { useAuthStore } from "./stores/auth";
import { useAppStore } from "./stores/app";
import { useRoute } from "vue-router";

const route = useRoute();
const authStore = useAuthStore();
const appStore = useAppStore();

const isLoggingOut = ref<boolean>(false);

function logout() {
  isLoggingOut.value = true;
  authStore.logout();
  setTimeout(() => {
    isLoggingOut.value = false;
  }, 1000);
}

const items = computed((): BreadcrumbItem[] => {
  if (route.meta?.showBreadcrumb === false) return [] as BreadcrumbItem[];

  const baseItems = route.meta?.base ?? [];
  const params: RouteParams = route.params ?? {};

  const paramItems = Object.values(params).reduce(
    (acc: BreadcrumbItem[], param: RouteParamValue | RouteParamValue[]) => {
      const name = route.name;
      const paramList: RouteParamValue[] = Array.isArray(param)
        ? param
        : [param];
      if (name) {
        for (const p of paramList) {
          acc.push({
            text: p,
            to: { name: name, params: route.params },
            active: true,
          });
        }
      }
      return acc;
    },
    [] as BreadcrumbItem[]
  );

  return baseItems.concat(paramItems);
});

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

const branding = computed(() => {
  return parseBool(import.meta.env.VITE_BRANDING);
});

const version = computed(() => {
  return import.meta.env.VITE_APP_VERSION;
});

onMounted(() => {
  authStore.init();
});
</script>

<template>
  <div id="app">
    <div>
      <b-navbar toggleable="sm" size="sm" variant="dark">
        <!--
        <RouterLink :to="{ name: 'HomeRoute' }">
          <span class="title">
            {{ branding ? "LDAP Manager" : "Home" }}
          </span>
        </RouterLink>
        -->

        

        <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>

        <b-collapse class="navbar" id="nav-collapse" is-nav>
          <b-navbar-nav v-if="username !== null">
            <b-nav-item
              :to="{ name: 'HomeRoute' }"
              ><span class="title">{{ branding ? "LDAP Manager" : "Home" }}</span></b-nav-item
            >

            <b-nav-item
              :to="{
                name: 'EditUserRoute',
                params: { username: username },
              }"
              >My account</b-nav-item
            >
            <b-nav-item v-if="isAdmin" :to="{ name: 'UsersRoute' }"
              >Users</b-nav-item
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
              <b-dropdown-item
                :to="{
                  name: 'EditUserRoute',
                  params: { username: username },
                }"
                >My account</b-dropdown-item
              >
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
          style="z-index: 99999; top: 50px; left: -200px"
        ></b-container>
        <div v-if="pendingConfirmation !== null">
          <confirmation-component
            :message="pendingConfirmation.message"
            :ackMessage="pendingConfirmation.ack"
            v-on:cancel="cancelConfirmation"
            v-on:confirm="confirmConfirmation"
          ></confirmation-component>
        </div>
        <div class="logout-container" v-if="isLoggingOut">
          <p>You are being logged out...</p>
          <p><b-spinner label="Logging out..."></b-spinner></p>
        </div>
        <div v-else>
          <b-breadcrumb v-if="items.length > 0" :items="items"></b-breadcrumb>
          <RouterView />
        </div>
      </div>
    </div>
    <div class="footer">
      <a v-if="branding" href="https://github.com/romnn/ldap-manager"
        >LDAPManager {{ version }}</a
      >
    </div>
  </div>
</template>

<style lang="sass">
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

.navbar.bg-dark
  .nav-link
    color: white
    &.router-link-exact-active
      color: #0dcaf0

.navbar
  .title
    font-weight: 600
</style>
