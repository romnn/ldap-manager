<template>
  <div class="home-container">
    <b-jumbotron>
      <template v-slot:header>LDAP Manager</template>

      <template v-slot:lead>
        The cloud-native LDAP management web interface
      </template>

      <hr class="my-4" />

      <p>
        What do you want to do?
      </p>

      <p>
        <b-button class="m-2" variant="info" :to="{ name: 'NewAccountRoute' }"
          >Create a new account</b-button
        >
        <b-button class="m-2" variant="light" :to="{ name: 'AccountsRoute' }"
          >Edit an existing account</b-button
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
    </b-jumbotron>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";

@Component
export default class Home extends Vue {}
</script>

<style lang="sass" scoped></style>
