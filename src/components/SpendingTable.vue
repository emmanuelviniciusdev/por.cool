<template>
  <div>
    <FilterByDate />

    <b-table :data="data.table" :mobile-cards="false" hoverable paginated :per-page="4">
      <template slot-scope="props">
        <b-table-column field="expense_name" label="#">{{ props.row.expense_name }}</b-table-column>
        <b-table-column field="amount" label="gasto">{{ props.row.amount }}</b-table-column>
        <b-table-column field="status" label="status">
          <b-tag :type="status_types[props.row.status]">
              <b-tooltip v-if="props.row.status === 'partially_paid'" label="R$ 50" type="is-dark">{{ status_labels[props.row.status] }}</b-tooltip>
              <span v-else>{{ status_labels[props.row.status] }}</span>
          </b-tag>
        </b-table-column>
        <b-table-column field="type" label="tipo">
          <b-tag size="is-medium">{{ types_labels[props.row.type] }}</b-tag>
        </b-table-column>
        <b-table-column field="action">
          <button class="button is-warning is-small btn-table-action">
            <b-icon icon="pencil-alt"></b-icon>
          </button>
          <button class="button is-danger is-small btn-table-action">
            <b-icon icon="trash"></b-icon>
          </button>
        </b-table-column>
      </template>
    </b-table>
  </div>
</template>

<script>
import FilterByDate from "./FilterByDate";

export default {
  name: "SpendingTable",
  components: {
    FilterByDate
  },
  data() {
    return {
      data: {
        table: [
          {
            expense_name: "Netflix",
            amount: 29.9,
            status: "pending",
            type: "invoice"
          },
          {
            expense_name: "RP no LoL",
            amount: 79.9,
            status: "paid",
            type: "expense"
          },
          {
            expense_name: "Passe Ônibus",
            amount: 150,
            status: "partially_paid",
            type: "expense",
            extra: { current_paid_amout: 50 }
          }
        ]
      },
      status_types: {
        paid: "is-success",
        partially_paid: "is-warning",
        pending: "is-danger"
      },
      status_labels: {
        paid: "pago",
        partially_paid: "parcialmente pago",
        pending: "pendente"
      },
      types_labels: {
        invoice: "fatura",
        expense: "gasto"
      }
    };
  }
};
</script>

<style lang="scss" scoped>
.btn-table-action {
  margin-right: 5px;
}
</style>