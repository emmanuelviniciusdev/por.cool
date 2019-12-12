<template>
  <div>
    <div class="columns">
      <div class="column">
        <div class="introduction">
          <h1 class="title">criar uma conta</h1>
          <p>para a alegria de alguns e tristeza de muitos, vivemos em uma sociedade capitalista.</p>
          <p>então, você terá que desembolsar uma bagatela de R$ 10,00 por mês para utilizar o porcool.</p>
        </div>
        <form @submit.prevent="signUp()">
          <div class="columns">
            <div class="column">
              <b-field
                label="nome"
                :type="{'is-danger': hasInputErrorAndDirty('name')}"
                :message="{
                  'insira o seu nome': isInvalidInputMsg('name', 'required'),
                  'o nome é muito curto': isInvalidInputMsg('name', 'minLength'),
                  'o nome é muito grande': isInvalidInputMsg('name', 'maxLength'),
                  }"
              >
                <b-input
                  placeholder="seu nome"
                  maxlength="50"
                  v-model.trim="form.name"
                  @change.native="$v.form.name.$model = $event.target.value"
                ></b-input>
              </b-field>

              <b-field
                label="sobrenome"
                :type="{'is-danger': hasInputErrorAndDirty('lastname')}"
                :message="{
                  'insira o seu sobrenome': isInvalidInputMsg('lastname', 'required'),
                  'o sobrenome é muito curto': isInvalidInputMsg('lastname', 'minLength'),
                  'o sobrenome é muito grande': isInvalidInputMsg('lastname', 'maxLength'),
                  }"
              >
                <b-input
                  placeholder="seu sobrenome"
                  maxlength="50"
                  v-model.trim="form.lastname"
                  @change.native="$v.form.lastname.$model = $event.target.value"
                ></b-input>
              </b-field>

              <b-field
                label="e-mail"
                :type="{'is-danger': hasInputErrorAndDirty('email')}"
                :message="{
                  'insira o seu e-mail': isInvalidInputMsg('email', 'required'),
                  'insira um e-mail válido...': isInvalidInputMsg('email', 'email')
                  }"
              >
                <b-input
                  type="email"
                  placeholder="seu@email.com"
                  v-model.trim="form.email"
                  @change.native="$v.form.email.$model = $event.target.value"
                ></b-input>
              </b-field>
            </div>
            <div class="column">
              <b-field
                label="senha"
                :type="{'is-danger': hasInputErrorAndDirty('password')}"
                :message="{
                  'insira uma senha segura': isInvalidInputMsg('password', 'required'),
                  'sua senha está muito curta': isInvalidInputMsg('password', 'minLength')
                  }"
              >
                <b-input
                  type="password"
                  placeholder="**********"
                  v-model.trim="form.password"
                  @change.native="$v.form.password.$model = $event.target.value"
                ></b-input>
              </b-field>

              <b-field
                label="confirmar senha"
                :type="{'is-danger': hasInputErrorAndDirty('cPassword')}"
                :message="{
                  'insira uma senha segura': isInvalidInputMsg('cPassword', 'required'),
                  'sua senha está muito curta': isInvalidInputMsg('cPassword', 'minLength'),
                  'as duas senhas não batem': isInvalidInputMsg('cPassword', 'sameAsPassword'),
                  }"
              >
                <b-input
                  type="password"
                  placeholder="**********"
                  v-model.trim="form.cPassword"
                  @change.native="$v.form.cPassword.$model = $event.target.value"
                ></b-input>
              </b-field>

              <b-checkbox class="terms-checkbox" v-model="form.termsOfUse">
                li e aceito os
                <router-link to>termos de uso</router-link>
              </b-checkbox>

              <b-field>
                <button
                  class="button is-primary btn-signup"
                  :disabled="!form.termsOfUse"
                >criar conta</button>
              </b-field>
            </div>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import {
  required,
  email,
  minLength,
  maxLength,
  sameAs
} from "vuelidate/lib/validators";

export default {
  name: "SignUp",
  data() {
    return {
      form: {
        name: null,
        lastname: null,
        email: null,
        password: null,
        cPassword: null,
        termsOfUse: false
      }
    };
  },
  validations: {
    form: {
      name: { required, minLength: minLength(2), maxLength: maxLength(50) },
      lastname: { required, minLength: minLength(2), maxLength: maxLength(50) },
      email: { required, email },
      password: { required, minLength: minLength(6) },
      cPassword: {
        required,
        minLength: minLength(6),
        sameAsPassword: sameAs("password")
      }
    }
  },
  methods: {
    signUp() {
      const { $invalid } = this.$v.form;

      if ($invalid) {
        this.$buefy.toast.open({
          message: 'ei, você, preencha todos os campos corretamente',
          type: 'is-danger',
          position: 'is-bottom',
          duration: 5000,
        });

        return false;
      }

      if (!$invalid) {
        console.log('signup');
      }
    },
    // Validation checks
    hasInputErrorAndDirty(input) {
      return this.$v.form[input].$error && this.$v.form[input].$dirty;
    },
    isInvalidInputMsg(input, role) {
      return !this.$v.form[input][role] && this.$v.form[input].$error;
    }
  }
};
</script>

<style scoped>
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

.btn-signup {
  width: 100%;
}

@media screen and (min-width: 1024px) {
  .introduction {
    width: auto;
    margin-left: 0;
    text-align: left;
  }

  form {
    width: 70%;
    margin-top: -20px;
  }
}
</style>