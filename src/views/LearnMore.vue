<template>
  <div>
    <div class="columns is-multiline">
        <div class="column">
            <h1 class="title">sobre o porcool</h1>
            <p>antes de mais nada, meu nome é emmanuel e fui eu que dei vida a esta coisa estranha que você está vendo aqui. prazer.</p>
            <p>continuando, o porcool veio de uma necessidade minha de transformar uma planilha excel, que eu mesmo utilizava para fazer um pequeno controle dos meus gastos pessoais, para alguma coisa mais prática de ser utilizada porque, convenhamos, ficar editando planilha é MUITO CHATO, e eu já não aguentava mais. e foi aí que eu pensei: "hmmm, acho que não seria lá uma má ideia se eu juntasse os meus superpoderes de programador com a minha preguiça e ódio em editar planilhas para desenvolver alguma coisa que faça isso de maneira totalmente automática para mim".</p>
            <br>
            <p>e foi daí que surgiu o porcool. :)</p>
            <br>
            <p>eu sei que isto aqui é uma coisa simples, e pode parecer patético para algumas pessoas.</p>
            <p><i>mas eu desenvolvi isto daqui PARA MIM utilizar, tá ok? se não gostou, caia fora.</i></p>
            <br>
            <p>brincadeiras à parte, seja muito bem-vindo de qualquer jeito...</p>
            <br>
            <router-link :to="{name: 'signup'}">criar uma conta</router-link><br>
            <router-link :to="{name: 'signin'}">entrar</router-link>
        </div>
        <div class="column">
            <img src="../assets/images/pig2.png" alt="">
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
