<template>
  <div class="edit-group-container">
    <group-c :name="name" :title="'Edit ' + name" :all="all"></group-c>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import { AppModule } from "../../store/modules/app";
import GroupC from "../../components/Group.vue";

@Component({ components: { GroupC } })
export default class EditGroup extends Vue {
  @Prop() private name!: string;

  get all() {
    return AppModule.all;
  }
}
</script>

<style lang="sass" scoped></style>
