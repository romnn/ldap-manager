import axios from "axios";
import { defineStore } from "pinia";
import { computed, ref } from "vue";
import { API_ENDPOINT } from "../constants";

export interface PendingConfirmation {
  message: string;
  ack: string;
  resolve?: () => void;
  reject?: () => void;
}

export const useAppStore = defineStore("app", () => {
  const allowAll = ref(true);
  const pendingConfirmation: PendingConfirmation | null = ref(null);

  return {
    newConfirmation: (message: string, ack: string) => {
      return new Promise<void>((resolve, reject) => {
        pendingConfirmation.value = {
          message: req.message,
          ack: req.ack,
          resolve: resolve,
          reject: reject,
        };
      });
    },

    cancelConfirmation: () => {
      if (this.pendingConfirmation && this.pendingConfirmation.reject)
        this.pendingConfirmation.reject();
      this.pendingConfirmation = null;
    },

    confirmConfirmation: () => {
      if (this.pendingConfirmation && this.pendingConfirmation.resolve)
        this.pendingConfirmation.resolve();
      this.pendingConfirmation = null;
    },

    allowAll: computed(
      () =>
        this.allowAll &&
        Object.keys(router.currentRoute?.query ?? {}).includes("all")
    ),
  };
});
