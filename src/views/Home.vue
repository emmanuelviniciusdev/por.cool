<template>
  <div class="columns is-desktop">
    <div class="column is-4-desktop">
      <LateralHeader />
    </div>
    <div class="column">
      <div class="title-and-help">
        <h1 class="title has-text-black">meus gastos</h1>
        <Help>
          <template v-slot:title>meus gastos</template>
          <template v-slot:body>
            <div class="content">
              <p
                class="is-size-5 has-text-weight-normal"
              >Nesta tela, você terá uma listagem de todos os seus gastos que foram adicionados, referentes ao respectivo ano e mês selecionados.</p>
              <p class="is-size-5 has-text-weight-normal">
                <b>Dicas:</b>
                <ul>
                  <li>
                    <p class="is-size-6 has-text-weight-normal">
                      Você sabia que se clicar 2x em cima do botão de status do seu gasto, o status muda? Mas atenção: ele só vai de "pendente" para "pago" e vice-versa.
                      <b>Não funciona caso o status esteja como "parcialmente pago".</b>
                    </p>
                    <img src="https://i.imgur.com/u5xNOv2.png" alt="clique 2x em cima do botão de status e veja a mágica acontecer">
                  </li>
                </ul>
              </p>
            </div>
          </template>
        </Help>
      </div>
      <FinishCurrentSpendingDate v-if="showSpendingDateWarning" />
      <FilterByDate @onDateChange="loadExpense" />
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
import Help from "../components/Help";

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
    FinishCurrentSpendingDate,
    Help
  },
  computed: {
    ...mapState({
      userData: state => state.user.user,
      expensesData: state => state.expenses
    }),
    showSpendingDateWarning() {
      return moment().isAfter(this.userData.lookingAtSpendingDate, "months");
    }
  },
  methods: {
    loadExpense(spendingDate = null) {
      this.$store.dispatch("expenses/setExpenses", {
        userUid: this.userData.uid,
        spendingDate: spendingDate
          ? spendingDate
          : this.userData.lookingAtSpendingDate
      });
    }
  },
  created() {
    this.loadExpense();
  }
};
</script>

<style lang="scss" scoped>
.spending-effective-date-btn {
  margin: 10px 0 15px 0;
}

.title-and-help {
  h1 {
    display: inline;
  }

  margin-bottom: 15px;
}

@media screen and (min-width: 1024px) {
  .columns {
    margin-top: 30px;
  }
}
</style>