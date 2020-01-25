// Services
import balancesService from '../../services/balances';

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
        },

        SET_CURRENT_BALANCE(state, currentBalance) {
            state.currentBalance = currentBalance;
        },
        SET_LAST_MONTH_BALANCE(state, lastMonthBalance) {
            state.lastMonthBalance = lastMonthBalance;
        }
    },

    actions: {
        toggleShowBalance({ commit }) {
            commit('TOGGLE_SHOW_BALANCE');
        },

        async setCurrentBalance({ commit }, { userUid, spendingDate }) {
            const remainingBalance = await balancesService.calculate({ userUid, spendingDate });
            commit('SET_CURRENT_BALANCE', remainingBalance);
        },
        setLastMonthBalance({ commit }, lastMonthBalance) {
            commit('SET_LAST_MONTH_BALANCE', lastMonthBalance);
        },
    }
}