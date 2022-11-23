import { createRouter, createWebHistory } from "vue-router";
// see https://github.com/vuejs/router/blob/main/packages/router/src/types/index.ts
import type { RouteLocationNormalized, NavigationGuardNext } from "vue-router";
import type { BreadcrumbItem } from "bootstrap-vue-3";
import { useAuthStore } from "../stores/auth";

declare module "vue-router" {
  interface RouteMeta {
    base?: BreadcrumbItem[];
  }
}

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
      path: "/users",
      alias: "/user",
      name: "UsersRoute",
      redirect: { name: "ListUsersRoute" },
      beforeEnter: requireAuth,
      meta: { base: [] },
      component: () => import("../views/users/BaseView.vue"),
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
          component: () => import("../views/users/NewUserView.vue"),
        },

        {
          path: "list",
          name: "ListUsersRoute",
          meta: {
            base: [
              {
                text: "Users",
                to: { name: "UsersRoute" },
              },
              {
                text: "List",
                to: { name: "ListUsersRoute" },
                active: true,
              },
            ],
          },
          beforeEnter: requireAuth,
          component: () => import("../views/users/ListView.vue"),
        },
        {
          path: ":username",
          name: "EditUserRoute",
          props: true,
          meta: {
            base: [
              {
                text: "Users",
                to: { name: "UsersRoute" },
              },
            ],
          },
          beforeEnter: requireAuth,
          component: () => import("../views/users/EditView.vue"),
        },
      ],
    },
    {
      path: "/groups",
      alias: "/group",
      redirect: { name: "ListGroupsRoute" },
      name: "GroupsRoute",
      component: () => import("../views/groups/BaseView.vue"),
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
