<template>
  <div>
    <b-modal :active="isModalOpened" full-screen :canCancel="false">
      <div class="modal-card" style="width: auto">
        <header class="modal-card-head">
          <p
            class="modal-card-title"
          >{{ `${formatedUserLookingAtSpendingDate.month} de ${formatedUserLookingAtSpendingDate.year}` }}</p>
        </header>
        <section class="modal-card-body">
          <h1 class="title">antes de prosseguir</h1>
          <h2 class="subtitle">
            o que você pretende fazer com estes gastos que ficaram
            <i>pendentes</i> ou
            <i>parcialmente pagos</i>?
          </h2>
          <b-table :data="modifiedExpenses" :mobile-cards="true" hoverable paginated :per-page="10">
            <template slot-scope="props">
              <b-table-column
                field="expenseName"
                label="#"
                :class="{'has-text-weight-bold': props.row.type !== 'expense'}"
              >{{ props.row.expenseName }}</b-table-column>
              <b-table-column field="amount" label="gasto">{{ props.row.amount | sumAmounts(props.row) | currency }}</b-table-column>
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
              <b-table-column field="status" label="já está pago?">
                <div class="field">
                  <div class="control has-icons-left">
                    <div class="select">
                      <select
                        @change="setExpenseTemporaryValue(props.index, 'hasUserPaid', $event.target.value === 'true')"
                      >
                        <option value="false" selected>não</option>
                        <option value="true">sim</option>
                      </select>
                    </div>
                    <div class="icon is-small is-left">
                      <i
                        :class="{
                        'fas fa-thumbs-up': props.row.temporary.hasUserPaid,
                        'fas fa-thumbs-down': !props.row.temporary.hasUserPaid,
                      }"
                      ></i>
                    </div>
                  </div>
                </div>
              </b-table-column>
              <b-table-column field="action">
                <div class="field" style="width: 300px;" v-if="!props.row.temporary.hasUserPaid">
                  <div class="control has-icons-left">
                    <div class="select">
                      <select
                        @change="setExpenseTemporaryValue(props.index, 'action', $event.target.value)"
                      >
                        <option value="nothing" selected>o que pretende fazer, então?</option>
                        <option
                          value="move_on"
                          v-if="props.row.type === 'expense'"
                        >passar pro próximo mês</option>
                        <option
                          value="move_on_with_difference"
                          v-if="props.row.type !== 'expense'"
                        >passar pro próximo mês (com a diferença)</option>
                        <option
                          value="move_on_without_difference"
                          v-if="props.row.type !== 'expense'"
                        >passar pro próximo mês (sem a diferença)</option>
                        <option value="discard">descartar</option>
                      </select>
                    </div>
                    <div class="icon is-small is-left">
                      <i
                        :class="{
                        'fas fa-question-circle': props.row.temporary.action === 'nothing',
                        'fas fa-dollar-sign': props.row.temporary.action === 'move_on_with_difference',
                        'fas fa-angle-double-right': ['move_on', 'move_on_without_difference'].includes(props.row.temporary.action ),
                        'fas fa-trash-alt': props.row.temporary.action === 'discard',
                      }"
                      ></i>
                    </div>
                  </div>
                </div>
              </b-table-column>
            </template>
          </b-table>
        </section>
        <footer class="modal-card-foot">
          <button class="button" @click="closeModal()">cancelar</button>
          <b-button
            type="is-primary"
            @click="finishCurrentSpendingDate()"
            :disabled="!haveExpensesSettedActions"
            :loading="loadingFinishSpendingDate"
          >fechar gastos</b-button>
        </footer>
      </div>
    </b-modal>

    <div class="hero is-warning spendingDateWarning">
      <div class="hero-body">
        <div class="container">
          <h1 class="title">{{ this.userData.displayName | capitalizeName }},</h1>
          <h2 class="subtitle">
            você já pode fechar os gastos para
            <b>{{ `${formatedUserLookingAtSpendingDate.month} de ${formatedUserLookingAtSpendingDate.year}` }}</b>.
          </h2>
          <b-button
            type="is-light"
            @click="mayOpenModal()"
            :loading="loadingFinishSpendingDate"
          >fechar gastos</b-button>
          <div v-if="showResetExpensesWarning">
            <hr />
            <p>
              <i>
                Notamos que você não utiliza o porcool já faz alguns meses.
                Se você já não se lembra mais do que gastou dentro deste intervalo de tempo e gostaria de fazer um reset de tudo para recomeçar
                a fazer o seu controle financeiro do zero,
                <a
                  href="#"
                >clique aqui</a> ou acesse
                <b>minha conta > recomeçar do zero.</b>
              </i>
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { mapState } from "vuex";
import moment from "moment";

// Filters
import filters from "../filters";

// Services
import expensesService from "../services/expenses";

// Helpers
import dateAndTimeHelper from "../helpers/dateAndTime";

// Mixins
import SpendingTableMixin from "../mixins/SpendingTable";

import Vue from "vue";

export default {
  name: "FinishCurrentSpendingDate",
  props: {
    expenses: {
      type: Array,
      required: true
    }
  },
  mixins: [SpendingTableMixin],
  data() {
    return {
      isModalOpened: false,
      loadingFinishSpendingDate: false,
      modifiedExpenses: []
    };
  },
  computed: {
    ...mapState({
      userData: state => state.user.user
    }),
    haveExpensesSettedActions() {
      return (
        this.modifiedExpenses.filter(
          ({ temporary }) =>
            !temporary.hasUserPaid && temporary.action === "nothing"
        ).length === 0
      );
    },
    formatedUserLookingAtSpendingDate() {
      return dateAndTimeHelper.extractOnly(
        this.userData.lookingAtSpendingDate,
        ["year", "month"]
      );
    },
    newSpendingDate() {
      return dateAndTimeHelper.startOfMonthAndDay(
        moment(this.userData.lookingAtSpendingDate)
          .add(1, "months")
          .toDate()
      );
    },
    showResetExpensesWarning() {
      return moment().diff(this.userData.lookingAtSpendingDate, "months") >= 2;
    }
  },
  filters: {
    capitalizeName: filters.capitalizeName,
    sumAmounts: filters.sumAmounts
  },
  methods: {
    onLoadingFinishSpendingDate(state = true) {
      this.loadingFinishSpendingDate = state;
    },
    setExpenseTemporaryValue(index, key, value) {
      this.modifiedExpenses[index].temporary[key] = value;
      this.modifiedExpenses.splice(index, 1, {
        ...this.modifiedExpenses[index]
      });
    },
    async mayOpenModal() {
      const filteredExpenses = this.expenses.filter(
        ({ status }) => status !== "paid"
      );

      if (filteredExpenses.length > 0) {
        this.modifiedExpenses = [
          ...filteredExpenses.map(expense => {
            expense.temporary = {};
            expense.temporary.action = "nothing";
            expense.temporary.hasUserPaid = false;
            return expense;
          })
        ];

        this.isModalOpened = true;

        return;
      }

      this.onLoadingFinishSpendingDate();

      // Finish current spending date with auto clone
      await expensesService.finishCurrentSpendingDate(
        this.userData.uid,
        this.userData.lookingAtSpendingDate,
        { autoClone: true }
      );

      this.updateGeneral();
      this.onLoadingFinishSpendingDate(false);
    },
    closeModal() {
      this.isModalOpened = false;
    },
    async finishCurrentSpendingDate() {
      if (!this.haveExpensesSettedActions) return;

      this.onLoadingFinishSpendingDate();

      const paidExpenses = [];
      const changingExpenses = {
        expensesToUpdateNow: [],
        expensesToClone: []
      };

      this.modifiedExpenses.forEach(expense => {
        const { hasUserPaid, action } = expense.temporary;

        // This props will be the same to each changing expense, so we putted it
        // into an object.
        const expensePropsToChange = {
          alreadyPaidAmount: 0,
          status: "pending",
          created: new Date(),
          spendingDate: this.newSpendingDate
        };

        if (hasUserPaid) {
          paidExpenses.push({
            ...expense,
            status: "paid",
            alreadyPaidAmount: 0
          });

          // Check if expense has type "invoice" or "savings".
          // If so, we'll clone the expense if it's in force.
          if (expense.type !== "expense") {
            const isInForce = expense.indeterminateValidity
              ? true
              : moment(
                  dateAndTimeHelper.transformSecondsToDate(
                    expense.validity.seconds
                  )
                ).isSameOrAfter(moment(this.newSpendingDate));
            if (isInForce) {
              changingExpenses.expensesToClone.push({
                ...expense,
                ...expensePropsToChange,
                differenceAmount: 0
              });
            }
          }
        } else {
          // Every expense will have its amount changed to 0 or to 'alreadyPaidAmount' value
          // in the current spending date.
          changingExpenses.expensesToUpdateNow.push({
            ...expense,
            amount:
              expense.status === "partially_paid"
                ? expense.alreadyPaidAmount
                : 0,
            alreadyPaidAmount: 0,
            differenceAmount: 0,
            status: "paid"
          });

          // Only expense with "expense" type will have "move_on" action.
          if (action === "move_on") {
            changingExpenses.expensesToClone.push({
              ...expense,
              ...expensePropsToChange,
              amount: expense.amount - expense.alreadyPaidAmount
            });
          } else if (action === "move_on_without_difference") {
            changingExpenses.expensesToClone.push({
              ...expense,
              ...expensePropsToChange,
              differenceAmount: 0
            });
          } else if (action === "move_on_with_difference") {
            const pastDifferenceAmount =
              expense.differenceAmount !== undefined
                ? expense.differenceAmount
                : 0;

            changingExpenses.expensesToClone.push({
              ...expense,
              ...expensePropsToChange,
              differenceAmount:
                expense.amount -
                expense.alreadyPaidAmount +
                pastDifferenceAmount
            });
          }
        }
      });

      // Remove useless props without reflect on the original object.
      // Just in case.
      function removeUselessProps(obj) {
        const newObj = { ...obj };
        delete newObj.temporary;
        return newObj;
      }

      // Update "paid" expenses
      await expensesService.bulkUpdate(paidExpenses.map(removeUselessProps));

      // Update "amount" and "status" of expenses before cloning them to the
      // next spending date
      await expensesService.bulkUpdate(
        changingExpenses.expensesToUpdateNow.map(removeUselessProps)
      );

      // Clone expenses to the next spending date
      await expensesService.insert(
        changingExpenses.expensesToClone.map(removeUselessProps)
      );

      // Finish current spending date without auto clone
      await expensesService.finishCurrentSpendingDate(
        this.userData.uid,
        this.userData.lookingAtSpendingDate
      );

      this.updateGeneral();
      this.onLoadingFinishSpendingDate(false);
    },
    updateGeneral() {
      this.closeModal();

      this.$store.dispatch("expenses/setExpenses", {
        userUid: this.userData.uid,
        spendingDate: this.newSpendingDate
      });

      // IT MUST BE AT THE END!
      // Because if we update user's spending date first, the 'newSpendingDate' will
      // automatically react to the new setted 'lookingAtSpendingDate' and get the wrong
      // expenses with the wrong new spending date.
      this.$store.dispatch("user/update", {
        lookingAtSpendingDate: this.newSpendingDate
      });
    }
  }
};
</script>

<style lang="scss" scoped>
.spendingDateWarning {
  margin-bottom: 20px;
}
</style>