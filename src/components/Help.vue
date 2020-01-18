<template>
  <div class="help">
    <b-tooltip :label="getTooltipText()" :position="this.tooltipPosition" type="is-warning">
      <button class="help-btn" @click="openModal()">
        <b-icon type="is-warning" icon="question-circle" size="is-medium"></b-icon>
      </button>
    </b-tooltip>

    <b-modal :active="showModal" has-modal-card aria-role="dialog" aria-modal :canCancel="false">
      <div class="modal-card">
        <header class="modal-card-head">
          <span class="modal-card-title is-size-6-mobile">
            <slot name="title"></slot>
          </span>
        </header>
        <section class="modal-card-body">
          <slot name="body"></slot>
        </section>
        <footer class="modal-card-foot">
          <b-button type="is-primary" @click="closeModal()">entendido!</b-button>
        </footer>
      </div>
    </b-modal>
  </div>
</template>

<script>
export default {
  name: "Help",
  props: {
    tooltipText: {
      type: String
    },
    tooltipPosition: {
      type: String,
      default: "is-top"
    }
  },
  data() {
    return {
      showModal: false
    };
  },
  methods: {
    openModal() {
      this.showModal = true;
    },
    closeModal() {
      this.showModal = false;
    },
    getTooltipText() {
      let tooltipText = this.tooltipText;
      return tooltipText !== undefined && tooltipText.trim() !== ""
        ? this.tooltipText
        : "se precisar de ajuda é só clicar aqui";
    }
  }
};
</script>

<style lang="scss" scoped>
.help {
  display: inline;

  .help-btn {
    cursor: pointer;
    background: none;
    border: none;
  }
}
</style>