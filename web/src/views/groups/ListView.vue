<script setup lang="ts">
import { ref, watch, computed, onMounted } from "vue";
import TableView from "../../components/TableView.vue";
import type { Group, GroupList } from "ldap-manager";
import { GatewayError } from "../../constants";

import { useToast } from "bootstrap-vue-3";
import { useGroupsStore } from "../../stores/groups";
import { useAppStore } from "../../stores/app";

const toast = useToast();
const appStore = useAppStore();
const groupsStore = useGroupsStore();

const groups = ref<Group[]>([]);
const deleted = ref<string[]>([]);
const error = ref<string | undefined>(undefined);
const search = ref<string>("");
const loading = ref<boolean>(true);
const processing = ref<boolean>(false);
const currentPage = ref<number>(1);
const total = ref<number>(100);
const perPage = ref<number>(40);

const count = computed(() => groups.value.length);

const pendingConfirmation = computed(() => appStore.pendingConfirmation);

/* async function submitSearch() { */
/*   await loadGroups(); */
/* } */

function startSearch(s: string) {
  search.value = s;
}

function isDeleted(username: string) {
  return deleted.value.includes(username);
}

function errorAlert(message: string) {
  toast?.danger(
    {
      title: "Error",
      body: message,
    },
    {
      autoHide: true,
      delay: 5000,
    }
  );
}

async function loadGroups() {
  try {
    loading.value = true;
    error.value = undefined;
    groups.value = [];
    const request = {
      page: currentPage.value,
      perPage: perPage.value,
      search: search.value,
    };
    const list: GroupList | undefined = await groupsStore.getGroups(request);
    if (!list) {
      error.value = "invalid group list";
      return;
    }
    groups.value = list.groups;
  } catch (err: unknown) {
    if (err instanceof GatewayError) {
      error.value = err.message;
    } else {
      throw err;
    }
  } finally {
    loading.value = false;
  }
}

async function deleteGroup(name: string) {
  await appStore.newConfirmation({
    message: "Are you sure?",
    ack: "Yes, delete",
  });
  try {
    processing.value = true;
    await groupsStore.deleteGroup(name);
    deleted.value.push(name);
  } catch (err: unknown) {
    if (err instanceof GatewayError) {
      errorAlert(err.message);
    } else {
      throw err;
    }
  } finally {
    processing.value = false;
  }
}

watch(currentPage, async () => {
  await loadGroups();
});

onMounted(async () => {
  await loadGroups();
});
</script>

<template>
  <div class="list-group-container">
    <table-view
      :inactive="pendingConfirmation !== null"
      :error="error"
      :loading="loading"
      :processing="processing"
      v-on:search="startSearch"
      searchLabel="Name:"
    >
      <!-- No results -->
      <div class="setup-group m-5" v-if="count < 1">
        <div v-if="search.length < 1">
          <p>There are no groups yet</p>
          <p>
            <router-link :to="{ name: 'NewGroupRoute' }"
              ><b-button size="sm" variant="primary"
                >Create a new group</b-button
              ></router-link
            >
          </p>
        </div>
        <div v-else>
          <p>Did not find any groups</p>
        </div>
      </div>
      <div v-else>
        <table class="groups-table striped-table">
          <thead>
            <td>Name</td>
            <td></td>
          </thead>
          <tr
            v-for="(group, idx) in groups"
            v-bind:key="group.GID"
            :class="{
              even: idx % 2 == 0,
              deleted: isDeleted(group.name),
            }"
          >
            <td>{{ group.name }}</td>
            <td>
              <span v-if="isDeleted(group.name)">Deleted</span>
              <div v-else>
                <b-button
                  pill
                  @click="deleteGroup(group.name)"
                  size="sm"
                  class="mr-2 float-end"
                  variant="outline-danger"
                  >Delete</b-button
                >
                <router-link
                  :to="{
                    name: 'EditGroupRoute',
                    params: { name: group.name },
                  }"
                  ><b-button
                    pill
                    size="sm"
                    class="mr-2 float-end"
                    variant="outline-info"
                    >Edit</b-button
                  ></router-link
                >
              </div>
            </td>
          </tr>
          <tr>
            <td></td>
            <td>
              <div>
                <router-link
                  :to="{
                    name: 'NewGroupRoute',
                  }"
                  ><b-button
                    pill
                    size="sm"
                    class="mr-2 float-end"
                    variant="outline-primary"
                    >Create</b-button
                  ></router-link
                >
              </div>
            </td>
          </tr>
        </table>

        <b-pagination
          class="group-pagination"
          size="sm"
          v-model="currentPage"
          :total-rows="total"
          :per-page="perPage"
          aria-controls="group-table"
        ></b-pagination>
      </div>
    </table-view>
  </div>
</template>

<style lang="sass" scoped>
.groups-table
  table-layout: fixed
  width: 100%
  td
    word-wrap: break-word

.confirmation
  border: 2px #e9ecef solid
  border-radius: 15px
  padding: 15px
  background-color: #ffffff
  z-index: 999999
  position: fixed
  top: 50%
  left: 50%
  transform: translate(-50%, -50%)

.group-list
  z-index: 100
  &.inactive
    opacity: 0.2

  .setup-account
    padding: 30px

.group-pagination
  margin: 20px
  float: right
</style>
