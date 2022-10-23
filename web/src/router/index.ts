import Vue from "vue";
import VueRouter, { RouteConfig } from "vue-router";
import { Route, RawLocation } from "vue-router";
import { AuthModule } from "@/store/modules/auth";

Vue.use(VueRouter);

const checkAuthenticated = (): boolean => {
  return AuthModule.isAuthenticated;
};

const checkNotAlreadyAuthenticated = (
  to: Route,
  from: Route,
  next: (to?: RawLocation | false | void) => void
) => {
  if (!checkAuthenticated()) {
    next();
    return;
  }
  next({ name: "HomeRoute" });
};

const requireAuth = (
  to: Route,
  from: Route,
  next: (to?: RawLocation | false | void) => void
) => {
  if (checkAuthenticated()) {
    next();
    return;
  }
  AuthModule.logout();
  next({ name: "LoginRoute" });
};

const routes: Array<RouteConfig> = [
  {
    path: "/",
    name: "HomeRoute",
    beforeEnter: requireAuth,
    meta: {
      base: []
    },
    component: () =>
      import(/* webpackChunkName: "homeAccount" */ "../views/Home.vue")
  },
  {
    path: "/login",
    name: "LoginRoute",
    meta: {
      base: []
    },
    beforeEnter: checkNotAlreadyAuthenticated,
    component: () =>
      import(/* webpackChunkName: "login" */ "../views/Login.vue")
  },
  {
    path: "/accounts",
    alias: "/account",
    name: "AccountsRoute",
    redirect: { name: "ListAccountsRoute" },
    beforeEnter: requireAuth,
    component: () =>
      import(/* webpackChunkName: "accounts" */ "../views/accounts/Base.vue"),
    children: [
      {
        path: "new",
        name: "NewAccountRoute",
        meta: {
          base: [
            {
              text: "Accounts",
              to: { name: "AccountsRoute" }
            },
            {
              text: "New",
              to: { name: "NewAccountRoute" },
              active: true
            }
          ]
        },
        beforeEnter: requireAuth,
        component: () =>
          import(
            /* webpackChunkName: "newAccount" */ "../views/accounts/New.vue"
          )
      },
      {
        path: "list",
        name: "ListAccountsRoute",
        meta: {
          base: [
            {
              text: "Accounts",
              to: { name: "AccountsRoute" }
            },
            {
              text: "List",
              to: { name: "ListAccountsRoute" },
              active: true
            }
          ]
        },
        beforeEnter: requireAuth,
        component: () =>
          import(
            /* webpackChunkName: "listAccounts" */ "../views/accounts/List.vue"
          )
      },
      {
        path: ":username",
        name: "EditAccountRoute",
        props: true,
        meta: {
          base: [
            {
              text: "Accounts",
              to: { name: "AccountsRoute" }
            }
          ]
        },
        beforeEnter: requireAuth,
        component: () =>
          import(
            /* webpackChunkName: "editAccount" */ "../views/accounts/Edit.vue"
          )
      }
    ]
  },
  {
    path: "/groups",
    alias: "/group",
    redirect: { name: "ListGroupsRoute" },
    name: "GroupsRoute",
    component: () =>
      import(/* webpackChunkName: "groups" */ "../views/groups/Base.vue"),
    beforeEnter: requireAuth,
    children: [
      {
        path: "new",
        name: "NewGroupRoute",
        meta: {
          base: [
            {
              text: "Groups",
              to: { name: "GroupsRoute" }
            },
            {
              text: "New",
              to: { name: "NewGroupRoute" },
              active: true
            }
          ]
        },
        beforeEnter: requireAuth,
        component: () =>
          import(/* webpackChunkName: "newGroup" */ "../views/groups/New.vue")
      },
      {
        path: "list",
        name: "ListGroupsRoute",
        meta: {
          base: [
            {
              text: "Groups",
              to: { name: "GroupsRoute" }
            },
            {
              text: "List",
              to: { name: "ListGroupsRoute" },
              active: true
            }
          ]
        },
        beforeEnter: requireAuth,
        component: () =>
          import(
            /* webpackChunkName: "listGroups" */ "../views/groups/List.vue"
          )
      },
      {
        path: ":name",
        name: "EditGroupRoute",
        props: true,
        meta: {
          base: [
            {
              text: "Groups",
              to: { name: "GroupsRoute" }
            }
          ]
        },
        beforeEnter: requireAuth,
        component: () =>
          import(/* webpackChunkName: "editGroup" */ "../views/groups/Edit.vue")
      }
    ]
  }
];

const router = new VueRouter({
  mode: "history",
  base: process.env.BASE_URL,
  routes
});

export default router;
