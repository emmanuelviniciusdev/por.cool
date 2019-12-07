<template>
  <div class="columns">
    <div class="column">
      <h1 class="title has-text-black">para continuar, informe a sua renda fixa mensal</h1>
      <money
        v-model="income"
        v-bind="{
          decimal: ',',
          thousands: '.',
          prefix: 'R$',
          precision: 2
      }"
        v-if="!noFixedIncome"
        class="input-text-no-border"
      ></money>
      <div class="centralize">
        <b-checkbox class="no-fixed-income" v-model="noFixedIncome">eu não tenho uma renda fixa</b-checkbox>
      </div>
      <b-button type="is-primary">continuar</b-button>
    </div>
  </div>
</template>

<script>
import { Money } from "v-money";

export default {
  name: "DefineMonthlyIncome",
  components: {
    Money
  },
  data: () => ({
    income: 0,
    noFixedIncome: false
  })
};
</script>

<style scoped>
.v-money {
  width: 100%;
  /* background: gray; */
  font-size: 50px;
  color: gray;
  text-align: center;
}

.no-fixed-income {
  margin-top: 20px;
  /* margin: 0 auto; */
}

.centralize {
  text-align: center;
}

button {
  display: block;
  width: 70%;
  margin: 0 auto;
  margin-top: 20px;
}

@media screen and (min-width: 769px) {
    button {
        width: 50%;
    }
}

@media screen and (min-width: 1024px) {
  .v-money {
    text-align: left;
  }
  button {
    width: 150px;
    margin-left: 0;
  }
  .centralize {
      text-align: left;
  }
}
</style>