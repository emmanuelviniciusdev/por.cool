<template>
  <div>
    <b-modal :active="showModal" :canCancel="false">
      <div class="modal-card" style="width: auto;">
        <div class="modal-card-head">
          <p class="modal-card-title">
            editando
            <b>{{bankToEdit.nome ? bankToEdit.nome : 'insira o nome...'}}</b>
          </p>
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
              v-model.trim="$v.bankToEdit.nome.$model"
              placeholder="ex: Banco do Brasil, Nubank, XP Investimentos"
              maxlength="100"
              :has-counter="true"
            ></b-input>
          </b-field>

          <div class="field">
            <label class="label">tipos de operações</label>
            <div class="control checkboxes-control">
              <div class="checkbox-wrapper">
                <b-checkbox v-model="bankToEdit.cartaoCredito">cartão de crédito</b-checkbox>
              </div>
              <div class="checkbox-wrapper">
                <b-checkbox v-model="bankToEdit.movimentacaoDinheiro">movimentação de dinheiro</b-checkbox>
              </div>
              <div class="checkbox-wrapper">
                <b-checkbox v-model="bankToEdit.investimentos">investimentos</b-checkbox>
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
              v-model.trim="$v.bankToEdit.observacoes.$model"
              placeholder="ex: conta corrente nº 12345-6"
              maxlength="200"
              type="textarea"
              :has-counter="true"
            ></b-input>
          </b-field>
        </section>
        <footer class="modal-card-foot">
          <b-button @click="closeModal()">cancelar</b-button>
          <b-button type="is-warning" @click="save()" :loading="isLoading">salvar</b-button>
        </footer>
      </div>
    </b-modal>
  </div>
</template>

<script>
import { required, maxLength } from "vuelidate/lib/validators";
import { mapState } from "vuex";

// Services
import banksService from "../services/banks";

export default {
  name: "EditBank",
  props: {
    bank: {
      required: true,
      type: Object
    }
  },
  data() {
    return {
      bankToEdit: {},
      isLoading: false,
      showModal: false
    };
  },
  computed: {
    ...mapState({
      userData: state => state.user.user
    })
  },
  watch: {
    bank() {
      this.bankToEdit = { ...this.bank };
      this.showModal = true;
      this.$nextTick(() => {
        if (this.$v) {
          this.$v.$reset();
        }
      });
    }
  },
  validations: {
    bankToEdit: {
      nome: { required, maxLength: maxLength(100) },
      observacoes: { maxLength: maxLength(200) }
    }
  },
  methods: {
    closeModal() {
      this.showModal = false;
      this.$v.$reset();
    },
    async save() {
      this.$v.bankToEdit.$touch();
      if (!this.bankToEdit || this.$v.bankToEdit.$invalid) return;

      this.onLoading();

      try {
        await banksService.updateBank(this.bankToEdit);

        this.$store.dispatch("banks/setBanksList", {
          userUid: this.userData.uid
        });

        this.$buefy.toast.open({
          message: "atualizado com sucesso",
          position: "is-bottom",
          type: "is-success"
        });

        this.closeModal();
      } catch (err) {
        this.$buefy.toast.open({
          message: "ocorreu um erro ao tentar atualizar",
          position: "is-bottom",
          type: "is-danger"
        });
      } finally {
        this.onLoading(false);
      }
    },
    onLoading(state = true) {
      this.isLoading = state;
    },
    hasInputErrorAndDirty(input) {
      return (
        this.$v.bankToEdit[input].$error &&
        this.$v.bankToEdit[input].$dirty
      );
    },
    isInvalidInputMsg(input, role) {
      return !this.$v.bankToEdit[input][role] && this.$v.bankToEdit[input].$error;
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
