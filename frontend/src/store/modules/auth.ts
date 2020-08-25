import { VuexModule, Module, Action, getModule } from "vuex-module-decorators";
import store from "@/store";
import { GatewayError } from "../../types";

export interface AuthState {}

@Module({ dynamic: true, store, name: "auth" })
class AuthMod extends VuexModule implements AuthState {
  /*
  let token = response.data.token;
  localStorage.setItem("user-token", token);
  Vue.axios.defaults.headers.common["Authorization"] = token;
  commit("AUTH_SUCCESS", token);
  // dispatch("USER_REQUEST");
  resolve(response);
  */
  /*
  // commit("AUTH_ERROR", err);
          localStorage.removeItem("user-token");
          if (err.response) {
            let error_response = err.response.data;
            reject({
              status: err.response.status,
              error: error_response.error,
              message: error_response.message
            });
          }
          // Reject promise
          reject({
            status: null,
            error: "Authentication failed",
            message:
              "Could not properly connect to the server. Please try again later."
          });// commit("AUTH_ERROR", err);
          localStorage.removeItem("user-token");
          if (err.response) {
            let error_response = err.response.data;
            reject({
              status: err.response.status,
              error: error_response.error,
              message: error_response.message
            });
          }
          // Reject promise
          reject({
            status: null,
            error: "Authentication failed",
            message:
              "Could not properly connect to the server. Please try again later."
          });
  */
}

export const AuthModule = getModule(AuthMod);
