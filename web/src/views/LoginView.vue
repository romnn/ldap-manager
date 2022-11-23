<script setup lang="ts">
import { ref, defineProps, computed, onMounted } from "vue";
import { GatewayError } from "../constants";
import type { Login } from "../stores/auth";

import { useAuthStore } from "../stores/auth";
import { useRouter } from "vue-router";

const auth = useAuthStore();
const router = useRouter();

const error = ref<string | undefined>(undefined);
const processing = ref<boolean>(false);

const form = ref<Login>({
  username: "",
  password: "",
  remember: true,
});

async function onSubmit() {
  try {
    error.value = undefined;
    processing.value = true;

    const request = {
      username: form.value.username,
      password: form.value.password,
      remember: form.value.remember,
    };
    await auth.login(request);
    router.push({
      name: "EditAccountRoute",
      params: { username: form.value.username },
    });
  } catch (err: unknown) {
    if (err instanceof GatewayError) {
      error.value = err.message;
    } else {
      throw err;
    }
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
