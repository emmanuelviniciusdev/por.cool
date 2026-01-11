<template>
  <div>
    <div class="columns is-desktop">
      <div class="column is-4-desktop">
        <LateralHeader />
      </div>
      <div class="column">
        <h1 class="title has-text-black">
          minha conta
          <Help>
            <template v-slot:title>minha conta</template>
            <template v-slot:body>
              <div class="content">
                <p class="is-size-5 has-text-weight-normal">
                  Por aqui, você terá algumas informações inerentes à sua conta
                  e poderá:
                </p>
                <ul>
                  <li>
                    <p class="is-size-6 has-text-weight-normal">
                      fazer alteração do e-mail e senha
                    </p>
                  </li>
                  <li>
                    <p class="is-size-6 has-text-weight-normal">
                      executar o reset financeiro (recomeçar do zero)
                    </p>
                  </li>
                  <li>
                    <p class="is-size-6 has-text-weight-normal">
                      deletar a sua conta
                    </p>
                  </li>
                </ul>
              </div>
            </template>
          </Help>
        </h1>

        <!-- Modal change email -->
        <b-modal :active="openModalChangeEmail" :canCancel="false">
          <div class="modal-card" style="width: auto;">
            <div class="modal-card-head">
              <h1 class="modal-card-title">alterar e-mail</h1>
            </div>

            <section class="modal-card-body">
              <b-field
                label="novo e-mail"
                :type="{
                  'is-danger': hasInputErrorAndDirty(
                    'formChangeEmail',
                    'newEmail'
                  )
                }"
                :message="{
                  'insira o seu e-mail': isInvalidInputMsg(
                    'formChangeEmail',
                    'newEmail',
                    'required'
                  ),
                  'insira um e-mail válido...': isInvalidInputMsg(
                    'formChangeEmail',
                    'newEmail',
                    'email'
                  )
                }"
              >
                <b-input
                  placeholder="seu@novo.email"
                  type="email"
                  v-model.trim="formChangeEmail.newEmail"
                  @change.native="
                    $v.formChangeEmail.newEmail.$model = $event.target.value
                  "
                ></b-input>
              </b-field>
              <b-field
                label="repita seu novo e-mail"
                :type="{
                  'is-danger': hasInputErrorAndDirty(
                    'formChangeEmail',
                    'confirmNewEmail'
                  )
                }"
                :message="{
                  'insira o seu e-mail': isInvalidInputMsg(
                    'formChangeEmail',
                    'confirmNewEmail',
                    'required'
                  ),
                  'insira um e-mail válido...': isInvalidInputMsg(
                    'formChangeEmail',
                    'confirmNewEmail',
                    'email'
                  ),
                  'e-mails não batem': isInvalidInputMsg(
                    'formChangeEmail',
                    'confirmNewEmail',
                    'sameAs'
                  )
                }"
              >
                <b-input
                  placeholder="seu@novo.email"
                  type="email"
                  v-model.trim="formChangeEmail.confirmNewEmail"
                  @change.native="
                    $v.formChangeEmail.confirmNewEmail.$model =
                      $event.target.value
                  "
                ></b-input>
              </b-field>
              <b-field
                label="confirmar senha"
                :type="{
                  'is-danger': hasInputErrorAndDirty(
                    'formChangeEmail',
                    'password'
                  )
                }"
                :message="{
                  'insira a sua senha': isInvalidInputMsg(
                    'formChangeEmail',
                    'password',
                    'required'
                  )
                }"
              >
                <b-input
                  placeholder="******"
                  type="password"
                  v-model.trim="formChangeEmail.password"
                  @change.native="
                    $v.formChangeEmail.password.$model = $event.target.value
                  "
                ></b-input>
              </b-field>
            </section>

            <footer class="modal-card-foot">
              <b-button @click="openModalChangeEmail = false"
                >cancelar</b-button
              >
              <b-button
                type="is-primary"
                @click="changeEmail()"
                :loading="loadingChangeEmail"
                >alterar</b-button
              >
            </footer>
          </div>
        </b-modal>

        <!-- Modal change password -->
        <b-modal :active="openModalChangePassword" :canCancel="false">
          <div class="modal-card" style="width: auto;">
            <div class="modal-card-head">
              <h1 class="modal-card-title">alterar senha</h1>
            </div>

            <section class="modal-card-body">
              <b-field
                label="senha atual"
                :type="{
                  'is-danger': hasInputErrorAndDirty(
                    'formChangePassword',
                    'password'
                  )
                }"
                :message="{
                  'insira uma senha segura': isInvalidInputMsg(
                    'formChangePassword',
                    'password',
                    'required'
                  )
                }"
              >
                <b-input
                  placeholder="******"
                  type="password"
                  v-model.trim="formChangePassword.password"
                  @change.native="
                    $v.formChangePassword.password.$model = $event.target.value
                  "
                ></b-input>
              </b-field>
              <b-field
                label="nova senha"
                :type="{
                  'is-danger': hasInputErrorAndDirty(
                    'formChangePassword',
                    'newPassword'
                  )
                }"
                :message="{
                  'insira uma senha segura': isInvalidInputMsg(
                    'formChangePassword',
                    'newPassword',
                    'required'
                  ),
                  'no mínimo, sua senha deve conter 6 caracteres': isInvalidInputMsg(
                    'formChangePassword',
                    'newPassword',
                    'minLength'
                  )
                }"
              >
                <b-input
                  placeholder="******"
                  type="password"
                  v-model.trim="formChangePassword.newPassword"
                  @change.native="
                    $v.formChangePassword.newPassword.$model =
                      $event.target.value
                  "
                ></b-input>
              </b-field>
              <b-field
                label="confirmar nova senha"
                :type="{
                  'is-danger': hasInputErrorAndDirty(
                    'formChangePassword',
                    'confirmNewPassword'
                  )
                }"
                :message="{
                  'insira uma senha segura': isInvalidInputMsg(
                    'formChangePassword',
                    'confirmNewPassword',
                    'required'
                  ),
                  'no mínimo, sua senha deve conter 6 caracteres': isInvalidInputMsg(
                    'formChangePassword',
                    'confirmNewPassword',
                    'minLength'
                  ),
                  'as duas senhas não batem': isInvalidInputMsg(
                    'formChangePassword',
                    'confirmNewPassword',
                    'sameAsPassword'
                  )
                }"
              >
                <b-input
                  placeholder="******"
                  type="password"
                  v-model.trim="formChangePassword.confirmNewPassword"
                  @change.native="
                    $v.formChangePassword.confirmNewPassword.$model =
                      $event.target.value
                  "
                ></b-input>
              </b-field>
            </section>

            <footer class="modal-card-foot">
              <b-button @click="openModalChangePassword = false"
                >cancelar</b-button
              >
              <b-button
                type="is-primary"
                @click="changePassword()"
                :loading="loadingChangePassword"
                >alterar</b-button
              >
            </footer>
          </div>
        </b-modal>

        <div class="notification user-info">
          <p class="is-size-5 has-text-weight-bold">
            {{ fullName | capitalizeName }}
          </p>
          <p class="is-size-5 has-text-weight-bold">
            {{ this.userData.email }}
          </p>
          <p class="is-size-6 has-text-weight-normal">
            {{ this.paymentRemainingDays }} dias restantes de utilização
          </p>
        </div>

        <div class="notification" v-if="userData.admin">
          <p class="is-size-5 has-text-weight-bold">Sincronizações</p>
          <table
            v-if="syncMetadata.length > 0"
            class="table is-fullwidth is-striped sync-table"
          >
            <thead>
              <tr>
                <th>Nome</th>
                <th>Última sincronização</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(sync, index) in syncMetadata" :key="index">
                <td>{{ sync.name }}</td>
                <td>{{ formatSyncDate(sync.latestSyncDatetime) }}</td>
              </tr>
            </tbody>
          </table>
          <p v-else class="is-size-6 has-text-weight-normal" style="margin-top: 10px;">
            Nenhuma sincronização realizada até o momento.
          </p>
        </div>

        <div class="notification">
          <p class="is-size-5 has-text-weight-bold">e-mail e senha</p>
          <b-button
            style="margin-top: 10px; margin-right: 5px;"
            @click="openModalChangeEmail = true"
            >alterar e-mail</b-button
          >
          <b-button
            style="margin-top: 10px;"
            @click="openModalChangePassword = true"
            >alterar senha</b-button
          >
        </div>

        <div class="notification is-warning">
          <p class="is-size-5 has-text-weight-bold">Recomeçar do zero</p>
          <p class="is-size-6 has-text-weight-normal">
            Se você ficou muito tempo sem utilizar o porcool ou simplesmente se
            perdeu em suas finanças, você pode executar um
            <i>reset financeiro</i> e recomeçar tudo do zero. Ao executar esta
            ação, <b>todos</b> os seus saldos e gastos serão excluídos
            <b>para sempre</b>.
          </p>
          <b-button
            style="margin-top: 10px;"
            @click="startOver()"
            :loading="loadingStartOver"
            >recomeçar do zero</b-button
          >
        </div>

        <div class="notification is-danger">
          <p class="is-size-5 has-text-weight-bold">Deletar minha conta</p>
          <p class="is-size-6 has-text-weight-normal">
            Deleta todos os seus dados e vínculos com o porcool.
          </p>
          <p class="is-size-6 has-text-weight-bold">
            Esta ação não tem mais volta.
          </p>
          <b-button
            style="margin-top: 10px;"
            @click="deleteAccount()"
            :loading="loadingDeleteAccount"
            >deletar</b-button
          >
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { mapState } from "vuex";
import {
  required,
  email,
  minLength,
  maxLength,
  sameAs
} from "vuelidate/lib/validators";

// Components
import LateralHeader from "../components/LateralHeader";
import Help from "../components/Help";

// Filters
import filters from "../filters";

// Services
import paymentService from "../services/payment";
import authService from "../services/auth";
import userService from "../services/user";
import settingsService from "../services/settings";

// Helpers
import dateAndTimeHelper from "../helpers/dateAndTime";

export default {
  name: "MyAccount",
  components: {
    LateralHeader,
    Help
  },
  filters: {
    capitalizeName: filters.capitalizeName
  },
  data() {
    return {
      paymentRemainingDays: 0,
      syncMetadata: [],

      openModalChangePassword: false,
      loadingChangePassword: false,
      formChangePassword: {
        password: "",
        newPassword: "",
        confirmNewPassword: ""
      },

      openModalChangeEmail: false,
      loadingChangeEmail: false,
      formChangeEmail: {
        newEmail: "",
        confirmNewEmail: "",
        password: ""
      },

      loadingStartOver: false,
      loadingDeleteAccount: false
    };
  },
  validations: {
    formChangePassword: {
      password: { required },
      newPassword: { required, minLength: minLength(6) },
      confirmNewPassword: {
        required,
        minLength: minLength(6),
        sameAsPassword: sameAs("newPassword")
      }
    },

    formChangeEmail: {
      newEmail: { required, email },
      confirmNewEmail: { required, email, sameAs: sameAs("newEmail") },
      password: { required }
    }
  },
  computed: {
    ...mapState({
      userData: state => state.user.user
    }),
    fullName() {
      return this.userData.name + " " + this.userData.lastName;
    }
  },
  methods: {
    async changePassword() {
      if (this.$v.formChangePassword.$invalid) return;

      this.loadingChangePassword = true;

      try {
        const { password, newPassword } = this.formChangePassword;
        const changePassword = await authService.changePassword(
          password,
          newPassword
        );

        this.$buefy.toast.open({
          message: changePassword.message,
          type: changePassword.error ? "is-danger" : "is-success",
          position: "is-bottom"
        });

        if (!changePassword.error) {
          this.$v.formChangePassword.$reset();
          this.formChangePassword.password = "";
          this.formChangePassword.newPassword = "";
          this.formChangePassword.confirmNewPassword = "";

          this.openModalChangePassword = false;
        }
      } catch (err) {
        this.$buefy.toast.open({
          message: "ocorreu um erro ao alterar a sua senha [2]",
          type: "is-danger",
          position: "is-bottom"
        });
      } finally {
        this.loadingChangePassword = false;
      }
    },
    async changeEmail() {
      if (this.$v.formChangeEmail.$invalid) return;

      this.loadingChangeEmail = true;

      try {
        const { newEmail, password } = this.formChangeEmail;
        const changeEmail = await authService.changeEmail(password, newEmail);

        this.$buefy.toast.open({
          message: changeEmail.message,
          type: changeEmail.error ? "is-danger" : "is-success",
          position: "is-bottom"
        });

        if (!changeEmail.error) {
          this.$store.dispatch("user/update", {
            email: this.formChangeEmail.newEmail
          });

          this.$v.formChangeEmail.$reset();
          this.formChangeEmail.newEmail = "";
          this.formChangeEmail.confirmNewEmail = "";
          this.formChangeEmail.password = "";

          this.openModalChangeEmail = false;
        }
      } catch (err) {
        this.$buefy.toast.open({
          message: "ocorreu um erro ao alterar o seu e-mail [2]",
          type: "is-danger",
          position: "is-bottom"
        });
      } finally {
        this.loadingChangeEmail = false;
      }
    },
    async startOver() {
      this.$buefy.dialog.prompt({
        title: "confirmar senha",
        message: "Por favor, confirme sua senha",
        inputAttrs: {
          placeholder: "******",
          type: "password"
        },
        type: "is-warning",
        cancelText: "cancelar",
        confirmText: "recomeçar do zero",
        canCancel: ["button", "escape"],
        onConfirm: async password => {
          this.loadingStartOver = true;

          try {
            // Authenticate password
            try {
              await authService.reauthenticate(password);
            } catch (err) {
              let message = "a senha está incorreta";

              if (err.code === "auth/too-many-requests")
                message =
                  "Você excedeu o limite de tentativas. Por favor, tente novamente mais tarde.";

              this.$buefy.toast.open({
                message,
                type: "is-danger",
                position: "is-bottom"
              });

              return false;
            }

            await userService.startOver(this.userData.uid);

            this.$store.dispatch("user/update", {
              lookingAtSpendingDate: dateAndTimeHelper.startOfMonthAndDay(
                new Date()
              )
            });
            this.$store.dispatch("balances/setBalances", {
              userUid: this.userData.uid,
              spendingDate: this.userData.lookingAtSpendingDate
            });

            this.$buefy.toast.open({
              message: "reset realizado com sucesso! agora sim, sem bagunça...",
              type: "is-success",
              position: "is-bottom"
            });
          } catch {
            this.$buefy.toast.open({
              message: "ocorreu um erro ao resetar suas finanças",
              type: "is-danger",
              position: "is-bottom"
            });
          } finally {
            this.loadingStartOver = false;
          }
        }
      });
    },
    deleteAccount() {
      this.$buefy.dialog.prompt({
        title: "confirmar senha",
        message: "Por favor, confirme sua senha",
        inputAttrs: {
          placeholder: "******",
          type: "password"
        },
        type: "is-danger",
        cancelText: "cancelar",
        confirmText: "deletar conta",
        canCancel: ["button", "escape"],
        onConfirm: async password => {
          this.loadingDeleteAccount = true;

          try {
            // Authenticate password
            try {
              await authService.reauthenticate(password);
            } catch (err) {
              let message = "a senha está incorreta";

              if (err.code === "auth/too-many-requests")
                message =
                  "Você excedeu o limite de tentativas. Por favor, tente novamente mais tarde.";

              this.$buefy.toast.open({
                message,
                type: "is-danger",
                position: "is-bottom"
              });

              return false;
            }

            await authService.deleteAccount(password);

            this.$router.push({ name: "goodbye" });
          } catch {
            this.$buefy.toast.open({
              message: "ocorreu um erro ao deletar sua conta",
              type: "is-danger",
              position: "is-bottom"
            });
          } finally {
            this.loadingDeleteAccount = false;
          }
        }
      });
    },
    hasInputErrorAndDirty(form, input) {
      return this.$v[form][input].$error && this.$v[form][input].$dirty;
    },
    isInvalidInputMsg(form, input, role) {
      return !this.$v[form][input][role] && this.$v[form][input].$error;
    },
    formatSyncDate(date) {
      if (!date) return "-";
      const d = date.toDate ? date.toDate() : new Date(date);
      return d.toLocaleString("pt-BR", {
        day: "2-digit",
        month: "2-digit",
        year: "numeric",
        hour: "2-digit",
        minute: "2-digit"
      });
    }
  },
  async created() {
    const lastPayment = await paymentService.lastPaymentInfo(this.userData.uid);
    this.paymentRemainingDays = lastPayment.remainingDays;

    if (this.userData.admin) {
      this.syncMetadata = await settingsService.getSyncMetadata();
    }
  }
};
</script>

<style lang="scss" scoped>
.user-info {
  p:last-child {
    margin-top: 5px;
  }
}

.sync-table {
  margin-top: 10px;
  background: transparent;
}

@media screen and (min-width: 769px) {
  .modal-card {
    width: 500px !important;
    margin: 0 auto !important;
  }
}

@media screen and (min-width: 1024px) {
  .columns {
    margin-top: 30px;
  }
}
</style>
