<template>
  <div>
    <b-table
      :data="banks"
      :loading="loadingBanksList"
      :mobile-cards="false"
      hoverable
      paginated
      :per-page="10"
    >
      <template slot-scope="props">
        <b-table-column field="nome" label="nome">{{
          props.row.nome
        }}</b-table-column>
        <b-table-column
          field="cartaoCredito"
          label="cartão de crédito"
          centered
        >
          <b-icon
            :icon="props.row.cartaoCredito ? 'check' : 'times'"
            :type="props.row.cartaoCredito ? 'is-success' : 'is-danger'"
          ></b-icon>
        </b-table-column>
        <b-table-column
          field="movimentacaoDinheiro"
          label="movimentação de dinheiro"
          centered
        >
          <b-icon
            :icon="props.row.movimentacaoDinheiro ? 'check' : 'times'"
            :type="props.row.movimentacaoDinheiro ? 'is-success' : 'is-danger'"
          ></b-icon>
        </b-table-column>
        <b-table-column field="investimentos" label="investimentos" centered>
          <b-icon
            :icon="props.row.investimentos ? 'check' : 'times'"
            :type="props.row.investimentos ? 'is-success' : 'is-danger'"
          ></b-icon>
        </b-table-column>
        <b-table-column field="observacoes" label="observações">
          <span class="observacoes-text">{{
            truncateText(props.row.observacoes, 15)
          }}</span>
        </b-table-column>
        <b-table-column field="action" label="#">
          <b-tooltip label="editar" type="is-dark" class="action-button">
            <button
              class="button is-warning is-small"
              @click="editBank(props.row)"
            >
              <b-icon icon="edit"></b-icon>
            </button>
          </b-tooltip>
          <b-tooltip label="remover" type="is-dark" class="action-button">
            <button
              class="button is-danger is-small"
              @click="deleteBank(props.row)"
            >
              <b-icon icon="trash"></b-icon>
            </button>
          </b-tooltip>
        </b-table-column>
      </template>

      <template slot="empty">
        <section class="section">
          <div class="content has-text-black has-text-centered">
            <div v-if="!loadingBanksListError">
              <div class="notification">
                Nenhum banco ou instituição foi encontrado por aqui...
              </div>
            </div>

            <div class="notification is-danger" v-if="loadingBanksListError">
              <p>Não foi possível carregar os seus bancos e instituições</p>
            </div>
          </div>
        </section>
      </template>
    </b-table>
  </div>
</template>

<script>
import { mapState } from "vuex";

// Services
import banksService from "../services/banks";

export default {
  name: "BankListTable",
  props: {
    banks: {
      required: true,
      type: Array
    },
    loadingBanksList: {
      required: true,
      type: Boolean
    },
    loadingBanksListError: {
      required: true,
      type: Boolean
    }
  },
  computed: {
    ...mapState({
      userData: state => state.user.user
    })
  },
  methods: {
    truncateText(text, maxLength) {
      if (!text) return "-";
      if (text.length <= maxLength) return text;
      return text.substring(0, maxLength) + "...";
    },
    editBank(bank) {
      this.$emit("edit-bank", bank);
    },
    deleteBank({ id, nome }) {
      this.$buefy.dialog.confirm({
        title: "deletar banco ou instituição",
        message: `Você está prestes a deletar <b>${nome}</b> PARA SEMPRE. Você tem certeza de que deseja continuar?`,
        confirmText: "Sim. Deletar.",
        type: "is-danger",
        hasIcon: true,
        onConfirm: async () => {
          try {
            await banksService.removeBank(id);

            this.$store.dispatch("banks/setBanksList", {
              userUid: this.userData.uid
            });

            this.$buefy.toast.open({
              message: "banco deletado com sucesso",
              type: "is-success",
              position: "is-bottom"
            });
          } catch {
            this.$buefy.toast.open({
              message: "ocorreu um erro ao tentar deletar o banco",
              type: "is-danger",
              position: "is-bottom"
            });
          }
        }
      });
    }
  }
};
</script>

<style scoped>
.action-button {
  margin-right: 0.5rem;
}

.action-button:last-child {
  margin-right: 0;
}
</style>
