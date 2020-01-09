// General
import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";
import Buefy from "buefy";
import Vuelidate from 'vuelidate';

// Styles
import "./assets/scss/app.scss";

// Firebase
import firebase from 'firebase/app';
import 'firebase/auth';
import 'firebase/analytics';

// Services
import userService from './services/user';

// Plugins
Vue.use(Buefy, {
  defaultIconPack: "fas"
});
Vue.use(Vuelidate);

firebase.initializeApp({
  apiKey: "REDACTED_PROD_API_KEY",
  authDomain: "REDACTED.firebaseapp.com",
  databaseURL: "https://REDACTED.firebaseio.com",
  projectId: "REDACTED",
  storageBucket: "REDACTED.appspot.com",
  messagingSenderId: "REDACTED_PROD_SENDER",
  appId: "1:REDACTED_PROD_SENDER:web:2322a618907f2774b8973f",
  measurementId: "REDACTED_PROD_MEASUREMENT"
});
firebase.analytics();

firebase.auth().onAuthStateChanged(async user => {
  if (user) {
    const loggedUser = await userService.get(user.uid);

    // Set user data in vuex
    store.dispatch('user/setUser', {
      uid: user.uid,
      displayName: user.displayName,
      name: loggedUser.name,
      lastName: loggedUser.lastName,
      email: loggedUser.email,
    });
  }
});

Vue.config.productionTip = false;

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount("#app");