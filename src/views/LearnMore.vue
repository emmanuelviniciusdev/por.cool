<template>
  <div>
    <div class="columns is-multiline">
      <div class="column">
        <h1 class="title">sobre o porcool</h1>
        <p class="is-italic has-text-weight-bold">História...</p>
        <p>
          <i
            >Antes de continuar, deixa eu me apresentar: meu nome é Emmanuel e
            fui eu que dei vida ao porcool.</i
          >
        </p>
        <p>
          O porcool surgiu de uma necessidade de transformar uma simples
          planilha de controle financeiro pessoal para alguma coisa mais
          "elaborada", que fizesse algumas coisas de maneira automatizada para
          mim.
        </p>
        <p>
          Editar planilhas é ruim, chato, horrível. Você pode até automatizá-las
          com suas funções matemáticas mirabolantes e deixá-las uns 2% mais
          usuais. Mas elas continuarão sendo planilhas.
        </p>
        <p>
          Então eu pensei:
          <i
            >eu poderia criar alguma coisa simples, que faça as coisas de
            maneira mais fácil e prática pro meu dia a dia...</i
          >
        </p>
        <p>
          No começo, a minha ideia era desenvolver alguma coisa só para mim
          utilizar, mas no meio do caminho me veio na cabeça que poderiam
          existir mais pessoas que também abonimassem planilhas, então eu
          desenvolvi a ferramenta de uma forma que tornasse possível sua
          utilização também por outras pessoas.
        </p>
        <p class="is-italic has-text-weight-bold">
          O que vai dar pra mim fazer com o porcool?
        </p>
        <p>
          Bom, por enquanto, o porcool é uma ferramenta muito simples e,
          enquanto estiver em sua versão beta, você conseguirá gerenciar os seus
          gastos do dia a dia e a sua renda mensal.
          <i>Tudo o que uma simples planilha faria, só que melhor.</i>
        </p>
        <p class="is-italic has-text-weight-bold">Tudo na nuvem</p>
        <p>
          Como o porcool é um aplicativo <i>web</i>, todos os seus dados ficarão
          disponíveis na nuvem e acessíveis de qualquer dispositivo que possua
          acessso à internet.
        </p>
        <br />
        <router-link :to="{ name: 'signup' }">criar uma conta</router-link
        ><br />
        <router-link :to="{ name: 'signin' }">entrar</router-link>
      </div>
      <div class="column">
        <img src="../assets/images/pig2.png" alt="" />
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "LearnMore"
};
</script>

<style scoped>
@media screen and (min-width: 769px) {
  img {
    margin-top: 80px;
  }
}
</style>
