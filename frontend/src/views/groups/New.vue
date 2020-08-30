<template>
  <div class="new-group-container">
    <group-c :create="true" title="New group"></group-c>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import GroupC from "../../components/Group.vue";
import { AppModule } from "../../store/modules/app";

@Component({ components: { GroupC } })
export default class NewGroup extends Vue {
  get all() {
    return AppModule.all;
  }
}
</script>

<style lang="sass" scoped></style>
