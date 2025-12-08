<template>
  <div>
    <div class="columns is-desktop">
      <div class="column is-4-desktop">
        <LateralHeader />
      </div>
      <div class="column">
        <h1 class="title has-text-black">
          bancos e instituições
          <Help>
            <template v-slot:title>bancos e instituições</template>
            <template v-slot:body>
              <div class="content">
                <h1 class="subtitle has-text-weight-normal">O que eu posso fazer por aqui?</h1>
                <p
                  class="is-size-5 has-text-weight-normal"
                >Esta tela é para você gerenciar os seus bancos e instituições financeiras. Por aqui, você poderá cadastrar os bancos, corretoras e outras instituições que você utiliza, além de indicar quais tipos de operações você faz em cada uma delas (cartão de crédito, movimentação de dinheiro, investimentos).</p>
                <p class="is-size-5 has-text-weight-normal">Você pode adicionar, editar e remover bancos e instituições a qualquer momento.</p>
              </div>
            </template>
          </Help>
        </h1>

        <AddBank />

        <BankListTable
          :banks="banks.banksList"
          :loadingBanksList="banks.loadingBanksList"
          :loadingBanksListError="banks.loadingBanksListError"
          @edit-bank="onEditBank"
        />

        <EditBank :bank="bankToEdit" />
      </div>
    </div>
  </div>
</template>

<script>
import { mapState } from "vuex";

// Components
import LateralHeader from "../components/LateralHeader";
import Help from "../components/Help";
import AddBank from "../components/AddBank";
import BankListTable from "../components/BankListTable";
import EditBank from "../components/EditBank";

export default {
  name: "Banks",
  components: {
    LateralHeader,
    Help,
    AddBank,
    BankListTable,
    EditBank
  },
  data() {
    return {
      bankToEdit: {}
    };
  },
  computed: {
    ...mapState({
      userData: state => state.user.user,
      banks: state => state.banks
    })
  },
  methods: {
    loadBanks() {
      this.$store.dispatch("banks/setBanksList", {
        userUid: this.userData.uid
      });
    },
    onEditBank(bank) {
      this.bankToEdit = { ...bank };
    }
  },
  created() {
    this.loadBanks();
  }
};
</script>

<style lang="scss" scoped>
.AddBank {
  margin-bottom: 15px;
}

@media screen and (min-width: 1024px) {
  .columns {
    margin-top: 30px;
  }
}
</style>
