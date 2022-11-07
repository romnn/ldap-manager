import axios from "axios";
import {defineStore} from "pinia";
import {computed, ref} from "vue";

import {API_ENDPOINT} from "../constants";

const AUTH_TOKEN_KEY = "auth/token";
const AUTH_IS_ADMIN_KEY = "auth/admin";
const AUTH_USERNAME_KEY = "auth/username";
const AUTH_DISPLAY_NAME_KEY = "auth/displayname";

export const useAuthStore = defineStore("auth", () => {
  const token = ref(null);
  const isAdmin = ref(null);
  const username = ref(null);
  const displayName = ref(null);

  // todo: only if you want to remember
  // token = useLocalStorage('auth/token', null);
  // isAdmin = useLocalStorage('auth/admin', null);
  // username = useLocalStorage('auth/username', null);
  // displayName = useLocalStorage('auth/displayname', null);

  // const authToken = computed(
  //   () => token ?? localStorage.getItem("x-user-token")
  // );

  // const activeIsAdmin = computed(
  //   () => isAdmin ?? localStorage.getItem("x-user-admin") !== null
  // );

  // const activeUsername = computed(
  //   () => username ?? localStorage.getItem("x-user-name")
  // );

  // const activeDisplayName = computed(
  //   () => displayName ?? localStorage.getItem("x-user-display-name")
  // );

  const isAuthenticated =
      computed(() => token.value != undefined && token.value != null &&
                     token.value.length > 0);

  // async function handleTokenResponse(auth: TokenResponse) {
  //   console.log(auth);
  //   // if (remember ?? false) {
  //   //   localStorage.setItem("x-user-token", auth.token);
  //   //   localStorage.setItem("x-user-name", auth.username);
  //   //   if (auth.is_admin)
  //   //     localStorage.setItem(IS_ADMIN_KEY, "true");
  //   //   localStorage.setItem("x-user-display-name", auth.display_name);
  //   // }
  //   token.value = auth.token;
  //   isAdmin.value = auth.is_admin;
  //   displayName.value = auth.display_name;
  //   username.value = auth.username;
  // }

  async function logout() {
    token.value = null;
    isAdmin.value = null;
    displayName.value = null;
    username.value = null;

    localStorage.removeItem(IS_ADMIN_KEY);
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(USERNAME_KEY);
    localStorage.removeItem(DISPLAY_NAME_KEY);
    // router.push({name : "LoginRoute"}).catch(() => {
    //   // Ignore
    // });
    // Vue.nextTick(function() {
    //   router.push({ name: "LoginRoute" }).catch(() => {
    //     // Ignore
    //   });
    // });
  }

  async function login(
      request: {username: string; password : string; remember?: boolean;}) {
    const {username, password, remember} = request;

    try {
      console.log(request);
      const response = await axios.post(API_ENDPOINT + "/login", request)

      console.log(response);
      token.value = response.data.token;
      isAdmin.value = response.data.is_admin;
      displayName.value = response.data.display_name;
      username.value = response.data.username;

    if (remember ?? false) {
      localStorage.setItem(TOKEN_KEY, auth.token);
      localStorage.setItem(USERNAME_KEY, auth.username);
      if (auth.is_admin)
        localStorage.setItem(IS_ADMIN_KEY, "true");
      localStorage.setItem(DISPLAY_NAME_KEY, auth.display_name);
    }

    // await handleTokenResponse({
    //   auth: response.data,
    //   remember: remember ?? false,
    // });
    return response.data;
    } catch (err) {
      console.log(err);
      return err;
    }
  }

  return {
    token,
    displayName,
    username,
    isAdmin,
    // activeIsAdmin,
    // activeUsername,
    // activeDisplayName,
    isAuthenticated,
    logout,
    login,
  };
});
