<template>
  <div>
    <div class="columns is-desktop">
      <div class="column is-4-desktop">
        <LateralHeader />
      </div>
      <div class="column">
        <h1 class="title has-text-black">
          adicionar novos gastos
          <Help>
              <template v-slot:title>como adicionar gastos?</template>
              <template v-slot:body>
                  <h1 class="subtitle">É muito simples</h1>
              </template>
          </Help>
        </h1>

        <InsertExpenses />
      </div>
    </div>
  </div>
</template>

<script>
import LateralHeader from "../components/LateralHeader";
import InsertExpenses from "../components/InsertExpenses";
import Help from "../components/Help";

export default {
  name: "AddExpenses",
  components: {
    LateralHeader,
    InsertExpenses,
    Help
  }
};
</script>

<style lang="scss" scoped>
@media screen and (min-width: 1024px) {
  .columns {
    margin-top: 30px;
  }
}
</style>