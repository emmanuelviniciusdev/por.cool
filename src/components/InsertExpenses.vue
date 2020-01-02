<template>
  <div>
    <div class="expense" v-for="expense in expenses" :key="expense.key">
      <b-tooltip label="remover" type="is-dark" class="btn-delete-expense">
        <button class="button is-small is-danger" @click="removeExpense(expense.key)">
          <b-icon icon="trash" size="is-small"></b-icon>
        </button>
      </b-tooltip>
      <b-field grouped>
        <b-field
          label="gasto"
          :type="{'is-danger': expense.expenseNameError && expense.expenseName.length <= 0}"
        >
          <b-input
            class="input-expense-name"
            placeholder="conta de..."
            v-model.trim="expense.expenseName"
            maxlength="50"
            :has-counter="false"
          ></b-input>
        </b-field>
        <b-field label="valor">
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
            :type="{'is-danger': expense.validityError && !expenses.validity && !expense.indeterminateValidity}"
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
      <b-tooltip label="salvar" type="is-dark">
        <button class="button is-success" @click="saveExpenses()">
          <b-icon icon="save" size="is-small"></b-icon>
        </button>
      </b-tooltip>
      <b-tooltip label="adicionar" type="is-dark">
        <button class="button is-warning" @click="insertExpense()">
          <b-icon icon="plus" size="is-small"></b-icon>
        </button>
      </b-tooltip>
    </div>
  </div>
</template>

<script>
import { Money } from "v-money";

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
          status: "pending",
          type: "invoice",
          validity: null,
          indeterminateValidity: false
        }
      ]
    };
  },
  methods: {
    insertExpense() {
      this.expenses.push({
        key: Math.random(),
        expenseName: "",
        amount: "",
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
    saveExpenses() {
      let validationError = false;

      this.expenses = this.expenses.map(expense => {
        const { expenseName, validity, type, indeterminateValidity } = expense;

        // Simple validation for 'expenseName' and 'validity'
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

        // Make 'validity' null if 'type' is invoice or savings or 'indeterminateValidity' is true
        expense.validity =
          type === "expense" || indeterminateValidity
            ? null
            : validity;

        return expense;
      });

      if (validationError) return;

      console.log(this.expenses);
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
      0 4px 8px rgba(0, 0, 0, 0.07), 0 8px 16px rgba(0, 0, 0, 0.07),
      0 16px 32px rgba(0, 0, 0, 0.07), 0 32px 64px rgba(0, 0, 0, 0.07);
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
    position: absolute;
    margin-top: -8px !important;
  }
}

.controls {
  margin-top: 10px;

  button {
    margin-right: 5px;
  }
}
</style>