export default {
    namespaced: true,

    state: {
        showBalance: false,
        currentBalance: 0,
        lastMonthBalance: 0
    },

    mutations: {
        TOGGLE_SHOW_BALANCE(state) {
            state.showBalance = !state.showBalance;
        }
    },

    actions: {
        toggleShowBalance({ commit }) {
            commit('TOGGLE_SHOW_BALANCE');
        }
    }
}