<template>
  <div>
    <b-table
      :data="data.table"
      :mobile-cards="false"
      hoverable
      paginated
      :per-page="10"
      :loading="loading"
    >
      <template slot-scope="props">
        <b-table-column
          field="expenseName"
          label="#"
          :class="{'has-text-weight-bold': props.row.type !== 'expense'}"
        >{{ props.row.expenseName }}</b-table-column>
        <b-table-column field="amount" label="gasto">{{ props.row.amount | currency }}</b-table-column>
        <b-table-column field="status" label="status">
          <b-tag :type="status_types[props.row.status]">
            <b-tooltip
              v-if="props.row.status === 'partially_paid'"
              :label="props.row.alreadyPaidAmount | currency"
              type="is-dark"
            >{{ status_labels[props.row.status] }}</b-tooltip>
            <span v-else>{{ status_labels[props.row.status] }}</span>
          </b-tag>
        </b-table-column>
        <b-table-column field="type" label="tipo">
          <b-tag size="is-medium">{{ types_labels[props.row.type] }}</b-tag>
        </b-table-column>
        <b-table-column field="action">
          <b-tooltip label="editar" type="is-dark">
            <button
              class="button is-warning is-small btn-table-action"
              :disabled="!canDeleteOrUpdateExpense(props.row.spendingDate)"
            >
              <b-icon icon="pencil-alt"></b-icon>
            </button>
          </b-tooltip>
          <b-tooltip label="remover" type="is-dark">
            <button
              class="button is-danger is-small btn-table-action"
              @click="removeExpense(props.row)"
              :disabled="!canDeleteOrUpdateExpense(props.row.spendingDate)"
            >
              <b-icon icon="trash"></b-icon>
            </button>
          </b-tooltip>
        </b-table-column>
      </template>

      <template slot="empty">
        <section class="section">
          <div class="content has-text-black has-text-centered">
            <div class="notification is-primary" v-if="!data.hasLoadingError">
              <p>
                Nenhum gasto adicionado para o mês de janeiro.
                <b>
                  <i>Legal.</i>
                </b>
              </p>
              <!-- TODO: Put a pig image here -->
            </div>

            <div class="notification is-danger" v-if="data.hasLoadingError">
              <p>Não foi possível carregar os seus gastos</p>
              <!-- TODO: Put a pig image here -->
            </div>
          </div>
        </section>
      </template>
    </b-table>
  </div>
</template>

<script>
import firebase from "firebase/app";
import "firebase/auth";
import { mapState } from "vuex";
import moment from "moment";
import FilterByDate from "./FilterByDate";

// Services
import expensesService from "../services/expenses";

// Helpers
import dateAndTimeHelper from "../helpers/dateAndTime";

export default {
  name: "SpendingTable",
  data() {
    return {
      data: {
        table: [],
        hasLoadingError: false
      },
      status_types: {
        paid: "is-success",
        partially_paid: "is-warning",
        pending: "is-danger"
      },
      status_labels: {
        paid: "pago",
        partially_paid: "parcialmente pago",
        pending: "pendente"
      },
      types_labels: {
        invoice: "fatura",
        expense: "gasto",
        savings: "poupança"
      },
      loading: false
    };
  },
  methods: {
    onLoading(state = true) {
      this.loading = state;
    },
    removeExpense(expense) {
      const {
        id: expenseDocId,
        expenseName,
        type: expenseType,
        spendingDate
      } = expense;

      if (!this.canDeleteOrUpdateExpense(spendingDate)) return;

      this.$buefy.dialog.confirm({
        title: `deletar ${this.types_labels[expenseType]}`,
        message: `Hey, você está prestes a deletar '<b>${expenseName}</b>' PARA SEMPRE. Você tem certeza de que deseja continuar?`,
        confirmText: "Sim. Deletar.",
        type: "is-danger",
        hasIcon: true,
        onConfirm: async () => {
          await expensesService.remove(expenseDocId);
          this.loadExpenses();
        }
      });
    },
    canDeleteOrUpdateExpense(spendingDate) {
      spendingDate = dateAndTimeHelper.transformSecondsToDate(
        spendingDate.seconds
      );
      return moment(this.userData.lookingAtSpendingDate).isSame(
        spendingDate,
        "month"
      );
    },
    async loadExpenses() {
      this.onLoading();

      try {
        const userExpenses = await expensesService.getAll(this.userData.uid);
        this.data.table = userExpenses;
      } catch (err) {
        // console.error(err);
        this.data.hasLoadingError = true;
      } finally {
        this.onLoading(false);
      }
    }
  },
  computed: {
    ...mapState({
      userData: state => state.user.user
    })
  },
  created() {
    this.loadExpenses();
  }
};
</script>

<style lang="scss" scoped>
.btn-table-action {
  margin: 2px 5px 2px 0;
}
</style>