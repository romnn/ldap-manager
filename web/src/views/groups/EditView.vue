<script setup lang="ts">
import { defineProps } from "vue";
import Group from "../../components/Group.vue";

const props = defineProps({
  name: {
    type: String,
  },
});
</script>

<template>
  <div class="edit-group-container">
    <group :name="name" :title="'Edit ' + props.name"></group>
  </div>
</template>

<style lang="sass" scoped></style>
