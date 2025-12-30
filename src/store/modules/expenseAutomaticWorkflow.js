import expenseAutomaticWorkflowService from "../../services/expenseAutomaticWorkflow";
import preSavedDescriptionService from "../../services/expenseAutomaticWorkflowPreSavedDescription";

export default {
  namespaced: true,

  state: {
    workflowEntries: [],
    loadingWorkflowEntries: false,
    loadingWorkflowEntriesError: false,

    preSavedDescriptions: [],
    loadingPreSavedDescriptions: false
  },

  mutations: {
    SET_WORKFLOW_ENTRIES(state, entries) {
      state.workflowEntries = entries;
    },
    SET_LOADING_WORKFLOW_ENTRIES(state, loading = true) {
      state.loadingWorkflowEntries = loading;
    },
    SET_LOADING_WORKFLOW_ENTRIES_ERROR(state, hasError = true) {
      state.loadingWorkflowEntriesError = hasError;
    },

    SET_PRE_SAVED_DESCRIPTIONS(state, descriptions) {
      state.preSavedDescriptions = descriptions;
    },
    SET_LOADING_PRE_SAVED_DESCRIPTIONS(state, loading = true) {
      state.loadingPreSavedDescriptions = loading;
    }
  },

  actions: {
    async setWorkflowEntries({ commit }, { userUid, spendingDate }) {
      commit("SET_LOADING_WORKFLOW_ENTRIES");
      commit("SET_LOADING_WORKFLOW_ENTRIES_ERROR", false);

      try {
        const entries = await expenseAutomaticWorkflowService.getAll(
          userUid,
          spendingDate
        );
        commit("SET_WORKFLOW_ENTRIES", entries);
      } catch (err) {
        console.error("Erro ao carregar workflow entries:", err);
        commit("SET_LOADING_WORKFLOW_ENTRIES_ERROR");
      } finally {
        commit("SET_LOADING_WORKFLOW_ENTRIES", false);
      }
    },

    async setPreSavedDescriptions({ commit }, { userUid }) {
      commit("SET_LOADING_PRE_SAVED_DESCRIPTIONS");

      try {
        const descriptions = await preSavedDescriptionService.getAll(userUid);
        commit("SET_PRE_SAVED_DESCRIPTIONS", descriptions);
      } finally {
        commit("SET_LOADING_PRE_SAVED_DESCRIPTIONS", false);
      }
    }
  }
};
