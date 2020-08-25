<template>
  <div class="new-account-container">
    <account-c :account="username" :title="'Edit ' + username"></account-c>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import { AccountModule, UserList } from "../../store/modules/accounts";
import AccountC from "../../components/Account.vue";

@Component({ components: { AccountC } })
export default class EditAccount extends Vue {
  @Prop() private username!: string;
}
</script>

<style lang="sass" scoped></style>
