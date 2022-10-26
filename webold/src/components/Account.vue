<template>
  <div class="account-container">
    <div v-if="error !== null">
      <b-alert show variant="danger">
        <h4 class="alert-heading">Error</h4>
        <hr />
        <p class="mb-0">
          {{ error }}
        </p>
      </b-alert>
    </div>
    <div v-else>
      <b-overlay :show="processing" rounded="sm">
        <b-card
          class="login"
          header-tag="header"
          footer-tag="footer"
          :aria-hidden="processing ? 'true' : null"
        >
          <template v-slot:header>
            <b-row fluid class="text-center">
              <b-col sm="2"></b-col>
              <b-col>{{ title }}</b-col>
              <b-col sm="2"
                ><b-button
                  v-if="!create"
                  @click="deleteAccount"
                  pill
                  variant="outline-danger"
                  size="sm"
                  >Delete</b-button
                ></b-col
              >
            </b-row>
          </template>
          <b-card-body>
            <p class="text-left">
              If you are not sure about some values, just leave them blank. You
              will be prompted to enter missing values when you submit.
            </p>
            <b-form @submit.prevent="onSubmit">
              <b-form-group
                label-size="sm"
                label-cols-sm="3"
                label="First name:"
                class="account-label"
                label-for="login-input-username"
              >
                <b-form-input
                  autocomplete="off"
                  id="login-input-first-name"
                  size="sm"
                  v-model="form.first_name"
                  type="text"
                  required
                  placeholder="Max"
                ></b-form-input>
              </b-form-group>

              <b-form-group
                label-size="sm"
                label-cols-sm="3"
                label="Last name:"
                class="account-label"
                label-for="login-input-last-name"
              >
                <b-form-input
                  autocomplete="off"
                  id="login-input-last-name"
                  size="sm"
                  v-model="form.last_name"
                  type="text"
                  required
                  placeholder="Mustermann"
                ></b-form-input>
              </b-form-group>

              <b-form-group
                v-if="all && activeIsAdmin"
                label-size="sm"
                label-cols-sm="3"
                label="UID:"
                class="account-label"
                label-for="login-input-uid"
              >
                <b-form-input
                  autocomplete="off"
                  id="login-input-uid"
                  size="sm"
                  v-model="form.uid"
                  type="number"
                  placeholder="2004"
                  aria-describedby="login-input-uid-help-block"
                ></b-form-input>
                <b-form-text id="login-input-uid-help-block">
                  Is optional. If you leave this empty, /bin/bash will be used
                </b-form-text>
              </b-form-group>

              <b-form-group
                v-if="all && activeIsAdmin"
                label-size="sm"
                label-cols-sm="3"
                label="GID:"
                class="account-label"
                label-for="login-input-gid"
              >
                <b-form-input
                  autocomplete="off"
                  id="login-input-gid"
                  size="sm"
                  v-model="form.gid"
                  type="number"
                  placeholder="2001"
                  aria-describedby="login-input-gid-help-block"
                ></b-form-input>
                <b-form-text id="login-input-gid-help-block">
                  Is optional. If you leave this empty, /bin/bash will be used
                </b-form-text>
              </b-form-group>

              <b-form-group
                v-if="all"
                label-size="sm"
                label-cols-sm="3"
                label="Login shell:"
                class="account-label"
                label-for="login-input-shell"
              >
                <b-form-input
                  autocomplete="off"
                  id="login-input-shell"
                  size="sm"
                  v-model="form.login_shell"
                  type="text"
                  placeholder="/bin/bash"
                  aria-describedby="login-input-shell-help-block"
                ></b-form-input>
                <b-form-text id="login-input-shell-help-block">
                  Is optional. If you leave this empty, /bin/bash will be used
                </b-form-text>
              </b-form-group>

              <b-form-group
                v-if="all"
                label-size="sm"
                label-cols-sm="3"
                label="Home directory:"
                class="account-label"
                label-for="login-input-home-dir"
              >
                <b-form-input
                  autocomplete="off"
                  id="login-input-home-dir"
                  size="sm"
                  v-model="form.home_dir"
                  type="text"
                  placeholder="/home/max123"
                  aria-describedby="login-input-home-dir-help-block"
                ></b-form-input>
                <b-form-text id="login-input-home-dir-help-block">
                  Is optional. If you leave this empty, the /home/USERNAME is
                  chosen
                </b-form-text>
              </b-form-group>

              <b-form-group
                v-if="activeIsAdmin"
                label-size="sm"
                label-cols-sm="3"
                label="Username:"
                class="account-label"
                label-for="login-input-username"
              >
                <b-form-input
                  autocomplete="off"
                  id="login-input-username"
                  size="sm"
                  v-model="form.username"
                  type="text"
                  required
                  placeholder="max123"
                ></b-form-input>
              </b-form-group>

              <b-form-group
                label-size="sm"
                label-cols-sm="3"
                label="Email:"
                class="account-label"
                label-for="login-input-email"
              >
                <b-form-input
                  autocomplete="off"
                  id="login-input-email"
                  size="sm"
                  v-model="form.email"
                  :state="validEmail"
                  type="email"
                  required
                  placeholder="max.mustermann@example.com"
                ></b-form-input>
                <b-form-invalid-feedback :state="validEmail">
                  This is not a valid email
                </b-form-invalid-feedback>
                <b-form-valid-feedback :state="validEmail">
                  All good
                </b-form-valid-feedback>
              </b-form-group>

              <b-form-group>
                <b-form-group
                  label-size="sm"
                  label-cols-sm="3"
                  label="Password:"
                  class="account-label"
                  label-for="login-input-password"
                >
                  <b-form-input
                    autocomplete="off"
                    id="login-input-password"
                    size="sm"
                    v-model="form.password"
                    type="password"
                    :required="create"
                    placeholder=""
                    aria-describedby="login-input-password-help-block"
                  ></b-form-input>
                  <b-form-text id="login-input-password-help-block">
                    Good passwords must be 8-20 characters long, contain letters
                    and numbers, and must not contain spaces, special
                    characters, or emoji.
                  </b-form-text>
                </b-form-group>

                <b-form-group
                  label-size="sm"
                  label-cols-sm="3"
                  label="Confirm:"
                  class="account-label"
                  label-for="login-input-confirm-password"
                >
                  <b-form-input
                    autocomplete="off"
                    id="login-input-confirm-password"
                    size="sm"
                    v-model="form.password_confirm"
                    :state="passwordsMatch"
                    type="password"
                    :required="create"
                    placeholder="Confirm password"
                  ></b-form-input>
                  <b-form-invalid-feedback :state="passwordsMatch">
                    Passwords do not match
                  </b-form-invalid-feedback>
                  <b-form-valid-feedback :state="passwordsMatch">
                    All good
                  </b-form-valid-feedback>
                </b-form-group>

                <b-row align-h="end">
                  <b-col cols="9"
                    ><b-progress max="6" v-if="enteredPassword">
                      <b-progress-bar
                        :value="passwordStrength"
                        :label="passwordStrengthLabel"
                        :variant="passwordStrengthVariant"
                      ></b-progress-bar></b-progress
                  ></b-col>
                </b-row>
              </b-form-group>

              <b-form-group
                v-if="activeIsAdmin"
                label-size="sm"
                label-cols-sm="3"
                label="Groups:"
                class="account-label"
                label-for="account-input-groups"
              >
                <b-form-tags
                  autocomplete="off"
                  input-id="account-input-groups"
                  :invalid-tag-text="invalidGroupText"
                  duplicate-tag-text="already in group"
                  v-model="groups"
                  tag-variant="primary"
                  tag-pills
                  size="sm"
                  separator=" "
                  :state="groupsState"
                  :tag-validator="groupValidator"
                  placeholder="Enter group names separated by spaces"
                  class="mb-2"
                ></b-form-tags>
                <template v-slot:invalid-feedback>
                  No such group
                </template>
                <template v-slot:description>
                  <div v-if="checkingGroup">
                    Checking group...
                  </div>
                </template>
                <b-alert
                  class="text-left"
                  dismissible
                  :show="groupMemberError !== null"
                  variant="danger"
                >
                  {{ groupMemberError }}
                </b-alert>
              </b-form-group>

              <b-form-group class="mb-0">
                <b-button
                  class="float-right"
                  size="sm"
                  variant="primary"
                  @click="create ? createAccount() : updateAccount()"
                  >{{ create ? "Create account" : "Update" }}
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
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue, Watch } from "vue-property-decorator";
import {
  AccountModule,
  Account,
  RemoteAccount
} from "../store/modules/accounts";
import { GroupModule, GroupList } from "../store/modules/groups";
import { AppModule } from "../store/modules/app";
import { GatewayError } from "../types";
import { GroupMemberModule } from "../store/modules/members";
import { Codes } from "../constants";
import { AuthModule } from "../store/modules/auth";
import { AxiosError } from "axios";

@Component
export default class AccountC extends Vue {
  @Prop() private account!: string;
  @Prop({ default: "Account" }) private title!: string;
  @Prop({ default: false }) private all!: boolean;
  @Prop({ default: false }) private create!: boolean;

  protected invalidGroupText = "no such group";
  protected error: string | null = null;
  protected groupMemberError: string | null = null;
  protected submissionError: string | null = null;

  protected watchGroups = false;
  protected totalAvailableGroups?: number;
  protected availableGroups: string[] = [];
  protected groups: string[] = [];
  protected processing = false;
  protected checkingGroup = false;
  protected emailRegex = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;

  protected groupFetchRequest = {
    page: 1,
    perPage: 100,
    search: ""
  };

  protected form: Account = {
    /* eslint-disable-next-line @typescript-eslint/camelcase */
    first_name: "",
    /* eslint-disable-next-line @typescript-eslint/camelcase */
    last_name: "",
    uid: undefined,
    gid: undefined,
    /* eslint-disable-next-line @typescript-eslint/camelcase */
    login_shell: "",
    /* eslint-disable-next-line @typescript-eslint/camelcase */
    home_directory: "",
    username: "",
    email: "",
    password: "",
    /* eslint-disable-next-line @typescript-eslint/camelcase */
    password_confirm: ""
  };

  get activeIsAdmin() {
    return AuthModule.activeIsAdmin;
  }

  get groupsState() {
    return null;
  }

  @Watch("groups")
  onGroupsChanged(after: string[], before: string[]) {
    const a = new Set(after);
    const b = new Set(before);
    const added = new Set([...a].filter(x => !b.has(x)));
    const removed = new Set([...b].filter(x => !a.has(x)));
    if (this.watchGroups) {
      // .then(() => this.successAlert(`${this.account} was added to ${group}`)).catch((err: GatewayError) => this.errorAlert(err.message))
      added.forEach(group => this.addToGroup(group));
      removed.forEach(group => this.removeFromGroup(group));
    }
  }

  groupValidator(name: string): boolean | null {
    // we only have to search if total is > perPage
    // console.log(this.groupFetchRequest.perPage, this.totalAvailableGroups, this.availableGroups, name)
    if (this.groupFetchRequest.perPage <= (this.totalAvailableGroups ?? 0)) {
      this.invalidGroupText = "searching for group";
      if (!name.includes(this.groupFetchRequest.search)) {
        // set the search string for the group
        this.groupFetchRequest.search = name;
        this.fetchAvailableGroups()
          .then(() => this.groupValidator(name))
          .catch(() => {
            // Ignore
          });
        return null; // cant tell
      } else if (
        (this.groupFetchRequest.page - 1) * this.groupFetchRequest.perPage <
        (this.totalAvailableGroups ?? 0)
      ) {
        // increase the page
        this.groupFetchRequest.page = this.groupFetchRequest.page + 1;
        this.fetchAvailableGroups()
          .then(() => this.groupValidator(name))
          .catch(() => {
            // Ignore
          });
        return null; // cant tell
      }

      this.invalidGroupText = "no such group";
      return false; // no such group
    }

    this.invalidGroupText = "no such group";
    // we have all relevant groups and can safely filter
    const matches = this.availableGroups.filter((g: string) =>
      g.toLowerCase().includes(name.toLowerCase())
    );
    const found =
      matches.find((g: string) => g.toLowerCase() == name.toLowerCase()) !==
      undefined;
    if (!found && matches.length > 0) {
      this.invalidGroupText = `have group ${matches[0]}, but not`;
    }
    return found;
  }

  get passwordStrengthVariant() {
    if (this.passwordStrength < 3) return "danger";
    if (this.passwordStrength < 6) return "warning";
    return "success";
  }

  get passwordStrength() {
    return (
      1 +
      Number(/.{8,}/.test(this.form.password)) /* at least 8 characters */ *
        (Number(/.{12,}/.test(this.form.password)) /* bonus if longer */ +
        Number(/[a-z]/.test(this.form.password)) /* a lower letter */ +
        Number(/[A-Z]/.test(this.form.password)) /* a upper letter */ +
        Number(/\d/.test(this.form.password)) /* a digit */ +
          Number(
            /[^A-Za-z0-9]/.test(this.form.password)
          )) /* a special character */
    );
  }

  get enteredPassword() {
    return (this.form.password + this.form.password_confirm).length > 0;
  }

  get passwordsMatch(): boolean | null {
    return this.enteredPassword
      ? this.form.password == this.form.password_confirm
      : null;
  }

  get validEmail() {
    return this.emailRegex.test(this.form.email);
  }

  get passwordStrengthLabel() {
    if (this.passwordStrength < 3) return "weak!";
    if (this.passwordStrength < 6) return "fair enough";
    return "good";
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

  errorAlert(message: string, append = true) {
    this.$bvToast.toast(message, {
      title: "Error",
      autoHideDelay: 5000,
      appendToast: append,
      variant: "danger",
      solid: true
    });
  }

  deleteAccount() {
    AppModule.newConfirmation({ message: "Are you sure?", ack: "Yes, delete" })
      .then(() => {
        this.processing = true;
        AccountModule.deleteAccount(this.account)
          .then(() => {
            this.$router.push({ name: "LoginRoute" });
            this.successAlert(`${this.account} was successfully deleted`);
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
        // Ignore
      });
  }

  createAccount() {
    if (this.form.password !== this.form.password_confirm) return;
    this.submissionError = null;
    this.processing = true;
    AccountModule.newAccount(this.form)
      .then(() => {
        this.successAlert(`${this.form.username} was created`);
      })
      .catch((err: GatewayError) => {
        if (err.code == Codes.Unauthenticated) return AuthModule.logout();
        this.submissionError = err.message;
      })
      .finally(() => (this.processing = false));
  }

  removeFromGroup(group: string) {
    this.groupMemberError = null;
    this.watchGroups = false;
    GroupMemberModule.deleteGroupMember({
      username: this.account,
      group: group
    })
      .then(() => {
        this.successAlert(`${this.account} was removed from ${group}`);
      })
      .catch((err: GatewayError) => {
        if (err.code == Codes.Unauthenticated) return AuthModule.logout();
        this.groupMemberError = err.message;
        this.groups.push(group);
      })
      .finally(() => (this.watchGroups = true));
  }

  addToGroup(group: string) {
    this.groupMemberError = null;
    this.watchGroups = false;
    GroupMemberModule.addGroupMember({
      username: this.account,
      group: group
    })
      .then(() => {
        this.successAlert(`${this.account} was added to ${group}`);
      })
      .catch((err: GatewayError) => {
        if (err.code == Codes.Unauthenticated) return AuthModule.logout();
        this.groupMemberError = err.message;
        this.groups = this.groups.filter(g => g !== group);
      })
      .finally(() => (this.watchGroups = true));
  }

  updateAccount() {
    if (this.form.password !== this.form.password_confirm) return;
    this.submissionError = null;
    this.processing = true;
    AccountModule.updateAccount({ update: this.form, username: this.account })
      .then(() => {
        this.successAlert(`${this.account} was updated`);
        AccountModule.getAccount(this.account)
          .then((acc: RemoteAccount) => {
            AuthModule.setActiveDisplayName(
              (acc.data?.givenName ?? "") + " " + (acc.data?.sn ?? "")
            );
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

  onSubmit() {
    this.create ? this.createAccount() : this.updateAccount();
  }

  fetchAvailableGroups(): Promise<void> {
    return new Promise<void>((resolve, reject) => {
      this.checkingGroup = true;
      GroupModule.getGroups(this.groupFetchRequest)
        .then((list: GroupList) => {
          this.availableGroups = list.groups;
          this.totalAvailableGroups = Number(list.total);
          resolve();
        })
        .catch((err: AxiosError<GatewayError>) => {
          if (err.response?.data?.code == Codes.Unauthenticated)
            return AuthModule.logout();
          this.error = `${err.response?.data?.message ?? err}`;
          reject();
        })
        .finally(() => (this.checkingGroup = false));
    });
  }

  loadAccountData(account: string) {
    this.watchGroups = false;
    // Populate the form with the account data
    AccountModule.getAccount(account)
      .then((acc: RemoteAccount) => {
        /* eslint-disable-next-line @typescript-eslint/camelcase */
        this.form.first_name = acc.data?.givenName ?? "";
        /* eslint-disable-next-line @typescript-eslint/camelcase */
        this.form.last_name = acc.data?.sn ?? "";
        this.form.uid = Number(acc.data?.uidNumber);
        this.form.gid = Number(acc.data?.gidNumber);
        this.form.email = acc.data?.mail ?? "";
        /* eslint-disable-next-line @typescript-eslint/camelcase */
        this.form.login_shell = acc.data?.loginShell ?? "";
        /* eslint-disable-next-line @typescript-eslint/camelcase */
        this.form.home_directory = acc.data?.homeDirectory ?? "";
        this.form.username = acc.data?.uid ?? "";
      })
      .catch((err: AxiosError<GatewayError>) => {
        if (err.response?.data?.code == Codes.Unauthenticated)
          return AuthModule.logout();
        this.error = `${err.response?.data?.message ?? err}`;
      })
      .then(() => {
        // Get the accounts groups
        GroupModule.getUserGroups(this.account)
          .then((groups: GroupList) => {
            this.groups = groups.groups;
            this.$nextTick(function() {
              this.watchGroups = true;
            });
          })
          .catch((err: AxiosError<GatewayError>) => {
            if (err.response?.data?.code == Codes.Unauthenticated)
              return AuthModule.logout();
            this.error = `${err.response?.data?.message ?? err}`;
          });
      });
  }

  mounted() {
    this.error = null;

    if (!this.activeIsAdmin && this.create) {
      this.error =
        "Only admin users can create accounts. Please login as an admin user.";
      return;
    }

    // Fetch all available groups used for validating groups to join
    this.fetchAvailableGroups()
      .then(() => {
        if (!this.create) this.loadAccountData(this.account);
      })
      .catch(() => {
        // Ignore
      });
  }
}
</script>

<style scoped lang="sass">
.account-label
    text-align: right
    font-weight: bold
</style>
