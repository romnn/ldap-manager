import axios from "axios";
import {defineStore} from "pinia";
import {computed, ref} from "vue";

import {API_ENDPOINT, handleError} from "../constants";
import {
  Token,
  SortOrder,
  ChangePasswordRequest,
  UpdateUserRequest,
  NewUserRequest,
  GetUserListRequest,
  User,
  UserList
} 
from "ldap-manager";

import {useAuthStore} from "./auth";

export const useAccountsStore = defineStore("accounts", () => {
  async function listAccounts({
    page,
    perPage,
    search,
  }: {page: number; perPage : number; search : string;}): Promise<UserList | undefined> {
    const request: GetUserListRequest = {
      start : (page - 1) * perPage,
      end : page * perPage,
      sortOrder : SortOrder.ASCENDING,
      sortKey : "",
      filter : [],
    };
    if (search.length > 0) {
      request.filter.push(`uid=${search}`);
    }

    try {
      const response = await axios.get(API_ENDPOINT + "/users", {
        params : request,
      });
      const users = UserList.fromJSON(response.data);
      return users;
    } catch (err: unknown) {
      handleError(err);
    }
  }

  async function getAccount(username: string): Promise<User | undefined> {
    try {
      const response = await axios.get(API_ENDPOINT + "/user/" + username, {});
      const user = User.fromJSON(response.data);
      return user;
    } catch (err: unknown) {
      handleError(err);
    }
  }

  async function newAccount(request: NewUserRequest): Promise<void> {
    try {
      await axios.post(API_ENDPOINT + "/user", request);
    } catch (err: unknown) {
      handleError(err);
    }
  }

  async function updateAccount(request: UpdateUserRequest): Promise<void> {
    try {
      const response = await axios.post(
          API_ENDPOINT + "/user/" + request.username + "/update", request);
      const newToken = Token.fromJSON(response.data);
      const authStore = useAuthStore();
      authStore.updateToken({newToken});
    } catch (err: unknown) {
      handleError(err);
    }
  }

  async function deleteAccount(username: string): Promise<void> {
    try {
      axios.delete(API_ENDPOINT + "/user/" + username, {});
    } catch (err: unknown) {
      handleError(err);
    }
  }

  async function changePassword(request: ChangePasswordRequest): Promise<void> {
    try {
      await axios.post(API_ENDPOINT + "/user/password", request);
    } catch (err: unknown) {
      handleError(err);
    }
  }

  return {
    listAccounts,
    getAccount,
    newAccount,
    updateAccount,
    deleteAccount,
    changePassword,
  };
});

// export interface UserList {
//   users?: {
//     data?: {
//       givenName: string;
//       mail: string;
//       sn: string;
//       uid: string;
//       [key: string]: string;
//     };
//   }[];
//   total?: string;
// }

// export interface RemoteAccount {
//   data?: {
//     cn?: string;
//     displayName?: string;
//     gidNumber?: string;
//     givenName?: string;
//     homeDirectory?: string;
//     loginShell?: string;
//     mail?: string;
//     sn?: string;
//     uid?: string;
//     uidNumber?: string;
//   };
// }

// export interface Account {
//   firstName: string;
//   lastName: string;
//   uid?: number;
//   gid?: number;
//   loginShell?: string;
//   homeDirectory?: string;
//   username: string;
//   email: string;
//   password: string;
//   passwordConfirm: string;

//   // first_name: string;
//   // last_name: string;
//   // uid?: number;
//   // gid?: number;
//   // login_shell?: string;
//   // home_directory?: string;
//   // username: string;
//   // email: string;
//   // password: string;
//   // password_confirm: string;
// }
