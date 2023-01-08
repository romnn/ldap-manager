<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import TableView from "../../components/TableView.vue";
import { RouterLink } from "vue-router";
import type { UserList } from "ldap-manager";
import { GatewayError } from "../../constants";
import { AxiosError } from "axios";

import { useToast } from "bootstrap-vue-3";
import { useAppStore } from "../../stores/app";
import { useAccountsStore } from "../../stores/accounts";

const toast = useToast();
const appStore = useAppStore();
const accountStore = useAccountsStore();

const list = ref<UserList>({ users: [], total: 0 });
const deleted = ref<string[]>([]);
const error = ref<string | undefined>(undefined);
const search = ref<string>("");
const loading = ref<boolean>(true);
const processing = ref<boolean>(false);

const currentPage = ref<number>(1);
const total = ref<number>(100);
const perPage = ref<number>(40);

const count = computed(() => list.value.users?.length ?? 0);

const pendingConfirmation = computed(() => {
  return appStore.pendingConfirmation;
});

/* function submitSearch() { */
/*   loadAccounts(); */
/* } */

function startSearch(s: string) {
  search.value = s;
}

function isDeleted(username: string) {
  return deleted.value.includes(username);
}

async function loadUsers() {
  error.value = undefined;
  list.value = { users: [], total: 0 };
  try {
    const users: UserList | undefined = await accountStore.listAccounts({
      page: currentPage.value,
      perPage: perPage.value,
      search: search.value,
    });
    if (!users) {
      error.value = "invalid user list";
      return;
    }
    list.value = users;
  } catch (err: unknown) {
    if (err instanceof GatewayError) {
      error.value = err.message;
    } else if (err instanceof AxiosError) {
      error.value = err.message;
    } else {
      throw err;
    }
  } finally {
    loading.value = false;
  }
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

async function deleteUser(username: string) {
  try {
    await appStore.newConfirmation({
      message: "Are you sure?",
      ack: "Yes, delete",
    });
  } catch (err: unknown) {
    return;
  }

  try {
    processing.value = true;
    await accountStore.deleteAccount(username);
    deleted.value.push(username);
  } catch (err: unknown) {
    if (err instanceof GatewayError) {
      errorAlert(err.message);
    } else if (err instanceof AxiosError) {
      errorAlert(err.message);
    } else {
      throw err;
    }
  } finally {
    processing.value = false;
  }
}

onMounted(() => {
  loadUsers();
});
</script>

<template>
  <div class="list-users-container">
    <TableView
      :inactive="pendingConfirmation !== null"
      :error="error"
      :loading="loading"
      :processing="processing"
      v-on:search="startSearch"
      searchLabel="Username:"
    >
      <!-- No results -->
      <div class="setup-user m-5" v-if="count < 1">
        <div v-if="search.length < 1">
          <p>There are no users yet</p>
          <p>
            <RouterLink :to="{ name: 'NewUserRoute' }"
              ><b-button size="sm" variant="primary"
                >Create a user</b-button
              ></RouterLink
            >
          </p>
        </div>
        <div v-else>
          <p>Did not find any users</p>
        </div>
      </div>
      <div v-else>
        <table class="users-table striped-table">
          <thead>
            <td>Username</td>
            <td>First Name</td>
            <td>Last Name</td>
            <td>E-Mail</td>
            <td></td>
          </thead>
          <tr
            v-for="(user, idx) in list.users"
            v-bind:key="user.UID"
            :class="{
              even: idx % 2 == 0,
              deleted: isDeleted(user.username),
            }"
          >
            <td>{{ user.username }}</td>
            <td>{{ user.firstName }}</td>
            <td>{{ user.lastName }}</td>
            <td>{{ user.email }}</td>
            <td>
              <span v-if="isDeleted(user.username)">Deleted</span>
              <div v-else>
                <b-button
                  pill
                  @click="deleteUser(user.username)"
                  size="sm"
                  class="mr-2 float-end"
                  variant="outline-danger"
                  >Delete</b-button
                >
                <RouterLink
                  :to="{
                    name: 'EditUserRoute',
                    params: { username: user.username },
                  }"
                  ><b-button
                    pill
                    size="sm"
                    class="mr-2 float-end"
                    variant="outline-info"
                    >Edit</b-button
                  ></RouterLink
                >
              </div>
            </td>
          </tr>
          <tr>
            <td></td>
            <td></td>
            <td></td>
            <td></td>
            <td>
              <div>
                <RouterLink
                  :to="{
                    name: 'NewUserRoute',
                  }"
                  ><b-button
                    pill
                    size="sm"
                    class="mr-2 float-end"
                    variant="outline-primary"
                    >Create</b-button
                  ></RouterLink
                >
              </div>
            </td>
          </tr>
        </table>

        <b-pagination
          class="users-pagination"
          size="sm"
          v-model="currentPage"
          :total-rows="total"
          :per-page="perPage"
          aria-controls="users-table"
        ></b-pagination>
      </div>
    </TableView>
  </div>
</template>

<style lang="sass" scoped>
.users-table
  table-layout: fixed
  width: 100%
  td
    word-wrap: break-word

/* is this used? */
.users-list
  z-index: 100
  &.inactive
    opacity: 0.2

  .setup-user
    padding: 30px

.users-pagination
  margin: 20px
  float: right
</style>
