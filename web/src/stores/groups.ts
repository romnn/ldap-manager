import axios from "axios";
import {defineStore} from "pinia";
import {computed, ref} from "vue";

import {API_ENDPOINT, handleError} from "../constants";
import {
  SortOrder,
  GetGroupListRequest,
  Group,
  GroupList,
  NewGroupRequest,
  UpdateGroupRequest
} from "ldap-manager";

export const useGroupsStore = defineStore("groups", () => {
  async function newGroup(request: NewGroupRequest) {
    try {
      await axios.put(API_ENDPOINT + "/group", request);
    } catch (err: unknown) {
      handleError(err);
    }
  }

  async function deleteGroup(name: string) {
    try {
      await axios.delete(API_ENDPOINT + "/group/" + name, {});
    } catch (err: unknown) {
      handleError(err);
    }
  }

  async function updateGroup(request: UpdateGroupRequest) {
    try {
      await axios.post(API_ENDPOINT + "/group/" + request.name + "/update",
                       request);
    } catch (err: unknown) {
      handleError(err);
    }
  }

  async function getGroup(name: string): Promise<Group | undefined> {
    try {
      const response = await axios.get(API_ENDPOINT + "/group/" + name, {});
      const group = Group.fromJSON(response.data);
      return group;
    } catch (err: unknown) {
      handleError(err);
    }
  }

  async function getGroups({
    page,
    perPage,
    search,
  }: {page: number; perPage : number; search : string;}) {
    try {
      const params: GetGroupListRequest = {
        start : (page - 1) * perPage,
        end : page * perPage,
        filter : [],
        sortOrder: SortOrder.ASCENDING,
        sortKey: "",
      };
      if (search.length > 0) {
        params.filter.push(`(cn=*${search}*)`);
      }

      const response = await axios.get(API_ENDPOINT + "/groups", {params});
      const groups = GroupList.fromJSON(response.data);
      return groups;
    } catch (err: unknown) {
      handleError(err);
    }
  }

  async function getUserGroups(username: string) {
    try {
      const response =
          await axios.get(API_ENDPOINT + "/user/" + username + "/groups", {});
      const groups = GroupList.fromJSON(response.data);
      return groups;
    } catch (err: unknown) {
      handleError(err);
    }
  }

  return {
    newGroup,
    deleteGroup,
    updateGroup,
    getGroup,
    getGroups,
    getUserGroups,
  };
});
