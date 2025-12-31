<template>
  <b-modal
    :active="isOpen"
    has-modal-card
    aria-role="dialog"
    aria-modal
    :canCancel="false"
  >
    <div class="modal-card" style="width: auto; max-width: 90vw;">
      <header class="modal-card-head">
        <p class="modal-card-title is-size-6-mobile">detalhes da extração</p>
      </header>
      <section class="modal-card-body">
        <div class="processing-message-section" v-if="processingMessage">
          <b-message :type="messageType">
            {{ processingMessage }}
          </b-message>
        </div>

        <div class="description-section" v-if="description">
          <b-message type="is-info" has-icon>
            <p><strong>Descrição</strong></p>
            <p class="description-text">{{ description }}</p>
          </b-message>
        </div>

        <div
          class="extracted-content-section"
          v-if="parsedExtractedContent && parsedExtractedContent.length > 0"
        >
          <h2 class="subtitle is-5">gastos extraídos</h2>
          <b-table
            :data="parsedExtractedContent"
            :mobile-cards="false"
            hoverable
            paginated
            :per-page="5"
          >
            <template slot-scope="props">
              <b-table-column field="storeName" label="nome">
                {{ props.row.storeName }}
              </b-table-column>

              <b-table-column field="spendingAmount" label="total">
                {{ props.row.spendingAmount | currency }}
              </b-table-column>

              <b-table-column field="spendingCurrency" label="moeda">
                {{ props.row.spendingCurrency }}
              </b-table-column>

              <b-table-column field="spendingDate" label="data do gasto">
                {{ props.row.spendingDate }}
              </b-table-column>

              <b-table-column field="currentDate" label="data vigente">
                {{ props.row.currentDate }}
              </b-table-column>

              <b-table-column field="rawText" label="texto extraído">
                <span class="raw-text-cell">{{ props.row.rawText }}</span>
              </b-table-column>
            </template>

            <template slot="empty">
              <section class="section">
                <div class="content has-text-centered has-text-grey">
                  <p>nenhum gasto extraído</p>
                </div>
              </section>
            </template>
          </b-table>
        </div>

        <div v-else class="has-text-centered has-text-grey">
          <p>nenhum dado de extração disponível</p>
        </div>
      </section>
      <footer class="modal-card-foot">
        <b-button type="is-primary" @click="close">fechar</b-button>
      </footer>
    </div>
  </b-modal>
</template>

<script>
export default {
  name: "AutomaticWorkflowDetailsModal",
  props: {
    isOpen: {
      type: Boolean,
      default: false
    },
    processingMessage: {
      type: String,
      default: ""
    },
    syncStatus: {
      type: String,
      default: ""
    },
    description: {
      type: String,
      default: ""
    },
    extractedContent: {
      type: [Array, String],
      default: () => []
    }
  },
  computed: {
    messageType() {
      if (this.syncStatus === "success") return "is-success";
      if (this.syncStatus === "error") return "is-danger";
      return "is-info";
    },
    parsedExtractedContent() {
      if (!this.extractedContent) return [];
      if (Array.isArray(this.extractedContent)) return this.extractedContent;
      if (typeof this.extractedContent === "string") {
        try {
          return JSON.parse(this.extractedContent);
        } catch {
          return [];
        }
      }
      return [];
    }
  },
  methods: {
    close() {
      this.$emit("close");
    }
  }
};
</script>

<style lang="scss" scoped>
.processing-message-section {
  margin-bottom: 20px;
}

.description-section {
  margin-bottom: 20px;

  .description-text {
    white-space: pre-wrap;
    word-wrap: break-word;
    margin-top: 5px;
  }
}

.extracted-content-section {
  margin-top: 10px;
}

.raw-text-cell {
  display: block;
  max-width: 250px;
  word-wrap: break-word;
  white-space: normal;
  line-height: 1.4;
}
</style>
