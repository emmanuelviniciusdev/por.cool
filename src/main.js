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
  apiKey: "AIzaSyAMFqGceZwJ5srrPePLCXuDgaJlVXurVvI",
  authDomain: "por.cool",
  databaseURL: "https://porcool.firebaseio.com",
  projectId: "porcool",
  storageBucket: "porcool.appspot.com",
  messagingSenderId: "6802874030",
  appId: "1:6802874030:web:2322a618907f2774b8973f",
  measurementId: "G-NS5TW4HLVT"
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
  created() {
    firebase.auth().onAuthStateChanged(async user => {
      if (user) {
        const loggedUser = await userService.get(user.uid);
        this.$store.dispatch('user/set', { uid: user.uid, displayName: user.displayName, ...loggedUser });
      }

      this.isReadyToRender = true;
    });
  },
  render(h) {
    if (this.isReadyToRender)
      return h(App);
  },
}).$mount("#app");