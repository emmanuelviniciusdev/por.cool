<template>
  <div class="columns is-desktop">
    <div class="column is-4-desktop">
      <LateralHeader />
    </div>
    <div class="column">
      <h1 class="title has-text-black">meus gastos</h1>

      <FilterByDate />

      <button
        class="button is-warning spending-effective-date-btn"
        v-if="user.lookingAtSpendingDate"
      >
        <b-icon icon="hand-holding-usd"></b-icon>
        <span>fechar as contas para {{user.lookingAtSpendingDate.month + ' de ' + user.lookingAtSpendingDate.year }}</span>
      </button>

      <SpendingTable />
    </div>
  </div>
</template>

<script>
import firebase from "firebase/app";
import "firebase/auth";
import userService from "../services/user";
import dateAndTimeHelper from "../helpers/dateAndTime";
import LateralHeader from "../components/LateralHeader";
import SpendingTable from "../components/SpendingTable";
import FilterByDate from "../components/FilterByDate";

import moment from "moment";

export default {
  name: "Home",
  components: {
    LateralHeader,
    SpendingTable,
    FilterByDate
  },
  data() {
    return {
      user: {
        displayName: null,
        lookingAtSpendingDate: null
      }
    };
  },
  created() {
    const unsubscribe = firebase.auth().onAuthStateChanged(async user => {
      unsubscribe();

      if (user) {
        const loggedUser = await userService.get(user.uid);

        const lookingAtSpendingDate = dateAndTimeHelper.transformSecondsToDate(
          loggedUser.lookingAtSpendingDate.seconds
        );

        this.user.lookingAtSpendingDate = dateAndTimeHelper.extractOnly(
          lookingAtSpendingDate,
          ["year", "month"]
        );
        this.user.displayName = user.displayName;
      }
    });
  }
};
</script>

<style lang="scss" scoped>
.spending-effective-date-btn {
  margin: 10px 0 15px 0;
}

@media screen and (min-width: 1024px) {
  .columns {
    margin-top: 30px;
  }
}
</style>