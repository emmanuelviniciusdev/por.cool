<template>
  <div>
    <div class="columns is-desktop">
      <div class="column is-4-desktop">
        <LateralHeader />
      </div>
      <div class="column">
        <h1 class="title has-text-black">
          novo gasto (workflow automático)
          <Help>
            <template v-slot:title>workflow automático de gastos</template>
            <template v-slot:body>
              <div class="content">
                <h1 class="subtitle has-text-weight-normal">como funciona?</h1>
                <p class="is-size-5 has-text-weight-normal">
                  esta funcionalidade permite inserir gastos automaticamente
                  através de um fluxo que utiliza inteligência artificial.
                </p>
                <p class="is-size-5 has-text-weight-normal">
                  basta fazer o upload de um screenshot contendo a descrição dos
                  gastos (por exemplo: notificações de aplicativos de banco).
                </p>
                <p class="is-size-5 has-text-weight-normal">
                  a IA irá processar a imagem e extrair automaticamente as
                  informações dos gastos, como nome da loja, valor, moeda e
                  data.
                </p>
                <h2 class="subtitle has-text-weight-normal is-6">
                  descrições pré-salvas
                </h2>
                <p class="is-size-6 has-text-weight-normal">
                  você pode salvar descrições para reutilizar em futuros
                  uploads. isso ajuda a categorizar e organizar seus gastos
                  automáticos.
                </p>
              </div>
            </template>
          </Help>
        </h1>

        <FilterByDate @onDateChange="loadEntries" />

        <AutomaticWorkflowForm @submitted="loadEntries" />

        <AutomaticWorkflowListTable
          :entries="workflowData.workflowEntries"
          :loading="workflowData.loadingWorkflowEntries"
          :loadingError="workflowData.loadingWorkflowEntriesError"
        />
      </div>
    </div>
  </div>
</template>

<script>
import { mapState } from "vuex";

// Components
import LateralHeader from "../components/LateralHeader";
import Help from "../components/Help";
import FilterByDate from "../components/FilterByDate";
import AutomaticWorkflowForm from "../components/AutomaticWorkflowForm";
import AutomaticWorkflowListTable from "../components/AutomaticWorkflowListTable";

export default {
  name: "AutomaticExpenseWorkflow",
  components: {
    LateralHeader,
    Help,
    FilterByDate,
    AutomaticWorkflowForm,
    AutomaticWorkflowListTable
  },
  computed: {
    ...mapState({
      userData: state => state.user.user,
      workflowData: state => state.expenseAutomaticWorkflow
    })
  },
  methods: {
    loadEntries(spendingDate = null) {
      this.$store.dispatch("expenseAutomaticWorkflow/setWorkflowEntries", {
        userUid: this.userData.uid,
        spendingDate: spendingDate
          ? spendingDate
          : this.userData.lookingAtSpendingDate
      });
    },
    loadPreSavedDescriptions() {
      this.$store.dispatch("expenseAutomaticWorkflow/setPreSavedDescriptions", {
        userUid: this.userData.uid
      });
    }
  },
  created() {
    this.loadEntries();
    this.loadPreSavedDescriptions();
  }
};
</script>

<style lang="scss" scoped>
@media screen and (min-width: 1024px) {
  .columns {
    margin-top: 30px;
  }
}
</style>
