<template>
  <div class="list-account-container">
    <table-view-c
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
            <router-link :to="{ name: 'NewAccountRoute' }"
              ><b-button size="sm" variant="primary"
                >Create an account</b-button
              ></router-link
            >
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
              deleted: isDeleted(user.data.uid)
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
                <router-link
                  :to="{
                    name: 'EditAccountRoute',
                    params: { username: user.data.uid }
                  }"
                  ><b-button
                    pill
                    size="sm"
                    class="mr-2 float-right"
                    variant="outline-info"
                    >Edit</b-button
                  ></router-link
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
                <router-link
                  :to="{
                    name: 'NewAccountRoute'
                  }"
                  ><b-button
                    pill
                    size="sm"
                    class="mr-2 float-right"
                    variant="outline-primary"
                    >Create</b-button
                  ></router-link
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
    </table-view-c>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Watch } from "vue-property-decorator";
import { AccountModule, UserList } from "../../store/modules/accounts";
import { AppModule } from "../../store/modules/app";
import TableViewC from "../../components/TableView.vue";
import { GatewayError } from "../../types";
import { Codes } from "../../constants";
import { AuthModule } from "../../store/modules/auth";
import { AxiosError } from "axios";

@Component({
  components: { TableViewC }
})
export default class AccountList extends Vue {
  list: UserList = { users: [] };
  deleted: string[] = [];
  error: string | null = null;
  search = "";
  loading = true;
  processing = false;

  currentPage = 1;
  total = 100;
  perPage = 40;

  get count() {
    return this.list?.users?.length ?? 0;
  }

  get pendingConfirmation() {
    return AppModule.pendingConfirmation;
  }

  startSearch(search: string) {
    this.search = search;
  }

  submitSearch() {
    this.loadAccounts();
  }

  isDeleted(username: string) {
    return this.deleted.includes(username);
  }

  get listOptions() {
    return {
      page: this.currentPage,
      perPage: this.perPage,
      search: this.search
    };
  }

  @Watch("listOptions")
  handleCurrentPageChange() {
    this.loadAccounts();
  }

  loadAccounts() {
    this.error = null;
    this.list = { users: [] };
    AccountModule.listAccounts(this.listOptions)
      .then((list: UserList) => {
        this.list = list;
      })
      .catch((err: AxiosError<GatewayError>) => {
        if (err.response?.data?.code == Codes.Unauthenticated)
          return AuthModule.logout();
        this.error = `${err.response?.data?.message ?? err}`;
      })
      .finally(() => (this.loading = false));
  }

  errorAlert(message: string, append = true) {
    this.$bvToast.toast(message, {
      title: "Error",
      autoHideDelay: 5000,
      appendToast: append,
      variant: "danger",
      solid: true
    });
  }

  deleteAccount(username: string) {
    AppModule.newConfirmation({ message: "Are you sure?", ack: "Yes, delete" })
      .then(() => {
        this.processing = true;
        AccountModule.deleteAccount(username)
          .then(() => this.deleted.push(username))
          .catch((err: GatewayError) => {
            if (err.code == Codes.Unauthenticated) return AuthModule.logout();
            this.errorAlert(err.message);
          })
          .finally(() => (this.processing = false));
      })
      .catch(() => {
        // Ignore
      });
  }

  mounted() {
    this.loadAccounts();
  }
}
</script>

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
