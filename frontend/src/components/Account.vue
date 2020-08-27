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
                aria-describedby="login-input-uid-help-block"
              ></b-form-input>
              <b-form-text id="login-input-uid-help-block">
                Is optional. If you leave this empty, /bin/bash will be used
              </b-form-text>
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
                  required
                  placeholder=""
                  aria-describedby="login-input-password-help-block"
                ></b-form-input>
                <b-form-text id="login-input-password-help-block">
                  Good passwords must be 8-20 characters long, contain letters
                  and numbers, and must not contain spaces, special characters,
                  or emoji.
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
                  required
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
                  ><b-progress max="6">
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
import { AppModule } from "../store/modules/app";
import { GatewayError } from "../types";

// TODO: Form feedback helpers?

@Component
export default class AccountC extends Vue {
  @Prop() private account!: string;
  @Prop({ default: "Account" }) private title!: string;
  @Prop({ default: false }) private all!: boolean;
  @Prop({ default: false }) private create!: boolean;

  protected processing = false;
  protected emailRegex = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;

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

  get passwordsMatch() {
    return (
      (this.form.password + this.form.password_confirm).length > 0 &&
      this.form.password == this.form.password_confirm
    );
  }

  get validEmail() {
    return this.emailRegex.test(this.form.email);
  }

  get passwordStrengthLabel() {
    if (this.passwordStrength < 3) return "weak!";
    if (this.passwordStrength < 6) return "fair enough";
    return "good";
  }

  deleteAccount(username: string) {
    AppModule.newConfirmation({ message: "Are you sure?", ack: "Yes, delete" })
      .then(() => {
        this.processing = true;
        AccountModule.deleteAccount(username)
          .then(() => this.$router.push({ name: "AccountsRoute" }))
          .catch((err: GatewayError) => alert(err.message))
          .finally(() => (this.processing = false));
      })
      .catch(() => {
        // Ignore
      });
  }

  createAccount() {
    if (this.form.password !== this.form.password_confirm) return;
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
}
</script>

<style scoped lang="sass">
.account-label
    text-align: right
    font-weight: bold
</style>
