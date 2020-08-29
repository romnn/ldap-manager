import {
  VuexModule,
  Module,
  Mutation,
  Action,
  getModule
} from "vuex-module-decorators";
import store from "@/store";
import Vue from "vue";
import { API_ENDPOINT } from "../../constants";
import { GatewayError } from "../../types";

export interface TokenResponse {
  token: string;
  display_name: string;
  expiration: string;
  username: string;
}

export interface AuthState {
  token: string | null;
  displayName: string | null;
  username: string | null;
}

@Module({ dynamic: true, store, name: "auth" })
class AuthMod extends VuexModule implements AuthState {
  token: string | null = null;
  displayName: string | null = null;
  username: string | null = null;

  get authToken(): string | null {
    return this.token ?? localStorage.getItem("x-user-token");
  }

  get activeUsername(): string | null {
    return this.username ?? localStorage.getItem("x-user-name");
  }

  get activeDisplayName(): string | null {
    return this.displayName ?? localStorage.getItem("x-user-display-name");
  }

  get isAuthenticated(): boolean {
    return (
      this.authToken != undefined &&
      this.authToken != null &&
      this.authToken.length > 0
    );
  }

  @Mutation
  public setActiveDisplayName(name: string | null) {
    this.displayName = name;
  }

  @Mutation
  public setActiveUsername(name: string | null) {
    this.username = name;
  }

  @Mutation
  public setToken(token: string | null) {
    this.token = token;
    Vue.axios.defaults.headers.common["x-user-token"] = token;
  }

  @Action({ rawError: true })
  public async login(req: {
    username: string;
    password: string;
    remember?: boolean;
  }): Promise<TokenResponse> {
    return new Promise<TokenResponse>((resolve, reject) => {
      Vue.axios
        .post(API_ENDPOINT + "/login", {
          username: req.username,
          password: req.password
        })
        .then(
          response => {
            const authResponse = response.data as TokenResponse;
            if (req.remember ?? false) {
              localStorage.setItem("x-user-token", authResponse.token);
              localStorage.setItem("x-user-name", authResponse.username);
              localStorage.setItem(
                "x-user-display-name",
                authResponse.display_name
              );
            }
            this.setToken(authResponse.token);
            this.setActiveDisplayName(authResponse.display_name);
            this.setActiveUsername(authResponse.username);
            resolve(authResponse);
          },
          err => {
            reject(err.response?.data as GatewayError);
          }
        );
    });
  }

  @Action({ rawError: true })
  public async logout(): Promise<void> {
    return new Promise<void>(resolve => {
      this.setToken(null);
      this.setActiveUsername(null);
      this.setActiveDisplayName(null);
      localStorage.removeItem("x-user-token");
      localStorage.removeItem("x-user-name");
      localStorage.removeItem("x-user-display-name");
      resolve();
    });
  }
}

export const AuthModule = getModule(AuthMod);
