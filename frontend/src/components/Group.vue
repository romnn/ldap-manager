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
          If you are not sure, just leave blank
          <b-form @submit.prevent="onSubmit" @reset.prevent="onReset">
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
              v-if="all"
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
              ></b-form-input>
            </b-form-group>

            <b-form-group
              v-if="all"
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
              ></b-form-input>
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
              ></b-form-input>
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
              ></b-form-input>
            </b-form-group>

            <b-form-group
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
                type="email"
                required
                placeholder="max.mustermann@example.com"
              ></b-form-input>
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
                  required
                  placeholder=""
                ></b-form-input>
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
                  type="password"
                  required
                  placeholder="Confirm password"
                ></b-form-input>
              </b-form-group>

              <b-row align-h="end">
                <b-col cols="9"
                  ><b-progress max="100">
                    <b-progress-bar
                      :value="passwordStrength"
                      :label="passwordStrengthLabel"
                      :variant="passwordStrengthVariant"
                    ></b-progress-bar> </b-progress
                ></b-col>
              </b-row>
            </b-form-group>

            <b-form-group>
              <b-button
                class="float-right"
                size="sm"
                type="submit"
                variant="primary"
                >{{ create ? "Create account" : "Update" }}
              </b-button>
            </b-form-group>
          </b-form>
        </b-card-body>
      </b-card>
    </b-overlay>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from "vue-property-decorator";
import { AccountModule, Account } from "../store/modules/accounts";
import { GatewayError } from "../types";

// TODO: Form feedback helpers?

@Component
export default class GroupC extends Vue {
  @Prop() private account!: string;
  @Prop({ default: "Account" }) private title!: string;
  @Prop({ default: false }) private all!: boolean;
  @Prop({ default: false }) private create!: boolean;

  protected processing = false;

  protected form: Account = {
    first_name: "",
    last_name: "",
    uid: undefined,
    gid: undefined,
    login_shell: "",
    home_directory: "",
    username: "",
    email: "",
    password: "",
    password_confirm: ""
  };

  get passwordStrengthVariant() {
    return "success"; // "warning", "danger"
  }

  get passwordStrength() {
    return this.form.password.length;
  }

  get passwordStrengthLabel() {
    return "weak!";
  }

  deleteAccount() {
    this.processing = true;
    AccountModule.deleteAccount(this.account)
      .catch((err: GatewayError) => alert(err.message))
      .finally(() => (this.processing = false));
  }

  createAccount() {
    this.processing = true;
    AccountModule.newAccount(this.form)
      .catch((err: GatewayError) => alert(err.message))
      .finally(() => (this.processing = false));
  }

  updateAccount() {
    this.processing = true;
  }

  onSubmit() {
    this.create ? this.createAccount() : this.updateAccount();
    // TODO: Perform a auth request, save in cookie, show error if invalid
    // continues with run_checks
    // Password for <?php print $LDAP['admin_bind_dn'] <-- inject these into the frontend somehow
  }

  onReset() {}
}
</script>

<style scoped lang="sass">
.account-label
    text-align: right
    font-weight: bold
</style>
