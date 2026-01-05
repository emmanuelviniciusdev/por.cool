<template>
  <div class="columns payment-wrapper">
    <b-loading :active="!user"></b-loading>

    <div class="column" v-if="user">
      <div class="welcome" v-if="user.isNewUser">
        <h1 class="title has-text-black">
          É isso aí, {{ user.displayName | capitalizeName }}! Você já está quase
          lá.
        </h1>
        <h2 class="subtitle has-text-black">
          Para continuar, efetue o pagamento no valor de R$ 10,00.
        </h2>
        <div class="notification is-warning">
          <b
            >O pagamento deverá ser realizado a cada 30 dias, mas não existe
            nenhum tipo de vínculo que te prenda e te obrigue a pagar todo mês.
            Você só paga quando quiser utilizar.</b
          >
        </div>
        <div class="notification is-info" v-if="!user.requestedPayment">
          Em breve, uma solicitação de pagamento via
          <b>paypal</b> será enviada para o seu e-mail e, assim que aprovado o
          pagamento, a sua conta será liberada.
        </div>
      </div>

      <div class="levy" v-if="!user.isNewUser">
        <h1 class="title has-text-black">
          Oh, {{ user.displayName | capitalizeName }}. Os seus 30 dias de
          utilização se expiraram e você ainda não efetuou um novo pagamento
          para continuar utilizando o porcool.
        </h1>
        <h2 class="subtitle has-text-black">
          Sem a ajuda do porcool, a sua vida financeira fica uma bagunça 😱!1!
          Não perca tempo e PAGUE agora mesmo!11!!!
        </h2>
        <div class="notification is-info" v-if="!user.requestedPayment">
          Uma solicitação de pagamento via
          <b>paypal</b> será enviada em breve para o seu e-mail.
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import firebase from "firebase/app";
import "firebase/auth";
import "firebase/firestore";
import paymentHelper from "../helpers/payment";
import filters from "../filters";

export default {
  name: "Payment",
  data() {
    return {
      user: null
    };
  },
  filters: {
    capitalizeName: filters.capitalizeName
  },
  beforeCreate() {
    const unsubscribe = firebase.auth().onAuthStateChanged(async user => {
      unsubscribe();

      if (user) {
        // Check if user payment is ok. If so, we'll redirect user to home page.
        const lastUserPayment = await firebase
          .firestore()
          .collection("payments")
          .where("user", "==", user.uid)
          .orderBy("paymentDate", "desc")
          .limit(1)
          .get();

        if (!lastUserPayment.empty) {
          const remainingDays = paymentHelper.remainingDays(
            lastUserPayment.docs[0].data().paymentDate
          );

          if (remainingDays > 0) {
            this.$router.push({ name: "home" });
          }
        }
      }
    });
  },
  created() {
    /**
     * (pendingPayment && !requestedPayment && !paidPayment) => não fazer nada...
     * (pendingPayment && requestedPayment && !paidPayment) => não fazer nada...
     * (!pendingPayment && requestedPayment && paidPayment) => inserir dados em 'payments'...
     * (!pendingPayment && !requestedPayment && !paidPayment && remainingDays <= 0) => pendingPayment = true
     */
    firebase.auth().onAuthStateChanged(async user => {
      if (user) {
        const users = firebase.firestore().collection("users");
        const payments = firebase.firestore().collection("payments");

        const userInfo = await users.doc(user.uid).get();
        const {
          monthlyIncome,
          pendingPayment,
          requestedPayment,
          paidPayment
        } = userInfo.data();

        const lastUserPayment = await payments
          .where("user", "==", user.uid)
          .orderBy("paymentDate", "desc")
          .limit(1)
          .get();

        if (!pendingPayment && requestedPayment && paidPayment) {
          const { FieldValue } = firebase.firestore;

          await payments.add({ user: user.uid, paymentDate: new Date(), onPremiseSyncDatetime: null });
          await users.doc(user.uid).update({
            pendingPayment: FieldValue.delete(),
            requestedPayment: FieldValue.delete(),
            paidPayment: FieldValue.delete()
          });

          this.$router.push({
            name: !monthlyIncome ? "define-monthly-income" : "home"
          });

          return;
        }

        const remainingDays = !lastUserPayment.empty
          ? paymentHelper.remainingDays(
              lastUserPayment.docs[0].data().paymentDate
            )
          : 0;

        if (
          !pendingPayment &&
          !requestedPayment &&
          !paidPayment &&
          remainingDays <= 0
        ) {
          await users.doc(user.uid).update({ pendingPayment: true });
        }

        this.user = user;
        // If user has no registered payments, then it's a new user and we'll show a
        // welcome message
        this.user.isNewUser = lastUserPayment.empty;
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

.payment-wrapper {
  width: 100%;
  max-width: 900px;
}
</style>
