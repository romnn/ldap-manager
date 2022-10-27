<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { RouterLink, RouterView } from "vue-router";
import TableView from "../../components/TableView.vue";
import { useAuthStore } from "../../stores/auth";
import { useAppStore } from "../../stores/app";
import { useAccountsStore } from "../../stores/accounts";
import { Codes } from "../../constants";

const list: UserList = ref({ users: [] });
const deleted: string[] = ref([]);
const error: string | null = ref(null);
const search = ref("");
const loading = ref(true);
const processing = ref(false);

const currentPage = ref(1);
const total = ref(100);
const perPage = ref(40);

const count = computed(() => list?.users?.length ?? 0);

/* const listOptions = computed(() => { */
/*   return { */
/*     page: currentPage, */
/*     perPage: perPage, */
/*     search: search, */
/*   }; */
/* }); */

const pendingConfirmation = computed(() => {
  const app = useAppStore();
  return app.pendingConfirmation;
});

function submitSearch(search: string) {
  loadAccounts();
}

function startSearch(search: string) {
  search.value = search;
}

function isDeleted(username: string) {
  return deleted.includes(username);
}

async function loadAccounts() {
  error.value = null;
  list.value = { users: [] };
  const auth = useAuthStore();
  const accounts = useAccountsStore();
  try {
    list.value = await accounts.listAccounts({
      page: currentPage.value,
      perPage: perPage.value,
      search: search.value,
    });
    console.log(list.value);
  } catch (error: AxiosError<GatewayError>) {
    console.log(error);
    if (error.response?.data?.code == Codes.Unauthenticated)
      return auth.logout();
    error.value = `${error.response?.data?.message ?? error}`;
  } finally {
    loading.value = false;
  }
}

function errorAlert(message: string, append = true) {
  /* this.$bvToast.toast(message, { */
  /*   title: "Error", */
  /*   autoHideDelay: 5000, */
  /*   appendToast: append, */
  /*   variant: "danger", */
  /*   solid: true */
  /* }); */
}

async function deleteAccount(username: string) {
  const app = useAppStore();
  const accounts = useAccountsStore();
  try {
    await app.newConfirmation({ message: "Are you sure?", ack: "Yes, delete" });
    processing.value = true;
    try {
      await accounts.deleteAccount(username);
      deleted.value.push(username);
    } catch (error: GatewayError) {
      if (error.code == Codes.Unauthenticated) return auth.logout();
      errorAlert(error.message);
    } finally {
      processing.value = false;
    }
  } catch (_) {
    // Ignore
  }
}

onMounted(() => {
  loadAccounts();
});
</script>

<template>
  <div class="list-account-container">
    <TableView
      :inactive="pendingConfirmation !== null"
      :error="error"
      :loading="loading"
      :processing="processing"
      v-on:search="startSearch"
      searchLabel="Username:"
    >
      <!-- No results -->
      <div class="setup-account m-5" v-if="count < 1">
        <div v-if="search.length < 1">
          <p>There are no accounts yet</p>
          <p>
            <!--
            <RouterLink :to="{ name: 'NewAccountRoute' }"
              ><b-button size="sm" variant="primary"
                >Create an account</b-button
              ></RouterLink
            >
            -->
          </p>
        </div>
        <div v-else>
          <p>Did not find any accounts</p>
        </div>
      </div>
      <div v-else>
        <table class="accounts-table striped-table">
          <thead>
            <td>Username</td>
            <td>First Name</td>
            <td>Last Name</td>
            <td>E-Mail</td>
            <td></td>
          </thead>
          <tr
            v-for="(user, idx) in list.users"
            v-bind:key="user.data.uid"
            :class="{
              even: idx % 2 == 0,
              deleted: isDeleted(user.data.uid),
            }"
          >
            <td>{{ user.data.uid }}</td>
            <td>{{ user.data.givenName }}</td>
            <td>{{ user.data.sn }}</td>
            <td>{{ user.data.mail }}</td>
            <td>
              <span v-if="isDeleted(user.data.uid)">Deleted</span>
              <div v-else>
                <b-button
                  pill
                  @click="deleteAccount(user.data.uid)"
                  size="sm"
                  class="mr-2 float-right"
                  variant="outline-danger"
                  >Delete</b-button
                >
                <RouterLink
                  :to="{
                    name: 'EditAccountRoute',
                    params: { username: user.data.uid },
                  }"
                  ><b-button
                    pill
                    size="sm"
                    class="mr-2 float-right"
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
                    name: 'NewAccountRoute',
                  }"
                  ><b-button
                    pill
                    size="sm"
                    class="mr-2 float-right"
                    variant="outline-primary"
                    >Create</b-button
                  ></RouterLink
                >
              </div>
            </td>
          </tr>
        </table>

        <b-pagination
          class="account-pagination"
          size="sm"
          v-model="currentPage"
          :total-rows="total"
          :per-page="perPage"
          aria-controls="accounts-table"
        ></b-pagination>
      </div>
    </TableView>
  </div>
</template>

<style lang="sass" scoped>
.accounts-table
  table-layout: fixed
  width: 100%
  td
    word-wrap: break-word

.account-list
  z-index: 100
  &.inactive
    opacity: 0.2

  .setup-account
    padding: 30px

.account-pagination
  margin: 20px
  float: right
</style>
