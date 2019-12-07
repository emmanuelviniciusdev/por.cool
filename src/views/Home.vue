<template>
  <div class="columns is-desktop">
    <div class="column is-4-desktop">
      <LateralHeader />
    </div>
    <div class="column">opa</div>
  </div>
</template>

<script>
import LateralHeader from "../components/LateralHeader";

export default {
  name: "Home",
  components: {
    LateralHeader
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