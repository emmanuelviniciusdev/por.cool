<template>
  <b-modal
    :active="isOpen"
    has-modal-card
    aria-role="dialog"
    aria-modal
    :canCancel="false"
  >
    <div class="modal-card" style="width: auto;">
      <header class="modal-card-head">
        <p class="modal-card-title is-size-6-mobile">visualizar imagem</p>
      </header>
      <section class="modal-card-body image-modal-body">
        <img
          v-if="base64Image"
          :src="base64Image"
          alt="Imagem do gasto"
          class="workflow-image"
        />
        <p v-else class="has-text-centered has-text-grey">
          nenhuma imagem disponível
        </p>
      </section>
      <footer class="modal-card-foot">
        <b-button type="is-primary" @click="close">fechar</b-button>
      </footer>
    </div>
  </b-modal>
</template>

<script>
export default {
  name: "AutomaticWorkflowImageModal",
  props: {
    base64Image: {
      type: String,
      default: ""
    },
    isOpen: {
      type: Boolean,
      default: false
    }
  },
  methods: {
    close() {
      this.$emit("close");
    }
  }
};
</script>

<style lang="scss" scoped>
.image-modal-body {
  display: block;
  overflow: auto;
  max-height: 70vh;
  text-align: center;
}

.workflow-image {
  max-width: 100%;
  height: auto;
}
</style>
