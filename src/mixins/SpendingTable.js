export default {
  data() {
    return {
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
        expense: "gasto",
        savings: "poupança"
      }
    };
  }
};
