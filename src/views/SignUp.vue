<template>
  <div>
    <div class="columns">
      <div class="column">
        <div class="introduction">
          <h1 class="title">criar uma conta</h1>
          <p>para a alegria de alguns e tristeza de muitos, vivemos em uma sociedade capitalista.</p>
          <p>então, você terá que desembolsar uma bagatela de R$ 10,00 por mês para utilizar o porcool.</p>
        </div>
        <form>
          <div class="columns">
            <div class="column">
              <b-field label="nome">
                <b-input placeholder="seu nome"></b-input>
              </b-field>
              <b-field label="sobrenome">
                <b-input placeholder="seu sobrenome"></b-input>
              </b-field>
              <b-field label="e-mail">
                <b-input type="email" placeholder="seu@email.com"></b-input>
              </b-field>
            </div>
            <div class="column">
              <b-field label="senha">
                <b-input type="password" placeholder="**********"></b-input>
              </b-field>
              <b-field label="confirmar senha">
                <b-input type="password" placeholder="**********"></b-input>
              </b-field>
              <b-checkbox class="terms-checkbox">
                li e aceito os
                <router-link to>termos de uso</router-link>
              </b-checkbox>
              <b-field>
                <b-button type="is-primary" expanded>criar conta</b-button>
              </b-field>
            </div>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "SignUp"
};
</script>

<style>
.terms-checkbox {
  margin-top: 10px;
  margin-bottom: 20px !important;
}

.introduction {
    text-align: center;
    width: 70%;
    margin: 0 auto;
    margin-bottom: 20px;
}

@media screen and (min-width: 1024px) {
    .introduction {
        width: auto;
        margin-left: 0;
        text-align: left;
    }

    form {
        width: 70%;
    }
}
</style>