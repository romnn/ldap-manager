import { VuexModule, Module, Action, getModule } from "vuex-module-decorators";
import store from "@/store";
import Vue from "vue";
import { API_ENDPOINT } from "../../constants";
import { TokenResponse, AuthModule } from "./auth";

export interface UserList {
  users?: {
    data?: {
      givenName: string;
      mail: string;
      sn: string;
      uid: string;
      [key: string]: string;
    };
  }[];
  total?: string;
}

export interface RemoteAccount {
  data?: {
    cn?: string;
    displayName?: string;
    gidNumber?: string;
    givenName?: string;
    homeDirectory?: string;
    loginShell?: string;
    mail?: string;
    sn?: string;
    uid?: string;
    uidNumber?: string;
  };
}

export interface Account {
  first_name: string;
  last_name: string;
  uid?: number;
  gid?: number;
  login_shell?: string;
  home_directory?: string;
  username: string;
  email: string;
  password: string;
  password_confirm: string;
  // [key: string]: string | number | undefined;
}

@Module({ dynamic: true, store, name: "accounts" })
class AccountMod extends VuexModule {
  @Action({ rawError: true })
  public async listAccounts(req: {
    page: number;
    perPage: number;
    search: string;
  }): Promise<UserList> {
    // we will not configure sort_key or sort_order
    const request: { start?: number; end?: number; filter?: string[] } = {
      start: (req.page - 1) * req.perPage,
      end: req.page * req.perPage
    };
    if (req.search.length > 0) {
      request.filter = ["uid=" + req.search];
    }

    return new Promise<UserList>((resolve, reject) => {
      Vue.axios.get(API_ENDPOINT + "/accounts", { params: request }).then(
        response => {
          resolve(response.data);
        },
        err => {
          reject(err.response);
        }
      );
    });
  }

  // TODO: AuthenticateUser

  @Action({ rawError: true })
  public async getAccount(username: string): Promise<RemoteAccount> {
    return new Promise<RemoteAccount>((resolve, reject) => {
      Vue.axios.get(API_ENDPOINT + "/account/" + username, {}).then(
        response => {
          resolve(response.data);
        },
        err => {
          reject(err.response);
        }
      );
    });
  }

  @Action({ rawError: true })
  public async newAccount(account: Account): Promise<void> {
    return new Promise<void>((resolve, reject) => {
      Vue.axios.put(API_ENDPOINT + "/account", { account }).then(
        () => {
          resolve();
        },
        err => {
          reject(err.response);
        }
      );
    });
  }

  @Action({ rawError: true })
  public async updateAccount(req: {
    update: Account;
    username: string;
  }): Promise<TokenResponse> {
    return new Promise<TokenResponse>((resolve, reject) => {
      Vue.axios
        .post(API_ENDPOINT + "/account/" + req.username + "/update", {
          update: req.update
        })
        .then(
          response => {
            // refreshed token for potentially updated account
            AuthModule.handleTokenResponse({
              auth: response.data,
              remember: localStorage.getItem("x-user-token") !== null
            });
            resolve();
          },
          err => {
            reject(err.response);
          }
        );
    });
  }

  @Action({ rawError: true })
  public async deleteAccount(username: string): Promise<void> {
    return new Promise<void>((resolve, reject) => {
      Vue.axios.delete(API_ENDPOINT + "/account/" + username, {}).then(
        () => {
          resolve();
        },
        err => {
          reject(err.response);
        }
      );
    });
  }

  @Action({ rawError: true })
  public async changePassword(
    username: string,
    newPassword: string
  ): Promise<void> {
    return new Promise<void>((resolve, reject) => {
      Vue.axios
        .post(API_ENDPOINT + "/account/password", { username, newPassword })
        .then(
          () => {
            resolve();
          },
          err => {
            reject(err.response);
          }
        );
    });
  }
}

export const AccountModule = getModule(AccountMod);
