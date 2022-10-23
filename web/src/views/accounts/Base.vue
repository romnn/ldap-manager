<template>
  <div class="account-container">
    <!-- This is just a proxy that migh be extended in the future -->
    <router-view />
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";

@Component
export default class AccountBase extends Vue {}
</script>

<style lang="sass" scoped></style>
