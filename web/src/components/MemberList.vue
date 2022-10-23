<template>
  <div class="group-member-list-container">
    <b-card no-body>
      <!-- Header -->
      <template v-slot:header>{{ title }}</template>

      <!-- Body -->
      <b-card-body class="p-0">
        <b-form-group class="controls">
          <b-form-row>
            <b-col cols="8">
              <b-form-input
                @keyup.enter="submitSearch"
                autocomplete="off"
                id="account-search-input"
                size="sm"
                v-model="search"
                type="text"
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
        <b-overlay :show="loading" rounded="sm">
          <div class="group-member-list-content">
            <slot></slot>
          </div>
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
  @Prop() private title!: string;
  search = "";

  submitSearch() {
    this.$emit("search", this.search);
  }
}
</script>

<style lang="sass" scoped>
.controls
  padding: 20px
  margin: 0
.group-member-list-content
  height: 300px
  overflow-y: scroll
</style>
