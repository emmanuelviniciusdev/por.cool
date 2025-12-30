<template>
  <div>
    <div class="mobile-menu">
      <b-button
        type="is-dark"
        outlined
        class="btn-open-menu"
        @click="isMenuOpened = !isMenuOpened"
        >menu</b-button
      >
      <MenuLinks v-if="isMenuOpened" />
    </div>
    <div class="desktop-menu">
      <MenuLinks />
    </div>
  </div>
</template>

<script>
import MenuLinks from "./MenuLinks";

export default {
  name: "Menu",
  components: {
    MenuLinks
  },
  data: () => ({
    isMenuOpened: false
  })
};
</script>

<style lang="scss" scoped>
.mobile-menu {
  margin-top: 20px;

  .btn-open-menu {
    width: 100%;
    margin-bottom: 15px;
  }
}

.desktop-menu {
  display: none;
}

@media screen and (min-width: 769px) {
  .mobile-menu {
    width: 50%;
    margin: 0 auto;
    margin-top: 20px;
  }
}

@media screen and (min-width: 1024px) {
  .mobile-menu {
    display: none;
  }

  .desktop-menu {
    display: block;
  }
}
</style>
