<script setup lang="ts">
import { ref } from "vue";

const search = ref("");
const emit = defineEmits(["search"]);

const props = defineProps<{
  searchLabel?: string;
  error?: string;
  loading?: boolean;
  inactive?: boolean;
}>();

function submitSearch() {
  emit("search", search.value);
}
</script>

<template>
  <div class="table-view-list-container">
    <div class="table-view-list" :class="{ inactive: props.inactive }">
      <b-card no-body>
        <!-- Header -->
        <template v-slot:header>
          <b-form-group
            label-size="sm"
            label-cols-sm="2"
            :label="props.searchLabel"
            class="m-0"
            label-for="account-search-input"
          >
            <b-form-group inline class="mb-0">
              <b-form-row>
                <b-col>
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
                  <b-form-text
                    class="text-right"
                    id="account-search-input-help-block"
                  >
                    Confirm search with the enter key
                  </b-form-text>
                </b-col>
                <b-col cols="2">
                  <b-button @click="submitSearch" size="sm" variant="primary"
                    >Search</b-button
                  >
                </b-col>
              </b-form-row>
            </b-form-group>
          </b-form-group>
        </template>

        <!-- Body -->
        <b-card-body class="p-0">
          <div v-if="props.loading" class="m-5">
            <b-spinner label="Loading..."></b-spinner>
          </div>
          <div v-else>
            <!-- Error -->
            <b-alert
              class="text-left m-5"
              v-if="props.error !== undefined"
              :show="props.error !== undefined"
              variant="danger"
            >
              <h4 class="alert-heading">Error</h4>
              <hr />
              <p>
                {{ props.error }}
              </p>
            </b-alert>

            <div v-else>
              <!-- User content here -->
              <slot></slot>
            </div>
          </div>
        </b-card-body>
      </b-card>
    </div>
  </div>
</template>

<style lang="sass" scoped>
.table-view-list
  z-index: 100
  &.inactive
    opacity: 0.2
</style>
