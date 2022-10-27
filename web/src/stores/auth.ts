import { defineStore } from "pinia";
import { computed, ref } from "vue";

export const useAuthStore = defineStore("auth", () => {
  const token = ref(null);
  const displayName = ref(null);
  const username = ref(null);
  const isAdmin = ref(null);

  return {
    token,
    displayName,
    username,
    isAdmin,
    authToken: computed(
      () => token ?? localStorage.getItem("x-user-token")
    ),
    activeIsAdmin: computed(
      () => isAdmin ?? localStorage.getItem("x-user-admin") !== null
    ),
    activeUsername: computed(
      () => username ?? localStorage.getItem("x-user-name")
    ),
    activeDisplayName: computed(
      () => displayName ?? localStorage.getItem("x-user-display-name")
    ),
    isAuthenticated: computed(
      () =>
        authToken != undefined &&
        authToken != null &&
        authToken.length > 0
    ),
    login: async (username: string, password: string, remember?: boolean) => {
      axios
        .post(API_ENDPOINT + "/login", {
          username,
          password,
        })
        .then(
          (response) => {
            handleTokenResponse({
              auth: response.data,
              remember: remember ?? false,
            });
            return response.data;
          },
          (error) => {
            return error;
          }
        );
    },
    handleTokenResponse: async (auth: TokenResponse, remember: boolean) => {
      if (remember ?? false) {
        localStorage.setItem("x-user-token", auth.token);
        localStorage.setItem("x-user-name", auth.username);
        if (auth.is_admin) localStorage.setItem("x-user-admin", "true");
        localStorage.setItem("x-user-display-name", auth.display_name);
      }
      token = auth.token;
      isAdmin = auth.is_admin;
      activeDisplayName = auth.display_name;
      activeUsername = auth.username;
    },
    logout: async () => {
      token = null;
      isAdmin = null;
      activeDisplayName = null;
      activeUsername = null;
      localStorage.removeItem("x-user-admin");
      localStorage.removeItem("x-user-token");
      localStorage.removeItem("x-user-name");
      localStorage.removeItem("x-user-display-name");
      router.push({ name: "LoginRoute" }).catch(() => {
        // Ignore
      });

      // Vue.nextTick(function() {
      //   router.push({ name: "LoginRoute" }).catch(() => {
      //     // Ignore
      //   });
      // });
    },
  };
});
