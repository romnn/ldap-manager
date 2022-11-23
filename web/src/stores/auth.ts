import axios from "axios";
import { LoginRequest, Token, User } from "ldap-manager";
import { defineStore } from "pinia";
import { computed, ref } from "vue";
import { useRouter } from "vue-router";

import { API_ENDPOINT, handleError } from "../constants";

const AUTH_TOKEN_KEY = "auth/token";
const AUTH_IS_ADMIN_KEY = "auth/admin";
const AUTH_USERNAME_KEY = "auth/username";
const AUTH_DISPLAY_NAME_KEY = "auth/displayname";

export interface Login {
  username: string;
  password: string;
  remember?: boolean;
}
export const useAuthStore = defineStore("auth", () => {
  const token = ref<string | null>(null);
  const isAdmin = ref<boolean | null>(null);
  const username = ref<string | null>(null);
  const displayName = ref<string | null>(null);

  const router = useRouter();

  const isAuthenticated = computed(
    () => token.value != null && token.value.length > 0
  );

  function init() {
    // read from local storage
    token.value = localStorage.getItem(AUTH_TOKEN_KEY);
    isAdmin.value = localStorage.getItem(AUTH_IS_ADMIN_KEY) !== null;
    username.value = localStorage.getItem(AUTH_USERNAME_KEY);
    displayName.value = localStorage.getItem(AUTH_DISPLAY_NAME_KEY);

    axios.defaults.headers.common["x-user-token"] = token.value;
  }

  async function logout() {
    token.value = null;
    isAdmin.value = null;
    displayName.value = null;
    username.value = null;
    localStorage.removeItem(AUTH_TOKEN_KEY);
    localStorage.removeItem(AUTH_IS_ADMIN_KEY);
    localStorage.removeItem(AUTH_USERNAME_KEY);
    localStorage.removeItem(AUTH_DISPLAY_NAME_KEY);
    router.push({ name: "LoginRoute" });
  }

  function updateUser(user: User) {
    displayName.value = user.displayName;
    username.value = user.username;
  }

  function updateToken({
    newToken,
    remember,
  }: {
    newToken: Token;
    remember?: boolean;
  }) {
    token.value = newToken.token;
    isAdmin.value = newToken.isAdmin;
    displayName.value = newToken.displayName;
    username.value = newToken.username;

    axios.defaults.headers.common["x-user-token"] = token.value;

    const defaultRemember = localStorage.getItem(AUTH_TOKEN_KEY) !== null;
    if (remember ?? defaultRemember) {
      localStorage.setItem(AUTH_TOKEN_KEY, newToken.token);
      localStorage.setItem(AUTH_USERNAME_KEY, newToken.username);
      if (isAdmin.value) {
        localStorage.setItem(AUTH_IS_ADMIN_KEY, "yes");
      } else {
        localStorage.removeItem(AUTH_IS_ADMIN_KEY);
      }
      localStorage.setItem(AUTH_DISPLAY_NAME_KEY, newToken.displayName);
    }
  }

  async function login({ username, password, remember }: Login) {
    try {
      const request: LoginRequest = {
        username,
        password,
      };
      const response = await axios.post(API_ENDPOINT + "/login", request);
      const newToken = Token.fromJSON(response.data);
      updateToken({ newToken, remember });
      return newToken;
    } catch (err: unknown) {
      handleError(err, false);
    }
  }

  return {
    token,
    displayName,
    username,
    isAdmin,
    isAuthenticated,
    logout,
    login,
    updateUser,
    updateToken,
    init,
  };
});
