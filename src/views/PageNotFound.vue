<template>
  <div>
    <div class="has-text-black has-text-centered">
      <h1 class="huge-404">404</h1>
      <img
        src="../assets/images/error.png"
        class="error-img"
        alt="página não encontrada"
      />
      <br />
      <p class="is-size-4">página não encontrada</p>
      <p class="is-size-6">{{ redirectMessage }}</p>
      <p class="is-size-5">
        <i>perdido?</i>
      </p>
    </div>
  </div>
</template>

<script>
export default {
  name: "PageNotFound",
  created() {
    const snackbarRef = this.$buefy.snackbar.open({
      message: "dentro de alguns segundos, você será redirecionado...",
      type: "is-warning",
      indefinite: true,
      onAction: () => {
        this.$router.push({ name: "signin" });
        clearInterval(redirectTimer);
      }
    });

    let redirectTimeInSec = 15;
    const redirectTimer = setInterval(() => {
      if (redirectTimeInSec <= 0) {
        clearInterval(redirectTimer);
        snackbarRef.close();
        this.$router.push({ name: "signin" });
        return;
      }
      redirectTimeInSec--;
    }, 1000);
  }
};
</script>

<style lang="scss" scoped>
.huge-404 {
  font-size: 150px;
  margin-top: -30px;
}

.error-img {
  width: 150px;
}

@media screen and (min-width: 1024px) {
  .huge-404 {
    font-size: 200px;
  }
}
</style>
