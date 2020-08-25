<template>
  <div class="logout-container">
    <p>You are being logged out...</p>
    <p><b-spinner label="Logging out..."></b-spinner></p>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";

@Component
export default class Logout extends Vue {
  mounted() {
    // TODO: Start the logout process
  }
}
</script>

<style lang="sass" scoped></style>
