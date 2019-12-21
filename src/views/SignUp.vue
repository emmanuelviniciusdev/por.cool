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
                :type="{'is-danger': hasInputErrorAndDirty('lastName')}"
                :message="{
                  'insira o seu sobrenome': isInvalidInputMsg('lastName', 'required'),
                  'o sobrenome é muito curto': isInvalidInputMsg('lastName', 'minLength'),
                  'o sobrenome é muito grande': isInvalidInputMsg('lastName', 'maxLength'),
                  }"
              >
                <b-input
                  placeholder="seu sobrenome"
                  maxlength="50"
                  v-model.trim="form.lastName"
                  @change.native="$v.form.lastName.$model = $event.target.value"
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
                  'no mínimo, sua senha deve conter 6 caracteres': isInvalidInputMsg('password', 'minLength')
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
                  'no mínimo, sua senha deve conter 6 caracteres': isInvalidInputMsg('cPassword', 'minLength'),
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
                <!-- <button class="button is-primary btn-signup" :disabled="form.termsOfUse" :loading="true">criar conta</button> -->
                <b-button
                  type="is-primary"
                  expanded
                  :disabled="!form.termsOfUse"
                  :loading="formLoading"
                  native-type="submit"
                >criar conta</b-button>
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
import * as firebase from "firebase/app";
import "firebase/auth";
import "firebase/database";
import settings from "../services/settings";

export default {
  name: "SignUp",
  data() {
    return {
      form: {
        name: null,
        lastName: null,
        email: null,
        password: null,
        cPassword: null,
        termsOfUse: false
      },
      formLoading: false
    };
  },
  validations: {
    form: {
      name: { required, minLength: minLength(2), maxLength: maxLength(50) },
      lastName: { required, minLength: minLength(2), maxLength: maxLength(50) },
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
    async signUp() {
      this.loading();

      if (this.$v.form.$invalid) {
        this.$buefy.toast.open({
          message: "ei, você, preencha todos os campos corretamente",
          type: "is-danger",
          position: "is-bottom",
          duration: 5000
        });

        this.loading(false);

        return;
      }

      // Check if system is under maintenance
      const {
        blockUserRegistration,
        maintenance
      } = await settings.checkMaintenances();

      if (blockUserRegistration || maintenance) {
        let message = blockUserRegistration
          ? blockUserRegistration
          : maintenance
          ? maintenance
          : "Ocorreu um erro inesperado. Por favor, tenta novamente mais tarde.";

        this.$buefy.toast.open({
          message,
          type: "is-warning",
          position: "is-bottom",
          duration: 5000
        });

        this.loading(false);

        return;
      }

      try {
        const { name, lastName, email, password } = this.form;

        const user = await firebase
          .auth()
          .createUserWithEmailAndPassword(email, password);

        if (!user) {
          this.$buefy.toast.open({
            message: "desculpe, ocorreu um erro ao efetuar o seu cadastro",
            type: "is-danger",
            position: "is-bottom",
            duration: 5000
          });

          this.loading(false);

          return;
        }

        // Set additional user data
        await firebase.auth().currentUser.updateProfile({ displayName: name });
        await firebase
          .database()
          .ref(`users/${user.user.uid}`)
          .set({
            name: name.toLowerCase(),
            lastName: lastName.toLowerCase()
          });

        this.clearForm();

        this.$buefy.toast.open({
          message: "bem-vindx ao porcool!!",
          type: "is-success",
          position: "is-bottom",
          duration: 5000
        });

        // TODO: Redirect user to payment page
      } catch (err) {
        console.log(err);
        this.$buefy.toast.open({
          message: this.authErrorMessages(err.code),
          type: "is-danger",
          position: "is-bottom",
          duration: 5000
        });
      } finally {
        this.loading(false);
      }
    },
    hasInputErrorAndDirty(input) {
      return this.$v.form[input].$error && this.$v.form[input].$dirty;
    },
    isInvalidInputMsg(input, role) {
      return !this.$v.form[input][role] && this.$v.form[input].$error;
    },
    authErrorMessages(errorCode) {
      const errorMessages = {
        "auth/email-already-in-use": "este e-mail já está sendo utilizado",
        "auth/weak-password": "no mínimo, sua senha deve conter 6 caracteres"
      };

      return errorMessages[errorCode]
        ? errorMessages[errorCode]
        : "parece que ocorreu um erro ao tentar criar a sua conta";
    },
    loading(loading = true) {
      this.formLoading = loading;
    },
    clearForm() {
      Object.keys(this.form).forEach(k => (this.form[k] = null));
      this.$v.form.$reset();
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