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
              <div class="content">
                <h1 class="subtitle has-text-weight-normal">É muito simples!</h1>
                <h2
                  class="subtitle has-text-weight-normal"
                >Você só precisa preencher alguns campos. São eles:</h2>
                <ul type="1">
                  <li>
                    <p class="is-size-5 has-text-weight-normal">
                      <b>gasto</b> — É o nome que você vai dar ao seu gasto, por exemplo: "conta de luz".
                    </p>
                  </li>
                  <li>
                    <p class="is-size-5 has-text-weight-normal">
                      <b>valor total</b> — Valor total do seu gasto.
                    </p>
                  </li>
                  <li>
                    <p class="is-size-5 has-text-weight-normal">
                      <b>status</b> — É onde você vai colocar o estado do seu gasto para o mês atual. Eles podem ser classificados como:
                      "pendente", quando você ainda não efetuou o pagamento; "parcialmente pago", quando você já efetuou apenas uma parte
                      do pagamento; e "pago", quando você já efetuou o pagamento.
                      <b>
                        <i>
                          Mas atenção: este status não influenciará EM NADA
                          no cálculo dos seus saldos, ele só servirá como uma espécie de referência e organização mesmo.
                        </i>
                      </b>
                    </p>
                  </li>
                  <li>
                    <p class="is-size-5 has-text-weight-normal">
                      <b>valor que já foi pago</b> — O valor que já está parcialmente pago.
                    </p>
                  </li>
                  <li>
                    <p class="is-size-5 has-text-weight-normal">
                      <b>tipo</b> — Aqui vai uma coisa bem interessante: você pode classificar os seus gastos por um determinado tipo.
                      São eles: "gasto", "fatura" e "poupança". Quando você define o tipo "gasto", você está dizendo que o seu
                      gasto só irá valer para o mês atual e não irá se repetir para os próximos meses. Entretanto, ao definir
                      o tipo "fatura" ou "poupança", você tem a possibilidade de definir uma vigência, ou seja, um prazo de validade
                      para o seu gasto, e ele se repetirá durante todos os meses na sua lista de gastos até esse prazo de validade
                      se expirar.
                      <i>Legal, né?</i> E tem mais: você pode definir que este prazo de validade seja por
                      <b>tempo indeterminado</b>
                      e o seu gasto irá se repetir para sempre.
                      <i>Ou, até quando você quiser.</i>
                    </p>
                    <p class="is-size-5 has-text-weight-normal">
                      Você pode estar se perguntando qual a diferença entre os tipos "fatura" e "poupança"... E Eu lhe respondo:
                      nenhuma. Por enquanto, nenhuma! Eles só estão ali para te ajudar na organização mesmo... Mas por enquanto, tanto
                      um gasto que possui o tipo "fatura" quanto um gasto que possui o tipo "poupança", ambos se comportarão da mesma
                      forma no sistema.
                    </p>
                  </li>
                </ul>
              </div>
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