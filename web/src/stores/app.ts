import axios from "axios";
import {defineStore} from "pinia";
import {computed, ref} from "vue";
import {useRoute} from 'vue-router'

import {API_ENDPOINT} from "../constants";

export interface PendingConfirmation {
  message: string;
  ack: string;
  resolve?: () => void;
  reject?: () => void;
}

export const useAppStore = defineStore("app", () => {
  const route = useRoute();

  // const allowAll = ref(true);
  const pendingConfirmation: PendingConfirmation|null = ref(null);

  return {
    newConfirmation: (message: string, ack: string) => {
    return new Promise<void>((resolve, reject) => {
      pendingConfirmation.value = {
        message : req.message,
        ack : req.ack,
        resolve : resolve,
        reject : reject,
      };
    });
    },

    cancelConfirmation: () => {
    if (pendingConfirmation.value && pendingConfirmation.value.reject)
      pendingConfirmation.value.reject();
    pendingConfirmation.value = null;
    },

    confirmConfirmation: () => {
    if (pendingConfirmation.value && pendingConfirmation.value.resolve)
      pendingConfirmation.value.resolve();
    pendingConfirmation.value = null;
    },

    allowAll: computed(
      () => Object.keys(route?.query ?? {}).includes("all")
    ),
  };
});
