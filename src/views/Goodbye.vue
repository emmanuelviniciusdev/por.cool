<template>
  <div>
    <div class="has-text-centered">
      <p class="is-size-5 goodbye-text">
        Muito obrigado, usuário(a), por ter feito parte do porcool. Não sabemos
        o motivo pelo qual você está partindo, mas saiba que você e cada um de
        nossos usuários significam muito 💚
      </p>
      <h1 class="title is-size-1 has-text-black">adeus</h1>
      <img
        src="../assets/images/crying-pig.png"
        alt="porco chorando"
        width="200"
      />
    </div>
  </div>
</template>

<script>
export default {
  name: "Goodbye"
};
</script>

<style lang="scss" scoped>
div .goodbye-text {
  width: 75% !important;
  margin: 0 auto;
  margin-bottom: 10px;
}
</style>
