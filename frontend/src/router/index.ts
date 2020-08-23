import Vue from "vue";
import VueRouter, { RouteConfig } from "vue-router";

Vue.use(VueRouter);

const checkAuthenticated = (): boolean => {
  return true
  /*
  let token = Vue.cookies?.get("user-token") as string;
  let email = Vue.cookies?.get("user-email") as string;
  let userID = Vue.cookies?.get("user-id") as string;
  AuthModule.setAuthToken(token);
  AuthModule.setUserID(userID);
  UserManagementModule.setCurrentUserEmail(email);
  return (
    (token != undefined && token != null && token.length > 0) ||
    !authenticationRequired
  );
  */
};

const checkNotAlreadyAuthenticated = (to: any, from: any, next: any) => {
  if (!checkAuthenticated()) {
    next();
    return;
  }
  next({ name: "MyAccountRoute" });
};

const requireAdmin = (to: any, from: any, next: any) => {
  if (checkAuthenticated()) {
    next();
    return;
  }
  next({ name: "LoginRoute" });
};

const requireAuth = (to: any, from: any, next: any) => {
  if (checkAuthenticated()) {
    next();
    return;
  }
  next({ name: "LoginRoute" });
};

const routes: Array<RouteConfig> = [
  {
    path: "/",
    name: "MyAccountRoute",
    beforeEnter: requireAuth,
    component: () =>
      import(/* webpackChunkName: "editMyAccount" */ "../views/accounts/Edit.vue")
  },
  {
    path: "/login",
    name: "LoginRoute",
    beforeEnter: checkNotAlreadyAuthenticated,
    component: () =>
      import(/* webpackChunkName: "login" */ "../views/Login.vue")
  },
  {
    path: "/accounts",
    name: "AccountsRoute",
    redirect: { name: "ListAccountsRoute" },
    beforeEnter: requireAdmin,
    children: [
      {
        path: "new",
        name: "NewAccountRoute",
        beforeEnter: requireAdmin,
        component: () =>
          import(/* webpackChunkName: "newAccount" */ "../views/accounts/New.vue")
      },
      {
        path: "edit",
        name: "EditAccountRoute",
        beforeEnter: requireAdmin,
        component: () =>
          import(/* webpackChunkName: "editAccount" */ "../views/accounts/Edit.vue")
      },
      {
        path: "list",
        name: "ListAccountsRoute",
        beforeEnter: requireAdmin,
        component: () =>
          import(/* webpackChunkName: "listAccounts" */ "../views/accounts/List.vue")
      },
    ]
  },
  {
    path: "/groups",
    name: "GroupsRoute",
    beforeEnter: requireAdmin,
    children: [
      {
        path: "new",
        name: "NewGroupRoute",
        beforeEnter: requireAdmin,
        component: () =>
          import(/* webpackChunkName: "newGroup" */ "../views/groups/New.vue")
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
