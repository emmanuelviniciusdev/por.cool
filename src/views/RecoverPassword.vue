<template>
  <div>
    <h1 class="title">recuperar senha</h1>
    <p>
      se você se esqueceu de sua senha, você pode recuperá-la inserindo abaixo o
      e-mail registrado em sua conta
    </p>
    <form @submit.prevent="submit()">
      <b-field
        label="e-mail"
        :type="{ 'is-danger': hasInputErrorAndDirty('email') }"
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
      <b-field>
        <b-button
          type="is-primary"
          native-type="submit"
          :loading="loadingRecoverPassword"
          >recuperar senha</b-button
        >
      </b-field>
    </form>
  </div>
</template>

<script>
import { required, email } from "vuelidate/lib/validators";
import firebase from "firebase/app";
import "firebase/auth";

// Services
import authService from "../services/auth";

export default {
  name: "RecoverPassword",
  data() {
    return {
      form: {
        email: ""
      },
      loadingRecoverPassword: false
    };
  },
  validations: {
    form: {
      email: { required, email }
    }
  },
  methods: {
    async submit() {
      if (this.$v.form.$invalid) return;

      this.loadingRecoverPassword = true;

      try {
        const recoverPassword = await authService.recoverPassword(
          this.form.email
        );

        if (recoverPassword.error) {
          this.$buefy.toast.open({
            message: recoverPassword.message,
            type: "is-danger",
            position: "is-bottom"
          });

          return;
        }

        this.$buefy.toast.open({
          message: recoverPassword.message,
          type: "is-success",
          position: "is-bottom",
          duration: 5000
        });

        this.$v.form.$reset();
        this.form.email = "";
      } catch (err) {
        let message = "ocorreu um erro ao enviar a recuperação de senha";

        if (err.code === "auth/user-not-found")
          message = "e-mail não encontrado no porcool";

        this.$buefy.toast.open({
          message,
          type: "is-danger",
          position: "is-bottom"
        });
      } finally {
        this.loadingRecoverPassword = false;
      }
    },
    hasInputErrorAndDirty(input) {
      return this.$v.form[input].$error && this.$v.form[input].$dirty;
    },
    isInvalidInputMsg(input, role) {
      return !this.$v.form[input][role] && this.$v.form[input].$error;
    }
  },
  beforeCreate() {
    const unsubscribe = firebase.auth().onAuthStateChanged(user => {
      unsubscribe();
      if (user) {
        this.$router.push({ name: "home" });
      }
    });
  }
};
</script>

<style lang="scss" scoped>
@media screen and (min-width: 769px) {
  p {
    width: 500px;
  }

  form {
    width: 400px;
    margin-top: 20px;
  }
}
</style>
