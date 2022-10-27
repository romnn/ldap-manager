import axios from "axios";
import { defineStore } from "pinia";
import { computed, ref } from "vue";
import { API_ENDPOINT } from "../constants";

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

export const useAccountsStore = defineStore("accounts", () => {
  return {
    listAccounts: async ({
      page,
      perPage,
      search,
    }: {
      page: number;
      perPage: number;
      search: string;
    }) => {
      // we will not configure sort_key or sort_order
      const request: { start?: number; end?: number; filter?: string[] } = {
        start: (page - 1) * perPage,
        end: page * perPage,
      };
      if (search.length > 0) {
        request.filter = ["uid=" + search];
      }

      try {
        const response = await axios.get(API_ENDPOINT + "/accounts", {
          params: request,
        });
        return response.data;
      } catch (error) {
        return error.response;
      }
    },

    getAccount: async (username: string) => {
      try {
        const response = await axios.get(
          API_ENDPOINT + "/account/" + username,
          {}
        );
        return response.data;
      } catch (error) {
        return error.response;
      }
    },

    newAccount: async (account: Account) => {
      try {
        const response = await axios.get(API_ENDPOINT + "/account", {
          account,
        });
        return null;
      } catch (error) {
        return error.response;
      }
    },

    updateAccount: async (update: Account, username: string) => {
      try {
        const response = axios.post(
          API_ENDPOINT + "/account/" + req.username + "/update",
          {
            update: req.update,
          }
        );
        // refreshed token for potentially updated account
        // AuthModule.handleTokenResponse({
        //   auth: response.data,
        //   remember: localStorage.getItem("x-user-token") !== null
        // });
        return null;
      } catch (error) {
        return error.response;
      }
    },

    deleteAccount: async (username: string) => {
      try {
        const response = axios.delete(
          API_ENDPOINT + "/account/" + username,
          {}
        );
        return null;
      } catch (error) {
        return error.response;
      }
    },

    changePassword: async (username: string, newPassword: string) => {
      try {
        axios.post(API_ENDPOINT + "/account/password", {
          username,
          newPassword,
        });
      } catch (error) {
        return error.response;
      }
    },
  };
});
