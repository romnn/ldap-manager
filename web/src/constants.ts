import axios from "axios";
import { useAuthStore } from "./stores/auth";
import type { AxiosError } from "axios";
export const API_ENDPOINT = "/api/v1";

export const Codes = {
  Unauthenticated: 16,
};

export interface GatewayErrorI {
  code: number;
  message: string;
}

export function isGatewayError(error: any): error is GatewayErrorI {
  return "code" in error && "message" in error;
}

export class GatewayError {
  code: number;
  message: string;

  unauthenticated() {
    return this.code == Codes.Unauthenticated;
  }

  constructor({ code, message }: GatewayErrorI) {
    this.message = message;
    this.code = code;
  }
}

export function handleError(error: unknown, logout: boolean = true) {
  if (axios.isAxiosError(error)) {
    if (error.response?.data) {
      if (isGatewayError(error.response?.data)) {
        const gatewayError = new GatewayError(error.response?.data);
        if (logout && gatewayError.unauthenticated()) {
          const authStore = useAuthStore();
          authStore.logout();
          return;
        }
        throw gatewayError as GatewayError;
      }
    }
    throw error as AxiosError;
  } else {
    throw error as unknown;
  }
}
