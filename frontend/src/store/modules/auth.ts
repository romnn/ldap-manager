import {
  VuexModule,
  Module,
  Mutation,
  Action,
  getModule
} from "vuex-module-decorators";
import store from "@/store";
import router from "@/router";
import Vue from "vue";
import { API_ENDPOINT } from "../../constants";
import { GatewayError } from "../../types";

export interface TokenResponse {
  token: string;
  display_name: string;
  expiration: string;
  username: string;
  is_admin: boolean;
}

export interface AuthState {
  token: string | null;
  displayName: string | null;
  username: string | null;
  isAdmin: boolean | null;
}

@Module({ dynamic: true, store, name: "auth" })
class AuthMod extends VuexModule implements AuthState {
  token: string | null = null;
  displayName: string | null = null;
  username: string | null = null;
  isAdmin: boolean | null = null;

  get authToken(): string | null {
    return this.token ?? localStorage.getItem("x-user-token");
  }

  get activeIsAdmin(): boolean {
    return this.isAdmin ?? localStorage.getItem("x-user-admin") !== null;
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
  public setIsAdmin(isAdmin: boolean | null) {
    this.isAdmin = isAdmin;
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
            this.handleTokenResponse({
              auth: response.data,
              remember: req.remember ?? false
            });
            resolve(response.data);
          },
          err => {
            reject(err.response?.data as GatewayError);
          }
        );
    });
  }

  @Action({ rawError: true })
  public handleTokenResponse(req: { auth: TokenResponse; remember: boolean }) {
    if (req.remember ?? false) {
      localStorage.setItem("x-user-token", req.auth.token);
      localStorage.setItem("x-user-name", req.auth.username);
      if (req.auth.is_admin) localStorage.setItem("x-user-admin", "true");
      localStorage.setItem("x-user-display-name", req.auth.display_name);
    }
    this.setToken(req.auth.token);
    this.setIsAdmin(req.auth.is_admin);
    this.setActiveDisplayName(req.auth.display_name);
    this.setActiveUsername(req.auth.username);
  }

  @Action({ rawError: true })
  public logout() {
    this.setToken(null);
    this.setIsAdmin(null);
    this.setActiveUsername(null);
    this.setActiveDisplayName(null);
    localStorage.removeItem("x-user-admin");
    localStorage.removeItem("x-user-token");
    localStorage.removeItem("x-user-name");
    localStorage.removeItem("x-user-display-name");
    Vue.nextTick(function() {
      router.push({ name: "LoginRoute" }).catch(() => {
        // Ignore
      });
    });
  }
}

export const AuthModule = getModule(AuthMod);
