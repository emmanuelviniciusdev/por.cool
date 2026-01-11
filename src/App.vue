<template>
  <div class="hero is-fullheight">
    <section class="hero-body is-align-items-start">
      <div class="container">
        <router-link :to="{ name: 'signin' }" class="title app-title">
          por.<span>cool</span>
        </router-link>
        <transition name="slide" mode="out-in">
          <router-view></router-view>
        </transition>
      </div>
    </section>
    <footer class="hero-foot footer is-primary">
      <div class="columns">
        <div class="column">
          <p class="title is-6"><b>v1.3.0</b> — porcool — produção independente</p>
        </div>
        <div class="column">
          <p class="title is-6">fique por dentro</p>
          <router-link :to="{ name: 'terms-of-use' }">Termos de Uso</router-link>
          <br>
          <router-link :to="{ name: 'privacy-policy' }">Política de Privacidade</router-link>
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
  // Garante que o container ocupe toda a largura para transições funcionarem
  .container {
    width: 100%;
  }

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
