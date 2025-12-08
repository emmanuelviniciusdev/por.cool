// Services
import banksService from '../../services/banks';

export default {
    namespaced: true,

    state: {
        banksList: [],
        loadingBanksList: false,
        loadingBanksListError: false,
    },

    mutations: {
        SET_BANKS_LIST(state, banksList) {
            state.banksList = banksList;
        },
        SET_LOADING_BANKS_LIST(state, isLoading = true) {
            state.loadingBanksList = isLoading;
        },
        SET_LOADING_BANKS_LIST_ERROR(state, hasError = true) {
            state.loadingBanksListError = hasError;
        },
    },

    actions: {
        async setBanksList({ commit }, { userUid }) {
            commit('SET_LOADING_BANKS_LIST');
            commit('SET_LOADING_BANKS_LIST_ERROR', false);

            try {
                const banksList = await banksService.getBanks(userUid);
                commit('SET_BANKS_LIST', banksList);
            } catch (err) {
                console.error('Error in setBanksList:', err);
                commit('SET_LOADING_BANKS_LIST_ERROR');
            } finally {
                commit('SET_LOADING_BANKS_LIST', false);
            }
        }
    }
}