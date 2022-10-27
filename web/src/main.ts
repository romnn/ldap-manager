import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import axios from "axios";
import bootstrap from 'bootstrap-vue-3'
// import bootstrap, { IconsPlugin } from 'bootstrap-vue-3'
// import { BootstrapVue, IconsPlugin } from "bootstrap-vue";

import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue-3/dist/bootstrap-vue-3.css'
import './assets/main.sass'

const app = createApp(App)

app.use(createPinia())
app.use(bootstrap)
// app.use(icons);
app.use(router)

app.mount('#app')
