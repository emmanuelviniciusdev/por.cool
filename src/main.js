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

// Services
import userService from './services/user';

// Helpers
import dateAndTime from './helpers/dateAndTime';

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

// REVIEW
// It displays an error related to the buefy's components. Until now, there is no solution.
// https://github.com/vuetifyjs/vuetify/issues/9999
const ignoreWarnMessage = 'The .native modifier for v-on is only valid on components but it was used on <div>.';
Vue.config.warnHandler = function (msg, vm, trace) {
  // `trace` is the component hierarchy trace
  if (msg === ignoreWarnMessage) {
    msg = null;
    vm = null;
    trace = null;
  }
}

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

        if (loggedUser !== undefined) {
          const userLookingAtSpendingDate = dateAndTime.transformSecondsToDate(loggedUser.lookingAtSpendingDate.seconds);

          this.$store.dispatch('user/set', { uid: user.uid, displayName: user.displayName, ...loggedUser });
          this.$store.dispatch('expenses/setSpendingDatesList', {
            userUid: user.uid,
            lookingAtSpendingDate: userLookingAtSpendingDate
          });
          this.$store.dispatch('balances/setBalances', {
            userUid: user.uid,
            spendingDate: userLookingAtSpendingDate
          });
        }
      }

      this.isReadyToRender = true;
    });
  },
  render(h) {
    if (this.isReadyToRender)
      return h(App);
  },
}).$mount("#app");