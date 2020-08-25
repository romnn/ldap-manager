import { VuexModule, Module, Action, getModule } from "vuex-module-decorators";
import store from "@/store";
import Vue from "vue";
import { API_ENDPOINT } from "../../constants";
import { Group } from "./groups";
import { GatewayError } from "../../types";

export interface GroupMemberState {}

@Module({ dynamic: true, store, name: "groups" })
class GroupMemberMod extends VuexModule implements GroupMemberState {
  @Action({ rawError: true })
  public async getGroup(name: string): Promise<Group> {
    return new Promise<Group>((resolve, reject) => {
      Vue.axios.get(API_ENDPOINT + "/group/" + name, {}).then(
        response => {
          resolve(response.data);
        },
        err => {
          reject(err.response?.data as GatewayError);
        }
      );
    });
  }

  @Action({ rawError: true })
  public async addGroupMember(group: string, member: string): Promise<void> {
    return new Promise<void>((resolve, reject) => {
      Vue.axios
        .put(API_ENDPOINT + "/group/" + group + "/members", { member })
        .then(
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
  public async deleteGroupMember(group: string, member: string): Promise<void> {
    return new Promise<void>((resolve, reject) => {
      Vue.axios
        .delete(API_ENDPOINT + "/group/" + group + "/member/" + member, {})
        .then(
          () => {
            resolve();
          },
          err => {
            reject(err.response?.data as GatewayError);
          }
        );
    });
  }
}

export const GroupMemberModule = getModule(GroupMemberMod);
