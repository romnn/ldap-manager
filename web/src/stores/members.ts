import axios from "axios";
import { defineStore } from "pinia";
import { API_ENDPOINT, handleError } from "../constants";

import type { GroupMember } from "ldap-manager";

export const useMembersStore = defineStore("members", () => {
  async function addGroupMember(member: GroupMember) {
    try {
      await axios.put(
        API_ENDPOINT + "/group/" + member.group + "/members",
        member
      );
    } catch (err: unknown) {
      handleError(err);
    }
  }

  async function removeGroupMember(member: GroupMember) {
    try {
      await axios.delete(
        API_ENDPOINT + "/group/" + member.group + "/member/" + member.username,
        {}
      );
    } catch (err: unknown) {
      handleError(err);
    }
  }

  return {
    addGroupMember,
    removeGroupMember,
  };
});
