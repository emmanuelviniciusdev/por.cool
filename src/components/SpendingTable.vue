<template>
  <div>
    <EditExpense :expense="expenseToEdit" />

    <b-table
      :data="expenses"
      :loading="loadingExpenses"
      :mobile-cards="false"
      hoverable
      paginated
      :per-page="10"
    >
      <template slot-scope="props">
        <b-table-column
          field="expenseName"
          label="#"
          :class="{'has-text-weight-bold': props.row.type !== 'expense'}"
        >{{ props.row.expenseName }}</b-table-column>
        <b-table-column
          field="amount"
          label="gasto"
        >{{ props.row.amount | sumAmounts(props.row) | currency }}</b-table-column>
        <b-table-column field="status" label="status">
          <b-button
            :type="status_types[props.row.status]"
            size="is-small"
            @dblclick="fastChangeStatusExpense(props.row)"
            :loading="loadingFastChangeStatus"
            :disabled="!canDeleteOrUpdateExpense(props.row.spendingDate)"
          >
            <b-tooltip
              v-if="props.row.status === 'partially_paid'"
              :label="props.row.alreadyPaidAmount | currency"
              type="is-dark"
            >{{ status_labels[props.row.status] }}</b-tooltip>
            <span v-else>{{ status_labels[props.row.status] }}</span>
          </b-button>
        </b-table-column>
        <b-table-column field="type" label="tipo">
          <b-tag size="is-medium">{{ types_labels[props.row.type] }}</b-tag>
        </b-table-column>
        <b-table-column field="action">
          <b-tooltip label="editar" type="is-dark">
            <button
              class="button is-warning is-small btn-table-action"
              :disabled="!canDeleteOrUpdateExpense(props.row.spendingDate)"
              @click.prevent="editExpense(props.row)"
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
            <div class="stoincs" v-if="!loadingExpensesError">
              <div class="notification">
                <p>
                  Nenhum gasto encontrado para {{this.userData.lookingAtSpendingDate | extractFromDateOnly('month') }} de {{this.userData.lookingAtSpendingDate | extractFromDateOnly('year') }}.
                  <b>
                    <i>Legal.</i>
                  </b>
                </p>
              </div>
              <img src="../assets/images/stoincs.png" alt="stoincs" />
            </div>

            <div class="notification is-danger" v-if="loadingExpensesError">
              <p>Não foi possível carregar os seus gastos</p>
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

// Components
import EditExpense from "./EditExpense";

// Services
import expensesService from "../services/expenses";

// Helpers
import dateAndTimeHelper from "../helpers/dateAndTime";

// Filters
import filters from "../filters";

// Mixins
import SpendingTableMixin from "../mixins/SpendingTable";

export default {
  name: "SpendingTable",
  props: {
    expenses: {
      required: true,
      type: Array
    },
    loadingExpenses: {
      required: true,
      type: Boolean
    },
    loadingExpensesError: {
      required: true,
      type: Boolean
    }
  },
  components: {
    EditExpense
  },
  data() {
    return {
      expenseToEdit: {},
      loadingFastChangeStatus: false
    };
  },
  mixins: [SpendingTableMixin],
  methods: {
    editExpense(expense) {
      if (!this.canDeleteOrUpdateExpense(expense.spendingDate)) return;

      this.expenseToEdit = expense;
      // It is to Vue capture changes in 'this.expenseToEdit' and open edit modal
      // without any problem.
      this.expenseToEdit = { ...this.expenseToEdit, watchKey: Math.random() };
    },
    async fastChangeStatusExpense(expense) {
      if (!this.canDeleteOrUpdateExpense(expense.spendingDate)) return;

      const { status } = expense;

      if (status === "partially_paid") return;

      this.loadingFastChangeStatus = true;

      expense.status = status === "paid" ? "pending" : "paid";

      try {
        await expensesService.update(expense);
        this.$store.dispatch("expenses/setExpenses", {
          userUid: this.userData.uid,
          spendingDate: this.userData.lookingAtSpendingDate
        });
      } catch (err) {
        console.log(err);
        this.$buefy.toast.open({
          message: "ocorreu um erro ao fazer a troca rápida de status",
          position: "is-bottom",
          type: "is-danger"
        });
      } finally {
        this.loadingFastChangeStatus = false;
      }
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
          this.$store.dispatch("expenses/setExpenses", {
            userUid: this.userData.uid,
            spendingDate: this.userData.lookingAtSpendingDate
          });
        }
      });
    },
    // When expense date is different to 'lookingAtSpendingDate',
    // that is the current spending month of the user, it returns false
    // and means that user can't do anything at all with his expenses.
    canDeleteOrUpdateExpense(spendingDate) {
      spendingDate = dateAndTimeHelper.transformSecondsToDate(
        spendingDate.seconds
      );
      return moment(this.userData.lookingAtSpendingDate).isSame(
        spendingDate,
        "month"
      );
    }
  },
  computed: {
    ...mapState({
      userData: state => state.user.user
    })
  },
  filters: {
    extractFromDateOnly: filters.extractFromDateOnly,
    sumAmounts: filters.sumAmounts
  }
};
</script>

<style lang="scss" scoped>
.btn-table-action {
  margin: 2px 5px 2px 0;
}

.stoincs img {
  width: 450px;
}
</style>