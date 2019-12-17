import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";
import Buefy from "buefy";
import Vuelidate from 'vuelidate';
import firebase from 'firebase/app';
import 'firebase/analytics';
import "./assets/scss/app.scss";

// Plugins
Vue.use(Buefy, {
  defaultIconPack: "fas"
});
Vue.use(Vuelidate);

// Initialize Firebase
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

Vue.config.productionTip = false;

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount("#app");
