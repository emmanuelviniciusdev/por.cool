import moment from 'moment';

// Services
import balancesService from '../../services/balances';

export default {
    namespaced: true,

    state: {
        showBalance: false,
        currentBalance: 0,
        lastMonthBalance: 0,

        balancesList: [],
        loadingBalancesList: false,
        loadingBalancesListError: false,
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

        SET_BALANCES_LIST(state, balancesList) {
            state.balancesList = balancesList;
        },
        SET_LOADING_BALANCES_LIST(state, isLoading = true) {
            state.loadingBalancesList = isLoading;
        },
        SET_LOADING_BALANCES_LIST_ERROR(state, hasError = true) {
            state.loadingBalancesListError = hasError;
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
        },

        async setBalancesList({ commit }, { userUid, spendingDate }) {
            commit('SET_LOADING_BALANCES_LIST');
            commit('SET_LOADING_BALANCES_LIST_ERROR', false);

            try {
                const balancesList = await balancesService.getAdditionalBalances({ userUid, spendingDate });
                commit('SET_BALANCES_LIST', balancesList);
            } catch {
                commit('SET_LOADING_BALANCES_LIST_ERROR');
            } finally {
                commit('SET_LOADING_BALANCES_LIST', false);
            }
        }
    }
}