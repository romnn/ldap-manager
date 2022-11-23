<script setup lang="ts">
import GroupComponent from "../../components/GroupComponent.vue";

const props = withDefaults(
  defineProps<{
    name?: string;
  }>(),
  {}
);
</script>

<template>
  <div class="edit-group-container">
    <group-component
      :name="props.name"
      :title="'Edit ' + props.name"
    ></group-component>
  </div>
</template>

<style lang="sass" scoped></style>
