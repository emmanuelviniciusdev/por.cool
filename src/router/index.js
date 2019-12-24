import Vue from "vue";
import VueRouter from "vue-router";
import NProgress from 'nprogress';
import firebase from 'firebase/app';
import 'firebase/auth';

Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    name: "signin",
    component: () => import("../views/SignIn"),
    meta: {
      isPublic: true
    }
  },
  {
    path: "/saber-mais",
    name: "learn-more",
    component: () => import("../views/LearnMore"),
    meta: {
      isPublic: true
    }
  },
  {
    path: '/criar-uma-conta',
    name: 'signup',
    component: () => import("../views/SignUp"),
    meta: {
      isPublic: true
    }
  },
  {
    path: '/pagar',
    name: 'payment',
    component: () => import("../views/Payment")
  },
  {
    path: '/definir-renda',
    name: 'define-monthly-income',
    component: () => import("../views/DefineMonthlyIncome")
  },
  {
    path: '/meus-gastos',
    name: 'home',
    component: () => import("../views/Home")
  }
];

const router = new VueRouter({
  routes
});

// Intercepting routes
router.beforeEach((to, from, next) => {
  // Start loading animation
  if (to.name) NProgress.start();

  // Check if user is authenticated only if router is not public
  const requiresAuth = to.matched.some(record => !record.meta.isPublic);
  if (requiresAuth) {
    const unsubscribe = firebase.auth().onAuthStateChanged(user => {
      // If we don't do this, the 'onAuthStateChanged()' will be executed a lot of times
      unsubscribe();

      if (!user) next({name: 'signin'});
    });
  }

  next();
});
router.afterEach((to, from) => {
  // Stop loading animation
  NProgress.done();
});

export default router;
