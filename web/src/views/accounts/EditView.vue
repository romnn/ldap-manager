<script setup lang="ts">
import Account from "../../components/Account.vue";
/* import { Component, Vue, Prop } from "vue-property-decorator"; */
/* import { AppModule } from "../../store/modules/app"; */
/* import { useAppStore } from "../stores/app"; */
/* const app = useAppStore(); */

const props = defineProps({
  username: {
    type: String,
  },
});
</script>

<template>
  <div class="edit-account-container">
    <account :account="username" :title="'Edit ' + username"></account>
  </div>
</template>

<style lang="sass" scoped></style>
