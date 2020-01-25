<template>
  <div>
    <AvailableBalance />
    <Menu class="menu" />
  </div>
</template>

<script>
// Components
import Menu from "./Menu";
import AvailableBalance from "./AvailableBalance";

export default {
  name: "LateralHeader",
  components: {
    Menu,
    AvailableBalance
  }
};
</script>

<style lang="scss" scoped>
.menu {
  clear: both;
}
</style>