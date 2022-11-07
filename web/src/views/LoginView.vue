<script setup lang="ts">
import { AxiosError } from "axios";
import { ref, defineProps, computed, onMounted } from "vue";
import { useAuthStore } from "../stores/auth";
/* import { GatewayError } from "../types"; */

const auth = useAuthStore();

const error: string | null = ref(null);
const processing: boolean = ref(false);

const form: {
  username: string;
  password: string;
  remember: boolean;
} = ref({
  username: "",
  password: "",
  remember: true,
});

async function onSubmit() {
  try {
    processing.value = true;
    const request = {
      username: form.value.username,
      password: form.value.password,
    };
    console.log(request);
    await auth.login(request);
    /* this.$router */
    /*   .push({ */
    /*     name: "EditAccountRoute", */
    /*     params: { username: this.form.username } */
    /*   }) */
  } catch (err: AxiosError<GatewayError>) {
    console.log(err);
    error.value = `${err.response?.data?.message ?? err}`;
  } finally {
    processing.value = false;
  }
}
</script>

<template>
  <div class="login-container">
    <b-overlay :show="processing" rounded="sm">
      <b-card
        class="login"
        header-tag="header"
        footer-tag="footer"
        :aria-hidden="processing ? 'true' : null"
      >
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

            <b-form-group inline class="mb-0">
              <b-form-row>
                <b-col>
                  <b-form-checkbox size="sm" v-model="form.remember"
                    >Remember me</b-form-checkbox
                  >
                </b-col>
                <b-col>
                  <b-button
                    size="sm"
                    class="float-right"
                    @click="onSubmit"
                    type="submit"
                    variant="primary"
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
    </b-overlay>
  </div>
</template>

<style lang="sass" scoped>
.login-error
  margin-top: 40px
</style>