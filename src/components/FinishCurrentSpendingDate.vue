<template>
  <div>
    <b-modal :active="isModalOpened" full-screen :canCancel="false">
      <div class="modal-card" style="width: auto">
        <header class="modal-card-head">
          <p class="modal-card-title">
            {{
              `${formatedUserLookingAtSpendingDate.month} de ${formatedUserLookingAtSpendingDate.year}`
            }}
          </p>
        </header>
        <section class="modal-card-body">
          <h1 class="title">
            antes de prosseguir
            <Help tooltipPosition="is-bottom">
              <template v-slot:title
                >como funciona o fechamento de gastos?</template
              >
              <template v-slot:body>
                <div class="content">
                  <h1 class="subtitle has-text-weight-normal">
                    É bizarro de simples.
                  </h1>
                  <p class="is-size-5 has-text-weight-normal">
                    Esta tela funciona assim: toda vez que termina um mês e
                    entra outro, você precisa "fechar os gastos" para o mês
                    antigo antes de começar a adicionar novos gastos para o novo
                    mês.
                    <i>
                      Se você continuar adicionando gastos sem executar o
                      fechamento de gastos, você estará adicionando gastos para
                      o mês antigo, e não para o novo mês.
                    </i>
                  </p>
                  <p class="is-size-5 has-text-weight-normal">
                    Se você caiu nesta tela, quer dizer que quando você foi
                    executar o fechamento de gastos você ainda tinha alguns
                    gastos cujo os status estavam "pendentes" ou "parcialmente
                    pagos".
                    <i>
                      Se todos os seus gastos estivessem com status "pago", você
                      não cairia nesta tela e o fechamento seria realizado de
                      forma automática.
                    </i>
                  </p>
                  <p class="is-size-5 has-text-weight-bold">
                    Entendido como e por que você foi parar aqui, vamos ao que
                    interessa:
                  </p>
                  <p class="is-size-5 has-text-weight-normal">
                    Se o seu gasto já está pago e você apenas se esqueceu de
                    atualizá-lo ao clicar em fechar gastos, você pode marcar a
                    opção de
                    <b>"já está pago?"</b> para <b>sim</b>, e o status deste
                    gasto simplesmente se atualizará para "pago" e o fechamento
                    ocorrerá normalmente, de forma automática.
                    <img src="https://i.imgur.com/insJhSu.png" />
                  </p>
                  <p class="is-size-5 has-text-weight-normal">
                    Porém, se você de fato não efetuou o pagamento do seu gasto
                    em questão, você possui algumas alternativas:
                  </p>
                  <ul>
                    <li>
                      <p class="is-size-6 has-text-weight-normal">
                        <b>passar pro próximo mês:</b> irá passar o gasto com o
                        valor que ainda não foi pago para o próximo mês.
                      </p>
                    </li>
                    <li>
                      <p class="is-size-6 has-text-weight-normal">
                        <b>passar pro próximo mês (com a diferença):</b> irá
                        passar a fatura ou poupança para o próximo mês com a
                        diferença de valor que ainda não foi paga.
                      </p>
                    </li>
                    <li>
                      <p class="is-size-6 has-text-weight-normal">
                        <b>passar pro próximo mês (sem a diferença):</b> irá
                        passar a fatura ou poupança para o próximo mês sem
                        diferença nenhuma de valor. Ou seja: com o seu valor
                        original.
                      </p>
                    </li>
                    <li>
                      <p class="is-size-6 has-text-weight-normal">
                        <b>descartar:</b> não irá passar para o próximo mês.
                      </p>
                    </li>
                  </ul>
                  <p class="is-size-5 has-text-weight-normal">
                    <b>Observação:</b> todas as opções acima irão zerar os
                    valores dos seus gastos para o mês antigo se o status
                    estiver como "pendente". Se o status for "parcialmente
                    pago", elas irão considerar e definir os valores dos gastos
                    para os valores que já haviam sido pagos.
                  </p>
                </div>
              </template>
            </Help>
          </h1>
          <h2 class="subtitle">
            o que você pretende fazer com estes gastos que ficaram
            <i>pendentes</i> ou <i>parcialmente pagos</i>?
          </h2>
          <b-table
            :data="
              modifiedExpenses.filter(
                ({ status, type }) => !(status === 'paid' && type !== 'expense')
              )
            "
            :mobile-cards="true"
            hoverable
            paginated
            :per-page="10"
          >
            <template slot-scope="props">
              <b-table-column
                field="expenseName"
                label="#"
                :class="{
                  'has-text-weight-bold': props.row.type !== 'expense'
                }"
                >{{ props.row.expenseName }}</b-table-column
              >
              <b-table-column field="amount" label="gasto">{{
                props.row.amount | sumAmounts(props.row) | currency
              }}</b-table-column>
              <b-table-column field="status" label="status">
                <b-tag :type="status_types[props.row.status]">
                  <b-tooltip
                    v-if="props.row.status === 'partially_paid'"
                    :label="props.row.alreadyPaidAmount | currency"
                    type="is-dark"
                    >{{ status_labels[props.row.status] }}</b-tooltip
                  >
                  <span v-else>{{ status_labels[props.row.status] }}</span>
                </b-tag>
              </b-table-column>
              <b-table-column field="type" label="tipo">
                <b-tag size="is-medium">{{
                  types_labels[props.row.type]
                }}</b-tag>
              </b-table-column>
              <b-table-column field="status" label="já está pago?">
                <div class="field">
                  <div class="control has-icons-left">
                    <div class="select">
                      <select
                        @change="
                          setExpenseTemporaryValue(
                            props.row.temporary.pseudoId,
                            'hasUserPaid',
                            $event.target.value === 'true'
                          )
                        "
                      >
                        <option value="false" selected>não</option>
                        <option value="true">sim</option>
                      </select>
                    </div>
                    <div class="icon is-small is-left">
                      <i
                        :class="{
                          'fas fa-thumbs-up': props.row.temporary.hasUserPaid,
                          'fas fa-thumbs-down': !props.row.temporary.hasUserPaid
                        }"
                      ></i>
                    </div>
                  </div>
                </div>
              </b-table-column>
              <b-table-column field="action">
                <div
                  class="field"
                  style="width: 300px;"
                  v-if="!props.row.temporary.hasUserPaid"
                >
                  <div class="control has-icons-left">
                    <div class="select">
                      <select
                        @change="
                          setExpenseTemporaryValue(
                            props.row.temporary.pseudoId,
                            'action',
                            $event.target.value
                          )
                        "
                      >
                        <option value="nothing" selected
                          >o que pretende fazer, então?</option
                        >
                        <option
                          value="move_on"
                          v-if="props.row.type === 'expense'"
                          >passar pro próximo mês</option
                        >
                        <option
                          value="move_on_with_difference"
                          v-if="props.row.type !== 'expense'"
                          >passar pro próximo mês (com a diferença)</option
                        >
                        <option
                          value="move_on_without_difference"
                          v-if="props.row.type !== 'expense'"
                          >passar pro próximo mês (sem a diferença)</option
                        >
                        <option value="discard">descartar</option>
                      </select>
                    </div>
                    <div class="icon is-small is-left">
                      <i
                        :class="{
                          'fas fa-question-circle':
                            props.row.temporary.action === 'nothing',
                          'fas fa-dollar-sign':
                            props.row.temporary.action ===
                            'move_on_with_difference',
                          'fas fa-angle-double-right': [
                            'move_on',
                            'move_on_without_difference'
                          ].includes(props.row.temporary.action),
                          'fas fa-trash-alt':
                            props.row.temporary.action === 'discard'
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
            >fechar gastos</b-button
          >
        </footer>
      </div>
    </b-modal>

    <div class="hero is-warning spendingDateWarning">
      <div class="hero-body">
        <div class="container">
          <h1 class="title">
            {{ this.userData.displayName | capitalizeName }},
          </h1>
          <h2 class="subtitle">
            você já pode fechar os gastos para
            <b>{{
              `${formatedUserLookingAtSpendingDate.month} de ${formatedUserLookingAtSpendingDate.year}`
            }}</b
            >.
          </h2>
          <b-button
            type="is-light"
            @click="mayOpenModal()"
            :loading="loadingFinishSpendingDate"
            >fechar gastos</b-button
          >
          <div v-if="showResetExpensesWarning">
            <hr />
            <p>
              <i>
                Notamos que você não utiliza o porcool já faz alguns meses. Se
                você já não se lembra mais do que gastou dentro deste intervalo
                de tempo e gostaria de fazer um reset de tudo para recomeçar a
                fazer o seu controle financeiro do zero,
                <a href="#">clique aqui</a> ou acesse
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

// Components
import Help from "./Help";

// Filters
import filters from "../filters";

// Services
import expensesService from "../services/expenses";
import balancesServices from "../services/balances";

// Helpers
import dateAndTimeHelper from "../helpers/dateAndTime";

// Mixins
import SpendingTableMixin from "../mixins/SpendingTable";

import Vue from "vue";

export default {
  name: "FinishCurrentSpendingDate",
  components: {
    Help
  },
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
    setExpenseTemporaryValue(pseudoId, key, value) {
      const triggeredExpense = this.modifiedExpenses.filter(
        ({ temporary }) => temporary.pseudoId === pseudoId
      )[0];
      const indexOfTriggeredExpense = this.modifiedExpenses.indexOf(
        triggeredExpense
      );

      if (indexOfTriggeredExpense !== -1) {
        this.modifiedExpenses[indexOfTriggeredExpense].temporary[key] = value;

        this.modifiedExpenses.splice(indexOfTriggeredExpense, 1, {
          ...this.modifiedExpenses[indexOfTriggeredExpense]
        });
      }
    },
    async mayOpenModal() {
      this.onLoadingFinishSpendingDate();

      // When we're going to finish spending date we can't access vuex's expenses
      // because they are mutable and can be different from the expenses of current
      // spending date that user is looking at. So, unfortunately, we have to make
      // another server request.
      const expenses = await expensesService.getAll(
        this.userData.uid,
        this.userData.lookingAtSpendingDate
      );

      const filteredExpenses = expenses.filter(
        ({ status, type }) => status !== "paid" || type !== "expense"
      );

      if (filteredExpenses.length > 0) {
        this.modifiedExpenses = [
          ...filteredExpenses.map(expense => {
            const hasAutoClone =
              expense.status === "paid" && expense.type !== "expense";

            expense.temporary = {};
            expense.temporary.pseudoId = Math.random();
            expense.temporary.action = !hasAutoClone
              ? "nothing"
              : "nothing_autoclone";
            expense.temporary.hasUserPaid = hasAutoClone;

            return expense;
          })
        ];

        // Check if there is only expenses to auto clone.
        // If so, we don't need to open modal.
        const modifiedExpensesHasOnlyAutoClone =
          this.modifiedExpenses.filter(
            ({ temporary }) => temporary.action === "nothing"
          ).length === 0;

        if (!modifiedExpensesHasOnlyAutoClone) {
          this.isModalOpened = true;

          this.onLoadingFinishSpendingDate(false);

          return;
        }
      }

      // Finish current spending date with auto clone
      await expensesService.finishCurrentSpendingDate(
        this.userData.uid,
        this.userData.lookingAtSpendingDate,
        { autoClone: true }
      );

      // Record the user's balance history
      await balancesServices.recordHistory({
        userUid: this.userData.uid,
        spendingDate: this.userData.lookingAtSpendingDate
      });

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

      // Record the user's balance history
      await balancesServices.recordHistory({
        userUid: this.userData.uid,
        spendingDate: this.userData.lookingAtSpendingDate
      });

      this.updateGeneral();
      this.onLoadingFinishSpendingDate(false);
    },
    async updateGeneral() {
      this.closeModal();

      this.$store.dispatch("expenses/setExpenses", {
        userUid: this.userData.uid,
        spendingDate: this.newSpendingDate
      });
      this.$store.dispatch("balances/setBalances", {
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
