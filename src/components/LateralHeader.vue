<template>
  <div>
    <div class="balance-info" @click="showBalance = !showBalance">
      <div class="level available-balance">
        <div class="level-item">
          <b-icon icon="eye" size="is-medium" class="has-text-black" v-if="showBalance"></b-icon>
          <b-icon icon="eye-slash" size="is-medium" class="has-text-black" v-else></b-icon>
        </div>
        <div class="level-item">
          <span class="has-text-black">saldo disponível</span>
        </div>
      </div>
      <b-tooltip
        type="is-dark"
        multilined
        position="is-bottom"
        label="R$ 100,75 (mês passado) + R$ 100,00 (agora)"
        v-if="showBalance"
      >
        <h1 class="title has-text-black">R$ 200,75</h1>
      </b-tooltip>
      <b-tooltip type="is-dark" multilined position="is-bottom" label="..." v-else>
        <h1 class="title has-text-black">...</h1>
      </b-tooltip>
    </div>

    <Menu class="menu" />
  </div>
</template>

<script>
import Menu from "./Menu";

export default {
  name: "LateralHeader",
  components: {
    Menu
  },
  data: () => ({
    showBalance: false
  })
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

.menu {
  clear: both;
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