import { createRouter, createWebHistory } from "vue-router";
import type { RouteLocationNormalized, NavigationGuardNext } from "vue-router";
import HomeView from "../views/HomeView.vue";
import { useAuthStore } from "../stores/auth";

const checkAuthenticated = (): boolean => {
  return false;
  //
  // return AuthModule.isAuthenticated;
};

const checkNotAlreadyAuthenticated = (
  to: RouteLocationNormalized,
  from: RouteLocationNormalized,
  next: NavigationGuardNext
) => {
  const authStore = useAuthStore();
  if (!authStore.isAuthenticated) {
    next();
    return;
  }
  next({ name: "HomeRoute" });
};

const requireAuth = (
  to: RouteLocationNormalized,
  from: RouteLocationNormalized,
  next: NavigationGuardNext
) => {
  const authStore = useAuthStore();
  if (authStore.isAuthenticated) {
    next();
    return;
  }
  authStore.logout();
  next({ name: "LoginRoute" });
  // next({name : "HomeRoute"});
};

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "HomeRoute",
      beforeEnter: requireAuth,
      meta: { base: [] },
      component: () => import("../views/HomeView.vue"),
    },
    {
      path: "/login",
      name: "LoginRoute",
      meta: { base: [] },
      beforeEnter: checkNotAlreadyAuthenticated,
      component: () => import("../views/LoginView.vue"),
    },
    {
      path: "/accounts",
      alias: "/account",
      name: "AccountsRoute",
      redirect: { name: "ListAccountsRoute" },
      beforeEnter: requireAuth,
      meta: { base: [] },
      component: () => import("../views/accounts/BaseView.vue"),
      children: [
        {
          path: "new",
          name: "NewUserRoute",
          meta: {
            base: [
              { text: "Users", to: { name: "UsersRoute" } },
              { text: "New", to: { name: "NewUserRoute" }, active: true },
            ],
          },
          beforeEnter: requireAuth,
          component: () => import("../views/accounts/NewUserView.vue"),
        },

        {
          path: "list",
          name: "ListAccountsRoute",
          meta: {
            base: [
              {
                text: "Accounts",
                to: { name: "AccountsRoute" },
              },
              {
                text: "List",
                to: { name: "ListAccountsRoute" },
                active: true,
              },
            ],
          },
          beforeEnter: requireAuth,
          component: () => import("../views/accounts/ListView.vue"),
        },
        {
          path: ":username",
          name: "EditAccountRoute",
          props: true,
          meta: {
            base: [
              {
                text: "Accounts",
                to: { name: "AccountsRoute" },
              },
            ],
          },
          beforeEnter: requireAuth,
          component: () => import("../views/accounts/EditView.vue"),
        },
      ],
    },
    {
      path: "/groups",
      alias: "/group",
      redirect: { name: "ListGroupsRoute" },
      name: "GroupsRoute",
      component: () => import("../views/groups/Base.vue"),
      beforeEnter: requireAuth,
      children: [
        {
          path: "new",
          name: "NewGroupRoute",
          meta: {
            base: [
              {
                text: "Groups",
                to: { name: "GroupsRoute" },
              },
              {
                text: "New",
                to: { name: "NewGroupRoute" },
                active: true,
              },
            ],
          },
          beforeEnter: requireAuth,
          component: () => import("../views/groups/NewGroupView.vue"),
        },
        {
          path: "list",
          name: "ListGroupsRoute",
          meta: {
            base: [
              {
                text: "Groups",
                to: { name: "GroupsRoute" },
              },
              {
                text: "List",
                to: { name: "ListGroupsRoute" },
                active: true,
              },
            ],
          },
          beforeEnter: requireAuth,
          component: () => import("../views/groups/ListView.vue"),
        },
        {
          path: ":name",
          name: "EditGroupRoute",
          props: true,
          meta: {
            base: [
              {
                text: "Groups",
                to: { name: "GroupsRoute" },
              },
            ],
          },
          beforeEnter: requireAuth,
          component: () => import("../views/groups/EditView.vue"),
        },
      ],
    },
  ],
});

export default router;
