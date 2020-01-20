<template>
  <div>
    <b-modal
      :active="openModal"
      :custom-class="{'modify-modal': expenseToEdit.type !== 'expense'}"
      :canCancel="false"
    >
      <div class="modal-card" style="width: auto;">
        <div class="modal-card-head">
          <p class="modal-card-title">
            editando
            <b>{{expenseToEdit.expenseName ? expenseToEdit.expenseName : 'insira o nome do gasto...'}}</b>
          </p>
        </div>
        <section class="modal-card-body">
          <b-field grouped>
            <b-field label="gasto" :type="{'is-danger': hasInputErrorAndDirty('expenseName')}">
              <b-input
                v-model.trim="$v.expenseToEdit.expenseName.$model"
                class="input-expense-name"
                placeholder="conta de..."
                maxlength="50"
                :has-counter="false"
              ></b-input>
            </b-field>
            <b-field label="valor total">
              <money
                v-model.trim="expenseToEdit.amount"
                v-bind="{
                decimal: ',',
                thousands: '.',
                prefix: 'R$',
                precision: 2
            }"
                class="input input-amount"
              ></money>
            </b-field>
            <b-field label="valor diferencial" v-if="expenseToEdit.differenceAmount !== undefined">
              <money
                v-model.trim="expenseToEdit.differenceAmount"
                v-bind="{
                decimal: ',',
                thousands: '.',
                prefix: 'R$',
                precision: 2
            }"
                class="input input-amount"
              ></money>
            </b-field>
            <b-field label="valor que já foi pago" v-if="expenseToEdit.status === 'partially_paid'">
              <money
                ref="alreadyPaidAmountInput"
                v-model.trim="expenseToEdit.alreadyPaidAmount"
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
                  v-model.trim="expenseToEdit.status"
                  @change="onAlreadyPaidAmountFocus($event.target.value)"
                >
                  <option value="pending">pendente</option>
                  <option value="partially_paid">parcialmente pago</option>
                  <option value="paid">pago</option>
                </select>
              </div>
            </b-field>

            <b-field label="tipo">
              <div class="select">
                <select v-model.trim="expenseToEdit.type">
                  <option value="expense">gasto</option>
                  <option value="invoice">fatura</option>
                  <option value="savings">poupança</option>
                </select>
              </div>
            </b-field>

            <div
              v-if="expenseToEdit.type === 'invoice' || expenseToEdit.type === 'savings'"
              style="display: inline-block !important;"
            >
              <b-field
                :label="expenseToEdit.type === 'invoice' ? 'esta fatura vai até' : 'esta poupança vai até'"
              >
                <b-datepicker
                  editable
                  type="month"
                  placeholder="Selecione uma data"
                  class="input-date"
                  :disabled="expenseToEdit.indeterminateValidity"
                  v-model.trim="expenseToEdit.validity"
                ></b-datepicker>
              </b-field>
              <b-field>
                <b-checkbox
                  class="input-checkbox"
                  v-model.trim="expenseToEdit.indeterminateValidity"
                >tempo indeterminado</b-checkbox>
              </b-field>
            </div>
          </b-field>
        </section>
        <footer class="modal-card-foot">
          <b-button @click="openModal = false">cancelar</b-button>
          <b-button type="is-warning" @click="save()" :loading="isLoading">salvar</b-button>
        </footer>
      </div>
    </b-modal>
  </div>
</template>

<script>
import { Money } from "v-money";
import { required } from "vuelidate/lib/validators";
import { mapState } from "vuex";
import moment from "moment";

// Services
import expenseService from "../services/expenses";

// Helpers
import dateAndTime from "../helpers/dateAndTime";

export default {
  name: "EditExpense",
  props: {
    expense: {
      required: true,
      type: Object
    },
    openModal: {
      type: Boolean
    }
  },
  components: {
    Money
  },
  data() {
    return {
      expenseToEdit: {},
      isLoading: false,
      openModal: false
    };
  },
  computed: {
    ...mapState({
      userData: state => state.user.user
    })
  },
  watch: {
    expense() {
      this.expenseToEdit = { ...this.expense };
      const { validity } = this.expenseToEdit;

      delete this.expenseToEdit.watchKey;

      this.expenseToEdit.validity =
        validity === null
          ? validity
          : dateAndTime.transformSecondsToDate(validity.seconds);

      this.openModal = true;
    }
  },
  validations: {
    expenseToEdit: {
      expenseName: { required }
    }
  },
  methods: {
    async save() {
      if (!this.expenseToEdit || this.$v.expenseToEdit.$invalid) return;

      this.onLoading();

      const {
        indeterminateValidity,
        amount,
        alreadyPaidAmount,
        differenceAmount,
        validity
      } = this.expenseToEdit;

      // Validations
      const errorData = { error: false, msg: "" };

      if (indeterminateValidity) this.expenseToEdit.validity = null;

      const owedAmount =
        differenceAmount !== undefined ? amount + differenceAmount : amount;

      if (alreadyPaidAmount > owedAmount) {
        errorData.error = true;
        errorData.msg = `você não pode pagar um valor maior do que aquele que você deve (R$ ${owedAmount}) 😞`;
      }

      if (validity !== null && moment(validity).isBefore(moment(this.userData.lookingAtSpendingDate))) {
        errorData.error = true;
        errorData.msg = `não foi possível adicionar uma fatura/poupança no passado...`;
      }

      if (errorData.error) {
        this.$buefy.toast.open({
          message: errorData.msg,
          position: "is-bottom",
          type: "is-danger"
        });

        this.onLoading(false);

        return;
      }

      // Update
      await expenseService.update(this.expenseToEdit);

      this.$store.dispatch("expenses/setExpenses", {
        userUid: this.userData.uid,
        spendingDate: this.userData.lookingAtSpendingDate
      });

      this.$buefy.toast.open({
        message: "atualizado com sucesso",
        position: "is-bottom",
        type: "is-success"
      });

      this.openModal = false;
      this.onLoading(false);
    },
    onLoading(state = true) {
      this.isLoading = state;
    },
    hasInputErrorAndDirty(input) {
      return (
        this.$v.expenseToEdit[input].$error &&
        this.$v.expenseToEdit[input].$dirty
      );
    },
    onAlreadyPaidAmountFocus(value) {
      this.$nextTick(() => {
        if (value === "partially_paid")
          this.$refs.alreadyPaidAmountInput.$el.focus();
        else this.expenseToEdit.alreadyPaidAmount = 0;
      });
    }
  }
};
</script>

<style lang="scss" scoped>
.field {
  display: inline-block;
  margin-right: 5px;
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

@media screen and (min-width: 1024px) {
  .modify-modal .modal-card-body {
    height: 300px !important;
  }
}
</style>