<template>
  <div class="new-account-container">
    <account-c create="true" title="New account"></account-c>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import AccountC from "../../components/Account.vue";
import { AppModule } from "../../store/modules/app";

@Component({ components: { AccountC } })
export default class NewAccount extends Vue {
  get all() {
    return AppModule.all;
  }
}
</script>

<style lang="sass" scoped></style>
