<script setup lang="ts">
import { ref } from "vue";

const search = ref<string>("");

const props = withDefaults(
  defineProps<{
    loading?: boolean;
    title?: string;
  }>(),
  {}
);

const emit = defineEmits(["search"]);

function submitSearch() {
  emit("search", search.value);
}
</script>

<template>
  <div class="group-member-list-container">
    <b-card no-body>
      <!-- Header -->
      <template v-slot:header>{{ props.title }}</template>

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
        <b-overlay :show="props.loading" rounded="sm">
          <div class="group-member-list-content">
            <slot></slot>
          </div>
        </b-overlay>
      </b-card-body>
    </b-card>
  </div>
</template>

<style lang="sass" scoped>
.controls
  padding: 20px
  margin: 0
.group-member-list-content
  height: 300px
  overflow-y: scroll
</style>
