<template>
  <div class="group-container">
    <!-- This is just a proxy that migh be extended in the future -->
    <router-view />
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";

@Component
export default class GroupBase extends Vue {}
</script>

<style lang="sass" scoped></style>
