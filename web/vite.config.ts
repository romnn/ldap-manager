import vue from '@vitejs/plugin-vue'
import {fileURLToPath, URL} from 'node:url'
import {defineConfig} from 'vite'

// https://vitejs.dev/config/
export default defineConfig({
  plugins : [ vue() ],
  server : {
    port : 8080,
    proxy : {
      '/api' : {target : 'http://localhost:8090'},
    }
  },
  resolve : {alias : {'@' : fileURLToPath(new URL('./src', import.meta.url))}}
})
