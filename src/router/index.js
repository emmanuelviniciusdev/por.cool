import Vue from "vue";
import VueRouter from "vue-router";
import NProgress from 'nprogress';
import firebase from 'firebase/app';
import 'firebase/auth';
import 'firebase/firestore';

// Services
import userService from '../services/user';

// Helpers
import paymentHelper from '../helpers/payment';

Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    name: "signin",
    component: () => import("../views/SignIn"),
    meta: {
      title: 'entrar',
      isPublic: true
    }
  },
  {
    path: "/saber-mais",
    name: "learn-more",
    component: () => import("../views/LearnMore"),
    meta: {
      title: 'saiba mais!',
      isPublic: true
    }
  },
  {
    path: '/criar-uma-conta',
    name: 'signup',
    component: () => import("../views/SignUp"),
    meta: {
      title: 'criar uma conta',
      isPublic: true
    }
  },
  {
    path: '/contribuir',
    name: 'payment',
    component: () => import("../views/Payment"),
    meta: {
      title: 'contribuição'
    }
  },
  {
    path: '/definir-renda',
    name: 'define-monthly-income',
    component: () => import("../views/DefineMonthlyIncome"),
    meta: {
      title: 'definir renda mensal'
    }
  },
  {
    path: '/meus-gastos',
    name: 'home',
    component: () => import("../views/Home"),
    meta: {
      title: 'meus gastos'
    }
  },
  {
    path: '/novo-gasto',
    name: 'add-expenses',
    component: () => import("../views/AddExpenses"),
    meta: {
      title: 'adicionar gastos'
    }
  },
  {
    path: '*',
    name: 'page-not-found',
    component: () => import("../views/PageNotFound"),
    meta: {
      title: 'página não encontrada'
    }
  },
];

const router = new VueRouter({
  routes
});

// Intercepting routes
router.beforeEach((to, from, next) => {
  // Start loading animation
  if (to.name) NProgress.start();

  // Change page's informations
  document.title = to.meta.title ? `porcool — ${to.meta.title}` : 'porcool';

  // Check route roles only if route is not public
  const isPrivate = to.matched.some(record => !record.meta.isPublic);
  if (isPrivate) {
    const unsubscribe = firebase.auth().onAuthStateChanged(async user => {
      // If we don't do this, the 'onAuthStateChanged()' will be executed a lot of times
      unsubscribe();

      // Check if user is logged in
      if (!user) {
        next({ name: 'signin' });
        return;
      }

      const loggedUser = await userService.get(user.uid);

      const payments = firebase.firestore().collection('payments');
      const userPayments = await payments
        .where('user', '==', user.uid)
        .orderBy('paymentDate', 'desc')
        .limit(1)
        .get();

      // Check if user payment is ok
      const remainingDays = !userPayments.empty ? paymentHelper.remainingDays(userPayments.docs[0].data().paymentDate) : 0;
      
      const userPaymentIsOk = (!userPayments.empty && remainingDays > 0);
      const userPaymentIsNotOk = (userPayments.empty || (!userPayments.empty && remainingDays <= 0));

      if (userPaymentIsNotOk) {
        next({ name: 'payment' });
        return;
      }

      // Check if user has a defined monthly income
      if (loggedUser.monthlyIncome === undefined && userPaymentIsOk) {
        next({ name: 'define-monthly-income' });
        return;
      }
    });
  }

  next();
});
router.afterEach((to, from) => {
  // Stop loading animation
  NProgress.done();
});

export default router;
