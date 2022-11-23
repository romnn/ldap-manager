<script setup lang="ts">
import AccountComponent from "../../components/AccountComponent.vue";

const props = defineProps({
  username: {
    type: String,
  },
});
</script>

<template>
  <div class="edit-account-container">
    <account-component
      :account="props.username"
      :title="'Edit ' + username"
    ></account-component>
  </div>
</template>

<style lang="sass" scoped></style>
