<template>
  <div class="AddBalance">
    <b-modal :active="openModal" :canCancel="false">
      <div class="modal-card" style="width: auto">
        <div class="modal-card-head">
          <p class="modal-card-title">
            adionando saldo para
            <b>{{formatedUserLookingAtSpendingDate.month + ' de ' + formatedUserLookingAtSpendingDate.year}}</b>
          </p>
        </div>
        <section class="modal-card-body">
          <b-field
            label="valor do saldo"
            :type="{'is-danger': hasInputErrorAndDirty('balance')}"
            :message="{
              'insira o valor do saldo': isInvalidInputMsg('balance', 'required'),
            }"
          >
            <money
              class="input"
              style="width: 200px;"
              v-model.trim="$v.form.balance.$model"
              v-bind="{
                    decimal: ',',
                    thousands: '.',
                    prefix: 'R$',
                    precision: 2
                }"
            ></money>
          </b-field>
          <b-field
            label="descrição (opcional)"
            :type="{'is-danger': hasInputErrorAndDirty('description')}"
            :message="{
              'a descrição é muito grande': isInvalidInputMsg('description', 'maxLength'),
            }"
          >
            <b-input
              style="width: 300px;"
              v-model.trim="$v.form.description.$model"
              placeholder="descrição"
              maxlength="20"
              :has-counter="true"
            ></b-input>
          </b-field>
        </section>
        <footer class="modal-card-foot">
          <b-button type="is-default" @click="openModal = false">cancelar</b-button>
          <b-button type="is-primary" @click="addBalance()" :loading="loading">adicionar</b-button>
        </footer>
      </div>
    </b-modal>

    <b-button type="is-success" icon-left="dollar-sign" @click="onOpenModal()">novo saldo</b-button>
  </div>
</template>

<script>
import { mapState } from "vuex";
import { Money } from "v-money";
import { required, maxLength } from "vuelidate/lib/validators";

// Services
import balancesService from "../services/balances";

// Helpers
import dateAndTimeHelper from "../helpers/dateAndTime";

export default {
  name: "AddBalance",
  components: {
    Money
  },
  data() {
    return {
      openModal: false,
      loading: false,

      form: {
        balance: 0,
        description: ""
      }
    };
  },
  validations: {
    form: {
      balance: { required },
      description: { maxLength: maxLength(20) }
    }
  },
  computed: {
    ...mapState({
      userData: state => state.user.user
    }),
    formatedUserLookingAtSpendingDate() {
      return dateAndTimeHelper.extractOnly(
        this.userData.lookingAtSpendingDate,
        ["year", "month"]
      );
    }
  },
  methods: {
    onOpenModal(state = true) {
      this.openModal = state;
    },
    onLoading(state = true) {
      this.loading = state;
    },
    async addBalance() {
      if (this.$v.form.$invalid) return;

      this.onLoading();

      try {
        await balancesService.addAdditionalBalance({
          balance: this.form.balance,
          description: this.form.description,
          spendingDate: this.userData.lookingAtSpendingDate,
          userUid: this.userData.uid
        });

        this.$store.dispatch("balances/setBalances", {
          userUid: this.userData.uid,
          spendingDate: this.userData.lookingAtSpendingDate
        });
        this.$store.dispatch("balances/setBalancesList", {
          userUid: this.userData.uid,
          spendingDate: this.userData.lookingAtSpendingDate
        });

        this.$buefy.toast.open({
          message: "saldo adicionado com sucesso",
          type: "is-success",
          position: "is-bottom"
        });

        this.form.balance = 0;
        this.form.description = "";

        this.onOpenModal(false);
      } catch (err) {
        this.$buefy.toast.open({
          message: "ocorreu um erro ao tentar adicionar saldo adicional",
          type: "is-danger",
          position: "is-bottom"
        });
      } finally {
        this.onLoading(false);
      }
    },
    hasInputErrorAndDirty(input) {
      return this.$v.form[input].$error && this.$v.form[input].$dirty;
    },
    isInvalidInputMsg(input, role) {
      return !this.$v.form[input][role] && this.$v.form[input].$error;
    }
  }
};
</script>

<style>
</style>