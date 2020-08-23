<template>
  <div class="login-container">
    <b-card class="login" header-tag="header" footer-tag="footer">
      <template v-slot:header>
        <h6 class="mb-0">Login</h6>
      </template>
      <b-card-body>
        <b-form @submit="onSubmit" @reset="onReset">
          <b-form-group
            label-size="sm"
            label-cols-sm="4"
            label-cols-lg="3"
            label="Username:"
            label-for="login-input-username"
          >
            <b-form-input
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
              id="login-input-password"
              size="sm"
              v-model="form.password"
              type="password"
              required
              placeholder="Enter password"
            ></b-form-input>
          </b-form-group>

          <b-form-group
            inline
            label-cols-sm="3"
            class="mb-0"
          >
            <b-form-row>
            <b-col>
              <b-form-checkbox size="sm" v-model="form.remember">Remember me</b-form-checkbox>
            </b-col>
            <b-col>
              <b-button size="sm" type="submit" variant="primary">Log in</b-button>
            </b-col>
          </b-form-row>
          </b-form-group>
        </b-form>
      </b-card-body>
    </b-card>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";

@Component({
  components: {
  }
})
export default class Login extends Vue {
  show = true
  form = {
    username: "",
    password: "",
    remember: true,
  }
  foods = ["a", "b"]

  onSubmit() {
    // TODO: Perform a auth request, save in cookie, show error if invalid
    // continues with run_checks
    // Password for <?php print $LDAP['admin_bind_dn'] <-- inject these into the frontend somehow
  }

  onReset() {
    
  }
}
</script>

<style lang="sass" scoped>
.login
  position: relative
  top: 100px
  min-width: 600px
  width: 50%
  margin: 0 auto
</style>
