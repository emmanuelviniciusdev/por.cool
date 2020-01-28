<template>
  <div class="hero">
    <section class="hero-body">
      <div class="container">
        <router-link :to="{ name: 'signin' }" class="title app-title">
          por.<span>cool</span>
        </router-link>
        <transition name="slide" mode="out-in">
          <router-view></router-view>
        </transition>
      </div>
    </section>
    <footer class="footer is-primary">
      <div class="columns">
        <div class="column">
          <h6 class="title is-6">informações</h6>
          <p>
            <b>closed beta</b> — porcool — produção independente
          </p>
          <p>
            dúvidas, ajuda? entre em contato pelo
            <br />
            <b>porcool@gmail.com</b>
          </p>
        </div>

        <div class="column">
          <h6 class="title is-6">fique por dentro</h6>
          <router-link :to="{name: 'terms-of-use'}">Termos de Uso</router-link><br>
          <router-link :to="{name: 'privacy-policy'}">Política de Privacidade</router-link>
        </div>
      </div>
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
  .app-title {
    width: auto;
    text-align: center;
    font-size: 28px;
    display: inline-block;
    cursor: pointer;
    user-select: none;
    // background: red;

    span {
      color: #1bb450;
    }
  }
}

@media screen and (min-width: 1024px) {
  .hero .hero-body .app-title {
    text-align: left;
  }
}
</style>
