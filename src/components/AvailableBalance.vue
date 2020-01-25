<template>
  <div class="balance-info" @click="$store.dispatch('balances/toggleShowBalance')">
    <div class="level available-balance">
      <div class="level-item">
        <b-icon
          :icon="balances.showBalance ? 'eye' : 'eye-slash'"
          size="is-medium"
          class="has-text-black"
        ></b-icon>
      </div>
      <div class="level-item">
        <span class="has-text-black">saldo disponível</span>
      </div>
    </div>
    <b-tooltip
      type="is-dark"
      multilined
      position="is-bottom"
      :label="balances.showBalance ? `${$options.filters.currency(balances.lastMonthBalance)} (do mês passado) + ${$options.filters.currency(balances.currentBalance)} (deste mês)` : 'ver saldo'"
    >
      <h1
        class="title has-text-black"
      >{{balances.showBalance ? $options.filters.currency(balances.currentBalance) : '...'}}</h1>
    </b-tooltip>
  </div>
</template>

<script>
import { mapState } from "vuex";

export default {
  name: "AvailableBalance",
  computed: {
    ...mapState(["balances"])
  },
};
</script>

<style lang="scss" scoped>
.balance-info {
  text-align: center;
  cursor: pointer;
  user-select: none;

  .available-balance {
    width: 180px;
    // background: green;
    margin: 0 auto;
    font-size: 18px;

    .level-item:nth-child(2) {
      margin-top: -4px;
      margin-left: 10px;
    }
  }

  .b-tooltip {
    clear: both;
    display: block;
    text-align: center;
  }
}

@media screen and (min-width: 1024px) {
  .balance-info {
    text-align: left;
    margin-bottom: 20px;
    // background: blue;
    position: relative;

    .available-balance {
      float: left;
    }

    .b-tooltip {
      text-align: left;
      padding-top: 7px;
    }
  }
}
</style>