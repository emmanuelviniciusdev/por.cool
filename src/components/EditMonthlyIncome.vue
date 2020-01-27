<template>
  <div class="EditMonthlyIncome">
    <b-modal :active="openModal" :canCancel="false">
      <div class="modal-card" style="width: auto">
        <div class="modal-card-head">
          <p class="modal-card-title">
            sua renda atual:
            <b>{{userData.monthlyIncome | currency}}</b>
          </p>
        </div>
        <section class="modal-card-body">
          <b-field label="insira a nova renda">
            <money
              class="input"
              style="width: 300px;"
              v-model.trim="monthlyIncome"
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
          <b-button type="is-default" @click="onOpenModal(false)">cancelar</b-button>
          <b-button type="is-primary" @click="updateMonthlyIncome()" :loading="loading">alterar</b-button>
        </footer>
      </div>
    </b-modal>

    <div class="notification is-primary">
      <p class="is-size-5">sua renda fixa atual é de <b>{{ this.userData.monthlyIncome | currency }}</b></p>
      <b-button @click="onOpenModal()" style="margin-top: 10px;">clique para alterar</b-button>
    </div>
  </div>
</template>

<script>
import { mapState } from "vuex";
import { Money } from "v-money";

// Services
import userService from "../services/user";

export default {
  name: "EditMonthlyIncome",
  components: {
    Money
  },
  data() {
    return {
      openModal: false,
      loading: false
    };
  },
  computed: {
    ...mapState({
      userData: state => state.user.user
    })
  },
  methods: {
    onOpenModal(state = true) {
      this.openModal = state;
    },
    onLoading(state = true) {
      this.loading = state;
    },
    async updateMonthlyIncome() {
      this.onLoading();

      try {
        await userService.update(this.userData.uid, {
          monthlyIncome: this.monthlyIncome
        });

        this.$store.dispatch("user/update", {
          monthlyIncome: this.monthlyIncome
        });
        this.$store.dispatch("balances/setBalances", {
          userUid: this.userData.uid,
          spendingDate: this.userData.lookingAtSpendingDate,
        });

        this.$buefy.toast.open({
          message: "renda alterada com sucesso",
          type: "is-success",
          position: "is-bottom"
        });

        this.onOpenModal(false);
      } catch (err) {
        this.$buefy.toast.open({
          message: "ocorreu um erro ao tentar alterar a sua renda fixa mensal",
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