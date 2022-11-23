<script setup lang="ts">
import AccountComponent from "../../components/AccountComponent.vue";
</script>

<template>
  <div class="new-group-container">
    <account-component :create="true" title="New user"></account-component>
  </div>
</template>

<style lang="sass" scoped></style>
