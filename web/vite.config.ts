import vue from "@vitejs/plugin-vue";
import { fileURLToPath, URL } from "node:url";
import * as path from "path";
import { defineConfig } from "vite";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  server: {
    port: 8080,
    proxy: {
      "/api": { target: "http://localhost:8090" },
    },
  },
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
      // "ldap-manager": "file:./generated",
      "ldap-manager": path.resolve(
        __dirname,
        "./generated/src/ldap_manager.ts"
      ),
    },
  },
});
