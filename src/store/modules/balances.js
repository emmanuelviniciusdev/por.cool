import moment from 'moment'; 

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

        SET_BALANCES(state, { currentBalance, lastMonthBalance }) {
            state.currentBalance = currentBalance;
            state.lastMonthBalance = lastMonthBalance;
        },
        SET_CURRENT_BALANCE(state, currentBalance) {
            state.currentBalance = currentBalance;
        },
    },

    actions: {
        toggleShowBalance({ commit }) {
            commit('TOGGLE_SHOW_BALANCE');
        },

        async setBalances({ commit }, { userUid, spendingDate }) {
            const lastMonthSpendingDate = moment(spendingDate).subtract(1, 'months').toDate();

            const currentBalance = await balancesService.calculate({ userUid, spendingDate });
            const lastMonthBalance = await balancesService.getHistoryByDate({ userUid, spendingDate: lastMonthSpendingDate });

            commit('SET_BALANCES', {
                currentBalance,
                lastMonthBalance: lastMonthBalance.balance ? lastMonthBalance.balance : 0
            });
        },
        async setCurrentBalance({ commit }, { userUid, spendingDate }) {
            const remainingBalance = await balancesService.calculate({ userUid, spendingDate });
            commit('SET_CURRENT_BALANCE', remainingBalance);
        }
    }
}