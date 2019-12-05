import Vue from "vue";
import VueRouter from "vue-router";
import NProgress from 'nprogress';

Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    name: "signin",
    component: () => import("../views/SignIn")
  },
  {
    path: "/saber-mais",
    name: "learn-more",
    component: () => import("../views/LearnMore")
  },
  {
    path: '/criar-uma-conta',
    name: 'signup',
    component: () => import("../views/SignUp")
  }
];

const router = new VueRouter({
  routes
});

// Intercepting routes
router.beforeEach((to, from, next) => {
  // Start loading animation
  if (to.name) NProgress.start();

  next();
});
router.afterEach((to, from) => {
  // Stop loading animation
  NProgress.done();
});

export default router;
