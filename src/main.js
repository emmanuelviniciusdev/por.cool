// General
import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";
import Buefy from "buefy";
import Vuelidate from 'vuelidate';
import VueCurrencyFilter from 'vue-currency-filter';

// Styles
import "./assets/scss/app.scss";

// Firebase
import firebase from 'firebase/app';
import 'firebase/analytics';
import 'firebase/auth';

// Plugins
Vue.use(Buefy, {
  defaultIconPack: "fas"
});
Vue.use(Vuelidate);
Vue.use(VueCurrencyFilter, {
  symbol: 'R$',
  thousandsSeparator: '.',
  fractionCount: 2,
  fractionSeparator: ',',
  symbolPosition: 'front',
  symbolSpacing: true
});

// Services
import userService from './services/user';

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
  data() {
    return {
      isReadyToRender: false
    }
  },
  beforeCreate() {
    const unsubscribe = firebase.auth().onAuthStateChanged(async user => {
      unsubscribe();

      if (user) {
        const loggedUser = await userService.get(user.uid);
        this.$store.dispatch('user/set', { uid: user.uid, ...loggedUser });
        this.isReadyToRender = true;
      }
    });
  },
  render(h) {
    if (this.isReadyToRender)
      return h(App);
  },
}).$mount("#app");