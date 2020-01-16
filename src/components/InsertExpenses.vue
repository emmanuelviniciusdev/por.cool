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
            <select v-model="expense.status">
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
import moment from 'moment';

// Services
import userService from "../services/user";
import expenses from "../services/expenses";

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
          indeterminateValidity: false,
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
        indeterminateValidity: false,
      });
    },
    removeExpense(key) {
      this.expenses = this.expenses.filter(expense => expense.key !== key);
      if (this.expenses.length === 0) this.insertExpense();
    },
    generateSpendingDate(lookingAtSpendingDate) {
      return moment(lookingAtSpendingDate).set('date', 1).startOf('day').toDate();
    },
    async saveExpenses() {
      this.onLoading();

      let validationError = false;

      this.expenses = this.expenses.map(expense => {
        const { expenseName, validity, type, indeterminateValidity } = expense;

        // Validation
        if (
          expenseName === "" ||
          (validity === null && type !== "expense" && !indeterminateValidity)
        ) {
          this.$buefy.toast.open({
            message:
              "Tenha certeza de que você definiu um nome e um prazo para todos os seus gastos",
            type: "is-danger",
            position: "is-bottom"
          });

          validationError = true;

          return expense;
        }

        // Make 'validity' null if 'type' is not invoice or savings or 'indeterminateValidity' is true
        expense.validity =
          type === "expense" || indeterminateValidity ? null : validity;
        expense.user = this.userData.uid;
        // TODO: Use 'dateAndTime.startOfMonthAndDay' helper instead of this local function
        expense.spendingDate = this.generateSpendingDate(this.userData.lookingAtSpendingDate);
        expense.created = new Date();
        delete expense.key;

        return expense;
      });

      if (validationError) {
        this.onLoading(false);
        return;
      }

      await expenses.insert(this.expenses);
      this.$router.push({ name: "home" });
      
      this.onLoading(false);
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