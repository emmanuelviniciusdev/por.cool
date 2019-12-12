<template>
  <div class="columns is-multiline">
    <div class="column is-12 is-6-desktop">
      <div class="introduction">
        <h1 class="has-text-black">
          fazer o controle dos seus gastos pessoais financeiros nunca foi tão
          fácil. e prático.
        </h1>

        <img src="../assets/images/pig1.png" alt="porcool <3" />

        <router-link :to="{ name: 'learn-more' }">
          <b-button type="is-link" outlined>clique para saber mais</b-button>
        </router-link>
      </div>
    </div>

    <div class="divider"></div>

    <div class="column">
      <div class="signin">
        <h1 class="title has-text-black">entre com a sua conta</h1>
        <form @submit.prevent="signin()">
          <b-field
            label="e-mail"
            :type="{'is-danger': hasInputErrorAndDirty('email')}"
            :message="{'insira o seu e-mail': isInvalidInputMsg('email', 'required'), 'insira um e-mail válido...': isInvalidInputMsg('email', 'email')}"
          >
            <b-input
              type="email"
              placeholder="seu@email.com"
              v-model.trim="form.email"
              @change.native="$v.form.email.$model = $event.target.value"
            ></b-input>
          </b-field>
          <b-field
            label="senha"
            :type="{'is-danger': hasInputErrorAndDirty('password')}"
            :message="{'insira a sua senha': isInvalidInputMsg('password', 'required')}"
          >
            <b-input type="password" placeholder="*******" v-model.trim="$v.form.password.$model"></b-input>
          </b-field>
          <button class="button is-primary">entrar</button>
        </form>

        <router-link :to="{name: 'signup'}">criar uma conta</router-link>
        <router-link to>recuperar minha senha</router-link>

        <div class="social-media-signin">
          <h1 class="title has-text-black">entre com a sua rede social</h1>
          <button class="button facebook">facebook</button>
          <button class="button google">google</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { required, email } from "vuelidate/lib/validators";

export default {
  name: "SignIn",
  data() {
    return {
      form: {
        email: null,
        password: null
      }
    };
  },
  validations: {
    form: {
      email: { required, email },
      password: { required }
    }
  },
  methods: {
    test(str) {
      alert(str);
    },
    signin() {
      if (!this.$v.form.$invalid) {
        console.log("signin");
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

<style lang="scss" scoped>
.columns {
  margin-top: 30px !important;
}

.introduction {
  text-align: center;

  h1 {
    font-size: 50px;
    line-height: 1.3;
  }

  img {
    display: none;
  }

  button {
    margin-top: 30px;
    width: 100%;
  }
}

.signin {
  text-align: center;

  form {
    margin-bottom: 20px;
    text-align: left;

    button {
      width: 100% !important;
      height: 45px !important;
      font-size: 20px;
    }
  }

  a {
    display: block;
    text-align: left;
  }

  .social-media-signin {
    margin-top: 20px;
    // background: #ccc;

    .title {
      font-size: 20px !important;
    }

    button {
      margin-top: -10px;
      margin-right: 5px;
    }
    .button.facebook {
      background: #4267b2;
      color: #fff;
      border: none;
    }
    .button.google {
      background: #db4437;
      color: #fff;
      border: none;
    }
  }
}

@media screen and (min-width: 769px) {
  .signin {
    width: 80%;
    // background: #ccc;
    margin: 0 auto;
  }

  .introduction {
    h1 {
      width: 80%;
      margin: 0 auto;
    }
    button {
      width: 200px;
    }
  }
}

@media screen and (min-width: 1024px) {
  .columns .divider {
    display: none;
  }

  .introduction {
    text-align: left;

    h1 {
      margin-left: 0;
      width: 100%;
      font-size: 55px;
    }
    img {
      display: block;
      margin-top: 20px;
      width: 500px;
    }
    button {
      margin-top: 10px;
    }
  }

  .signin {
    text-align: left;
    width: 70%;
    // background: #ccc;
  }
}
</style>
