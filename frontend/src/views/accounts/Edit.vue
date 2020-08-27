<template>
  <div class="edit-account-container">
    <account-c :account="username" :title="'Edit ' + username"></account-c>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import AccountC from "../../components/Account.vue";
import { AppModule } from "../../store/modules/app";

@Component({ components: { AccountC } })
export default class EditAccount extends Vue {
  @Prop() private username!: string;
  get all() {
    return AppModule.all;
  }
}
</script>

<style lang="sass" scoped></style>
