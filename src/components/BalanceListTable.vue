<template>
  <div>
    <b-table
      :data="balances"
      :loading="loadingBalancesList"
      :mobile-cards="false"
      hoverable
      paginated
      :per-page="10"
    >
      <template slot-scope="props">
        <b-table-column field="balance" label="saldo">{{
          props.row.balance | currency
        }}</b-table-column>
        <b-table-column field="description" label="descrição">{{
          props.row.description
        }}</b-table-column>
        <b-table-column field="action" label="#">
          <b-tooltip label="remover" type="is-dark">
            <button
              class="button is-danger is-small"
              :disabled="!canDeleteOrUpdateBalance(props.row.spendingDate)"
              @click="deleteBalance(props.row)"
            >
              <b-icon icon="trash"></b-icon>
            </button>
          </b-tooltip>
        </b-table-column>
      </template>

      <template slot="empty">
        <section class="section">
          <div class="content has-text-black has-text-centered">
            <div v-if="!loadingBalancesListError">
              <div class="notification">
                Nenhum saldo foi encontrado por aqui...
              </div>
            </div>

            <div class="notification is-danger" v-if="loadingBalancesListError">
              <p>Não foi possível carregar os seus saldos</p>
            </div>
          </div>
        </section>
      </template>
    </b-table>
  </div>
</template>

<script>
import { mapState } from "vuex";
import moment from "moment";

// Filters
import filters from "../filters";

// Helpers
import dateAndTimeHelper from "../helpers/dateAndTime";

// Services
import balancesService from "../services/balances";

export default {
  name: "BalanceListTable",
  props: {
    balances: {
      required: true,
      type: Array
    },
    loadingBalancesList: {
      required: true,
      type: Boolean
    },
    loadingBalancesListError: {
      required: true,
      type: Boolean
    }
  },
  computed: {
    ...mapState({
      userData: state => state.user.user
    })
  },
  filters: {
    extractFromDateOnly: filters.extractFromDateOnly
  },
  methods: {
    canDeleteOrUpdateBalance(spendingDate) {
      spendingDate = dateAndTimeHelper.transformSecondsToDate(
        spendingDate.seconds
      );
      return moment(this.userData.lookingAtSpendingDate).isSame(
        spendingDate,
        "month"
      );
    },
    deleteBalance({ id, balance, spendingDate }) {
      if (!this.canDeleteOrUpdateBalance(spendingDate)) return;

      this.$buefy.dialog.confirm({
        title: "deletar saldo",
        message: `Hey, você está prestes a deletar um saldo de <b>${this.$options.filters.currency(
          balance
        )}</b> PARA SEMPRE. Você tem certeza de que deseja continuar?`,
        confirmText: "Sim. Deletar.",
        type: "is-danger",
        hasIcon: true,
        onConfirm: async () => {
          try {
            await balancesService.removeAdditionalBalance(id);

            this.$store.dispatch("balances/setBalances", {
              userUid: this.userData.uid,
              spendingDate: this.userData.lookingAtSpendingDate
            });
            this.$store.dispatch("balances/setBalancesList", {
              userUid: this.userData.uid,
              spendingDate: this.userData.lookingAtSpendingDate
            });
          } catch {
            this.$buefy.toast.open({
              message: "ocorreu um erro ao tentar deletar o seu saldo",
              type: "is-danger",
              position: "is-bottom"
            });
          }
        }
      });
    }
  }
};
</script>

<style></style>
