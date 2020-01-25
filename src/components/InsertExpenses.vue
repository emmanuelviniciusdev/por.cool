<template>
  <div>
    <div class="expense" v-for="expense in expenses" :key="expense.key">
      <b-tooltip label="remover" type="is-dark" class="btn-delete-expense">
        <button
          class="button is-small is-danger"
          @click="removeExpense(expense.key)"
          :disabled="loading"
        >
          <b-icon icon="trash" size="is-small"></b-icon>
        </button>
      </b-tooltip>
      <b-field grouped>
        <b-field label="gasto">
          <b-input
            class="input-expense-name"
            placeholder="conta de..."
            v-model.trim="expense.expenseName"
            maxlength="50"
            :has-counter="false"
          ></b-input>
        </b-field>
        <b-field label="valor total">
          <money
            v-model="expense.amount"
            v-bind="{
                decimal: ',',
                thousands: '.',
                prefix: 'R$',
                precision: 2
            }"
            class="input input-amount"
          ></money>
        </b-field>
        <b-field label="valor que já foi pago" v-if="expense.status === 'partially_paid'">
          <money
            :ref="'alreadyPaidAmountInput_' + expense.key"
            v-model="expense.alreadyPaidAmount"
            v-bind="{
                decimal: ',',
                thousands: '.',
                prefix: 'R$',
                precision: 2
            }"
            class="input input-amount"
          ></money>
        </b-field>
        <b-field label="status">
          <div class="select">
            <select
              v-model="expense.status"
              @change="onAlreadyPaidAmountFocus($event.target.value, expense.key)"
            >
              <option value="pending">pendente</option>
              <option value="partially_paid">parcialmente pago</option>
              <option value="paid">pago</option>
            </select>
          </div>
        </b-field>
        <b-field label="tipo">
          <div class="select">
            <select v-model="expense.type">
              <option value="expense">gasto</option>
              <option value="invoice">fatura</option>
              <option value="savings">poupança</option>
            </select>
          </div>
        </b-field>

        <div
          v-if="expense.type === 'invoice' || expense.type === 'savings'"
          style="display: inline-block !important;"
        >
          <b-field
            :label="expense.type === 'invoice' ? 'esta fatura vai até' : 'esta poupança vai até'"
          >
            <b-datepicker
              type="month"
              placeholder="Selecione uma data"
              class="input-date"
              :disabled="expense.indeterminateValidity"
              v-model="expense.validity"
            ></b-datepicker>
          </b-field>
          <b-field>
            <b-checkbox
              class="input-checkbox"
              v-model="expense.indeterminateValidity"
            >tempo indeterminado</b-checkbox>
          </b-field>
        </div>
      </b-field>
    </div>

    <div class="controls">
      <b-tooltip :label="loading ? 'aguarde' : 'salvar'" type="is-dark">
        <b-button type="is-success" :loading="loading" @click="saveExpenses()">
          <b-icon icon="save" size="is-small"></b-icon>
        </b-button>
      </b-tooltip>
      <b-tooltip label="adicionar" type="is-dark">
        <button class="button is-warning" @click="insertExpense()" :disabled="loading">
          <b-icon icon="plus" size="is-small"></b-icon>
        </button>
      </b-tooltip>
    </div>
  </div>
</template>

<script>
import { Money } from "v-money";
import { mapState } from "vuex";
import firebase from "firebase/app";
import "firebase/auth";
import moment from "moment";

// Services
import userService from "../services/user";
import balancesService from "../services/balances";
import expenses from "../services/expenses";

// Helpers
import dateAndTime from "../helpers/dateAndTime";

export default {
  name: "InsertExpenses",
  components: {
    Money
  },
  data() {
    return {
      expenses: [
        {
          key: Math.random(),
          expenseName: "",
          amount: 0,
          alreadyPaidAmount: 0,
          status: "pending",
          type: "expense",
          validity: null,
          indeterminateValidity: false
        }
      ],
      loading: false
    };
  },
  computed: {
    ...mapState({
      userData: state => state.user.user
    })
  },
  methods: {
    onLoading(state = true) {
      this.loading = state;
    },
    insertExpense() {
      this.expenses.push({
        key: Math.random(),
        expenseName: "",
        amount: 0,
        alreadyPaidAmount: 0,
        status: "pending",
        type: "expense",
        validity: null,
        indeterminateValidity: false
      });
    },
    removeExpense(key) {
      this.expenses = this.expenses.filter(expense => expense.key !== key);
      if (this.expenses.length === 0) this.insertExpense();
    },
    onAlreadyPaidAmountFocus(value, expenseKey) {
      this.$nextTick(() => {
        if (value === "partially_paid") {
          this.$refs[`alreadyPaidAmountInput_${expenseKey}`][0].$el.focus();
        } else {
          this.expense = this.expenses.map(expense => {
            if (expense.key === expenseKey) expense.alreadyPaidAmount = 0;
            return expense;
          });
        }
      });
    },
    async saveExpenses() {
      this.onLoading();

      const errorData = { error: false, msg: "" };

      const expensesToInsert = [];

      this.expenses.forEach(expense => {
        // Add more data
        expense.validity =
          expense.type === "expense" || expense.indeterminateValidity
            ? null
            : expense.validity;
        expense.user = this.userData.uid;
        expense.spendingDate = dateAndTime.startOfMonthAndDay(
          this.userData.lookingAtSpendingDate
        );
        expense.created = new Date();

        const {
          expenseName,
          validity,
          type,
          indeterminateValidity,
          amount,
          alreadyPaidAmount
        } = expense;

        // Validation
        if (
          expenseName === "" ||
          (validity === null && type !== "expense" && !indeterminateValidity)
        ) {
          errorData.error = true;
          errorData.msg =
            "tenha certeza de que preencheu todos os campos corretamente";
        }

        if (alreadyPaidAmount > amount) {
          errorData.error = true;
          errorData.msg = "você não pode pagar mais do que deve";
        }

        if (
          validity !== null &&
          moment(validity).isBefore(moment(this.userData.lookingAtSpendingDate))
        ) {
          errorData.error = true;
          errorData.msg = `não foi possível adicionar uma fatura/poupança no passado...`;
        }

        // Delete some data
        delete expense.key;

        expensesToInsert.push(expense);
      });

      if (errorData.error) {
        this.$buefy.toast.open({
          message: errorData.msg,
          position: "is-bottom",
          type: "is-danger"
        });
        this.onLoading(false);
        return;
      }

      await expenses.insert(expensesToInsert);

      this.refreshRemainingBalance();

      this.$router.push({ name: "home" });
      this.onLoading(false);
    },
    async refreshRemainingBalance() {
      const remainingBalance = await balancesService.calculate({
        userUid: this.userData.uid,
        spendingDate: this.userData.lookingAtSpendingDate,
      });
      this.$store.dispatch('balances/setCurrentBalance', remainingBalance);
    }
  }
};
</script>

<style lang="scss" scoped>
.expense {
  padding: 10px;
  padding-bottom: 15px;
  border: solid #ccc 1px;
  border-radius: 5px;
  transition: all 0.2s;
  margin-bottom: 10px;

  &:hover {
    transition: all 0.2s;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.07), 0 2px 4px rgba(0, 0, 0, 0.07),
      0 4px 8px rgba(0, 0, 0, 0.07);
  }

  .btn-delete-expense {
    float: right;
  }

  .field {
    display: inline-block;
  }

  .input-expense-name {
    width: 230px !important;
  }
  .input-amount {
    width: 120px !important;
  }
  .input-checkbox {
    position: relative;
    top: 13px;
    margin-bottom: 13px !important;
  }
}

.controls {
  margin-top: 10px;

  button {
    margin-right: 5px;
  }
}
</style>