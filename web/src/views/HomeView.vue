<script setup lang="ts">
import { useAppStore } from "../stores/app";
const appStore = useAppStore();
</script>

<template>
  <div class="home-container">
    <div class="h-100 p-5 bg-light border rounded-3">
      <h1 v-if="appStore.showBranding">LDAP Manager</h1>
      <p>The cloud-native LDAP management web interface</p>

      <hr class="my-4" />

      <p>What do you want to do?</p>

      <p>
        <b-button class="m-2" variant="info" :to="{ name: 'NewUserRoute' }"
          >Create a new user</b-button
        >
        <b-button class="m-2" variant="light" :to="{ name: 'UsersRoute' }"
          >Edit an existing user</b-button
        >
      </p>

      <p>
        <b-button class="m-2" variant="primary" :to="{ name: 'NewGroupRoute' }"
          >Create a new group</b-button
        >
        <b-button class="m-2" variant="light" :to="{ name: 'GroupsRoute' }"
          >Edit an existing group</b-button
        >
      </p>
    </div>
  </div>
</template>
