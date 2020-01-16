<template>
  <div class="columns is-desktop">
    <div class="column is-4-desktop">
      <LateralHeader />
    </div>
    <div class="column">
      <h1 class="title has-text-black">meus gastos</h1>
      <!-- TODO: check if this component will be rendered here and not inside it -->
      <FinishCurrentSpendingDate
        :expenses="expensesData.expenses"
        :showSpendingDateWarning="showSpendingDateWarning"
        :showResetExpensesWarning="showResetExpensesWarning"
      />
      <FilterByDate />
      <SpendingTable
        :expenses="expensesData.expenses"
        :loadingExpenses="expensesData.loadingExpenses"
        :loadingExpensesError="expensesData.loadingExpensesError"
      />
    </div>
  </div>
</template>

<script>
import firebase from "firebase/app";
import "firebase/auth";
import { mapState } from "vuex";
import moment from "moment";

// Components
import LateralHeader from "../components/LateralHeader";
import SpendingTable from "../components/SpendingTable";
import FilterByDate from "../components/FilterByDate";
import FinishCurrentSpendingDate from "../components/FinishCurrentSpendingDate";

// Services
import userService from "../services/user";

// Helpers
import dateAndTimeHelper from "../helpers/dateAndTime";

// Filters
import filters from "../filters";

export default {
  name: "Home",
  components: {
    LateralHeader,
    SpendingTable,
    FilterByDate,
    FinishCurrentSpendingDate
  },
  computed: {
    ...mapState({
      userData: state => state.user.user,
      expensesData: state => state.expenses
    })
  },
  data() {
    return {
      showSpendingDateWarning: false,
      showResetExpensesWarning: false
    };
  },
  created() {
    this.$store.dispatch("expenses/setExpenses", {
      userUid: this.userData.uid,
      spendingDate: this.userData.lookingAtSpendingDate
    });

    this.showSpendingDateWarning = moment().isAfter(
      this.userData.lookingAtSpendingDate,
      "months"
    );
    this.showResetExpensesWarning =
      moment().diff(this.userData.lookingAtSpendingDate, "months") >= 2;

    // setInterval(() => console.log(this.userData), 2000);
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