<template>
  <div class="hero">
    <section class="hero-body">
      <div class="container">
        <router-link :to="{ name: 'signin' }" class="title"
          >por.<span>cool</span></router-link
        >
        <transition name="slide" mode="out-in">
          <router-view></router-view>
        </transition>
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
    display: block;

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
