<template>
  <div class="group-member-list-container">
    <b-card no-body bg-variant="dark" text-variant="white">
      <!-- Header -->
      <template v-slot:header>
      </template>

      <!-- Body -->
      <b-card-body class="p-0">
        <b-overlay :show="loading" rounded="sm">
          <b-form-group>
            <b-form-row>
              <b-col cols="8">
                <b-form-input
                  @keyup.enter="submitSearch"
                  autocomplete="off"
                  id="account-search-input"
                  size="sm"
                  v-model="search"
                  type="text"
                  required
                  placeholder=""
                  aria-describedby="account-search-input-help-block"
                ></b-form-input>
              </b-col>
              <b-col cols="2">
                <b-button @click="submitSearch" size="sm" variant="primary"
                  >Search</b-button
                >
              </b-col>
            </b-form-row>
          </b-form-group>
          <slot></slot>
        </b-overlay>
      </b-card-body>
    </b-card>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";

@Component
export default class MemberListC extends Vue {
  @Prop() private loading!: boolean;
  search = "";

  submitSearch() {
    this.$emit("search", this.search);
  }
}
</script>

<style lang="sass" scoped>
.table-view-list
  z-index: 100
  &.inactive
    opacity: 0.2
</style>
