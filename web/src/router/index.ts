import {createRouter, createWebHistory} from "vue-router";
import HomeView from "../views/HomeView.vue";

const checkAuthenticated = (): boolean => {
  return false;
  // return AuthModule.isAuthenticated;
};

const checkNotAlreadyAuthenticated =
    (to: Route, from: Route, next: (to?: RawLocation|false|void) => void) => {
      if (!checkAuthenticated()) {
        next();
        return;
      }
      next({name : "HomeRoute"});
    };

const requireAuth =
    (to: Route, from: Route, next: (to?: RawLocation|false|void) => void) => {
      if (checkAuthenticated()) {
        next();
        return;
      }
      // AuthModule.logout();
      // next({name : "LoginRoute"});
      next({name : "HomeRoute"});
    };

const router = createRouter({
  history : createWebHistory(import.meta.env.BASE_URL),
  routes : [
    {
      path : "/",
      name : "HomeRoute",
      // beforeEnter : requireAuth,
      meta : {base : []},
      component : () => import("../views/HomeView.vue"),
    },
    {
      path : "/login",
      name : "LoginRoute",
      meta : {base : []},
      // beforeEnter : checkNotAlreadyAuthenticated,
      component : () => import("../views/LoginView.vue")
    },
    {
      path : "/accounts",
      alias : "/account",
      name : "AccountsRoute",
      redirect : {name : "ListAccountsRoute"},
      // beforeEnter : requireAuth,
      meta : {base : []},
      component : () => import("../views/accounts/BaseView.vue"),
      children : [
        {
          path : "new",
          name : "NewUserRoute",
          meta : {
            base : [
              {text : "Users", to : {name : "UsersRoute"}},
              {text : "New", to : {name : "NewUserRoute"}, active : true}
            ]
          },
          // beforeEnter : requireAuth,
          component : () => import("../views/accounts/NewUserView.vue")
        },

        {
          path : "list",
          name : "ListAccountsRoute",
          meta : {
            base : [
              {
                text : "Accounts",
                to : {name : "AccountsRoute"},
              },
              {
                text : "List",
                to : {name : "ListAccountsRoute"},
                active : true,
              },
            ],
          },
          // beforeEnter: requireAuth,
          component : () => import("../views/accounts/ListView.vue"),
        },
      ],
    },
  ],
});

export default router;
