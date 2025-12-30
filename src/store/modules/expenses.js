import expensesService from "../../services/expenses";

export default {
  namespaced: true,

  state: {
    expenses: [],
    loadingExpenses: false,
    loadingExpensesError: false,

    spendingDatesList: [],
    loadingSpendingDatesList: false,
    loadingSpendingDatesListError: false
  },

  mutations: {
    SET_EXPENSES(state, expenses) {
      state.expenses = expenses;
    },
    SET_LOADING_EXPENSES(state, loading = true) {
      state.loadingExpenses = loading;
    },
    SET_LOADING_EXPENSES_ERROR(state, hasError = true) {
      state.loadingExpensesError = hasError;
    },

    SET_SPENDING_DATES_LIST(state, spendingDatesList) {
      state.spendingDatesList = spendingDatesList;
    }
  },

  actions: {
    async setExpenses({ commit }, { userUid, spendingDate }) {
      commit("SET_LOADING_EXPENSES");
      commit("SET_LOADING_EXPENSES_ERROR", false);

      try {
        const userExpenses = await expensesService.getAll(
          userUid,
          spendingDate
        );
        commit("SET_EXPENSES", userExpenses);
      } catch {
        commit("SET_LOADING_EXPENSES_ERROR");
      } finally {
        commit("SET_LOADING_EXPENSES", false);
      }
    },

    async setSpendingDatesList({ commit }, { userUid, lookingAtSpendingDate }) {
      const spendingDatesList = await expensesService.getSpendingDatesList({
        userUid,
        lookingAtSpendingDate
      });
      commit("SET_SPENDING_DATES_LIST", spendingDatesList);
    }
  }
};
