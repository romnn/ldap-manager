<template>
  <div class="account-container">
    <b-overlay :show="processing" rounded="sm">
      <b-card
        class="login"
        header-tag="header"
        footer-tag="footer"
        :aria-hidden="processing ? 'true' : null"
      >
        <template v-slot:header>
          <b-row class="text-center">
            <b-col></b-col>
            <b-col cols="8">{{ title }}</b-col>
            <b-col
              ><b-button
                v-if="!create"
                @click="deleteGroup"
                pill
                variant="outline-danger"
                size="sm"
                >Delete</b-button
              ></b-col
            >
          </b-row>
        </template>
        <b-card-body>
          <b-form>
            <b-form-group
              label-size="sm"
              label-cols-sm="3"
              label="Group name:"
              class="group-label"
              label-for="login-input-name"
            >
              <b-form-input
                autocomplete="off"
                id="login-input-name"
                size="sm"
                v-model="form.name"
                type="text"
                required
                placeholder="My group"
              ></b-form-input>
            </b-form-group>

            <b-form-group
              v-if="all"
              label-size="sm"
              label-cols-sm="3"
              label="GID:"
              class="group-label"
              label-for="group-input-gid"
            >
              <b-form-input
                autocomplete="off"
                id="group-input-gid"
                size="sm"
                v-model="form.gid"
                type="number"
                placeholder="2001"
                aria-describedby="group-input-gid-help-block"
              ></b-form-input>
              <b-form-text id="group-input-gid-help-block">
                Is optional. If you leave this empty, will be auto calculated
              </b-form-text>
            </b-form-group>

            <b-form-group>
              <b-row>
                <b-col>
                  <member-list-c
                    title="Members"
                    :loading="loadingMembers"
                    v-on:search="updateMemberSearch"
                  >
                    <div v-if="form.members.length < 1">
                      No members yet
                    </div>
                    <table v-else class="striped-table">
                      <thead>
                        <td>Username</td>
                        <td></td>
                      </thead>
                      <tr
                        v-for="(member, idx) in filteredMembers"
                        v-bind:key="member"
                        :class="{
                          even: idx % 2 == 0
                        }"
                      >
                        <td>{{ member }}</td>
                        <td>
                          <b-button
                            pill
                            @click="removeAccount(member)"
                            size="sm"
                            class="mr-2 float-end"
                            variant="outline-danger"
                            >Remove</b-button
                          >
                        </td>
                      </tr>
                    </table>
                  </member-list-c>
                </b-col>
                <b-col>
                  <member-list-c
                    title="All users"
                    :loading="loadingAvailableAccounts"
                    v-on:search="updateAvailableSearch"
                  >
                    <div v-if="available.users.length < 1">
                      No users available
                    </div>
                    <table v-else class="striped-table">
                      <thead>
                        <td>Username</td>
                        <td></td>
                      </thead>
                      <tr
                        v-for="(user, idx) in available.users"
                        v-bind:key="user.data.uid"
                        :class="{
                          even: idx % 2 == 0,
                          isMember: isMember(user.data.uid)
                        }"
                      >
                        <td>{{ user.data.uid }}</td>
                        <td>
                          <span v-if="isMember(user.data.uid)">
                            <i>member already</i>
                          </span>
                          <div v-else>
                            <b-button
                              pill
                              @click="addAccount(user.data.uid)"
                              size="sm"
                              class="mr-2 float-end"
                              variant="outline-primary"
                              >Add</b-button
                            >
                          </div>
                        </td>
                      </tr>
                    </table>
                  </member-list-c>
                  <b-alert
                    class="text-left"
                    dismissible
                    :show="loadingAvailableError !== null"
                    variant="danger"
                  >
                    {{ loadingAvailableError }}
                  </b-alert>
                </b-col>
              </b-row>
            </b-form-group>

            <b-alert
              class="text-left"
              dismissible
              :show="loadingGroupError !== null"
              variant="danger"
            >
              {{ loadingGroupError }}
            </b-alert>

            <b-alert
              class="text-left"
              dismissible
              :show="groupMemberOperationError !== null"
              variant="danger"
            >
              {{ groupMemberOperationError }}
            </b-alert>

            <b-form-group>
              <b-button
                class="float-end"
                size="sm"
                variant="primary"
                @click="create ? createGroup() : updateGroup()"
                >{{ create ? "Create group" : "Update" }}
              </b-button>
            </b-form-group>

            <b-alert
              class="text-left"
              dismissible
              :show="submissionError !== null"
              variant="danger"
            >
              {{ submissionError }}
            </b-alert>
          </b-form>
        </b-card-body>
      </b-card>
    </b-overlay>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from "vue-property-decorator";
import { GatewayError } from "../types";
import { GroupModule, Group } from "../store/modules/groups";
import { AppModule } from "../store/modules/app";
import { AccountModule, UserList } from "../store/modules/accounts";
import MemberListC from "./MemberList.vue";
import { Codes } from "../constants";
import { AuthModule } from "../store/modules/auth";
import { GroupMemberModule } from "../store/modules/members";

@Component({
  components: { MemberListC }
})
export default class GroupC extends Vue {
  @Prop() private name!: string;
  @Prop({ default: "Group" }) private title!: string;
  @Prop({ default: false }) private all!: boolean;
  @Prop({ default: false }) private create!: boolean;

  protected processing = false;
  protected loadingMembers = false;
  protected loadingAvailableAccounts = false;

  protected loadingAvailableError: string | null = null;
  protected loadingGroupError: string | null = null;
  protected groupMemberOperationError: string | null = null;
  protected submissionError: string | null = null;

  protected available: UserList = { users: [], total: "0" };

  protected membersSearch = "";
  protected availableSearch = "";

  protected form: {
    members: string[];
    name: string;
    gid: number;
  } = {
    members: [],
    name: "",
    gid: 0
  };

  get filteredMembers() {
    return this.form.members.filter(member => {
      return member.includes(this.membersSearch);
    });
  }

  updateMemberSearch(search: string) {
    this.membersSearch = search;
  }

  updateAvailableSearch(search: string) {
    this.availableSearch = search;
    this.loadAvailableAccounts();
  }

  isMember(username: string) {
    return this.form.members.includes(username);
  }

  successAlert(message: string, append = true) {
    this.$bvToast.toast(message, {
      title: "Success",
      autoHideDelay: 5000,
      appendToast: append,
      variant: "success",
      solid: true
    });
  }

  deleteGroup() {
    AppModule.newConfirmation({ message: "Are you sure?", ack: "Yes, delete" })
      .then(() => {
        this.processing = true;
        GroupModule.deleteGroup(this.name)
          .then(() => {
            this.$router.push({ name: "GroupsRoute" });
            this.successAlert(`${this.name} was deleted`);
          })
          .catch(() => {
            // Ignore
          })
          .catch((err: GatewayError) => {
            if (err.code == Codes.Unauthenticated) return AuthModule.logout();
            this.submissionError = err.message;
          })
          .finally(() => (this.processing = false));
      })
      .catch(() => {
        // Ingore
      });
  }

  createGroup() {
    this.processing = true;
    GroupModule.newGroup(this.form)
      .then(() => {
        this.successAlert(`${this.form.name} was created`);
        this.$router
          .push({
            name: "EditGroupRoute",
            params: { name: this.form.name }
          })
          .catch(() => {
            // Ignore
          });
      })
      .catch((err: GatewayError) => {
        if (err.code == Codes.Unauthenticated) return AuthModule.logout();
        this.submissionError = err.message;
      })
      .finally(() => (this.processing = false));
  }

  removeAccount(username: string) {
    if (this.create) {
      this.form.members = this.form.members.filter(
        member => member !== username
      );
      return;
    }
    this.processing = true;
    this.groupMemberOperationError = null;
    GroupMemberModule.deleteGroupMember({
      username: username,
      group: this.name
    })
      .then(() => {
        this.successAlert(`${username} was removed from ${this.name}`);
        this.form.members = this.form.members.filter(
          member => member !== username
        );
      })
      .catch((err: GatewayError) => {
        if (err.code == Codes.Unauthenticated) return AuthModule.logout();
        this.groupMemberOperationError = err.message;
      })
      .finally(() => (this.processing = false));
  }

  addAccount(username: string) {
    if (this.create) {
      this.form.members.push(username);
      return;
    }
    this.processing = true;
    this.groupMemberOperationError = null;
    GroupMemberModule.addGroupMember({
      username: username,
      group: this.name
    })
      .then(() => {
        this.successAlert(`${username} was added to ${this.name}`);
        this.form.members.push(username);
      })
      .catch((err: GatewayError) => {
        if (err.code == Codes.Unauthenticated) return AuthModule.logout();
        this.groupMemberOperationError = err.message;
      })
      .finally(() => (this.processing = false));
  }

  updateGroup() {
    this.processing = true;
    GroupModule.updateGroup({
      name: this.name,
      /* eslint-disable-next-line @typescript-eslint/camelcase */
      new_name: this.form.name,
      gid: this.form.gid
    })
      .then(() => {
        this.successAlert(`${this.name} was updated`);
        this.$router
          .push({
            name: "EditGroupRoute",
            params: { name: this.form.name }
          })
          .catch(() => {
            // Ignore
          });
      })
      .catch((err: GatewayError) => (this.submissionError = err.message))
      .finally(() => (this.processing = false));
  }

  loadAvailableAccounts() {
    this.loadingAvailableAccounts = true;
    this.loadingAvailableError = null;
    this.available = { users: [], total: "0" };
    AccountModule.listAccounts({
      search: this.availableSearch,
      page: 1,
      perPage: 50
    })
      .then((list: UserList) => {
        this.available.users = list?.users ?? [];
        this.available.total = list?.total ?? "0";
      })
      .catch((err: GatewayError) => {
        if (err.code == Codes.Unauthenticated) return AuthModule.logout();
        this.loadingAvailableError = err.message;
      })
      .finally(() => (this.loadingAvailableAccounts = false));
  }

  loadGroupData() {
    // Load members of the group
    this.loadingMembers = true;
    this.loadingGroupError = null;
    GroupMemberModule.getGroup(this.name)
      .then((group: Group) => {
        this.form.gid = group.gid;
        this.form.name = group.name;
        this.form.members = group.members;
      })
      .catch((err: GatewayError) => {
        if (err.code == Codes.Unauthenticated) return AuthModule.logout();
        this.loadingGroupError = err.message;
      })
      .finally(() => (this.loadingMembers = false));
  }

  mounted() {
    this.loadAvailableAccounts();
    if (!this.create) this.loadGroupData();
  }
}
</script>

<style scoped lang="sass">
group-label
    text-align: right
    font-weight: bold
</style>
