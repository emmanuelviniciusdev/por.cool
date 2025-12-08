<template>
  <div class="AddBank">
    <b-modal :active="openModal" :canCancel="false">
      <div class="modal-card" style="width: auto">
        <div class="modal-card-head">
          <p class="modal-card-title">adicionar banco ou instituição</p>
        </div>
        <section class="modal-card-body">
          <b-field
            label="nome"
            :type="{'is-danger': hasInputErrorAndDirty('nome')}"
            :message="{
              'insira o nome do banco ou instituição': isInvalidInputMsg('nome', 'required'),
              'o nome é muito grande': isInvalidInputMsg('nome', 'maxLength'),
            }"
          >
            <b-input
              style="width: 400px;"
              v-model.trim="$v.form.nome.$model"
              placeholder="Nome do banco ou da instituição financeira"
              maxlength="100"
              :has-counter="true"
            ></b-input>
          </b-field>

          <div class="field">
            <label class="label">tipos de operações</label>
            <div class="control checkboxes-control">
              <div class="checkbox-wrapper">
                <b-checkbox v-model="form.cartaoCredito">cartão de crédito</b-checkbox>
              </div>
              <div class="checkbox-wrapper">
                <b-checkbox v-model="form.movimentacaoDinheiro">movimentação de dinheiro</b-checkbox>
              </div>
              <div class="checkbox-wrapper">
                <b-checkbox v-model="form.investimentos">investimentos</b-checkbox>
              </div>
            </div>
          </div>

          <b-field
            label="observações (opcional)"
            :type="{'is-danger': hasInputErrorAndDirty('observacoes')}"
            :message="{
              'as observações são muito grandes': isInvalidInputMsg('observacoes', 'maxLength'),
            }"
          >
            <b-input
              style="width: 400px;"
              v-model.trim="$v.form.observacoes.$model"
              placeholder="Observações gerais..."
              maxlength="200"
              type="textarea"
              :has-counter="true"
            ></b-input>
          </b-field>
        </section>
        <footer class="modal-card-foot">
          <b-button type="is-default" @click="onOpenModal(false)">cancelar</b-button>
          <b-button type="is-primary" @click="addBank()" :loading="loading">adicionar</b-button>
        </footer>
      </div>
    </b-modal>

    <b-button type="is-success" icon-left="plus" @click="onOpenModal()">adicionar</b-button>
  </div>
</template>

<script>
import { mapState } from "vuex";
import { required, maxLength } from "vuelidate/lib/validators";

// Services
import banksService from "../services/banks";

export default {
  name: "AddBank",
  data() {
    return {
      openModal: false,
      loading: false,

      form: {
        nome: "",
        cartaoCredito: false,
        movimentacaoDinheiro: false,
        investimentos: false,
        observacoes: ""
      }
    };
  },
  validations: {
    form: {
      nome: { required, maxLength: maxLength(100) },
      observacoes: { maxLength: maxLength(200) }
    }
  },
  computed: {
    ...mapState({
      userData: state => state.user.user
    })
  },
  methods: {
    onOpenModal(state = true) {
      this.openModal = state;
      if (!state) {
        this.$v.$reset();
      }
    },
    onLoading(state = true) {
      this.loading = state;
    },
    async addBank() {
      this.$v.form.$touch();
      if (this.$v.form.$invalid) return;

      this.onLoading();

      try {
        await banksService.addBank({
          nome: this.form.nome,
          cartaoCredito: this.form.cartaoCredito,
          movimentacaoDinheiro: this.form.movimentacaoDinheiro,
          investimentos: this.form.investimentos,
          observacoes: this.form.observacoes,
          userUid: this.userData.uid
        });

        this.$store.dispatch("banks/setBanksList", {
          userUid: this.userData.uid
        });

        this.$buefy.toast.open({
          message: "banco adicionado com sucesso",
          type: "is-success",
          position: "is-bottom"
        });

        this.form.nome = "";
        this.form.cartaoCredito = false;
        this.form.movimentacaoDinheiro = false;
        this.form.investimentos = false;
        this.form.observacoes = "";

        this.onOpenModal(false);
      } catch (err) {
        console.error('Error in addBank:', err);
        this.$buefy.toast.open({
          message: `Erro: ${err.message || 'ocorreu um erro ao tentar adicionar banco ou instituição'}`,
          type: "is-danger",
          position: "is-bottom",
          duration: 5000
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

<style scoped>
.field {
  margin-bottom: 0.75rem;
}

.checkboxes-control .checkbox-wrapper {
  display: inline-block;
  margin-right: 1rem;
  margin-bottom: 0.5rem;
}

.checkboxes-control .checkbox-wrapper:last-child {
  margin-right: 0;
}
</style>
