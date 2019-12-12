<template>
  <div class="columns is-desktop">
    <div class="column is-4-desktop">
      <LateralHeader />
    </div>
    <div class="column">
      <h1 class="title has-text-black">meus gastos</h1>
      <SpendingTable />
    </div>
  </div>
</template>

<script>
import LateralHeader from "../components/LateralHeader";
import SpendingTable from "../components/SpendingTable";

export default {
  name: "Home",
  components: {
    LateralHeader,
    SpendingTable
  }
};
</script>

<style lang="scss" scoped>
@media screen and (min-width: 1024px) {
  .columns {
    margin-top: 30px;
  }
}
</style>