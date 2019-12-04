<template>
  <div class="hero">
    <section class="hero-body">
      <div class="container">
        <div class="title">por.<span>cool</span></div>
        <router-view></router-view>
      </div>
    </section>
    <footer class="footer is-primary">
      <p>por.cool — produção independente</p>
      <p>dúvidas, ajuda? entre em contato pelo<br /><b>porcool@gmail.com</b></p>
    </footer>
  </div>
</template>

<script>
export default {
  name: "App"
};
</script>

<style lang="scss">
.hero .hero-body {
  .title {
    text-align: center;
    font-size: 28px;

    span {
      color: #1bb450;
    }
  }
}

@media screen and (min-width: 1024px) {
  .hero .hero-body .title {
    text-align: left;
  }
}
</style>
