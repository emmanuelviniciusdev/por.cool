<template>
  <div>
    <b-menu>
      <b-menu-list>
        <b-menu-item
          label="meus gastos"
          icon="hand-holding-usd"
          @click="redirectTo('home')"
          
        ></b-menu-item>
        <b-menu-item label="novo gasto" icon="wallet" @click="redirectTo('')"></b-menu-item>
        <b-menu-item label="adicionar saldo" icon="dollar-sign" @click="redirectTo('')"></b-menu-item>
        <b-menu-item label="minha conta" icon="user" @click="redirectTo('')"></b-menu-item>
        <b-menu-item label="sair do porcool" icon="sad-tear" @click="redirectTo('')"></b-menu-item>
      </b-menu-list>
    </b-menu>
  </div>
</template>

<script>
export default {
  name: "MenuLinks",
  methods: {
    redirectTo(name) {
      if (name.trim() !== "" && this.currentRoute.name !== name)
        this.$router.push({ name });
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