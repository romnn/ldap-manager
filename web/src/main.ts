import "bootstrap/dist/css/bootstrap.css";
import "bootstrap-vue-3/dist/bootstrap-vue-3.css";
import "./assets/main.sass";

import { BToastPlugin } from "bootstrap-vue-3";
import bootstrap from "bootstrap-vue-3";
import { createPinia } from "pinia";
import { createApp } from "vue";

import App from "./App.vue";
import router from "./router";

const app = createApp(App);

app.use(createPinia());
app.use(bootstrap);
app.use(BToastPlugin);
// app.use(icons);
app.use(router);

app.mount("#app");
