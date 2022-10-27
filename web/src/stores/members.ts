import axios from "axios";
import { defineStore } from "pinia";
import { computed, ref } from "vue";
import { API_ENDPOINT } from "../constants";

export const useMembersStore = defineStore("members", () => {
  return {
    getGroup: async (name: string) => {
      try {
        const response = await axios.get(API_ENDPOINT + "/group/" + name, {});
        return response.data;
      } catch (error) {
        return error.response;
      }
    },

    addGroupMember: async (group: string, username: string) => {
      try {
        const response = await axios.put(
          API_ENDPOINT + "/group/" + req.group + "/members",
          {
            username: req.username,
          }
        );
        return null;
      } catch (error) {
        return error.response;
      }
    },

    deleteGroupMember: async (group: string, username: string) => {
      try {
        const response = await axios.delete(
          API_ENDPOINT + "/group/" + req.group + "/member/" + req.username,
          {}
        );
        return null;
      } catch (error) {
        return error.response;
      }
    },
  };
});
