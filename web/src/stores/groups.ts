import axios from "axios";
import { defineStore } from "pinia";
import { computed, ref } from "vue";
import { API_ENDPOINT } from "../../constants";

export interface Group {
  name: string;
  members: string[];
  total?: number;
  gid: number;
}

export interface GroupList {
  groups: string[];
  total?: string;
}

export const useGroupsStore = defineStore("groups", () => {
  return {
    newGroup: async (group: Group) => {
      try {
        const response = await axios.put(API_ENDPOINT + "/group", group);
        return null;
      } catch (error) {
        return error.response;
      }
    },

    deleteGroup: async (name: string) => {
      try {
        const response = await axios.delete(
          API_ENDPOINT + "/group/" + name,
          {}
        );
        return null;
      } catch (error) {
        return error.response;
      }
    },

    updateGroup: async (name: string, new_name?: string, gid?: number) => {
      try {
        const response = await axios.post(
          API_ENDPOINT + "/group/" + name + "/update",
          {
            /* eslint-disable-next-line @typescript-eslint/camelcase */
            new_name: new_name,
            gid: gid,
          }
        );
        return null;
      } catch (error) {
        return error.response;
      }
    },

    getGroups: async (page: number, perPage: number, search: string) => {
      try {
        const request: { start?: number; end?: number; filters?: string } = {
          start: (req.page - 1) * req.perPage,
          end: req.page * req.perPage,
        };
        if (req.search.length > 0) {
          request.filters = `(cn=*${req.search}*)`;
        }

        const response = await axios.get(API_ENDPOINT + "/groups", {
          params: request,
        });
        return response.data;
      } catch (error) {
        return error.response;
      }
    },

    getUserGroups: async (page: number, perPage: number, search: string) => {
      try {
        const response = await axios.get(
          API_ENDPOINT + "/account/" + username + "/groups",
          {}
        );
        return response.data;
      } catch (error) {
        return error.response;
      }
    },
  };
});
