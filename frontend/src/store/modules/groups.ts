import { VuexModule, Module, Action, getModule } from "vuex-module-decorators";
import store from "@/store";
import Vue from "vue";
import { API_ENDPOINT } from "../../constants";
import { GatewayError } from "../../types";

export interface Group {}

export interface GroupList {
  groups: Group[];
}

export interface GroupState {}

@Module({ dynamic: true, store, name: "groups" })
class GroupMod extends VuexModule implements GroupState {
  @Action({ rawError: true })
  public async newGroup(group: Group): Promise<void> {
    return new Promise<void>((resolve, reject) => {
      Vue.axios.put(API_ENDPOINT + "/group", group).then(
        () => {
          resolve();
        },
        err => {
          reject(err.response?.data as GatewayError);
        }
      );
    });
  }

  @Action({ rawError: true })
  public async deleteGroup(name: string): Promise<void> {
    return new Promise<void>((resolve, reject) => {
      Vue.axios.delete(API_ENDPOINT + "/group/" + name, {}).then(
        () => {
          resolve();
        },
        err => {
          reject(err.response?.data as GatewayError);
        }
      );
    });
  }

  @Action({ rawError: true })
  public async renameGroup(name: string): Promise<void> {
    return new Promise<void>((resolve, reject) => {
      Vue.axios.post(API_ENDPOINT + "/group/rename/" + name, {}).then(
        () => {
          resolve();
        },
        err => {
          reject(err.response?.data as GatewayError);
        }
      );
    });
  }

  @Action({ rawError: true })
  public async getGroups(): Promise<GroupList> {
    return new Promise<GroupList>((resolve, reject) => {
      Vue.axios.get(API_ENDPOINT + "/groups/", {}).then(
        response => {
          resolve(response.data);
        },
        err => {
          reject(err.response?.data as GatewayError);
        }
      );
    });
  }
}

export const GroupModule = getModule(GroupMod);
