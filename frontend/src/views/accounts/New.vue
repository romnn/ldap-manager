<template>
  <div class="new-account-container">
    <account-c create="true" title="New account"></account-c>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { AccountModule, UserList } from "../../store/modules/accounts";
import AccountC from "../../components/Account.vue";

@Component({ components: { AccountC } })
export default class NewAccount extends Vue {}
</script>

<style lang="sass" scoped></style>
