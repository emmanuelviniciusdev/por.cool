<template>
  <div class="AutomaticWorkflowForm">
    <div class="box">
      <h2 class="subtitle is-5">novo registro</h2>

      <b-field
        label="imagem"
        :type="{ 'is-danger': imageError }"
        :message="{ 'selecione uma imagem': imageError }"
      >
        <b-upload
          v-model="imageFile"
          accept="image/*"
          @input="onImageSelected"
          drag-drop
          expanded
        >
          <section class="section">
            <div class="content has-text-centered">
              <p>
                <b-icon icon="upload" size="is-large"></b-icon>
              </p>
              <p>arraste a imagem aqui ou clique para selecionar</p>
            </div>
          </section>
        </b-upload>
      </b-field>

      <div v-if="imagePreview" class="image-preview">
        <img :src="imagePreview" alt="Preview" class="preview-image" />
        <b-button
          type="is-danger"
          size="is-small"
          @click="clearImage"
          class="clear-image-btn"
        >
          <b-icon icon="times" size="is-small"></b-icon>
        </b-button>
      </div>

      <b-field label="descrição (opcional)" class="description-field-wrapper">
        <div class="description-field">
          <b-autocomplete
            v-model="form.description"
            :data="filteredDescriptions"
            placeholder="digite ou selecione uma descrição"
            :open-on-focus="true"
            @select="option => (form.description = option ? option : '')"
            clearable
            expanded
          >
            <template slot="empty">nenhuma descrição encontrada</template>
          </b-autocomplete>
          <b-button
            type="is-success"
            size="is-small"
            @click="saveDescription"
            :disabled="!canSaveDescription"
            :loading="savingDescription"
            class="save-description-btn"
          >
            <b-icon icon="save" size="is-small"></b-icon>
          </b-button>
        </div>
      </b-field>

      <div class="saved-descriptions" v-if="preSavedDescriptions.length > 0">
        <p class="is-size-7 has-text-grey">descrições salvas:</p>
        <div class="tags">
          <span
            v-for="desc in preSavedDescriptions"
            :key="desc.id"
            class="tag is-light"
          >
            {{ desc.description }}
            <button
              class="delete is-small"
              @click="removeDescription(desc.id)"
            ></button>
          </span>
        </div>
      </div>

      <div class="form-actions">
        <b-button
          type="is-primary"
          @click="submit"
          :loading="loading"
          :disabled="!canSubmit"
          icon-left="plus"
        >
          adicionar
        </b-button>
      </div>
    </div>
  </div>
</template>

<script>
import { mapState } from "vuex";
import expenseAutomaticWorkflowService from "../services/expenseAutomaticWorkflow";
import preSavedDescriptionService from "../services/expenseAutomaticWorkflowPreSavedDescription";

export default {
  name: "AutomaticWorkflowForm",
  data() {
    return {
      imageFile: null,
      imagePreview: null,
      imageBase64: null,
      imageError: false,
      loading: false,
      savingDescription: false,
      form: {
        description: ""
      }
    };
  },
  computed: {
    ...mapState({
      userData: state => state.user.user,
      preSavedDescriptions: state =>
        state.expenseAutomaticWorkflow.preSavedDescriptions
    }),
    filteredDescriptions() {
      if (!this.form.description) {
        return this.preSavedDescriptions.map(d => d.description);
      }
      return this.preSavedDescriptions
        .map(d => d.description)
        .filter(desc =>
          desc.toLowerCase().includes(this.form.description.toLowerCase())
        );
    },
    canSubmit() {
      return this.imageBase64 !== null;
    },
    canSaveDescription() {
      const trimmedDesc = this.form.description.trim();
      if (!trimmedDesc) return false;
      const alreadyExists = this.preSavedDescriptions.some(
        d => d.description.toLowerCase() === trimmedDesc.toLowerCase()
      );
      return !alreadyExists;
    }
  },
  methods: {
    onImageSelected(file) {
      if (!file) {
        this.clearImage();
        return;
      }

      this.imageError = false;
      this.compressImage(file);
    },
    compressImage(file) {
      const maxWidth = 1200;
      const maxHeight = 1200;
      const quality = 0.7;

      const reader = new FileReader();
      reader.onload = e => {
        const img = new Image();
        img.onload = () => {
          const canvas = document.createElement("canvas");
          let width = img.width;
          let height = img.height;

          if (width > maxWidth) {
            height = (height * maxWidth) / width;
            width = maxWidth;
          }
          if (height > maxHeight) {
            width = (width * maxHeight) / height;
            height = maxHeight;
          }

          canvas.width = width;
          canvas.height = height;

          const ctx = canvas.getContext("2d");
          ctx.drawImage(img, 0, 0, width, height);

          const compressedBase64 = canvas.toDataURL("image/jpeg", quality);
          this.imageBase64 = compressedBase64;
          this.imagePreview = compressedBase64;
        };
        img.src = e.target.result;
      };
      reader.readAsDataURL(file);
    },
    clearImage() {
      this.imageFile = null;
      this.imagePreview = null;
      this.imageBase64 = null;
    },
    async saveDescription() {
      if (!this.canSaveDescription) return;

      this.savingDescription = true;
      try {
        await preSavedDescriptionService.insert({
          user: this.userData.uid,
          description: this.form.description.trim()
        });

        this.$store.dispatch(
          "expenseAutomaticWorkflow/setPreSavedDescriptions",
          {
            userUid: this.userData.uid
          }
        );

        this.$buefy.toast.open({
          message: "descrição salva com sucesso",
          type: "is-success",
          position: "is-bottom"
        });
      } catch (err) {
        console.error("Erro ao salvar descrição:", err);
        this.$buefy.toast.open({
          message: "erro ao salvar descrição",
          type: "is-danger",
          position: "is-bottom"
        });
      } finally {
        this.savingDescription = false;
      }
    },
    async removeDescription(descId) {
      try {
        await preSavedDescriptionService.remove(descId);

        this.$store.dispatch(
          "expenseAutomaticWorkflow/setPreSavedDescriptions",
          {
            userUid: this.userData.uid
          }
        );

        this.$buefy.toast.open({
          message: "descrição removida",
          type: "is-success",
          position: "is-bottom"
        });
      } catch (err) {
        this.$buefy.toast.open({
          message: "erro ao remover descrição",
          type: "is-danger",
          position: "is-bottom"
        });
      }
    },
    async submit() {
      if (!this.imageBase64) {
        this.imageError = true;
        return;
      }

      this.loading = true;

      try {
        const spendingDate = this.userData.lookingAtSpendingDate;

        await expenseAutomaticWorkflowService.insert({
          user: this.userData.uid,
          description: this.form.description.trim(),
          base64_image: this.imageBase64,
          spendingDate: spendingDate
        });

        this.$buefy.toast.open({
          message: "registro adicionado com sucesso",
          type: "is-success",
          position: "is-bottom"
        });

        this.clearImage();
        this.form.description = "";

        this.$emit("submitted");
      } catch (err) {
        console.error("Erro ao adicionar registro:", err);
        this.$buefy.toast.open({
          message: "erro ao adicionar registro",
          type: "is-danger",
          position: "is-bottom"
        });
      } finally {
        this.loading = false;
      }
    }
  }
};
</script>

<style lang="scss" scoped>
.AutomaticWorkflowForm {
  margin-bottom: 20px;
}

.image-preview {
  position: relative;
  display: inline-block;
  margin-bottom: 15px;

  .preview-image {
    max-width: 200px;
    max-height: 200px;
    border-radius: 4px;
    border: 1px solid #dbdbdb;
  }

  .clear-image-btn {
    position: absolute;
    top: -8px;
    right: -8px;
  }
}

.description-field-wrapper {
  width: 100%;
}

.description-field {
  display: flex;
  gap: 10px;
  align-items: center;
  width: 100%;

  .b-autocomplete,
  .autocomplete {
    flex: 1;
    width: 100%;
  }

  .save-description-btn {
    flex-shrink: 0;
  }
}

.saved-descriptions {
  margin-top: 10px;
  margin-bottom: 15px;

  .tags {
    margin-top: 5px;
  }
}

.form-actions {
  margin-top: 20px;
}
</style>
