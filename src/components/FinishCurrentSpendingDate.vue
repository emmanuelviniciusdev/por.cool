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
          <button
            class="button is-primary"
            @click="finishCurrentSpendingDate()"
            :disabled="!haveExpensesSettedActions"
          >fechar gastos</button>
        </footer>
      </div>
    </b-modal>

    <div class="hero is-warning spendingDateWarning" v-if="showSpendingDateWarning">
      <div class="hero-body">
        <div class="container">
          <h1 class="title">{{ this.userData.displayName | capitalizeName }},</h1>
          <h2 class="subtitle">
            você já pode fechar os gastos para
            <b>{{ `${formatedUserLookingAtSpendingDate.month} de ${formatedUserLookingAtSpendingDate.year}` }}</b>.
          </h2>
          <button class="button is-light" @click="mayOpenModal()">fechar gastos</button>
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
    showSpendingDateWarning: {
      type: Boolean,
      required: true
    },
    showResetExpensesWarning: {
      type: Boolean,
      required: true
    },
    expenses: {
      type: Array,
      required: true
    }
  },
  mixins: [SpendingTableMixin],
  data() {
    return {
      formatedUserLookingAtSpendingDate: null,
      isModalOpened: false,
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
    }
  },
  filters: {
    capitalizeName: filters.capitalizeName
  },
  methods: {
    setExpenseTemporaryValue(index, key, value) {
      this.modifiedExpenses[index].temporary[key] = value;
      this.modifiedExpenses.splice(index, 1, {
        ...this.modifiedExpenses[index]
      });
    },
    mayOpenModal(state = true) {
      console.log(this.expenses);
      const filteredExpenses = this.expenses.filter(
        ({ status }) => status === "pending" || status === "partially_paid"
      );

      if (state && filteredExpenses.length > 0) {
        this.modifiedExpenses = [
          ...filteredExpenses.map(expense => {
            expense.temporary = {};
            expense.temporary.action = "nothing";
            expense.temporary.hasUserPaid = false;
            return expense;
          })
        ];

        this.isModalOpened = state;

        return;
      }

      this.finishCurrentSpendingDate();
    },
    closeModal() {
      this.isModalOpened = false;
    },
    async finishCurrentSpendingDate() {
      if (!this.haveExpensesSettedActions) return;

      const paidExpenses = [];
      const notPaidExpenses = {
        moveOn: {
          // Expenses that we're gonna update only 'amount' to zero
          expensesToUpdateNow: [],
          // Expenses that we're gonna clone to the next month
          expensesToClone: []
        },
        moveOnWithDifference: {
          // Expenses that we're gonna update only 'amount' to zero
          expensesToUpdateNow: [],
          // Expenses that we're gonna clone to the next month
          expensesToClone: []
        },
        moveOnWithoutDifference: {
          // Expenses that we're gonna update only 'amount' to zero
          expensesToUpdateNow: [],
          // Expenses that we're gonna clone to the next month
          expensesToClone: []
        },
        discard: {
          // Expenses that we're gonna update only 'amount' to zero
          expensesToUpdateNow: []
        }
      };

      this.modifiedExpenses.forEach(expense => {
        const { hasUserPaid, action } = expense.temporary;

        // TODO: Improve this

        if (hasUserPaid) {
          paidExpenses.push({ ...expense, status: "paid" });
        } else {
          const newSpendingDate = dateAndTimeHelper.startOfMonthAndDay(
            moment(this.userData.lookingAtSpendingDate).add(1, "months")
          );

          if (action === "move_on") {
            notPaidExpenses.moveOn.expensesToUpdateNow.push({
              docId: expense.id,
              amount:
                expense.status === "partially_paid"
                  ? expense.alreadyPaidAmount
                  : 0
            });

            notPaidExpenses.moveOn.expensesToClone.push({
              ...expense,
              amount: expense.amount - expense.alreadyPaidAmount,
              alreadyPaidAmount: 0,
              status: "pending",
              created: new Date(),
              spendingDate: newSpendingDate
            });
          } else if (action === "move_on_without_difference") {
            notPaidExpenses.moveOnWithoutDifference.expensesToUpdateNow.push({
              docId: expense.id,
              amount:
                expense.status === "partially_paid"
                  ? expense.alreadyPaidAmount
                  : 0
            });

            notPaidExpenses.moveOnWithoutDifference.expensesToClone.push({
              ...expense,
              alreadyPaidAmount: 0,
              status: "pending",
              created: new Date(),
              spendingDate: newSpendingDate
            });
          } else if (action === "move_on_with_difference") {
            notPaidExpenses.moveOnWithDifference.expensesToUpdateNow.push({
              docId: expense.id,
              amount:
                expense.status === "partially_paid"
                  ? expense.alreadyPaidAmount
                  : 0
            });

            notPaidExpenses.moveOnWithDifference.expensesToClone.push({
              ...expense,
              differenceAmount: expense.amount - expense.alreadyPaidAmount,
              alreadyPaidAmount: 0,
              status: "pending",
              created: new Date(),
              spendingDate: newSpendingDate
            });
          } else if (action === "discard") {
            notPaidExpenses.discard.expensesToUpdateNow.push({
              docId: expense.id,
              amount:
                expense.status === "partially_paid"
                  ? expense.alreadyPaidAmount
                  : 0
            });
          }
        }
      });

      console.log(notPaidExpenses);
    }
  },
  created() {
    this.formatedUserLookingAtSpendingDate = dateAndTimeHelper.extractOnly(
      this.userData.lookingAtSpendingDate,
      ["year", "month"]
    );
  }
};
</script>

<style lang="scss" scoped>
.spendingDateWarning {
  margin-bottom: 20px;
}
</style>