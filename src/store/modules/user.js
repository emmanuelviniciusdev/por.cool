import dateAndTimeHelper from '../../helpers/dateAndTime';

export default {
    namespaced: true,

    state: {
        user: {}
    },

    mutations: {
        SET(state, user) {
            state.user = user;
        },
        UPDATE(state, userData) {
            Object.keys(userData).map(data => {
                state.user[data] = userData[data];
            });
        }
    },

    actions: {
        set({ commit }, user) {
            user.lookingAtSpendingDate = dateAndTimeHelper.transformSecondsToDate(user.lookingAtSpendingDate.seconds);
            commit('SET', user);
        },
        update({ commit }, userData) {
            commit('UPDATE', userData);
        }
    }
}