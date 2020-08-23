import Vue from "vue";
import VueRouter, { RouteConfig } from "vue-router";
import Login from "../views/Login.vue";

Vue.use(VueRouter);

const routes: Array<RouteConfig> = [
  {
    path: "/",
    name: "HomeRoute",
    component: Login
  },
  {
    path: "/login",
    name: "LoginRoute",
    component: () =>
      import(/* webpackChunkName: "login" */ "../views/Login.vue")
  },
  {
    path: "/account",
    name: "AccountsRoute",
    children: [
      {
        path: "new",
        name: "NewAccountRoute",
        component: () =>
          import(/* webpackChunkName: "new" */ "../views/accounts/New.vue")
      },
    ]
  }
];

const router = new VueRouter({
  mode: "history",
  base: process.env.BASE_URL,
  routes
});

export default router;
