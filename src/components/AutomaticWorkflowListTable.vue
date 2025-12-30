<template>
  <div class="AutomaticWorkflowListTable">
    <AutomaticWorkflowImageModal
      :base64Image="selectedImage"
      :isOpen="isImageModalOpen"
      @close="closeImageModal"
    />

    <AutomaticWorkflowDetailsModal
      :isOpen="isDetailsModalOpen"
      :processingMessage="selectedEntry.processingMessage"
      :syncStatus="selectedEntry.syncStatus"
      :extractedContent="selectedEntry.extracted_expense_content_from_image"
      @close="closeDetailsModal"
    />

    <b-table
      :data="entries"
      :loading="loading"
      :mobile-cards="false"
      hoverable
      paginated
      :per-page="10"
    >
      <template slot-scope="props">
        <b-table-column field="image" label="imagem" centered>
          <b-button
            type="is-info"
            size="is-small"
            @click="openImageModal(props.row.base64_image)"
          >
            <b-icon icon="image" size="is-small"></b-icon>
          </b-button>
        </b-table-column>

        <b-table-column
          field="syncStatus"
          label="status de sincronização"
          centered
        >
          <b-tag :type="getSyncStatusType(props.row.syncStatus)">
            {{ getSyncStatusLabel(props.row.syncStatus) }}
          </b-tag>
        </b-table-column>

        <b-table-column field="description" label="descrição">
          {{ props.row.description || "-" }}
        </b-table-column>

        <b-table-column field="created" label="criado em">
          {{ formatDate(props.row.created) }}
        </b-table-column>

        <b-table-column field="details" label="#" centered>
          <b-tooltip
            :label="
              props.row.syncStatus === 'pending'
                ? 'aguardando processamento'
                : 'ver detalhes'
            "
            type="is-dark"
          >
            <b-button
              type="is-warning"
              size="is-small"
              :disabled="props.row.syncStatus === 'pending'"
              @click="openDetailsModal(props.row)"
            >
              <b-icon icon="info-circle" size="is-small"></b-icon>
            </b-button>
          </b-tooltip>
        </b-table-column>
      </template>

      <template slot="empty">
        <section class="section">
          <div class="content has-text-centered has-text-grey">
            <div v-if="!loadingError">
              <p>nenhum registro encontrado para este período</p>
            </div>
            <div class="notification is-danger" v-if="loadingError">
              <p>não foi possível carregar os registros</p>
            </div>
          </div>
        </section>
      </template>
    </b-table>
  </div>
</template>

<script>
import moment from "moment";
import AutomaticWorkflowImageModal from "./AutomaticWorkflowImageModal";
import AutomaticWorkflowDetailsModal from "./AutomaticWorkflowDetailsModal";
import dateAndTimeHelper from "../helpers/dateAndTime";

export default {
  name: "AutomaticWorkflowListTable",
  components: {
    AutomaticWorkflowImageModal,
    AutomaticWorkflowDetailsModal
  },
  props: {
    entries: {
      type: Array,
      required: true
    },
    loading: {
      type: Boolean,
      required: true
    },
    loadingError: {
      type: Boolean,
      required: true
    }
  },
  data() {
    return {
      isImageModalOpen: false,
      selectedImage: "",
      isDetailsModalOpen: false,
      selectedEntry: {}
    };
  },
  methods: {
    getSyncStatusType(status) {
      const types = {
        pending: "is-warning",
        success: "is-success",
        error: "is-danger"
      };
      return types[status] || "is-info";
    },
    getSyncStatusLabel(status) {
      const labels = {
        pending: "pendente",
        success: "sucesso",
        error: "falha"
      };
      return labels[status] || status;
    },
    formatDate(date) {
      if (!date) return "-";
      if (date.seconds) {
        date = dateAndTimeHelper.transformSecondsToDate(date.seconds, {
          dontSetMidnight: true
        });
      }
      return moment(date).format("DD/MM/YYYY HH:mm");
    },
    openImageModal(base64Image) {
      this.selectedImage = base64Image;
      this.isImageModalOpen = true;
    },
    closeImageModal() {
      this.isImageModalOpen = false;
      this.selectedImage = "";
    },
    openDetailsModal(entry) {
      if (entry.syncStatus === "pending") return;
      this.selectedEntry = entry;
      this.isDetailsModalOpen = true;
    },
    closeDetailsModal() {
      this.isDetailsModalOpen = false;
      this.selectedEntry = {};
    }
  }
};
</script>

<style lang="scss" scoped>
.AutomaticWorkflowListTable {
  margin-top: 20px;
}
</style>
