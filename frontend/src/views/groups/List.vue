<template>
  <div class="list-group-container">
    <table-view-c
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
        <table class="striped-table">
          <thead>
            <td>Name</td>
            <td></td>
          </thead>
          <tr
            v-for="(group, idx) in groups.groups"
            v-bind:key="group"
            :class="{
              even: idx % 2 == 0,
              deleted: isDeleted(group)
            }"
          >
            <td>{{ group }}</td>
            <td>
              <span v-if="isDeleted(group)">Deleted</span>
              <div v-else>
                <b-button
                  pill
                  @click="deleteGroup(group)"
                  size="sm"
                  class="mr-2 float-right"
                  variant="outline-danger"
                  >Delete</b-button
                >
                <router-link
                  :to="{
                    name: 'EditGroupRoute',
                    params: { name: group }
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
            <td>
              <div>
                <router-link
                  :to="{
                    name: 'NewGroupRoute'
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
          class="group-pagination"
          size="sm"
          v-model="currentPage"
          :total-rows="total"
          :per-page="perPage"
          aria-controls="group-table"
        ></b-pagination>
      </div>
    </table-view-c>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Watch } from "vue-property-decorator";
import { GroupModule, GroupList } from "../../store/modules/groups";
import { AppModule } from "../../store/modules/app";
import { GatewayError } from "../../types";
import TableViewC from "../../components/TableView.vue";
import { Codes } from "../../constants";
import { AuthModule } from "../../store/modules/auth";

@Component({
  components: { TableViewC }
})
export default class GroupListView extends Vue {
  groups: GroupList = { groups: [] };
  deleted: string[] = [];
  error: string | null = null;
  search = "";
  loading = true;
  processing = false;

  currentPage = 1;
  total = 100;
  perPage = 40;

  @Watch("currentPage")
  handleCurrentPageChange() {
    this.loadGroups();
  }

  get count() {
    return this.groups?.groups?.length ?? 0;
  }

  get pendingConfirmation() {
    return AppModule.pendingConfirmation;
  }

  startSearch(search: string) {
    this.search = search;
  }

  submitSearch() {
    this.loadGroups();
  }

  isDeleted(username: string) {
    return this.deleted.includes(username);
  }

  loadGroups() {
    this.error = null;
    this.groups = { groups: [] };
    GroupModule.getGroups({
      page: this.currentPage,
      perPage: this.perPage,
      search: this.search
    })
      .then((list: GroupList) => {
        this.groups = list;
      })
      .catch((err: GatewayError) => {
        if (err.code == Codes.Unauthenticated) return AuthModule.logout();
        this.error = err.message;
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

  deleteGroup(name: string) {
    AppModule.newConfirmation({ message: "Are you sure?", ack: "Yes, delete" })
      .then(() => {
        this.processing = true;
        GroupModule.deleteGroup(name)
          .then(() => this.deleted.push(name))
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
    this.loadGroups();
  }
}
</script>

<style lang="sass" scoped>

.confirmation
  border: 2px #e9ecef solid
  border-radius: 15px
  padding: 15px
  background-color: #ffffff
  z-index: 999999
  position: fixed
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);

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
