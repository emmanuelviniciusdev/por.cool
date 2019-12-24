<template>
  <div class="columns">
    <b-loading :active="!user"></b-loading>

    <div class="column" v-if="user">
      <div class="welcome" v-if="user && user.isNewUser">
        <h1 class="title has-text-black">É isso aí, {{user.displayName}}! Você já está quase lá.</h1>
        <h2 class="subtitle has-text-black">Para continuar, efetue o pagamento no valor de R$ 10,00.</h2>
        <div class="notification is-warning">
          <b>O pagamento deverá ser realizado a cada 30 dias, mas não existe nenhum tipo de vínculo que te prenda e te obrigue a pagar todo mês. Você só paga quando quiser utilizar.</b>
        </div>
        <div class="notification is-info">
          Em breve, uma solicitação de pagamento via <b>paypal</b> será enviada para o seu e-mail e, assim que aprovado o pagamento, a sua conta será liberada.
        </div>
      </div>

      <div class="levy" v-if="user && !user.isNewUser">
        <h1
          class="title has-text-black"
        >Oh, {{user.displayName}}. Os seus 30 dias de utilização se expiraram e você ainda não efetuou um novo pagamento para continuar utilizando o porcool.</h1>
        <h2
          class="subtitle has-text-black"
        >Sem a ajuda do porcool, a sua vida financeira fica uma bagunça 😱!1! Não perca tempo e PAGUE agora mesmo!11!!!</h2>
        <div class="notification is-info">
          Uma solicitação de pagamento via <b>paypal</b> será enviada em breve para o seu e-mail.
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import firebase from "firebase/app";
import "firebase/auth";
import "firebase/database";
import "firebase/firestore";

export default {
  name: "Payment",
  data() {
    return {
      user: null
    };
  },
  methods: {
    capitalizeName(name) {
      return name
        .split(" ")
        .map(namePart => {
          if (namePart !== "de" && namePart !== "do" && namePart !== "da")
            namePart = namePart.charAt(0).toUpperCase() + namePart.slice(1);
          return namePart;
        })
        .join(" ");
    }
  },
  created() {  
    firebase.auth().onAuthStateChanged(async user => {
      if (user) {
        // If user has no registered payments, then it's a new user
        const payments = await firebase
          .firestore()
          .collection("payments")
          .where("user", "==", user.uid)
          .get();

        this.user = user;
        this.user.displayName = this.capitalizeName(this.user.displayName);
        this.user.isNewUser = payments.empty;
      }
    });
  }
};
</script>

<style lang="scss" scoped>
.welcome p {
  font-weight: bold;
  color: red;
}
</style>