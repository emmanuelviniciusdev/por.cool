<template>
  <div>
    <b-menu>
      <b-menu-list>
        <b-menu-item
          label="novo gasto"
          icon="wallet"
          @click="redirectTo('add-expenses')"
          :active="isActiveRoute('add-expenses')"
        ></b-menu-item>
        <b-menu-item
          label="meus gastos"
          icon="dollar-sign"
          @click="redirectTo('home')"
          :active="isActiveRoute('home')"
        ></b-menu-item>
        <b-menu-item
          label="meus saldos"
          icon="hand-holding-usd"
          @click="redirectTo('balances')"
          :active="isActiveRoute('balances')"
        ></b-menu-item>
        <b-menu-item
          label="bancos e instituições"
          icon="university"
          @click="redirectTo('banks')"
          :active="isActiveRoute('banks')"
        ></b-menu-item>
        <b-menu-item
          label="minha conta"
          icon="user"
          @click="redirectTo('my-account')"
          :active="isActiveRoute('my-account')"
        ></b-menu-item>
        <b-menu-item label="sair" icon="sad-tear" @click="signOut()"></b-menu-item>
      </b-menu-list>
    </b-menu>
  </div>
</template>

<script>
import firebase from "firebase/app";
import "firebase/auth";

export default {
  name: "MenuLinks",
  methods: {
    redirectTo(name) {
      if (name.trim() !== "" && this.currentRoute.name !== name)
        this.$router.push({ name });
    },
    isActiveRoute(name) {
      return this.currentRoute.name === name;
    },
    signOut() {
      firebase.auth().signOut();
      this.$store.dispatch("user/set", {});
      this.$router.push({ name: "signin" });
    }
  },
  computed: {
    currentRoute() {
      return this.$router.currentRoute;
    }
  }
};
</script>

<style>
</style>