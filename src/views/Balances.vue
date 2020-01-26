<template>
  <div>
    <div class="columns is-desktop">
      <div class="column is-4-desktop">
        <LateralHeader />
      </div>
      <div class="column">
        <h1 class="title has-text-black">
          meus saldos
          <Help>
            <template v-slot:title>meu saldos</template>
            <template v-slot:body>
              <div class="content">
                <h1 class="subtitle has-text-weight-normal">O que eu posso fazer aqui?</h1>
                <p
                  class="is-size-5 has-text-weight-normal"
                >Esta tela é para você que não recebe uma renda fixa mensal ou fez um dinheirinho extra durante o mês. Aqui, você poderá gerenciar todas as suas rendas extras.</p>
                <p
                  class="is-size-5 has-text-weight-normal"
                >Inclusive, é por aqui que você vai editar o valor da sua renda fixa mensal.</p>
              </div>
            </template>
          </Help>
        </h1>

        <FilterByDate @onDateChange="loadBalances" />

        <EditMonthlyIncome />

        <AddBalance />

        <BalanceListTable
          :balances="balances.balancesList"
          :loadingBalancesList="balances.loadingBalancesList"
          :loadingBalancesListError="balances.loadingBalancesListError"
        />
      </div>
    </div>
  </div>
</template>

<script>
import { mapState } from "vuex";

// Components
import LateralHeader from "../components/LateralHeader";
import Help from "../components/Help";
import FilterByDate from "../components/FilterByDate";
import AddBalance from "../components/AddBalance";
import EditMonthlyIncome from "../components/EditMonthlyIncome";
import BalanceListTable from "../components/BalanceListTable";

export default {
  name: "Balances",
  components: {
    LateralHeader,
    Help,
    FilterByDate,
    AddBalance,
    EditMonthlyIncome,
    BalanceListTable
  },
  data() {
    return {
      filterSpendingDate: null
    }
  },
  computed: {
    ...mapState({
      userData: state => state.user.user,
      balances: state => state.balances
    })
  },
  methods: {
    loadBalances(spendingDate = null) {
      this.$store.dispatch("balances/setBalancesList", {
        userUid: this.userData.uid,
        spendingDate: spendingDate
          ? spendingDate
          : this.userData.lookingAtSpendingDate
      });
    }
  },
  created() {
    this.loadBalances();
  }
};
</script>

<style lang="scss" scoped>
.AddBalance,
.EditMonthlyIncome {
  margin-bottom: 15px;
}

@media screen and (min-width: 1024px) {
  .columns {
    margin-top: 30px;
  }
}
</style>