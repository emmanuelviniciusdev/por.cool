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
          <b-field label="valor do saldo">
            <money
              class="input"
              style="width: 300px;"
              v-model.trim="balance"
              v-bind="{
                    decimal: ',',
                    thousands: '.',
                    prefix: 'R$',
                    precision: 2
                }"
            ></money>
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
      balance: 0
    };
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
      this.onLoading();

      try {
        await balancesService.addAdditionalBalance({
          balance: this.balance,
          spendingDate: this.userData.lookingAtSpendingDate,
          userUid: this.userData.uid
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

        this.balance = 0;

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
    }
  }
};
</script>

<style>
</style>