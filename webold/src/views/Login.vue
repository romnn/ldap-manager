<template>
  <div class="login-container">
    <b-card class="login" header-tag="header" footer-tag="footer">
      <template v-slot:header>
        <h6 class="mb-0">Login</h6>
      </template>
      <b-card-body>
        <b-form @submit.prevent="onSubmit">
          <b-form-group
            label-size="sm"
            label-cols-sm="4"
            label-cols-lg="3"
            label="Username:"
            label-for="login-input-username"
          >
            <b-form-input
              autocomplete="off"
              id="login-input-username"
              size="sm"
              v-model="form.username"
              type="text"
              required
              placeholder="Enter username"
            ></b-form-input>
          </b-form-group>

          <b-form-group
            label-size="sm"
            label-cols-sm="4"
            label-cols-lg="3"
            label="Password:"
            label-for="login-input-password"
          >
            <b-form-input
              autocomplete="off"
              id="login-input-password"
              size="sm"
              v-model="form.password"
              type="password"
              required
              placeholder="Enter password"
            ></b-form-input>
          </b-form-group>

          <b-form-group inline label-cols-sm="3" class="mb-0">
            <b-form-row>
              <b-col>
                <b-form-checkbox size="sm" v-model="form.remember"
                  >Remember me</b-form-checkbox
                >
              </b-col>
              <b-col>
                <b-button size="sm" @click="onSubmit" type="submit" variant="primary"
                  >Log in</b-button
                >
              </b-col>
            </b-form-row>
          </b-form-group>

          <b-alert
            class="login-error"
            dismissible
            :show="error !== null"
            variant="danger"
          >
            <h4>Login failed</h4>
            <hr />
            {{ error }}
          </b-alert>
        </b-form>
      </b-card-body>
    </b-card>
  </div>
</template>

<script lang="ts">
import { AxiosError } from "axios";
import { Component, Vue } from "vue-property-decorator";
import { AuthModule } from "../store/modules/auth";
import { GatewayError } from "../types";

@Component({
  components: {}
})
export default class Login extends Vue {
  error: string | null = null;
  form = {
    username: "",
    password: "",
    remember: true
  };

  onSubmit() {
    AuthModule.login(this.form)
      .then(() => {
        this.$router
          .push({
            name: "EditAccountRoute",
            params: { username: this.form.username }
          })
          .catch(() => {
            // Ignore
          });
      })
      .catch((err: AxiosError<GatewayError>) => {
        this.error = `${err.response?.data?.message ?? err}`;
      });
  }
}
</script>

<style lang="sass" scoped>
.login-error
  margin-top: 40px
</style>
