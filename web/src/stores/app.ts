import { defineStore } from "pinia";
import { computed, ref } from "vue";
import { useRoute } from "vue-router";
import { parseBool } from "../utils";

export interface PendingConfirmation {
  message: string;
  ack: string;
  resolve?: () => void;
  reject?: () => void;
}

export const useAppStore = defineStore("app", () => {
  const route = useRoute();
  const pendingConfirmation = ref<PendingConfirmation | null>(null);

  function newConfirmation({ message, ack }: { message: string; ack: string }) {
    return new Promise<void>((resolve, reject) => {
      pendingConfirmation.value = {
        message,
        ack,
        resolve: resolve,
        reject: reject,
      };
    });
  }

  function cancelConfirmation() {
    if (pendingConfirmation.value && pendingConfirmation.value.reject)
      pendingConfirmation.value.reject();
    pendingConfirmation.value = null;
  }

  function confirmConfirmation() {
    if (pendingConfirmation.value && pendingConfirmation.value.resolve)
      pendingConfirmation.value.resolve();
    pendingConfirmation.value = null;
  }

  const allowAll = computed(() =>
    Object.keys(route?.query ?? {}).includes("all")
  );

  const showBranding = computed(() => parseBool(import.meta.env.VITE_BRANDING));

  const version = computed(() => {
    return import.meta.env.VITE_APP_VERSION;
  });

  return {
    pendingConfirmation,
    newConfirmation,
    cancelConfirmation,
    confirmConfirmation,
    allowAll,
    showBranding,
    version,
  };
});
