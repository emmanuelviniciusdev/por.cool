<template>
  <div>
    <div class="hero is-warning spendingDateWarning" v-if="showSpendingDateWarning">
      <div class="hero-body">
        <div class="container">
          <h1 class="title">{{ this.userData.displayName | capitalizeName }},</h1>
          <h2 class="subtitle">
            você já pode fechar as contas para
            <b>{{ `${formatedUserLookingAtSpendingDate.month} de ${formatedUserLookingAtSpendingDate.year}` }}</b>.
          </h2>
          <button class="button is-light" @click="finishCurrentSpendingDate()">fechar contas</button>
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

// Filters
import filters from "../filters";

// Services
import expensesService from "../services/expenses";

// Helpers
import dateAndTimeHelper from "../helpers/dateAndTime";

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
    }
  },
  data() {
    return {
      formatedUserLookingAtSpendingDate: null
    };
  },
  computed: {
    ...mapState({
      userData: state => state.user.user
    })
  },
  filters: {
    capitalizeName: filters.capitalizeName
  },
  methods: {
    async finishCurrentSpendingDate() {
      const update = await expensesService.finishCurrentSpendingDate(
        this.userData.uid,
        this.userData.lookingAtSpendingDate
      );
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