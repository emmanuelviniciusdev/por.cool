import dateAndTimeHelper from "../../helpers/dateAndTime";

export default {
  namespaced: true,

  state: {
    user: {}
  },

  mutations: {
    SET(state, user) {
      if (
        user.lookingAtSpendingDate &&
        !(user.lookingAtSpendingDate instanceof Date)
      )
        user.lookingAtSpendingDate = dateAndTimeHelper.transformSecondsToDate(
          user.lookingAtSpendingDate.seconds
        );

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
      commit("SET", user);
    },
    update({ commit }, userData) {
      commit("UPDATE", userData);
    }
  }
};
